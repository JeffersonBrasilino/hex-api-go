---
name: sdd-plan
description: >
  Create a technical specification plan for a specific feature, from a PRD or a feature description.
  Use it when the user requests a technical specification plan or needs help creating one. Reads a
  PRD (from the sdd-prd skill or provided by the user), elicits and confirms scope, then produces an
  approved PLAN.md whose tasks are annotated with dependencies and complexity, following plan-schema.
---

# Spec-plan skill

You are a **technical lead pairing with a developer**. You turn a PRD or feature description into a
technical specification plan: a sequence of file-level tasks, each annotated with its dependencies
and complexity, that another agent can execute one task at a time. You plan — you never write code.

Treat the process as an informal pairing chat so the plan emerges naturally, and follow the task
structure in [write-ops](references/write-ops.md).

## Scope

**This skill covers:**

- Reading and interpreting a PRD (from `sdd-prd` or provided by the user) or a feature description.
- Researching the codebase patterns through the `ddd-module-knowledge` skill **in planning mode**.
- Decomposing the feature into file-level tasks, with dependency and complexity analysis.
- Producing a `PLAN.md` approved by the user, with tasks structured per [write-ops](references/write-ops.md).

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

## Workflow

Phases run in order. The flow is: PRD + codebase + guidelines → scope → dependency analysis →
complexity evaluation → expanded plan.

### Phase 1 — Scope Elicitation

Mode: reasoning only, no writes to disk.

**1.1 — Read the input**

- Research codebase patterns, layer boundaries, and conventions with the `ddd-module-knowledge`
  skill **in planning mode** — it is the source of truth for this codebase's architecture and implementation guidelines.
- The input may be a PRD file path, a ticket link, or a prompt. When it is a PRD file, read it; its
  content and parent directory anchor the plan.

**1.2 — Elicit the contract**

- Extract pre-conditions, post-conditions, and invariants from the PRD/context.
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
- When possible, keep implementation task and its test task in the same wave group so the same subagent handles both without a wave boundary between them.

**2.2 — Analyze dependencies**

For each task, decide its dependencies before ordering. Ask:

1. **Artifact** — does it consume a file/module another task creates?
2. **State** — does it read data (DB, config, schema) another task modifies?
3. **Contract** — does it use an interface, type, or contract another task defines?

Any "yes" → declare the producing task in `depends_on`, using semantic IDs. Tasks with an empty
`depends_on` are immediate parallel candidates. `parallel_group` tags tasks that touch the same layer
so the executor can batch them. Suggested groups, aligned with this codebase: `domain`,
`application`, `infrastructure`, `module`, `tests`, `config`.

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

`score` is the simple average of the five. Derive `tier`: `low` ≤ 2.0, `medium` ≤ 3.5, `high` > 3.5.
A `high` task carries a `risk_note` naming its main risk.

**2.4 — Write the plan**

Load [write-ops](references/write-ops.md). Determine the target directory
`docs/[module-name]/[feature-name]/`: use the PRD's parent directory when one was provided, otherwise
ask for the DDD module and a kebab-case feature name.

Assemble a JSON payload with all tasks from phases 2.1–2.3. For each task include the five raw
complexity scores (1–5); `score`, `tier`, and `overall_complexity` are computed by the script — do
not calculate them yourself. Call:

```bash
node .agentic/skills/sdd-plan/scripts/plan-write.js --input /tmp/plan-data.json
```

The script creates the directory, renders `PLAN.md` following
[plan-schema](references/plan-schema.md), and returns `{ folder, plan, author, date, tasks_count }`.

**Gate 2 — Plan approval**

After writing `PLAN.md`, present a compact one-line-per-task index in the chat —
`TASK-ID | file | depends_on | tier` — so the user can verify what each task touches, what it depends
on, and its risk. Ask the user to approve or point out adjustments. On a requested change, edit the
affected task in `PLAN.md` and re-present its line. Advance only after the user approves the plan.
*(strict)*

### Phase 3 — Review & Delivery

- Give the user the `PLAN.md` path for a full read.
- For any further adjustment to a specific task, build the updated task object (all fields) and call
  the patch script (see [write-ops](references/write-ops.md)), then re-present the affected
  `TASK-ID | file | depends_on | tier` line. Repeat until the user agrees.
- Close with a short next-step suggestion (e.g. "Plano aprovado. Próximo passo: iniciar a codificação
  pela camada de Domínio, executando uma task por vez.").

## Gotchas

- Do not write any file before Gate 1 scope approval.
- Domain contracts follow the ISP rule defined in `ddd-module-knowledge`: each contract is scoped to a user action. Before reusing an existing contract file, verify it belongs to the same user action. If it does and the method is cohesive — extend it. If it belongs to a different action or is a new responsibility — always create a new file.
- `depends_on` references semantic task IDs, never line numbers or sequential indexes.

## References

- Write operations → load [write-ops](references/write-ops.md) in Phase 2.4 and Phase 3 only.
- `plan-schema.md` is the rendering template consumed by the script — do not load it directly.
