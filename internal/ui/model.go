package ui

import (
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode int

const (
	ViewModeProjectList ViewMode = iota
	ViewModeProjectDetail
)

type Model struct {
	sidebar             SidebarModel
	mainPanel           MainPanelModel
	titleBar            TitleBarModel
	statusBar           StatusBarModel
	focusIndex          int
	focusableCount      int
	width               int
	height              int
	projects            []terraform.Project
	mode                terraform.Mode
	viewMode            ViewMode
	selectedProject     *terraform.Project
	varFiles            []terraform.VarFile
	selectedVarFile     *terraform.VarFile
	backendVarFiles     []terraform.BackendVarFile
	selectedBackendFile *terraform.BackendVarFile
	backendState        terraform.BackendState
}

func NewModel(projects []terraform.Project, mode terraform.Mode) Model {
	var sidebar SidebarModel
	var viewMode ViewMode
	var selectedProject *terraform.Project
	var varFiles []terraform.VarFile
	var backendVarFiles []terraform.BackendVarFile

	var backendState terraform.BackendState

	if mode == terraform.ModeSingleProject {
		// Single-project mode setup
		viewMode = ViewModeProjectDetail
		selectedProject = &projects[0]

		// Discover var files and backends for the project
		varFiles, _ = terraform.DiscoverVarFiles(selectedProject.Path)
		sidebarItems := terraform.GetVarFileDisplayNames(varFiles)
		backendVarFiles, _ = terraform.DiscoverBackendVarFiles(selectedProject.Path)

		// Detect current backend initialization state
		backendState = terraform.DetectCurrentBackend(selectedProject.Path, backendVarFiles)

		sidebar = NewSidebar(sidebarItems...)
		sidebar.Title = selectedProject.Name

	} else {
		// Multi-project mode setup
		viewMode = ViewModeProjectList
		selectedProject = nil
		varFiles = nil
		backendVarFiles = nil

		// Build sidebar with project names
		sidebarItems := make([]string, len(projects))
		for i, project := range projects {
			sidebarItems[i] = project.Name
		}
		sidebar = NewSidebar(sidebarItems...)
	}

	sidebar.IsFocused = true
	mainPanel := NewMainPanel()
	titleBar := NewTitleBar()
	statusBar := NewStatusBar()

	m := Model{
		sidebar:             sidebar,
		mainPanel:           mainPanel,
		titleBar:            titleBar,
		statusBar:           statusBar,
		focusIndex:          0,
		focusableCount:      2,
		width:               0,
		height:              0,
		projects:            projects,
		mode:                mode,
		viewMode:            viewMode,
		selectedProject:     selectedProject,
		varFiles:            varFiles,
		selectedVarFile:     nil,
		backendVarFiles:     backendVarFiles,
		selectedBackendFile: nil,
		backendState:        backendState,
	}

	// Set initial status bar text
	m.statusBar.SetText(m.buildStatusText())

	return m
}

func (m *Model) updateFocusStates() {
	m.sidebar.IsFocused = (m.focusIndex == 0)
	m.mainPanel.IsFocused = (m.focusIndex == 1)
}

// buildStatusText creates dynamic status bar text based on current state
func (m Model) buildStatusText() string {
	var parts []string

	// Add project info if selected
	if m.selectedProject != nil {
		parts = append(parts, "Project: "+m.selectedProject.Name)
	}

	// Add environment info if var file selected
	if m.selectedVarFile != nil {
		parts = append(parts, "Env: "+m.selectedVarFile.EnvName)
	}

	// Add separator if we have context
	if len(parts) > 0 {
		parts = append(parts, "│")
	}

	// Add help keys based on view mode
	if m.viewMode == ViewModeProjectDetail && m.mode == terraform.ModeMultiProject {
		parts = append(parts, "Backspace: back")
		parts = append(parts, "│")
	}

	parts = append(parts, "Tab: switch", "↑↓/jk: navigate", "Enter: select", "q: quit")

	return strings.Join(parts, " ")
}

func (m Model) View() string {
	title := m.titleBar.View()

	content := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), m.mainPanel.View())

	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, title, content, status)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % m.focusableCount
			m.updateFocusStates()
		case "backspace", "esc", "p":
			// Only go back if we're in ViewModeProjectDetail
			if m.viewMode == ViewModeProjectDetail {
				m.viewMode = ViewModeProjectList
				m.selectedProject = nil
				m.varFiles = nil
				m.selectedVarFile = nil
				sidebarItems := make([]string, len(m.projects))
				for i, project := range m.projects {
					sidebarItems[i] = project.Name
				}
				m.sidebar.Items = sidebarItems
				m.sidebar.Title = "Projects"
				m.sidebar.SelectedIndex = 0
				m.statusBar.SetText(m.buildStatusText())
				return m, nil
			}
		case "enter":
			// Handle enter based on focus and view mode
			if m.focusIndex == 0 { // Sidebar is focused
				if m.viewMode == ViewModeProjectList {
					// Select project
					if m.sidebar.SelectedIndex < len(m.projects) {
						return m, func() tea.Msg {
							return ProjectSelectedMsg{
								ProjectName: m.sidebar.Items[m.sidebar.SelectedIndex],
								Index:       m.sidebar.SelectedIndex,
							}
						}
					}
				} else if m.viewMode == ViewModeProjectDetail {
					// Select var file
					if m.sidebar.SelectedIndex < len(m.varFiles) {
						return m, func() tea.Msg {
							return VarFileSelectedMsg{
								VarFileName: m.sidebar.Items[m.sidebar.SelectedIndex],
								Index:       m.sidebar.SelectedIndex,
							}
						}
					}
				}
			}
		}

		var cmd tea.Cmd
		// Delegate to focused component if not global key
		switch m.focusIndex {
		case 0:
			m.sidebar, cmd = m.sidebar.Update(msg)
		}
		return m, cmd

	case ProjectSelectedMsg:
		// Get the selected project
		selectedProject := &m.projects[msg.Index]
		m.selectedProject = selectedProject

		// Discover var files and backends for the selected project
		m.varFiles, _ = terraform.DiscoverVarFiles(selectedProject.Path)
		sidebarItems := terraform.GetVarFileDisplayNames(m.varFiles)
		m.backendVarFiles, _ = terraform.DiscoverBackendVarFiles(selectedProject.Path)

		// Detect current backend initialization state
		m.backendState = terraform.DetectCurrentBackend(selectedProject.Path, m.backendVarFiles)

		// Update sidebar
		m.sidebar.Items = sidebarItems
		m.sidebar.Title = selectedProject.Name
		m.sidebar.SelectedIndex = 0 // Reset to first item

		// Switch to detail view
		m.viewMode = ViewModeProjectDetail

		// Update main panel with title and content
		m.mainPanel.Title = "📋 Project Details"
		backendStateInfo := terraform.FormatBackendState(m.backendState)
		m.mainPanel.Content = "Project: " + selectedProject.Name + "\n" +
			"Path: " + selectedProject.Path + "\n\n" +
			"--- Backend Status ---\n" +
			backendStateInfo + "\n" +
			"--- Available Environments ---\n" +
			"Var Files: " + string(rune(len(m.varFiles)+'0')) + "\n" +
			"Backend Configs: " + string(rune(len(m.backendVarFiles)+'0'))

		// Update status bar
		m.statusBar.SetText(m.buildStatusText())

		return m, nil

	case VarFileSelectedMsg:
		// Get the selected var file
		selectedVarFile := &m.varFiles[msg.Index]
		m.selectedVarFile = selectedVarFile

		// Find matching backend configs (returns a slice, may be 0, 1, or multiple)
		matchedBackends := terraform.MatchBackendToVarFile(*selectedVarFile, m.backendVarFiles)
		backendInfo := terraform.FormatBackendInfo(matchedBackends)

		// Check if this environment is currently initialized
		var currentStatus string
		if m.backendState.IsInitialized && m.backendState.DetectedEnv == selectedVarFile.EnvName {
			currentStatus = "✅ This environment is currently initialized"
		} else if m.backendState.IsInitialized {
			currentStatus = "⚠️  Different environment is initialized (" + m.backendState.DetectedEnv + ")"
		} else {
			currentStatus = "❌ Not initialized"
		}

		// Update main panel with title and content
		m.mainPanel.Title = "🌍 Environment Details"
		m.mainPanel.Content = "Environment: " + selectedVarFile.EnvName + "\n" +
			"Status: " + currentStatus + "\n\n" +
			"Var File: " + selectedVarFile.Name + "\n" +
			"Full Path: " + selectedVarFile.FullPath + "\n\n" +
			"Backend Configuration:\n" + backendInfo

		// Update status bar
		m.statusBar.SetText(m.buildStatusText())

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		titleBarHeight := 1
		statusBarHeight := 1
		contentHeight := m.height - titleBarHeight - statusBarHeight - 2

		sidebarWidth := m.width / 4
		mainPanelWidth := m.width - sidebarWidth - 4

		m.titleBar.Width = m.width
		m.statusBar.Width = m.width
		m.sidebar.Width = sidebarWidth
		m.sidebar.Height = contentHeight
		m.mainPanel.Width = mainPanelWidth
		m.mainPanel.Height = contentHeight
	}
	return m, nil
}
