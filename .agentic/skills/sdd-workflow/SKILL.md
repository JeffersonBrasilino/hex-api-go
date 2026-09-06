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

Each phase's instructions live in their own reference file, loaded only once the phase is known —
load just the one file that matches the detected phase, never the others.

## Scripts disponíveis

| Script | Quando usar | exit 0 (stdout) | exit 1 (stderr) |
|--------|-------------|-----------------|-----------------|
| `detect-state.cjs` | Step 1 — detecção de fase a partir dos artefatos em disco | `PHASE` + artefatos + tasks | — |
| `write-state.cjs` | gravar/atualizar STATE.md (todas as fases) | Campos atualizados vs. preservados | Fase/campo inválido |
| `compute-wave.cjs` | fase `implement` — computar a próxima wave executável | `WAVE`/`TOTAL_WAVES` (do campo `wave:` persistido) + tabela de tasks | Deadlock detectado |
| `update-plan.cjs` | fase `implement` — marcar checkboxes e Status do PLAN.md | Tasks marcadas / já corretas | TASK-ID não encontrado / Status inválido |

Todos os scripts suportam `--help`. Exit code `2` = erro de uso/runtime (argumento faltando,
arquivo não encontrado) — trate como um bloqueio a reportar ao usuário, não como um bug a corrigir
silenciosamente.

## Principles

- **Detect from disk, not from memory of the last run.** `STATE.md` is a cache of what
  `detect-state.cjs` found — never the source of truth by itself. Re-derive the phase from
  `PRD.md`/`PRD.cache.md`/`NOTES.md`/`PLAN.md` every invocation; `STATE.md` only backs the fields
  files can't encode (`retry_counts`, `failed_tasks`, `current_wave_groups`). Never rely on
  conversation memory for any of it.
- **One phase-appropriate action per invocation, but fast-forward to reach it.** If artifacts for
  several phases already exist on disk, route through all of them in this same invocation — don't
  make the user re-run `/sdd-workflow` just to notice a file already on disk. Stop only once you
  reach a phase whose required artifact is genuinely missing, or an unapproved plan.
- **Deterministic work goes through scripts, not reasoning.** Computing executable tasks, writing
  STATE.md fields, and marking PLAN.md checkboxes/Status are fixed-schema operations — never
  hand-roll them, even for a single-field change; call the corresponding script and use its stdout.
- **Wave numbers come from the plan, never from a counter.** `compute-wave.cjs` reads the
  `**Wave:**` field `sdd-plan`'s `compute-waves.cjs` wrote at plan time — never increment a session
  counter or re-derive boundaries from `depends_on`. If a `PLAN.md` predates this (no `wave:`
  field), the scripts exit 2 with instructions to run `compute-waves.cjs` first — don't estimate
  wave boundaries yourself.
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

`STATE.md` can go stale or be absent even though `PRD.md`/`NOTES.md`/`PLAN.md` already exist. Never
trust file *existence* alone to infer the phase — e.g. a `PLAN.md` with `Status: Planning` means
implementation may not start yet. Run the detection script to compute the furthest phase actually
reachable from what is on disk:

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

Omit any flag whose value the detection script reported as `(none)` or that wasn't reported at all.
This call only touches the fields named by the flags passed — `retry_counts`, `failed_tasks`, and
`current_wave_groups` survive untouched automatically.

Then go directly to **Step 3 — Phase routing** using this freshly detected `phase` — do not stop
here even if this required fast-forwarding through multiple phases in one invocation.

## Step 3 — Phase routing

Load only the reference file matching the `phase` just detected in Step 1 and follow it:

| Phase | Reference |
|-------|-----------|
| `prd` | [references/phase-prd.md](references/phase-prd.md) |
| `plan` | [references/phase-plan.md](references/phase-plan.md) |
| `implement` or `verify` | [references/phase-implement.md](references/phase-implement.md) |
| `done` | [references/phase-done.md](references/phase-done.md) |

`references/phase-implement.md` covers both `implement` (Step 4 — execute wave) and `verify`
(Step 5 — gate-only check), since `verify` is only ever reached from inside that same execution
loop. `references/phase-done.md` covers the `done` phase, the PR gate, and archiving.
