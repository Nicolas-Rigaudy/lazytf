## task_001_current_flow_correctness_and_tests - Current flow correctness and tests
> From version: 0.1.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: Nicolas Rigaudy
> Indicators reviewed: 2026-10-19 15:07:08

# AI Context
- Summary: Deliver item_000 in small waves: per-env plan file, config defaults merge, -input=false, then tests for every non-UI package and a fixture integration test.
- Keywords: current, flow, correctness, tests
- Use when: Working on milestone 0.2 correctness fixes or test coverage.
- Skip when: The work changes the env/backend model (item_001) or adds UI snapshot tests.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_000_current_flow_correctness_and_tests`

# Acceptance criteria
- AC1: Planning env A, switching to env B and pressing apply never offers env A's plan.
- AC2: A `config.yml` containing only `search_paths` still ignores `.git`, `node_modules` and `.terraform`.
- AC3: A plan with a missing required variable fails immediately with Terraform's error shown in the output panel.
- AC4: `go test ./...` covers config, discovery, backend matching and env inference, and the fixture integration test runs green in CI.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_001_current_flow_correctness_and_tests.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_001_current_flow_correctness_and_tests.md` after implementation.

# Validation
- `go vet ./... && go test -race ./...` passed locally.
- `go test -tags integration -run Integration -v ./internal/terraform` passed locally (Terraform 1.15.9).
- CI run 39257292170 green: build/vet/fmt/test, integration, golangci-lint.
- command: `go vet ./... && go test -race ./... && go test -tags integration -run Integration ./internal/terraform` | result: passed | date: 2026-10-19
- Finish workflow executed on 2026-10-19.
- Linked backlog/request close verification passed.

# Report
- 3d9c901: plan file keyed by absolute project path + env under the user cache dir (0700); apply menu flags a saved plan older than the last init; post-plan summary reads the planned env's file. (AC1)
- ac7b91c: config.yml merged over defaults field by field. (AC2)
- f4aafed: -input=false on plan, apply and destroy. Checked on Terraform 1.15.9: approval is still prompted, variable prompts are not, so the stdin "yes" can no longer fill a missing variable. (AC3)
- 0db40c0: failed command and reason in the status bar; start failures no longer reported as "exit code 0".
- 4856efb: plan summary groups in stable Terraform order by module_address; replace counted as add + destroy; -/+ and +/- symbols.
- 61725db: ignore patterns such as "*/.git" never matched (base-name only); now matched on trailing segments. Without this AC2 was moot.
- 476ccf9, 1d9e379: table tests for var files, backend discovery/matching, env inference, backend detection, executor. (AC4)
- ffa13d1: integration test (init / fail-fast plan / plan / apply saved plan / stdin-confirmed destroy) on a local-backend null-provider fixture, CI job with Terraform 1.15.9 pinned by checksum. (AC4)
- Finished on 2026-10-19.
- Linked backlog item(s): `item_000_current_flow_correctness_and_tests`
- Related request(s): `req_000_lazytf_1_0_public_release`

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# AC Traceability
- request-AC1 -> This task. Proof: Contributes only (milestone 0.2 of road_001): current-flow correctness fixes 3d9c901..61725db and tests 476ccf9, 1d9e379, ffa13d1; validated with go test ./..., the integration test, and green CI run 39257292170. Does not by itself satisfy the release-level criterion. Source: `ffa13d1`
- request-AC2 -> This task. Proof: Contributes only (milestone 0.2 of road_001): current-flow correctness fixes 3d9c901..61725db and tests 476ccf9, 1d9e379, ffa13d1; validated with go test ./..., the integration test, and green CI run 39257292170. Does not by itself satisfy the release-level criterion. Source: `ffa13d1`
- request-AC3 -> This task. Proof: Contributes only (milestone 0.2 of road_001): current-flow correctness fixes 3d9c901..61725db and tests 476ccf9, 1d9e379, ffa13d1; validated with go test ./..., the integration test, and green CI run 39257292170. Does not by itself satisfy the release-level criterion. Source: `ffa13d1`
