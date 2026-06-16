package terraform

import (
	"slices"
	"testing"
)

func TestBuildDestroyArgs(t *testing.T) {
	tests := []struct {
		name    string
		varFile VarFile
		want    []string
	}{
		{
			name:    "with var file",
			varFile: VarFile{FullPath: "/home/user/infra/dev.tfvars"},
			want:    []string{"destroy", "-var-file=/home/user/infra/dev.tfvars"},
		},
		{
			name:    "without var file",
			varFile: VarFile{},
			want:    []string{"destroy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDestroyArgs(tt.varFile)
			if !slices.Equal(got, tt.want) {
				t.Errorf("buildDestroyArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
