package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestRenderScreens writes the main screen as ANSI text, e.g.
//
//	LAZYTF_RENDER_DIR=/path go test ./internal/ui -run TestRenderScreens
func TestRenderScreens(t *testing.T) {
	outDir := os.Getenv("LAZYTF_RENDER_DIR")
	if outDir == "" {
		t.Skip("LAZYTF_RENDER_DIR not set")
	}
	lipgloss.SetColorProfile(termenv.TrueColor)

	proj := t.TempDir()
	for _, f := range []string{"main.tf", "dev.tfvars", "staging.tfvars", "prod.tfvars"} {
		if err := os.WriteFile(filepath.Join(proj, f), []byte("# sample\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan := &terraform.PlanSummary{ResourceChanges: []terraform.ResourceChange{
		rc("aws_s3_bucket.logs", "", "create"),
		rc("aws_s3_bucket.assets", "", "update"),
		rc("aws_iam_role.old", "", "delete"),
		rc("module.vpc.aws_subnet.a", "module.vpc", "delete", "create"),
		rc("module.vpc.aws_subnet.b", "module.vpc", "create"),
	}}

	for _, size := range [][2]int{{120, 35}, {80, 24}} {
		m := NewModel([]terraform.Project{{Name: "app", Path: proj}}, terraform.ModeSingleProject, config.DefaultLayout())
		next, _ := m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		m = next.(Model)
		vm := NewPlanSummaryViewModel(plan)
		vm.ModuleGroups[0].Expanded = true
		vm.SelectedResourceIndex = 1
		vm.updateContent()
		vm.SetSize(m.mainPanel.viewport.Width, m.mainPanel.viewport.Height)
		m.mainPanel.PlanSummary = vm
		m.mainPanel.ViewMode = ViewModeSummary
		m.mainPanel.Title = vm.Title()
		m.mainPanel.IsFocused = true
		m.sidebar.IsFocused = false
		name := filepath.Join(outDir, fmt.Sprintf("main-%dx%d.ans", size[0], size[1]))
		if err := os.WriteFile(name, []byte(m.View()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
