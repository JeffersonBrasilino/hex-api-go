#!/usr/bin/env node
'use strict';

// Scans <iteration_path>/eval-*/{with_skill,without_skill}/grading.json, validates
// that each has an assertion_results array, computes summary.{passed,failed,total,
// pass_rate} and writes it back — idempotently, skipping files whose summary is
// already correct. Lets the grader subagent focus purely on semantic PASS/FAIL/
// evidence judgment instead of also being responsible for schema-correct arithmetic.
//
// Exit codes:
//   0 — scan complete (stdout: files scanned/updated/unchanged)
//   1 — a grading.json has an invalid shape (stderr: file + expected shape)
//   2 — usage or runtime error (iteration path not found)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/finalize-grading.cjs <iteration_path>

Scans <iteration_path>/eval-*/{with_skill,without_skill}/grading.json. Each file
must contain an "assertion_results" array of { text, passed, evidence } entries.
Computes summary = { passed, failed, total, pass_rate } (pass_rate rounded to 3
decimal places) and writes it into the file, preserving assertion_results as-is.
Files whose summary already matches the computed value are left untouched.

Exit codes:
  0   Scan complete
  1   Logic violation — a grading.json is missing assertion_results or a required field
  2   Usage or runtime error (iteration path not found)

Examples:
  node scripts/finalize-grading.cjs .agentic/skills/foo/evals/workspace/iteration-3
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

const [, , iterationArg] = process.argv;

if (!iterationArg) {
  fail2([
    'Error: <iteration_path> is required.',
    'Usage: node scripts/finalize-grading.cjs <iteration_path>',
  ]);
}

const iterationPath = iterationArg.replace(/\/+$/, '');

if (!fs.existsSync(iterationPath) || !fs.statSync(iterationPath).isDirectory()) {
  fail2([
    `Error: iteration path not found: "${iterationPath}"`,
    'Action: confirm scaffold-eval-workspace.cjs was run for this iteration.',
  ]);
}

const evalDirs = fs.readdirSync(iterationPath)
  .filter(name => name.startsWith('eval-'))
  .filter(name => fs.statSync(path.join(iterationPath, name)).isDirectory());

const gradingFiles = [];
for (const evalDir of evalDirs) {
  for (const config of ['with_skill', 'without_skill']) {
    const p = path.join(iterationPath, evalDir, config, 'grading.json');
    if (fs.existsSync(p)) gradingFiles.push(p);
  }
}

let updated = 0;
let unchanged = 0;

for (const file of gradingFiles) {
  let data;
  try {
    data = JSON.parse(fs.readFileSync(file, 'utf8'));
  } catch (e) {
    fail1([
      `Error: "${file}" is not valid JSON.`,
      `Details: ${e.message}`,
      'Action: have the grader subagent rewrite this file with valid JSON.',
    ]);
  }

  if (!Array.isArray(data.assertion_results) || data.assertion_results.length === 0) {
    fail1([
      'Error: grading.json missing "assertion_results" array or entry missing required field.',
      `File: ${file}`,
      'Expected: { "assertion_results": [{ "text", "passed", "evidence" }, ...] }',
      'Action: have the grader subagent rewrite this file with the correct shape before retrying.',
    ]);
  }

  const badEntry = data.assertion_results.find(
    a => typeof a.text !== 'string' || typeof a.passed !== 'boolean' || typeof a.evidence !== 'string'
  );
  if (badEntry) {
    fail1([
      'Error: grading.json missing "assertion_results" array or entry missing required field.',
      `File: ${file}`,
      'Expected: { "assertion_results": [{ "text", "passed", "evidence" }, ...] }',
      'Action: have the grader subagent rewrite this file with the correct shape before retrying.',
    ]);
  }

  const passed = data.assertion_results.filter(a => a.passed).length;
  const total = data.assertion_results.length;
  const failed = total - passed;
  const passRate = Number((passed / total).toFixed(3));
  const summary = { passed, failed, total, pass_rate: passRate };

  const current = data.summary;
  const isSame = current
    && current.passed === summary.passed
    && current.failed === summary.failed
    && current.total === summary.total
    && current.pass_rate === summary.pass_rate;

  if (isSame) {
    unchanged++;
    continue;
  }

  data.summary = summary;
  fs.writeFileSync(file, JSON.stringify(data, null, 2) + '\n');
  updated++;
}

process.stdout.write(`STATUS: ok\n`);
process.stdout.write(`ITERATION: ${iterationPath}\n`);
process.stdout.write(`FILES_SCANNED: ${gradingFiles.length}\n`);
process.stdout.write(`UPDATED: ${updated}\n`);
process.stdout.write(`UNCHANGED: ${unchanged}\n`);
process.exit(0);
