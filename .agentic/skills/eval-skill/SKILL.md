---
name: eval-skill
description: >
  Run structured evaluations (evals) for any skill that has an evals/evals.json file.
  Executes each test case twice — once with the skill loaded (with_skill) and once without
  it (without_skill baseline) — then grades outputs against assertions, generates per-eval
  grading.json, timing.json, benchmark.json, and a feedback.json for human review.
  Use when the user says "eval this skill", "run evals for", "test this skill", or provides
  a skill path and asks to evaluate it. Requires the target skill to have evals/evals.json.

execution_profile:
  tier: reasoning
  strategy: agent
---

# eval-skill

Runs the agentskills.io eval workflow for any skill. Produces a with/without baseline
comparison, per-assertion grading, benchmark delta, and a feedback file for human review
and skill iteration.

## Language

All output — grading evidence, benchmark analysis, report content, feedback files, and the final message — must be written in **Brazilian Portuguese (pt-BR)**. This applies to every file produced by this skill and to every text response. Technical terms (PASS, FAIL, JSON keys, file paths, code identifiers) are kept in their original form.

## Scope

**This skill covers:**
- Reading `evals/evals.json` from the target skill directory
- Running with_skill and without_skill passes for all test cases
- Grading each assertion with PASS/FAIL + concrete evidence
- Writing per-eval `grading.json` and `timing.json` per configuration
- Writing `benchmark.json` and `feedback.json` at the iteration root
- Producing a structured gap analysis to guide skill iteration

**This skill does NOT cover:**
- Creating or editing `evals.json` — do that separately before running evals
- Fixing the target skill based on results — review the report and iterate manually
- Running evals for skills without an `evals/evals.json` file

## Input

The user must provide the **skill path** — the directory containing the target `SKILL.md`. Examples:
- `.agentic/skills/ddd-module-knowledge`
- `/absolute/path/to/my-skill`

Resolve relative paths from the repository root.

**Optional**: the user may specify a **model** for the subagent runs (e.g., "eval using haiku", "use opus"). If not specified, use the current conversation model. The model ID is recorded in `timing.json` and used for cost calculation.

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
            └── feedback.json           ← human review notes per eval slug
```

**Eval slug**: derive from the `description` field (lowercase, spaces and special chars → hyphens, max 60 chars). Fall back to `eval-{id}` if no description. Examples:
- `"Module structure — verify layout"` → `eval-module-structure-verify-layout`
- eval with id 3 and no description → `eval-3`

## Workflow

### PHASE 1 — Load and Validate

1. Read `{skill-path}/SKILL.md`
2. Read `{skill-path}/evals/evals.json`
3. Validate: `evals` array non-empty; each entry has `id`, `prompt`, `assertions`
4. If invalid, stop and report what is missing
5. Read all files under `{skill-path}/references/` — this is the full knowledge injected into the with_skill run
6. Compute workspace root: `{skill-path}/evals/workspace/`
7. Determine iteration: count existing `iteration-N/` dirs in workspace root, use N+1 (start at 1)
8. Create directory tree:
   - `{workspace}/iteration-{N}/`
   - `{workspace}/iteration-{N}/eval-{slug}/with_skill/outputs/` for each eval
   - `{workspace}/iteration-{N}/eval-{slug}/without_skill/outputs/` for each eval

### PHASE 2 — Execute Runs (parallel)

**File writing rule**: always use the dedicated Write tool to create files. Never use shell scripts, Python scripts, bash heredocs, or any other scripting mechanism to write files. The Write tool is faster, does not require a shell process, and keeps the run time low.

Each eval should run in a **clean context** — no state from previous evals. In Claude Code, subagents provide this isolation naturally. Spawn **two subagents in parallel** — one per configuration — each receiving all eval prompts at once to minimise subagent count while preserving config isolation.

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

After each subagent completes, immediately capture from the task completion notification — these are not persisted anywhere else:
- `model` — model ID string (e.g., `"claude-sonnet-4-6"`)
- `input_tokens` and `output_tokens` — capture separately when available; fall back to `total_tokens` only if the split is absent
- `duration_ms`

Compute `estimated_cost_usd` using the pricing table in `references/pricing-table.md`. Apply fallback rules from that reference if model or token split is unknown.

For each eval write:
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/outputs/output.json` → `{ "id": N, "answer": "..." }`
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/outputs/output.json` → same
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/timing.json` → see [timing schema](references/grading-schema.md)
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/timing.json` → same

When batching evals into one subagent, the timing values are shared across all evals in the batch. Add `"note": "batch run — timing shared across {N} evals"` to each `timing.json`.

If a subagent cannot write files, capture its answer array from the task result and write all files yourself before continuing.

### PHASE 3 — Grade

For each eval and each configuration, grade every assertion against the answer in `outputs/output.json`:
- **PASS**: assertion clearly satisfied — provide the exact quote or observation as evidence
- **FAIL**: assertion not satisfied — state precisely what was missing or wrong

Grading rules:
- Require concrete evidence for PASS — no benefit of the doubt
- Negative assertions ("does NOT suggest X"): PASS only if X is genuinely absent from the output
- Redirect assertions ("redirects to skill Y"): PASS only if the skill name is explicitly present
- Grade on substance only — not tone, style, or length

Write per-eval per-config grading files:
- `{workspace}/iteration-{N}/eval-{slug}/with_skill/grading.json`
- `{workspace}/iteration-{N}/eval-{slug}/without_skill/grading.json`

Schema → see [reference](references/grading-schema.md).

### PHASE 4 — Benchmark

Aggregate timing and pass rates across all evals for both configurations. Compute delta.

**Costs**: sum `estimated_cost_usd` from each eval's `timing.json` per configuration; compute `mean_per_eval` and `delta.estimated_cost_usd`.

**value_tier**: classify `delta.pass_rate` using the thresholds in [benchmark-schema.md](references/benchmark-schema.md):
- `"forte"` ≥ 0.40 · `"moderado"` 0.20–0.39 · `"fraco"` 0.05–0.19 · `"sem_valor"` 0.00–0.04 · `"negativo"` < 0.00

**delta_vs_prev_iteration**: when `N > 1`, read `iteration-(N-1)/benchmark.json` and compute `pass_rate_delta_change`, `value_tier_change`, and `cost_delta_change_usd`. Set to `null` for iteration 1.

Write `{workspace}/iteration-{N}/benchmark.json` — full schema → see [reference](references/benchmark-schema.md).

Pattern analysis to include in the benchmark:
- Assertions that **always pass in both configs** — these inflate the with_skill rate without measuring skill value; flag for review
- Assertions that **always fail in both configs** — broken assertion or task too hard; flag for fixing
- Assertions that **pass with skill, fail without** — where the skill clearly adds value
- High `stddev` across evals — signal of ambiguous or flaky skill instructions

### PHASE 5 — Feedback File

Write `{workspace}/iteration-{N}/feedback.json` pre-populated with empty strings per eval slug:

```json
{
  "eval-module-structure-verify-layout": "",
  "eval-naming-convention-command-dir": ""
}
```

Instruct the user to fill in specific, actionable notes for each eval where the output missed the point — even if assertions passed. Empty string = output was acceptable. Schema → see [reference](references/feedback-schema.md).

### PHASE 6 — Report

Write a `report.md` file at `{workspace}/iteration-{N}/report.md` using the Write tool. Do NOT print the report to the prompt. The file must be written in Brazilian Portuguese (pt-BR).

Structure of `report.md`:

```markdown
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

_(Se aplicável)_ Nota de precificação: `{pricing_note}`

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
...

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
...

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `slug`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
...
_(Nenhuma — pontuação perfeita)_ se aplicável

## Recomendação

Uma frase — pronta / precisa de iteração / precisa de revisão maior.
```

After writing `report.md`, output **only** this message (nothing else):

```
o relatório da avaliação da skill foi criado, local: {workspace}/iteration-{N}/report.md
```

## Gotchas

- **Never use scripts** (Python, bash, shell heredocs) to write files — use the Write tool exclusively; scripts add latency and increase total run time significantly
- All written content must be in **pt-BR**; only JSON keys, file paths, code identifiers, PASS/FAIL labels stay in their original form
- The report goes to `report.md` inside the iteration dir — never printed to the conversation prompt
- Workspace is **alongside** the skill dir — `my-skill-workspace/` next to `my-skill/`, not inside it
- Always read ALL `references/` files before building the with_skill prompt — critical conventions often live in references/, not SKILL.md body
- Eval slug must be filesystem-safe: lowercase, hyphens only, max 60 chars — truncate if needed
- Iteration number comes from filesystem scan — never assume 1 if `{workspace}/` already exists
- Capture `model`, `input_tokens`, `output_tokens`, and `duration_ms` immediately from task completion notification — they are not available later
- Token split (`input_tokens` / `output_tokens`) may not always be present in the notification; apply the 75/25 fallback from pricing-table.md and record `pricing_note` — never skip cost computation
- `value_tier` is derived solely from `delta.pass_rate` — do not adjust it based on cost or context
- `delta_vs_prev_iteration` requires reading the previous iteration's `benchmark.json` — skip only if iteration 1
- If subagent is denied write permissions, answers still arrive in task result — write all files yourself
- Negative assertions are highest-signal — grade strictly, they catch regressions
- `stddev` in benchmark is only meaningful with multiple runs per eval; with single runs, focus on raw delta
- Baseline-confirmed evals are not waste — they validate no regression in universal knowledge; do not suggest removing them without user input

## References

- Grading output schema → see [reference](references/grading-schema.md)
- Benchmark output schema → see [reference](references/benchmark-schema.md)
- Feedback file schema → see [reference](references/feedback-schema.md)
- Model pricing table → see [reference](references/pricing-table.md)
