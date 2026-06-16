package terraform

import (
	"slices"
	"testing"
)

func TestBuildPlanArgs(t *testing.T) {
	// The table: each row is one named scenario with its inputs and expected output.
	tests := []struct {
		name        string
		projectPath string
		varFile     VarFile
		want        []string
	}{
		{
			name:        "with var file",
			projectPath: "/home/user/infra",
			varFile:     VarFile{FullPath: "/home/user/infra/dev.tfvars"},
			want:        []string{"plan", "-var-file=/home/user/infra/dev.tfvars", "-out=/tmp/lazytf-infra.tfplan"},
		},
		{
			name:        "without var file",
			projectPath: "/home/user/infra",
			varFile:     VarFile{}, // empty FullPath -> no -var-file flag
			want:        []string{"plan", "-out=/tmp/lazytf-infra.tfplan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPlanArgs(tt.projectPath, tt.varFile)
			if !slices.Equal(got, tt.want) {
				t.Errorf("buildPlanArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
