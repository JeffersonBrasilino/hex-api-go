# Pricing Table — Claude Models

Use this table to compute `estimated_cost_usd` from `input_tokens` and `output_tokens`.
Prices are in USD per million tokens.

| model_id                      | input ($/M) | output ($/M) |
|-------------------------------|-------------|--------------|
| claude-haiku-4-5-20251001     | 0.80        | 4.00         |
| claude-haiku-4-5              | 0.80        | 4.00         |
| claude-sonnet-4-6             | 3.00        | 15.00        |
| claude-opus-4-7               | 15.00       | 75.00        |

## Computation

```
estimated_cost_usd = (input_tokens / 1_000_000 * input_price)
                   + (output_tokens / 1_000_000 * output_price)
```

Round to 6 decimal places.

## Fallback rules

1. **Unknown model**: use Sonnet pricing and add `"pricing_note": "unknown model — sonnet pricing used as fallback"`.
2. **Token split unavailable** (only `total_tokens` captured): assume 75 % input / 25 % output and add `"pricing_note": "token split estimated 75/25 — actual split unavailable"`.
3. Both fallbacks can apply simultaneously; concatenate both notes with "; ".
