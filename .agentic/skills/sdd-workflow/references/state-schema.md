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
  current_wave: 1
  completed_tasks: []             # IDs confirmed by dev after each wave
  pending_tasks:   []             # IDs extracted from PLAN.md at the plan→implement transition
  current_wave_groups: []         # parallel_groups being executed in the current wave (for resume)

verify:
  retry_counts: {}                # map task_id → attempt count (max 3)
  failed_tasks: []                # objects: {id, reason} — tasks that exhausted retries
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
- Never overwrite the entire STATE.md; update only the changed fields.

## Wave group logic

Within a single wave, tasks are split by `parallel_group`. Each group runs as an independent
subagent. Groups within the same wave touch different files, so there is no conflict risk.

Example — Wave 1 of a login feature:

```
Wave 1
  group "domain"      → 7 tasks (all contract files)       → subagent A
  group "application" → 5 tasks (command DTOs only)        → subagent B
```

Both subagents are spawned sequentially and their results are collected before presenting
to the dev for approval.

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
  current_wave: 2
  current_wave_groups:
    - application
    - infrastructure
  completed_tasks:
    - TASK-DOM-LOGIN-DATASOURCE
    - TASK-DOM-LOGIN-SESSION-REPO
    - TASK-DOM-REFRESH-SESSION-REPO
    - TASK-DOM-LOGOUT-SESSION-REPO
    - TASK-DOM-REVOKE-SESSION-REPO
    - TASK-DOM-TOKEN-SERVICE
    - TASK-DOM-RATE-LIMITER
    - TASK-APP-LOGIN-COMMAND
    - TASK-APP-REFRESH-COMMAND
    - TASK-APP-LOGOUT-COMMAND
    - TASK-APP-REVOKEONE-COMMAND
    - TASK-APP-REVOKEALL-COMMAND
  pending_tasks:
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
