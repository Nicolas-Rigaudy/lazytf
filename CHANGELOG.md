# Changelog

All notable changes to LazyTF are listed here. Release notes on GitHub are generated from commit subjects; this file keeps the human summary.

## Unreleased

### Changed
- The saved plan is keyed by project path and environment, so two environments or two projects with the same directory name never share a plan. Plans live in the user cache directory (0700), not in `/tmp`.
- `config.yml` is merged over the defaults key by key: a file that only sets `search_paths` keeps the default ignore patterns.
- Plan, apply and destroy run with `-input=false`; the approval prompt still reads the confirmation sent by LazyTF.
- The default search path is the current directory only. LazyTF no longer scans `~/Projects` and `~/Documents` unless `search_paths` names them.
- `terraform.tfvars` and `*.auto.tfvars` are no longer listed as environments.
- The AWS SSO login is no longer on by default. Add `auth: {provider: aws-sso}` to the config to keep it.

### Added
- Projects written only in `*.tf.json` are detected; a directory with Terraform files opens as a project whether or not it is initialised.
- Layout configuration: `var_files`, `backend_configs`, `exclude`, `include`, `env_pattern` and `backend_env_pattern`, in `config.yml` or a per-repo `.lazytf.yml`.
- A project with no var files gets a `(no var-file)` environment; a project with no backend config files is initialised without `-backend-config`, so local-backend projects work.
- LazyTF remembers the environment and backend file of its last `init` instead of guessing the environment from the state key.
- `terraform:` setting to run another executable, such as `tofu`.
- Auth providers `aws-sso`, `gcloud`, `az` and `custom`, with per-environment variables (`auth.env`), and retry of the failed command after login.
- `--version`, `--help`, `--config` and a path argument; clear messages when no project or no `terraform` binary is found.
- Release archives for macOS (Intel and Apple Silicon).

### Fixed
- Ignore patterns such as `*/.git` never matched; they now match trailing path segments.
- Failed commands name the command and the reason in the status bar; start failures are no longer reported as "exit code 0".
- The header box was two columns wider than the terminal, so its right border was cut off.
