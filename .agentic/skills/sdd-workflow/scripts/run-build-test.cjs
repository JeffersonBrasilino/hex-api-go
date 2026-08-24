#!/usr/bin/env node
'use strict';

// Runs `go build ./...` and `go test ./...` once and returns a compact, structured
// summary instead of raw command output — so the orchestrator calling this never
// accumulates full build/test logs in its own context (sdd-workflow's own
// rule: never accumulate phase content). Used once per wave group in Step A.5 and
// once per retry in Step C, so every verify-code agent in a group consumes the same
// precomputed result instead of re-running the suite itself.
//
// Exit codes:
//   0 — build and tests both passed (stdout: structured summary)
//   1 — build or tests failed (stdout: structured summary with STATUS: failed and
//       the relevant failure output — this is a normal outcome, not a script error)
//   2 — usage or runtime error (go not on PATH, etc.)

const { spawnSync } = require('child_process');

const HELP = `
Usage: node scripts/run-build-test.cjs [module_dir]

Runs \`go build ./...\` then, if it passes, \`go test ./...\` in [module_dir]
(default: current directory) and prints a compact structured summary — never the
raw logs — so callers don't accumulate build/test output in their own context.

Exit codes:
  0   Build and tests both passed
  1   Build or tests failed (this is a normal outcome, not a script error —
      the summary on stdout still reports STATUS/details)
  2   Usage or runtime error (e.g. \`go\` not found on PATH)

Examples:
  node scripts/run-build-test.cjs
  node scripts/run-build-test.cjs .
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , moduleDirArg] = process.argv;
const cwd = moduleDirArg || '.';

function run(args) {
  const result = spawnSync('go', args, { cwd, encoding: 'utf8' });
  if (result.error) {
    process.stderr.write(`Error: failed to run "go ${args.join(' ')}": ${result.error.message}\n`);
    process.exit(2);
  }
  return result;
}

// --- build ---------------------------------------------------------------

const build = run(['build', './...']);

if (build.status !== 0) {
  const errorLines = (build.stderr || build.stdout || '').trim().split('\n');
  process.stdout.write('STATUS: failed\n');
  process.stdout.write('BUILD: failed\n');
  process.stdout.write('TESTS: skipped (build failed)\n');
  process.stdout.write('---\n');
  process.stdout.write('Build errors:\n');
  errorLines.forEach(l => process.stdout.write(`  ${l}\n`));
  process.exit(1);
}

// --- tests -----------------------------------------------------------------

const test = run(['test', './...', '-cover']);
const testOutputLines = (test.stdout || '').trim().split('\n').filter(Boolean);

const packageLines = testOutputLines.filter(l => /^(ok|FAIL|---)/.test(l));
const failedPackages = packageLines.filter(l => l.startsWith('FAIL'));
const passedPackages = packageLines.filter(l => l.startsWith('ok'));

const testsPass = test.status === 0;

process.stdout.write(`STATUS: ${testsPass ? 'ok' : 'failed'}\n`);
process.stdout.write('BUILD: ok\n');
process.stdout.write(`TESTS: ${testsPass ? 'ok' : 'failed'}\n`);
process.stdout.write(`PACKAGES_PASSED: ${passedPackages.length}\n`);
process.stdout.write(`PACKAGES_FAILED: ${failedPackages.length}\n`);
process.stdout.write('---\n');
for (const l of packageLines) process.stdout.write(`${l}\n`);

if (!testsPass) {
  // Include failure detail (assertion output) beyond the per-package summary lines,
  // capped so a single runaway test doesn't reintroduce the log-dump problem this
  // script exists to avoid.
  const detailLines = testOutputLines.filter(l => !/^(ok|FAIL|---)/.test(l)).slice(0, 40);
  if (detailLines.length > 0) {
    process.stdout.write('---\n');
    process.stdout.write('Failure detail (capped at 40 lines):\n');
    detailLines.forEach(l => process.stdout.write(`  ${l}\n`));
  }
  process.exit(1);
}

process.exit(0);
