## item_003_first_run_error_messages_and_ui_polish - First run, error messages and UI polish
> From version: 0.1.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 70%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-19 15:31:18

# AI Context
- Summary: Make the first five minutes work for a new user: CLI flags, empty states, clear errors and consistent keybinds.
- Keywords: first run, cli flags, empty state, errors, keybinds, help, performance
- Use when: Changing startup, CLI flags, empty states or help text.
- Skip when: The work is about the env model or packaging.

# Problem
- A new user hits several rough edges: no `--version`/`--help`, no path argument, an uninitialised project directory triggers a scan of `~/Projects` and `~/Documents`, errors for a missing `terraform` binary or empty project list are generic, and keybinds differ between README (`L`), status bar (`shift+l`) and help.

# Scope
- In:
  - `--version` (set by GoReleaser ldflags), `--help`, `--config`, and a positional path.
  - Single-project mode when the cwd has `.tf` files, initialised or not; only scan configured search paths, never home subfolders by default.
  - Empty and error states: no projects, no environments, `terraform` not found, init failed.
  - One source of truth for keybinds shared by help modal, status bar and README.
  - Checks on 80x24 terminals and light terminal themes; discovery time on a large home dir.
  - Follow the terminal theme: replace the hardcoded Catppuccin hex colours with the terminal's 16 ANSI colour slots, and paint no backgrounds, so LazyTF matches whatever theme the terminal uses and follows theme switches.
- Out:
  - A LazyTF-specific theme system or theme picker (colours come from the terminal).
  - Mouse support.

# Acceptance criteria
- AC1: `lazytf --version` prints the release version; `lazytf ~/infra/app` opens that project.
- AC2: Running in an uninitialised project directory opens that project, not a home-wide scan.
- AC3: Without `terraform` on PATH the user sees one clear message saying what to install.
- AC4: A first-time tester plans a sample repo in under five minutes using only the README.
- AC5: No hex or 256-colour values remain in the UI styles. After switching the terminal theme (e.g. Catppuccin Macchiato to a light theme), LazyTF's colours change with it and stay readable.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: contributes to a stranger installing LazyTF and planning their own repo.
- request-AC2 -> This backlog slice. Proof: contributes to an automated, licensed, multi-channel release.
- request-AC3 -> This backlog slice. Proof: linked to roadmap milestone 0.4 - User-facing polish.

# Decision framing
- Product framing: Covered by `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture decision(s): (none yet)
- Request: `req_000_lazytf_1_0_public_release`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Decides whether a stranger keeps the tool after the first run; waits for the 0.3 model to settle.

# Notes
- Roadmap milestone: 0.4 - User-facing polish (`road_001_lazytf_to_1_0`).

# Tasks
- `task_004_first_run_error_messages_and_ui_polish`
