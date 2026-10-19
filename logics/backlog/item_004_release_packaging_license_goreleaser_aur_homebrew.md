## item_004_release_packaging_license_goreleaser_aur_homebrew - Release packaging: LICENSE, GoReleaser, AUR, Homebrew
> From version: 0.1.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-19 15:57:23

# AI Context
- Summary: Ship 1.0 through one tag push: LICENSE, GoReleaser builds for linux/darwin amd64/arm64 on GitHub releases, AUR lazytf-bin, a Homebrew tap, go install, and CI gating build, test and lint.
- Keywords: release, license, mit, goreleaser, aur, lazytf-bin, homebrew, tap, go install, ci, lint
- Use when: Cutting a release or changing build, CI or packaging.
- Skip when: The work is feature development.

# Problem
- There is no LICENSE file, although the README says MIT, so the code cannot legally be reused or packaged.
- `.goreleaser.yaml` builds Linux only and sets no version, and there is no package-manager channel. v0.1.0 is tagged and one local commit is not pushed.

# Scope
- In:
  - Remove the v0.1 legacy plan-file cleanup (`internal/terraform/legacyplans.go` and its startup call) before tagging 1.0 (decided 2026-10-19).
  - LICENSE (MIT proposed; owner confirms) and a matching README section.
  - GoReleaser: `goos` linux and darwin, `goarch` amd64 and arm64, `-ldflags -X main.version`, checksums, changelog grouped by commit prefix.
  - `aurs:` block publishing `lazytf-bin` (needs an AUR SSH key as a repo secret).
  - `brews:` (or `homebrew_casks:`) block publishing to `Nicolas-Rigaudy/homebrew-tap` (needs a tap repo and a token secret).
  - `go install github.com/Nicolas-Rigaudy/lazytf/cmd/lazytf@latest` verified and documented.
  - CI: keep build, vet, gofmt, test and golangci-lint; add `goreleaser check` and a snapshot build on pull requests; macOS runner for tests.
  - CHANGELOG, issue templates and SECURITY.md.
- Out:
  - Windows, deb/rpm, Nix and Scoop packages (post-1.0 if asked).
  - Code signing and notarisation for macOS (document the Gatekeeper workaround instead).

# Acceptance criteria
- AC1: Pushing tag `v1.0.0` creates a GitHub release with linux and darwin, amd64 and arm64 archives plus checksums.
- AC2: `yay -S lazytf-bin`, `brew install Nicolas-Rigaudy/tap/lazytf` and `go install ...@latest` each give a working `lazytf --version` on a clean machine.
- AC3: A pull request that breaks build, tests, lint or the GoReleaser config fails CI.
- AC4: The repository has a LICENSE file matching the README.

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
- Rationale: Distribution is the 1.0 headline; the LICENSE alone blocks any packaging.

# Notes
- Roadmap milestone: 1.0 - Public release (`road_001_lazytf_to_1_0`).
