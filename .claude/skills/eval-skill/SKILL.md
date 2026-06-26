---
name: eval-skill
description: >
  Run structured evaluations (evals) for any skill that has an evals/evals.json file.
  Executes each test case with and without the skill loaded, grades outputs against assertions,
  and produces grading, benchmark, and report files. Use when asked to eval, test, or validate a skill.
execution_profile:
  tier: speed
  strategy: plan
---

Load the full skill instructions by reading `.agentic/skills/eval-skill/SKILL.md` before proceeding with any task.