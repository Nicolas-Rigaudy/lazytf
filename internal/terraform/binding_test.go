package terraform

import "testing"

const s3KeyTeamDev = `{"backend":{"type":"s3","config":{"bucket":"b","key":"team/dev.tfstate"}}}`

func TestRecordedBindingBeatsInference(t *testing.T) {
	cache, root := t.TempDir(), t.TempDir()
	writeFiles(t, root, map[string]string{".terraform/terraform.tfstate": s3KeyTeamDev})
	backends := []BackendVarFile{{Name: "dev.hcl", Path: "backends", EnvName: "dev"}}

	// Without a record the key guess ("team") is wrong.
	if got := detectCurrentBackend(cache, root, backends).DetectedEnv; got != "team" {
		t.Fatalf("inferred env = %q, want the guess %q", got, "team")
	}

	if err := recordInit(cache, root, "dev", backends[0]); err != nil {
		t.Fatal(err)
	}
	got := detectCurrentBackend(cache, root, backends)
	if got.DetectedEnv != "dev" {
		t.Errorf("DetectedEnv = %q, want dev from the record", got.DetectedEnv)
	}
	if len(got.MatchedBackends) != 1 || got.MatchedBackends[0].Name != "dev.hcl" {
		t.Errorf("MatchedBackends = %+v, want dev.hcl", got.MatchedBackends)
	}
}

func TestBindingIgnoredAfterManualReinit(t *testing.T) {
	cache, root := t.TempDir(), t.TempDir()
	writeFiles(t, root, map[string]string{".terraform/terraform.tfstate": s3KeyTeamDev})
	if err := recordInit(cache, root, "dev", BackendVarFile{}); err != nil {
		t.Fatal(err)
	}
	// Someone re-inits by hand against another state key.
	writeFiles(t, root, map[string]string{".terraform/terraform.tfstate": `{"backend":{"type":"s3","config":{"bucket":"b","key":"prod/terraform.tfstate"}}}`})
	if got := detectCurrentBackend(cache, root, nil).DetectedEnv; got != "prod" {
		t.Errorf("DetectedEnv = %q, want the inferred prod, not the stale record", got)
	}
}

func TestBindingForLocalBackend(t *testing.T) {
	cache, root := t.TempDir(), t.TempDir()
	writeFiles(t, root, map[string]string{".terraform/providers/x": ""}) // initialised, no backend block
	devVars := []BackendVarFile{{Name: "backend_dev.tfvars", EnvName: "dev"}}
	if got := detectCurrentBackend(cache, root, devVars).DetectedEnv; got != "" {
		t.Fatalf("DetectedEnv = %q before any record", got)
	}
	if err := recordInit(cache, root, "dev", BackendVarFile{}); err != nil {
		t.Fatal(err)
	}
	got := detectCurrentBackend(cache, root, devVars)
	if got.DetectedEnv != "dev" || !got.EnvReady("dev") || got.EnvReady("prod") {
		t.Errorf("state = %+v, want ready for dev only", got)
	}
}

func TestEnvReady(t *testing.T) {
	tests := []struct {
		name  string
		state BackendState
		env   string
		want  bool
	}{
		{"not initialised", BackendState{}, "dev", false},
		{"not initialised even with no backends", BackendState{AnyEnv: true}, "dev", false},
		{"matching env", BackendState{IsInitialized: true, DetectedEnv: "dev"}, "dev", true},
		{"other env", BackendState{IsInitialized: true, DetectedEnv: "dev"}, "prod", false},
		{"no backend configs serve every env", BackendState{IsInitialized: true, AnyEnv: true}, "prod", true},
	}
	for _, tt := range tests {
		if got := tt.state.EnvReady(tt.env); got != tt.want {
			t.Errorf("%s: EnvReady(%q) = %v, want %v", tt.name, tt.env, got, tt.want)
		}
	}
}

func TestWithDefaultEnv(t *testing.T) {
	got := WithDefaultEnv(nil)
	if len(got) != 1 || got[0].EnvName != NoVarFileEnv || got[0].FullPath != "" {
		t.Errorf("WithDefaultEnv(nil) = %+v, want one stand-in without a var file", got)
	}
	in := []VarFile{{EnvName: "dev"}}
	if got := WithDefaultEnv(in); len(got) != 1 || got[0].EnvName != "dev" {
		t.Errorf("WithDefaultEnv(in) = %+v, want in unchanged", got)
	}
}
