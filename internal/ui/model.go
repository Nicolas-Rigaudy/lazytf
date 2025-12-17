package ui

import (
	"path/filepath"
	"strings"
	"time"

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
	header              HeaderModel
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
	modal               Modal // Modal component
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
		sidebar.InitializedEnv = backendState.DetectedEnv

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
	header := NewHeader()
	modal := NewModal()

	m := Model{
		sidebar:             sidebar,
		mainPanel:           mainPanel,
		titleBar:            titleBar,
		statusBar:           statusBar,
		header:              header,
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
		modal:               modal,
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

	// Add help keys based on view mode
	if m.viewMode == ViewModeProjectDetail {
		if m.mode == terraform.ModeMultiProject {
			parts = append(parts, "Backspace: back", "│")
		}
		parts = append(parts, "i: init", "│")
	}

	parts = append(parts, "Tab: switch", "↑↓/jk: navigate", "Enter: select", "q: quit")

	return strings.Join(parts, " ")
}

func (m Model) View() string {
	// Build base UI
	title := m.titleBar.View()
	header := m.header.View(InfoHeaderData{
		ProjectName: func() string {
			if m.selectedProject != nil {
				return m.selectedProject.Name
			} else {
				return "No Project"
			}
		}(),
		EnvName: func() string {
			if m.selectedVarFile != nil {
				return m.selectedVarFile.EnvName
			} else {
				return "No Env"
			}
		}(),
		IsInitialized:   m.backendState.IsInitialized && m.selectedVarFile != nil && m.backendState.DetectedEnv == m.selectedVarFile.EnvName,
		LastCommand:     "",
		LastCommandTime: time.Time{},
	})
	content := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), m.mainPanel.View())
	status := m.statusBar.View()
	baseUI := lipgloss.JoinVertical(lipgloss.Left, title, header, content, status)

	// If modal is active, render it (modal handles its own rendering)
	if m.modal.IsActive() {
		return m.modal.View(m.width, m.height)
	}

	return baseUI
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If modal is active, it gets priority for input handling
	if m.modal.IsActive() {
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // ctrl+c always quits (emergency exit)
			return m, tea.Quit
		case "q": // 'q' only quits when modal is NOT active
			if !m.modal.IsActive() {
				return m, tea.Quit
			}
		case "t": // TEST: Show test confirm modal (remove later)
			m.modal.Show(ModalState{
				Type:    ModalConfirm,
				Title:   "Test Confirmation Modal",
				Message: "This is a test modal!\n\nBackend: variables/backend/local/backend_dev2.tfvars\nVar-file: variables/dev2.tfvars\n\nDoes it look good?",
				OnConfirm: func() tea.Msg {
					m.statusBar.SetText("✅ Modal confirmed!")
					return nil
				},
			})
			return m, nil
		case "s": // TEST: Show test select modal (remove later)
			items := []string{"dev1", "dev2", "int", "prod"}
			m.modal.Show(ModalState{
				Type:     ModalSelect,
				Title:    "Select Environment",
				Items:    items,
				Selected: 0,
				OnSelect: func(selectedIndex int) tea.Msg {
					m.statusBar.SetText("✅ Selected: " + items[selectedIndex])
					return nil
				},
			})
			return m, nil

		case "i":
			if m.viewMode == ViewModeProjectDetail && m.selectedProject != nil {
				envNames := terraform.GetVarFileDisplayNames(m.varFiles)

				m.modal.Show(ModalState{
					Type:    ModalSelect,
					Title:   "Terraform Init",
					Message: "Choose an environment to init for " + m.selectedProject.Name,
					Items:   envNames,
					OnSelect: func(index int) tea.Msg {
						selectedEnv := envNames[index]

						backends := terraform.MatchBackendsForEnv(selectedEnv, m.backendVarFiles)

						return InitEnvironmentSelectedMsg{
							EnvName:  selectedEnv,
							Backends: backends,
						}
					},
				})
				return m, nil
			}
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
		m.sidebar.InitializedEnv = m.backendState.DetectedEnv

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

		// Find matching backend configs
		matchedBackends := terraform.MatchBackendsForEnv(selectedVarFile.EnvName, m.backendVarFiles)
		backendInfo := terraform.FormatBackendInfo(matchedBackends)

		// Determine initialization status for this environment
		isThisEnvInitialized := m.backendState.IsInitialized && m.backendState.DetectedEnv == selectedVarFile.EnvName

		// Build status string
		var currentStatus string
		if isThisEnvInitialized {
			currentStatus = "✅ This environment is currently initialized"
		} else if m.backendState.IsInitialized {
			currentStatus = "⚠️  Different environment is initialized (" + m.backendState.DetectedEnv + ")"
		} else {
			currentStatus = "❌ Not initialized"
		}

		// If not initialized, show modal asking if they want to init
		if !isThisEnvInitialized {
			m.modal.Show(ModalState{
				Type:  ModalConfirm,
				Title: "⚠️  Environment Not Initialized",
				Message: "Environment \"" + selectedVarFile.EnvName + "\" is not initialized.\n\n" +
					"Terraform commands won't work until you initialize it.\n\n" +
					"[Enter] Initialize now    [Esc] View details anyway",
				OnConfirm: func() tea.Msg {
					// Start init workflow
					backends := terraform.MatchBackendsForEnv(selectedVarFile.EnvName, m.backendVarFiles)
					return InitEnvironmentSelectedMsg{
						EnvName:  selectedVarFile.EnvName,
						Backends: backends,
					}
				},
				OnCancel: func() tea.Msg {
					// Show details anyway - trigger display update
					return nil
				},
			})
		}

		m.mainPanel.Title = "🌍 Environment Details"
		m.mainPanel.Content = "Environment: " + selectedVarFile.EnvName + "\n" +
			"Status: " + currentStatus + "\n\n" +
			"Var File: " + selectedVarFile.Name + "\n" +
			"Full Path: " + selectedVarFile.FullPath + "\n\n" +
			"Backend Configuration:\n" + backendInfo

		// Update status bar
		m.statusBar.SetText(m.buildStatusText())

		return m, nil

	case InitEnvironmentSelectedMsg:
		m.selectedVarFile, m.sidebar.SelectedIndex = terraform.FindVarFileByEnvName(msg.EnvName, m.varFiles)
		switch len(msg.Backends) {
		case 0:
			m.modal.Show(ModalState{
				Type:    ModalError,
				Title:   "❌ No Backend Config Found",
				Message: "No backend configuration found for environment: " + msg.EnvName,
				ErrorText: "Please create a backend configuration file (e.g., backend_" + msg.EnvName + ".tfvars) " +
					"to initialize this environment.",
			})
			return m, nil
		case 1:
			m.modal.Show(ModalState{
				Type:    ModalConfirm,
				Title:   "Confirm Terraform Init",
				Message: "Initialize project " + m.selectedProject.Name + " with environment " + msg.EnvName + "?\n\nUsing backend: " + msg.Backends[0].Name,
				OnConfirm: func() tea.Msg {
					return RunInitMsg{
						ProjectPath: m.selectedProject.Path,
						Options: terraform.InitOptions{
							BackendConfigFile: msg.Backends[0],
							Reconfigure:       true,
							Upgrade:           true,
							Input:             false,
						},
					}
				},
			})
			return m, nil
		default:
			var backendNames []string
			for _, b := range msg.Backends {
				backendNames = append(backendNames, b.Name+" ("+filepath.Base(b.Path)+")")

			}
			m.modal.Show(ModalState{
				Type:    ModalSelect,
				Title:   "Select Backend Config",
				Message: "Multiple backend configurations found for environment " + msg.EnvName + ". Please select one:",
				Items:   backendNames,
				OnSelect: func(backendIndex int) tea.Msg {
					return InitBackendSelectedMsg{
						EnvName: msg.EnvName,
						Backend: msg.Backends[backendIndex],
					}
				},
			})
			return m, nil
		}

	case InitBackendSelectedMsg:
		m.modal.Show(ModalState{
			Type:    ModalConfirm,
			Title:   "Confirm Terraform Init",
			Message: "Initialize project " + m.selectedProject.Name + " with environment " + msg.EnvName + "?\n\nUsing backend: " + msg.Backend.Name,
			OnConfirm: func() tea.Msg {
				return RunInitMsg{
					ProjectPath: m.selectedProject.Path,
					Options: terraform.InitOptions{
						BackendConfigFile: msg.Backend,
						Reconfigure:       true,
						Upgrade:           true,
						Input:             false,
					},
				}
			},
		})
		return m, nil

	case RunInitMsg:
		cmd := terraform.RunInit(msg.ProjectPath, msg.Options)
		return m, cmd

	case terraform.CommandCompletedMsg:
		// Update content with command output
		m.mainPanel.Content = m.mainPanel.Content + "\n\n" + string(msg.Output)
		m.backendState = terraform.DetectCurrentBackend(m.selectedProject.Path, m.backendVarFiles)
		m.sidebar.InitializedEnv = m.backendState.DetectedEnv

		return m, func() tea.Msg {
			return VarFileSelectedMsg{
				VarFileName: m.selectedVarFile.Name,
				Index:       m.sidebar.SelectedIndex,
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		titleBarHeight := 1
		statusBarHeight := 1
		contentHeight := m.height - titleBarHeight - m.header.Height - statusBarHeight - 2

		sidebarWidth := m.width / 4
		mainPanelWidth := m.width - sidebarWidth - 4

		m.titleBar.Width = m.width
		m.statusBar.Width = m.width
		m.sidebar.Width = sidebarWidth
		m.sidebar.Height = contentHeight
		m.mainPanel.Width = mainPanelWidth
		m.mainPanel.Height = contentHeight
		m.header.Width = m.width
	}
	return m, nil
}
