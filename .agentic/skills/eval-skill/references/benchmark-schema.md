# benchmark.json Schema

One `benchmark.json` per iteration, at the iteration root alongside all eval directories.

```json
{
  "run_summary": {
    "model": "claude-sonnet-4-6",
    "with_skill": {
      "pass_rate": { "mean": 0.83, "stddev": 0.06 },
      "time_seconds": { "mean": 45.0, "stddev": 12.0 },
      "tokens": { "mean": 3800, "stddev": 400 },
      "estimated_cost_usd": { "total": 0.4336, "mean_per_eval": 0.0434 }
    },
    "without_skill": {
      "pass_rate": { "mean": 0.33, "stddev": 0.10 },
      "time_seconds": { "mean": 32.0, "stddev": 8.0 },
      "tokens": { "mean": 2100, "stddev": 300 },
      "estimated_cost_usd": { "total": 0.2520, "mean_per_eval": 0.0252 }
    },
    "delta": {
      "pass_rate": 0.50,
      "time_seconds": 13.0,
      "tokens": 1700,
      "estimated_cost_usd": 0.1816
    },
    "value_tier": "forte",
    "delta_vs_prev_iteration": null
  }
}
```

## Field rules

- `model`: model ID used for subagent runs — same value as in `timing.json`
- `pass_rate.mean`: average pass_rate across all evals for this configuration
- `pass_rate.stddev`: standard deviation — only meaningful with multiple runs; set to `null` for single runs
- `time_seconds.mean`: average of `duration_ms / 1000` per eval
- `tokens.mean`: average of `total_tokens` per eval
- `estimated_cost_usd.total`: sum of `estimated_cost_usd` across all eval timing files for this configuration
- `estimated_cost_usd.mean_per_eval`: `total / number_of_evals`
- `delta.pass_rate`: `with_skill.pass_rate.mean - without_skill.pass_rate.mean`
- `delta.time_seconds`: positive means skill costs more time
- `delta.tokens`: positive means skill uses more tokens
- `delta.estimated_cost_usd`: `with_skill.total - without_skill.total` — the extra cost of loading the skill

## value_tier

Classify the skill's value based on `delta.pass_rate`:

| tier          | condition                              |
|---------------|----------------------------------------|
| `"forte"`     | delta.pass_rate ≥ 0.40                 |
| `"moderado"`  | 0.20 ≤ delta.pass_rate < 0.40         |
| `"fraco"`     | 0.05 ≤ delta.pass_rate < 0.20         |
| `"sem_valor"` | 0.00 ≤ delta.pass_rate < 0.05         |
| `"negativo"`  | delta.pass_rate < 0.00 (skill harms)  |

## delta_vs_prev_iteration

When `N > 1`, read `iteration-(N-1)/benchmark.json` and compute:

```json
{
  "delta_vs_prev_iteration": {
    "pass_rate_delta_change": 0.12,
    "value_tier_change": "fraco → forte",
    "cost_delta_change_usd": -0.05
  }
}
```

- `pass_rate_delta_change`: current `delta.pass_rate` − previous `delta.pass_rate` — positive means improvement
- `value_tier_change`: `"{prev_tier} → {current_tier}"` string, or `"sem mudança"` if tier did not change
- `cost_delta_change_usd`: current `delta.estimated_cost_usd` − previous — negative means the skill got cheaper

Set to `null` when this is iteration 1 (no prior benchmark to compare against).

## Interpretation

- `value_tier "forte"` + acceptable cost delta → skill is worth it
- `value_tier "fraco"` or `"sem_valor"` + high cost → skill may not justify the overhead; simplify
- `delta_vs_prev_iteration.pass_rate_delta_change > 0` → iteration is improving the skill
- `value_tier_change` showing regression (e.g., `"forte → moderado"`) → last edit degraded the skill

`stddev` for pass_rate is only meaningful with multiple runs per eval; focus on raw `mean` and `delta` in early iterations.
