#!/usr/bin/env node
'use strict';

// Parses PLAN.md and prints the standardized Gate-2 handoff table: one row per
// task with its ID, a <=150-char summary (derived from Reason), which wave it
// belongs to, and that wave's layers/dependencies. This is the single fixed
// format presented at every plan approval — never rebuilt by hand, so it never
// drifts between features or sessions.
//
// Exit codes:
//   0 — index generated (stdout: formatted plain-text table)
//   1 — no tasks found in the plan (stderr: what was expected)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/generate-task-index.cjs <path/to/PLAN.md>

Parses a PLAN.md and prints the standardized handoff table for Gate 2 plan
approval: one row per task with ID, a <=150-char summary (from its Reason
field), its wave number, and that wave's layers/dependencies — read from the
**Wave:** field written by compute-waves.cjs. Run compute-waves.cjs first if
any task is missing a wave value.

Exit codes:
  0   Index generated
  1   No tasks found in the plan
  2   Usage or runtime error (file not found, missing argument, wave missing)

Examples:
  node scripts/generate-task-index.cjs docs/auth/login/PLAN.md
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , planPath] = process.argv;

if (!planPath) {
  process.stderr.write([
    'Error: <path/to/PLAN.md> is required.',
    'Usage: node scripts/generate-task-index.cjs <path/to/PLAN.md>',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(planPath)) {
  process.stderr.write([
    `Error: file not found: "${planPath}"`,
    `Action: run scaffold.cjs to initialise PLAN.md, then write its content first.`,
  ].join('\n') + '\n');
  process.exit(2);
}

const content = fs.readFileSync(planPath, 'utf8');
const section1Match = content.match(/## 1\. Plan([\s\S]*?)(?=\n## 3\. Execution|$)/);
const section1       = section1Match ? section1Match[1] : '';

const TASK_START = /^- \[[ x/]\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts = [];
let m;
while ((m = TASK_START.exec(section1)) !== null) {
  starts.push({ id: m[1], index: m.index });
}

if (starts.length === 0) {
  process.stderr.write([
    `Error: no tasks found in "${planPath}".`,
    `Expected: at least 1 task block matching "- [ ] **TASK-[LAYER]-[CONCERN] — ..."`,
    `Action: write the task blocks in PLAN.md before generating the index.`,
  ].join('\n') + '\n');
  process.exit(1);
}

const blocks = starts.map((s, i) => {
  const end = i + 1 < starts.length ? starts[i + 1].index : section1.length;
  return { id: s.id, text: section1.slice(s.index, end) };
});

function summarize(reason) {
  const collapsed = reason.replace(/\s+/g, ' ').trim();
  if (collapsed.length <= 150) return collapsed;
  return collapsed.slice(0, 149).trimEnd() + '…';
}

const missingWave = [];
const rows = blocks.map(block => {
  const reasonMatch = block.text.match(/\*\*Reason:\*\*\s*([^\n]+)/);
  const depsMatch  = block.text.match(/`depends_on:`\s*`\[([^\]]*)\]`/);
  const groupMatch = block.text.match(/`parallel_group:`\s*`\[([^\]]*)\]`/);
  const waveMatch  = block.text.match(/`wave:`\s*`?(\d+)`?/);
  const tierMatch  = block.text.match(/\|\s*\*\*tier\*\*\s*\|\s*([a-z]+)/i);

  if (!waveMatch) missingWave.push(block.id);

  return {
    id: block.id,
    summary: reasonMatch ? summarize(reasonMatch[1]) : '?',
    depends_on: depsMatch && depsMatch[1].trim()
      ? depsMatch[1].split(',').map(s => s.trim()).filter(Boolean)
      : [],
    parallel_group: groupMatch && groupMatch[1].trim()
      ? groupMatch[1].split(',').map(s => s.trim()).filter(Boolean)[0]
      : 'default',
    wave: waveMatch ? parseInt(waveMatch[1], 10) : null,
    tier: tierMatch ? tierMatch[1].toLowerCase() : '?',
  };
});

if (missingWave.length > 0) {
  process.stderr.write([
    `Error: ${missingWave.length} task(s) have no \`wave:\` value: ${missingWave.join(', ')}`,
    'Action: run compute-waves.cjs on this PLAN.md before generating the handoff table.',
    '  node .agentic/skills/sdd-plan/scripts/compute-waves.cjs ' + planPath,
  ].join('\n') + '\n');
  process.exit(2);
}

const byId = {};
for (const r of rows) byId[r.id] = r;

const byWave = {};
for (const r of rows) {
  if (!byWave[r.wave]) byWave[r.wave] = [];
  byWave[r.wave].push(r);
}
const waveNumbers = Object.keys(byWave).map(Number).sort((a, b) => a - b);

function waveDescription(waveNum) {
  const tasksInWave = byWave[waveNum];
  const layers = [...new Set(tasksInWave.map(t => t.parallel_group))].sort();
  const depWaves = new Set();
  for (const t of tasksInWave) {
    for (const dep of t.depends_on) {
      if (byId[dep]) depWaves.add(byId[dep].wave);
    }
  }
  depWaves.delete(waveNum);
  const depsLabel = depWaves.size ? [...depWaves].sort((a, b) => a - b).map(w => `Wave ${w}`).join(', ') : 'nenhuma';
  return `camadas: ${layers.join(', ')} — depende de: ${depsLabel}`;
}

const waveDescById = {};
for (const w of waveNumbers) waveDescById[w] = waveDescription(w);

const idWidth       = Math.max(...rows.map(r => r.id.length));
const summaryWidth  = Math.max(...rows.map(r => r.summary.length));
const waveDescWidth = Math.max(...waveNumbers.map(w => waveDescById[w].length));

process.stdout.write(`TASK_INDEX: ${rows.length} tasks, ${waveNumbers.length} wave(s)\n`);
process.stdout.write('---\n');
process.stdout.write(
  `${'Task ID'.padEnd(idWidth)} | ${'Resumo'.padEnd(summaryWidth)} | Wave | ${'Onda (camadas / depende de)'.padEnd(waveDescWidth)} | Tier\n`
);
for (const w of waveNumbers) {
  for (const r of byWave[w]) {
    process.stdout.write(
      `${r.id.padEnd(idWidth)} | ${r.summary.padEnd(summaryWidth)} | ${String(r.wave).padEnd(4)} | ${waveDescById[w].padEnd(waveDescWidth)} | ${r.tier}\n`
    );
  }
}
process.exit(0);
