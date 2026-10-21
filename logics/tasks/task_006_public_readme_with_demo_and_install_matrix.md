## task_006_public_readme_with_demo_and_install_matrix - Public README with demo and install matrix
> From version: 0.1.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: night
> Indicators reviewed: 2026-10-21 21:04:54

# AI Context
- Summary: Rewrite the README for strangers: screenshots of discover, init, plan, summary and apply on a sample repo (recorded by docs/demo/record.py), an install matrix with AUR and Homebrew marked as coming with 1.0, layout-neutral configuration, safety notes, comparison and contributing. The animated GIF is left as a marked slot for the owner.
- Keywords: public, readme, demo, install, matrix
- Use when: Changing the README, the demo screenshots or the sample repo under docs/demo/.
- Skip when: The work is code or packaging.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_005_public_readme_with_demo_and_install_matrix`

# Acceptance criteria
- AC1: The README shows a demo of discover, init, plan, summary and apply on a sample repo.
- AC2: Every install method in the README works as written on a clean machine.
- AC3: No section assumes AWS or the `variables/backend/` layout.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_006_public_readme_with_demo_and_install_matrix.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_006_public_readme_with_demo_and_install_matrix.md` after implementation.

# Validation
- `go test ./...` and `go vet ./...` pass; `internal/ui/layout_test.go` checks every line of the view fits widths 80, 110 and 160.
- `docs/demo/record.py` runs end to end on `docs/demo/sample` (discover, init, plan, summary, apply) with Terraform 1.15.
- README: every `#anchor` matches a heading, every relative link and image exists, code fences are balanced; rendered with Python-Markdown (4 tables, 18 code blocks, 10 images).
- `go install ...@main` gives a working `lazytf --version`; the Linux release tarball and `go install ...@latest` install v0.1.0, which has no `--version` (the README says so).

# Report
Wave 1 done (2026-10-21):
- af0161f: the header box was two columns wider than the terminal, so its right border was cut off; found while taking the screenshots.
- d3ba131: `docs/demo/sample` (two projects, `terraform_data` only, local backend, nothing cloud), `docs/demo/record.py` (pty + pyte, writes one ANSI screen per step; masks the local path as ~/infra) and six optimised PNG screenshots (about 85 KB). (AC1)
- e0108b9: README rewritten for strangers: demo, install matrix with AUR and Homebrew marked "coming with 1.0", layout-neutral examples (gcloud, envs/ + backends/*.hcl, OpenTofu), safety section, comparison with terraform-tui and web UIs, contributing. (AC1, AC3)

Not done, needs the owner:
- The animated GIF: marked slot (HTML comment) above the hero screenshot in README.md.
- AC2 on a clean machine, and the AUR and Homebrew rows once task_005 publishes them.

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
