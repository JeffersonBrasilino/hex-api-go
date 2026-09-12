#!/usr/bin/env node
'use strict';

// Aggregates grading.json + timing.json across all evals in an iteration into
// benchmark.json (mean/stddev/delta/value_tier/delta_vs_prev_iteration, per
// references/benchmark-schema.md) and seeds feedback.json (per
// references/feedback-schema.md) if it does not already exist. Replaces the
// single "benchmark aggregator" subagent that previously had to receive
// benchmark-schema.md + feedback-schema.md + pricing-table.md + every raw
// grading/timing file just to do fixed arithmetic.
//
// Exit codes:
//   0 — benchmark.json written (stdout: key aggregate figures)
//   1 — a grading.json is missing its summary (finalize-grading.cjs not run yet)
//   2 — usage or runtime error (iteration path not found)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/compute-benchmark.cjs <iteration_path> [--prev <prev_iteration_path>]

Reads every eval-*/{with_skill,without_skill}/{grading.json,timing.json} under
<iteration_path>, computes per-configuration aggregates (pass_rate, time_seconds,
tokens, estimated_cost_usd — mean/stddev/total as applicable), the with/without
delta, and value_tier. Writes benchmark.json (always) and feedback.json (only if
it does not already exist — never overwrites human review).

If --prev is passed, or a sibling "iteration-{N-1}" directory exists, reads its
benchmark.json and computes delta_vs_prev_iteration.

Requires finalize-grading.cjs to have already added "summary" to every
grading.json in this iteration.

Exit codes:
  0   benchmark.json written
  1   Logic violation — a grading.json has no "summary" field yet
  2   Usage or runtime error (iteration path not found)

Examples:
  node scripts/compute-benchmark.cjs .agentic/skills/foo/evals/workspace/iteration-3
  node scripts/compute-benchmark.cjs .../iteration-3 --prev .../iteration-2
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

function fail1(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(1);
}

function fail2(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(2);
}

const argv = process.argv.slice(2);
const iterationArg = argv[0] && !argv[0].startsWith('--') ? argv[0] : null;

if (!iterationArg) {
  fail2([
    'Error: <iteration_path> is required.',
    'Usage: node scripts/compute-benchmark.cjs <iteration_path> [--prev <prev_iteration_path>]',
  ]);
}

const iterationPath = iterationArg.replace(/\/+$/, '');

if (!fs.existsSync(iterationPath) || !fs.statSync(iterationPath).isDirectory()) {
  fail2([
    `Error: iteration path not found: "${iterationPath}"`,
    'Action: confirm the iteration was scaffolded and evals were run before computing benchmark.',
  ]);
}

const prevIdx = argv.indexOf('--prev');
let prevPath = prevIdx !== -1 ? argv[prevIdx + 1] : null;
if (!prevPath) {
  const m = iterationPath.match(/^(.*)\/iteration-(\d+)$/);
  if (m) {
    const n = parseInt(m[2], 10);
    if (n > 1) {
      const candidate = `${m[1]}/iteration-${n - 1}`;
      if (fs.existsSync(path.join(candidate, 'benchmark.json'))) prevPath = candidate;
    }
  }
}

const evalDirs = fs.readdirSync(iterationPath)
  .filter(name => name.startsWith('eval-'))
  .filter(name => fs.statSync(path.join(iterationPath, name)).isDirectory());

if (evalDirs.length === 0) {
  fail2([
    `Error: no eval-* directories found under "${iterationPath}"`,
    'Action: confirm scaffold-eval-workspace.cjs and the runner subagents ran for this iteration.',
  ]);
}

function readJson(p, evalDir, config, kind) {
  if (!fs.existsSync(p)) {
    fail1([
      `Error: ${evalDir}/${config}/${kind}.json not found.`,
      'Action: confirm the runner and grader subagents completed for every eval before computing the benchmark.',
    ]);
  }
  return JSON.parse(fs.readFileSync(p, 'utf8'));
}

const perEval = { with_skill: [], without_skill: [] };

for (const evalDir of evalDirs) {
  for (const config of ['with_skill', 'without_skill']) {
    const grading = readJson(path.join(iterationPath, evalDir, config, 'grading.json'), evalDir, config, 'grading');
    if (!grading.summary || typeof grading.summary.pass_rate !== 'number') {
      fail1([
        `Error: ${evalDir}/${config}/grading.json has no "summary" field.`,
        'Action: run finalize-grading.cjs for this iteration before compute-benchmark.cjs.',
      ]);
    }
    const timing = readJson(path.join(iterationPath, evalDir, config, 'timing.json'), evalDir, config, 'timing');
    perEval[config].push({
      slug: evalDir,
      pass_rate: grading.summary.pass_rate,
      time_seconds: timing.duration_ms / 1000,
      tokens: timing.total_tokens,
      estimated_cost_usd: timing.estimated_cost_usd,
      model: timing.model,
    });
  }
}

function mean(values) {
  return values.reduce((a, b) => a + b, 0) / values.length;
}

function stddev(values) {
  if (values.length < 2) return null;
  const m = mean(values);
  const variance = values.reduce((a, v) => a + (v - m) ** 2, 0) / (values.length - 1);
  return Number(Math.sqrt(variance).toFixed(4));
}

function round(n, places) {
  return Number(n.toFixed(places));
}

function aggregate(rows) {
  const passRates = rows.map(r => r.pass_rate);
  const times = rows.map(r => r.time_seconds);
  const tokens = rows.map(r => r.tokens);
  const costs = rows.map(r => r.estimated_cost_usd);
  const costTotal = round(costs.reduce((a, b) => a + b, 0), 6);
  return {
    pass_rate: { mean: round(mean(passRates), 3), stddev: stddev(passRates) },
    time_seconds: { mean: round(mean(times), 1), stddev: stddev(times) },
    tokens: { mean: round(mean(tokens), 0), stddev: stddev(tokens) },
    estimated_cost_usd: { total: costTotal, mean_per_eval: round(costTotal / rows.length, 6) },
  };
}

const withSkill = aggregate(perEval.with_skill);
const withoutSkill = aggregate(perEval.without_skill);
const model = perEval.with_skill[0].model;

const delta = {
  pass_rate: round(withSkill.pass_rate.mean - withoutSkill.pass_rate.mean, 3),
  time_seconds: round(withSkill.time_seconds.mean - withoutSkill.time_seconds.mean, 1),
  tokens: round(withSkill.tokens.mean - withoutSkill.tokens.mean, 0),
  estimated_cost_usd: round(withSkill.estimated_cost_usd.total - withoutSkill.estimated_cost_usd.total, 6),
};

function classifyTier(deltaPassRate) {
  if (deltaPassRate >= 0.40) return 'forte';
  if (deltaPassRate >= 0.20) return 'moderado';
  if (deltaPassRate >= 0.05) return 'fraco';
  if (deltaPassRate >= 0.00) return 'sem_valor';
  return 'negativo';
}

const valueTier = classifyTier(delta.pass_rate);

let deltaVsPrev = null;
if (prevPath) {
  const prevBenchmark = JSON.parse(fs.readFileSync(path.join(prevPath, 'benchmark.json'), 'utf8'));
  const prevSummary = prevBenchmark.run_summary;
  deltaVsPrev = {
    pass_rate_delta_change: round(delta.pass_rate - prevSummary.delta.pass_rate, 3),
    value_tier_change: prevSummary.value_tier === valueTier
      ? 'sem mudança'
      : `${prevSummary.value_tier} → ${valueTier}`,
    cost_delta_change_usd: round(delta.estimated_cost_usd - prevSummary.delta.estimated_cost_usd, 6),
  };
}

const benchmark = {
  run_summary: {
    model,
    with_skill: withSkill,
    without_skill: withoutSkill,
    delta,
    value_tier: valueTier,
    delta_vs_prev_iteration: deltaVsPrev,
  },
};

fs.writeFileSync(path.join(iterationPath, 'benchmark.json'), JSON.stringify(benchmark, null, 2) + '\n');

const feedbackPath = path.join(iterationPath, 'feedback.json');
let feedbackStatus;
if (fs.existsSync(feedbackPath)) {
  feedbackStatus = 'preserved (already existed)';
} else {
  const feedback = {};
  for (const evalDir of evalDirs) feedback[evalDir] = '';
  fs.writeFileSync(feedbackPath, JSON.stringify(feedback, null, 2) + '\n');
  feedbackStatus = 'created';
}

process.stdout.write(`STATUS: ok\n`);
process.stdout.write(`ITERATION: ${iterationPath}\n`);
process.stdout.write(`EVAL_COUNT: ${evalDirs.length}\n`);
process.stdout.write(`WITH_SKILL_PASS_RATE_MEAN: ${withSkill.pass_rate.mean.toFixed(3)}\n`);
process.stdout.write(`WITHOUT_SKILL_PASS_RATE_MEAN: ${withoutSkill.pass_rate.mean.toFixed(3)}\n`);
process.stdout.write(`DELTA_PASS_RATE: ${delta.pass_rate.toFixed(3)}\n`);
process.stdout.write(`VALUE_TIER: ${valueTier}\n`);
process.stdout.write(`FEEDBACK_FILE: ${feedbackStatus}\n`);
process.exit(0);
