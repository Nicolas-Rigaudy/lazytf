package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	addStyle      = lipgloss.NewStyle().Foreground(theme.Current.Green).Bold(true)
	deleteStyle   = lipgloss.NewStyle().Foreground(theme.Current.Red).Bold(true)
	changeStyle   = lipgloss.NewStyle().Foreground(theme.Current.Yellow).Bold(true)
	selectedStyle = lipgloss.NewStyle().Background(theme.Current.Surface1)
)

type ModuleGroup struct {
	Name        string
	Resources   []terraform.ResourceChange
	AddCount    int
	ChangeCount int
	DeleteCount int
	Expanded    bool
}

type PlanSummaryViewModel struct {
	ModuleGroups          []ModuleGroup
	SelectedModuleIndex   int
	SelectedResourceIndex int
	ResourceExpandedMap   map[string]bool
	AddCount              int
	ChangeCount           int
	DeleteCount           int
	viewport              viewport.Model
}

func NewPlanSummaryViewModel(planSummary *terraform.PlanSummary) *PlanSummaryViewModel {
	vp := viewport.New(0, 0)
	vp.SetHorizontalStep(6)

	moduleGroups := buildModuleGroups(planSummary)

	vm := &PlanSummaryViewModel{
		ModuleGroups:          moduleGroups,
		SelectedModuleIndex:   0,
		SelectedResourceIndex: -1,
		ResourceExpandedMap:   make(map[string]bool),
		viewport:              vp,
	}

	for _, group := range moduleGroups {
		vm.AddCount += group.AddCount
		vm.ChangeCount += group.ChangeCount
		vm.DeleteCount += group.DeleteCount
	}

	vm.updateContent()
	return vm
}

func (vm *PlanSummaryViewModel) SetSize(width, height int) {
	vm.viewport.Width = width
	vm.viewport.Height = height
	vm.updateContent()
}

func (vm *PlanSummaryViewModel) updateContent() {
	var lines []string

	for moduleIdx, group := range vm.ModuleGroups {
		// Module header line
		arrow := "▶"
		if group.Expanded {
			arrow = "▼"
		}

		// Format counts with colors: (2+, 3~, 1-)
		addCountStr := addStyle.Render(fmt.Sprintf("%d+", group.AddCount))
		changeCountStr := changeStyle.Render(fmt.Sprintf("%d~", group.ChangeCount))
		deleteCountStr := deleteStyle.Render(fmt.Sprintf("%d-", group.DeleteCount))
		counts := fmt.Sprintf("(%s, %s, %s)", addCountStr, changeCountStr, deleteCountStr)

		moduleLine := fmt.Sprintf("%s %s %s", arrow, group.Name, counts)

		// Highlight if module header is selected
		if moduleIdx == vm.SelectedModuleIndex && vm.SelectedResourceIndex == -1 {
			moduleLine = selectedStyle.Render(moduleLine)
		}

		lines = append(lines, moduleLine)

		// If module is expanded, show resources
		if group.Expanded {
			for resourceIdx, rc := range group.Resources {
				resArrow := "  ▶"
				if vm.ResourceExpandedMap[rc.Address] {
					resArrow = "  ▼"
				}

				// Determine action prefix and style
				prefix := "~"
				style := changeStyle
				if len(rc.Change.Actions) > 0 {
					switch rc.Change.Actions[0] {
					case "create":
						prefix = "+"
						style = addStyle
					case "delete":
						prefix = "-"
						style = deleteStyle
					}
				}

				// Show short name (strip module prefix for non-root modules)
				shortName := rc.Address
				if group.Name != "(root)" {
					shortName = rc.Address[len(group.Name)+1:] // +1 for the dot
				}

				resourceLine := resArrow + " " + style.Render(prefix+" "+shortName)

				// Highlight if this resource is selected
				if moduleIdx == vm.SelectedModuleIndex && resourceIdx == vm.SelectedResourceIndex {
					resourceLine = selectedStyle.Render(resourceLine)
				}

				lines = append(lines, resourceLine)

				// If resource is expanded, show before/after
				if vm.ResourceExpandedMap[rc.Address] {
					beforeJSON, _ := json.MarshalIndent(rc.Change.Before, "      ", "  ")
					afterJSON, _ := json.MarshalIndent(rc.Change.After, "      ", "  ")
					lines = append(lines, "      Before: "+string(beforeJSON))
					lines = append(lines, "      After:  "+string(afterJSON))
				}
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	vm.viewport.SetContent(content)
}

func (vm *PlanSummaryViewModel) Title() string {
	return fmt.Sprintf("Plan Summary - %d to add, %d to change, %d to destroy",
		vm.AddCount, vm.ChangeCount, vm.DeleteCount)
}

func (vm *PlanSummaryViewModel) SelectedLineNumber() int {
	lineNumber := 0

	for moduleIdx, group := range vm.ModuleGroups {
		if moduleIdx == vm.SelectedModuleIndex && vm.SelectedResourceIndex == -1 {
			return lineNumber
		}
		lineNumber++ // module header

		if group.Expanded {
			for resourceIdx, rc := range group.Resources {
				if moduleIdx == vm.SelectedModuleIndex && resourceIdx == vm.SelectedResourceIndex {
					return lineNumber
				}
				lineNumber++ // resource line
				if vm.ResourceExpandedMap[rc.Address] {
					lineNumber += 2 // before/after lines
				}
			}
		}
	}
	return lineNumber
}

func (vm *PlanSummaryViewModel) ensureSelectedVisible() {
	selectedLine := vm.SelectedLineNumber()
	yOffset := vm.viewport.YOffset
	height := vm.viewport.Height

	// If selected line is above visible area, scroll up
	if selectedLine < yOffset {
		vm.viewport.SetYOffset(selectedLine)
	}
	// If selected line is below visible area, scroll down
	if selectedLine >= yOffset+height {
		vm.viewport.SetYOffset(selectedLine - height + 1)
	}
}

func extractModuleName(address string) string {
	parts := strings.Split(address, ".")
	if len(parts) <= 2 {
		return "(root)"
	}
	return strings.Join(parts[:len(parts)-2], ".")
}

func buildModuleGroups(planSummary *terraform.PlanSummary) []ModuleGroup {
	moduleMap := make(map[string]*ModuleGroup)

	for _, resourceChange := range planSummary.ResourceChanges {
		if len(resourceChange.Change.Actions) == 0 {
			continue
		}
		// Skip no-op resources
		if resourceChange.Change.Actions[0] == "no-op" {
			continue
		}
		moduleName := extractModuleName(resourceChange.Address)
		group, exists := moduleMap[moduleName]
		if !exists {
			group = &ModuleGroup{Name: moduleName}
			moduleMap[moduleName] = group
		}
		group.Resources = append(group.Resources, resourceChange)
		switch resourceChange.Change.Actions[0] {
		case "create":
			group.AddCount++
		case "delete":
			group.DeleteCount++
			// Check for replace (delete then create)
			if len(resourceChange.Change.Actions) > 1 && resourceChange.Change.Actions[1] == "create" {
				group.AddCount++
			}
		case "update":
			group.ChangeCount++
		}
	}

	var moduleGroups []ModuleGroup
	for _, group := range moduleMap {
		moduleGroups = append(moduleGroups, *group)
	}

	return moduleGroups
}

func (vm *PlanSummaryViewModel) View() string {
	return vm.viewport.View()
}

func (vm *PlanSummaryViewModel) moveNext() {
	group := &vm.ModuleGroups[vm.SelectedModuleIndex]

	if vm.SelectedResourceIndex == -1 {
		// On module header
		if group.Expanded && len(group.Resources) > 0 {
			// Go to first resource
			vm.SelectedResourceIndex = 0
		} else if vm.SelectedModuleIndex < len(vm.ModuleGroups)-1 {
			// Go to next module
			vm.SelectedModuleIndex++
		}
	} else if vm.SelectedResourceIndex < len(group.Resources)-1 {
		// Go to next resource
		vm.SelectedResourceIndex++
	} else if vm.SelectedModuleIndex < len(vm.ModuleGroups)-1 {
		// Go to next module
		vm.SelectedModuleIndex++
		vm.SelectedResourceIndex = -1
	}
}

func (vm *PlanSummaryViewModel) movePrev() {
	if vm.SelectedResourceIndex == -1 {
		// On module header
		if vm.SelectedModuleIndex > 0 {
			vm.SelectedModuleIndex--
			prevGroup := &vm.ModuleGroups[vm.SelectedModuleIndex]
			if prevGroup.Expanded && len(prevGroup.Resources) > 0 {
				// Go to last resource of previous module
				vm.SelectedResourceIndex = len(prevGroup.Resources) - 1
			}
		}
	} else if vm.SelectedResourceIndex > 0 {
		// Go to previous resource
		vm.SelectedResourceIndex--
	} else {
		// Go to module header
		vm.SelectedResourceIndex = -1
	}
}

func (vm *PlanSummaryViewModel) moveNextModule() {
	if vm.SelectedModuleIndex < len(vm.ModuleGroups)-1 {
		vm.SelectedModuleIndex++
		vm.SelectedResourceIndex = -1
	}
}

func (vm *PlanSummaryViewModel) movePrevModule() {
	if vm.SelectedModuleIndex > 0 {
		vm.SelectedModuleIndex--
		vm.SelectedResourceIndex = -1
	}
}

func (vm *PlanSummaryViewModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+k", "ctrl+up":
			// Move to previous item (module or resource)
			vm.movePrev()
			vm.updateContent()
			vm.ensureSelectedVisible()
			return nil
		case "ctrl+j", "ctrl+down":
			// Move to next item (module or resource)
			vm.moveNext()
			vm.updateContent()
			vm.ensureSelectedVisible()
			return nil
		case "alt+k", "alt+up":
			// Jump to previous module
			vm.movePrevModule()
			vm.updateContent()
			vm.ensureSelectedVisible()
			return nil
		case "alt+j", "alt+down":
			// Jump to next module
			vm.moveNextModule()
			vm.updateContent()
			vm.ensureSelectedVisible()
			return nil
		case "enter", " ":
			if vm.SelectedResourceIndex == -1 {
				// Toggle module expansion
				vm.ModuleGroups[vm.SelectedModuleIndex].Expanded = !vm.ModuleGroups[vm.SelectedModuleIndex].Expanded
			} else {
				// Toggle resource expansion
				rc := vm.ModuleGroups[vm.SelectedModuleIndex].Resources[vm.SelectedResourceIndex]
				vm.ResourceExpandedMap[rc.Address] = !vm.ResourceExpandedMap[rc.Address]
			}
			vm.updateContent()
			vm.ensureSelectedVisible()
			return nil
		}
	}

	// Delegate all other keys to viewport (j/k, arrows, pgup/pgdn, etc)
	vm.viewport, cmd = vm.viewport.Update(msg)
	return cmd
}
