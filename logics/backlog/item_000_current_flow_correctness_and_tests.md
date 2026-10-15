## item_000_current_flow_correctness_and_tests - Current flow correctness and tests
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
- Summary: Fix the correctness bugs in the existing init/plan/apply/destroy flow and add tests for every non-UI package before generalising anything.
- Keywords: correctness, tests, plan file, config defaults, env detection, integration test, terraform
- Use when: Fixing bugs in the current flow or adding tests.
- Skip when: The work is about new layouts, auth providers or packaging.

# Problem
- The current flow has bugs that would hit any user and only one package has tests:
  - The saved plan path `/tmp/lazytf-<basename>.tfplan` is shared across environments and across projects with the same directory name, so apply can offer the wrong plan.
  - `config.Load` replaces defaults wholesale: a `config.yml` with only `search_paths` loses the ignore patterns and walks `.git`, `node_modules` and `.terraform`.
  - Plan runs without `-input=false`; a missing variable fails late with a confusing message because stdin is not wired.
  - `go test ./...` only covers the arg builders in `internal/terraform`; discovery, var-file and backend matching, env inference, config and the executor are untested.

# Scope
- In:
  - Plan file keyed by absolute project path and env name, stored under the user cache dir; warn when a saved plan belongs to another env or is older than the last init.
  - Merge user config over defaults field by field.
  - `-input=false` on plan, and on apply/destroy when no stdin is wired; surface exit codes and errors in the status bar.
  - Table-driven tests for `DiscoverVarFiles`, `DiscoverBackendVarFiles`, `MatchBackendsForEnv`, `inferEnvFromBackendConfig`, `DiscoverProjects`, `config.Load` and plan-summary grouping.
  - One integration test in CI: init/plan/apply against a fixture project with the local backend and the `null` provider, using a pinned `terraform` binary.
- Out:
  - Changing the env/backend model itself (item_001).
  - UI snapshot tests.

# Acceptance criteria
- AC1: Planning env A, switching to env B and pressing apply never offers env A's plan.
- AC2: A `config.yml` containing only `search_paths` still ignores `.git`, `node_modules` and `.terraform`.
- AC3: A plan with a missing required variable fails immediately with Terraform's error shown in the output panel.
- AC4: `go test ./...` covers config, discovery, backend matching and env inference, and the fixture integration test runs green in CI.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: contributes to a stranger installing LazyTF and planning their own repo.
- request-AC2 -> This backlog slice. Proof: contributes to an automated, licensed, multi-channel release.
- request-AC3 -> This backlog slice. Proof: linked to roadmap milestone 0.2 - Current flow correctness and tests.

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
- Rationale: Apply on the wrong plan is the worst failure a Terraform tool can have; this gates every later milestone.

# Notes
- Roadmap milestone: 0.2 - Current flow correctness and tests (`road_001_lazytf_to_1_0`).
