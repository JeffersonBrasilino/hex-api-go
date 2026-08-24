#!/usr/bin/env node
'use strict';

// Reads/creates <feature_path>/STATE.md and applies only the fields explicitly
// requested via flags, following the fixed schema in references/state-schema.md.
// Never rewrites the whole file from scratch when it already exists, and never
// touches verify.retry_counts / verify.failed_tasks / implement.current_wave_groups
// unless a flag explicitly targets them.
//
// Exit codes:
//   0 — STATE.md written/updated (stdout: which fields changed vs. were preserved)
//   1 — logic violation (invalid --phase, non-numeric --current-wave, etc.)
//   2 — usage or runtime error (missing argument, feature directory not found)

const fs = require('fs');
const path = require('path');

const HELP = `
Usage: node scripts/write-state.cjs <feature_path> [flags]

Creates or updates <feature_path>/STATE.md, touching only the fields named by
the flags passed. Never overwrites the whole file, and never clears
retry_counts / failed_tasks / current_wave_groups unless explicitly told to.

Flags:
  --phase <value>            prd | plan | implement | verify | done
  --prd <path>                Path to PRD.md
  --notes <path>               Path to NOTES.md
  --plan <path>                 Path to PLAN.md
  --current-wave <N>            Current wave number
  --completed <IDs>              Replace the whole completed_tasks list (comma-separated)
  --pending <IDs>                 Replace the whole pending_tasks list (comma-separated)
  --add-completed <ID>             Append one ID to completed_tasks (idempotent)
  --remove-pending <ID>              Remove one ID from pending_tasks
  --add-pending <ID>                   Append one ID to pending_tasks (idempotent)
  --remove-completed <ID>                Remove one ID from completed_tasks (task reopening)
  --add-failed <ID:reason>             Append one entry to failed_tasks
  --increment-retry <ID>                 Increment retry_counts[ID] by 1
  --clear-wave-groups                      Clear current_wave_groups
  --wave-groups <groups>                     Replace current_wave_groups (comma-separated)
  --prd-gate-wave <N>                        Wave number the PRD gate status below refers to
  --prd-gate-status <value>                    pending | passed | failed | skipped
  --pr-status <value>                          not_started | blocked_missing_tool |
                                                blocked_protected_branch | created | failed
  --pr-url <url>                               PR/MR URL once created
  --pr-branch <name>                           Branch used for the PR/MR

Exit codes:
  0   STATE.md written/updated
  1   Logic violation — invalid flag value (stderr: found/expected/action)
  2   Usage or runtime error (missing argument, feature directory not found)

Examples:
  node scripts/write-state.cjs docs/auth/login --phase plan --prd docs/auth/login/PRD.md
  node scripts/write-state.cjs docs/auth/login --add-completed TASK-DOM-A --remove-pending TASK-DOM-A
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const VALID_PHASES = ['prd', 'plan', 'implement', 'verify', 'done'];
const VALID_PRD_GATE_STATUSES = ['pending', 'passed', 'failed', 'skipped'];
const VALID_PR_STATUSES = ['not_started', 'blocked_missing_tool', 'blocked_protected_branch', 'created', 'failed'];

function fail1(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(1);
}

function fail2(lines) {
  process.stderr.write(lines.join('\n') + '\n');
  process.exit(2);
}

// --- argv parsing -----------------------------------------------------------

const argv = process.argv.slice(2);
const featurePath = argv[0] && !argv[0].startsWith('--') ? argv[0].replace(/\/+$/, '') : null;

if (!featurePath) {
  fail2([
    'Error: <feature_path> is required.',
    'Usage: node scripts/write-state.cjs <feature_path> [flags]',
  ]);
}

const flags = {};
for (let i = 1; i < argv.length; i++) {
  const tok = argv[i];
  if (!tok.startsWith('--')) continue;
  const name = tok.slice(2);
  const BOOLEAN_FLAGS = new Set(['clear-wave-groups']);
  if (BOOLEAN_FLAGS.has(name)) {
    flags[name] = true;
    continue;
  }
  const value = argv[i + 1];
  if (value === undefined || value.startsWith('--')) {
    fail2([`Error: flag --${name} requires a value.`, 'Usage: node scripts/write-state.cjs <feature_path> [flags]']);
  }
  if (!flags[name]) flags[name] = [];
  flags[name].push(value);
  i++;
}

if (!fs.existsSync(featurePath) || !fs.statSync(featurePath).isDirectory()) {
  fail2([
    `Error: feature directory not found: "${featurePath}"`,
    'Action: confirm the path exists or run scaffold.cjs first.',
  ]);
}

const statePath = path.join(featurePath, 'STATE.md');

// --- fixed-schema parse/serialize -------------------------------------------

function stripQuotes(s) {
  const t = (s || '').trim();
  const m = t.match(/^"(.*)"$/);
  return m ? m[1] : t;
}

function indentOf(line) {
  return line.match(/^(\s*)/)[1].length;
}

function defaultState(feature) {
  return {
    feature,
    phase: 'prd',
    artifacts: { prd: '', notes: '', plan: '' },
    implement: { current_wave: 1, completed_tasks: [], pending_tasks: [], current_wave_groups: [] },
    verify: { retry_counts: {}, failed_tasks: [], prd_gate: { wave: null, status: 'pending' } },
    pr: { status: 'not_started', url: '', branch: '' },
  };
}

function parseState(content, feature) {
  const state = defaultState(feature);
  const lines = content.split('\n');
  let i = 0;

  while (i < lines.length) {
    const line = lines[i].replace(/\r$/, '');
    if (line.trim() === '' || line.trim().startsWith('#')) { i++; continue; }
    const indent = indentOf(line);
    if (indent !== 0) { i++; continue; }

    const m = line.trim().match(/^([a-zA-Z_]+):\s*(.*)$/);
    if (!m) { i++; continue; }
    const key = m[1];
    const val = m[2];

    if (key === 'feature') { state.feature = stripQuotes(val); i++; continue; }
    if (key === 'phase') { state.phase = stripQuotes(val); i++; continue; }

    if (key === 'artifacts') {
      i++;
      while (i < lines.length && (lines[i].trim() === '' || indentOf(lines[i]) > indent)) {
        const lt = lines[i].trim();
        const mm = lt.match(/^(prd|notes|plan):\s*(.*)$/);
        if (mm) state.artifacts[mm[1]] = stripQuotes(mm[2]);
        i++;
      }
      continue;
    }

    if (key === 'pr') {
      i++;
      while (i < lines.length && (lines[i].trim() === '' || indentOf(lines[i]) > indent)) {
        const lt = lines[i].trim();
        const mm = lt.match(/^(status|url|branch):\s*(.*)$/);
        if (mm) state.pr[mm[1]] = stripQuotes(mm[2]);
        i++;
      }
      continue;
    }

    if (key === 'implement') {
      i++;
      while (i < lines.length && (lines[i].trim() === '' || indentOf(lines[i]) > indent)) {
        const l = lines[i];
        const lt = l.trim();
        if (lt === '') { i++; continue; }
        const lIndent = indentOf(l);
        const mm = lt.match(/^(current_wave|completed_tasks|pending_tasks|current_wave_groups):\s*(.*)$/);
        if (!mm) { i++; continue; }
        const subKey = mm[1];
        const rest = mm[2];
        if (subKey === 'current_wave') {
          state.implement.current_wave = parseInt(rest, 10) || 1;
          i++;
          continue;
        }
        if (rest.trim() === '[]') {
          state.implement[subKey] = [];
          i++;
          continue;
        }
        if (rest.trim() !== '') {
          const inner = rest.trim().replace(/^\[|\]$/g, '');
          state.implement[subKey] = inner ? inner.split(',').map(s => s.trim()).filter(Boolean) : [];
          i++;
          continue;
        }
        i++;
        const items = [];
        while (i < lines.length && indentOf(lines[i]) > lIndent && lines[i].trim().startsWith('- ')) {
          items.push(lines[i].trim().slice(2).trim());
          i++;
        }
        state.implement[subKey] = items;
      }
      continue;
    }

    if (key === 'verify') {
      i++;
      while (i < lines.length && (lines[i].trim() === '' || indentOf(lines[i]) > indent)) {
        const l = lines[i];
        const lt = l.trim();
        if (lt === '') { i++; continue; }
        const lIndent = indentOf(l);
        const mm = lt.match(/^(retry_counts|failed_tasks|prd_gate):\s*(.*)$/);
        if (!mm) { i++; continue; }
        const subKey = mm[1];
        const rest = mm[2];

        if (subKey === 'prd_gate') {
          i++;
          const gate = { wave: null, status: 'pending' };
          while (i < lines.length && indentOf(lines[i]) > lIndent && lines[i].trim() !== '') {
            const gm = lines[i].trim().match(/^(wave|status):\s*(.*)$/);
            if (gm) {
              if (gm[1] === 'wave') {
                const n = parseInt(gm[2], 10);
                gate.wave = Number.isInteger(n) ? n : null;
              } else {
                gate.status = stripQuotes(gm[2]);
              }
            }
            i++;
          }
          state.verify.prd_gate = gate;
          continue;
        }

        if (subKey === 'retry_counts') {
          if (rest.trim() === '{}') { state.verify.retry_counts = {}; i++; continue; }
          i++;
          const map = {};
          while (i < lines.length && indentOf(lines[i]) > lIndent && lines[i].trim() !== '') {
            const km = lines[i].trim().match(/^([A-Za-z0-9_-]+):\s*(\d+)/);
            if (km) map[km[1]] = parseInt(km[2], 10);
            i++;
          }
          state.verify.retry_counts = map;
          continue;
        }

        if (subKey === 'failed_tasks') {
          if (rest.trim() === '[]') { state.verify.failed_tasks = []; i++; continue; }
          i++;
          const items = [];
          while (i < lines.length && indentOf(lines[i]) > lIndent && lines[i].trim().startsWith('- ')) {
            const idMatch = lines[i].trim().match(/^- id:\s*(.*)$/);
            const id = idMatch ? stripQuotes(idMatch[1]) : '';
            i++;
            let reason = '';
            if (i < lines.length && lines[i].trim().startsWith('reason:')) {
              reason = stripQuotes(lines[i].trim().replace(/^reason:\s*/, ''));
              i++;
            }
            items.push({ id, reason });
          }
          state.verify.failed_tasks = items;
          continue;
        }
      }
      continue;
    }

    i++;
  }

  return state;
}

function renderList(key, arr, indent) {
  if (!arr || arr.length === 0) return `${indent}${key}: []`;
  return [`${indent}${key}:`, ...arr.map(item => `${indent}  - ${item}`)].join('\n');
}

function renderRetryCounts(map, indent) {
  const keys = Object.keys(map || {});
  if (keys.length === 0) return `${indent}retry_counts: {}`;
  return [`${indent}retry_counts:`, ...keys.map(k => `${indent}  ${k}: ${map[k]}`)].join('\n');
}

function renderFailedTasks(arr, indent) {
  if (!arr || arr.length === 0) return `${indent}failed_tasks: []`;
  const out = [`${indent}failed_tasks:`];
  for (const item of arr) {
    out.push(`${indent}  - id: ${item.id}`);
    out.push(`${indent}    reason: "${item.reason}"`);
  }
  return out.join('\n');
}

function renderPrdGate(gate, indent) {
  const g = gate || { wave: null, status: 'pending' };
  return [
    `${indent}prd_gate:`,
    `${indent}  wave: ${g.wave === null || g.wave === undefined ? 'null' : g.wave}`,
    `${indent}  status: ${g.status || 'pending'}`,
  ].join('\n');
}

function serialize(state) {
  return [
    '# SDD State',
    `feature: ${state.feature}`,
    `phase: ${state.phase}`,
    '',
    'artifacts:',
    `  prd:   "${state.artifacts.prd}"`,
    `  notes: "${state.artifacts.notes}"`,
    `  plan:  "${state.artifacts.plan}"`,
    '',
    'implement:',
    `  current_wave: ${state.implement.current_wave}`,
    renderList('completed_tasks', state.implement.completed_tasks, '  '),
    renderList('pending_tasks', state.implement.pending_tasks, '  '),
    renderList('current_wave_groups', state.implement.current_wave_groups, '  '),
    '',
    'verify:',
    renderRetryCounts(state.verify.retry_counts, '  '),
    renderFailedTasks(state.verify.failed_tasks, '  '),
    renderPrdGate(state.verify.prd_gate, '  '),
    '',
    'pr:',
    `  status: ${state.pr.status}`,
    `  url:    "${state.pr.url}"`,
    `  branch: "${state.pr.branch}"`,
    '',
  ].join('\n');
}

// --- load or initialise ------------------------------------------------------

const defaultFeature = featurePath.replace(/^docs\//, '');
const state = fs.existsSync(statePath)
  ? parseState(fs.readFileSync(statePath, 'utf8'), defaultFeature)
  : defaultState(defaultFeature);

const updated = [];
const PROTECTED = ['retry_counts', 'failed_tasks', 'current_wave_groups', 'prd_gate'];

// --- apply flags --------------------------------------------------------------

if (flags['phase']) {
  const value = flags['phase'][flags['phase'].length - 1];
  if (!VALID_PHASES.includes(value)) {
    fail1([
      `Error: --phase "${value}" is not a valid SDD phase.`,
      `Received: "${value}"`,
      `Expected: one of ${VALID_PHASES.join(' | ')}`,
      'Action: use a valid phase value.',
    ]);
  }
  state.phase = value;
  updated.push('phase');
}

if (flags['prd']) {
  state.artifacts.prd = flags['prd'][flags['prd'].length - 1];
  updated.push('prd');
}
if (flags['notes']) {
  state.artifacts.notes = flags['notes'][flags['notes'].length - 1];
  updated.push('notes');
}
if (flags['plan']) {
  state.artifacts.plan = flags['plan'][flags['plan'].length - 1];
  updated.push('plan');
}

if (flags['current-wave']) {
  const raw = flags['current-wave'][flags['current-wave'].length - 1];
  const n = parseInt(raw, 10);
  if (!Number.isInteger(n) || String(n) !== raw.trim() || n < 1) {
    fail1([
      `Error: --current-wave "${raw}" is not a valid wave number.`,
      `Received: "${raw}"`,
      'Expected: a positive integer.',
      'Action: pass a positive integer wave number.',
    ]);
  }
  state.implement.current_wave = n;
  updated.push('current_wave');
}

if (flags['completed']) {
  const raw = flags['completed'][flags['completed'].length - 1];
  state.implement.completed_tasks = raw.trim() ? raw.split(',').map(s => s.trim()).filter(Boolean) : [];
  updated.push('completed_tasks');
}

if (flags['pending']) {
  const raw = flags['pending'][flags['pending'].length - 1];
  state.implement.pending_tasks = raw.trim() ? raw.split(',').map(s => s.trim()).filter(Boolean) : [];
  updated.push('pending_tasks');
}

if (flags['add-completed']) {
  for (const id of flags['add-completed']) {
    if (!state.implement.completed_tasks.includes(id)) state.implement.completed_tasks.push(id);
  }
  updated.push('completed_tasks');
}

if (flags['remove-pending']) {
  for (const id of flags['remove-pending']) {
    state.implement.pending_tasks = state.implement.pending_tasks.filter(t => t !== id);
  }
  updated.push('pending_tasks');
}

if (flags['add-pending']) {
  for (const id of flags['add-pending']) {
    if (!state.implement.pending_tasks.includes(id)) state.implement.pending_tasks.push(id);
  }
  updated.push('pending_tasks');
}

if (flags['remove-completed']) {
  for (const id of flags['remove-completed']) {
    state.implement.completed_tasks = state.implement.completed_tasks.filter(t => t !== id);
  }
  updated.push('completed_tasks');
}

if (flags['add-failed']) {
  for (const raw of flags['add-failed']) {
    const idx = raw.indexOf(':');
    if (idx === -1) {
      fail1([
        `Error: --add-failed "${raw}" is not in "ID:reason" format.`,
        `Received: "${raw}"`,
        'Expected: "TASK-ID:reason text"',
        'Action: pass the task ID and reason separated by a colon.',
      ]);
    }
    const id = raw.slice(0, idx).trim();
    const reason = raw.slice(idx + 1).trim();
    state.verify.failed_tasks = state.verify.failed_tasks.filter(t => t.id !== id);
    state.verify.failed_tasks.push({ id, reason });
  }
  updated.push('failed_tasks');
}

if (flags['increment-retry']) {
  for (const id of flags['increment-retry']) {
    state.verify.retry_counts[id] = (state.verify.retry_counts[id] || 0) + 1;
  }
  updated.push('retry_counts');
}

if (flags['clear-wave-groups']) {
  state.implement.current_wave_groups = [];
  updated.push('current_wave_groups');
}

if (flags['wave-groups']) {
  const raw = flags['wave-groups'][flags['wave-groups'].length - 1];
  state.implement.current_wave_groups = raw.trim() ? raw.split(',').map(s => s.trim()).filter(Boolean) : [];
  updated.push('current_wave_groups');
}

if (flags['prd-gate-wave'] || flags['prd-gate-status']) {
  if (flags['prd-gate-status']) {
    const value = flags['prd-gate-status'][flags['prd-gate-status'].length - 1];
    if (!VALID_PRD_GATE_STATUSES.includes(value)) {
      fail1([
        `Error: --prd-gate-status "${value}" is not valid.`,
        `Received: "${value}"`,
        `Expected: one of ${VALID_PRD_GATE_STATUSES.join(' | ')}`,
        'Action: use a valid PRD gate status value.',
      ]);
    }
    state.verify.prd_gate.status = value;
  }
  if (flags['prd-gate-wave']) {
    const raw = flags['prd-gate-wave'][flags['prd-gate-wave'].length - 1];
    const n = parseInt(raw, 10);
    if (!Number.isInteger(n) || String(n) !== raw.trim() || n < 1) {
      fail1([
        `Error: --prd-gate-wave "${raw}" is not a valid wave number.`,
        `Received: "${raw}"`,
        'Expected: a positive integer.',
        'Action: pass a positive integer wave number.',
      ]);
    }
    state.verify.prd_gate.wave = n;
  }
  updated.push('prd_gate');
}

if (flags['pr-status']) {
  const value = flags['pr-status'][flags['pr-status'].length - 1];
  if (!VALID_PR_STATUSES.includes(value)) {
    fail1([
      `Error: --pr-status "${value}" is not valid.`,
      `Received: "${value}"`,
      `Expected: one of ${VALID_PR_STATUSES.join(' | ')}`,
      'Action: use a valid PR status value.',
    ]);
  }
  state.pr.status = value;
  updated.push('pr_status');
}
if (flags['pr-url']) {
  state.pr.url = flags['pr-url'][flags['pr-url'].length - 1];
  updated.push('pr_url');
}
if (flags['pr-branch']) {
  state.pr.branch = flags['pr-branch'][flags['pr-branch'].length - 1];
  updated.push('pr_branch');
}

fs.writeFileSync(statePath, serialize(state));

const updatedUnique = [...new Set(updated)];
const preserved = PROTECTED.filter(f => !updatedUnique.includes(f));

process.stdout.write(`STATUS: ok\n`);
process.stdout.write(`FILE: ${statePath}\n`);
process.stdout.write(`PHASE: ${state.phase}\n`);
process.stdout.write(`UPDATED: ${updatedUnique.length ? updatedUnique.join(', ') : '(none)'}\n`);
process.stdout.write(`PRESERVED: ${preserved.length ? preserved.join(', ') : '(none)'}\n`);
process.exit(0);
