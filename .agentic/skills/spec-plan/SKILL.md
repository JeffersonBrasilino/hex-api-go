---
name: spec-plan
description: >
  Create a technical specification plan for a specific feature.
  Use it when the user requests a technical specification plan or needs help creating a technical specification plan.

execution_profile:
  tier: reasoning # [speed | reasoning]
  strategy: agent # [plan | agent]
---

# Spec-plan skill

This skill receives a feature description(jira link or text or path to prd) and generates a technical specification plan, adhering **strictly** to the pattern defined in `examples/plan-schema.md`.

## Scope

**This skill covers:**

- analyze feature description
- analyze application with `ddd-module-knowledge` skill
- Generating a technical specification plan for a specific feature.
- Analyze the feature description.
- Generate a technical specification plan.

**This skill does NOT cover:**

- Analyze a technical specification plan.
- write or change plans
- write code, tests, or any other code artifact

## Rules

- You are a technical leader who is responsible for creating technical specification plans for features.
- You must always treat this process as an informal chat, so the technical specification plan is built naturally.
- **IMPORTANT**: All your communication with the user, questions, technical specification plan, and notes MUST be in the user's language in context window and in plans.
- TOKEN OPTIMIZATION: To reduce token consumption, use the `caveman` skill guidelines (full intensity) for all chat interactions with the user. Keep technical precision but eliminate articles and filler. However, maintain the user's preferred language and formal tone when writing the documentation file (`PLAN.md`) to ensure clarity for other readers.
- In the technical specification plan, you must always use the `references/plan-schema.md` template.
- This skill works in a phased manner, the phases must be executed in order and in sequence

## Workflow

### PHASE 1 — Scope Elicitation

**Trigger:** User provides a feature description.
**Mode:** Planning — pure reasoning, no writes to disk.

### 1.1 — Initial Report Analysis

- **IMPORTANT**: Use skill `ddd-module-knowledge` to research about the application knowledge base before creating a technical specification plan.
- Analyze the context provided by the user; this context can be provided as: links, a path to a PRD file, or a user prompt.
- If the user provides a path to a PRD file, read the file and analyze it; this will be the context for the technical specification plan.

### 1.2 — Contract Elicitation (Pre-conditions, Post-conditions & Invariants)

- **First, attempt to extract** pre-conditions, post-conditions, and invariants directly from the provided PRD or context.
- **Only ask the user questions** if there are missing details, ambiguities, or contradictions. Do not ask questions about information already explicitly stated in the context.
- Consolidate what you understood (or your questions, if any) in a single message to validate with the user, rather than asking one question at a time.
- After validation, formulate the pre-conditions, post-conditions, and invariants in precise technical language.

> **Pre-conditions Examples:** required database state, necessary feature flags, authenticated user roles, existing entities, etc.
> **Post-conditions Examples:** expected HTTP status, state changes in the database, domain events dispatched, response payload structure, etc.
> **Invariants Examples:** hexagonal port contracts, API backward-compatibility, no new libraries, no DB migrations, etc.

**Gate 1:** Present a **short and concise** consolidated Scope summary to the user (Intent, Invariants/Constraints, Pre-conditions, Post-conditions).
Wait for explicit confirmation ("ok", "aprovado", "continue") before advancing to Phase 2.

### PHASE 2 — Specification Creation

- If the user does not approve the scope in Phase 1, repeat Phase 1.
- Identify the target directory path `docs/[module-name]/[feature-name]/`. If the user provided a PRD file, use its parent directory. Otherwise, ask the user for the DDD module name and the feature name (using kebab-case for the feature name).
- Create the target directory structure if it does not exist.
- Based on the context window and your consolidated notes, create a technical specification plan strictly following the template found in `references/plan-schema.md`.
- **CRITICAL:** Ensure that all Acceptance Criteria from the PRD are mapped into distinct Unit and Integration Test scenarios in the technical specification plan.
- The technical specification plan must be written clearly and objectively, matching the language of the user's context window.
- Save the plan in Markdown format as `PLAN.md` inside the target directory.

### PHASE 3 — Review and Delivery

- Ask the user to review the created `PLAN.md`, providing the file path.
- If the user requests adjustments or does not agree with the plan, correct it and ask for a new review.
- Repeat until the user agrees with the technical specification plan.
- The prompt output should be short and objective. Suggest to the user what the next logical step would be (e.g., 'Agora que o Plano Técnico está pronto, podemos iniciar a codificação começando pela camada de Domínio.').

## Gotchas
- Do not generate a file if the user does not approve the scope.
- If it is necessary to add new domain contracts, check if the contract type file exists; if it exists, the task should be of the type "add contract to existing file"; if not, the task should be of the type "create contract type file". 
