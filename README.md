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

# With LazyTF
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
- **AWS SSO login** — auto-discovers sessions from `~/.aws/config`, opens browser automatically
- **Smart credential detection** — detects AWS auth errors and prompts to login mid-command
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

### From source

```bash
git clone https://github.com/Nicolas-Rigaudy/lazytf
cd lazytf
go build -o lazytf cmd/lazytf/main.go
sudo mv lazytf /usr/local/bin/
```

### Requirements

- Terraform installed and on `$PATH`
- AWS CLI (optional, for SSO login)
- Go 1.25+ (only for building from source)

---

## Usage

Run `lazytf` in or above a directory containing Terraform projects. It auto-detects whether you're inside a single project or a monorepo with multiple projects.

```bash
cd ~/infra
lazytf
```

### First run — how LazyTF finds your stuff

You don't configure anything to start. LazyTF discovers projects and environments by convention:

- **Single vs. multi-project** — if the current directory is already initialized (has a `.terraform/` folder), LazyTF opens it directly. Otherwise it scans your search paths and lists every project it finds in the sidebar. A "project" is any directory containing `.tf` files.
- **Environments** — inside a project, each `*.tfvars` file is an environment. LazyTF looks in the project root and in `variables/`, `env/`, and `tfvars/`. `dev.tfvars` shows up as the `dev` environment.
- **Backends** — backend configs live under `variables/backend/**/*.tfvars`. A file named `backend_dev.tfvars` is auto-matched to the `dev` environment; a plain `backend.tfvars` applies to all.

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

By default LazyTF searches `.`, `~/Projects`, and `~/Documents`. To point it at your infra directory instead, create `~/.config/lazytf/config.yml`:

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

The file is optional — sensible defaults apply if it's absent.

### Key Bindings

| Key | Action |
|-----|--------|
| `i` | Initialize selected environment |
| `p` | Run `terraform plan` |
| `a` | Apply (uses saved plan if available) |
| `d` | Destroy (with confirmation gate) |
| `L` | AWS SSO login |
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
│   ├── aws/            # AWS SSO session discovery and login
│   ├── terraform/      # Command builders: init, plan, apply, destroy
│   └── ui/             # Bubble Tea components: model, sidebar, panels, modals
```

---

## Tech Stack

- **[Go](https://go.dev/)** — fast, portable binary
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Elm-inspired MVU TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** — declarative terminal styling
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — reusable TUI components (viewport, etc.)

---

## License

MIT
