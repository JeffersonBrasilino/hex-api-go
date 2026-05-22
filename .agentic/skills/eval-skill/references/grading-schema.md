# grading.json Schema

One `grading.json` per eval per configuration (`with_skill/` and `without_skill/`).

```json
{
  "assertion_results": [
    {
      "text": "The output includes a bar chart image file",
      "passed": true,
      "evidence": "Found chart.png (45KB) in outputs directory"
    },
    {
      "text": "The chart shows exactly 3 months",
      "passed": true,
      "evidence": "Chart displays bars for March, July, and November"
    },
    {
      "text": "Both axes are labeled",
      "passed": false,
      "evidence": "Y-axis is labeled 'Revenue ($)' but X-axis has no label"
    },
    {
      "text": "The chart title or caption mentions revenue",
      "passed": true,
      "evidence": "Chart title reads 'Top 3 Months by Revenue'"
    }
  ],
  "summary": {
    "passed": 3,
    "failed": 1,
    "total": 4,
    "pass_rate": 0.75
  }
}
```

## Field rules

- `text`: copy verbatim from `evals.json assertions[]` — do not paraphrase
- `passed`: boolean — no partial credit
- `evidence`: must quote or precisely reference the output — never "seems correct" or "looks good"
- `summary.pass_rate`: `passed / total` rounded to 3 decimal places

---

# timing.json Schema

One `timing.json` per eval per configuration. Values come from the task completion notification.

```json
{
  "model": "claude-sonnet-4-6",
  "input_tokens": 63000,
  "output_tokens": 1852,
  "total_tokens": 64852,
  "duration_ms": 23332,
  "estimated_cost_usd": 0.216780
}
```

When evals are batched into a single subagent run, add a `note` field:

```json
{
  "model": "claude-sonnet-4-6",
  "input_tokens": 63000,
  "output_tokens": 1852,
  "total_tokens": 64852,
  "duration_ms": 23332,
  "estimated_cost_usd": 0.216780,
  "note": "batch run — timing shared across 15 evals"
}
```

When the token split is unavailable (only `total_tokens` captured), add `pricing_note`:

```json
{
  "model": "claude-sonnet-4-6",
  "input_tokens": null,
  "output_tokens": null,
  "total_tokens": 84852,
  "duration_ms": 23332,
  "estimated_cost_usd": 0.254556,
  "pricing_note": "token split estimated 75/25 — actual split unavailable"
}
```

## Field rules

- `model`: model ID string exactly as reported by the subagent or set by the runner — never infer
- `input_tokens` / `output_tokens`: from task completion notification; `null` if unavailable
- `total_tokens`: always present — `input_tokens + output_tokens` or directly from notification
- `estimated_cost_usd`: computed from pricing-table.md; see that reference for fallback rules
- `pricing_note`: present only when a fallback was applied; omit the field entirely otherwise

Capture token values and model immediately from the task completion notification — they are not available later.
