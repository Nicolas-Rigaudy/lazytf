package ui

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/auth"
	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
)

func providerFor(t *testing.T, cfg auth.Config) *auth.Provider {
	t.Helper()
	p, err := auth.Resolve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func newAuthModel(p *auth.Provider, output string) Model {
	m := Model{modal: NewModal(), authProvider: p}
	m.mainPanel.Content = output
	return m
}

func TestExpiredGCloudCredentialsPromptLoginAndRetry(t *testing.T) {
	m := newAuthModel(providerFor(t, auth.Config{Provider: auth.GCloud}),
		"Error: google: could not find default credentials.")
	plan := RunPlanMsg{ProjectPath: "/p", VarFile: terraform.VarFile{EnvName: "dev"}}
	m.lastCommandMsg = plan

	if !m.detectCommonErrors() || !m.modal.IsActive() {
		t.Fatal("expired gcloud credentials did not open the login prompt")
	}
	if !strings.Contains(m.modal.state.Message, "Google Cloud") || strings.Contains(m.modal.state.Message, "AWS") {
		t.Errorf("prompt = %q, want the gcloud label and no AWS wording", m.modal.state.Message)
	}
	if got := m.modal.state.OnConfirm(); got != (TriggerLoginMsg{}) {
		t.Errorf("confirming the prompt sent %#v, want TriggerLoginMsg", got)
	}
	// After the login finishes, the failed plan is the command that runs again.
	if m.pendingTerraformCommand == nil || m.pendingTerraformCommand() != tea.Msg(plan) {
		t.Errorf("pending command does not retry the failed plan")
	}
}

func TestNoProviderNeverMentionsLogin(t *testing.T) {
	m := newAuthModel(nil, "Error: No valid credential sources found")
	if m.detectCommonErrors() || m.modal.IsActive() {
		t.Error("with no auth provider an AWS credential error opened a login prompt")
	}

	if cmd := m.triggerLogin(); cmd != nil {
		t.Error("triggerLogin with no provider returned a command")
	}
	if !m.modal.IsActive() || strings.Contains(m.modal.state.Title+m.modal.state.ErrorText, "AWS SSO") {
		t.Errorf("modal = %+v, want a neutral 'no login configured' error", m.modal.state)
	}
}

func TestAWSProviderKeepsItsErrorsAndLabel(t *testing.T) {
	m := newAuthModel(providerFor(t, auth.Config{Provider: auth.AWSSSO}), "Error: Unable to locate credentials")
	if !m.detectCommonErrors() || !strings.Contains(m.modal.state.Message, "AWS credentials") {
		t.Errorf("aws-sso preset: modal = %+v", m.modal.state)
	}
}

func TestLoginKeyOnlyWithProvider(t *testing.T) {
	hasLogin := func(groups []KeyBindingGroup) bool {
		for _, g := range groups {
			for _, kb := range g.KeyBindings {
				if kb.Key == "L" {
					return true
				}
			}
		}
		return false
	}
	if hasLogin(DefaultKeyBindings(false)) || !hasLogin(DefaultKeyBindings(true)) {
		t.Error("help lists the login key when no provider is set, or hides it when one is")
	}

	m := Model{viewMode: ViewModeProjectDetail}
	if strings.Contains(m.buildStatusText(), "login") {
		t.Errorf("status = %q, want no login hint without a provider", m.buildStatusText())
	}
	m.authProvider = providerFor(t, auth.Config{Provider: auth.Az})
	if !strings.Contains(m.buildStatusText(), "shift+l: login") || strings.Contains(m.buildStatusText(), "aws") {
		t.Errorf("status = %q, want a neutral login hint", m.buildStatusText())
	}
}

func TestLoadProjectResolvesAuthAndBinary(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/.lazytf.yml", "terraform: tofu\nauth:\n  provider: gcloud\n"); err != nil {
		t.Fatal(err)
	}
	setup := loadProject(dir, defaultLayoutForTest())
	if setup.bin != "tofu" || setup.auth == nil || setup.auth.Name != auth.GCloud || setup.note != "" {
		t.Errorf("setup = %+v, want tofu + gcloud and no note", setup)
	}
	if len(setup.varFiles) != 1 || setup.varFiles[0].EnvName != terraform.NoVarFileEnv {
		t.Errorf("varFiles = %+v, want the stand-in env", setup.varFiles)
	}

	if err := writeFile(dir+"/.lazytf.yml", "auth:\n  provider: okta\n"); err != nil {
		t.Fatal(err)
	}
	if setup := loadProject(dir, defaultLayoutForTest()); !strings.Contains(setup.note, "okta") {
		t.Errorf("note = %q, want it to name the unknown provider", setup.note)
	}
}

func writeFile(path, content string) error { return os.WriteFile(path, []byte(content), 0o644) }

func defaultLayoutForTest() config.Layout { return config.DefaultLayout() }
