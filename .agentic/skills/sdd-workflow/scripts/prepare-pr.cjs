#!/usr/bin/env node
'use strict';

// Prepares (and, where possible, pushes) the git branch/commit for a feature whose
// SDD pipeline just reached phase "done", ahead of PR/MR creation. Deterministic —
// no LLM judgment calls; scope of files, branch naming, and provider are all read
// from PLAN.md and sdd-workflow.config.json, never guessed.
//
// This script never creates the PR/MR itself — it only gets the branch, commit, and
// push ready and reports the structured task/PRD data that a separate lightweight
// agent (pr-writer) uses to write the PR title/body and call `gh`/`glab`. If the
// required CLI tool is missing, this script does NOT touch git at all — it reports
// blocked_missing_tool and stops, so a later re-run (once the tool is installed)
// starts from a clean slate.
//
// Exit codes:
//   0 — completed a step; stdout STATUS line says which (see below)
//   1 — logic violation (invalid config value, malformed PLAN.md)
//   2 — usage or runtime error (missing argument, feature directory not found)

const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');

const HELP = `
Usage: node scripts/prepare-pr.cjs <feature_path>

Gets the feature branch and commit ready for a PR/MR once a feature's SDD pipeline
reaches phase "done": validates the gh/glab CLI is available, creates or reuses the
feature branch (per sdd-workflow.config.json), commits only the files touched by
tasks marked Done/[y] in PLAN.md, and pushes. Reports structured data for the
pr-writer agent to consume — it does not create the PR/MR itself.

STATUS values (stdout, exit 0 in all cases below):
  blocked_missing_tool      Required CLI (gh/glab) not found — nothing was touched.
  blocked_protected_branch  Current branch is the base branch and auto_create_branch
                             is false — nothing was touched.
  already_open              A PR/MR for this branch already exists (resume case) —
                             PR_URL is reported, no new commit/push attempted.
  ready                     Branch/commit/push done (or already up to date); task
                             list + PRD path reported for pr-writer to act on.

Exit codes:
  0   See STATUS above.
  1   Logic violation — invalid config value, malformed PLAN.md.
  2   Usage or runtime error (missing argument, feature directory not found).
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

const featurePath = process.argv[2] && !process.argv[2].startsWith('--')
  ? process.argv[2].replace(/\/+$/, '')
  : null;

if (!featurePath) {
  fail2(['Error: <feature_path> is required.', 'Usage: node scripts/prepare-pr.cjs <feature_path>']);
}
if (!fs.existsSync(featurePath) || !fs.statSync(featurePath).isDirectory()) {
  fail2([`Error: feature directory not found: "${featurePath}"`]);
}

const planPath = path.join(featurePath, 'PLAN.md');
if (!fs.existsSync(planPath)) {
  fail2([`Error: file not found: "${planPath}"`]);
}

const prdPath = path.join(featurePath, 'PRD.md');

function git(args) {
  return execFileSync('git', args, { encoding: 'utf8' }).trim();
}
function gitOrNull(args) {
  try {
    return git(args);
  } catch {
    return null;
  }
}
function commandExists(cmd) {
  try {
    execFileSync('which', [cmd], { stdio: 'ignore' });
    return true;
  } catch {
    return false;
  }
}

// --- config -------------------------------------------------------------

const CONFIG_PATH = path.join(__dirname, '..', 'assets', 'sdd-workflow.config.json');
const DEFAULT_CONFIG = {
  git: {
    provider: 'auto',
    base_branch: 'main',
    auto_create_branch: true,
    branch_pattern: '{type}/{feature-slug}',
    default_type: 'feat',
  },
};

const CONVENTIONAL_TYPES = ['feat', 'fix', 'refactor', 'perf', 'chore', 'docs', 'test', 'build', 'ci'];

function loadConfig() {
  if (!fs.existsSync(CONFIG_PATH)) return DEFAULT_CONFIG;
  let raw;
  try {
    raw = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
  } catch (e) {
    fail1([
      `Error: "${CONFIG_PATH}" is not valid JSON.`,
      `Received: ${e.message}`,
      'Action: fix the JSON or delete the file to fall back to defaults.',
    ]);
  }
  return {
    git: {
      provider: raw.git?.provider || DEFAULT_CONFIG.git.provider,
      base_branch: raw.git?.base_branch || DEFAULT_CONFIG.git.base_branch,
      auto_create_branch: raw.git?.auto_create_branch !== undefined
        ? raw.git.auto_create_branch
        : DEFAULT_CONFIG.git.auto_create_branch,
      branch_pattern: raw.git?.branch_pattern || DEFAULT_CONFIG.git.branch_pattern,
      default_type: raw.git?.default_type || DEFAULT_CONFIG.git.default_type,
    },
  };
}

const config = loadConfig();
const VALID_PROVIDERS = ['auto', 'github', 'gitlab', 'none'];
if (!VALID_PROVIDERS.includes(config.git.provider)) {
  fail1([
    `Error: sdd-workflow.config.json git.provider "${config.git.provider}" is not valid.`,
    `Expected: one of ${VALID_PROVIDERS.join(' | ')}`,
  ]);
}
if (!CONVENTIONAL_TYPES.includes(config.git.default_type)) {
  fail1([
    `Error: sdd-workflow.config.json git.default_type "${config.git.default_type}" is not valid.`,
    `Expected: one of ${CONVENTIONAL_TYPES.join(' | ')}`,
  ]);
}

function detectProvider() {
  if (config.git.provider !== 'auto') return config.git.provider;
  const remote = gitOrNull(['remote', 'get-url', 'origin']) || '';
  if (/github\.com/.test(remote)) return 'github';
  if (/gitlab/.test(remote)) return 'gitlab';
  return 'none';
}

const provider = detectProvider();
const CLI_BY_PROVIDER = { github: 'gh', gitlab: 'glab' };

// --- Step 1: tool availability (before touching git at all) --------------

if (provider !== 'none') {
  const cli = CLI_BY_PROVIDER[provider];
  if (!commandExists(cli)) {
    process.stdout.write('STATUS: blocked_missing_tool\n');
    process.stdout.write(`PROVIDER: ${provider}\n`);
    process.stdout.write(`MISSING_CLI: ${cli}\n`);
    process.stdout.write(
      provider === 'github'
        ? 'INSTALL: brew install gh && gh auth login\n'
        : 'INSTALL: brew install glab && glab auth login\n'
    );
    process.exit(0);
  }
}

// --- Step 2: Conventional Commits type (from PLAN.md header, else config default) ---

const planContent = fs.readFileSync(planPath, 'utf8');
const typeMatch = planContent.match(/^\*\*Type:\*\*\s*`([a-z]+)`/m);
let type = typeMatch ? typeMatch[1] : config.git.default_type;
if (!CONVENTIONAL_TYPES.includes(type)) {
  fail1([
    `Error: PLAN.md "**Type:**" is "${type}", which is not a valid Conventional Commits type.`,
    `File: ${planPath}`,
    `Expected: one of ${CONVENTIONAL_TYPES.join(' | ')}`,
    'Action: fix the Type field in the PLAN.md header before retrying.',
  ]);
}

// --- Step 3: branch ------------------------------------------------------

const currentBranch = git(['rev-parse', '--abbrev-ref', 'HEAD']);
const remoteHead = gitOrNull(['symbolic-ref', '--short', 'refs/remotes/origin/HEAD'])
  ?.replace(/^origin\//, '') || config.git.base_branch;

const onBaseBranch = currentBranch === remoteHead || currentBranch === config.git.base_branch;

let branch = currentBranch;
let branchCreated = false;

if (onBaseBranch) {
  if (!config.git.auto_create_branch) {
    process.stdout.write('STATUS: blocked_protected_branch\n');
    process.stdout.write(`CURRENT_BRANCH: ${currentBranch}\n`);
    process.stdout.write('MESSAGE: current branch is the base branch and auto_create_branch is false in sdd-workflow.config.json — create a feature branch manually and re-run.\n');
    process.exit(0);
  }
  const slug = featurePath.replace(/\/+$/, '').split('/').pop();
  branch = config.git.branch_pattern.replace('{type}', type).replace('{feature-slug}', slug);
  git(['checkout', '-b', branch]);
  branchCreated = true;
}

// --- Step 4: file list from PLAN.md Done/[y] tasks ----------------------
const taskBlockRe = /^- \[([xy])\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—\s*(.*?)\*\*\s*$/gm;

const tasks = [];
let m;
while ((m = taskBlockRe.exec(planContent)) !== null) {
  const [, , id, title] = m;
  const blockStart = m.index + m[0].length;
  const nextMatch = planContent.slice(blockStart).match(/^- \[[ x/y]\] \*\*TASK-/m);
  const blockEnd = nextMatch ? blockStart + nextMatch.index : planContent.length;
  const block = planContent.slice(blockStart, blockEnd);
  const fileMatch = block.match(/\*\*File:\*\*\s*`([^`]+)`/);
  if (!fileMatch) continue;
  tasks.push({ id, title, file: fileMatch[1] });
}

if (tasks.length === 0) {
  fail1([
    `Error: no Done/[y] task with a File: field found in "${planPath}".`,
    'Action: confirm the pipeline actually completed implementation before this step runs.',
  ]);
}

const files = [...new Set(tasks.map(t => t.file))].filter(f => fs.existsSync(f));

// --- Step 5: commit (only if there is something staged) -----------------

if (files.length > 0) {
  git(['add', ...files]);
}
// `git diff --cached --quiet` exits 1 (throws, caught as null) when there IS a
// staged diff, and exits 0 (returns '') when there is none.
const hasStagedChanges = files.length > 0 && gitOrNull(['diff', '--cached', '--quiet']) === null;
let committed = false;
if (hasStagedChanges) {
  const parts = featurePath.split('/').filter(Boolean);
  const module = parts[1] || parts[0];
  const slug = parts[parts.length - 1];
  const commitTitle = `${type}(${module}): implement ${slug}`;
  const commitBody = tasks.map(t => `- ${t.id} — ${t.title}`).join('\n');
  git(['commit', '-m', commitTitle, '-m', commitBody]);
  committed = true;
}

// --- Step 6: push ---------------------------------------------------------

git(['push', '-u', 'origin', branch]);

// --- Step 7: check for an already-open PR (resume case) -------------------

let existingUrl = null;
if (provider === 'github') {
  try {
    existingUrl = execFileSync('gh', ['pr', 'view', branch, '--json', 'url', '-q', '.url'], { encoding: 'utf8' }).trim() || null;
  } catch {
    existingUrl = null;
  }
} else if (provider === 'gitlab') {
  try {
    const out = execFileSync('glab', ['mr', 'view', branch], { encoding: 'utf8' });
    const urlMatch = out.match(/https?:\/\/\S+/);
    existingUrl = urlMatch ? urlMatch[0] : null;
  } catch {
    existingUrl = null;
  }
}

if (existingUrl) {
  process.stdout.write('STATUS: already_open\n');
  process.stdout.write(`PR_URL: ${existingUrl}\n`);
  process.stdout.write(`BRANCH: ${branch}\n`);
  process.exit(0);
}

// --- Report -----------------------------------------------------------

process.stdout.write('STATUS: ready\n');
process.stdout.write(`PROVIDER: ${provider}\n`);
process.stdout.write(`TYPE: ${type}\n`);
process.stdout.write(`BRANCH: ${branch}\n`);
process.stdout.write(`BRANCH_CREATED: ${branchCreated}\n`);
process.stdout.write(`BASE_BRANCH: ${config.git.base_branch}\n`);
process.stdout.write(`COMMITTED: ${committed}\n`);
process.stdout.write(`PRD_PATH: ${fs.existsSync(prdPath) ? prdPath : '(none)'}\n`);
process.stdout.write('TASKS:\n');
for (const t of tasks) {
  process.stdout.write(`  ${t.id} — ${t.title} (${t.file})\n`);
}
process.exit(0);
