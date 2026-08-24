#!/usr/bin/env node
'use strict';

// Computes estimated_cost_usd from model + token counts using the pricing table
// mirrored from references/pricing-table.md, applies the documented fallbacks
// (unknown model -> sonnet pricing; missing token split -> 75/25 estimate), and
// writes timing.json at <out_path>. Replaces manual arithmetic previously done
// by Claude after each runner subagent completes.
//
// Exit codes:
//   0 — timing.json written (stdout: model, total tokens, cost, notes)
//   2 — usage or runtime error (missing required flag)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/write-timing.cjs <out_path> --model <id> --duration-ms <N>
         [--input-tokens <N> --output-tokens <N>] [--total-tokens <N>] [--note <text>]

Computes estimated_cost_usd from the pricing table (mirrors
references/pricing-table.md) and writes timing.json at <out_path>.

If --input-tokens/--output-tokens are given, total_tokens = input + output.
If only --total-tokens is given, applies a 75% input / 25% output fallback and
records a pricing_note. If the model is not in the pricing table, falls back to
Sonnet pricing and records a pricing_note (fallbacks concatenate with "; ").

Exit codes:
  0   timing.json written
  2   Usage or runtime error (missing required flag)

Examples:
  node scripts/write-timing.cjs out/timing.json --model claude-sonnet-4-6 \\
    --input-tokens 63000 --output-tokens 1852 --duration-ms 23332

  node scripts/write-timing.cjs out/timing.json --model claude-sonnet-4-6 \\
    --total-tokens 84852 --duration-ms 23332 --note "batch run — timing shared across 5 evals"
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

function fail2(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(2);
}

// Mirrors references/pricing-table.md — keep both in sync if the table changes.
const PRICING = {
  'claude-haiku-4-5-20251001': { input: 0.80, output: 4.00 },
  'claude-haiku-4-5': { input: 0.80, output: 4.00 },
  'claude-sonnet-4-6': { input: 3.00, output: 15.00 },
  'claude-opus-4-7': { input: 15.00, output: 75.00 },
};
const FALLBACK_MODEL = 'claude-sonnet-4-6';

const argv = process.argv.slice(2);
const outPath = argv[0] && !argv[0].startsWith('--') ? argv[0] : null;

if (!outPath) {
  fail2([
    'Error: <out_path> is required.',
    'Usage: node scripts/write-timing.cjs <out_path> --model <id> --duration-ms <N> [...]',
  ]);
}

const flags = {};
for (let i = 1; i < argv.length; i++) {
  const tok = argv[i];
  if (!tok.startsWith('--')) continue;
  const name = tok.slice(2);
  const value = argv[i + 1];
  if (value === undefined || value.startsWith('--')) {
    fail2([`Error: flag --${name} requires a value.`, 'Usage: node scripts/write-timing.cjs <out_path> --model <id> --duration-ms <N> [...]']);
  }
  flags[name] = value;
  i++;
}

const REQUIRED = ['model', 'duration-ms'];
for (const name of REQUIRED) {
  if (flags[name] === undefined) {
    fail2([
      `Error: missing required flag --${name}.`,
      'Action: pass all required flags — see --help for the full list.',
    ]);
  }
}

if (flags['total-tokens'] === undefined && (flags['input-tokens'] === undefined || flags['output-tokens'] === undefined)) {
  fail2([
    'Error: pass either --total-tokens, or both --input-tokens and --output-tokens.',
    'Action: capture token counts from the task completion notification before calling this script.',
  ]);
}

const durationMs = parseInt(flags['duration-ms'], 10);
if (!Number.isInteger(durationMs) || durationMs < 0) {
  fail2([
    `Error: --duration-ms "${flags['duration-ms']}" is not a valid non-negative integer.`,
    'Action: pass the duration in milliseconds as reported by the task completion notification.',
  ]);
}

const model = flags['model'];
const notes = [];

let inputTokens;
let outputTokens;
let totalTokens;

if (flags['input-tokens'] !== undefined && flags['output-tokens'] !== undefined) {
  inputTokens = parseInt(flags['input-tokens'], 10);
  outputTokens = parseInt(flags['output-tokens'], 10);
  totalTokens = inputTokens + outputTokens;
} else {
  totalTokens = parseInt(flags['total-tokens'], 10);
  inputTokens = Math.round(totalTokens * 0.75);
  outputTokens = totalTokens - inputTokens;
  notes.push('token split estimated 75/25 — actual split unavailable');
}

let pricing = PRICING[model];
if (!pricing) {
  pricing = PRICING[FALLBACK_MODEL];
  notes.push('unknown model — sonnet pricing used as fallback');
}

const estimatedCostUsd = Number(
  ((inputTokens / 1_000_000) * pricing.input + (outputTokens / 1_000_000) * pricing.output).toFixed(6)
);

const timing = {
  model,
  input_tokens: inputTokens,
  output_tokens: outputTokens,
  total_tokens: totalTokens,
  duration_ms: durationMs,
  estimated_cost_usd: estimatedCostUsd,
};

if (flags['note']) timing.note = flags['note'];
if (notes.length > 0) timing.pricing_note = notes.join('; ');

fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, JSON.stringify(timing, null, 2) + '\n');

process.stdout.write(`STATUS: ok\n`);
process.stdout.write(`FILE: ${outPath}\n`);
process.stdout.write(`MODEL: ${model}\n`);
process.stdout.write(`TOTAL_TOKENS: ${totalTokens}\n`);
process.stdout.write(`ESTIMATED_COST_USD: ${estimatedCostUsd.toFixed(6)}\n`);
process.stdout.write(`NOTE: ${notes.length > 0 ? notes.join('; ') : 'none'}\n`);
process.exit(0);
