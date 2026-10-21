## req_001_r_d_what_to_do_after_the_1_0_blockers_and_what_to_cut_or_defer - R&D: what to do after the 1.0 blockers, and what to cut or defer
> From version: 0.1.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 70%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research note written when every remaining road_001 item was waiting on the owner. It recommends what to do first once they are unblocked, and what to cut or defer from 1.1 to 1.3. Everything below is a proposal; no backlog item has been created from it.
- Keywords: roadmap, prioritisation, plan review, tfui, workspaces, opentofu, terragrunt, 1.0, research
- Use when: Deciding what to build after the 1.0 exit signals, or reordering road_001.
- Skip when: You need execution detail for an existing item.

# Needs
- Decide where effort goes once the owner-gated items (colour style, GIF, clean-machine installs, publishing) are cleared.
- Avoid the roadmap growing past what one maintainer can keep correct: LazyTF runs apply and destroy, so each feature adds blast radius.

# Context
State on 2026-10-22: 0.2 and 0.3 are done (task_001 to task_003). Remaining: task_004 (50%: ANSI colour pick, README five-minute test, manual small/light terminal checks), task_005 (50%: publishing, needs secrets), task_006 (80%: GIF, clean-machine installs). The 0.4 "plan review from `terraform show -json`" scope in road_001 is not covered by any backlog item and is visual, so it was never started.

What similar tools do (from the product brief's comparison and the README; not re-verified online tonight):
- `lmarqs/terraform-ui` (tfui): plan review with risk badges, blast radius, attribute diff, targeted apply, terraform/tofu/terragrunt. Closest competitor, already distributed.
- `idoavrah/terraform-tui`: state browser, taint/untaint/delete. State-first, no environment navigation.
- Terragrunt, Atlantis, Terraform Cloud: pipeline or monorepo tools; LazyTF's non-goals say it is not one of them.
- lazygit (the model): value is fast navigation plus safe defaults, not feature breadth.

# Recommendation (proposals)
Priority order after the owner-gated items:
1. **Plan review MVP (0.4), before the 1.0 tag.** It is the one feature the roadmap promises that no item owns, and the only one that directly reduces wrong-apply risk: counts of create/update/destroy/replace from `terraform show -json`, repeated in the apply confirm, with destroy/replace highlighted. Keep it to counts plus a per-resource action list. Defer per-attribute diff to 1.1. Needs a visual decision, so it waits on the owner.
2. **Fix the hint bar wrapping below about 140 columns** (reported by the README job). Small, but it hurts the first five minutes on an 80-column terminal, which the 0.4 exit signal tests.
3. **Cut a 0.1.1 or 1.0.0-rc release soon after the owner approves.** `go install @latest` and the release archive still give v0.1.0 without `--version`, so the README install table is only true for `@main`. This is the cheapest way to make the public page honest.
4. **1.1: OpenTofu first, workspaces second.** OpenTofu is mostly a configurable binary (done in 0.3) plus a CI job running the integration fixture with `tofu`; low cost, and it is a stated owner need. Workspaces change what an "environment" means in discovery, binding and the saved-plan key, so they touch the correctness core; do them only after OpenTofu and with fixtures first.
5. **1.2: only the safe subset.** `validate`/`fmt` shortcuts and read-only state list/show are cheap and low risk. `-target`/`-replace` from the plan summary is useful and builds on plan review. `state rm`/`mv` should be deferred: it is destructive, a rarely wanted TUI action, and tfui/terraform-tui already cover state surgery.

Cut or defer:
- **1.3 Terragrunt and multi-project combined runs: defer past 1.x until there is user demand.** It needs a second runner abstraction and combined-summary UX, and Terragrunt users already have `run-all` and its own tooling. Revisit only if issues ask for it.
- **Per-project run history:** defer; no user need shown yet.
- **Windows:** keep out until the plan-file and shell assumptions are removed; revisit on demand.
- **Do not chase tfui feature parity** (attribute diff, blast-radius graphs). The differentiator remains multi-project and environment navigation with wrong-environment guards.
- **Possible later ideas, not on the roadmap:** a `--no-color`/`NO_COLOR` honouring pass; a read-only mode flag for demos and review sessions (disables apply/destroy). The read-only mode is cheap and fits the safety story.

Reasoning in one line: correctness and honesty of what ships (items 1 to 3) beat breadth; breadth (1.1 to 1.3) should be pulled by users after 1.0, not pushed on a guess.

# Acceptance criteria
- AC1: The owner reads the recommendations and accepts, edits or rejects each.
- AC2: Accepted proposals are turned into backlog items through the workflow (none exist yet), and road_001 is edited to match.
- AC3: Nothing here is treated as committed scope until then.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable (owner decision).
- [x] Dependencies and known risks are listed (see road_001 risks).

# Companion docs
- Product brief(s): `prod_001_lazytf_a_keyboard_driven_terraform_cockpit`
- Architecture decision(s): (none yet)

# References
- `logics/roadmap/road_001_lazytf_to_1_0.md`
- `logics/request/req_000_lazytf_1_0_public_release.md`

# Backlog
- none (proposals only)
