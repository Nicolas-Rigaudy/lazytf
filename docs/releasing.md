# Releasing

A release is one tag push. `.github/workflows/release.yml` runs GoReleaser on any `v*` tag.

## What a tag publishes

| Channel | Published | Needs |
|---------|-----------|-------|
| GitHub release (linux and darwin, amd64 and arm64, `checksums.txt`) | always | `GITHUB_TOKEN` (provided by Actions) |
| AUR `lazytf-bin` | only when `AUR_KEY` is set | AUR account and SSH key |
| Homebrew cask `Nicolas-Rigaudy/tap/lazytf` | only when `HOMEBREW_TAP_GITHUB_TOKEN` is set | tap repository and token |

The AUR and Homebrew steps are skipped while their secret is unset, so a tag without them still produces the GitHub release.

## One-time setup

### AUR

1. Create an account on https://aur.archlinux.org and add an SSH public key to it.
2. Reserve the package name: `git clone ssh://aur@aur.archlinux.org/lazytf-bin.git` (an empty clone creates it on first push).
3. Add the private key as the repository secret `AUR_KEY` (Settings > Secrets and variables > Actions).

GoReleaser writes the `PKGBUILD` and `.SRCINFO` and pushes them to the AUR repository.

### Homebrew

1. Create the public repository `Nicolas-Rigaudy/homebrew-tap` with a `Casks/` directory.
2. Create a fine-grained token with Contents read and write on that repository only.
3. Add it as the repository secret `HOMEBREW_TAP_GITHUB_TOKEN`.

GoReleaser writes `Casks/lazytf.rb` and pushes it to the tap. The cask clears the macOS quarantine flag on install, since the binaries are not notarised.

## Cutting a release

1. Move the `Unreleased` section of `CHANGELOG.md` under the new version and date.
2. Check that CI is green on `main`, including the `goreleaser check and snapshot` job.
3. Tag and push: `git tag v1.0.0 && git push origin v1.0.0`.
4. Check the release page, then `yay -S lazytf-bin`, `brew install Nicolas-Rigaudy/tap/lazytf` and `go install github.com/Nicolas-Rigaudy/lazytf/cmd/lazytf@latest`, each followed by `lazytf --version`.

## Local check

```bash
goreleaser check
goreleaser release --snapshot --clean
```

The snapshot build writes the archives, the `PKGBUILD` and the cask to `dist/` without publishing anything.
