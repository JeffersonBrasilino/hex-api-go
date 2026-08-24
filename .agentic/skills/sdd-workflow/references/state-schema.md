# STATE.md — Schema

`STATE.md` is the single source of truth for the orchestrator between sessions.
It never contains artifact content — only paths and progress metadata.

## Template

```yaml
# SDD State
feature: {module}/{feature-name}
phase: prd                        # prd | plan | implement | verify | done

artifacts:
  prd:   ""                       # filled when PRD.md is detected
  notes: ""                       # filled when NOTES.md is detected
  plan:  ""                       # filled when PLAN.md is detected

implement:
  current_wave: 1                 # CACHE ONLY — see "current_wave is a cache" below
  completed_tasks: []             # IDs confirmed by dev after each wave
  pending_tasks:   []             # IDs extracted from PLAN.md at the plan→implement transition
  current_wave_groups: []         # parallel_groups being executed in the current wave (for resume)

verify:
  retry_counts: {}                # map task_id → attempt count (max 3)
  failed_tasks: []                # objects: {id, reason} — tasks that exhausted retries
  prd_gate:                       # status of the once-per-wave verify-wave-prd gate
    wave: null                    # which wave this status refers to
    status: pending                # pending | passed | failed | skipped

pr:
  status: not_started             # not_started | blocked_missing_tool |
                                   # blocked_protected_branch | created | failed
  url: ""                         # filled once the PR/MR is created
  branch: ""                      # filled once prepare-pr.cjs decides/creates the branch
```

## Rules

- `phase` advances only when the expected artifact exists on disk and the dev confirms.
- `completed_tasks` is updated **only after dev approval** (Step 4.3), never before.
- `pending_tasks` is fully populated at the `plan → implement` transition; tasks migrate
  to `completed_tasks` as waves are approved.
- `current_wave_groups` is populated at Step 4.1 and cleared when the wave is fully approved.
  It allows the orchestrator to resume a wave that was interrupted mid-execution.
- `retry_counts` is incremented by the orchestrator each time a verify failure triggers a
  correction cycle. Never reset after increment.
- `failed_tasks` is populated when a task reaches 3 failed verify attempts. Each entry is an
  object `{id: "TASK-ID", reason: "..."}`. Tasks in this list are never re-executed.
- `prd_gate` tracks the once-per-wave `verify-wave-prd` gate: `status` moves
  `pending` → `passed`/`failed`/`skipped` as the wave's gate resolves. A task only reaches
  `completed_tasks` after both its own `verify-code`
  checklist **and** this gate pass (or the gate is `skipped` because the wave has no checkable
  behavior). If the gate finds a defect traced to an **already-completed, earlier-wave** task, that
  task is moved back out of `completed_tasks` into `pending_tasks` (`--remove-completed` /
  `--add-pending`) for a correction cycle — this is the one case where a task leaves
  `completed_tasks` after entering it.
- Never overwrite the entire STATE.md; update only the changed fields.
- `pr` is bookkeeping for the `done`-phase PR/MR gate (`prepare-pr.cjs` / `pr-writer`), not derived
  from disk. `status` only reaches `created` once a PR/MR actually exists — `blocked_missing_tool`
  and `blocked_protected_branch` are resumable stops, not failures, and a re-run of
  `/sdd-workflow` re-enters this gate automatically while `status` is anything but `created`.
  `archive-spec` is only spawned after `status` reaches `created`.

## `current_wave` is a cache, not the source of truth

Wave numbers are decided once, at plan time, by `sdd-plan`'s `compute-waves.cjs`, and written
into each task's `**Wave:**` field in `PLAN.md` — never in `STATE.md`. `implement.current_wave` here
is only a display cache of the last value `compute-wave.cjs` reported; `sdd-workflow` never
increments it by arithmetic (`N+1`) and never trusts it as authoritative. Every real decision —
"what is the next wave", "is this task in this wave" — is recomputed fresh from `PLAN.md`'s `wave:`
fields and `implement.completed_tasks` each time `compute-wave.cjs` runs. If this field and the
plan's actual `wave:` fields ever disagree (e.g. after a manual edit), the plan wins.

## Wave group logic

A wave is a strict boundary: nothing in wave N+1 starts before every task in wave N is in
`completed_tasks`. Within a single wave, tasks are additionally split by `parallel_group` — each
group runs as an independent subagent, and groups within the same wave run sequentially (their tasks
run in parallel within the group). A wave commonly contains more than one `parallel_group` when
several layers all depend on nothing later than the previous wave — that is expected, not a bug. What
must never happen is a task from a genuinely later wave executing early.

Example — a login feature's wave map (as computed by `compute-waves.cjs`, domain/config floor=1,
application floor=2, infrastructure/module floor=3):

```
Wave 1 — camadas: domain, config     → 7 domain contract tasks + 2 config tasks   → 2 subagents (groups), sequential
Wave 2 — camadas: application        → 5 command DTO tasks (depends_on: [], but application floor=2)
Wave 3 — camadas: application, infrastructure, tests → handlers + adapters + their unit tests,
                                        since they all depend only on waves 1-2
```

Domain contracts (wave 1) and application command DTOs (wave 2) are never bundled into the same wave
even though the DTOs have `depends_on: []` — the application layer's floor keeps them a wave apart.
This is the property that makes `implement-wave=1` execute only the domain/config layer, never the
application layer alongside it.

All subagents within a wave's groups are spawned sequentially per group (parallel within the group)
and their results are collected before presenting to the dev for approval.

## Filled example (verify phase, 1 failed task after 3 retries)

```yaml
# SDD State
feature: user/login
phase: verify

artifacts:
  prd:   docs/user/login/PRD.md
  plan:  docs/user/login/PLAN.md

implement:
  current_wave: 3
  current_wave_groups: []
  completed_tasks:
    - TASK-DOM-LOGIN-DATASOURCE
    - TASK-APP-LOGIN-COMMAND
  pending_tasks: []

verify:
  retry_counts:
    TASK-INFRA-JWT-SERVICE: 3
  failed_tasks:
    - id: TASK-INFRA-JWT-SERVICE
      reason: "Token expiry field does not match PRD requirement (field 'exp_seconds' missing from JWTClaims). 3 correction attempts exhausted."
```

---

## Filled example (implement phase, wave 2 in progress)

```yaml
# SDD State
feature: user/login
phase: implement

artifacts:
  prd:   docs/user/login/PRD.md
  notes: docs/user/login/NOTES.md
  plan:  docs/user/login/PLAN.md

implement:
  current_wave: 2                 # cache — the plan's wave: fields are the source of truth
  current_wave_groups:
    - application
  completed_tasks:
    # Wave 1 (floor: domain=1) — fully done, so wave 2 is now executable.
    - TASK-DOM-LOGIN-DATASOURCE
    - TASK-DOM-LOGIN-SESSION-REPO
    - TASK-DOM-REFRESH-SESSION-REPO
    - TASK-DOM-LOGOUT-SESSION-REPO
    - TASK-DOM-REVOKE-SESSION-REPO
    - TASK-DOM-TOKEN-SERVICE
    - TASK-DOM-RATE-LIMITER
  pending_tasks:
    # Wave 2 (floor: application=2) — in progress, group "application" above.
    - TASK-APP-LOGIN-COMMAND
    - TASK-APP-REFRESH-COMMAND
    - TASK-APP-LOGOUT-COMMAND
    - TASK-APP-REVOKEONE-COMMAND
    - TASK-APP-REVOKEALL-COMMAND
    # Wave 3+ (depends_on wave-2 commands, or floor: infrastructure=3) — not yet executable.
    - TASK-APP-LOGIN-HANDLER
    - TASK-APP-REFRESH-HANDLER
    - TASK-APP-LOGOUT-HANDLER
    - TASK-APP-REVOKEONE-HANDLER
    - TASK-APP-REVOKEALL-HANDLER
    - TASK-INFRA-REDIS-SESSION-REPO
    - TASK-INFRA-JWT-SERVICE
    - TASK-INFRA-REDIS-RATELIMITER
    - TASK-INFRA-LOGIN-HTTP
    - TASK-INFRA-REFRESH-HTTP
    - TASK-INFRA-LOGOUT-HTTP
    - TASK-INFRA-REVOKEONE-HTTP
    - TASK-INFRA-REVOKEALL-HTTP
    - TASK-TEST-LOGIN
    - TASK-TEST-REFRESH
    - TASK-TEST-LOGOUT
    - TASK-TEST-REVOKEONE
    - TASK-TEST-REVOKEALL
    - TASK-MOD-BOOTSTRAP
    - TASK-MOD-MAIN
```
