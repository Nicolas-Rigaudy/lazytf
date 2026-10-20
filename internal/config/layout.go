package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectFileName is the per-repo layout override, read from the project root.
const ProjectFileName = ".lazytf.yml"

// Layout says where a project keeps its var files and backend configs and how
// environment names are derived from them. It is read from the global config
// and overridden key by key by a project's .lazytf.yml.
//
// Globs are relative to the project root and may use "**" to span directories.
type Layout struct {
	// Terraform is the executable to run, "terraform" when empty. A name
	// such as "tofu" is looked up on PATH; a path is used as given.
	Terraform string `yaml:"terraform,omitempty"`
	// VarFiles are the globs whose matches are environments.
	VarFiles []string `yaml:"var_files"`
	// BackendConfigs are the globs of files passed as -backend-config to init.
	BackendConfigs []string `yaml:"backend_configs"`
	// Exclude removes var files from the environments. A glob without "/"
	// is matched against the file name, one with "/" against the full
	// relative path.
	Exclude []string `yaml:"exclude"`
	// Include adds var files that are listed even when Exclude matches them.
	Include []string `yaml:"include,omitempty"`
	// EnvPattern is a regular expression with one capture group, applied to
	// the var file name without its extension. Empty means the whole name.
	EnvPattern string `yaml:"env_pattern,omitempty"`
	// BackendEnvPattern is EnvPattern for backend config files. Empty means
	// the "backend_<env>" convention (a bare "backend" is shared by all envs).
	BackendEnvPattern string `yaml:"backend_env_pattern,omitempty"`
}

// DefaultLayout is the layout LazyTF understood before it became configurable:
// *.tfvars in ., variables/, env/ and tfvars/, backend configs under
// variables/backend/. Terraform loads terraform.tfvars and *.auto.tfvars on
// its own, so they are not environments.
func DefaultLayout() Layout {
	return Layout{
		VarFiles:       []string{"*.tfvars", "variables/*.tfvars", "env/*.tfvars", "tfvars/*.tfvars"},
		BackendConfigs: []string{"variables/backend/**/*.tfvars"},
		Exclude:        []string{"terraform.tfvars", "*.auto.tfvars"},
	}
}

// clone returns a copy that shares no slices with l.
func (l Layout) clone() Layout {
	l.VarFiles = cloneStrings(l.VarFiles)
	l.BackendConfigs = cloneStrings(l.BackendConfigs)
	l.Exclude = cloneStrings(l.Exclude)
	l.Include = cloneStrings(l.Include)
	return l
}

func cloneStrings(s []string) []string {
	if s == nil {
		return nil
	}
	return append([]string{}, s...)
}

// LoadLayout reads projectPath/.lazytf.yml over base: keys the file does not
// set keep the value from base. A missing file gives base unchanged.
func LoadLayout(projectPath string, base Layout) (Layout, error) {
	data, err := os.ReadFile(filepath.Join(projectPath, ProjectFileName))
	if os.IsNotExist(err) {
		return base.clone(), nil
	}
	if err != nil {
		return base.clone(), err
	}
	layout := base.clone()
	if err := yaml.Unmarshal(data, &layout); err != nil {
		return base.clone(), err
	}
	return layout, nil
}
