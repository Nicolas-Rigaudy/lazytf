## task_001_current_flow_correctness_and_tests - Current flow correctness and tests
> From version: 0.1.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 20%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: Nicolas Rigaudy
> Indicators reviewed: 2026-10-19 14:45:51

# AI Context
- Summary: Deliver item_000 in small waves: per-env plan file, config defaults merge, -input=false, then tests for every non-UI package and a fixture integration test.
- Keywords: current, flow, correctness, tests
- Use when: Working on milestone 0.2 correctness fixes or test coverage.
- Skip when: The work changes the env/backend model (item_001) or adds UI snapshot tests.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_000_current_flow_correctness_and_tests`

# Acceptance criteria
- AC1: Planning env A, switching to env B and pressing apply never offers env A's plan.
- AC2: A `config.yml` containing only `search_paths` still ignores `.git`, `node_modules` and `.terraform`.
- AC3: A plan with a missing required variable fails immediately with Terraform's error shown in the output panel.
- AC4: `go test ./...` covers config, discovery, backend matching and env inference, and the fixture integration test runs green in CI.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_001_current_flow_correctness_and_tests.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_001_current_flow_correctness_and_tests.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- Wave 1: plan file keyed by absolute project path + env under the user cache dir (private 0700 dir); apply menu flags a saved plan older than the last init; post-plan summary reads the plan of the env that was planned.

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
