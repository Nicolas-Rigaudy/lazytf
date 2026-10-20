package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    options
		wantErr bool
	}{
		{name: "none", args: nil, want: options{}},
		{name: "version", args: []string{"--version"}, want: options{showVersion: true}},
		{name: "path", args: []string{"~/infra/app"}, want: options{path: "~/infra/app"}},
		{name: "config then path", args: []string{"--config", "c.yml", "app"}, want: options{configPath: "c.yml", path: "app"}},
		{name: "two paths", args: []string{"a", "b"}, wantErr: true},
		{name: "unknown flag", args: []string{"--nope"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args, &bytes.Buffer{})
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseArgs() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestHelpDescribesFlagsAndPath(t *testing.T) {
	var out bytes.Buffer
	_, err := parseArgs([]string{"--help"}, &out)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("--help error = %v, want flag.ErrHelp", err)
	}
	for _, want := range []string{"lazytf [flags] [path]", "--version", "--config"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("help output lacks %q:\n%s", want, out.String())
		}
	}
}

func writeTree(t *testing.T, root string, files ...string) {
	t.Helper()
	for _, f := range files {
		p := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// isolate keeps resolveStart away from the real user config.
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
}

func TestResolveStartUninitialisedProjectOpensThatProject(t *testing.T) {
	isolate(t)
	cwd := t.TempDir()
	writeTree(t, cwd, "main.tf")
	st, err := resolveStart(options{}, cwd)
	if err != nil {
		t.Fatal(err)
	}
	if st.mode != terraform.ModeSingleProject || len(st.projects) != 1 || st.projects[0].Path != cwd {
		t.Errorf("start = %+v, want single project at %s", st, cwd)
	}
}

func TestResolveStartNeverScansHomeFolders(t *testing.T) {
	isolate(t)
	home, _ := os.UserHomeDir()
	writeTree(t, home, "Projects/elsewhere/main.tf", "Documents/other/main.tf")
	cwd := t.TempDir()
	writeTree(t, cwd, "stack/main.tf")
	st, err := resolveStart(options{}, cwd)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range st.projects {
		names = append(names, p.Name)
	}
	if !slices.Equal(names, []string{"stack"}) {
		t.Errorf("projects = %v, want only the working directory's", names)
	}
}

func TestResolveStartPathArgument(t *testing.T) {
	isolate(t)
	cwd := t.TempDir()
	writeTree(t, cwd, "infra/app/main.tf", "infra/other/main.tf.json", "plain/readme.md")

	st, err := resolveStart(options{path: "infra/app"}, cwd)
	if err != nil || st.mode != terraform.ModeSingleProject || st.projects[0].Path != filepath.Join(cwd, "infra", "app") {
		t.Errorf("project path: start = %+v, err = %v", st, err)
	}

	st, err = resolveStart(options{path: filepath.Join(cwd, "infra")}, cwd)
	if err != nil || st.mode != terraform.ModeMultiProject || len(st.projects) != 2 {
		t.Errorf("directory of projects: start = %+v, err = %v; want two projects", st, err)
	}

	if _, err := resolveStart(options{path: "missing"}, cwd); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("missing path error = %v", err)
	}
	if _, err := resolveStart(options{path: "plain"}, cwd); err == nil || !strings.Contains(err.Error(), "no Terraform projects found") {
		t.Errorf("empty directory error = %v", err)
	}
}

func TestResolveStartNoProjectsExplainsWhatToDo(t *testing.T) {
	isolate(t)
	_, err := resolveStart(options{}, t.TempDir())
	if err == nil {
		t.Fatal("empty search paths gave no error")
	}
	for _, want := range []string{"no Terraform projects found", ".tf.json", "search_paths"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q lacks %q", err, want)
		}
	}
}

func TestResolveStartConfigFlag(t *testing.T) {
	isolate(t)
	cwd := t.TempDir()
	writeTree(t, cwd, "mono/a/main.tf")
	cfgPath := filepath.Join(t.TempDir(), "c.yml")
	if err := os.WriteFile(cfgPath, []byte("search_paths: ["+filepath.Join(cwd, "mono")+"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, err := resolveStart(options{configPath: cfgPath}, t.TempDir())
	if err != nil || len(st.projects) != 1 || st.projects[0].Name != "a" {
		t.Errorf("start = %+v, err = %v; want the project from the config's search path", st, err)
	}
	if _, err := resolveStart(options{configPath: filepath.Join(cwd, "nope.yml")}, cwd); err == nil {
		t.Error("a --config file that does not exist gave no error")
	}
}

func TestCheckBinary(t *testing.T) {
	bin := t.TempDir()
	fake := filepath.Join(bin, "terraform")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	if err := checkBinary(""); err != nil {
		t.Errorf("terraform on PATH: %v", err)
	}
	if err := checkBinary("tofu"); err == nil || !strings.Contains(err.Error(), `"tofu"`) {
		t.Errorf("missing tofu error = %v", err)
	}
	t.Setenv("PATH", t.TempDir())
	err := checkBinary("")
	if err == nil || !strings.Contains(err.Error(), "Install Terraform") {
		t.Errorf("missing terraform error = %v, want what to install", err)
	}
}

func TestVersionString(t *testing.T) {
	old := version
	defer func() { version = old }()
	version = "1.2.3"
	if got := versionString(); got != "1.2.3" {
		t.Errorf("versionString() = %q, want the ldflags version", got)
	}
}
