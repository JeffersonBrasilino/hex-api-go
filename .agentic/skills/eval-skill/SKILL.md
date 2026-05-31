---
name: eval-skill
description: >
  Run structured evaluations (evals) for any skill that has an evals/evals.json file.
  Executes each test case twice — once with the skill loaded (with_skill) and once without
  it (without_skill baseline) — then grades outputs against assertions via a dedicated haiku
  grader subagent, generates per-eval grading.json, timing.json, benchmark.json, and a
  feedback.json for human review. Grading and reporting are offloaded to subagents to keep
  the main context lean.
  Use when the user says "eval this skill", "run evals for", "test this skill", or provides
  a skill path and asks to evaluate it. Requires the target skill to have evals/evals.json.
---

# eval-skill

Runs the agentskills.io eval workflow for any skill. Produces a with/without baseline
comparison, per-assertion grading, benchmark delta, and a feedback file for human review
and skill iteration.

Grading and report generation are delegated to dedicated subagents to minimize tokens
consumed in the main context.

## Language

All output — grading evidence, benchmark analysis, report content, feedback files, and the final message — must be written in **Brazilian Portuguese (pt-BR)**. This applies to every file produced by this skill and to every text response. Technical terms (PASS, FAIL, JSON keys, file paths, code identifiers) are kept in their original form.

## Scope

**This skill covers:**
- Reading `evals/evals.json` from the target skill directory
- Optionally filtering evals by ID before running
- Running with_skill and without_skill passes for all selected test cases
- Grading each assertion via a dedicated haiku subagent (offloaded from main context)
- Writing per-eval `grading.json` and `timing.json` per configuration
- Writing `benchmark.json`, `feedback.json`, and `report.md` via a dedicated subagent

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
            │   │   ├── timing.json    ← { "total_tokens": N, "duration_ms": N }
            │   │   └── grading.json   ← assertion results for this eval
            │   └── without_skill/
            │       ├── outputs/
            │       ├── timing.json
            │       └── grading.json
            ├── eval-{slug}/
            │   └── ...
            ├── benchmark.json          ← aggregated statistics across all evals
            ├── feedback.json           ← human review notes per eval slug
            └── report.md              ← summary report in pt-BR
```

**Eval slug**: derive from the `description` field (lowercase, spaces and special chars → hyphens, max 60 chars). Fall back to `eval-{id}` if no description. Examples:
- `"Module structure — verify layout"` → `eval-module-structure-verify-layout`
- eval with id 3 and no description → `eval-3`

## Workflow

### PHASE 1 — Load and Validate

1. Read `{skill-path}/SKILL.md`
2. Read `{skill-path}/evals/evals.json`
3. Validate: `evals` array non-empty; each entry has `id`, `prompt`, `assertions`
4. If `--ids` was provided, filter the `evals` array to only entries whose `id` is in the list. If none match, stop and report.
5. If invalid, stop and report what is missing
6. Read all files under `{skill-path}/references/` — injected into the with_skill run
7. Compute workspace root: `{skill-path}/evals/workspace/`
8. Determine iteration: count existing `iteration-N/` dirs in workspace root, use N+1 (start at 1)
9. Create directory tree using the Write tool on placeholder files — never use shell mkdir:
   - For each eval: `{workspace}/iteration-{N}/eval-{slug}/with_skill/outputs/.keep`
   - For each eval: `{workspace}/iteration-{N}/eval-{slug}/without_skill/outputs/.keep`

### PHASE 2 — Execute Runner Subagents (parallel)

**File writing rule**: always use the dedicated Write tool to create files. Never use shell scripts, Python scripts, bash heredocs, or any other scripting mechanism to write files.

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

Compute `estimated_cost_usd` using the pricing table in `references/pricing-table.md`.

For each eval write:
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/outputs/output.json` → `{ "id": N, "answer": "..." }`
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/outputs/output.json` → same
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/timing.json` — see [timing schema](references/grading-schema.md), add `"note": "batch run — timing shared across {K} evals"`
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/timing.json` → same

### PHASE 3 — Grade via Haiku Subagent

Do NOT grade assertions in the main context. Spawn a **single grader subagent using `claude-haiku-4-5-20251001`** that receives all answers and assertions for both configurations and writes all `grading.json` files directly.

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
- `pass_rate`: passed / total rounded to 3 decimal places.

For each eval and each configuration, write the grading.json file at the path specified.

=== GRADING SCHEMA ===
{full content of references/grading-schema.md}

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

After the grader subagent completes, verify that all `grading.json` files exist. If any are missing (subagent denied write permissions), read the grader's output and write the missing files yourself.

### PHASE 4 — Benchmark and Report via Subagent

Do NOT build the benchmark or report in the main context. Spawn a **single benchmark subagent** that reads all grading and timing files and writes `benchmark.json`, `feedback.json`, and `report.md`.

Build the benchmark prompt as follows:

```
You are an eval benchmark aggregator. Read the grading and timing results below and produce
benchmark.json, feedback.json, and report.md for iteration {N} of skill "{skill-name}".

All written content (report.md, feedback values if non-empty) must be in Brazilian Portuguese (pt-BR).
JSON keys, file paths, PASS/FAIL labels, and code identifiers stay in their original form.

=== BENCHMARK SCHEMA ===
{full content of references/benchmark-schema.md}

=== FEEDBACK SCHEMA ===
{full content of references/feedback-schema.md}

=== PRICING TABLE ===
{full content of references/pricing-table.md}

{if N > 1:
=== PREVIOUS ITERATION BENCHMARK ===
{content of iteration-(N-1)/benchmark.json}
}

=== GRADING AND TIMING DATA ===
{for each eval:
--- {slug} ---
with_skill/grading.json:
{content}

without_skill/grading.json:
{content}

with_skill/timing.json:
{content}

without_skill/timing.json:
{content}
}

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

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
...

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

=== OUTPUT PATHS ===
Write these files using the Write tool:
- {workspace}/iteration-{N}/benchmark.json
- {workspace}/iteration-{N}/feedback.json   (pre-populated with empty strings per slug)
- {workspace}/iteration-{N}/report.md
```

After the benchmark subagent completes, output **only** this message:

```
o relatório da avaliação da skill foi criado, local: {workspace}/iteration-{N}/report.md
```

## Gotchas

- **Never grade or report in main context** — always delegate to subagents (Phases 3 and 4); doing it inline defeats the token-saving design
- **Grader uses haiku** (`claude-haiku-4-5-20251001`) — grading is mechanical and does not need a reasoning model
- **Eval filtering** (`--ids`) is applied before spawning runner subagents — runners only receive the selected evals
- **Never use scripts** (Python, bash, shell heredocs) to write files — use the Write tool exclusively
- All written content must be in **pt-BR**; only JSON keys, file paths, code identifiers, PASS/FAIL labels stay in their original form
- The report goes to `report.md` inside the iteration dir — never printed to the conversation prompt
- Always read ALL `references/` files before building the with_skill prompt — critical conventions often live in references/, not SKILL.md body
- Eval slug must be filesystem-safe: lowercase, hyphens only, max 60 chars — truncate if needed
- Iteration number comes from filesystem scan — never assume 1 if `{workspace}/` already exists
- Capture `model`, `input_tokens`, `output_tokens`, and `duration_ms` immediately from task completion notification — they are not available later
- Token split may not always be present; apply the 75/25 fallback from pricing-table.md and record `pricing_note`
- `value_tier` is derived solely from `delta.pass_rate` — do not adjust based on cost or context
- `delta_vs_prev_iteration` requires reading the previous iteration's `benchmark.json` — skip only if iteration 1
- If a subagent is denied write permissions, answers still arrive in task result — write all files yourself
- Negative assertions are highest-signal — grade strictly, they catch regressions
- Baseline-confirmed evals are not waste — they validate no regression in universal knowledge

## References

- Grading output schema → see [reference](references/grading-schema.md)
- Benchmark output schema → see [reference](references/benchmark-schema.md)
- Feedback file schema → see [reference](references/feedback-schema.md)
- Model pricing table → see [reference](references/pricing-table.md)
