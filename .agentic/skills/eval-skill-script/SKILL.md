---
name: eval-skill-script
description: >
  Run structured evaluations (evals) for any skill that has an evals/evals.json file.
  Executes each test case twice — once with the skill loaded (with_skill) and once without
  it (without_skill baseline) — then grades outputs against assertions via a dedicated haiku
  grader subagent, computes cost/timing/benchmark aggregation via scripts, and generates a
  lightweight report.md via a dedicated subagent. Grading is offloaded to a subagent to keep
  the main context lean; arithmetic and schema-fixed file writes are offloaded to scripts.
  Use when the user says "eval this skill", "run evals for", "test this skill", or provides
  a skill path and asks to evaluate it. Requires the target skill to have evals/evals.json.
allowed-tools: [Bash, Read, Write, Agent]
---

# eval-skill-script

Runs the agentskills.io eval workflow for any skill. Produces a with/without baseline
comparison, per-assertion grading, benchmark delta, and a feedback file for human review
and skill iteration.

Grading and report prose are delegated to dedicated subagents to minimize tokens consumed
in the main context. Arithmetic (cost, timing, statistical aggregation) and schema-fixed
file writes (`timing.json`, `grading.json` summaries, `benchmark.json`, `feedback.json`)
are delegated to Node.js scripts instead of subagents — see "Scripts disponíveis" below.

## Language

All output — grading evidence, report content, feedback files, and the final message — must be written in **Brazilian Portuguese (pt-BR)**. This applies to every file produced by this skill and to every text response. Technical terms (PASS, FAIL, JSON keys, file paths, code identifiers) are kept in their original form.

## Scripts disponíveis

| Script | Quando usar | exit 0 (stdout) | exit 1 (stderr) |
|--------|-------------|-----------------|-----------------|
| `scaffold-eval-workspace.cjs` | Fase 1 — validar, filtrar `--ids`, derivar slugs, criar workspace | Iteração + lista de evals | evals.json inválido ou `--ids` sem match |
| `write-timing.cjs` | Fase 2 — após cada subagente runner concluir | timing.json + custo | — (erros são de uso, exit 2) |
| `finalize-grading.cjs` | Fase 3 — após o subagente grader escrever os `grading.json` | Contagem atualizado/inalterado | grading.json com shape inválido |
| `compute-benchmark.cjs` | Fase 4 — antes de spawnar o subagente de relatório | benchmark.json + feedback.json | grading.json sem `summary` |

Todos os scripts suportam `--help`. Exit code `2` = erro de uso/runtime.

## Scope

**This skill covers:**
- Reading `evals/evals.json` from the target skill directory
- Optionally filtering evals by ID before running
- Running with_skill and without_skill passes for all selected test cases
- Grading each assertion via a dedicated haiku subagent (offloaded from main context)
- Finalizing `grading.json` summaries and writing `timing.json` via scripts
- Computing `benchmark.json` and seeding `feedback.json` via script
- Writing `report.md` via a lightweight dedicated subagent

**This skill does NOT cover:**
- Creating or editing `evals.json` — do that separately before running evals
- Fixing the target skill based on results — review the report and iterate manually
- Running evals for skills without an `evals/evals.json` file

## Input

The user must provide the **skill path** — the directory containing the target `SKILL.md`. Examples:
- `.agentic/skills/ddd-module-knowledge`
- `/absolute/path/to/my-skill`

Resolve relative paths from the repository root.

**Optional parameters:**
- `--ids N,M,...` — run only evals with these IDs (e.g., `--ids 1,3,5`). If omitted, all evals run.
- `--model <model-id>` — model for the runner subagents. If omitted, use the current conversation model. The model ID is recorded in `timing.json` and used for cost calculation.

## Workspace Layout

The workspace lives **inside** the skill directory, under `evals/workspace/`. Each iteration gets its own `iteration-N/` directory:

```
my-skill/
├── SKILL.md
└── evals/
    ├── evals.json
    └── workspace/
        └── iteration-1/
            ├── eval-{slug}/
            │   ├── with_skill/
            │   │   ├── outputs/       ← text or files produced by the run
            │   │   ├── timing.json    ← written by write-timing.cjs
            │   │   └── grading.json   ← assertion results, finalized by finalize-grading.cjs
            │   └── without_skill/
            │       ├── outputs/
            │       ├── timing.json
            │       └── grading.json
            ├── eval-{slug}/
            │   └── ...
            ├── benchmark.json          ← written by compute-benchmark.cjs
            ├── feedback.json           ← seeded by compute-benchmark.cjs, filled by human review
            └── report.md              ← summary report in pt-BR, written by report subagent
```

**Eval slug**: derived by `scaffold-eval-workspace.cjs` from the `description` field (lowercase, spaces and special chars → hyphens, max 60 chars). Falls back to the eval's `id` if no description.

## Workflow

### PHASE 1 — Scaffold Workspace

```bash
node .agentic/skills/eval-skill-script/scripts/scaffold-eval-workspace.cjs {skill_path} \
  [--ids {ids}]
```

This validates `evals.json`, applies `--ids` filtering, derives slugs, detects the next
iteration number, and creates the directory tree. Use its stdout (`ITERATION`, `WORKSPACE`,
the `id | eval-slug | path` table) directly to build the runner subagent prompts in Phase 2.

If exit 1 (evals.json malformed or `--ids` matches nothing): report the error to the user and stop.
If exit 2 (skill path or evals.json not found): report the error to the user and stop.

Also read all files under `{skill_path}/references/` — injected into the with_skill run.

### PHASE 2 — Execute Runner Subagents (parallel)

**File writing rule**: always use the dedicated Write tool to create output files. Never use shell scripts, Python scripts, bash heredocs, or any other scripting mechanism to write files.

Spawn **two subagents in parallel** — one per configuration — each receiving all selected eval prompts at once. Capture timing immediately from the task completion notification.

**with_skill agent** prompt:
```
You are an eval runner. Answer each question using ONLY the skill knowledge provided below.
Where the skill is explicit, do not override with general knowledge.

=== SKILL CONTENT ===
{full content of SKILL.md}

=== REFERENCES ===
{for each file in references/: --- {filename} ---\n{content}\n}

For each question, write { "id": N, "answer": "...", "slug": "{eval-slug}" }.
Return a JSON array of all answers.

Questions:
{for each eval: N. [slug: {slug}] {prompt}}
```

**without_skill agent** prompt:
```
You are an eval runner. Answer each question using ONLY your base domain knowledge.
Do NOT read any project files or skill files.

For each question, write { "id": N, "answer": "...", "slug": "{eval-slug}" }.
Return a JSON array of all answers.

Questions:
{for each eval: N. [slug: {slug}] {prompt}}
```

After each subagent completes, immediately capture from the task completion notification:
- `model` — model ID string
- `input_tokens` and `output_tokens` — capture separately when available; fall back to `total_tokens` only if the split is absent
- `duration_ms`

For each eval write:
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/outputs/output.json` → `{ "id": N, "answer": "..." }`
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/outputs/output.json` → same

Then call, for each eval directory, once per configuration:

```bash
node .agentic/skills/eval-skill-script/scripts/write-timing.cjs \
  {workspace}/iteration-{N}/eval-{slug}/with_skill/timing.json \
  --model {model} --input-tokens {input} --output-tokens {output} --duration-ms {duration} \
  --note "batch run — timing shared across {K} evals"
```

(repeat with `without_skill/timing.json` and the without_skill run's values). The script
computes `estimated_cost_usd` from the pricing table internally — no need to read
`references/pricing-table.md` for this.

### PHASE 3 — Grade via Haiku Subagent, then Finalize

Do NOT grade assertions in the main context. Spawn a **single grader subagent using `claude-haiku-4-5-20251001`** that receives all answers and assertions for both configurations and writes the assertion-level results directly.

Build the grader prompt as follows:

```
You are a strict eval grader. For each eval, grade every assertion against the given answer.

Grading rules:
- PASS: assertion clearly satisfied — provide the exact quote or observation as evidence.
- FAIL: assertion not satisfied — state precisely what was missing or wrong.
- Require concrete evidence for PASS — no benefit of the doubt.
- Negative assertions ("does NOT suggest X"): PASS only if X is genuinely absent.
- Redirect assertions ("redirects to skill Y"): PASS only if the skill name is explicitly present.
- Grade on substance only — not tone, style, or length.
- `text`: copy the assertion verbatim.

For each eval and each configuration, write a JSON file with exactly this shape:

{ "assertion_results": [ { "text": "...", "passed": true|false, "evidence": "..." }, ... ] }

Do NOT compute a summary field — that is done separately.

=== EVALS TO GRADE ===
{for each eval:
--- EVAL: {slug} ---
Assertions:
{JSON array of assertions}

with_skill answer:
{answer from with_skill output.json}

without_skill answer:
{answer from without_skill output.json}

Write to:
  {workspace}/iteration-{N}/eval-{slug}/with_skill/grading.json
  {workspace}/iteration-{N}/eval-{slug}/without_skill/grading.json
}
```

After the grader subagent completes, verify that all `grading.json` files exist. If any are missing (subagent denied write permissions), read the grader's output and write the missing files yourself — using the same `{ "assertion_results": [...] }` shape, no summary.

Then run:

```bash
node .agentic/skills/eval-skill-script/scripts/finalize-grading.cjs {workspace}/iteration-{N}
```

This computes and writes `summary.{passed,failed,total,pass_rate}` into every `grading.json`.
If exit 1 (invalid shape): identify the file named in the error and have the grader subagent
(or you, if the task is gone) rewrite it with the correct shape, then re-run the script.

### PHASE 4 — Benchmark via Script, Report via Lightweight Subagent

```bash
node .agentic/skills/eval-skill-script/scripts/compute-benchmark.cjs {workspace}/iteration-{N}
```

This computes all statistics (`pass_rate`/`time_seconds`/`tokens`/`estimated_cost_usd`
mean/stddev/delta, `value_tier`, `delta_vs_prev_iteration` when a prior iteration exists) and
writes `benchmark.json`. It also seeds `feedback.json` with empty strings per slug — but only
if `feedback.json` does not already exist, so a prior human review is never overwritten.

If exit 1 (a `grading.json` has no `summary`): run `finalize-grading.cjs` for this iteration
first, then retry.

Then spawn a **single, lightweight report subagent**. Unlike the raw benchmark data, it does
NOT need `benchmark-schema.md`, `pricing-table.md`, or `feedback-schema.md` inlined — those
files describe logic that now lives in the scripts. Give it only:
- the content of `{workspace}/iteration-{N}/benchmark.json` (already computed)
- for each eval, only the assertions from `grading.json` **that failed** (with_skill config) —
  not the full assertion set
- the content of `{workspace}/iteration-{N}/feedback.json`, for context if a prior review exists
- the report structure below

Build the report prompt as follows:

```
You are an eval report writer. Using the already-computed benchmark data and the failing
assertions below, write report.md for iteration {N} of skill "{skill-name}".

All written content must be in Brazilian Portuguese (pt-BR). JSON keys, file paths,
PASS/FAIL labels, and code identifiers stay in their original form.

=== BENCHMARK (already computed) ===
{content of benchmark.json}

{if N > 1: benchmark.json already includes delta_vs_prev_iteration — reference it in the report.}

=== FAILING ASSERTIONS (with_skill config only) ===
{for each eval with at least one failed assertion:
--- {slug} ---
{failed assertion text + evidence}
}

=== FEEDBACK (prior human review, if any) ===
{content of feedback.json}

=== REPORT STRUCTURE ===
Write report.md with the following structure (in pt-BR):

# Relatório de Avaliação: `{skill-name}` — iteração {N}

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | X.XXX |
| without_skill | X.XXX |
| **delta** | **+X.XXX** |
| **value_tier** | **forte / moderado / fraco / sem_valor / negativo** |

_(Se N > 1)_ Comparado à iteração anterior: delta passou de X.XXX → X.XXX (`{value_tier_change}`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $X.XXXXXX | $X.XXXXXX | {model} |
| without_skill | $X.XXXXXX | $X.XXXXXX | {model} |
| **custo adicional da skill** | **$X.XXXXXX** | — | — |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
...

## Baseline confirmado (ambos ≥ 0.95)

- `slug`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
...

## Recomendação

Uma frase — pronta / precisa de iteração / precisa de revisão maior.

=== OUTPUT PATH ===
Write this file using the Write tool:
- {workspace}/iteration-{N}/report.md
```

After the report subagent completes, output **only** this message:

```
o relatório da avaliação da skill foi criado, local: {workspace}/iteration-{N}/report.md
```

## Gotchas

- **Never grade or compute statistics in main context** — grading is a subagent (Phase 3), all arithmetic is a script (Phases 2 and 4); doing either inline defeats the token-saving design
- **Grader uses haiku** (`claude-haiku-4-5-20251001`) and only outputs `assertion_results` — it never computes `summary`; that's `finalize-grading.cjs`'s job
- **Eval filtering** (`--ids`) is applied by `scaffold-eval-workspace.cjs` before spawning runner subagents — runners only receive the selected evals
- **Never use scripts** (Python, bash, shell heredocs) to write *output* files (`output.json`) — use the Write tool exclusively; the four `.cjs` scripts in this skill are the only sanctioned scripted writes, and only for the fixed-schema files they own (`timing.json`, `grading.json` summary, `benchmark.json`, `feedback.json`)
- All written content must be in **pt-BR**; only JSON keys, file paths, code identifiers, PASS/FAIL labels stay in their original form
- The report goes to `report.md` inside the iteration dir — never printed to the conversation prompt
- Always read ALL `references/` files before building the with_skill prompt — critical conventions often live in references/, not SKILL.md body
- Iteration number comes from `scaffold-eval-workspace.cjs`'s filesystem scan — never assume 1 if `{workspace}/` already exists
- Capture `model`, `input_tokens`, `output_tokens`, and `duration_ms` immediately from task completion notification — `write-timing.cjs` needs them and they are not available later
- `finalize-grading.cjs` must run before `compute-benchmark.cjs` — the latter fails with exit 1 if any `grading.json` lacks a `summary`
- `feedback.json` is never overwritten by `compute-benchmark.cjs` once it exists — it may hold human review
- `value_tier` is derived solely from `delta.pass_rate` — the script does not adjust based on cost or context
- If a subagent is denied write permissions, answers/gradings still arrive in the task result — write the missing files yourself in the exact shape each script expects, then re-run the script
- Negative assertions are highest-signal — grade strictly, they catch regressions
- Baseline-confirmed evals are not waste — they validate no regression in universal knowledge

## References

- Grading output schema → see [reference](references/grading-schema.md) (documents the shape `finalize-grading.cjs` enforces)
- Benchmark output schema → see [reference](references/benchmark-schema.md) (documents the shape `compute-benchmark.cjs` produces)
- Feedback file schema → see [reference](references/feedback-schema.md) (documents the shape `compute-benchmark.cjs` seeds)
- Model pricing table → see [reference](references/pricing-table.md) (mirrored inside `write-timing.cjs` — keep both in sync if prices change)
