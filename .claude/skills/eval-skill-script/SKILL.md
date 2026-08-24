---
name: eval-skill-script
description: >
  Run structured evaluations (evals) for any skill that has an evals/evals.json file.
  Executes each test case twice — once with the skill loaded (with_skill) and once without
  it (without_skill baseline) — then grades outputs against assertions via a dedicated haiku
  grader subagent, computes benchmark/cost/timing via scripts, and generates a lightweight
  report.md via a dedicated subagent. Use when the user says "eval this skill", "run evals
  for", "test this skill", or provides a skill path and asks to evaluate it. Requires the
  target skill to have evals/evals.json. Script-enhanced variant: uses Node.js scripts for
  workspace scaffolding, cost/timing computation, grading summary finalization, and benchmark
  aggregation — subagents are reserved for semantic grading and report prose only.
---
Load the full skill by reading `.agentic/skills/eval-skill-script/SKILL.md` before proceeding
with any task. Its `references/` files and `scripts/` resolve relative to that directory.
