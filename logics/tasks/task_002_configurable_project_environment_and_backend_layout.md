## task_002_configurable_project_environment_and_backend_layout - Configurable project, environment and backend layout
> From version: 0.1.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-10-20 18:21:24

# AI Context
- Summary: Deliver item_001: *.tf.json detection, configurable var-file and backend globs (global config and per-repo .lazytf.yml), a recorded init binding, a stand-in env for projects without tfvars, and a configurable terraform binary.
- Keywords: configurable, project, environment, backend, layout
- Use when: Changing how projects, environments or backend configs are found or matched.
- Skip when: The work is about auth (item_002) or packaging.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_001_configurable_project_env_and_backend_layout`

# Acceptance criteria
- AC1: A local-backend project with no tfvars initialises, plans and applies from LazyTF.
- AC2: A GCS-backend project with `envs/dev.tfvars` and `backends/dev.hcl` works after a three-line `.lazytf.yml`.
- AC3: The owner's current layout still works with no config.
- AC4: `terraform.tfvars` and `*.auto.tfvars` are no longer listed as environments by default.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_002_configurable_project_environment_and_backend_layout.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_002_configurable_project_environment_and_backend_layout.md` after implementation.

# Validation
- command: `go vet ./... && go test -race ./... && go test -tags integration ./internal/terraform` | result: passed | date: 2026-10-20
- `go vet ./... && go test -race ./...` passed locally; `go test -tags integration ./internal/terraform` passed locally (Terraform 1.15.x).
- AC1: TestIntegrationNoVarFilesLocalBackend (init, plan, apply on a local backend without tfvars).
- AC2: TestDiscoverBackendVarFilesLayout and TestDiscoverVarFilesLayout (envs/ and backends/*.hcl via config keys); TestLoadLayout (three-line .lazytf.yml).
- AC3: the existing discovery tests run unchanged against config.DefaultLayout(); the first integration test still passes.
- AC4: TestDiscoverVarFilesLayout, "default layout skips terraform.tfvars and auto.tfvars".
- go vet, go test -race passed on 2026-10-20: 136 tests ok; integration tests with real terraform passed
- Finish workflow executed on 2026-10-20.
- Linked backlog/request close verification passed.

# Report
- 822eaed: projects written only in *.tf.json are detected.
- 14a0e78, 0efaeb5: layout config (global keys plus .lazytf.yml) and discovery through var_files / backend_configs / exclude / include / env patterns.
- 8ad28e3: init record with fingerprint, preferred to state-key inference.
- c182865, cc1275a: no-var-file stand-in env, init without -backend-config, integration test for a local backend without tfvars.
- eb5a85d: configurable terraform binary.
- 14a0e78: README.
- Layout lives in `config.Layout` (global `config.yml` keys plus `.lazytf.yml`); discovery uses a small glob matcher with `**`.
- The init record is keyed to a fingerprint of `.terraform/terraform.tfstate`, so a manual re-init invalidates it and inference is the fallback.
- A project with no backend config files treats one init as valid for every env (`BackendState.AnyEnv`).
- A project that declares no provider, module or backend gets no `.terraform/` from init; LazyTF treats such a configuration as ready (`needsInit`).
- Finished on 2026-10-20.
- Linked backlog item(s): `item_001_configurable_project_env_and_backend_layout`
- Related request(s): `req_000_lazytf_1_0_public_release`

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# AC Traceability
- request-AC1 -> This task. Proof: Contributes only (milestone 0.3 of road_001): a stranger's layouts and local backends now plan without renaming files; validated with go test -race and the integration test. Does not by itself satisfy the release-level criterion. Source: `14a0e78`
- request-AC2 -> This task. Proof: Contributes only (milestone 0.3 of road_001): layout config; packaging is item_004. Source: `14a0e78`
- request-AC3 -> This task. Proof: Contributes only: linked to roadmap milestone 0.3. Source: `14a0e78`
