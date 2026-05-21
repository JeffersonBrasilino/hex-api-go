# benchmark.json Schema

One `benchmark.json` per iteration, at the iteration root alongside all eval directories.

```json
{
  "run_summary": {
    "with_skill": {
      "pass_rate": { "mean": 0.83, "stddev": 0.06 },
      "time_seconds": { "mean": 45.0, "stddev": 12.0 },
      "tokens": { "mean": 3800, "stddev": 400 }
    },
    "without_skill": {
      "pass_rate": { "mean": 0.33, "stddev": 0.10 },
      "time_seconds": { "mean": 32.0, "stddev": 8.0 },
      "tokens": { "mean": 2100, "stddev": 300 }
    },
    "delta": {
      "pass_rate": 0.50,
      "time_seconds": 13.0,
      "tokens": 1700
    }
  }
}
```

## Field rules

- `pass_rate.mean`: average pass_rate across all evals for this configuration
- `pass_rate.stddev`: standard deviation — only meaningful with multiple runs per eval; with single runs, set to `null`
- `time_seconds.mean`: average of `duration_ms / 1000` per eval
- `tokens.mean`: average of `total_tokens` per eval
- `delta.pass_rate`: `with_skill.pass_rate.mean - without_skill.pass_rate.mean`
- `delta.time_seconds`: `with_skill.time_seconds.mean - without_skill.time_seconds.mean` — positive means skill costs more time
- `delta.tokens`: `with_skill.tokens.mean - without_skill.tokens.mean` — positive means skill uses more tokens

## Interpretation

The `delta` shows what the skill **costs** (time, tokens) vs what it **buys** (pass rate):
- High pass_rate delta + reasonable time/token cost → skill is worth it
- Low pass_rate delta + high time/token cost → skill may not be worth it; consider simplifying

`stddev` for pass_rate becomes useful when running each eval multiple times. With single runs per eval (first iterations), focus on the raw `mean` and `delta`.
