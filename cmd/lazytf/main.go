// Box drawing characters reference:
// ╔═╗  ╭─╮  ┌─┐  ┏━┓
// ║ ║  │ │  │ │  ┃ ┃
// ╚═╝  ╰─╯  └─┘  ┗━┛
// ─ │ ═ ║ ┼ ╬

package main

import (
	"fmt"
	"os"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Determine mode
	mode, project, err := terraform.DetermineMode()
	if err != nil {
		fmt.Printf("Error determining mode: %v\n", err)
		os.Exit(1)
	}

	var projects []terraform.Project
	switch mode {
	case terraform.ModeSingleProject:
		projects = []terraform.Project{*project}
	case terraform.ModeMultiProject:
		projects, err = terraform.DiscoverProjects(cfg)
		if err != nil {
			fmt.Printf("Error discovering projects: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize UI model
	m := ui.NewModel(projects, mode)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
