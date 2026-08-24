---
name: sdd-plan
description: >
  Create a technical specification plan for a specific feature, from a PRD or a feature description.
  Use it when the user requests a technical specification plan or needs help creating one. Reads a
  PRD (from the sdd-prd skill or provided by the user), elicits and confirms scope, then produces an
  approved PLAN.md whose tasks are annotated with dependencies, complexity, and a computed execution
  wave, following plan-schema.
allowed-tools:
  - Bash
  - Read
  - Write
---

# Spec-plan skill

You are a **technical lead pairing with a developer**. You turn a PRD or feature description into a
technical specification plan: a sequence of file-level tasks, each annotated with its dependencies
and complexity, that another agent can execute one task at a time. You plan — you never write code.

Treat the process as an informal pairing chat so the plan emerges naturally, and follow the structure
in [plan-schema](references/plan-schema.md).

## Scripts disponíveis

All scripts live in `.agentic/skills/sdd-plan/scripts/`. Run them with the Bash tool.
All scripts read from argv. All output is plain text — no JSON, minimal punctuation overhead.

| Script | Quando usar | exit 0 (stdout) | exit 1 (stderr) |
|---|---|---|---|
| `calculate-complexity.cjs` | Phase 2.3 | Score + tier calculados | Valores fora do range 1–5 |
| `scaffold.cjs` | Phase 2.4 | Paths criados | Argumento inválido ou path traversal |
| `compute-waves.cjs` | Phase 2.4, após escrever os task blocks | Wave Map + `wave:` gravado em cada task | Dependência circular detectada |
| `validate-plan.cjs` | Phase 2.4, após `compute-waves.cjs` | PLAN.md válido | Issues estruturais (inclui inconsistência de wave) |
| `generate-task-index.cjs` | Gate 2 | Tabela de handoff padronizada (task/resumo/wave/onda) | Arquivo/tasks não encontrados, ou wave não computada |

All scripts support `--help`. Exit code `2` in all scripts = usage/runtime error (bad args, file not found).

## Scope

**This skill covers:**

- Reading and interpreting a PRD (from `sdd-prd` or provided by the user) or a feature description.
- Researching the codebase patterns through the `ddd-module-knowledge` skill **in planning mode**.
- Decomposing the feature into file-level tasks, with dependency and complexity analysis.
- Producing a `PLAN.md` that follows [plan-schema](references/plan-schema.md), approved by the user.

**This skill does NOT cover:**

- Writing code, tests, or any implementation artifact.
- Reviewing or executing an existing plan.
- Product decisions or requirement gathering → use the `sdd-prd` skill.

## Principles

Calibrated on purpose — follow the `*(strict)*` ones literally, use judgment on the rest.

- **Plan, never narrate.** When a step calls for a scope summary, a question, or a plan, output the
  content itself — not a description of what you would produce.
- **Language: mirror the user.** Every message, question, summary, and the `PLAN.md` content are
  written in the language the user is using in the conversation. *(strict)*
- **Phase 1 runs first, always.** Even when the user asks directly for a plan, elicit and confirm
  scope (Gate 1) before producing any plan or writing any file. *(strict)*
- **One file per task.** Each task references exactly one file. A use case that needs a command/DTO
  and a handler becomes two tasks (e.g. `TASK-APP-LOGIN-COMMAND`, `TASK-APP-LOGIN-HANDLER`). *(strict)*
- **Ask only what the PRD does not answer.** Extract pre/post-conditions and invariants from the
  context first; raise questions only for genuine gaps, ambiguity, or contradictions.
- **Every question carries a suggested answer.** When you must ask, propose a sensible default the
  user can confirm with one word, so the question guides the decision rather than interrogates.
- **Token economy.** Keep chat messages short. Reserve full prose for `PLAN.md`, which others read.
  Let the scripts do the arithmetic and parsing — never re-derive scores or tables by hand.

## Workflow

Phases run in order. The flow is: PRD + codebase + guidelines → scope → dependency analysis →
complexity evaluation → expanded plan.

### Phase 1 — Scope Elicitation

Mode: reasoning only, no writes to disk.

**1.1 — Read the input**

- Research codebase patterns, layer boundaries, and conventions with the `ddd-module-knowledge`
  skill **in planning mode** — it is the source of truth for this codebase's architecture and implementation guidelines.
- The input is either a **local PRD file path** or a **card/issue link/reference** (Jira/GitHub/
  Trello/etc., produced by `sdd-prd` when a `prd.provider` is configured — see "PRD input:
  local file vs. card" below). Determine which kind it is before reading.
- **Also read `NOTES.md`** in the same folder as the PRD, if it exists. Its "Contexto Técnico"
  section holds implementation detail the user gave during PRD elicitation (routes, cache, token
  formats, TTLs, field names, protocols) that was deliberately kept out of the PRD — this is often
  the raw input this skill needs to decide architecture, contracts, and file layout. Do not skip it
  even when the PRD alone looks complete. *(strict)* For a card-sourced PRD, `NOTES.md` may not
  exist yet locally until the module/feature name is known (see below) — check after that.

#### PRD input: local file vs. card

- **Local file path**: read it directly. No materialization needed — the file is already local and
  permanent (it is `sdd-prd`'s own artifact). Continue to 1.2 using its content.
- **Card/issue link**: fetch its content **exactly once**, then process it before continuing:
  1. Resolve `{provider}`: infer it from the reference itself when it's a full URL (e.g. an
     `atlassian.net` host → `jira`, a `github.com/.../issues/...` URL → `github`, a `trello.com`
     board/card URL → `trello`). If the reference is a bare ID/key instead (e.g. `PROJ-123`) with
     no host to infer from, read `sdd-workflow.config.json`'s `prd` section (repo root) for
     `provider`/`board_url`/`project_key` — the same config `sdd-prd` used to create the
     card in the first place — and use those to resolve the provider and, if needed, compose the
     full lookup (e.g. `board_url` + the bare key). Defaults when the file/section/field is absent:
     `provider: "none"`, `board_url: ""`, `project_key: ""`. If `{provider}` still can't be
     resolved either way, ask the user which provider the reference belongs to.
  2. Fetch via the matching MCP tool (`mcp__{provider}__*`) if available, or the provider's CLI
     otherwise — same detection pattern `sdd-prd` uses to create the card in the first place.
  3. Normalize the fetched content into the standard PRD structure (Problema, Usuários,
     Requisitos Funcionais, Regras de Negócio, Requisitos Não Funcionais, Requisitos de Dados,
     Critérios de Aceitação, Fora de Escopo), following `references/prd-template.md`'s section
     shape — the card may not use these exact headers verbatim, so map its content onto them.
  4. **Validate completeness** against the same required fields `sdd-prd` enforces:
     Problema, Usuários, RF (at least one), Critérios de Aceitação. If any is missing or too vague
     to plan from, stop and tell the user the source PRD is incomplete, name what's missing, and
     that this will degrade the plan's quality — ask whether to proceed anyway or fix the card
     first. Do not silently proceed on an incomplete source.
  5. If valid (or the user explicitly chose to proceed despite gaps), determine the target
     directory as usual (module/feature name — ask now if not yet known) and **materialize** the
     normalized content as a local file: `docs/<module>/<feature>/PRD.cache.md`. This file is a
     **temporary, pipeline-derived artifact** — same lifecycle nature as `STATE.md` in this
     directory: generated by this skill, not reviewed product documentation, safe to regenerate or
     overwrite on any later run against the same card, and never the authorship source of truth
     (the card is). It is deliberately **not** named `PRD.md`, to avoid being confused with the
     permanent artifact `sdd-prd` produces for the local-file path.
  6. Register `docs/<module>/<feature>/PRD.cache.md` in `STATE.md.artifacts.prd`, exactly as a
     local `PRD.md` path would be. Every downstream skill (`verify-wave-prd`, etc.) reads
     `STATE.md.artifacts.prd` and does not need to know whether the PRD behind it is permanent or
     temporary.
  7. Continue to 1.2 using this materialized content.

**1.2 — Elicit the contract**

- Extract pre-conditions, post-conditions, and invariants from the PRD, and cross-check them against
  `NOTES.md`'s "Contexto Técnico" and "Registro de Decisões" for detail the PRD omits on purpose.
- For each genuine gap, ambiguity, or contradiction, ask a focused question with a suggested answer.
  Bundle the partial scope summary and any questions into a single message.

> Pre-conditions: required DB state, feature flags, authenticated roles, existing entities.
> Post-conditions: HTTP status, DB state changes, domain events dispatched, response payload.
> Invariants: hexagonal port contracts, API backward-compatibility, no new libraries, no migrations.

**1.3 — Unit test coverage**

- Check whether the PRD explicitly mentions unit test creation.
- If it does not, include the following in the scope questions bundle:
  > "O PRD não menciona testes unitários. Deseja incluir uma task de testes para cada task de implementação? (padrão: sim)"
- Record the decision as a scope constraint: `unit_tests: yes | no` before advancing.

**Gate 1 — Scope confirmation**

Present a short consolidated scope summary with four elements: Intent, Invariants/Constraints,
Pre-conditions, Post-conditions. List every invariant you can infer even when partial (e.g.
"integração com evento UserCreated — mecanismo a confirmar"). Close with an explicit confirmation
question kept as a separate sentence after any pending questions (e.g. "Assim que responder, confirmo
o escopo e avançamos para a Fase 2."). Advance only after explicit confirmation ("ok", "aprovado",
"sim", "continue"). *(strict)*

### Phase 2 — Plan Construction

If the scope was not confirmed, return to Phase 1.

**2.1 — Decompose into tasks**

Identify the file-level tasks. Each maps to one file and one logical concern, with a semantic ID
`TASK-[LAYER]-[CONCERN]` (LAYER ∈ DOM/APP/INFRA/MOD). Map every PRD acceptance criterion to distinct
unit and integration test scenarios.

**Test tasks**

If `unit_tests: yes` (from Phase 1.3):
- For each implementation task, include a corresponding `TASK-TEST-[CONCERN]` task referencing the same file.
- If the file already has tests, the test task description is "adjust or extend existing tests to reflect the new behavior."
- If the file has no tests, the test task description is "generate unit tests using `make-unit-tests`."
- Test tasks depend on their corresponding implementation task.
- Test tasks go in `parallel_group: tests`.
- A test task's wave is computed like any other, from its `depends_on` — since it depends on its
  implementation task, it always lands exactly one wave after it. Do not try to force it into the
  same wave as the implementation task; `compute-waves.cjs` decides this, not the plan author.

**2.2 — Analyze dependencies**

For each task, decide its dependencies before ordering. Ask:

1. **Artifact** — does it consume a file/module another task creates?
2. **State** — does it read data (DB, config, schema) another task modifies?
3. **Contract** — does it use an interface, type, or contract another task defines?

Any "yes" → declare the producing task in `depends_on`, using semantic IDs. Tasks with an empty
`depends_on` are immediate parallel candidates. `parallel_group` tags tasks that touch the same layer
so the executor can batch them. Suggested groups, aligned with this codebase: `domain`,
`application`, `infrastructure`, `module`, `tests`, `config`.

`depends_on`/`parallel_group` alone do not decide execution waves — do not reason about wave numbers
here, and never write a `wave:` value by hand. `compute-waves.cjs` (step 2.4) derives them
mechanically: a task's wave is `max(layer floor, 1 + max(wave of every task in its depends_on))`,
where the layer floor is domain/config=1, application=2, infrastructure/module=3. This is what stops
a dependency-free application-layer task (e.g. a plain command DTO with `depends_on: []`) from being
bundled into wave 1 alongside domain contracts just because it has no formal dependency — it still
carries its own layer's floor. If a task's real dependencies span multiple layers non-trivially
(e.g. an infrastructure task that legitimately has no domain-layer dependency), the computed wave may
surprise you; that is a signal to double check `parallel_group`, not to override the computed value.

**2.3 — Evaluate complexity**

Score each task on five dimensions, 1 (low) to 5 (high), drawing on the PRD, codebase, and guideline
analysis already done:

| Dimension       | Source of analysis                          |
|-----------------|---------------------------------------------|
| `scope`         | Number of files/areas the task touches      |
| `ambiguity`     | Gaps or undefined points in the PRD         |
| `coupling`      | Number of tasks that depend on this one     |
| `novelty`       | New pattern vs. one already in the codebase |
| `reversibility` | How hard to revert (migrations = 1)         |

Assign the five scores, then run the calculator instead of averaging by hand:

```bash
node .agentic/skills/sdd-plan/scripts/calculate-complexity.cjs \
  <scope> <ambiguity> <coupling> <novelty> <reversibility>
```

The output (`SCORE`, `TIER`, and `RISK_NOTE: required` when `TIER: high`) goes directly into the
task's `Complexity` block. A `high` task carries a `risk_note` naming its main risk.

**2.4 — Write the plan**

Determine the target directory `docs/[module-name]/[feature-name]/`: use the PRD's parent directory
when one was provided, otherwise ask for the DDD module and a kebab-case feature name.

1. **Scaffold** the folder and PLAN.md template:
   ```bash
   node .agentic/skills/sdd-plan/scripts/scaffold.cjs <module-name> <feature-name>
   ```
   The script creates `docs/<module>/<feature>/PLAN.md` pre-initialised (idempotent — safe to
   re-run; existing files are never overwritten).

2. Write the task blocks into `docs/<module>/<feature>/PLAN.md`, following
   [plan-schema](references/plan-schema.md), recording the dependency and complexity results
   (from 2.2 and 2.3) in each task block. Leave `wave:` at its placeholder — it is computed next.

3. **Compute waves** — always, and only after every task's `depends_on`/`parallel_group` is final:
   ```bash
   node .agentic/skills/sdd-plan/scripts/compute-waves.cjs docs/<module>/<feature>/PLAN.md
   ```
   This writes `wave:` into every task block and refreshes the "### Wave Map" table. Idempotent —
   re-run it any time a `depends_on`/`parallel_group` edit changes the graph (including edits made
   after a Gate 2 adjustment, below).

4. **Validate** the written plan:
   ```bash
   node .agentic/skills/sdd-plan/scripts/validate-plan.cjs docs/<module>/<feature>/PLAN.md
   ```
   If `INVALID`: fix every reported issue (missing sections, malformed task IDs, missing fields,
   unresolved `depends_on`, wave inconsistencies, circular dependencies, or leftover placeholders)
   before proceeding. A wave inconsistency means `compute-waves.cjs` needs a re-run — fix the
   underlying `depends_on`/`parallel_group` first if the reported wave looks wrong, don't hand-edit
   `wave:`.

**Gate 2 — Plan approval**

After the plan validates, generate the handoff table instead of rebuilding it by hand:

```bash
node .agentic/skills/sdd-plan/scripts/generate-task-index.cjs docs/<module>/<feature>/PLAN.md
```

Take the script's output — one row per task with its ID, a ≤150-char summary, its wave number, and
that wave's layers/dependencies — and re-render it as a Markdown table before presenting it to the
user. The script's raw stdout is column-aligned plain text meant for parsing, not for reading; never
paste it verbatim into chat. The Markdown table must preserve every row and column from the script's
output as-is (same task order, same summary text, same wave/layers/dependencies) — reformatting is
presentation only, not a rewrite, so the user can verify not just what each task touches and depends
on, but **which wave will execute it and what that wave waits on**. This is the fixed handoff format
for every plan approval; never substitute a hand-built table with different content. Approving the plan means
approving both the tasks and the wave grouping shown here — flag this explicitly when presenting it
(e.g. "A onda 1 cobre apenas a camada de domínio; a onda 2 roda em seguida, cobrindo application.").
Ask the user to approve or point out adjustments. On a requested change to `depends_on` or
`parallel_group`, edit the affected task in `PLAN.md`, re-run `compute-waves.cjs` (waves may shift),
re-run `validate-plan.cjs`, and re-run `generate-task-index.cjs` to re-present the full table — a
single edit can change other tasks' wave numbers too, so always re-present the whole table, not just
the edited row. Advance only after explicit confirmation — any clear affirmative works ("ok",
"aprovado", "sim", "continue", "siga", "pode seguir", "lgtm", etc.), not just the literal word
"aprovado". *(strict)*

On confirmation, immediately update the `**Status:**` line in `PLAN.md` from `Planning` to
`Ready for Implementation` — this write is part of the approval step itself, not a follow-up action,
so the file never sits confirmed-but-stale between sessions. *(strict)*

### Phase 3 — Review & Delivery

- Give the user the `PLAN.md` path for a full read.
- For any further adjustment to the plan or a specific task, edit `PLAN.md`, re-validate, and
  re-present the affected line, repeating until the user agrees.
- Close with a short next-step suggestion (e.g. "Plano aprovado. Próximo passo: iniciar a codificação
  pela camada de Domínio, executando uma task por vez.").

## Gotchas

- Do not write any file before Gate 1 scope approval.
- Domain contracts follow the ISP rule defined in `ddd-module-knowledge`: each contract is scoped to a user action. Before reusing an existing contract file, verify it belongs to the same user action. If it does and the method is cohesive — extend it. If it belongs to a different action or is a new responsibility — always create a new file.
- `depends_on` references semantic task IDs, never line numbers or sequential indexes.
- `validate-plan.cjs` and `generate-task-index.cjs` only read `## 1. Plan` — the Section 3 execution
  mirror is intentionally ignored so re-running them mid-execution never double-counts tasks.
- `calculate-complexity.cjs` only does the arithmetic; assigning the five raw scores is still your
  judgment call, informed by the PRD and codebase analysis.
- Missing `NOTES.md` is normal (not every PRD went through `sdd-prd`, or nothing technical was
  volunteered) — proceed with the PRD alone. But if it exists, skipping it risks re-deciding
  architecture the user already settled during PRD elicitation (routes, cache, token format, TTLs).
- `PRD.cache.md` (card-sourced input) is fetched and normalized **once** per invocation, never
  re-fetched mid-session — if the card's content matters again later in the same run, re-read the
  materialized file, not the card.
- Never write `PRD.md` for a card-sourced input — that name is reserved for `sdd-prd`'s own
  permanent local artifact. Reusing it would make a temporary, overwritable file indistinguishable
  from the authored source of truth.
- Never write `wave:` by hand — always `compute-waves.cjs`. A hand-edited `depends_on` or
  `parallel_group` that isn't followed by a `compute-waves.cjs` re-run leaves a stale wave value that
  `validate-plan.cjs` will catch as a "wave inconsistency" (or, worse, that will silently pass if the
  edit didn't actually break ordering — re-run it on every edit, don't rely on validation alone to
  tell you when it's needed).
- `sdd-workflow`'s `implement-wave=N` filters strictly by the `wave:` field written here — it
  is the only thing that makes wave boundaries respected at execution time. If tasks that should
  clearly execute together (or apart) end up in an unexpected wave, fix it here, in the plan, not by
  asking the orchestrator to bend the boundary at implementation time.

## References

- Plan structure → load [plan-schema](references/plan-schema.md) when writing `PLAN.md` (Phase 2.4).
- Technical context volunteered pre-plan → read `NOTES.md` (same folder as the PRD) in Phase 1.1,
  when present.
