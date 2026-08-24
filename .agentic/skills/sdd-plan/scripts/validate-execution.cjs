#!/usr/bin/env node
'use strict';

// Validates the "Execution — Validated Checklist" (Section 3) of a PLAN.md file
// for a given set of task IDs — the structural sub-checks of PLAN.md compliance
// that don't require reading the implementation (that judgment stays with the
// calling agent). Checks per task:
//   - An execution block exists in Section 3 for the task ID.
//   - Validation Status is set to "Implemented" by the implement-agent — the
//     precondition for verify-code to run its own checklist (not "Verified"/"Done":
//     those are written by verify-code and verify-wave-prd respectively, after this
//     script has passed).
//   - The File declared in the Section 1 task block appears in Section 3's
//     Files Modified list.
//
// Exit codes:
//   0 — all requested tasks pass every check (stdout: structured summary)
//   1 — one or more checks failed (stderr: numbered issue list)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/validate-execution.cjs <path/to/PLAN.md> <TASK-ID> [TASK-ID...]

Validates Section 3 ("Execution — Validated Checklist") of PLAN.md for the given
task IDs:
  - An execution block exists for the task ID.
  - Validation Status is "Implemented" (not pending/Blocked/Verification Failed/
    Failed/Done/unfilled).
  - The File declared in the task's Section 1 block appears in Section 3's
    Files Modified list.

Does NOT check whether sub-tasks are actually implemented in code — that
requires reading the implementation and stays a judgment call for the caller.

Exit codes:
  0   All tasks pass every check
  1   One or more checks failed — detailed issues on stderr
  2   Usage or runtime error (file not found, missing argument)

Examples:
  node scripts/validate-execution.cjs docs/auth/login/PLAN.md TASK-DOM-PASSWORD TASK-APP-LOGIN-HANDLER
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , planPath, ...taskIds] = process.argv;

if (!planPath || taskIds.length === 0) {
  process.stderr.write([
    'Error: <path/to/PLAN.md> and at least one <TASK-ID> are required.',
    'Usage: node scripts/validate-execution.cjs <path/to/PLAN.md> <TASK-ID> [TASK-ID...]',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(planPath)) {
  process.stderr.write(`Error: file not found: "${planPath}"\n`);
  process.exit(2);
}

const content = fs.readFileSync(planPath, 'utf8');

function extractSection(name, nextName) {
  const re = new RegExp(`## ${name}([\\s\\S]*?)(?=\\n## ${nextName}|$)`);
  const match = content.match(re);
  return match ? match[1] : '';
}

const section1 = extractSection('1\\. Plan', '3\\. Execution');
const section3 = extractSection('3\\. Execution', '4\\. Return');

// Section 1: map task ID -> declared File.
const TASK1_START = /^- \[[ x/]\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts1 = [];
let m1;
while ((m1 = TASK1_START.exec(section1)) !== null) starts1.push({ id: m1[1], index: m1.index });
const blocks1 = starts1.map((s, i) => {
  const end = i + 1 < starts1.length ? starts1[i + 1].index : section1.length;
  return { id: s.id, text: section1.slice(s.index, end) };
});
const fileById = {};
for (const block of blocks1) {
  const fileMatch = block.text.match(/\*\*File:\*\*\s*`([^`]+)`/);
  fileById[block.id] = fileMatch ? fileMatch[1].trim() : null;
}

// Section 3: map task ID -> execution block.
const TASK3_START = /^- \[[ x/]\] \*\*Execution — (TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+):/gm;
const starts3 = [];
let m3;
while ((m3 = TASK3_START.exec(section3)) !== null) starts3.push({ id: m3[1], index: m3.index });
const blocks3 = starts3.map((s, i) => {
  const end = i + 1 < starts3.length ? starts3[i + 1].index : section3.length;
  return { id: s.id, text: section3.slice(s.index, end) };
});
const executionById = {};
for (const block of blocks3) executionById[block.id] = block.text;

const issues = [];
const results = [];

for (const taskId of taskIds) {
  const execBlock = executionById[taskId];

  if (!execBlock) {
    issues.push({
      type: 'missing execution block',
      taskId,
      expected: `an "Execution — ${taskId}: ..." block in Section 3`,
      action: `Fill in the Section 3 execution block for ${taskId} before verifying.`,
    });
    results.push({ taskId, checklistFilled: false, statusReady: false, fileMatch: false });
    continue;
  }

  const statusMatch = execBlock.match(/\*Validation Status:\*\s*`?([^`\n]+)`?/);
  const status = statusMatch ? statusMatch[1].trim() : null;
  const isBlocked = !!status && /^Blocked\b/i.test(status);
  const isVerificationFailed = !!status && /^Verification Failed$/i.test(status);
  const statusReady = !!status && /^Implemented$/i.test(status);
  if (isBlocked) {
    issues.push({
      type: 'execution blocked',
      taskId,
      found: status,
      expected: 'Implemented',
      action: `Task ${taskId} is Blocked — resolve the contradiction described in its Agent Notes/Validation Evidence before verifying; do not attempt to verify a Blocked task.`,
    });
  } else if (isVerificationFailed) {
    issues.push({
      type: 'awaiting correction cycle',
      taskId,
      found: status,
      expected: 'Implemented',
      action: `Task ${taskId} previously failed verification and is awaiting the implement-agent's correction — it should re-run and overwrite Validation Status with "Implemented" before this task is verified again.`,
    });
  } else if (!statusReady) {
    issues.push({
      type: 'validation status not Implemented',
      taskId,
      found: status || '(unfilled)',
      expected: 'Implemented',
      action: `Set Validation Status to "Implemented" for ${taskId} once the implement-agent has finished (build + scoped tests passing), or do not attempt to verify it yet.`,
    });
  }

  const declaredFile = fileById[taskId];
  const filesModifiedMatch = execBlock.match(/\*Files Modified:\*([\s\S]*?)(?=\n\s*-\s*\*|\n---|\n- \[|$)/);
  const filesModifiedText = filesModifiedMatch ? filesModifiedMatch[1] : '';
  const fileMatch = !!declaredFile && filesModifiedText.includes(declaredFile);
  if (!declaredFile) {
    issues.push({
      type: 'missing File field in Section 1',
      taskId,
      expected: `a **File:** field on ${taskId}'s Section 1 task block`,
      action: `Add the File field to ${taskId} in Section 1 (this should have been caught by validate-plan.cjs).`,
    });
  } else if (!fileMatch) {
    issues.push({
      type: 'declared File not in Files Modified',
      taskId,
      found: filesModifiedText.trim() || '(empty)',
      expected: declaredFile,
      action: `Add "${declaredFile}" to ${taskId}'s Files Modified list in Section 3, or correct the mismatch.`,
    });
  }

  // Presence of a non-trivial checklist body (Agent Notes / Files Modified / Validation Evidence
  // all have some content) is treated as "the checklist is filled".
  const hasNotes = /\*Agent Notes:\*\s*\S/.test(execBlock);
  const hasEvidence = /\*Validation Evidence:\*\s*\S/.test(execBlock);
  const checklistFilled = hasNotes && hasEvidence && filesModifiedText.trim().length > 0;
  if (!checklistFilled) {
    issues.push({
      type: 'execution checklist not filled',
      taskId,
      expected: 'Agent Notes, Files Modified, and Validation Evidence all populated',
      action: `Fill in the remaining fields of ${taskId}'s execution block before verifying.`,
    });
  }

  results.push({ taskId, checklistFilled, statusReady, fileMatch: fileMatch || !declaredFile });
}

if (issues.length === 0) {
  process.stdout.write(`STATUS: valid\n`);
  process.stdout.write(`FILE: ${planPath}\n`);
  process.stdout.write(`TASKS_CHECKED: ${results.length}\n`);
  for (const r of results) {
    process.stdout.write(`TASK: ${r.taskId} CHECKLIST_FILLED: true VALIDATION_STATUS: Implemented FILE_MATCH: true\n`);
  }
  process.exit(0);
}

process.stderr.write(`STATUS: invalid\n`);
process.stderr.write(`FILE: ${planPath}\n`);
process.stderr.write(`ISSUES: ${issues.length}\n\n`);
issues.forEach((issue, i) => {
  process.stderr.write(`[${i + 1}] ${issue.type}\n`);
  if (issue.taskId)   process.stderr.write(`    Task:     ${issue.taskId}\n`);
  if (issue.found)    process.stderr.write(`    Found:    ${issue.found}\n`);
  if (issue.expected) process.stderr.write(`    Expected: ${issue.expected}\n`);
  if (issue.action)   process.stderr.write(`    Action:   ${issue.action}\n`);
  process.stderr.write('\n');
});
process.exit(1);
