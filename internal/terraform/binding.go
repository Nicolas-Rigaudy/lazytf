package terraform

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Binding records which environment (and backend config file) LazyTF last
// initialised a project with. It replaces guessing the environment from the
// state key, which fails for keys like "team/dev.tfstate".
//
// Fingerprint ties the record to the backend `terraform init` left in
// .terraform/terraform.tfstate: when someone re-inits by hand with another
// backend the fingerprint stops matching and the record is ignored.
type Binding struct {
	Env         string `json:"env"`
	BackendFile string `json:"backend_file,omitempty"`
	Fingerprint string `json:"fingerprint"`
}

// RecordInit saves the binding for a project after a successful init.
func RecordInit(projectPath, envName string, backend BackendVarFile) error {
	return recordInit(bindingCacheDir(), projectPath, envName, backend)
}

func recordInit(cacheDir, projectPath, envName string, backend BackendVarFile) error {
	b := Binding{Env: envName, Fingerprint: backendFingerprint(projectPath)}
	if backend.Name != "" {
		b.BackendFile = filepath.Join(backend.Path, backend.Name)
	}
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	path := bindingPath(cacheDir, projectPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// loadBinding returns the saved binding when it still describes the backend
// the project is initialised with.
func loadBinding(cacheDir, projectPath string) (Binding, bool) {
	data, err := os.ReadFile(bindingPath(cacheDir, projectPath))
	if err != nil {
		return Binding{}, false
	}
	var b Binding
	if err := json.Unmarshal(data, &b); err != nil || b.Env == "" {
		return Binding{}, false
	}
	if b.Fingerprint != backendFingerprint(projectPath) {
		return Binding{}, false
	}
	return b, true
}

// backendFingerprint hashes the backend section of .terraform/terraform.tfstate.
// A project with no backend block (local state) has no such file, which
// hashes to the same value every time.
func backendFingerprint(projectPath string) string {
	data, err := os.ReadFile(filepath.Join(projectPath, ".terraform", "terraform.tfstate"))
	if err != nil {
		return "none"
	}
	var tfstate struct {
		Backend json.RawMessage `json:"backend"`
	}
	if json.Unmarshal(data, &tfstate) != nil || len(tfstate.Backend) == 0 {
		return "none"
	}
	// Re-encode through a generic value so key order and whitespace do not matter.
	var generic any
	if json.Unmarshal(tfstate.Backend, &generic) != nil {
		return "none"
	}
	canonical, err := json.Marshal(generic)
	if err != nil {
		return "none"
	}
	return fmt.Sprintf("%x", sha256.Sum256(canonical))
}

func bindingCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "lazytf", "init")
}

func bindingPath(cacheDir, projectPath string) string {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		abs = filepath.Clean(projectPath)
	}
	sum := sha256.Sum256([]byte(abs))
	return filepath.Join(cacheDir, fmt.Sprintf("%s-%x.json", filepath.Base(abs), sum[:6]))
}
