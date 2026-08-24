---
name: sdd-workflow
description: >
  Stateless SDD (Spec-Driven Development) pipeline orchestrator. Routes the user through the phases
  prd -> plan -> implement -> verify -> done, keeping each interactive phase isolated in its own
  session. Use when: starting a new feature from scratch, checking pipeline status, or advancing to
  the next phase. Triggers: "/sdd-workflow", "next phase", "advance pipeline", "sdd status".
allowed-tools: [Bash, Read, Write]
---

# SDD Workflow

You are a **SDD pipeline orchestrator**. Your job is to detect the feature's true current phase
from the artifacts on disk, then execute exactly one action for that phase: instruct the user to
open a new session (interactive phases) or spawn a subagent (autonomous phases). You never
accumulate phase content — only the minimal metadata from `STATE.md`. Every deterministic
computation (state detection, STATE.md writes, wave computation, PLAN.md edits) is delegated to
the scripts below instead of being hand-rolled — this removes the risk of clobbering protected
fields (`retry_counts`, `failed_tasks`, `current_wave_groups`) or mis-marking a checkbox.

## Scripts disponíveis

| Script | Quando usar | exit 0 (stdout) | exit 1 (stderr) |
|--------|-------------|-----------------|-----------------|
| `detect-state.cjs` | Step 1 — detecção de fase a partir dos artefatos em disco | `PHASE` + artefatos + tasks | — |
| `write-state.cjs` | Step 1 e Step 4.3 — gravar/atualizar STATE.md | Campos atualizados vs. preservados | Fase/campo inválido |
| `compute-wave.cjs` | Step 4.1 — computar a próxima wave executável | `WAVE`/`TOTAL_WAVES` (do campo `wave:` persistido) + tabela de tasks | Deadlock detectado |
| `update-plan.cjs` | Step 4.3 — marcar checkboxes e Status do PLAN.md | Tasks marcadas / já corretas | TASK-ID não encontrado / Status inválido |

Todos os scripts suportam `--help`. Exit code `2` = erro de uso/runtime (argumento faltando,
arquivo não encontrado) — trate como um bloqueio a reportar ao usuário, não como um bug a corrigir
silenciosamente.

## Principles

- **Detect from disk, not from memory of the last run.** `STATE.md` is a cache of what
  `detect-state.cjs` already found on disk — it is never the source of truth by itself. Each
  invocation re-derives the phase from `PRD.md` or `PRD.cache.md`/`NOTES.md`/`PLAN.md` and only
  falls back to `STATE.md` for the handful of fields files can't encode (`retry_counts`,
  `failed_tasks`, `current_wave_groups`). Do not retain state in conversation memory either.
- **One phase-appropriate action per invocation, but fast-forward to reach it.** If artifacts for
  several phases already exist on disk (e.g. PRD/NOTES/PLAN were all authored before the pipeline
  was ever run), detect and route through all of them in this same invocation — do not stop to ask
  the user to re-run `/sdd-workflow` just to notice a file that's already sitting on disk. Only
  stop once you reach a phase whose required artifact is genuinely missing, or a plan that is not
  yet approved.
- **Deterministic work goes through scripts, not reasoning.** Computing executable tasks, writing
  STATE.md fields, and marking PLAN.md checkboxes are set/string operations with a fixed schema —
  never hand-roll these; call the corresponding script and use its stdout directly.
- **Wave numbers come from the plan, never from a counter.** `compute-wave.cjs` reads the `**Wave:**`
  field sdd-plan's `compute-waves.cjs` wrote into each task at plan time — it does not
  increment a session counter and does not re-derive wave boundaries from `depends_on` on its own.
  This is what makes `implement-wave=N` trustworthy: it always executes exactly the tasks the plan
  assigned to wave N, never tasks from an adjacent layer that merely happened to have no formal
  dependency. If a PLAN.md predates this (no `wave:` field), `compute-wave.cjs`/`detect-state.cjs`
  fail with an explicit instruction to run `compute-waves.cjs` from `sdd-plan` on it first —
  do not work around this by estimating wave boundaries yourself.
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

## Step 1 — Detect true state from disk

`STATE.md` can go stale or be absent even though `PRD.md`/`NOTES.md`/`PLAN.md` were already
produced (e.g. in a prior session, or authored by hand). Never trust file *existence* alone to
infer the phase — a `PLAN.md` sitting on disk with `Status: Planning` means the plan is still
awaiting dev approval, not that implementation may start. Run the detection script to compute the
furthest phase actually reachable from what is on disk:

```
node .agentic/skills/sdd-workflow/scripts/detect-state.cjs {feature_path}
```

The script reports `PHASE`, `PRD`, `NOTES`, `PLAN`, `PLAN_STATUS`, and — once a plan is approved —
`TOTAL_TASKS`/`COMPLETED_TASKS`/`PENDING_TASKS`/`CURRENT_WAVE` parsed directly from `PLAN.md`'s
Section 1 task checkboxes.

If `{feature_path}/STATE.md` does not yet exist, create the feature directory's `STATE.md` by
calling `write-state.cjs` with the fields the script reported (this both creates and populates it
in one call — the script initialises a fresh default `STATE.md` when none exists):

```
node .agentic/skills/sdd-workflow/scripts/write-state.cjs {feature_path} \
  --phase {PHASE} \
  --prd {PRD} \
  --notes {NOTES} \
  --plan {PLAN} \
  --current-wave {CURRENT_WAVE} \
  --completed "{COMPLETED_TASKS}" \
  --pending "{PENDING_TASKS}"
```

Omit any flag whose value the detection script reported as `(none)` or that wasn't reported at all
(e.g. `--current-wave`/`--completed`/`--pending` when phase is `prd` or `plan`). This call only
ever touches the fields named by the flags passed — `retry_counts`, `failed_tasks`, and
`current_wave_groups` are never touched here, so a pre-existing `STATE.md`'s memory-only fields
survive untouched automatically.

Then go directly to **Step 3 — Phase routing** using this freshly detected `phase` — do not stop
here even if this required fast-forwarding through multiple phases in one invocation.

## Step 3 — Phase routing

Follow the branch matching the `phase` just detected in Step 1.

---

### Phase: `prd`

Step 1 already confirmed no PRD artifact (`PRD.md` nor `PRD.cache.md`) exists yet for this feature
(that's the only way the script reports this phase). `sdd-prd` is a standalone tool for the
product team, not orchestrated by this pipeline — ask the user directly for the PRD source instead
of instructing them to run it:
```
Aguardando PRD.

Qual é a origem do PRD desta feature — caminho de um arquivo local (ex: docs/<module>/<feature>/PRD.md)
ou um link/referência de card (Jira/GitHub/Trello/etc.)?
```
- **Local path given**: confirm the file exists; if so, tell the user to invoke `/sdd-plan
  {path}` (a new session) to continue. If it doesn't exist yet, tell them to produce it first (with
  `sdd-prd` or by hand) and return once it's saved.
- **Card link given**: tell the user to invoke `/sdd-plan {link}` (a new session) — the plan
  skill fetches and normalizes the card content itself.
Stop.

---

### Phase: `plan`

Step 1 already confirmed either `PLAN.md` doesn't exist yet, or it exists but isn't approved
(`PLAN_STATUS` is `Draft`, `Planning`, or `Unknown`). Distinguish using the script's `PLAN` field:

**`PLAN: (none)`:**
```
PRD detectado: {PRD}

--- Avancando para fase: PLAN ---

Abra uma nova sessão e execute:

  /sdd-plan {PRD}

Quando o PLAN.md estiver aprovado e salvo, volte e execute /sdd-workflow novamente.
```

**`PLAN.md` exists but not approved:**
```
PLAN.md detectado, porém ainda não aprovado (Status: {PLAN_STATUS}).

Abra uma nova sessão e execute /sdd-plan {PRD} para revisar/aprovar o plano.
Quando o Status mudar para "Ready for Implementation", volte e execute /sdd-workflow novamente.
```
Stop.

`{PRD}` here is exactly the `PRD:` value `detect-state.cjs` reported in Step 1 — a local `PRD.md`
path, or a previously materialized `PRD.cache.md` path if this feature's plan was started from a
card. Never hardcode `{feature_path}/PRD.md`; the actual filename varies.

---

### Phase: `implement`

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

If `{target_wave}` is set, validate it before proceeding:
- Read the script's `WAVE` value as `current_wave`.
- If `{target_wave}` < `current_wave`: inform the user that the wave was already executed and stop.
- If `{target_wave}` > `current_wave` + gap of more than 1: warn that previous waves must complete
  first and stop.
- If `{target_wave}` == `current_wave` or `current_wave + 1`: proceed normally, but filter Step 4
  to only execute that wave.

Otherwise: go to **Step 4 — Execute wave**, using `compute-wave.cjs`'s output directly instead of
recomputing executable tasks by hand.

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

After presenting the summary above, run the PR gate (Step 3.1 below). Only once it resolves to
`created` or `skipped` does the archive step run — this runs on every arrival at `done`, even runs
with failed tasks (a failed task is still reported, not a reason to keep the spec's files in `docs/`
forever):

Read `sdd-workflow.config.json` from the repo root (if present) for its `archive_spec` section.
Default to `{destination_type: "mcp", destination_name: "obsidian"}` if the file, the section, or
either field is missing.

Load [archive-spec](.agentic/subagents/archive-spec.md) and inject:
- `{feature_path}` → the feature path for this pipeline
- `{prd_path}` → `STATE.md.artifacts.prd`
- `{plan_path}` → `STATE.md.artifacts.plan`
- `{notes_path}` → `STATE.md.artifacts.notes` (omit the line if empty)
- `{destination_type}` → `archive_spec.destination_type` (resolved above)
- `{destination_name}` → `archive_spec.destination_name` (resolved above)

Wait for its report, then:
- **`STATUS: archived`**: delete `{feature_path}/STATE.md` directly, then remove `{feature_path}`
  if it is now empty. Append: `Spec arquivada e removida de docs/.`
- **`STATUS: partial-failure`** or **`STATUS: no-destination-available`**: do not delete anything
  (the agent already left local files untouched). Append the agent's report verbatim so the user
  knows archiving didn't complete and why. `phase` stays `done` — this is not a pipeline failure.

---

### Step 3.1 — PR gate (runs before archive-spec)

Skip this entirely if `STATE.md.pr.status` is already `created` — proceed straight to the archive
step above. This makes the gate idempotent across re-runs (e.g. after installing `gh`).

1. Run:
   ```
   node .agentic/skills/sdd-workflow/scripts/prepare-pr.cjs {feature_path}
   ```
2. **`STATUS: blocked_missing_tool`**: run
   `write-state.cjs {feature_path} --pr-status blocked_missing_tool`, relay the script's `INSTALL`
   line to the user verbatim, and **stop the whole invocation** — do not spawn `archive-spec`. The
   next `/sdd-workflow {feature_path}` re-enters `phase: done`, re-runs this same step, and picks up
   automatically once the tool is installed.
3. **`STATUS: blocked_protected_branch`**: run
   `write-state.cjs {feature_path} --pr-status blocked_protected_branch`, relay the script's
   `MESSAGE` to the user, and stop the whole invocation — same resume behavior as above once the dev
   has switched to a feature branch.
4. **`STATUS: already_open`**: an open PR/MR already exists for this branch (resume case — a prior
   run got this far but failed before `archive-spec`). Run
   `write-state.cjs {feature_path} --pr-status created --pr-url {PR_URL} --pr-branch {BRANCH}`,
   append the URL to the summary already shown to the user, and proceed to `archive-spec` above.
5. **`STATUS: ready`**: branch/commit/push are done. Before spawning the writer, respect `{yolo}`
   (parsed in Step 0) exactly like the wave gate:
   - **`{yolo}` is `false`**: present `Branch: {BRANCH} -> {BASE_BRANCH}` and the task list, then wait
     for explicit confirmation (`sim`) before proceeding — pushing/opening a PR is visible to others
     and harder to reverse than a local wave.
   - **`{yolo}` is `true`**: skip confirmation, proceed immediately.

   Then load [pr-writer](.agentic/subagents/pr-writer.md) and inject:
   - `{provider}` → the script's `PROVIDER`
   - `{type}` → the script's `TYPE`
   - `{branch}` → the script's `BRANCH`
   - `{base_branch}` → the script's `BASE_BRANCH`
   - `{prd_path}` → the script's `PRD_PATH`
   - `{tasks_table}` → the script's `TASKS:` block, verbatim

   Wait for its report:
   - **`STATUS: created`**: run
     `write-state.cjs {feature_path} --pr-status created --pr-url {url} --pr-branch {branch}`, append
     the URL to the summary, then proceed to `archive-spec` above.
   - **`STATUS: failed`**: run `write-state.cjs {feature_path} --pr-status failed --pr-branch {branch}`,
     report the failure to the user, and stop — do not spawn `archive-spec`. Branch/commit/push
     already happened, so nothing is lost; a later re-run retries PR creation only.

---

## Step 4 — Execute wave

### 4.1 — Compute wave groups

Use the wave groups and task table already produced by `compute-wave.cjs` in the `implement`
phase branch above — do not recompute them by hand.

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

**If `{yolo}` is `true`:** skip confirmation entirely. Output instead:

```
[YOLO] Iniciando Wave {N} automaticamente — sem confirmação.
```

Then spawn immediately.

### 4.2 — Spawn one batched subagent per group; groups run sequentially

Process groups **sequentially** (group 1 fully resolves before group 2 starts).
Each group is handled by **one implement-agent and one verify-code agent, batched over every task
in that group** — not one agent per task. This amortizes the fixed cost of loading `AGENTS.md` and
the `ddd-module-knowledge`/`adjust-go-code`/`make-unit-tests` skills (previously paid once per task)
across every task in the group, which is what actually dominated session budget on waves with many
small tasks — the 200k context window itself was never the cost, the repeated bootstrap was.
Retry cycles run **sequentially**, scoped to the individual failed task, after the group's verdict
table is consolidated.

#### Step A — Implement (one batched agent per group)

Before spawning, compute `{has_dedicated_test_task}` for every task in the group **once**, without
loading PLAN.md in full — grep just the `TASK-TEST-*` blocks' `File:` lines and match them against
this group's task files:

```
grep -A 2 '^\- \[.\] \*\*TASK-TEST-' {plan_path} | grep 'File:'
```

A task's `{has_dedicated_test_task}` is `true` if its own `File:` (from `compute-wave.cjs`'s table)
appears in that output, else `false`. Build the `{tasks_table}` from this (one line per task:
`{task_id} — has_dedicated_test_task: {true|false}`) — this single grep replaces each task
scanning PLAN.md itself for this check.

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

All tasks in the group have finished writing by this point — the code is in one stable state
until the next group starts. Run the suite **once** for the whole group here, not once per task
in Step B: every verify-code in the group would otherwise check the exact same code state
redundantly.

```
node .agentic/skills/sdd-workflow/scripts/run-build-test.cjs
```

The script returns a compact structured summary (never raw logs, so this orchestrator's own
context stays free of build/test output). Capture its stdout as `{build_result}` — pass it
verbatim into every verify-code spawned in Step B; none of them re-run these commands
themselves.

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
     Inform the user.

Only after all retries for the group are resolved, proceed to the next group.

#### Step C.5 — Wave PRD gate (once per wave, after every group resolves)

Runs once, after **every** group in the wave has finished Step C — not per group. Rationale: PRD
acceptance criteria are written at the feature/flow level, not the file level, so checking them
once per task (the old design) produced redundant PRD reads and low-signal verdicts for tasks with
no observable behavior in isolation — see `verify-wave-prd.md`.

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
     — and warn the user which later, currently-completed tasks import or depend on the file being
     reopened (a blast-radius notice, not an automatic re-verification of those tasks — an
     automatic cascading re-verification could turn one PRD finding into a full pipeline re-run,
     which needs explicit dev sign-off, not a silent default). Then apply the same increment-retry
     / re-spawn-implement-agent / re-verify sequence as Step C.
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
Then refresh the wave cache — never assume the next wave is `{N+1}` by arithmetic; re-run
`compute-wave.cjs` and cache whatever `WAVE` it reports (wave numbers are not guaranteed dense, and
this is the same script `Step 3 — implement` will call next anyway, so this is just priming the cache
STATE.md keeps for display, not a computation you do by hand):
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
  Do not assume the next wave is `{N+1}` in this message — use the exact value `compute-wave.cjs` just
  reported (wave numbers are computed from the plan's dependency graph and are not guaranteed to be
  dense).
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
  - Proceed to Phase `done` (Step 3).
- If there are still pending tasks with unsatisfied dependencies (deadlock), report:
  ```
  Deadlock detectado: tasks pendentes têm dependências não satisfeitas.
  Pendentes: {list}
  Concluídas: {list}
  Verifique o PLAN.md e resolva manualmente.
  ```

## Gotchas

- Never advance a phase without verifying the expected artifact exists on disk — this is exactly
  what `detect-state.cjs` does; don't hand-roll existence checks that skip it.
- Files on disk (`PRD.md` or `PRD.cache.md`/`NOTES.md`/`PLAN.md` + its `Status` header and task
  checkboxes) are the source of truth for `phase`/`artifacts`/`completed_tasks`/`pending_tasks`/
  `current_wave`.
  `STATE.md` only retains what files can't encode: `retry_counts`, `failed_tasks`,
  `current_wave_groups`. Do not rely on conversation memory for any of it.
- A `PLAN.md` on disk does not imply the `implement` phase — check `PLAN_STATUS` from the
  detection script. Only `Ready for Implementation` (or later) allows advancing past `plan`.
- Never hand-edit `STATE.md` or `PLAN.md` checkboxes/Status directly — always go through
  `write-state.cjs` / `update-plan.cjs`, even for a single-field change. They're the only things
  that guarantee `retry_counts`/`failed_tasks`/`current_wave_groups` survive an unrelated update.
- `parallel_group` controls execution order **within** a wave: every task in the same group is
  handled by a single batched implement-agent / verify-code agent (not one agent per task); groups
  run sequentially. `wave` controls the coarser boundary — a task never executes before every task
  in an earlier wave is complete, and `implement-wave=N` never spills into a different wave's
  groups. Do not confuse the two: a single wave commonly contains several groups (e.g. wave 3 =
  `application, infrastructure, tests` all at once, because they all depend only on wave 1/2 and
  nothing blocks them running together), which is expected — what must never happen is a task from
  a *different, higher* wave slipping into the current one. Since tasks in the same wave can never
  depend on each other (a task's wave is always strictly greater than every one of its
  dependencies'), groups within a wave are independent of each other too — they are processed
  sequentially by choice (to bound concurrent agent/session load without worktree isolation), not
  because of a correctness requirement.
- If the user invokes `/sdd-workflow` mid-wave (interrupted session), re-present the current wave
  and ask whether to re-execute or skip.
- If `compute-wave.cjs`/`detect-state.cjs` exit 2 reporting a task with no `wave:` field, the PLAN.md
  predates wave computation (or a task was hand-added after `compute-waves.cjs` last ran). Do not
  estimate a wave number yourself — tell the user to open a plan session and run
  `sdd-plan`'s `compute-waves.cjs` on the plan, then retry.
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
- Phase `done` always attempts archiving (`archive-spec` agent), even when `failed_tasks` is
  non-empty. Only `STATE.md` is ever discarded unarchived — whatever `STATE.md.artifacts.prd` points
  to (`PRD.md` or a materialized `PRD.cache.md`), plus `PLAN.md`/`NOTES.md`, are only deleted
  locally after the agent confirms they were written to the destination. `archive-spec` treats
  `PRD.cache.md` exactly like `PRD.md` here — it reads whatever path `STATE.md.artifacts.prd` names
  and does not need to know it's a temporary, derived file. A failed archive
  attempt never blocks or reverts the pipeline result already reported to the user.
