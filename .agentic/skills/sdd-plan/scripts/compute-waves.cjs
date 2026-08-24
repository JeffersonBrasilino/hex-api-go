#!/usr/bin/env node
'use strict';

// Computes and persists a `**Wave:**` field for every task in PLAN.md's Section 1,
// then writes/refreshes a "### Wave Map" table summarising the result. This is the
// single source of truth for wave numbers — sdd-workflow's compute-wave.cjs
// and detect-state.cjs read the persisted field instead of recomputing anything.
//
// Wave(task) = max(layerFloor(parallel_group), 1 + max(Wave(dep) for dep in depends_on))
//   - layerFloor keeps tasks in a later architectural layer out of an earlier wave even
//     when they have no formal depends_on (e.g. an application-layer DTO with
//     depends_on=[] still lands in wave 2, never bundled into wave 1 with domain
//     contracts, because "application" has a higher floor than "domain").
//   - the depends_on term keeps a task strictly after every task it depends on,
//     regardless of layer (e.g. a domain task depending on another domain task).
//   - tasks with no depends_on get exactly layerFloor(parallel_group) as their wave.
//
// Exit codes:
//   0 — waves computed and written (stdout: wave map)
//   1 — circular dependency detected (stderr: cycle path) — nothing is written
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/compute-waves.cjs <path/to/PLAN.md>

Computes a deterministic wave number for every task in PLAN.md's "## 1. Plan"
section, from depends_on (a task always runs strictly after every task it
depends on) and parallel_group (a fixed per-layer floor: domain/config=1,
application=2, infrastructure/module=3, tests=1 — the depends_on term almost
always dominates for tests, since a test task depends on the task it covers).
Writes a **Wave:** N field into each task block (idempotent — updates the
value if already present) and a "### Wave Map" table right after the
Execution Roadmap intro, before the first task block.

Run this after depends_on/parallel_group are set for every task (end of
Phase 2.2) and before validate-plan.cjs. Re-run it whenever depends_on or
parallel_group changes on any task — it is idempotent and safe to re-run.

Exit codes:
  0   Waves computed and written
  1   Circular dependency detected — nothing written, fix depends_on first
  2   Usage or runtime error (file not found, missing argument, no tasks)

Examples:
  node scripts/compute-waves.cjs docs/user/access-control/PLAN.md
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , planPath] = process.argv;

if (!planPath) {
  process.stderr.write([
    'Error: <path/to/PLAN.md> is required.',
    'Usage: node scripts/compute-waves.cjs <path/to/PLAN.md>',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(planPath)) {
  process.stderr.write([
    `Error: file not found: "${planPath}"`,
    'Action: run scaffold.cjs and write the task blocks before computing waves.',
  ].join('\n') + '\n');
  process.exit(2);
}

const content = fs.readFileSync(planPath, 'utf8');
const section1Match = content.match(/## 1\. Plan([\s\S]*?)(?=\n## 3\. Execution|$)/);
if (!section1Match) {
  process.stderr.write([
    'Error: "## 1. Plan" section not found.',
    'Action: verify PLAN.md follows plan-schema.md before computing waves.',
  ].join('\n') + '\n');
  process.exit(2);
}
const section1 = section1Match[1];
const section1Start = section1Match.index + section1Match[0].indexOf(section1);

const TASK_RE = /^- \[([ x/])\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts = [];
let m;
while ((m = TASK_RE.exec(section1)) !== null) {
  starts.push({ mark: m[1], id: m[2], index: m.index });
}

if (starts.length === 0) {
  process.stderr.write([
    `Error: no tasks found in "${planPath}".`,
    'Expected: at least 1 task block matching "- [ ] **TASK-[LAYER]-[CONCERN] — ..."',
    'Action: write the task blocks in PLAN.md before computing waves.',
  ].join('\n') + '\n');
  process.exit(2);
}

const blocks = starts.map((s, i) => {
  const end = i + 1 < starts.length ? starts[i + 1].index : section1.length;
  return { ...s, text: section1.slice(s.index, end), start: s.index, end: i + 1 < starts.length ? starts[i + 1].index : section1.length };
});

const LAYER_FLOOR = {
  domain: 1,
  config: 1,
  tests: 1,
  application: 2,
  infrastructure: 3,
  module: 3,
  default: 1,
};

const tasksById = {};
for (const b of blocks) {
  const depsMatch = b.text.match(/`depends_on:`\s*`\[([^\]]*)\]`/);
  const groupMatch = b.text.match(/`parallel_group:`\s*`\[([^\]]*)\]`/);
  tasksById[b.id] = {
    id: b.id,
    depends_on: depsMatch && depsMatch[1].trim()
      ? depsMatch[1].split(',').map(s => s.trim()).filter(Boolean)
      : [],
    parallel_group: groupMatch && groupMatch[1].trim()
      ? groupMatch[1].split(',').map(s => s.trim()).filter(Boolean)[0]
      : 'default',
  };
}

// Unresolved depends_on references are not this script's concern — validate-plan.cjs
// is the authority for that. Here, treat an unresolved reference as "no dependency"
// so wave computation can still proceed; validate-plan.cjs will reject the plan anyway.

const waveCache = {};
let cyclePath = null;

function waveOf(id, stack) {
  if (waveCache[id] !== undefined) return waveCache[id];
  if (stack.includes(id)) {
    cyclePath = [...stack, id].join(' -> ');
    return 0;
  }
  const task = tasksById[id];
  if (!task) return 0; // unresolved reference — let validate-plan.cjs report it
  const floor = LAYER_FLOOR[task.parallel_group] || LAYER_FLOOR.default;
  let wave = floor;
  for (const dep of task.depends_on) {
    if (!tasksById[dep]) continue;
    const depWave = waveOf(dep, [...stack, id]);
    if (cyclePath) return 0;
    wave = Math.max(wave, depWave + 1);
  }
  waveCache[id] = wave;
  return wave;
}

for (const id of Object.keys(tasksById)) {
  waveOf(id, []);
  if (cyclePath) break;
}

if (cyclePath) {
  process.stderr.write('STATUS: cycle\n');
  process.stderr.write(`FILE: ${planPath}\n\n`);
  process.stderr.write(`Circular dependency: ${cyclePath}\n\n`);
  process.stderr.write('Error: cannot compute waves while a dependency cycle exists.\n');
  process.stderr.write('Action: break the cycle in depends_on (validate-plan.cjs will confirm), then re-run.\n');
  process.exit(1);
}

// --- write **Wave:** N into each task block (insert after parallel_group line, or
//     replace the existing value if already present) ------------------------------

let updatedSection1 = section1;
// Work backwards so earlier offsets in `blocks` stay valid as we splice.
for (let i = blocks.length - 1; i >= 0; i--) {
  const b = blocks[i];
  const wave = waveCache[b.id];
  const blockText = updatedSection1.slice(b.start, b.end);
  let newBlockText;
  // Matches both an already-computed numeric wave and the unfilled `[N]` placeholder,
  // so the first run replaces the placeholder in place instead of appending a duplicate line.
  const waveLineRe = /\n(\s*)- `wave:`[^\n]*\n/;
  if (waveLineRe.test(blockText)) {
    newBlockText = blockText.replace(waveLineRe, (full, indent) => `\n${indent}- \`wave:\` \`${wave}\`\n`);
  } else {
    const groupLineRe = /(\n(\s*)- `parallel_group:`[^\n]*\n)/;
    if (!groupLineRe.test(blockText)) {
      process.stderr.write([
        `Error: task ${b.id} has no \`parallel_group:\` line to anchor the wave field to.`,
        'Action: fix the Dependencies block to match plan-schema.md, then re-run.',
      ].join('\n') + '\n');
      process.exit(2);
    }
    newBlockText = blockText.replace(groupLineRe, (full, line, indent) => `${line}${indent}- \`wave:\` \`${wave}\`\n`);
  }
  updatedSection1 = updatedSection1.slice(0, b.start) + newBlockText + updatedSection1.slice(b.end);
}

// --- build the Wave Map table --------------------------------------------------

const byWave = {};
for (const id of Object.keys(tasksById)) {
  const w = waveCache[id];
  if (!byWave[w]) byWave[w] = [];
  byWave[w].push(id);
}
const waveNumbers = Object.keys(byWave).map(Number).sort((a, b) => a - b);

function waveDeps(waveNum) {
  const deps = new Set();
  for (const id of byWave[waveNum]) {
    for (const dep of tasksById[id].depends_on) {
      if (tasksById[dep] !== undefined) deps.add(waveCache[dep]);
    }
  }
  deps.delete(waveNum);
  return [...deps].sort((a, b) => a - b);
}

function waveLayers(waveNum) {
  return [...new Set(byWave[waveNum].map(id => tasksById[id].parallel_group))].sort();
}

const waveMapRows = waveNumbers.map(w => {
  const deps = waveDeps(w);
  const depsLabel = deps.length ? deps.map(d => `Wave ${d}`).join(', ') : '—';
  return `| ${w} | ${waveLayers(w).join(', ')} | ${depsLabel} | ${byWave[w].join(', ')} |`;
});

// No trailing newline here on purpose — the block's own closing separator
// ("\n\n---\n\n") is appended explicitly at each insertion/replacement site below,
// so this string never has to double as a shared separator with the next element
// (that ambiguity is exactly what previously let the WAVE_MAP_RE lookahead skip
// past an entire task block on a second run, deleting it).
const waveMapBlock = [
  '### Wave Map',
  '',
  '| Wave | Camadas (parallel_group) | Depende de | Tasks |',
  '|------|---------------------------|------------|-------|',
  ...waveMapRows,
].join('\n');

// Insert/replace the Wave Map right after the "### Execution Roadmap" blockquote,
// before the first task block, so it always sits just above the task list it
// summarises. The lookahead requires the Wave Map's own "---" separator to be
// immediately followed by a task heading ("- [") — never a generic "---\n" — so a
// re-run can never mistake a later task-to-task separator for this block's end and
// swallow an entire task block into the replacement.
const WAVE_MAP_RE = /### Wave Map\n[\s\S]*?\n\n---\n\n(?=- \[)/;
let finalSection1;
if (WAVE_MAP_RE.test(updatedSection1)) {
  finalSection1 = updatedSection1.replace(WAVE_MAP_RE, waveMapBlock + '\n\n---\n\n');
} else {
  // The trailing bare `\n` (outside the capture group) consumes the blank line that
  // separates the Roadmap's "---" from the first task heading, so no leftover blank
  // line survives to throw off WAVE_MAP_RE's exact-boundary lookahead on the next run.
  const anchorRe = /(### Execution Roadmap\n(?:.*\n)*?\n---\n)\n/;
  if (!anchorRe.test(updatedSection1)) {
    process.stderr.write([
      'Error: could not find "### Execution Roadmap" anchor to insert the Wave Map.',
      'Action: verify PLAN.md follows plan-schema.md structure, then re-run.',
    ].join('\n') + '\n');
    process.exit(2);
  }
  finalSection1 = updatedSection1.replace(anchorRe, `$1\n${waveMapBlock}\n\n---\n\n`);
}

const newContent = content.slice(0, section1Start) + finalSection1 + content.slice(section1Start + section1.length);
fs.writeFileSync(planPath, newContent, 'utf8');

// --- report --------------------------------------------------------------------

process.stdout.write('STATUS: ok\n');
process.stdout.write(`FILE: ${planPath}\n`);
process.stdout.write(`WAVES: ${waveNumbers.length}\n`);
process.stdout.write('---\n');
for (const w of waveNumbers) {
  const deps = waveDeps(w);
  const depsLabel = deps.length ? deps.map(d => `Wave ${d}`).join(', ') : 'nenhuma';
  process.stdout.write(`WAVE ${w} | camadas: ${waveLayers(w).join(', ')} | depende de: ${depsLabel} | ${byWave[w].length} task(s)\n`);
  for (const id of byWave[w]) {
    process.stdout.write(`  ${id}\n`);
  }
}
process.exit(0);
