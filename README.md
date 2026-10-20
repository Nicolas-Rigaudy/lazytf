# LazyTF

**A terminal UI for Terraform — plan, apply, and destroy without leaving your keyboard.**

LazyTF brings lazygit-style ergonomics to Terraform workflows. Browse your environments, run commands with a single keypress, and watch real-time streaming output — all from one panel.

---

## Why LazyTF?

Running Terraform day-to-day means repeating the same sequence of long commands, switching between your editor and terminal, and carefully typing `yes` for every apply. It's slow and error-prone.

LazyTF replaces that with a keyboard-driven TUI:

```bash
# Before LazyTF
aws sso login --sso-session my-session
terraform init -backend-config=variables/backend/local/backend_dev.tfvars -reconfigure
terraform plan -var-file=variables/dev.tfvars -out=dev.tfplan
terraform apply dev.tfplan

# With LazyTF (auth provider aws-sso)
L → select session → browser opens
i → select env → select backend → confirm
p → streams plan output
a → uses saved plan → confirm with y
```

---

## Features

- **Multi-environment navigation** — browse projects and `.tfvars` environments in a sidebar
- **Backend auto-discovery** — finds and matches backend configs to the selected environment
- **Streaming output** — real-time line-by-line output from every terraform command
- **Plan summary view** — structured, navigable diff grouped by module (toggle with `v`)
- **Apply from saved plan** — detects existing plan files and shows their age before applying
- **Destroy with safety gate** — confirmation modal before streaming destroy output
- **Optional login** — AWS SSO (sessions found in `~/.aws/config`), `gcloud`, `az` or your own command; see [Login](#login-optional)
- **Smart credential detection** — recognises expired-credential errors for the configured provider and prompts to log in mid-command
- **Command retry** — re-runs the last command automatically after a successful login
- **Collapsible sidebar** — more space for output when you need it (`b` to toggle)
- **Catppuccin theme** — easy on the eyes

---

## Demo

> *Screenshots / GIF coming soon*

---

## Installation

### Linux (prebuilt binary)

No Go toolchain needed. Downloads the latest release and installs to `/usr/local/bin`:

```bash
curl -sL https://github.com/Nicolas-Rigaudy/lazytf/releases/latest/download/lazytf_Linux_amd64.tar.gz | tar xz
sudo mv lazytf /usr/local/bin/
```

On ARM (e.g. a Graviton box or ARM laptop), swap `amd64` for `arm64`. Or grab a binary from the [Releases page](https://github.com/Nicolas-Rigaudy/lazytf/releases).

### macOS (prebuilt binary)

Archives for Intel (`amd64`) and Apple Silicon (`arm64`) are published from the first release after v0.1.0, named `lazytf_Darwin_amd64.tar.gz` and `lazytf_Darwin_arm64.tar.gz`. They are not code-signed, so Gatekeeper blocks the first run of a downloaded binary. Clear the quarantine flag once:

```bash
tar xzf lazytf_Darwin_arm64.tar.gz
xattr -d com.apple.quarantine lazytf
sudo mv lazytf /usr/local/bin/
```

Downloading with `curl` (as in the Linux example) does not set the flag.

### With `go install`

```bash
go install github.com/Nicolas-Rigaudy/lazytf/cmd/lazytf@latest
```

This puts `lazytf` in `$(go env GOPATH)/bin`; `lazytf --version` reports the module version.

### From source

```bash
git clone https://github.com/Nicolas-Rigaudy/lazytf
cd lazytf
go build -o lazytf cmd/lazytf/main.go
sudo mv lazytf /usr/local/bin/
```

### Requirements

- Terraform installed and on `$PATH`
- A login CLI such as `aws`, `gcloud` or `az` only if you configure that provider
- Go 1.25+ (only for building from source)

---

## Usage

Run `lazytf` in or above a directory containing Terraform projects. It auto-detects whether you're inside a single project or a monorepo with multiple projects.

```bash
cd ~/infra
lazytf
```

```bash
lazytf                 # the current directory: one project, or a scan for projects
lazytf ~/infra/app     # open that project
lazytf ~/infra         # scan that directory for projects
lazytf --config ./lazytf.yml
lazytf --version       # also: --help
```

If `terraform` is not on your `PATH`, LazyTF says so and exits before drawing anything (set `terraform: tofu`, or a full path, in the config to use another binary).

### First run — how LazyTF finds your stuff

You don't configure anything to start. LazyTF discovers projects and environments by convention:

- **Single vs. multi-project** — if the current directory holds Terraform files (or already has a `.terraform/` folder), LazyTF opens it directly, initialised or not. Otherwise it scans your search paths (by default just the current directory) and lists every project it finds in the sidebar. It never scans folders you did not name. A "project" is any directory containing `.tf` or `.tf.json` files (the JSON syntax that CDKTF and other generators write).
- **Environments** — inside a project, each `*.tfvars` file is an environment. LazyTF looks in the project root and in `variables/`, `env/`, and `tfvars/`. `dev.tfvars` shows up as the `dev` environment. `terraform.tfvars` and `*.auto.tfvars` are not environments: Terraform loads them by itself. A project with no var files gets one environment, `(no var-file)`, so it can still be planned.
- **Backends** — backend configs live under `variables/backend/**/*.tfvars`. A file named `backend_dev.tfvars` is auto-matched to the `dev` environment; a plain `backend.tfvars` applies to all. A project with no backend config files (local state, or a backend written in the `.tf` files) is initialised without `-backend-config`.
- **Which environment is initialised** — LazyTF remembers the environment and backend file of its last `init` (in `~/.cache/lazytf/init/`) and checks it against the backend Terraform recorded in `.terraform/terraform.tfstate`. If you re-init by hand, it falls back to guessing from the state key.

So a typical layout LazyTF understands out of the box:

```
my-project/
├── main.tf
├── dev.tfvars                       # → "dev" environment
├── prod.tfvars                      # → "prod" environment
└── variables/
    └── backend/
        ├── backend_dev.tfvars       # → backend for "dev"
        └── backend_prod.tfvars      # → backend for "prod"
```

### Configuration (optional)

By default LazyTF searches only the current directory. To point it at your infra directories instead, create `~/.config/lazytf/config.yml`:

```yaml
search_paths:
  - ~/infra
  - ~/work/terraform
ignore_patterns:
  - "*/node_modules"
  - "*/.git"
  - "*/vendor"
  - "*/.terraform"
```

The file is optional — sensible defaults apply if it's absent. Any key you leave out keeps its default, so a file with only `search_paths` still uses the ignore patterns above.

### Layout (optional)

If your repo does not follow the layout above, say where things are. Put a `.lazytf.yml` in the project root, or the same keys in `~/.config/lazytf/config.yml` to set them for every project (a key in `.lazytf.yml` wins). Globs are relative to the project root; `**` spans directories.

| Key | Meaning | Default |
|-----|---------|---------|
| `var_files` | Globs whose matches are environments | `*.tfvars`, `variables/*.tfvars`, `env/*.tfvars`, `tfvars/*.tfvars` |
| `backend_configs` | Globs of files passed as `-backend-config` | `variables/backend/**/*.tfvars` |
| `exclude` | Var files that are not environments (a glob without `/` matches the file name) | `terraform.tfvars`, `*.auto.tfvars` |
| `include` | Var files to list even when `exclude` matches | none |
| `env_pattern` | Regex with one capture group, applied to the var file name without extension | the whole name |
| `backend_env_pattern` | Same, for backend config files | `backend_<env>`; a bare `backend` is shared |
| `terraform` | Executable to run, e.g. `tofu` or a full path | `terraform` |

A key you leave out keeps its default; `exclude: []` clears one. Three layouts:

```yaml
# 1. GCS backend, envs/dev.tfvars + backends/dev.hcl
var_files: ["envs/*.tfvars"]
backend_configs: ["backends/*.hcl"]
```

```yaml
# 2. Local backend, one var file per environment: nothing to configure.
#    (No backend config files means init runs without -backend-config.)
```

```yaml
# 3. Var files named vars-<env>.tfvars next to the code, OpenTofu
var_files: ["vars-*.tfvars"]
env_pattern: '^vars-(.+)$'
terraform: tofu
```

### Login (optional)

LazyTF does nothing about credentials unless you name an auth provider, and then it only does two things: `L` runs the login command, and when a command fails with that provider's expired-credentials error it offers to log in and re-runs the failed command afterwards. Set it in `~/.config/lazytf/config.yml` or per repo in `.lazytf.yml`:

```yaml
auth:
  provider: aws-sso   # or gcloud, az, custom
```

| Provider | Login command | Recognised errors (examples) |
|----------|---------------|------------------------------|
| `aws-sso` | `aws sso login --sso-session <name>`, session picked from `~/.aws/config` | `No valid credential sources found`, `Unable to locate credentials`, `ExpiredToken` |
| `gcloud` | `gcloud auth application-default login` | `could not find default credentials`, `Reauthentication required` |
| `az` | `az login` | `Please run 'az login'`, `AADSTS70043` |
| `custom` | your `login` list | your `error_patterns` |

```yaml
auth:
  provider: custom
  login: ["my-sso-tool", "login", "--browser"]
  error_patterns: ["session expired"]   # matched without regard to case; also extends the presets
```

To use another account or profile per environment, set variables under `auth.env`; they are added to every terraform command and to the login command LazyTF runs for that environment:

```yaml
auth:
  provider: aws-sso
  env:
    dev:  { AWS_PROFILE: dev-profile }
    prod: { AWS_PROFILE: prod-profile }
```

`login` overrides a preset's command too. Upgrading from 0.1: the AWS SSO login used to be on by default; add the `aws-sso` provider above to keep it.

### Key Bindings

| Key | Action |
|-----|--------|
| `i` | Initialize selected environment |
| `p` | Run `terraform plan` |
| `a` | Apply (uses saved plan if available) |
| `d` | Destroy (with confirmation gate) |
| `L` | Log in with the configured auth provider (hidden when none is set) |
| `v` | Toggle plan summary / raw output |
| `b` | Collapse / expand sidebar |
| `j/k` or `↑/↓` | Navigate list |
| `h/l` or `←/→` | Scroll output horizontally |
| `Tab` | Switch panel focus |
| `Esc / Backspace` | Go back |
| `q` | Quit |

### Plan Summary View

After running `terraform plan`, press `v` to switch from raw streaming output to a structured summary:

```
┌─ Plan Summary — 3 to add, 2 to change, 1 to destroy ─────┐
│ ▶ module.metabase (2+, 0~, 1-)                            │
│ ▼ module.rds (0+, 1~, 0-)                                 │
│   ▶ ~ aws_db_instance.main                                │
│ ▶ (root) (1+, 1~, 0-)                                     │
└───────────────────────────────────────────────────────────┘
```

Navigate with `ctrl+j/k`, jump between modules with `alt+j/k`, expand resources with `Enter`.

---

## Project Structure

```
lazytf/
├── cmd/lazytf/         # Entry point
├── internal/
│   ├── config/         # Config loading (YAML, XDG-compliant)
│   ├── executor/       # Generic streaming command executor
│   ├── auth/           # Login providers: aws-sso, gcloud, az, custom
│   ├── aws/            # AWS SSO session discovery and login
│   ├── terraform/      # Command builders: init, plan, apply, destroy
│   └── ui/             # Bubble Tea components: model, sidebar, panels, modals
```

### Running the tests

```bash
go test ./...                                                       # unit tests
go test -tags integration -run Integration -v ./internal/terraform  # real init/plan/apply, needs terraform on PATH
```

The integration test copies `internal/terraform/testdata/fixture` (local backend, `null` provider) to a temp dir, so it never touches your state.

---

## Tech Stack

- **[Go](https://go.dev/)** — fast, portable binary
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Elm-inspired MVU TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** — declarative terminal styling
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — reusable TUI components (viewport, etc.)

---

## License

MIT
