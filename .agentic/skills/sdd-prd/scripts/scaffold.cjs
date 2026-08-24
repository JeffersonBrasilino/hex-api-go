#!/usr/bin/env node
'use strict';

// Creates docs/<module>/<feature>/ and initialises PRD.md + NOTES.md from templates.
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

Creates docs/<module>/<feature>/ with PRD.md and NOTES.md pre-initialised
from skill templates. Idempotent — existing files are never overwritten.

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
const targetDir   = path.resolve(projectRoot, 'docs', moduleName, featureName);
const skillDir    = path.resolve(projectRoot, '.agentic', 'skills', 'sdd-prd');
const today       = new Date().toISOString().slice(0, 10);

if (!targetDir.startsWith(projectRoot + path.sep)) {
  process.stderr.write([
    `Error: path traversal detected — target directory is outside project root.`,
    `Project root: "${projectRoot}"`,
    `Resolved target: "${targetDir}"`,
    `Action: use simple kebab-case names without ".." or absolute paths.`,
  ].join('\n') + '\n');
  process.exit(1);
}

function extractTemplateBlock(raw) {
  const m = raw.match(/```markdown\n([\s\S]*?)```/);
  return m ? m[1] : raw;
}

try {
  const dirCreated = !fs.existsSync(targetDir);
  if (dirCreated) fs.mkdirSync(targetDir, { recursive: true });

  const prdPath    = path.join(targetDir, 'PRD.md');
  const prdCreated = !fs.existsSync(prdPath);
  if (prdCreated) {
    const raw     = fs.readFileSync(path.join(skillDir, 'references', 'prd-template.md'), 'utf8');
    const content = extractTemplateBlock(raw)
      .replace('[YYYY-MM-DD]', today)
      .replace('[Nome]', 'Jefferson Brasilino')
      .replace('Draft / Em Revisão / Aprovado', 'Draft')
      .replace('[Versão]', '0.1')
      .replace('feat / fix / refactor / perf / chore / docs / test / build / ci', 'feat');
    fs.writeFileSync(prdPath, content, 'utf8');
  }

  const notesPath    = path.join(targetDir, 'NOTES.md');
  const notesCreated = !fs.existsSync(notesPath);
  if (notesCreated) {
    const raw     = fs.readFileSync(path.join(skillDir, 'references', 'notes-template.md'), 'utf8');
    const content = extractTemplateBlock(raw).replace(/\[YYYY-MM-DD\]/g, today);
    fs.writeFileSync(notesPath, content, 'utf8');
  }

  process.stdout.write(`DIR: docs/${moduleName}/${featureName}/ ${dirCreated ? '(criado)' : '(já existia)'}\n`);
  process.stdout.write(`PRD: docs/${moduleName}/${featureName}/PRD.md ${prdCreated ? '(criado)' : '(ignorado — já existia)'}\n`);
  process.stdout.write(`NOTES: docs/${moduleName}/${featureName}/NOTES.md ${notesCreated ? '(criado)' : '(ignorado — já existia)'}\n`);
  process.exit(0);
} catch (err) {
  process.stderr.write([
    `Error: filesystem operation failed.`,
    `Detail: ${err.message}`,
    `Action: check write permissions on "${path.resolve(projectRoot, 'docs')}"`,
  ].join('\n') + '\n');
  process.exit(2);
}
