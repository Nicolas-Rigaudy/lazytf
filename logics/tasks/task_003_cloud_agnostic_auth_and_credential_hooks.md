## task_003_cloud_agnostic_auth_and_credential_hooks - Cloud-agnostic auth and credential hooks
> From version: 0.1.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-10-20 18:24:47

# AI Context
- Summary: Deliver item_002: auth providers (aws-sso, gcloud, az, custom) with login commands, error patterns and per-env variables; no AWS wording or key unless a provider is configured.
- Keywords: cloud, agnostic, auth, credential, hooks
- Use when: Changing login, credential error detection or retry-after-login.
- Skip when: The work is about layout or packaging.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_002_cloud_agnostic_auth_and_credential_hooks`

# Acceptance criteria
- AC1: With no auth config, LazyTF starts and runs without the AWS CLI installed and shows no AWS wording.
- AC2: A GCP user with expired ADC is prompted to run the gcloud login and the failed plan retries after it.
- AC3: The owner's AWS SSO flow behaves as today with the `aws-sso` preset.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_003_cloud_agnostic_auth_and_credential_hooks.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_003_cloud_agnostic_auth_and_credential_hooks.md` after implementation.

# Validation
- command: `go vet ./... && go test -race ./... && go test -tags integration ./internal/terraform` | result: passed | date: 2026-10-20
- `go vet ./... && go test -race ./...` passed locally; `go test -tags integration ./internal/terraform` passed locally (Terraform 1.15.x).
- AC1: TestNoProviderNeverMentionsLogin, TestLoginKeyOnlyWithProvider (no login key, hint, prompt or AWS wording without a provider); the aws CLI is only run from the login workflow.
- AC2: TestExpiredGCloudCredentialsPromptLoginAndRetry (gcloud error opens the prompt, confirm triggers the login, the failed plan is the pending retry).
- AC3: TestAWSProviderKeepsItsErrorsAndLabel, TestResolveLoginCommands (aws-sso keeps session selection and its error strings).
- go vet, go test -race passed on 2026-10-20: all tests ok; integration tests passed
- Finish workflow executed on 2026-10-20.
- Linked backlog/request close verification passed.

# Report
- 3dd7dfa: auth package with presets and error matching.
- e37f47a: auth config in the layout, login key, prompt and status hint driven by the provider.
- 282b554: per-environment variables for terraform and login commands.
- Behaviour change for 0.1 users: with no `auth` config there is no login key or credential prompt any more; they add `auth: {provider: aws-sso}`. Documented in the README.
- Finished on 2026-10-20.
- Linked backlog item(s): `item_002_cloud_agnostic_auth_and_credential_hooks`
- Related request(s): `req_000_lazytf_1_0_public_release`

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# AC Traceability
- request-AC1 -> This task. Proof: Contributes only (milestone 0.3 of road_001): a stranger with GCP or Azure, or no cloud login, can run LazyTF; validated with go test -race. Source: `e37f47a`
- request-AC2 -> This task. Proof: Contributes only: packaging is item_004. Source: `e37f47a`
- request-AC3 -> This task. Proof: Contributes only: linked to roadmap milestone 0.3. Source: `e37f47a`
