## req_000_lazytf_1_0_public_release - LazyTF 1.0 public release that strangers can install
> From version: 0.1.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Umbrella request to take LazyTF from a v0.1.0 TUI that fits the owner's AWS SSO + tfvars layout to a public 1.0 that strangers install from a package manager and use on their own Terraform repos.
- Keywords: lazytf, terraform, tui, 1.0, release, goreleaser, aur, homebrew, layout, auth
- Use when: Scoping or sequencing work on the path to the LazyTF 1.0 public release.
- Skip when: You need details of a single change; read the linked backlog item instead.

# Needs
- LazyTF works today for one repo layout and one cloud: per-env `*.tfvars`, `variables/backend/**/backend_<env>.tfvars`, a remote backend whose state key starts with the env name, and AWS SSO. Someone else's repo usually cannot plan at all.
- The owner's GitHub is a portfolio: the goal is a tool that strangers install and use, which needs a licence, multi-platform releases and package-manager installs.
- Order the work so correctness and tests of the current flow come first, then removing owner-specific assumptions, then user-facing polish, then distribution.

# Context
- Current state (v0.1.0, last commit 2026-07-08, one local commit not pushed): Go 1.25, Bubble Tea; init/plan/apply/destroy with streaming output, plan summary, AWS SSO login and retry; CI (build, vet, gofmt, test, golangci-lint v2) and a GoReleaser Linux-only release workflow.
- 2026-10-16: `go build ./...` passes; `go test ./...` passes, but only `internal/terraform` has tests (arg builders); `config`, `executor`, `aws` and `ui` have none.
- Gaps found while reading the code:
  - Plan, apply and destroy are refused unless the env name inferred from the backend `key`/`path`/`prefix` (first path segment) equals the selected tfvars name. A local backend or a key such as team/dev.tfstate blocks them permanently.
  - The saved plan is `/tmp/lazytf-<dir basename>.tfplan`, shared across environments and across projects with the same name.
  - A partial `config.yml` drops the default ignore patterns, and the defaults scan ~/Projects and ~/Documents when the current directory is not initialised.
  - Every `*.tfvars` (including `terraform.tfvars`, `*.auto.tfvars`) is listed as an environment.
  - No LICENSE file, although the README says MIT; no macOS builds; no demo.
- Roadmap: `road_001_lazytf_to_1_0`; product brief: `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`.

# Acceptance criteria
- AC1: A new user can install LazyTF from a package manager or `go install`, run it in their own Terraform repo (any backend, any cloud) and plan and apply the environment they selected.
- AC2: Releases are automated from a tag: licensed, built for linux/darwin amd64/arm64, and published to GitHub releases, the AUR and a Homebrew tap, with CI gating build, test and lint.
- AC3: The work is split into backlog items linked to versioned roadmap milestones.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture decision(s): (none yet)

# References
- `cmd/lazytf/main.go`
- `internal/terraform/state.go`
- `internal/terraform/plan.go`
- `internal/terraform/varfiles.go`
- `internal/terraform/backend.go`
- `internal/config/models.go`
- `internal/ui/model.go`
- `.goreleaser.yaml`
- `README.md`

# Backlog
- `item_000_current_flow_correctness_and_tests`
- `item_001_configurable_project_env_and_backend_layout`
- `item_002_cloud_agnostic_auth_and_credential_hooks`
- `item_003_first_run_error_messages_and_ui_polish`
- `item_004_release_packaging_license_goreleaser_aur_homebrew`
- `item_005_public_readme_with_demo_and_install_matrix`
