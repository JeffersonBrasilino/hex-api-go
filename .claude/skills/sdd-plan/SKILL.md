---
name: sdd-plan
description: >
  Create a technical specification plan for a specific feature, from a PRD or a feature description.
  Use it when the user requests a technical specification plan or needs help creating one. Reads a
  PRD (from the sdd-prd skill or provided by the user), elicits and confirms scope, then produces an
  approved PLAN.md whose tasks are annotated with dependencies and complexity, following plan-schema.
allowed-tools:
  - Bash
  - Read
  - Write
---
Load the full skill by reading `.agentic/skills/sdd-plan/SKILL.md` before proceeding with any task. Its `references/` files (plan-schema.md) and `scripts/` resolve relative to that directory.
