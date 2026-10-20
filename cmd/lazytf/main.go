// Box drawing characters reference:
// ╔═╗  ╭─╮  ┌─┐  ┏━┓
// ║ ║  │ │  │ │  ┃ ┃
// ╚═╝  ╰─╯  └─┘  ┗━┛
// ─ │ ═ ║ ┼ ╬

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	opts, err := parseArgs(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "lazytf: %v\nTry 'lazytf --help'.\n", err)
		os.Exit(2)
	}
	if opts.showVersion {
		fmt.Println("lazytf " + versionString())
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lazytf: %v\n", err)
		os.Exit(1)
	}
	st, err := resolveStart(opts, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lazytf: %v\n", err)
		os.Exit(1)
	}
	if err := checkBinary(st.cfg.Terraform); err != nil {
		fmt.Fprintf(os.Stderr, "lazytf: %v\n", err)
		os.Exit(1)
	}

	// Initialize UI model
	m := ui.NewModel(st.projects, st.mode, st.cfg.Layout)
	if note := terraform.CleanupLegacyPlans(); note != "" {
		m.SetNotice(note)
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
