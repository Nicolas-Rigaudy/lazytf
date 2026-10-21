## task_004_first_run_error_messages_and_ui_polish - First run, error messages and UI polish
> From version: 0.1.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 50%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Deliver item_003 in waves. Wave 1 (non-visual): CLI flags and path argument, project detection, start-up errors, one keybind list, scan speed. Waves left need owner choices: ANSI colour mapping, light/small terminal checks, README timing test.
- Keywords: first, run, error, messages, polish
- Use when: Changing startup, CLI flags, empty states or help text.
- Skip when: The work is about the env model or packaging.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_003_first_run_error_messages_and_ui_polish`

# Acceptance criteria
- AC1: `lazytf --version` prints the release version; `lazytf ~/infra/app` opens that project.
- AC2: Running in an uninitialised project directory opens that project, not a home-wide scan.
- AC3: Without `terraform` on PATH the user sees one clear message saying what to install.
- AC4: A first-time tester plans a sample repo in under five minutes using only the README.
- AC5: No hex or 256-colour values remain in the UI styles. After switching the terminal theme (e.g. Catppuccin Macchiato to a light theme), LazyTF's colours change with it and stay readable.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_004_first_run_error_messages_and_ui_polish.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_004_first_run_error_messages_and_ui_polish.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
Wave 1 done (2026-10-20):
- b5335d6: --version (ldflags, or the module version for go install), --help, --config, a path argument; any directory with .tf/.tf.json opens as a project, initialised or not; default search path is only "."; clear errors for no projects and a missing terraform binary, before the TUI starts. (AC1, AC2, AC3)
- 69a4a38: status bar hints are generated from the help key list ("L" everywhere); the unbound validate/format entries are gone; "v" is listed as toggle view.
- 2e610d1: one ReadDir per directory when scanning; 40k-directory synthetic tree 1.05s -> 0.27s.

Not done, needs the owner:
- AC5 colours: mapping the 16 ANSI slots to UI roles and dropping the backgrounds is a visual choice.
- Light-terminal and 80x24 checks, and AC4 (a tester plans a sample repo in under five minutes from the README): manual.
- Empty states inside the TUI beyond the start-up messages, and the plan review from terraform show -json (road_001 0.4): UX design.

Wave 2, branch night/ansi-theme, not merged (2026-10-21): LAZYTF_STYLE=ansi-a (lazygit-like) or ansi-b (blue accent) switches the colours to the terminal's 16 ANSI slots with no painted backgrounds; unset keeps Catppuccin. Colour roles live in internal/ui/theme/roles.go. Renders of both, dark and light, at 120x35 and 80x24 are waiting for the owner's pick; TestRenderScreens (LAZYTF_RENDER_DIR) regenerates them.
- Plan rows in style B are bold blue without the extra marker, because the rows already start with a fold arrow.

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
