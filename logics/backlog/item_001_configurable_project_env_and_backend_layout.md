## item_001_configurable_project_env_and_backend_layout - Configurable project, environment and backend layout
> From version: 0.1.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Replace the hardcoded tfvars/backend layout and the state-key env inference with configurable globs and an explicit env/backend binding, so any repo can plan and apply.
- Keywords: layout, config, tfvars, backend-config, environment, local backend, binding, .lazytf.yml
- Use when: Changing how projects, environments or backend configs are found or matched.
- Skip when: The work is about auth or packaging.

# Problem
- LazyTF only understands the owner's layout:
  - Var files: every `*.tfvars` in `.`, `variables/`, `env/`, `tfvars/` is an environment, including `terraform.tfvars`, `*.auto.tfvars` and shared files.
  - Backends: only `variables/backend/**/*.tfvars`, matched by a `backend_<env>` prefix.
  - Env detection: `inferEnvFromBackendConfig` takes the first meaningful segment of the backend `key`/`path`/`prefix`, and plan/apply/destroy are refused unless it equals the selected env. Local backends (no `.terraform/terraform.tfstate`) and keys like `team/dev.tfstate` or `prod-network.tfstate` are blocked permanently.
  - A project with no tfvars has no environment and cannot be planned at all.

# Scope
- In:
  - Optional per-repo `.lazytf.yml` plus global config: var-file globs, backend-config globs, an env-name pattern, and files to always include or never treat as environments.
  - Record which env and backend file LazyTF last initialised (per project, in the cache dir) and use that instead of inferring from the state key; keep inference only as a hint.
  - A default "(no var-file)" environment so projects with no tfvars and local backends can plan.
  - Configurable `terraform` binary path (enables `tofu` by name; full OpenTofu support is 1.1).
  - Docs of the config keys with three example layouts.
- Out:
  - Terraform workspaces as environments (1.1).
  - Terragrunt (1.3).

# Acceptance criteria
- AC1: A local-backend project with no tfvars initialises, plans and applies from LazyTF.
- AC2: A GCS-backend project with `envs/dev.tfvars` and `backends/dev.hcl` works after a three-line `.lazytf.yml`.
- AC3: The owner's current layout still works with no config.
- AC4: `terraform.tfvars` and `*.auto.tfvars` are no longer listed as environments by default.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: contributes to a stranger installing LazyTF and planning their own repo.
- request-AC2 -> This backlog slice. Proof: contributes to an automated, licensed, multi-channel release.
- request-AC3 -> This backlog slice. Proof: linked to roadmap milestone 0.3 - Any layout, any cloud.

# Decision framing
- Product framing: Covered by `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture decision(s): (none yet)
- Request: `req_000_lazytf_1_0_public_release`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Without it most strangers' repos cannot plan at all; this is the real 1.0 blocker.

# Notes
- Roadmap milestone: 0.3 - Any layout, any cloud (`road_001_lazytf_to_1_0`).
