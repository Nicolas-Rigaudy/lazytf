package ui

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nicolas-Rigaudy/lazytf/internal/aws"
	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode int

const (
	ViewModeProjectList ViewMode = iota
	ViewModeProjectDetail
)

// InteractiveCommand tracks state for commands that need user confirmation
type InteractiveCommand struct {
	StdinWriter   io.WriteCloser // Writer to send input to the command
	OutputBuffer  string         // Buffered output before confirmation
	CommandName   string         // Name of the command (e.g., "apply", "destroy")
	IsWaiting     bool           // True when waiting for user confirmation
	ListenNext    tea.Cmd        // Next listener command to continue receiving output
}

type Model struct {
	sidebar                  SidebarModel
	mainPanel                MainPanelModel
	titleBar                 TitleBarModel
	statusBar                StatusBarModel
	header                   HeaderModel
	focusIndex               int
	focusableCount           int
	width                    int
	height                   int
	projects                 []terraform.Project
	mode                     terraform.Mode
	viewMode                 ViewMode
	selectedProject          *terraform.Project
	varFiles                 []terraform.VarFile
	selectedVarFile          *terraform.VarFile
	backendVarFiles          []terraform.BackendVarFile
	selectedBackendFile      *terraform.BackendVarFile
	backendState             terraform.BackendState
	modal                    Modal              // Modal component
	lastCommandMsg           tea.Msg            // Last terraform command message (for retry after AWS login)
	pendingTerraformCommand  func() tea.Msg     // Command to run after AWS login completes (set by detectCommonErrors)
	interactiveCmd           *InteractiveCommand // Tracks state for interactive commands (apply, destroy, etc.)
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

// prepareCommandExecution clears the main panel and sets it up to display command output
func (m *Model) prepareCommandExecution() {
	m.mainPanel.Title = "⏳ Running Command"
	m.mainPanel.SetContent("")
	m.mainPanel.HasError = false
}

// triggerAWSLogin initiates the AWS SSO login workflow
// Returns a tea.Cmd to execute, or nil if showing a modal
func (m *Model) triggerAWSLogin() tea.Cmd {
	sessions, err := aws.DiscoverSSOSessions()
	if err != nil || len(sessions) == 0 {
		m.modal.Show(ModalState{
			Type:  ModalError,
			Title: "❌ No AWS SSO Sessions Found",
			ErrorText: "No AWS SSO sessions found in your AWS config file. Please configure at least one SSO session in ~/.aws/config " +
				"before attempting to log in.",
		})
		return nil
	}

	if len(sessions) == 1 {
		// Only one session, run login directly
		m.prepareCommandExecution()
		return aws.RunSSOLogin(sessions[0])
	}

	// Multiple sessions, show selection modal
	var sessionNames []string
	for _, session := range sessions {
		sessionNames = append(sessionNames, session.Name)
	}
	m.modal.Show(ModalState{
		Type:    ModalSelect,
		Title:   "AWS SSO Login",
		Message: "Select an AWS SSO session to log in:",
		Items:   sessionNames,
		OnSelect: func(index int) tea.Msg {
			return RunAWSSSOLoginMsg{
				Session: sessions[index],
			}
		},
	})
	return nil
}

// detectCommonErrors checks command output for common error patterns
// Returns true if an error was detected and handled with a modal
func (m *Model) detectCommonErrors() bool {
	output := m.mainPanel.Content

	// Check for AWS credential errors
	if strings.Contains(output, "No valid credential sources found") ||
		strings.Contains(output, "Unable to locate credentials") {
		// Set pending command to retry after AWS login
		if m.lastCommandMsg != nil {
			m.pendingTerraformCommand = func() tea.Msg {
				return m.lastCommandMsg
			}
		}

		m.modal.Show(ModalState{
			Type:    ModalConfirm,
			Title:   "🔐 AWS Authentication Required",
			Message: "AWS credentials are not available or have expired.\n\nWould you like to log in now?",
			OnConfirm: func() tea.Msg {
				return TriggerAWSLoginMsg{}
			},
		})
		return true
	}

	// Check for backend not initialized
	if strings.Contains(output, "Backend initialization required") ||
		strings.Contains(output, "terraform init") {
		m.modal.Show(ModalState{
			Type:      ModalError,
			Title:     "❌ Terraform Not Initialized",
			ErrorText: "Please run terraform init (press 'i') before running this command.",
		})
		return true
	}

	return false
}

// buildStatusText creates dynamic status bar text based on current state
func (m Model) buildStatusText() string {
	var parts []string

	// Add help keys based on view mode
	if m.viewMode == ViewModeProjectDetail {
		if m.mode == terraform.ModeMultiProject {
			parts = append(parts, "Backspace: back", "│")
		}
		parts = append(parts, "i: init", "p: plan", "a: apply", "l: aws login", "│")
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

		case "l":
			// Trigger AWS SSO Login
			return m, m.triggerAWSLogin()

		case "p":
			// Trigger Terraform Plan
			if m.viewMode == ViewModeProjectDetail && m.selectedProject != nil && m.selectedVarFile != nil {
				// Check if the environment is initialized first
				if !m.backendState.IsInitialized || m.backendState.DetectedEnv != m.selectedVarFile.EnvName {
					m.modal.Show(ModalState{
						Type:      ModalError,
						Title:     "❌ Environment Not Initialized",
						ErrorText: "Please run terraform init (press 'i') before running plan.",
					})
					return m, nil
				}

				// Run plan - AWS credential check happens reactively if command fails
				m.prepareCommandExecution()
				return m, func() tea.Msg {
					return RunPlanMsg{
						ProjectPath: m.selectedProject.Path,
						VarFile:     *m.selectedVarFile,
					}
				}
			}

		case "a":
			// Trigger Terraform Apply
			if m.viewMode == ViewModeProjectDetail && m.selectedProject != nil && m.selectedVarFile != nil {
				// Check if the environment is initialized first
				if !m.backendState.IsInitialized || m.backendState.DetectedEnv != m.selectedVarFile.EnvName {
					m.modal.Show(ModalState{
						Type:      ModalError,
						Title:     "❌ Environment Not Initialized",
						ErrorText: "Please run terraform init (press 'i') before running apply.",
					})
					return m, nil
				}

				// Run apply - will show plan and wait for confirmation
				m.prepareCommandExecution()
				return m, func() tea.Msg {
					return RunApplyMsg{
						ProjectPath: m.selectedProject.Path,
						VarFile:     *m.selectedVarFile,
					}
				}
			}

		case "tab":
			m.focusIndex = (m.focusIndex + 1) % m.focusableCount
			m.updateFocusStates()

		case "backspace", "esc":
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
		content := "Project: " + selectedProject.Name + "\n" +
			"Path: " + selectedProject.Path + "\n\n" +
			"--- Backend Status ---\n" +
			backendStateInfo + "\n" +
			"--- Available Environments ---\n" +
			"Var Files: " + string(rune(len(m.varFiles)+'0')) + "\n" +
			"Backend Configs: " + string(rune(len(m.backendVarFiles)+'0'))
		m.mainPanel.SetContent(content)

		// Update status bar
		m.statusBar.SetText(m.buildStatusText())

		return m, nil

	case VarFileSelectedMsg:
		// Get the selected var file
		selectedVarFile := &m.varFiles[msg.Index]
		m.selectedVarFile = selectedVarFile

		// Determine initialization status for this environment
		isThisEnvInitialized := m.backendState.IsInitialized && m.backendState.DetectedEnv == selectedVarFile.EnvName

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
		// Store this command in case we need to retry after AWS login
		m.lastCommandMsg = msg
		m.prepareCommandExecution()
		cmd := terraform.RunInit(msg.ProjectPath, msg.Options)
		return m, cmd

	case ApplyConfirmedMsg:
		if m.interactiveCmd == nil || m.interactiveCmd.StdinWriter == nil {
			return m, nil
		}

		// User confirmed apply - send "yes" to terraform stdin
		m.mainPanel.Title = "🚀 Applying changes..."

		// Save the listener before modifying state
		listenNext := m.interactiveCmd.ListenNext

		// Write to stdin and close
		writer := m.interactiveCmd.StdinWriter
		writer.Write([]byte("yes\n"))
		writer.Close()

		// Clear interactiveCmd - we're done with the interactive flow
		m.interactiveCmd = nil

		// Continue listening for terraform's output
		return m, listenNext

	case ApplyCancelledMsg:
		// User cancelled apply - send "no" to terraform stdin to abort gracefully
		if m.interactiveCmd != nil && m.interactiveCmd.StdinWriter != nil {
			m.mainPanel.Title = "⛔ Apply cancelled"
			m.mainPanel.Content += "\n[Cancelling apply...]\n"
			m.mainPanel.SetContent(m.mainPanel.Content)

			// Save the listener before modifying state
			listenNext := m.interactiveCmd.ListenNext

			// Write to stdin in goroutine to avoid blocking
			writer := m.interactiveCmd.StdinWriter
			go func() {
				writer.Write([]byte("no\n"))
				writer.Close()
			}()

			// Clear interactiveCmd - we're done with the interactive flow
			m.interactiveCmd = nil

			// Continue listening for terraform's output
			return m, listenNext
		}
		return m, nil

	case TriggerAWSLoginMsg:
		// Trigger the AWS login workflow
		return m, m.triggerAWSLogin()

	case RunAWSSSOLoginMsg:
		m.prepareCommandExecution()
		return m, aws.RunSSOLogin(msg.Session)

	case RunPlanMsg:
		// Store this command in case we need to retry after AWS login
		m.lastCommandMsg = msg
		m.prepareCommandExecution()
		cmd := terraform.RunPlan(msg.ProjectPath, msg.VarFile)
		return m, cmd

	case RunApplyMsg:
		// Store this command in case we need to retry after AWS login
		m.lastCommandMsg = msg
		m.prepareCommandExecution()

		// Run terraform apply (without auto-approve) - returns stdin writer
		result := terraform.RunApply(msg.ProjectPath, msg.VarFile)

		// Store the interactive command state
		m.interactiveCmd = &InteractiveCommand{
			StdinWriter:  result.StdinWriter,
			OutputBuffer: "",
			CommandName:  "apply",
			IsWaiting:    false, // Will become true when we detect the confirmation prompt
		}

		return m, result.Cmd

	case executor.CommandOutputMsg:
		// Debug: show all output regardless of state
		m.mainPanel.Content += msg.Line + "\n"
		m.mainPanel.SetContent(m.mainPanel.Content)

		// If we have an interactive command running, also buffer the output
		if m.interactiveCmd != nil && !m.interactiveCmd.IsWaiting {
			m.interactiveCmd.OutputBuffer += msg.Line + "\n"
			// Store the listener so we can continue after modal interaction
			m.interactiveCmd.ListenNext = msg.ListenNext

			// Detect terraform's confirmation prompt
			// Note: "Enter a value:" never arrives as a complete line (no newline until we respond)
			// so we detect on "Only 'yes' will be accepted" which comes before the prompt
			lineLower := strings.ToLower(strings.TrimSpace(msg.Line))
			if strings.Contains(lineLower, "only 'yes' will be accepted") {
				// Switch to waiting state and show confirmation modal
				m.interactiveCmd.IsWaiting = true

				// Show confirmation modal
				envName := ""
				if m.selectedVarFile != nil {
					envName = m.selectedVarFile.EnvName
				}

				m.modal.Show(ModalState{
					Type:    ModalConfirm,
					Title:   fmt.Sprintf("🚀 Apply changes to %s?", envName),
					Message: "Confirm to proceed with apply, or cancel to abort.",
					OnConfirm: func() tea.Msg {
						return ApplyConfirmedMsg{}
					},
					OnCancel: func() tea.Msg {
						return ApplyCancelledMsg{}
					},
				})

				return m, msg.ListenNext
			}

			return m, msg.ListenNext
		}

		// Keep listening for more output
		return m, msg.ListenNext

	case executor.CommandCompletedMsg:
		// Command finished - update title and refresh state
		m.mainPanel.Title = "✅ Command Completed"
		m.mainPanel.HasError = false

		// Clean up interactive command state
		if m.interactiveCmd != nil {
			if m.interactiveCmd.StdinWriter != nil {
				m.interactiveCmd.StdinWriter.Close()
			}
			m.interactiveCmd = nil
		}

		// Refresh backend state to update sidebar indicators (only if we have a project)
		if m.selectedProject != nil {
			m.backendState = terraform.DetectCurrentBackend(m.selectedProject.Path, m.backendVarFiles)
			m.sidebar.InitializedEnv = m.backendState.DetectedEnv
		}

		// If there's a pending terraform command (after AWS login), execute it now
		if m.pendingTerraformCommand != nil {
			cmd := m.pendingTerraformCommand
			m.pendingTerraformCommand = nil // Clear it
			return m, cmd
		}

		return m, nil

	case executor.CommandErrorMsg:
		m.mainPanel.Title = fmt.Sprintf("❌ Command Failed (exit code %d)", msg.ExitCode)
		m.mainPanel.HasError = true

		// Clean up interactive command state
		if m.interactiveCmd != nil {
			if m.interactiveCmd.StdinWriter != nil {
				m.interactiveCmd.StdinWriter.Close()
			}
			m.interactiveCmd = nil
		}

		// Detect common errors and show helpful modals
		m.detectCommonErrors()

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate component heights
		titleBarHeight := 1
		statusBarHeight := 1
		headerHeight := m.header.Height // Default 5, updates dynamically in View()

		// JoinVertical adds newlines between components (title|header|content|status = 3 separators)
		separatorLines := 3

		// Available height for the content area (panels handle their own borders/padding)
		totalContentHeight := m.height - titleBarHeight - headerHeight - statusBarHeight - separatorLines

		// The panel height includes its border/padding, so we give it the full space
		panelHeight := totalContentHeight

		sidebarWidth := m.width / 4
		mainPanelWidth := m.width - sidebarWidth - 4

		m.titleBar.Width = m.width
		m.statusBar.Width = m.width
		m.header.Width = m.width
		m.sidebar.Width = sidebarWidth
		m.sidebar.Height = panelHeight
		m.mainPanel.SetSize(mainPanelWidth, panelHeight)
	}
	return m, nil
}
