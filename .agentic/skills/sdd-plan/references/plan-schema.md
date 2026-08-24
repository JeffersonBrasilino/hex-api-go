# 🚀 Feature Specification Plan: [Specification Name]

**Status:** `Draft` | `Planning` | `Ready for Implementation` | `Done`
**Date:** [YYYY-MM-DD]
**Author:** [Dev/Agent Name]
**Ticket/Reference:** [JIRA/GitHub Link or N/A]
**Complexity:** `High` | `Medium` | `Low`
**Type:** `feat` | `fix` | `refactor` | `perf` | `chore` | `docs` | `test` | `build` | `ci`

---

## 1. Plan — Solution Roadmap

> **Agent guidance:**
> 1. Base the decomposition solely on the Scope section above and the codebase structure — no code
>    changes yet.
> 2. Break the solution into discrete, file-level tasks. Each task maps to one file or one logical
>    concern; split anything that spans multiple files. Use semantic IDs (e.g. `TASK-DOM-PASSWORD`,
>    `TASK-APP-LOGIN-HANDLER`) so adding or removing tasks never forces a renumber.
> 3. For each task, record its dependencies (`depends_on`, `parallel_group`) and its complexity
>    (five dimensions, `score`, `tier`) in the task block.
> 4. Present the plan and wait for explicit dev approval. For a rejected task, ask at most 3 focused
>    questions, revise it, and re-present it. Move to execution only after every task is approved.

### Technical Strategy & Architecture

* **Architectural Approach:** [How the feature is designed — aggregates modified, new ports, adapters, etc.]
* **Affected Layer(s):** [Domain | Application | Infrastructure | Module]
* **Mapped Skills / Tools:** [e.g. `adjust-go-code`, `make-unit-tests`, `ddd-module-knowledge`]

### Execution Roadmap

> **Format per task:** each `Task` is a single file (or single isolated concern); each `Sub-task` is a
> specific change within that file. Mark `[ ]` → `[/]` when started → `[x]` when dev-approved.
> `depends_on` lists the semantic IDs of tasks that must complete first (`[]` when none); tasks with
> an empty `depends_on` are immediate parallel candidates. `tier` is derived from `score`
> (`low` ≤ 2.0, `medium` ≤ 3.5, `high` > 3.5); a `high` task carries a `risk_note`. `wave` is never
> decided by hand — it is computed and written by `compute-waves.cjs` from `depends_on` and
> `parallel_group` (a task can never share a wave with a task it depends on, and never lands in an
> earlier wave than its own layer's floor: domain/config=1, application=2, infrastructure/module=3).
> Re-run `compute-waves.cjs` after any `depends_on`/`parallel_group` edit — it is idempotent.

---

### Wave Map

> _(Generated and kept up to date by `compute-waves.cjs` — do not edit by hand. One row per wave:
> which layers it groups, which earlier wave(s) it depends on, and which tasks it contains.)_

| Wave | Camadas (parallel_group) | Depende de | Tasks |
|------|---------------------------|------------|-------|

---

- [ ] **TASK-[LAYER]-[CONCERN] — [Short description: file or concern name]**
  - **File:** `[relative/path/to/file.go]`
  - **Reason:** [Why this file needs to change, linked to the root cause.]
  - **Requirement:** [RF-0X / RNF-0X / RN-0X id(s) from the PRD this task implements — the specific requirement(s), not just the general area]
  - **Dependencies:**
    - `depends_on:` `[TASK-ID, TASK-ID]` <!-- semantic IDs, or [] if none -->
    - `parallel_group:` `[domain | application | infrastructure | module | tests | config]`
    - `wave:` `[N]` <!-- computed by compute-waves.cjs — leave any placeholder value, it will be overwritten -->
  - **Complexity:**
    | Dimension       | Score |
    |-----------------|-------|
    | scope           | X     |
    | ambiguity       | X     |
    | coupling        | X     |
    | novelty         | X     |
    | reversibility   | X     |
    | **score**       | X.X   |
    | **tier**        | low / medium / high |
    - `risk_note:` _(only when tier = high)_ [Main risk in one line.]
  - **Sub-tasks:**
    - [ ] [TASK-ID].1 — [Specific change: e.g., "Add nil-check guard before calling repository method."]
    - [ ] [TASK-ID].2 — [Specific change: e.g., "Update unit test to cover the new guard branch."]
  - **Completion criterion:** [How to know this task is done.]

---

- [ ] **TASK-[LAYER]-[CONCERN] — [Short description: file or concern name]**
  - **File:** `[relative/path/to/file.go]`
  - **Reason:** [Why this file needs to change.]
  - **Requirement:** [RF-0X / RNF-0X / RN-0X id(s) from the PRD this task implements — the specific requirement(s), not just the general area]
  - **Dependencies:**
    - `depends_on:` `[TASK-ID, TASK-ID]` <!-- semantic IDs, or [] if none -->
    - `parallel_group:` `[domain | application | infrastructure | module | tests | config]`
    - `wave:` `[N]` <!-- computed by compute-waves.cjs -->
  - **Complexity:**
    | Dimension       | Score |
    |-----------------|-------|
    | scope           | X     |
    | ambiguity       | X     |
    | coupling        | X     |
    | novelty         | X     |
    | reversibility   | X     |
    | **score**       | X.X   |
    | **tier**        | low / medium / high |
    - `risk_note:` _(only when tier = high)_ [Main risk in one line.]
  - **Sub-tasks:**
    - [ ] [TASK-ID].1 — [Specific change.]
    - [ ] [TASK-ID].2 — [Specific change.]
  - **Completion criterion:** [How to know this task is done.]

---

> _(Add more Task blocks as needed. One file = one task. Use semantic IDs: `TASK-[LAYER]-[CONCERN]`,
> where LAYER is DOM/APP/INFRA/MOD and CONCERN is a short kebab-case identifier.)_

---

## 3. Execution — Validated Checklist

> **Execution protocol** *(strict)*:
> - Execute tasks in dependency order: a task runs only after every ID in its `depends_on` is marked
>   `[x]`. Tasks sharing a `parallel_group` with no pending dependency may be batched.
> - Mark a task `[x]` and fill its execution block after: (1) `go build ./...` passes, (2) tests
>   pass or are explicitly skipped with reason, (3) completion criterion is met.
> - Do NOT stop between tasks to wait for dev input — approval happens at the wave level, controlled
>   by the orchestrator. Run all assigned tasks, then produce the consolidated report.
> - If a task fails the build or a test after implementation, fix it before advancing. If
>   unresolvable, record the blocker in the execution block and stop only that task.
>
> **`Validation Status` lifecycle** *(each value is written by a specific agent/script — never
> hand-edited)*:
> - `pending` — default, before the implement-agent starts this task.
> - `Implemented` — set by the implement-agent once build + scoped tests pass. This is the
>   structural precondition `validate-execution.cjs` checks before verify-code runs its own
>   checklist — it does not mean the task has been verified yet. The implement-agent also resets the
>   status back to `Implemented` after fixing a task during a retry cycle.
> - `Blocked: [reason]` — set by the implement-agent instead of `Implemented` when it hits a
>   technical contradiction or impossibility it cannot resolve. Stops only this task; verify-code
>   does not attempt to verify a `Blocked` task.
> - `Verified` — set by **verify-code** once build, scoped tests, PLAN structural check
>   (`validate-execution.cjs`), and architectural compliance all pass. **Not** the terminal state:
>   PRD conformance for the task's requirement has not been gated yet. A `Verified` task still
>   needs `verify-wave-prd` (run once per wave — see `.agentic/subagents/verify-wave-prd.md`) to
>   reach `Done`.
> - `Verification Failed` — set either by verify-code (build/test/plan/architecture checklist did
>   not all pass) or by `verify-wave-prd` (the task's output does not satisfy the PRD requirement it
>   was designed to fulfill — including a defect traced to a task from an **earlier, already-
>   completed wave**, discovered only once a later wave wires it into observable behavior). This is
>   a **transient** signal in both cases — it means "the orchestrator's retry cycle needs to re-run
>   the implement-agent for this task", not a terminal state. The implement-agent overwrites it with
>   `Implemented` once the fix lands.
> - `Failed` — the **terminal**, permanent failure state. Set only by the orchestrator via
>   `update-plan.cjs --mark-failed`, and only after the retry cycle exhausts its 3 attempts. Never
>   set directly by verify-code or verify-wave-prd — that would conflate a single failed attempt
>   with the give-up decision, which is the orchestrator's call (it owns `retry_counts`).
> - `Done` — set by **verify-wave-prd** once its PRD conformance check passes for the task, or by
>   the orchestrator directly when the whole wave has no checkable behavior (`wave_has_behavior:
>   false` — the gate is `skipped`, every `Verified` task in that wave is auto-approved). Terminal
>   state for the task.

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Section 1]**
  - *Agent Notes:* [Filled by agent: what exactly was changed and why.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Test output / log / diff snippet goes here.]
  - *Validation Status:* `pending` | `Implemented` | `Blocked: [reason]` | `Verification Failed` | `Failed` | `Done`

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Section 1]**
  - *Agent Notes:* [Fill during execution.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Fill during execution.]
  - *Validation Status:* `pending` | `Implemented` | `Blocked: [reason]` | `Verification Failed` | `Failed` | `Done`

---

> _(Mirror every task from Section 1 here. One execution block per plan task.)_

---

## 4. Return — Summary & Handover

> **Agent guidance:** fill this section only when every execution task is `Done` and marked
> `[x]`. It is the official record of what was done, why, and what the dev should be aware of going
> forward.

### Applied Solution Summary

[Brief, objective narrative of what was done to implement the feature — explain how each task was
addressed. Write this as if handing over to another engineer.]

### Modified Files

| File | Change Type | Description |
|------|-------------|-------------|
| `[relative/path/to/file.go]` | `Modified` / `Created` / `Deleted` | [What changed and why.] |

### Test Coverage

| Test File | Status | Notes |
|-----------|--------|-------|
| `[relative/path/to/file_test.go]` | `Added` / `Updated` / `Unchanged` | [What was covered.] |

### Side Effects & Warnings

* [e.g., "The refactored function now returns an additional error type — callers outside this spec scope may need updating."]
* [e.g., "Performance impact not benchmarked; recommend adding a benchmark before next release."]

### Spec-Driven Compliance Checklist

- [ ] Scope was fully defined before execution began.
- [ ] Plan was explicitly approved by the dev for each task.
- [ ] Execution followed the one-task-at-a-time checkpoint protocol.
- [ ] No constraint was violated during execution.
- [ ] All post-conditions from the approved Scope are met.
- [ ] All modified files were covered by tests.
- [ ] This Return section was only filled after full dev validation.
