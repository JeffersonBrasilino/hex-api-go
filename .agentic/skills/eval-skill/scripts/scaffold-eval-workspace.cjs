#!/usr/bin/env node
'use strict';

// Validates <skill_path>/evals/evals.json, optionally filters by --ids, derives a
// filesystem-safe slug per eval, detects the next iteration number by scanning
// evals/workspace/ for existing iteration-N/ directories, and scaffolds the
// with_skill/without_skill output directories for each selected eval. Replaces
// manual validation + slug derivation + directory creation previously done by
// Claude one Write-tool call at a time.
//
// Exit codes:
//   0 — workspace scaffolded (stdout: iteration + eval table)
//   1 — logic violation (evals.json invalid, or --ids matches nothing)
//   2 — usage or runtime error (skill path / evals.json not found)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/scaffold-eval-workspace.cjs <skill_path> [--ids N,M,...]

Reads <skill_path>/evals/evals.json, validates it, optionally filters entries by
--ids, derives a slug per eval, detects the next iteration number under
<skill_path>/evals/workspace/, and creates:

  evals/workspace/iteration-{N}/eval-{slug}/with_skill/outputs/.keep
  evals/workspace/iteration-{N}/eval-{slug}/without_skill/outputs/.keep

Exit codes:
  0   Workspace scaffolded (stdout: ITERATION, WORKSPACE, eval table)
  1   Logic violation — evals.json malformed, or --ids matches no entries
  2   Usage or runtime error (skill path or evals.json not found)

Examples:
  node scripts/scaffold-eval-workspace.cjs .agentic/skills/ddd-module-knowledge
  node scripts/scaffold-eval-workspace.cjs .agentic/skills/foo --ids 1,3,5
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
const skillArg = argv[0] && !argv[0].startsWith('--') ? argv[0].replace(/\/+$/, '') : null;

if (!skillArg) {
  fail2([
    'Error: <skill_path> is required.',
    'Usage: node scripts/scaffold-eval-workspace.cjs <skill_path> [--ids N,M,...]',
  ]);
}

let idsFilter = null;
const idsIdx = argv.indexOf('--ids');
if (idsIdx !== -1) {
  const raw = argv[idsIdx + 1];
  if (!raw || raw.startsWith('--')) {
    fail2(['Error: flag --ids requires a value.', 'Usage: node scripts/scaffold-eval-workspace.cjs <skill_path> [--ids N,M,...]']);
  }
  idsFilter = raw.split(',').map(s => s.trim()).filter(Boolean).map(Number);
}

const skillPath = skillArg;
const evalsJsonPath = path.join(skillPath, 'evals', 'evals.json');

if (!fs.existsSync(skillPath) || !fs.statSync(skillPath).isDirectory()) {
  fail2([
    `Error: skill directory not found: "${skillPath}"`,
    'Action: confirm the skill path is correct.',
  ]);
}

if (!fs.existsSync(evalsJsonPath)) {
  fail2([
    `Error: file not found: "${evalsJsonPath}"`,
    'Action: confirm the skill path is correct and evals.json exists before running evals.',
  ]);
}

let data;
try {
  data = JSON.parse(fs.readFileSync(evalsJsonPath, 'utf8'));
} catch (e) {
  fail2([
    `Error: "${evalsJsonPath}" is not valid JSON.`,
    `Details: ${e.message}`,
    'Action: fix the JSON syntax before retrying.',
  ]);
}

if (!Array.isArray(data.evals) || data.evals.length === 0) {
  fail1([
    'Error: evals.json has no entries in its "evals" array.',
    `File: ${evalsJsonPath}`,
    'Expected: { "evals": [ { "id", "prompt", "assertions" }, ... ] } with at least one entry.',
    'Action: add at least one eval entry before running evals.',
  ]);
}

const missingFields = [];
for (const [i, e] of data.evals.entries()) {
  const missing = ['id', 'prompt', 'assertions'].filter(f => e[f] === undefined);
  if (missing.length > 0) missingFields.push({ index: i, id: e.id, missing });
}
if (missingFields.length > 0) {
  fail1([
    'Error: one or more evals.json entries are missing required fields.',
    ...missingFields.map(m => `  entry index ${m.index} (id=${m.id ?? 'unknown'}): missing ${m.missing.join(', ')}`),
    'Expected: every entry has "id", "prompt", and "assertions".',
    'Action: fix the entries listed above before retrying.',
  ]);
}

let selected = data.evals;
if (idsFilter) {
  selected = data.evals.filter(e => idsFilter.includes(e.id));
  if (selected.length === 0) {
    fail1([
      `Error: evals.json has no entries matching --ids "${idsFilter.join(',')}".`,
      `Received IDs: ${idsFilter.join(', ')}`,
      `Available IDs: ${data.evals.map(e => e.id).join(', ')}`,
      'Action: pass valid eval IDs or omit --ids to run all evals.',
    ]);
  }
}

function slugify(description, id) {
  if (!description || !description.trim()) return String(id);
  const slug = description
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 60)
    .replace(/-+$/, '');
  return slug || String(id);
}

const withSlugs = selected.map(e => ({ ...e, slug: slugify(e.description, e.id) }));

const workspaceRoot = path.join(skillPath, 'evals', 'workspace');
let nextIteration = 1;
if (fs.existsSync(workspaceRoot)) {
  const existing = fs.readdirSync(workspaceRoot)
    .filter(name => /^iteration-\d+$/.test(name))
    .map(name => parseInt(name.slice('iteration-'.length), 10));
  if (existing.length > 0) nextIteration = Math.max(...existing) + 1;
}

const iterationRoot = path.join(workspaceRoot, `iteration-${nextIteration}`);

const rows = [];
for (const e of withSlugs) {
  const evalDir = path.join(iterationRoot, `eval-${e.slug}`);
  const withDir = path.join(evalDir, 'with_skill', 'outputs');
  const withoutDir = path.join(evalDir, 'without_skill', 'outputs');
  fs.mkdirSync(withDir, { recursive: true });
  fs.mkdirSync(withoutDir, { recursive: true });
  fs.writeFileSync(path.join(withDir, '.keep'), '');
  fs.writeFileSync(path.join(withoutDir, '.keep'), '');
  rows.push({ id: e.id, slug: e.slug, dir: evalDir });
}

const skillName = path.basename(skillPath);

process.stdout.write(`STATUS: ok\n`);
process.stdout.write(`SKILL: ${skillName}\n`);
process.stdout.write(`ITERATION: ${nextIteration}\n`);
process.stdout.write(`WORKSPACE: ${workspaceRoot}\n`);
process.stdout.write(`EVAL_COUNT: ${rows.length}\n`);
process.stdout.write('---\n');
for (const r of rows) {
  process.stdout.write(`${r.id} | eval-${r.slug} | ${r.dir}\n`);
}
process.exit(0);
