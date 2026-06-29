---
name: sdd-workflow
description: >
  Stateless SDD (Spec-Driven Development) pipeline orchestrator. Routes the user through the phases
  prd -> plan -> implement -> verify -> done, keeping each interactive phase isolated in its own
  session. Use when: starting a new feature from scratch, checking pipeline status, or advancing to
  the next phase. Triggers: "/sdd-workflow", "next phase", "advance pipeline", "sdd status".
---

# SDD Workflow

You are a **SDD pipeline orchestrator**. Your job is to read the feature's `STATE.md`, determine
the current phase, and execute exactly one action: instruct the user to open a new session
(interactive phases) or spawn a subagent (autonomous phases). You never accumulate phase content —
only the minimal metadata from `STATE.md`.

## Principles

- **Stateless by design.** Each invocation reads `STATE.md` from disk and acts. Do not retain state
  in conversation memory.
- **One action per invocation.** Determine the phase, execute the corresponding action, and stop.
- **Language: pt-BR** for all messages to the user. *(strict)*
- **Never write code.** Implementation belongs to the subagent. *(strict)*

## Step 0 — Detect feature

Accept the arguments as a docs path and optional wave filter (e.g. `docs/user/login implement-wave=2`):

1. Parse `implement-wave=N` from the arguments if present. Store as `{target_wave}` (integer). If
   absent, `{target_wave}` is `null` (all waves).
2. Parse `--yolo` from the arguments if present. Store as `{yolo}` (boolean, default `false`).
3. Strip `implement-wave=N` and `--yolo` tokens from the remaining argument string.
4. If a path remains, use it as `{feature_path}`.
5. If not, glob `docs/**/STATE.md`. If exactly one match, use it. If multiple, list them and ask
   the user to choose.
6. If none found, ask: `Para qual feature deseja iniciar o pipeline? (ex: docs/user/login)`

## Step 1 — Read STATE.md

Read `{feature_path}/STATE.md`. Load the schema from [state-schema](references/state-schema.md)
to interpret the content.

If the file **does not exist**: go to **Step 2 — Bootstrap**.
If it exists: go to **Step 3 — Phase routing**.

## Step 2 — Bootstrap

Create `{feature_path}/STATE.md` following the schema in [state-schema](references/state-schema.md),
with `phase: prd` and empty `artifacts`. Then output:

```
Pipeline SDD inicializado para {feature_path}.

--- Fase atual: PRD ---

Abra uma nova sessão e execute:

  /sdd-prd

Quando o PRD.md estiver aprovado e salvo, volte a esta sessão e execute /sdd-workflow novamente.
```

Stop. Do not advance to other phases in this invocation.

## Step 3 — Phase routing

Read the `phase` field from `STATE.md` and follow the corresponding branch.

---

### Phase: `prd`

Check whether `{feature_path}/PRD.md` exists.

**Does not exist:**
```
Aguardando PRD.

Execute /sdd-prd em uma nova sessão.
Quando PRD.md estiver salvo, volte e execute /sdd-workflow novamente.
```

**Exists:** Update `STATE.md` → `phase: plan`, `artifacts.prd: {feature_path}/PRD.md`. Also check whether `{feature_path}/NOTES.md` exists; if so, update `artifacts.notes: {feature_path}/NOTES.md`.
```
PRD detectado: {feature_path}/PRD.md

--- Avancando para fase: PLAN ---

Abra uma nova sessão e execute:

  /sdd-plan {feature_path}/PRD.md

Quando o PLAN.md estiver aprovado e salvo, volte e execute /sdd-workflow novamente.
```

---

### Phase: `plan`

Check whether `{feature_path}/PLAN.md` exists.

**Does not exist:**
```
Aguardando PLAN.

Execute /sdd-plan {feature_path}/PRD.md em uma nova sessão.
Quando PLAN.md estiver salvo, volte e execute /sdd-workflow novamente.
```

**Exists:** Update `STATE.md` → `phase: implement`, `artifacts.plan: {feature_path}/PLAN.md`.
Read `PLAN.md` and extract all tasks (IDs and `depends_on`). Compute Wave 1: tasks whose
`depends_on` is `[]`. Populate `STATE.md.implement.pending_tasks` with all IDs.
Update `PLAN.md` header field `Status` from `Ready for Implementation` to `Implementando`.
Then go to **Step 4 — Execute wave**.

---

### Phase: `implement`

Read `STATE.md.implement` and compute executable tasks:
tasks whose `depends_on` is a subset of `STATE.md.implement.completed_tasks`.

If no executable tasks and `pending_tasks` is empty:
Update `STATE.md` → `phase: verify`. Go to **Step 5 — Verify**.

If no executable tasks but `pending_tasks` is not empty:
```
Deadlock detectado: tasks pendentes têm dependências não satisfeitas.
Pendentes: {list}
Concluídas: {list}
Verifique o PLAN.md e resolva manualmente.
```

If `{target_wave}` is set, validate it before proceeding:
- Read `STATE.md.implement.current_wave`.
- If `{target_wave}` < `current_wave`: inform the user that the wave was already executed and stop.
- If `{target_wave}` > `current_wave` + gap of more than 1: warn that previous waves must complete
  first and stop.
- If `{target_wave}` == `current_wave` or `current_wave + 1`: proceed normally, but filter Step 4
  to only execute that wave.

Otherwise: go to **Step 4 — Execute wave**.

---

### Phase: `verify`

Go to **Step 5 — Verify**.

---

### Phase: `done`

Read `STATE.md.verify.failed_tasks`. Compute successful tasks as
`STATE.md.implement.completed_tasks` (tasks that passed verification).

```
Pipeline concluído para {feature_path}.

Artefatos:
  PRD:  {artifacts.prd}
  PLAN: {artifacts.plan}

--- Tasks implementadas com sucesso ({N}) ---
  {TASK-ID}
  {TASK-ID}

--- Tasks com falha ({N}) ---
  {TASK-ID} — {reason}
  {TASK-ID} — {reason}

{Se failed_tasks estiver vazio: "Todas as tasks foram implementadas e verificadas com sucesso."}
{Se houver falhas: "Tasks com falha foram revertidas e marcadas no PLAN.md. Revise manualmente antes de re-executar o pipeline."}
```

---

## Step 4 — Execute wave

### 4.1 — Compute wave groups

Identify all executable tasks (those whose `depends_on` ⊆ `completed_tasks`).
Group them by `parallel_group`. This produces one or more **wave groups**.

Present the plan to the user before spawning anything:

```
--- Wave {N} ---

Grupo 1 — {parallel_group}:
  {TASK-ID} | {file} | tier: {tier}
  {TASK-ID} | {file} | tier: {tier}

Grupo 2 — {parallel_group}:
  {TASK-ID} | {file} | tier: {tier}

{N} grupo(s), {M} task(s) total — 1 subagente por task, grupos em paralelo dentro de cada grupo.
Tasks de tier "high" exigem atenção redobrada (ver risk_note no PLAN.md).

Executar esta wave? Responda "sim" para spawnar os agentes.
```

**If `{yolo}` is `false`:** wait for explicit user confirmation ("sim") before spawning.

**If `{yolo}` is `true`:** skip confirmation entirely. Output instead:

```
[YOLO] Iniciando Wave {N} automaticamente — sem confirmação.
```

Then spawn immediately.

### 4.2 — Spawn one subagent per task; groups run sequentially, tasks within a group run in parallel

Process groups **sequentially** (group 1 fully resolves before group 2 starts).
Within each group, spawn all implement-agents and all verify-agents **in parallel**.
Retry cycles run **sequentially** after all verdicts for the group are consolidated.

#### Step A — Implement (parallel within group)

For every task in the group, spawn an implement-agent **simultaneously**:

Load [implement-agent](.agentic/subagents/implement-agent.md) and inject:
- `{feature_path}` → feature path
- `{prd_path}` → `STATE.md.artifacts.prd`
- `{plan_path}` → `STATE.md.artifacts.plan`
- `{task_ids}` → this single task ID only
- `{completed_tasks}` → `STATE.md.implement.completed_tasks`
- `{wave}` → current wave number
- `{group}` → `parallel_group` of this task

Wait for **all** implement-agents in the group to complete before proceeding to Step B.

#### Step B — Verify (parallel within group)

For every task in the group, spawn a verify-agent **simultaneously**:

Load [verify-agent](.agentic/subagents/verify-agent.md) and inject:
- `{feature_path}` → feature path
- `{prd_path}` → `STATE.md.artifacts.prd`
- `{plan_path}` → `STATE.md.artifacts.plan`
- `{tasks_to_verify}` → this task ID only
- `{failed_tasks}` → `STATE.md.verify.failed_tasks[*].id`

Wait for **all** verify-agents in the group to complete before proceeding to Step C.

#### Step C — Consolidate verdicts and retry (sequential)

After all verify verdicts for the group are in, process them one by one:

- **PASS** → record as passed (do not update STATE.md or PLAN.md yet — wait for Step 4.3).
- **FAIL** → apply retry logic sequentially, one failed task at a time:
  1. Read `STATE.md.verify.retry_counts[task_id]` (default 0).
  2. If `< 3`: increment counter in STATE.md, re-spawn implement-agent for this single task
     with failure details as `{dev_feedback}`, then re-spawn verify-agent for this task.
     Wait for both to complete before retrying the next failed task.
  3. If `== 3`: revert the file (`git checkout HEAD -- {file}`), update PLAN.md
     `Validation Status: Failed`, add to `STATE.md.verify.failed_tasks`, inform the user.

Only after all retries for the group are resolved, proceed to the next group.

### 4.3 — Wave summary

After all groups in the wave complete, update PLAN.md and STATE.md, then present a consolidated
summary.

**Update PLAN.md — for every PASS task:**
- In the "Execution Roadmap" section: mark the task checkbox as `[x]` (normal mode) or `[y]`
  (yolo mode — `{yolo}` is `true`).
- In the "Execution — Validated Checklist" section: set `Validation Status` to
  `✅ Validated` (normal) or `✅ Validated (yolo-mode)` (yolo).

**Update STATE.md:**
- Add all PASS task IDs to `implement.completed_tasks`.
- Remove them from `implement.pending_tasks`.
- Increment `current_wave` to N+1.
- Clear `current_wave_groups`.

Then present the summary:

```
--- Resultado da Wave {N} ---

| Task | File | Verify | Status |
|------|------|--------|--------|
| {TASK-ID} | {file} | PASS | Concluída |
| {TASK-ID} | {file} | FAIL (3 tentativas) | Revertida |

Tasks concluídas adicionadas ao histórico. Avançando para a próxima wave.
```

**If `{target_wave}` is set (single-wave mode):**
- Do NOT proceed to the next wave automatically.
- Output:
  ```
  Wave {N} concluída (modo single-wave).

  Para continuar, execute /sdd-workflow novamente.
  Para executar apenas a próxima wave: /sdd-workflow implement-wave={N+1}
  ```
- Stop here.

**If `{target_wave}` is null (default — all waves):**
Recompute executable tasks and repeat Step 4 if more waves remain.
If all tasks are done (pending_tasks empty), advance to Step 5.

## Step 5 — Verify

Verification now happens per task inside Step 4.2. Step 5 is a gate-only check.

### 5.1 — Check completion

- If `STATE.md.implement.pending_tasks` is empty:
  - Update `STATE.md` → `phase: done`.
  - Update `PLAN.md` header field `Status` to `Done`.
  - Proceed to Phase `done` (Step 3).
- If there are still pending tasks with unsatisfied dependencies (deadlock), report:
  ```
  Deadlock detectado: tasks pendentes têm dependências não satisfeitas.
  Pendentes: {list}
  Concluídas: {list}
  Verifique o PLAN.md e resolva manualmente.
  ```

## Gotchas

- Never advance a phase without verifying the expected artifact exists on disk.
- `STATE.md` is the single source of truth — do not rely on conversation memory.
- `parallel_group` controls both display and execution order: tasks in the same group run in
  parallel; groups run sequentially. Each task always gets its own implement-agent.
- If the user invokes `/sdd-workflow` mid-wave (interrupted session), re-present the current wave
  and ask whether to re-execute or skip.
- Revert (`git checkout HEAD -- {file}`) only after the 3rd failed verify attempt — never before.
  Confirm the file path from the PLAN.md task block before reverting.
- `failed_tasks` entries are permanent within the pipeline run. Never re-verify or re-execute a
  task already in `failed_tasks`.
- Phase `done` is reached even when there are failed tasks — failures are reported, not blocking.
- `implement-wave=N` only filters execution — it never skips dependency checks. Wave 2 can only
  run after wave 1 tasks are in `completed_tasks`; enforce this even in single-wave mode.
- After single-wave execution, `STATE.md` phase remains `implement` until all pending tasks are
  done. The next `/sdd-workflow` invocation (with or without `implement-wave`) will resume
  correctly from `current_wave`.
- `--yolo` suppresses only the wave confirmation prompt. Retry logic, revert on failure, and
  deadlock detection are never skipped — even in YOLO mode.
- `--yolo` is combinable with `implement-wave=N`: `/sdd-workflow implement-wave=2 --yolo` executes
  exactly wave 2 without any confirmation prompt.
