#!/usr/bin/env node
'use strict';

// Writes .agentic/skills/sdd-workflow/assets/sdd-workflow.config.json from validated flags.
// Called once by references/setup-config.md after the user answers the setup interview (or again
// if the user explicitly asks to reconfigure). Deterministic — validates enums the same way
// prepare-pr.cjs does, never trusts the caller's values blindly.
//
// Exit codes:
//   0 — config written; stdout reports the path.
//   1 — invalid value for one of the enum fields.
//   2 — usage error (unknown flag).

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/write-config.cjs [flags]

Writes assets/sdd-workflow.config.json from the given flags, defaulting any flag that is
omitted. Overwrites the file if it already exists (the caller decides whether that's wanted —
this script never asks).

Flags (all optional, defaults shown):
  --git-provider <auto|github|gitlab|none>        (auto)
  --base-branch <name>                            (main)
  --auto-create-branch <true|false>                (true)
  --branch-pattern <pattern>                       ({type}/{feature-slug})
  --default-type <feat|fix|refactor|perf|chore|docs|test|build|ci>  (feat)
  --archive-destination-type <mcp|skill|path|none> (none)
  --archive-destination-name <name>                ("")
                                                    (MCP/skill name for mcp|skill;
                                                     folder path for path; "" for none)
  --prd-provider <jira|github|trello|none>         (none)
  --prd-board-url <url>                            ("")
  --prd-project-key <key>                          ("")
  --prd-list-id <id>                               ("")
  --prd-repo <owner/repo>                          ("")

Exit codes:
  0   Config written. Path reported on stdout.
  1   Invalid value for an enum flag.
  2   Usage error (unrecognized flag).
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

const KNOWN_FLAGS = new Set([
  '--git-provider', '--base-branch', '--auto-create-branch', '--branch-pattern', '--default-type',
  '--archive-destination-type', '--archive-destination-name',
  '--prd-provider', '--prd-board-url', '--prd-project-key', '--prd-list-id', '--prd-repo',
]);

function parseFlags(argv) {
  const out = {};
  for (let i = 0; i < argv.length; i += 1) {
    const flag = argv[i];
    if (!flag.startsWith('--')) continue;
    if (!KNOWN_FLAGS.has(flag)) fail2([`Error: unrecognized flag "${flag}".`]);
    out[flag] = argv[i + 1] !== undefined ? argv[i + 1] : '';
    i += 1;
  }
  return out;
}

const flags = parseFlags(process.argv.slice(2));

const config = {
  git: {
    provider: flags['--git-provider'] || 'auto',
    base_branch: flags['--base-branch'] || 'main',
    auto_create_branch: flags['--auto-create-branch'] !== undefined
      ? flags['--auto-create-branch'] !== 'false'
      : true,
    branch_pattern: flags['--branch-pattern'] || '{type}/{feature-slug}',
    default_type: flags['--default-type'] || 'feat',
  },
  archive_spec: {
    destination_type: flags['--archive-destination-type'] || 'none',
    destination_name: flags['--archive-destination-name'] || '',
  },
  prd: {
    provider: flags['--prd-provider'] || 'none',
    board_url: flags['--prd-board-url'] || '',
    project_key: flags['--prd-project-key'] || '',
    list_id: flags['--prd-list-id'] || '',
    repo: flags['--prd-repo'] || '',
  },
};

const VALID_GIT_PROVIDERS = ['auto', 'github', 'gitlab', 'none'];
const VALID_COMMIT_TYPES = ['feat', 'fix', 'refactor', 'perf', 'chore', 'docs', 'test', 'build', 'ci'];
const VALID_ARCHIVE_TYPES = ['mcp', 'skill', 'path', 'none'];
const VALID_PRD_PROVIDERS = ['jira', 'github', 'trello', 'none'];

if (!VALID_GIT_PROVIDERS.includes(config.git.provider)) {
  fail1([`Error: git.provider "${config.git.provider}" is not valid.`, `Expected: one of ${VALID_GIT_PROVIDERS.join(' | ')}`]);
}
if (!VALID_COMMIT_TYPES.includes(config.git.default_type)) {
  fail1([`Error: git.default_type "${config.git.default_type}" is not valid.`, `Expected: one of ${VALID_COMMIT_TYPES.join(' | ')}`]);
}
if (!VALID_ARCHIVE_TYPES.includes(config.archive_spec.destination_type)) {
  fail1([`Error: archive_spec.destination_type "${config.archive_spec.destination_type}" is not valid.`, `Expected: one of ${VALID_ARCHIVE_TYPES.join(' | ')}`]);
}
if (!VALID_PRD_PROVIDERS.includes(config.prd.provider)) {
  fail1([`Error: prd.provider "${config.prd.provider}" is not valid.`, `Expected: one of ${VALID_PRD_PROVIDERS.join(' | ')}`]);
}

const ASSETS_DIR = path.join(__dirname, '..', 'assets');
const CONFIG_PATH = path.join(ASSETS_DIR, 'sdd-workflow.config.json');

fs.mkdirSync(ASSETS_DIR, { recursive: true });
fs.writeFileSync(CONFIG_PATH, JSON.stringify(config, null, 2) + '\n');

process.stdout.write(`STATUS: written\nPATH: ${CONFIG_PATH}\n`);
