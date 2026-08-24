#!/usr/bin/env node
'use strict';

// Creates docs/<module>/<feature>/ and initialises PLAN.md from plan-schema.md.
// Idempotent: skips existing files, never overwrites.
//
// Exit codes:
//   0 — scaffold complete (stdout: created paths + status of each file)
//   1 — invalid arguments (stderr: what was wrong + expected format)
//   2 — runtime or filesystem error (stderr: what failed)

const fs   = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/scaffold.cjs <module-name> <feature-name>

Creates docs/<module>/<feature>/ with PLAN.md pre-initialised from the
plan-schema.md template (Author, Date, Status filled in). Idempotent —
existing files are never overwritten.

Arguments:
  module-name    DDD module in kebab-case (e.g. auth, user, order)
  feature-name   Feature in kebab-case (e.g. login, reset-password)

Exit codes:
  0   Scaffold complete
  1   Invalid arguments (bad format, path traversal)
  2   Filesystem or runtime error

Examples:
  node scripts/scaffold.cjs auth login
  node scripts/scaffold.cjs user reset-password
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [,, moduleName, featureName] = process.argv;

if (!moduleName || !featureName) {
  process.stderr.write([
    `Error: both <module-name> and <feature-name> are required.`,
    `Received: ${JSON.stringify({ moduleName, featureName })}`,
    `Expected: kebab-case strings, e.g. "auth" "login"`,
    `Usage: node scripts/scaffold.cjs <module-name> <feature-name>`,
  ].join('\n') + '\n');
  process.exit(1);
}

const KEBAB = /^[a-z][a-z0-9-]*$/;
if (!KEBAB.test(moduleName)) {
  process.stderr.write([
    `Error: <module-name> must be kebab-case (lowercase letters, digits, hyphens).`,
    `Received: "${moduleName}"`,
    `Expected: e.g. "auth", "user-account", "order"`,
  ].join('\n') + '\n');
  process.exit(1);
}
if (!KEBAB.test(featureName)) {
  process.stderr.write([
    `Error: <feature-name> must be kebab-case (lowercase letters, digits, hyphens).`,
    `Received: "${featureName}"`,
    `Expected: e.g. "login", "reset-password", "create-order"`,
  ].join('\n') + '\n');
  process.exit(1);
}

const projectRoot = process.cwd();
const targetDir    = path.resolve(projectRoot, 'docs', moduleName, featureName);
const skillDir     = path.resolve(projectRoot, '.agentic', 'skills', 'sdd-plan');
const today        = new Date().toISOString().slice(0, 10);

if (!targetDir.startsWith(projectRoot + path.sep)) {
  process.stderr.write([
    `Error: path traversal detected — target directory is outside project root.`,
    `Project root: "${projectRoot}"`,
    `Resolved target: "${targetDir}"`,
    `Action: use simple kebab-case names without ".." or absolute paths.`,
  ].join('\n') + '\n');
  process.exit(1);
}

const CONVENTIONAL_TYPES = ['feat', 'fix', 'refactor', 'perf', 'chore', 'docs', 'test', 'build', 'ci'];

// Propagates the Conventional Commits classification chosen during the PRD interview
// (PRD.md's "Tipo" header row) into PLAN.md's "Type" header — this is never re-asked
// at plan time. Falls back to "feat" when the PRD doesn't exist yet or its Tipo row
// isn't one of the known types (e.g. still the unfilled template placeholder).
function readTypeFromPrd(prdPath) {
  if (!fs.existsSync(prdPath)) return 'feat';
  const prdContent = fs.readFileSync(prdPath, 'utf8');
  const m = prdContent.match(/\|\s*\*\*Tipo\*\*\s*\|\s*([a-z]+)\s*\|/i);
  const value = m ? m[1].toLowerCase() : null;
  return value && CONVENTIONAL_TYPES.includes(value) ? value : 'feat';
}

try {
  const dirCreated = !fs.existsSync(targetDir);
  if (dirCreated) fs.mkdirSync(targetDir, { recursive: true });

  const planPath    = path.join(targetDir, 'PLAN.md');
  const planCreated = !fs.existsSync(planPath);
  if (planCreated) {
    const type    = readTypeFromPrd(path.join(targetDir, 'PRD.md'));
    const raw     = fs.readFileSync(path.join(skillDir, 'references', 'plan-schema.md'), 'utf8');
    const content = raw
      .replace('**Date:** [YYYY-MM-DD]', `**Date:** ${today}`)
      .replace('**Author:** [Dev/Agent Name]', '**Author:** Jefferson Brasilino')
      .replace('**Status:** `Draft` | `Planning` | `Ready for Implementation` | `Done`', '**Status:** `Draft`')
      .replace('**Type:** `feat` | `fix` | `refactor` | `perf` | `chore` | `docs` | `test` | `build` | `ci`', `**Type:** \`${type}\``);
    fs.writeFileSync(planPath, content, 'utf8');
  }

  process.stdout.write(`DIR: docs/${moduleName}/${featureName}/ ${dirCreated ? '(criado)' : '(já existia)'}\n`);
  process.stdout.write(`PLAN: docs/${moduleName}/${featureName}/PLAN.md ${planCreated ? '(criado)' : '(ignorado — já existia)'}\n`);
  process.exit(0);
} catch (err) {
  process.stderr.write([
    `Error: filesystem operation failed.`,
    `Detail: ${err.message}`,
    `Action: check write permissions on "${path.resolve(projectRoot, 'docs')}"`,
  ].join('\n') + '\n');
  process.exit(2);
}
