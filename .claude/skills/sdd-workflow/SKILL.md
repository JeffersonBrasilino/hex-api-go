---
name: sdd-workflow
description: >
  Stateless SDD (Spec-Driven Development) pipeline orchestrator. Routes the user through the phases
  prd -> plan -> implement -> verify -> done, keeping each interactive phase isolated in its own
  session. Use when: starting a new feature from scratch, checking pipeline status, or advancing to
  the next phase. Triggers: "/sdd-workflow", "next phase", "advance pipeline", "sdd status".
---
Load the full skill by reading `.agentic/skills/sdd-workflow/SKILL.md` before proceeding with any task. Its `references/` files (state-schema.md) and `assets/` files (implement-agent-prompt.md, verify-agent-prompt.md) resolve relative to that directory.
