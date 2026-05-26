# 🚀 Feature Specification Plan: [Specification Name]

**Status:** `Draft` | `Planning` | `Ready for Implementation` | `Done`
**Date:** [YYYY-MM-DD]
**Author:** [Dev/Agent Name]
**Ticket/Reference:** [JIRA/GitHub Link or N/A]
**Complexity:** `High` | `Medium` | `Low`

---

## 1. Plan — Solution Roadmap

> **Agent Instruction:**
> 1. Perform a root cause analysis based **solely** on the Scope section above and the
>    codebase structure — do NOT execute any code change yet.
> 2. Decompose the solution into discrete, file-level tasks.
>    Each task MUST map to **one** file or one logical concern.
>    If a task spans multiple files, split it.
>    Use **semantic IDs** for tasks (e.g., `TASK-DOM-PASSWORD`, `TASK-APP-LOGIN-HANDLER`)
>    instead of sequential numbers to avoid cascading renumbering when tasks are added/removed.
> 3. **After drafting the full plan, present it to the dev and WAIT for explicit approval.**
>    For each task that is **rejected**, ask at most 3 focused questions to gather
>    missing context, revise that task, and re-present it for approval.
>    Only advance to Execution after ALL tasks are approved.

### Technical Strategy & Architecture

* **Architectural Approach:** [Fill during planning — Describe how the feature will be designed, e.g., which aggregates will be modified, new ports, adapters, etc.]
* **Affected Layer(s):** [e.g., Domain | Application | Adapter | Infrastructure]
* **Mapped Skills / Tools:** [e.g., `adjust_go_code`, `make_unit_tests`, `grep_search`]

### Execution Roadmap

> **Format per task:**
> Each `Task` refers to a **single file** (or single isolated concern).
> Each `Sub-task` refers to a specific change within that file.
> Mark `[ ]` → `[/]` when started → `[x]` when dev-approved.

---

- [ ] **TASK-[LAYER]-[CONCERN] — [Short description: file or concern name]**
  - **File:** `[relative/path/to/file.go]`
  - **Reason:** [Why this file needs to change, linked to the root cause.]
  - Sub-tasks:
    - [ ] [TASK-ID].1 — [Specific change: e.g., "Add nil-check guard before calling repository method."]
    - [ ] [TASK-ID].2 — [Specific change: e.g., "Update unit test to cover the new guard branch."]
  - **⏸️ Dev approval required before executing this task.**

---

- [ ] **TASK-[LAYER]-[CONCERN] — [Short description: file or concern name]**
  - **File:** `[relative/path/to/file.go]`
  - **Reason:** [Why this file needs to change.]
  - Sub-tasks:
    - [ ] [TASK-ID].1 — [Specific change.]
    - [ ] [TASK-ID].2 — [Specific change.]
  - **⏸️ Dev approval required before executing this task.**

---

> _(Add more Task blocks as needed. Keep one file = one task. Use semantic IDs: `TASK-[LAYER]-[CONCERN]` where LAYER is DOM/APP/INFRA/MOD and CONCERN is a short kebab-case identifier.)_

---

## 3. Execution — Validated Checklist

> **🚨 CRITICAL Implementation Agent Instruction:**
> * Execute **one task at a time**, strictly following the order defined in Section 2 (Plan).
> * After completing each task (or sub-task group), **STOP**, present the result
>   (diff, test output, log snippet), and **REQUIRE explicit dev validation** before
>   marking it `[x]` and advancing.
> * If a step produces an unexpected side effect or a failing test, revert or isolate
>   the change and return to Planning mode for that task.
> * Do NOT batch-execute the entire plan without checkpoints. Each `[x]` is a
>   dev-signed-off proof of correctness, not just a completion flag.

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Plan]**
  - *Agent Notes:* [Filled by agent: what exactly was changed and why.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Test output / log / diff snippet goes here.]
  - *Validation Status:* `⏳ Waiting for Dev` | `✅ Validated` | `❌ Rejected`

---

- [ ] **Execution — TASK-[LAYER]-[CONCERN]: [Mirror the task name from Plan]**
  - *Agent Notes:* [Fill during execution.]
  - *Files Modified:*
    - `[relative/path/to/file.go]`
  - *Validation Evidence:* [Fill during execution.]
  - *Validation Status:* `⏳ Waiting for Dev` | `✅ Validated` | `❌ Rejected`

---

> _(Mirror every task from Section 2 here. One execution block per plan task.)_

---

## 4. Return — Summary & Handover

> **Agent Instruction:**
> Fill this section **ONLY** when ALL execution tasks have status `✅ Validated`
> and are marked `[x]`.
> This section serves as the official record of what was done, why, and what
> the dev should be aware of going forward.

### Applied Solution Summary

[Brief, objective narrative of what was done to implement the feature — explain how each
task was addressed. Write this as if handing over to another engineer.]

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
- [ ] All post-conditions from Section 1.4 are met.
- [ ] All modified files were covered by tests.
- [ ] This Return section was only filled after full dev validation.
