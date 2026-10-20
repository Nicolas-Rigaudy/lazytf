//go:build integration

package terraform

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// run drives a command through the executor the way the UI does and returns
// its combined output and whether it succeeded.
func run(t *testing.T, step string, cmd tea.Cmd, timeout time.Duration) (string, bool) {
	t.Helper()
	var out strings.Builder
	deadline := time.After(timeout)
	for {
		msgCh := make(chan tea.Msg, 1)
		go func() { msgCh <- cmd() }()
		select {
		case <-deadline:
			t.Fatalf("%s: did not finish within %s; output so far:\n%s", step, timeout, out.String())
		case msg := <-msgCh:
			switch m := msg.(type) {
			case executor.CommandOutputMsg:
				out.WriteString(m.Line + "\n")
				cmd = m.ListenNext
			case executor.CommandCompletedMsg:
				return out.String(), true
			case executor.CommandErrorMsg:
				return out.String(), false
			default:
				t.Fatalf("%s: unexpected message %#v", step, msg)
			}
		}
	}
}

func mustRun(t *testing.T, step string, cmd tea.Cmd) {
	t.Helper()
	if out, ok := run(t, step, cmd, 3*time.Minute); !ok {
		t.Fatalf("%s failed:\n%s", step, out)
	}
}

// TestIntegrationInitPlanApply runs the real init/plan/apply/destroy flow
// against testdata/fixture: local backend, null provider.
func TestIntegrationInitPlanApply(t *testing.T) {
	if _, err := exec.LookPath("terraform"); err != nil {
		t.Fatal("terraform not on PATH; the integration test needs it")
	}
	t.Setenv("XDG_CACHE_HOME", t.TempDir()) // keep saved plans out of the real cache
	t.Setenv("TF_IN_AUTOMATION", "1")

	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("testdata/fixture")); err != nil {
		t.Fatal(err)
	}

	varFiles, _ := DiscoverVarFiles(dir, config.DefaultLayout())
	dev, _ := FindVarFileByEnvName("dev", varFiles)
	if dev == nil {
		t.Fatalf("dev var file not discovered: %+v", varFiles)
	}
	backends, _ := DiscoverBackendVarFiles(dir, config.DefaultLayout())
	devBackends := MatchBackendsForEnv("dev", backends)
	if len(devBackends) != 1 {
		t.Fatalf("dev backend config not matched: %+v", backends)
	}

	mustRun(t, "init", RunInit("", dir, InitOptions{BackendConfigFile: devBackends[0]}))
	state := DetectCurrentBackend(dir, backends)
	if state.BackendType != "local" || state.DetectedEnv != "dev" {
		t.Fatalf("after init: backend %q env %q, want local / dev", state.BackendType, state.DetectedEnv)
	}

	// AC3: a missing required variable fails at once with Terraform's error.
	out, ok := run(t, "plan without var file", RunPlan("", dir, VarFile{}), 30*time.Second)
	if ok || !strings.Contains(out, `"name" is not set`) {
		t.Fatalf("plan without var file: ok=%v, want a failure naming the variable; output:\n%s", ok, out)
	}

	mustRun(t, "plan", RunPlan("", dir, *dev))
	planFile := PlanFilePath(dir, "dev")
	if _, err := os.Stat(planFile); err != nil {
		t.Fatalf("saved plan missing: %v", err)
	}
	// AC1: another env never sees this plan.
	if _, err := os.Stat(PlanFilePath(dir, "prod")); !os.IsNotExist(err) {
		t.Fatalf("prod plan path exists after planning dev (err=%v)", err)
	}

	planJSON, err := RunShowPlanJSON("", dir, planFile)
	if err != nil {
		t.Fatalf("terraform show -json: %v", err)
	}
	plan, err := ParsePlanJSON(planJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.ResourceChanges) != 1 || plan.ResourceChanges[0].Address != "null_resource.x" ||
		plan.ResourceChanges[0].Change.Actions[0] != "create" {
		t.Fatalf("plan changes = %+v, want one create of null_resource.x", plan.ResourceChanges)
	}

	apply := RunApply("", dir, *dev, planFile)
	mustRun(t, "apply saved plan", apply.Cmd)
	if _, err := os.Stat(filepath.Join(dir, "states", "dev", "terraform.tfstate")); err != nil {
		t.Fatalf("state not written to the dev backend path: %v", err)
	}

	// Destroy confirms over stdin, as the UI does; -input=false must not
	// stop the approval prompt from reading it.
	destroy := RunDestroy("", dir, *dev)
	if _, err := io.WriteString(destroy.StdinWriter, "yes\n"); err != nil {
		t.Fatal(err)
	}
	out, ok = run(t, "destroy", destroy.Cmd, 3*time.Minute)
	if !ok || !strings.Contains(out, "Destroy complete") {
		t.Fatalf("destroy: ok=%v, output:\n%s", ok, out)
	}
}

// TestIntegrationNoVarFilesLocalBackend covers a project with no tfvars and no
// backend config: it gets the stand-in environment, inits without
// -backend-config, plans and applies.
func TestIntegrationNoVarFilesLocalBackend(t *testing.T) {
	if _, err := exec.LookPath("terraform"); err != nil {
		t.Fatal("terraform not on PATH; the integration test needs it")
	}
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	t.Setenv("TF_IN_AUTOMATION", "1")

	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("testdata/fixture-novars")); err != nil {
		t.Fatal(err)
	}

	layout := config.DefaultLayout()
	discovered, _ := DiscoverVarFiles(dir, layout)
	varFiles := WithDefaultEnv(discovered)
	env := varFiles[0]
	if env.EnvName != NoVarFileEnv {
		t.Fatalf("env = %+v, want the stand-in", env)
	}
	backends, _ := DiscoverBackendVarFiles(dir, layout)

	mustRun(t, "init", RunInit("", dir, InitOptions{}))
	if err := RecordInit(dir, env.EnvName, BackendVarFile{}); err != nil {
		t.Fatal(err)
	}
	state := DetectCurrentBackend(dir, backends)
	if !state.EnvReady(env.EnvName) {
		t.Fatalf("after init: state %+v is not ready for %q", state, env.EnvName)
	}

	mustRun(t, "plan", RunPlan("", dir, env))
	planFile := PlanFilePath(dir, env.EnvName)
	if _, err := os.Stat(planFile); err != nil {
		t.Fatalf("saved plan missing: %v", err)
	}
	apply := RunApply("", dir, env, planFile)
	mustRun(t, "apply saved plan", apply.Cmd)
	if _, err := os.Stat(filepath.Join(dir, "terraform.tfstate")); err != nil {
		t.Fatalf("local state not written: %v", err)
	}
}
