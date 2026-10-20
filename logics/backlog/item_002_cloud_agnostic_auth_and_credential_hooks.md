## item_002_cloud_agnostic_auth_and_credential_hooks - Cloud-agnostic auth and credential hooks
> From version: 0.1.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-20 18:24:48

# AI Context
- Summary: Turn the hardwired AWS SSO login and AWS error detection into optional, pluggable auth providers (AWS SSO, gcloud, az, custom command).
- Keywords: auth, aws sso, gcloud, azure, credentials, login, provider, hook
- Use when: Changing login, credential error detection or retry-after-login.
- Skip when: The work is about layout or packaging.

# Problem
- Auth is AWS-only: the `L` key runs `aws sso login --sso-session`, `IsLoggedIn` calls `aws sts get-caller-identity`, and credential errors are detected only from two AWS strings. GCP and Azure users get no help, and the AWS wording appears in the UI and status bar for everyone.

# Scope
- In:
  - An `auth` config section: provider (`aws-sso`, `gcloud`, `az`, `custom`, or none), login command, and error patterns that trigger the login prompt.
  - Built-in presets for AWS SSO (current behaviour), `gcloud auth application-default login` and `az login`.
  - Retry-after-login works for any provider; the login key and labels are hidden when no provider is configured.
  - Optional per-env auth (for example an AWS profile per env) set in `.lazytf.yml`.
- Out:
  - Storing or caching credentials.
  - Vault or cloud secret managers.

# Acceptance criteria
- AC1: With no auth config, LazyTF starts and runs without the AWS CLI installed and shows no AWS wording.
- AC2: A GCP user with expired ADC is prompted to run the gcloud login and the failed plan retries after it.
- AC3: The owner's AWS SSO flow behaves as today with the `aws-sso` preset.

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
- Primary task(s): `task_003_cloud_agnostic_auth_and_credential_hooks`

# Priority
- Priority: Medium
- Rationale: Needed for "any cloud" in 1.0, but it only matters once item_001 lets other repos plan.

# Notes
- Roadmap milestone: 0.3 - Any layout, any cloud (`road_001_lazytf_to_1_0`).
- Task `task_003_cloud_agnostic_auth_and_credential_hooks` was finished via `logics-manager flow finish task` on 2026-10-20.

# Tasks
- `task_003_cloud_agnostic_auth_and_credential_hooks`
