## task_005_release_packaging_license_goreleaser_aur_homebrew - Release packaging: LICENSE, GoReleaser, AUR, Homebrew
> From version: 0.1.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 50%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Deliver item_004 in waves. Wave 1 (no secrets, nothing published): GoReleaser darwin builds and version stamp, CI checks, install docs, CHANGELOG, issue templates. Left: AUR and Homebrew publishing, SECURITY.md, legacy plan cleanup, the tag.
- Keywords: release, packaging, license, goreleaser, aur, homebrew
- Use when: Cutting a release or changing build, CI or packaging.
- Skip when: The work is feature development.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_004_release_packaging_license_goreleaser_aur_homebrew`

# Acceptance criteria
- AC1: Pushing tag `v1.0.0` creates a GitHub release with linux and darwin, amd64 and arm64 archives plus checksums.
- AC2: `yay -S lazytf-bin`, `brew install Nicolas-Rigaudy/tap/lazytf` and `go install ...@latest` each give a working `lazytf --version` on a clean machine.
- AC3: A pull request that breaks build, tests, lint or the GoReleaser config fails CI.
- AC4: The repository has a LICENSE file matching the README.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_005_release_packaging_license_goreleaser_aur_homebrew.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_005_release_packaging_license_goreleaser_aur_homebrew.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
Wave 1 done (2026-10-20):
- 5f73a83: GoReleaser builds linux and darwin on amd64 and arm64, stamps `main.version`, ships LICENSE and README in the archives, groups the changelog by commit prefix. Checked locally with a pinned-free `goreleaser check` and `release --snapshot`: four archives plus checksums.txt, and the linux binary prints `lazytf --version`. (AC1 config side)
- 3f0ced1: CI job running `goreleaser check` and a snapshot build, and a macOS test job. Not run yet on GitHub. (AC3)
- 8da91ce: README install section: `go install`, macOS archives and the Gatekeeper workaround.
- d47d31f: CHANGELOG.md (Unreleased) and issue templates.
- LICENSE was already present and matches the README (AC4).

Not done, needs the owner:
- `aurs:` and `brews:` blocks: they need an AUR SSH key, a tap repo and a token as repo secrets; without them a tagged release would fail. (AC2)
- SECURITY.md: needs a private contact channel chosen by the owner.
- Removing the v0.1 legacy plan cleanup (decided 2026-10-19): deletes code files, left for the owner to confirm; do it just before the 1.0 tag.
- The v1.0.0 tag itself.

# Links
- Request: `req_000_lazytf_1_0_public_release`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
