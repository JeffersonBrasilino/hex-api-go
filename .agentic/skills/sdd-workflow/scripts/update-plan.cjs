#!/usr/bin/env node
'use strict';

// Marks task checkboxes ([x]/[y]) and updates the header Status field of a
// PLAN.md deterministically, instead of Claude hand-editing the file (risking
// touching adjacent content or using the wrong mark).
//
// Exit codes:
//   0 — PLAN.md updated (stdout: what changed, what was already correct)
//   1 — logic violation (unknown TASK-ID, invalid --status value)
//   2 — usage or runtime error (missing argument, file not found)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/update-plan.cjs <feature_path> [flags]

Marks task checkboxes and/or updates the header Status of <feature_path>/PLAN.md.
Updates checkboxes in both Section 1 (Plan) and Section 3 (Execution — Validated
Checklist) when a matching block exists in Section 3; Section 1 is authoritative
for validating that a TASK-ID exists.

Flags:
  --mark-done <IDs>      Mark [x] for the comma-separated TASK-IDs
  --mark-yolo <IDs>       Mark [y] for the comma-separated TASK-IDs
  --mark-failed <IDs>      Set Validation Status to "Failed" for the TASK-IDs — the
                            permanent give-up state, only after retries are exhausted
                            (in Section 3, if a mirrored block exists there)
  --status <value>          Update the header Status field
                            (Draft | Planning | Ready for Implementation |
                             Implementando | Done)

Exit codes:
  0   PLAN.md updated
  1   Logic violation — unknown TASK-ID or invalid --status value
  2   Usage or runtime error (feature_path missing, PLAN.md not found)

Examples:
  node scripts/update-plan.cjs docs/auth/login --mark-done TASK-DOM-A,TASK-APP-B --status Implementando
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const VALID_STATUSES = ['Draft', 'Planning', 'Ready for Implementation', 'Implementando', 'Done'];

function fail1(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(1);
}
function fail2(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(2);
}

const argv = process.argv.slice(2);
const featurePath = argv[0] && !argv[0].startsWith('--') ? argv[0].replace(/\/+$/, '') : null;

if (!featurePath) {
  fail2(['Error: <feature_path> is required.', 'Usage: node scripts/update-plan.cjs <feature_path> [flags]']);
}

const flags = {};
for (let i = 1; i < argv.length; i++) {
  const tok = argv[i];
  if (!tok.startsWith('--')) continue;
  const name = tok.slice(2);
  const value = argv[i + 1];
  if (value === undefined || value.startsWith('--')) {
    fail2([`Error: flag --${name} requires a value.`, 'Usage: node scripts/update-plan.cjs <feature_path> [flags]']);
  }
  flags[name] = value;
  i++;
}

const planPath = path.join(featurePath, 'PLAN.md');
if (!fs.existsSync(planPath)) {
  fail2([
    `Error: file not found: "${planPath}"`,
    'Action: confirm the feature path is correct and PLAN.md was created by sdd-plan.',
  ]);
}

let content = fs.readFileSync(planPath, 'utf8');

function escapeRe(s) {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function declaredTaskIds(text) {
  const ids = new Set();
  const re = /^- \[[ x/y]\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
  let m;
  while ((m = re.exec(text)) !== null) ids.add(m[1]);
  return ids;
}

const declared = declaredTaskIds(content);

function parseIdList(raw) {
  return raw.split(',').map(s => s.trim()).filter(Boolean);
}

function validateIds(ids) {
  const unknown = ids.filter(id => !declared.has(id));
  if (unknown.length > 0) {
    fail1([
      `Error: task not found in PLAN.md: "${unknown[0]}"`,
      `File: ${planPath}`,
      `Expected: a task block matching "- [ ] **${unknown[0]} —"`,
      'Action: confirm the task ID exists in Section 1 of PLAN.md before retrying.',
    ]);
  }
}

function markCheckbox(text, id, mark) {
  const re = new RegExp(`^(- )\\[([ x/y])\\](\\s*\\*\\*${escapeRe(id)}\\s*—)`, 'gm');
  let changed = false;
  let alreadyCorrect = false;
  const out = text.replace(re, (full, prefix, currentMark, suffix) => {
    if (currentMark === mark) {
      alreadyCorrect = true;
      return full;
    }
    changed = true;
    return `${prefix}[${mark}]${suffix}`;
  });
  return { text: out, changed, alreadyCorrect };
}

function setValidationStatus(text, id, value) {
  // Locate this task's block (from its checkbox line to the next checkbox line
  // or end of Section 3), and update **Validation Status:** if present there.
  const startRe = new RegExp(`^- \\[[ x/y]\\] \\*\\*${escapeRe(id)}\\s*—.*$`, 'm');
  const startMatch = text.match(startRe);
  if (!startMatch) return { text, found: false };

  const startIdx = startMatch.index;
  const rest = text.slice(startIdx + startMatch[0].length);
  const nextTaskMatch = rest.match(/^- \[[ x/y]\] \*\*TASK-/m);
  const blockEnd = nextTaskMatch ? startIdx + startMatch[0].length + nextTaskMatch.index : text.length;
  const block = text.slice(startIdx, blockEnd);

  const vsRe = /(\*\*Validation Status:\*\*\s*)\S.*$/m;
  if (!vsRe.test(block)) return { text, found: false };

  const updatedBlock = block.replace(vsRe, `$1${value}`);
  return { text: text.slice(0, startIdx) + updatedBlock + text.slice(blockEnd), found: true };
}

const markedDone = [];
const markedYolo = [];
const alreadyDone = [];
const markedFailed = [];
const failedNoBlock = [];

if (flags['mark-done']) {
  const ids = parseIdList(flags['mark-done']);
  validateIds(ids);
  for (const id of ids) {
    const result = markCheckbox(content, id, 'x');
    content = result.text;
    if (result.alreadyCorrect) alreadyDone.push(id);
    else if (result.changed) markedDone.push(id);
  }
}

if (flags['mark-yolo']) {
  const ids = parseIdList(flags['mark-yolo']);
  validateIds(ids);
  for (const id of ids) {
    const result = markCheckbox(content, id, 'y');
    content = result.text;
    if (result.alreadyCorrect) alreadyDone.push(id);
    else if (result.changed) markedYolo.push(id);
  }
}

if (flags['mark-failed']) {
  const ids = parseIdList(flags['mark-failed']);
  validateIds(ids);
  for (const id of ids) {
    const result = setValidationStatus(content, id, 'Failed');
    content = result.text;
    if (result.found) markedFailed.push(id);
    else failedNoBlock.push(id);
  }
}

let statusUpdated = null;
if (flags['status']) {
  const value = flags['status'];
  if (!VALID_STATUSES.includes(value)) {
    fail1([
      `Error: --status "${value}" is not a valid PLAN.md status.`,
      `Received: "${value}"`,
      `Expected: ${VALID_STATUSES.join(' | ')}`,
      'Action: use one of the valid status values.',
    ]);
  }
  const headerRe = /(\*\*Status:\*\*\s*`)[^`]*(`)/;
  if (headerRe.test(content)) {
    content = content.replace(headerRe, `$1${value}$2`);
    statusUpdated = value;
  }
}

fs.writeFileSync(planPath, content);

process.stdout.write('STATUS: ok\n');
process.stdout.write(`FILE: ${planPath}\n`);
if (markedDone.length) process.stdout.write(`MARKED_DONE: ${markedDone.join(', ')}\n`);
if (markedYolo.length) process.stdout.write(`MARKED_YOLO: ${markedYolo.join(', ')}\n`);
if (alreadyDone.length) process.stdout.write(`ALREADY_DONE: ${alreadyDone.join(', ')} (ignorado — já estava marcado)\n`);
if (markedFailed.length) process.stdout.write(`MARKED_FAILED: ${markedFailed.join(', ')}\n`);
if (failedNoBlock.length) process.stdout.write(`NO_SECTION3_BLOCK: ${failedNoBlock.join(', ')} (Validation Status não encontrado — apenas Section 1 existe ainda)\n`);
if (statusUpdated) process.stdout.write(`STATUS_UPDATED: ${statusUpdated}\n`);
process.exit(0);
