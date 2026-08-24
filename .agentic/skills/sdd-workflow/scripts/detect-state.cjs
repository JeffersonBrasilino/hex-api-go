#!/usr/bin/env node
'use strict';

// Detects the true current SDD phase for a feature directly from the artifacts on
// disk (PRD.md or PRD.cache.md, NOTES.md, PLAN.md — its Status header and task
// checkboxes), instead of trusting a stale or absent STATE.md. This lets
// /sdd-workflow jump straight to the furthest phase actually reached in a single
// invocation, rather than requiring one re-invocation per phase when PRD/PLAN were
// already produced out-of-band.
//
// PRD.cache.md is sdd-plan's materialized, temporary local copy of a
// card-sourced PRD (see sdd-plan/SKILL.md "PRD input: local file vs. card").
// It is recognized here on equal footing with PRD.md so a card-sourced feature
// advances past PHASE: prd once sdd-plan has run — detection must not depend
// on STATE.md to know this, per this script's own "detect from disk" contract. If
// both exist, PRD.md (the authored, permanent artifact) wins.
//
// CURRENT_WAVE is read directly from each task's `**Wave:**` field (written by
// sdd-plan's compute-waves.cjs at plan time) — never recomputed from
// depends_on here, so this script and compute-wave.cjs always agree on what "wave
// N" means for a given plan.
//
// Exit codes:
//   0 — detection succeeded (stdout: structured KEY: value report)
//   2 — usage or runtime error (stderr: what failed + usage hint, including a
//       pending task with no `wave:` field — run compute-waves.cjs first)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/detect-state.cjs <feature_path>

Inspects <feature_path>/{PRD.md or PRD.cache.md,NOTES.md,PLAN.md} and reports the
furthest SDD phase reachable from what already exists on disk:

  - No PRD.md and no PRD.cache.md      -> PHASE: prd
  - PRD (either file), no PLAN.md      -> PHASE: plan
  - PLAN.md, Status Draft/Planning     -> PHASE: plan (awaiting dev approval)
  - PLAN.md, Status Ready for Impl.    -> PHASE: implement (wave 1, not yet started)
  - PLAN.md, Status Implementando      -> PHASE: implement or verify, with
                                           COMPLETED_TASKS/PENDING_TASKS/CURRENT_WAVE
                                           parsed from Section 1 task checkboxes
  - PLAN.md, Status Done               -> PHASE: done

Exit codes:
  0   Detection succeeded
  2   Usage or runtime error (feature_path missing or not found)

Examples:
  node scripts/detect-state.cjs docs/users/access-control-with-scripts
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [, , featureArg] = process.argv;

if (!featureArg) {
  process.stderr.write([
    'Error: <feature_path> is required.',
    'Usage: node scripts/detect-state.cjs <feature_path>',
  ].join('\n') + '\n');
  process.exit(2);
}

const featurePath = featureArg.replace(/\/+$/, '');

if (!fs.existsSync(featurePath) || !fs.statSync(featurePath).isDirectory()) {
  process.stderr.write([
    `Error: directory not found: "${featurePath}"`,
    'Action: confirm the feature path or bootstrap it via /sdd-workflow first.',
  ].join('\n') + '\n');
  process.exit(2);
}

const prdLocalPath = path.join(featurePath, 'PRD.md');
const prdCachePath = path.join(featurePath, 'PRD.cache.md');
const notesPath = path.join(featurePath, 'NOTES.md');
const planPath = path.join(featurePath, 'PLAN.md');

const hasPrdLocal = fs.existsSync(prdLocalPath);
const hasPrdCache = fs.existsSync(prdCachePath);
const hasPrd = hasPrdLocal || hasPrdCache;
// PRD.md (authored, permanent) takes precedence over PRD.cache.md (derived,
// regenerable) when both exist.
const prdPath = hasPrdLocal ? prdLocalPath : prdCachePath;
const hasNotes = fs.existsSync(notesPath);
const hasPlan = fs.existsSync(planPath);

function emit(lines) {
  process.stdout.write(lines.join('\n') + '\n');
}

if (!hasPrd) {
  emit([
    'PHASE: prd',
    'PRD: (none)',
    `NOTES: ${hasNotes ? notesPath : '(none)'}`,
    'PLAN: (none)',
    'ACTION: Aguardando PRD. Informe ao usuário o caminho local ou o link do card e execute /sdd-plan com essa referência.',
  ]);
  process.exit(0);
}

if (!hasPlan) {
  emit([
    'PHASE: plan',
    `PRD: ${prdPath}`,
    `NOTES: ${hasNotes ? notesPath : '(none)'}`,
    'PLAN: (none)',
    `ACTION: Aguardando PLAN. Execute /sdd-plan ${prdPath} em uma nova sessão.`,
  ]);
  process.exit(0);
}

const planContent = fs.readFileSync(planPath, 'utf8');
const statusMatch = planContent.match(/\*\*Status:\*\*\s*`([^`]+)`/);
const status = statusMatch ? statusMatch[1] : 'Unknown';

const baseFields = [
  `PRD: ${prdPath}`,
  `NOTES: ${hasNotes ? notesPath : '(none)'}`,
  `PLAN: ${planPath}`,
  `PLAN_STATUS: ${status}`,
];

if (status === 'Draft' || status === 'Planning' || status === 'Unknown') {
  emit([
    'PHASE: plan',
    ...baseFields,
    'ACTION: PLAN.md existe mas ainda não foi aprovado pelo dev. Permanecer na fase plan até o Status mudar para "Ready for Implementation".',
  ]);
  process.exit(0);
}

// Status is "Ready for Implementation", "Implementando" or "Done" — parse tasks
// from Section 1 only (Section 3 mirrors the same IDs during execution and would
// double-count).
const section1Match = planContent.match(/## 1\. Plan([\s\S]*?)(?=\n## 3\. Execution|$)/);
const section1 = section1Match ? section1Match[1] : '';

const TASK_RE = /^- \[([ x/y])\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts = [];
let m;
while ((m = TASK_RE.exec(section1)) !== null) {
  starts.push({ mark: m[1], id: m[2], index: m.index });
}

if (starts.length === 0) {
  emit([
    'PHASE: plan',
    ...baseFields,
    'ACTION: PLAN.md não contém tasks reconhecíveis na Seção 1. Verifique o formato do plano antes de avançar.',
  ]);
  process.exit(0);
}

const blocks = starts.map((s, i) => {
  const end = i + 1 < starts.length ? starts[i + 1].index : section1.length;
  return { ...s, text: section1.slice(s.index, end) };
});

const waveById = {};
for (const b of blocks) {
  const waveMatch = b.text.match(/`wave:`\s*`?(\d+)`?/);
  waveById[b.id] = waveMatch ? parseInt(waveMatch[1], 10) : null;
}

const completed = blocks.filter(b => b.mark === 'x' || b.mark === 'y').map(b => b.id);
const pending = blocks.filter(b => b.mark !== 'x' && b.mark !== 'y').map(b => b.id);

const missingWave = blocks.filter(b => waveById[b.id] === null).map(b => b.id);
if (missingWave.length > 0) {
  process.stderr.write([
    `Error: ${missingWave.length} task(s) in "${planPath}" have no \`wave:\` field: ${missingWave.join(', ')}`,
    'Action: this PLAN.md predates wave computation — run sdd-plan\'s compute-waves.cjs on it:',
    `  node .agentic/skills/sdd-plan/scripts/compute-waves.cjs ${planPath}`,
  ].join('\n') + '\n');
  process.exit(2);
}

const taskFields = [
  `TOTAL_TASKS: ${blocks.length}`,
  `COMPLETED_TASKS: ${completed.join(',')}`,
  `PENDING_TASKS: ${pending.join(',')}`,
];

if (status === 'Done' || (status === 'Implementando' && pending.length === 0)) {
  const phase = status === 'Done' ? 'done' : 'verify';
  emit([
    `PHASE: ${phase}`,
    ...baseFields,
    ...taskFields,
  ]);
  process.exit(0);
}

const currentWave = pending.length === 0
  ? Math.max(...blocks.map(b => waveById[b.id])) + 1
  : Math.min(...pending.map(id => waveById[id]));

emit([
  'PHASE: implement',
  ...baseFields,
  ...taskFields,
  `CURRENT_WAVE: ${currentWave}`,
  status === 'Ready for Implementation'
    ? 'ACTION: Plano aprovado. Atualizar Status do PLAN.md para "Implementando" e iniciar Wave 1.'
    : 'ACTION: Continuar execução a partir das waves pendentes.',
]);
process.exit(0);
