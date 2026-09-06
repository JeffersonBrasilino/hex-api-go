# Phase: `implement` / `verify`

## Phase: `implement`

If `PLAN.md`'s header `Status` is still `Ready for Implementation` (first entry into this phase —
`implement.completed_tasks` is empty and `implement.current_wave_groups` is empty), flip it now,
before computing executable tasks:

```
node .agentic/skills/sdd-workflow/scripts/update-plan.cjs {feature_path} --status Implementando
```

Compute the executable tasks for this wave:

```
node .agentic/skills/sdd-workflow/scripts/compute-wave.cjs {feature_path}
```

- **Exit 0, `STATUS: complete`** (no pending tasks left): update `STATE.md` → `phase: verify` via
  `write-state.cjs --phase verify`. Go to **Step 5 — Verify**.
- **Exit 1, `STATUS: deadlock`**: relay the script's stderr output to the user as-is —
  ```
  Deadlock detectado: tasks pendentes têm dependências não satisfeitas.
  Pendentes: {list}
  Concluídas: {list}
  Verifique o PLAN.md e resolva manualmente.
  ```
  Stop.
- **Exit 0, `STATUS: ok`**: the script's stdout (`WAVE`, `GROUPS`, and the per-task table) is the
  wave to present. Continue below.

`implement-wave=N` only filters execution — it never skips dependency checks. If `{target_wave}` is
set, validate it before proceeding:
- Read the script's `WAVE` value as `current_wave`.
- If `{target_wave}` < `current_wave`: inform the user that the wave was already executed and stop.
- If `{target_wave}` > `current_wave` + gap of more than 1: warn that previous waves must complete
  first and stop.
- If `{target_wave}` == `current_wave` or `current_wave + 1`: proceed normally, but filter Step 4
  to only execute that wave.

Otherwise: go to **Step 4 — Execute wave**, using `compute-wave.cjs`'s output directly instead of
recomputing executable tasks by hand.

## Phase: `verify`

Go to **Step 5 — Verify**.

## Step 4 — Execute wave

### 4.1 — Compute wave groups

Use the wave groups and task table already produced by `compute-wave.cjs` in the `implement`
phase branch above — do not recompute them by hand. If the user invokes `/sdd-workflow` mid-wave
(interrupted session), re-present the current wave and ask whether to re-execute or skip.

Present the plan to the user before spawning anything. Always state the wave's position out of the
total (from `compute-wave.cjs`'s `TOTAL_WAVES`) so it is never ambiguous which layers this invocation
will and will not touch:

```
--- Wave {N} de {TOTAL_WAVES} ---

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

**If `{yolo}` is `true`:** skip confirmation entirely (only the confirmation prompt is skipped —
retry logic, revert on failure, and deadlock detection still run normally). Output instead:

```
[YOLO] Iniciando Wave {N} automaticamente — sem confirmação.
```

Then spawn immediately.

### 4.2 — Spawn one batched subagent per group; groups run sequentially

Process groups **sequentially** (group 1 fully resolves before group 2 starts — a deliberate choice
to bound concurrent agent load, not a correctness requirement, since groups within a wave never
depend on each other). Each group is handled by **one implement-agent and one verify-code agent,
batched over every task in that group** — not one agent per task. A single wave commonly contains
several groups at once; what must never happen is a task from a *different, higher* wave slipping
into the current one. Retry cycles run **sequentially**, scoped to the individual failed task,
after the group's verdict table is consolidated.

#### Step A — Implement (one batched agent per group)

Before spawning, compute `{has_dedicated_test_task}` for every task in the group **once**, without
loading PLAN.md in full — grep just the `TASK-TEST-*` blocks' `File:` lines and match them against
this group's task files:

```
grep -A 2 '^\- \[.\] \*\*TASK-TEST-' {plan_path} | grep 'File:'
```

A task's `{has_dedicated_test_task}` is `true` if its own `File:` (from `compute-wave.cjs`'s table)
appears in that output, else `false`. Build the `{tasks_table}` from this (one line per task:
`{task_id} — has_dedicated_test_task: {true|false}`).

Spawn **one** implement-agent for the whole group:

Load [implement-agent](.agentic/subagents/implement-agent.md) and inject:
- `{group_name}` → this group's `parallel_group` value
- `{prd_path}` → `STATE.md.artifacts.prd`
- `{plan_path}` → `STATE.md.artifacts.plan`
- `{task_id_list}` → every task ID in this group, in order
- `{tasks_table}` → built above
- `{completed_tasks}` → `STATE.md.implement.completed_tasks`

Wait for it to complete before proceeding to Step A.5.

#### Step A.5 — Build & test the group (once)

All tasks in the group have finished writing by this point. Run the suite **once** for the whole
group here, not once per task in Step B.

```
node .agentic/skills/sdd-workflow/scripts/run-build-test.cjs
```

The script returns a compact structured summary, never raw logs. Capture its stdout as
`{build_result}` — pass it verbatim into every verify-code spawned in Step B; none of them re-run
these commands themselves.

#### Step B — Verify code (one batched agent per group)

Build the `{tasks_table}` for verification (one line per task: `{task_id} — already_failed:
{true|false}`, from `STATE.md.verify.failed_tasks[*].id`).

Spawn **one** verify-code agent for the whole group:

Load [verify-code](.agentic/subagents/verify-code.md) and inject:
- `{group_name}` → this group's `parallel_group` value
- `{plan_path}` → `STATE.md.artifacts.plan`
- `{tasks_table}` → built above
- `{build_result}` → captured in Step A.5

Note this agent does **not** take `{prd_path}` — PRD conformance is checked once per wave in Step
C.5 below, not per task here. verify-code's terminal success state is `Verified`, not `Done`.

Wait for it to complete before proceeding to Step C.

#### Step C — Consolidate verdicts and retry (sequential)

The group's verify-code agent returns one verdict row per task. Process them one by one:

- **PASS (`Verified`)** → record as verified (do not update STATE.md or PLAN.md yet — this task
  still needs the wave-level PRD gate in Step C.5 before it can reach `completed_tasks`).
- **FAIL** → apply retry logic sequentially, one failed task at a time. Corrections are always
  scoped to the single failed task — never re-batch the whole group for a retry, only the task(s)
  that actually failed:
  1. Read `STATE.md.verify.retry_counts[task_id]` (default 0).
  2. If `< 3`: increment the counter via
     `write-state.cjs {feature_path} --increment-retry {task_id}`, re-spawn implement-agent with a
     batch of exactly this one task, failure details as `{dev_feedback}`. Once it completes,
     re-run `node .agentic/skills/sdd-workflow/scripts/run-build-test.cjs` — Step A.5's
     `{build_result}` is now stale for this task since the code just changed — and pass the
     fresh output as `{build_result}` into a re-spawned verify-code agent, also batched to just
     this one task. Wait for both to complete before retrying the next failed task.
  3. If `== 3`: revert the file (`git checkout HEAD -- {file}`), then run:
     ```
     node .agentic/skills/sdd-workflow/scripts/update-plan.cjs {feature_path} \
       --mark-failed {task_id}
     node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
       --add-failed "{task_id}:{reason}"
     ```
     Inform the user. This task is now permanently in `failed_tasks` — never re-verify or
     re-execute it within this pipeline run.

Only after all retries for the group are resolved, proceed to the next group.

#### Step C.5 — Wave PRD gate (once per wave, after every group resolves)

Runs once, after **every** group in the wave has finished Step C — not per group (PRD acceptance
criteria are written at the feature/flow level, not the file level; see `verify-wave-prd.md`).

1. Compute `wave_has_behavior` from the tasks with status `Verified` in this wave, using the
   `tier`/`parallel_group` data `compute-wave.cjs` already reported (no new parsing): `false` only
   if every one of them is `parallel_group: domain`/`config` **and** `tier: low`; `true` otherwise.
2. **If `wave_has_behavior` is `false`:**
   ```
   node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
     --prd-gate-wave {N} --prd-gate-status skipped
   ```
   Every `Verified` task in the wave is auto-approved — treat it as `Done` and continue to 4.3.
3. **Else:** spawn one [verify-wave-prd](.agentic/subagents/verify-wave-prd.md) agent for the whole
   wave and inject:
   - `{wave_number}` → `{N}`
   - `{plan_path}` → `STATE.md.artifacts.plan`
   - `{prd_path}` → `STATE.md.artifacts.prd`
   - `{verified_tasks}` → every task ID (+ `File:`) with status `Verified` in this wave
   ```
   node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
     --prd-gate-wave {N} --prd-gate-status pending
   ```
   Wait for it to complete, then for each row in its report:
   - **PASS** → task is now `Done`. Continue to 4.3.
   - **FAIL** → treat exactly like a Step C failure (same `retry_counts` cap, same
     re-spawn-implement-agent mechanism), with one addition: if `{task_id}` is already in
     `STATE.md.implement.completed_tasks` (i.e. it belongs to an **earlier, already-completed
     wave**), first reopen it —
     ```
     node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
       --remove-completed {task_id} --add-pending {task_id}
     ```
     — and warn the user which later, currently-completed tasks import or depend on the reopened
     file (a notice only — do not auto-re-verify those downstream tasks, that needs explicit dev
     sign-off). Then apply the same increment-retry / re-spawn-implement-agent / re-verify sequence
     as Step C.
   Once every task's verdict is resolved (PASS or exhausted retries → `Failed`):
   ```
   node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
     --prd-gate-wave {N} --prd-gate-status passed   # or --prd-gate-status failed if any task in
                                                      # this wave ended in Failed
   ```

### 4.3 — Wave summary

After Step C.5 resolves for the wave, update PLAN.md and STATE.md via the scripts, then present a
consolidated summary. Only tasks that reached `Done` (Step C.5) are marked here — a task that is
merely `Verified` has not cleared this wave yet.

**Update PLAN.md — for every `Done` task:**
```
node .agentic/skills/sdd-workflow/scripts/update-plan.cjs {feature_path} \
  --mark-done "{DONE_TASK_IDS}"      # normal mode
node .agentic/skills/sdd-workflow/scripts/update-plan.cjs {feature_path} \
  --mark-yolo "{DONE_TASK_IDS}"      # yolo mode ({yolo} is true)
```

**Update STATE.md — for every `Done` task:**
```
node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
  --add-completed "{TASK-ID}" \
  --remove-pending "{TASK-ID}"
```
**Refresh the wave cache — mandatory after every `Done` task, never skip this:** never assume the
next wave is `{N+1}` by arithmetic (wave numbers are not guaranteed dense); re-run `compute-wave.cjs`
and cache whatever `WAVE` it reports:
```
node .agentic/skills/sdd-workflow/scripts/compute-wave.cjs {feature_path}
```
- If it reports `STATUS: ok`, cache its `WAVE` value and clear the in-progress marker:
  ```
  node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
    --current-wave {WAVE_from_above} --clear-wave-groups
  ```
- If it reports `STATUS: complete`, there is no next wave — just clear the in-progress marker and let
  Step 5 handle the `phase: verify` transition:
  ```
  node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} --clear-wave-groups
  ```

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
  Para executar apenas a próxima wave: /sdd-workflow implement-wave={next WAVE reported by compute-wave.cjs above}
  ```
- Stop here.

**If `{target_wave}` is null (default — all waves):**
Re-run `compute-wave.cjs` and repeat Step 4 if more waves remain (`STATUS: ok`).
If it reports `STATUS: complete`, advance to Step 5.

## Step 5 — Verify

Verification now happens per task inside Step 4.2. Step 5 is a gate-only check.

### 5.1 — Check completion

- If `STATE.md.implement.pending_tasks` is empty:
  - Update `STATE.md` → `phase: done` via `write-state.cjs {feature_path} --phase done`.
  - Update `PLAN.md` header field `Status` to `Done` via
    `update-plan.cjs {feature_path} --status Done`.
  - Proceed to Phase `done` — load [phase-done](phase-done.md).
- If there are still pending tasks with unsatisfied dependencies (deadlock), report:
  ```
  Deadlock detectado: tasks pendentes têm dependências não satisfeitas.
  Pendentes: {list}
  Concluídas: {list}
  Verifique o PLAN.md e resolva manualmente.
  ```
