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
  "total_tokens": 84852,
  "duration_ms": 23332
}
```

When evals are batched into a single subagent run, add a note:

```json
{
  "total_tokens": 84852,
  "duration_ms": 23332,
  "note": "batch run — timing shared across 15 evals"
}
```

Capture these values immediately from the task completion notification — they are not persisted anywhere else.
