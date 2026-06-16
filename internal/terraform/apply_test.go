package terraform

import (
	"slices"
	"testing"
)

func TestBuildApplyArgs(t *testing.T) {
	tests := []struct {
		name      string
		varFile   VarFile
		planFile  string
		wantArgs  []string
		wantStdin bool
	}{
		{
			// Plan file present: apply it directly, no interactive confirmation needed.
			name:      "plan file disables stdin",
			varFile:   VarFile{FullPath: "/home/user/infra/dev.tfvars"}, // ignored when planFile set
			planFile:  "/tmp/lazytf-infra.tfplan",
			wantArgs:  []string{"apply", "/tmp/lazytf-infra.tfplan"},
			wantStdin: false,
		},
		{
			// No plan file but a var file: interactive apply, needs stdin for "yes".
			name:      "var file enables stdin",
			varFile:   VarFile{FullPath: "/home/user/infra/dev.tfvars"},
			planFile:  "",
			wantArgs:  []string{"apply", "-var-file=/home/user/infra/dev.tfvars"},
			wantStdin: true,
		},
		{
			// Neither: bare interactive apply, still needs stdin.
			name:      "no plan or var file still enables stdin",
			varFile:   VarFile{},
			planFile:  "",
			wantArgs:  []string{"apply"},
			wantStdin: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotStdin := buildApplyArgs(tt.varFile, tt.planFile)
			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("buildApplyArgs() args = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotStdin != tt.wantStdin {
				t.Errorf("buildApplyArgs() stdin = %v, want %v", gotStdin, tt.wantStdin)
			}
		})
	}
}
