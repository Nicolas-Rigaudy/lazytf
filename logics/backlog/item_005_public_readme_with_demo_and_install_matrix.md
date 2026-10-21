## item_005_public_readme_with_demo_and_install_matrix - Public README with demo and install matrix
> From version: 0.1.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 10%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-21 21:04:55

# AI Context
- Summary: Rewrite the README for strangers: what it is, a demo GIF or asciinema cast, an install matrix, configuration for common layouts, safety notes and a comparison with other Terraform TUIs.
- Keywords: readme, demo, gif, asciinema, vhs, install, configuration, docs
- Use when: Changing user-facing documentation.
- Skip when: The work is code.

# Problem
- The README is solid but written around the owner's layout and AWS SSO, says "Screenshots / GIF coming soon", documents only Linux tarball and source installs, and claims an MIT licence without a LICENSE file.

# Scope
- In:
  - Demo recorded with a script (vhs tape or asciinema) committed to the repo so it can be re-recorded each release.
  - Install matrix: AUR, Homebrew, `go install`, release tarball, from source.
  - Configuration reference with the three example layouts from item_001 and the auth presets from item_002.
  - "How it keeps you safe" section: confirm gates, saved-plan rules, no credential storage.
  - Short comparison with terraform-tui and terraform-ui; contributing section.
- Out:
  - A documentation website.

# Acceptance criteria
- AC1: The README shows a demo of discover, init, plan, summary and apply on a sample repo.
- AC2: Every install method in the README works as written on a clean machine.
- AC3: No section assumes AWS or the `variables/backend/` layout.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: contributes to a stranger installing LazyTF and planning their own repo.
- request-AC2 -> This backlog slice. Proof: contributes to an automated, licensed, multi-channel release.
- request-AC3 -> This backlog slice. Proof: linked to roadmap milestone 1.0 - Public release.

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
- Rationale: The README is the portfolio's front page; it decides whether a visitor installs.

# Notes
- Roadmap milestone: 1.0 - Public release (`road_001_lazytf_to_1_0`).

# Tasks
- `task_006_public_readme_with_demo_and_install_matrix`
