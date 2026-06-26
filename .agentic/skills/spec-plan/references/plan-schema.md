# 🚀 Feature Specification Plan: [Specification Name]

**Status:** `Draft` | `Planning` | `Ready for Implementation` | `Done`
**Date:** [YYYY-MM-DD]
**Author:** [Dev/Agent Name]
**Ticket/Reference:** [JIRA/GitHub Link or N/A]
**Complexity:** `High` | `Medium` | `Low`

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
> (`low` ≤ 2.0, `medium` ≤ 3.5, `high` > 3.5); a `high` task carries a `risk_note`.

---

- [ ] **TASK-[LAYER]-[CONCERN] — [Short description: file or concern name]**
  - **File:** `[relative/path/to/file.go]`
  - **Reason:** [Why this file needs to change, linked to the root cause.]
  - **Dependencies:**
    - `depends_on:` `[TASK-ID, TASK-ID]` <!-- semantic IDs, or [] if none -->
    - `parallel_group:` `[domain | application | infrastructure | module | tests | config]`
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
  - **Dependencies:**
    - `depends_on:` `[TASK-ID, TASK-ID]` <!-- semantic IDs, or [] if none -->
    - `parallel_group:` `[domain | application | infrastructure | module | tests | config]`
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

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Section 1]**
  - *Agent Notes:* [Filled by agent: what exactly was changed and why.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Test output / log / diff snippet goes here.]
  - *Validation Status:* `✅ Validated` | `❌ Failed` | `⚠️ Blocked: [reason]`

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Section 1]**
  - *Agent Notes:* [Fill during execution.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Fill during execution.]
  - *Validation Status:* `✅ Validated` | `❌ Failed` | `⚠️ Blocked: [reason]`

---

> _(Mirror every task from Section 1 here. One execution block per plan task.)_

---

## 4. Return — Summary & Handover

> **Agent guidance:** fill this section only when every execution task is `✅ Validated` and marked
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
