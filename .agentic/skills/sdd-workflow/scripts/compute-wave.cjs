#!/usr/bin/env node
'use strict';

// Reads <feature_path>/PLAN.md and <feature_path>/STATE.md and computes the next
// executable wave. The wave number is never guessed or incremented here — it is
// read from the `**Wave:**` field that sdd-plan's compute-waves.cjs writes
// into every task block at plan time. This script's only job is: find the lowest
// wave number that still has pending tasks, and return the tasks in it (a
// depends_on/completed_tasks check is kept as a deadlock safety net, since a
// correctly-computed wave graph should never leave a pending task's dependency
// unmet once every earlier wave is complete).
//
// Exit codes:
//   0 — executable tasks found, or pipeline complete (stdout: wave groups/table)
//   1 — deadlock detected: pending tasks remain but none are executable (stderr)
//   2 — usage or runtime error (missing argument, file not found, missing wave field)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/compute-wave.cjs <feature_path>

Computes the next executable wave from <feature_path>/PLAN.md (tasks,
depends_on, parallel_group, tier, file, wave) and <feature_path>/STATE.md
(implement.completed_tasks). The next wave is the lowest \`**Wave:**\` value
still present among pending tasks — never an incremented counter — so
\`implement-wave=N\` in sdd-workflow always executes exactly the tasks
sdd-plan assigned to wave N, nothing from an adjacent layer.

Exit codes:
  0   Executable tasks found, or pipeline complete (all tasks done)
  1   Deadlock — pending tasks remain but none are executable
  2   Usage or runtime error (feature_path missing, PLAN.md/STATE.md not found,
      or a pending task has no \`wave:\` field — run sdd-plan's
      compute-waves.cjs on the plan first)

Examples:
  node scripts/compute-wave.cjs docs/auth/login
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , featureArg] = process.argv;

if (!featureArg) {
  process.stderr.write([
    'Error: <feature_path> is required.',
    'Usage: node scripts/compute-wave.cjs <feature_path>',
  ].join('\n') + '\n');
  process.exit(2);
}

const featurePath = featureArg.replace(/\/+$/, '');
const planPath = path.join(featurePath, 'PLAN.md');
const statePath = path.join(featurePath, 'STATE.md');

if (!fs.existsSync(planPath)) {
  process.stderr.write([
    `Error: file not found: "${planPath}"`,
    'Action: run /sdd-plan to produce PLAN.md before computing a wave.',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(statePath)) {
  process.stderr.write([
    `Error: file not found: "${statePath}"`,
    'Action: run write-state.cjs first to initialise STATE.md for this feature.',
  ].join('\n') + '\n');
  process.exit(2);
}

// --- parse PLAN.md tasks (Section 1 only — Section 3 mirrors the same IDs) --

const planContent = fs.readFileSync(planPath, 'utf8');
const section1Match = planContent.match(/## 1\. Plan([\s\S]*?)(?=\n## 3\. Execution|$)/);
const section1 = section1Match ? section1Match[1] : '';

const TASK_RE = /^- \[([ x/y])\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts = [];
let m;
while ((m = TASK_RE.exec(section1)) !== null) {
  starts.push({ mark: m[1], id: m[2], index: m.index });
}

if (starts.length === 0) {
  process.stderr.write([
    `Error: no tasks found in "${planPath}".`,
    'Expected: at least 1 task block matching "- [ ] **TASK-[LAYER]-[CONCERN] — ..."',
    'Action: verify PLAN.md was generated correctly before computing a wave.',
  ].join('\n') + '\n');
  process.exit(2);
}

const blocks = starts.map((s, i) => {
  const end = i + 1 < starts.length ? starts[i + 1].index : section1.length;
  return { ...s, text: section1.slice(s.index, end) };
});

const tasksById = {};
for (const b of blocks) {
  const fileMatch = b.text.match(/\*\*File:\*\*\s*`([^`]*)`/);
  const depsMatch = b.text.match(/`depends_on:`\s*`\[([^\]]*)\]`/);
  const groupMatch = b.text.match(/`parallel_group:`\s*`\[([^\]]*)\]`/);
  const tierMatch = b.text.match(/\|\s*\*\*tier\*\*\s*\|\s*([a-z]+)/i);
  const waveMatch = b.text.match(/`wave:`\s*`?(\d+)`?/);

  tasksById[b.id] = {
    id: b.id,
    mark: b.mark,
    file: fileMatch ? fileMatch[1] : '?',
    depends_on: depsMatch && depsMatch[1].trim()
      ? depsMatch[1].split(',').map(s => s.trim()).filter(Boolean)
      : [],
    parallel_group: groupMatch && groupMatch[1].trim()
      ? groupMatch[1].split(',').map(s => s.trim()).filter(Boolean)[0]
      : 'default',
    tier: tierMatch ? tierMatch[1].toLowerCase() : '?',
    wave: waveMatch ? parseInt(waveMatch[1], 10) : null,
  };
}

// --- parse STATE.md (only what's needed: completed_tasks) --------------------

const stateContent = fs.readFileSync(statePath, 'utf8');

function parseList(content, blockKey, subKey) {
  const blockMatch = content.match(new RegExp(`${blockKey}:\\n([\\s\\S]*?)(?=\\n\\w|\\n$)`));
  const block = blockMatch ? blockMatch[1] : '';
  const lineMatch = block.match(new RegExp(`^[ \\t]*${subKey}:[ \\t]*(.*)$`, 'm'));
  if (!lineMatch) return [];
  const rest = lineMatch[1].trim();
  if (rest === '[]' || rest === '') {
    if (rest === '[]') return [];
    // block list form
    const idx = block.indexOf(lineMatch[0]);
    const after = block.slice(idx + lineMatch[0].length);
    const items = [];
    for (const l of after.split('\n')) {
      if (/^\s*-\s+/.test(l)) items.push(l.trim().slice(2).trim());
      else if (l.trim() === '') continue;
      else break;
    }
    return items;
  }
  return rest.replace(/^\[|\]$/g, '').split(',').map(s => s.trim()).filter(Boolean);
}

const completedTasks = new Set(parseList(stateContent, 'implement', 'completed_tasks'));

// --- compute the next wave from pending tasks' persisted `wave:` field --------

const allIds = blocks.map(b => b.id);
const pendingIds = allIds.filter(id => !completedTasks.has(id));

if (pendingIds.length === 0) {
  process.stdout.write('STATUS: complete\n');
  process.stdout.write('WAVE: done\n');
  process.stdout.write('EXECUTABLE_TASKS: 0\n');
  process.stdout.write('PENDING_TASKS: 0\n');
  process.stdout.write('ACTION: Todas as tasks concluídas. Avançar para fase verify/done.\n');
  process.exit(0);
}

const missingWave = pendingIds.filter(id => tasksById[id].wave === null);
if (missingWave.length > 0) {
  process.stderr.write([
    `Error: ${missingWave.length} pending task(s) have no \`wave:\` field: ${missingWave.join(', ')}`,
    'Action: this PLAN.md predates wave computation, or a task was added by hand after the fact.',
    '  Run sdd-plan\'s compute-waves.cjs on it before implementing:',
    `  node .agentic/skills/sdd-plan/scripts/compute-waves.cjs ${planPath}`,
  ].join('\n') + '\n');
  process.exit(2);
}

const totalWaves = Math.max(...blocks.map(b => tasksById[b.id].wave));
const nextWave = Math.min(...pendingIds.map(id => tasksById[id].wave));

const executable = pendingIds.filter(id => tasksById[id].wave === nextWave && tasksById[id].depends_on.every(d => completedTasks.has(d)));

if (executable.length === 0) {
  process.stderr.write('STATUS: deadlock\n');
  process.stderr.write(`PENDING_TASKS: ${pendingIds.length}\n`);
  process.stderr.write(`WAVE: ${nextWave}\n`);
  process.stderr.write('EXECUTABLE_TASKS: 0\n\n');
  process.stderr.write('Tasks bloqueadas:\n');
  for (const id of pendingIds.filter(id => tasksById[id].wave === nextWave)) {
    const unmet = tasksById[id].depends_on.filter(d => !completedTasks.has(d));
    process.stderr.write(`  ${id} — depende de: ${unmet.join(', ')} (não concluída)\n`);
  }
  process.stderr.write('\nError: a task in the lowest pending wave has an unmet dependency.\n');
  process.stderr.write('Action: review depends_on/wave in PLAN.md — run compute-waves.cjs again if they were hand-edited.\n');
  process.exit(1);
}

const groups = {};
for (const id of executable) {
  const t = tasksById[id];
  if (!groups[t.parallel_group]) groups[t.parallel_group] = [];
  groups[t.parallel_group].push(t);
}
const groupNames = Object.keys(groups);

process.stdout.write('STATUS: ok\n');
process.stdout.write(`WAVE: ${nextWave}\n`);
process.stdout.write(`TOTAL_WAVES: ${totalWaves}\n`);
process.stdout.write(`EXECUTABLE_TASKS: ${executable.length}\n`);
process.stdout.write(`GROUPS: ${groupNames.join(', ')}\n`);
process.stdout.write('---\n');
for (const g of groupNames) {
  for (const t of groups[g]) {
    process.stdout.write(`${g} | ${t.id} | depends_on=[${t.depends_on.join(', ')}] | tier=${t.tier} | file=${t.file}\n`);
  }
}
process.exit(0);
