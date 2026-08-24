#!/usr/bin/env node
'use strict';

// Averages the 5 complexity dimensions and derives the tier. Pure arithmetic —
// no reasoning required, so it never touches the model's context budget.
//
// Exit codes:
//   0 — score and tier calculated (stdout: SCORE, TIER, DIMENSIONS, RISK_NOTE)
//   1 — a dimension value is out of range (stderr: which one + expected range)
//   2 — usage error (stderr: missing arguments + usage hint)

const HELP = `
Usage: node scripts/calculate-complexity.cjs <scope> <ambiguity> <coupling> <novelty> <reversibility>

Averages the 5 complexity dimensions (each 1-5) into a score and derives the
tier: low <= 2.0, medium <= 3.5, high > 3.5. A high tier requires a risk_note
in the task block.

Arguments (all integers 1-5):
  scope           Number of files/areas the task touches
  ambiguity       Gaps or undefined points in the PRD
  coupling        Number of tasks that depend on this one
  novelty         New pattern vs. one already in the codebase
  reversibility   How hard to revert (migrations = 1)

Exit codes:
  0   Score and tier calculated
  1   A value is out of the 1-5 range
  2   Usage error (missing or non-numeric argument)

Examples:
  node scripts/calculate-complexity.cjs 2 3 3 3 2
  node scripts/calculate-complexity.cjs 4 4 5 3 2
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const DIMENSIONS = ['scope', 'ambiguity', 'coupling', 'novelty', 'reversibility'];
const rawArgs = process.argv.slice(2);

if (rawArgs.length !== 5) {
  process.stderr.write([
    `Error: exactly 5 arguments are required (${DIMENSIONS.join(', ')}).`,
    `Received: ${rawArgs.length} argument(s) — ${JSON.stringify(rawArgs)}`,
    `Expected: 5 integers between 1 and 5, e.g. "2 3 3 3 2"`,
    `Usage: node scripts/calculate-complexity.cjs <scope> <ambiguity> <coupling> <novelty> <reversibility>`,
  ].join('\n') + '\n');
  process.exit(2);
}

const values = {};
for (let i = 0; i < DIMENSIONS.length; i++) {
  const name = DIMENSIONS[i];
  const raw  = rawArgs[i];
  if (!/^-?\d+$/.test(raw)) {
    process.stderr.write([
      `Error: <${name}> must be an integer.`,
      `Received: "${raw}"`,
      `Expected: an integer between 1 and 5`,
    ].join('\n') + '\n');
    process.exit(2);
  }
  values[name] = parseInt(raw, 10);
}

const outOfRange = DIMENSIONS.filter(name => values[name] < 1 || values[name] > 5);
if (outOfRange.length > 0) {
  process.stderr.write([
    `Error: value(s) out of range — each dimension must be between 1 and 5.`,
    ...outOfRange.map(name => `Found:    ${name}=${values[name]}`),
    `Expected: 1 <= value <= 5`,
    `Action:   correct the out-of-range dimension(s) and re-run.`,
  ].join('\n') + '\n');
  process.exit(1);
}

const sum   = DIMENSIONS.reduce((acc, name) => acc + values[name], 0);
const score = Math.round((sum / DIMENSIONS.length) * 10) / 10;

let tier;
if (score <= 2.0) tier = 'low';
else if (score <= 3.5) tier = 'medium';
else tier = 'high';

process.stdout.write(`SCORE: ${score.toFixed(1)}\n`);
process.stdout.write(`TIER: ${tier}\n`);
process.stdout.write(`DIMENSIONS: ${DIMENSIONS.map(name => `${name}=${values[name]}`).join(' ')}\n`);
if (tier === 'high') {
  process.stdout.write(`RISK_NOTE: required\n`);
}
process.exit(0);
