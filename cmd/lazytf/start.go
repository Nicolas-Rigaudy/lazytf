package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
)

// version is set at release time by GoReleaser (-X main.version=...).
var version = "dev"

// versionString is the release version, or for `go install` builds the module
// version recorded in the binary.
func versionString() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return version
}

type options struct {
	showVersion bool
	configPath  string
	path        string
}

const usage = `lazytf - a keyboard-driven Terraform cockpit

Usage:
  lazytf [flags] [path]

With no path, LazyTF opens the current directory when it holds Terraform
files, and otherwise lists the projects found under the search paths of the
config file. A path to a project opens it; any other directory is scanned for
projects.

Flags:
  --config file   read this config file instead of ~/.config/lazytf/config.yml
  --version       print the version and exit
  --help          show this help
`

// parseArgs reads the command line. flag.ErrHelp is returned for -h/--help
// after the usage has been written to out.
func parseArgs(args []string, out io.Writer) (options, error) {
	var opts options
	fs := flag.NewFlagSet("lazytf", flag.ContinueOnError)
	fs.SetOutput(out)
	fs.Usage = func() { _, _ = io.WriteString(out, usage) }
	fs.BoolVar(&opts.showVersion, "version", false, "print the version and exit")
	fs.StringVar(&opts.configPath, "config", "", "read this config file instead of ~/.config/lazytf/config.yml")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	switch fs.NArg() {
	case 0:
	case 1:
		opts.path = fs.Arg(0)
	default:
		return opts, fmt.Errorf("expected at most one path, got %d", fs.NArg())
	}
	return opts, nil
}

// start is what LazyTF opens with.
type start struct {
	projects []terraform.Project
	mode     terraform.Mode
	cfg      config.Config
}

// resolveStart decides what to open from the options and the working
// directory. Its errors are meant to be shown to the user as they are.
func resolveStart(opts options, cwd string) (start, error) {
	var (
		cfg config.Config
		err error
	)
	if opts.configPath != "" {
		cfg, err = config.LoadFile(opts.configPath)
	} else {
		cfg, err = config.Load()
	}
	if err != nil {
		return start{}, fmt.Errorf("cannot load the config: %w", err)
	}

	dir := cwd
	explicit := opts.path != ""
	if explicit {
		if dir, err = expandArgPath(opts.path, cwd); err != nil {
			return start{}, err
		}
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			return start{}, fmt.Errorf("%s is not a directory", opts.path)
		}
	}

	mode, project, err := terraform.DetermineMode(dir)
	if err != nil {
		return start{}, fmt.Errorf("cannot read %s: %w", dir, err)
	}
	if mode == terraform.ModeSingleProject {
		return start{projects: []terraform.Project{*project}, mode: mode, cfg: cfg}, nil
	}

	if explicit {
		cfg.SearchPaths = []string{dir} // scan the directory that was asked for
	}
	for i, p := range cfg.SearchPaths {
		if p, err = expandArgPath(p, cwd); err != nil {
			return start{}, err
		}
		cfg.SearchPaths[i] = p // relative search paths are relative to the working directory
	}
	projects, err := terraform.DiscoverProjects(cfg)
	if err != nil {
		return start{}, fmt.Errorf("cannot scan for projects: %w", err)
	}
	if len(projects) == 0 {
		return start{}, noProjectsError(cfg, explicit)
	}
	return start{projects: projects, mode: mode, cfg: cfg}, nil
}

func expandArgPath(path, cwd string) (string, error) {
	path, err := terraform.ExpandPath(path)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	return filepath.Clean(path), nil
}

func noProjectsError(cfg config.Config, explicit bool) error {
	var b strings.Builder
	b.WriteString("no Terraform projects found in: " + strings.Join(cfg.SearchPaths, ", ") + "\n")
	b.WriteString("A project is a directory with .tf or .tf.json files\n")
	if explicit {
		b.WriteString("Pass the path of a project, or a directory that contains some")
	} else {
		b.WriteString("Run lazytf inside a project, pass a path, or list directories to scan as search_paths in the config file")
	}
	return errors.New(b.String())
}

// checkBinary fails with a message saying what to install when the
// terraform executable cannot be found.
func checkBinary(configured string) error {
	bin := terraform.Binary(configured)
	if _, err := exec.LookPath(bin); err != nil {
		if configured == "" {
			return errors.New("terraform was not found on your PATH\nInstall Terraform (https://developer.hashicorp.com/terraform/install) or OpenTofu, or set `terraform: tofu` (or a full path) in the config file")
		}
		return fmt.Errorf("%q was not found\nCheck the `terraform:` setting in your config", bin)
	}
	return nil
}
