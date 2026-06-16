package terraform

import (
	"slices"
	"testing"
)

func TestBuildInitArgs(t *testing.T) {
	tests := []struct {
		name    string
		options InitOptions
		want    []string
	}{
		{
			// Zero value: Input is false, so !Input is true -> -input=false appears.
			// This pins down the (counterintuitive) inverted Input flag.
			name:    "zero value emits input=false",
			options: InitOptions{},
			want:    []string{"init", "-input=false"},
		},
		{
			name:    "input enabled emits no flags",
			options: InitOptions{Input: true},
			want:    []string{"init"},
		},
		{
			name:    "reconfigure and upgrade",
			options: InitOptions{Reconfigure: true, Upgrade: true, Input: true},
			want:    []string{"init", "-reconfigure", "-upgrade"},
		},
		{
			name: "backend config only",
			options: InitOptions{
				BackendConfigFile: BackendVarFile{
					Name:     "backend_dev.tfvars",
					FullPath: "/home/user/infra/variables/backend/local/backend_dev.tfvars",
				},
				Input: true,
			},
			want: []string{"init", "-backend-config=/home/user/infra/variables/backend/local/backend_dev.tfvars"},
		},
		{
			// All options on: verifies flag ORDER is backend -> reconfigure -> upgrade.
			name: "everything on",
			options: InitOptions{
				BackendConfigFile: BackendVarFile{
					Name:     "backend_dev.tfvars",
					FullPath: "/x/backend_dev.tfvars",
				},
				Reconfigure: true,
				Upgrade:     true,
				Input:       true,
			},
			want: []string{"init", "-backend-config=/x/backend_dev.tfvars", "-reconfigure", "-upgrade"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildInitArgs(tt.options)
			if !slices.Equal(got, tt.want) {
				t.Errorf("buildInitArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
