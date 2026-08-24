# feedback.json Schema

One `feedback.json` per iteration, at the iteration root alongside `benchmark.json`.

The file maps eval slug → human review note. It is pre-populated by the eval-skill with empty strings, then filled in by the user after reviewing outputs.

```json
{
  "eval-top-months-chart": "The chart is missing axis labels and the months are in alphabetical order instead of chronological.",
  "eval-clean-missing-emails": ""
}
```

## Field rules

- Keys: eval slugs — must match the directory names under `iteration-N/`
- Values: free-text feedback string
  - Empty string `""` = output was acceptable, no complaints
  - Non-empty = specific, actionable complaint about the output

## Feedback quality

Good feedback (actionable):
- `"The chart is missing axis labels and the months are in alphabetical order instead of chronological."`
- `"Answer omits the file name — only struct and constructor were mentioned."`

Weak feedback (not actionable):
- `"looks bad"` — too vague
- `"could be better"` — no signal for improvement

## Usage in iteration

During the iteration step, provide `feedback.json` alongside `grading.json` results and the current `SKILL.md` to an LLM to propose improvements. Evals with non-empty feedback are the highest priority for skill instruction changes.
