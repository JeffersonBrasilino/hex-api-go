#!/usr/bin/env node
'use strict';

// Validates a PLAN.md file against the plan-schema structural requirements.
// Checks: required sections, task ID format, required fields per task,
// depends_on references resolve, no circular dependencies, no placeholders.
//
// Exit codes:
//   0 — plan is valid (stdout: structured validation summary)
//   1 — validation failed (stderr: numbered issue list with found/expected/action)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/validate-plan.cjs <path/to/PLAN.md>

Validates a PLAN.md file against the sdd-plan structural requirements:
  - Required sections present: "## 1. Plan", "## 3. Execution", "## 4. Return"
  - Every task ID matches TASK-[LAYER]-[CONCERN], LAYER in DOM/APP/INFRA/MOD/TEST
  - Every task has: File, Reason, Requirement (RF-0X/RNF-0X/RN-0X id(s)), depends_on,
    parallel_group, a 5-dimension complexity table with score/tier, Sub-tasks,
    Completion criterion
  - Every depends_on reference resolves to a task declared in the plan
  - No circular dependencies among tasks
  - No unfilled template placeholders remaining

Exit codes:
  0   Plan is valid
  1   Validation failed — detailed issues on stderr
  2   Usage or runtime error (file not found, missing argument)

Examples:
  node scripts/validate-plan.cjs docs/auth/login/PLAN.md
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [,, planPath] = process.argv;

if (!planPath) {
  process.stderr.write([
    'Error: <path/to/PLAN.md> is required.',
    'Usage: node scripts/validate-plan.cjs <path/to/PLAN.md>',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(planPath)) {
  process.stderr.write([
    `Error: file not found: "${planPath}"`,
    `Action: run scaffold.cjs to initialise PLAN.md, then write its content before validating.`,
  ].join('\n') + '\n');
  process.exit(2);
}

const content = fs.readFileSync(planPath, 'utf8');
const issues  = [];

const REQUIRED_SECTIONS = ['## 1. Plan', '## 3. Execution', '## 4. Return'];
const missingSections = REQUIRED_SECTIONS.filter(s => !content.includes(s));
if (missingSections.length > 0) {
  issues.push({
    type: 'missing sections',
    found: `${REQUIRED_SECTIONS.length - missingSections.length} of ${REQUIRED_SECTIONS.length}`,
    missing: missingSections,
    action: 'Add each missing section to the plan in the correct order.',
  });
}

// Isolate Section 1 (the actual task declarations, not the execution mirror in Section 3).
const section1Match = content.match(/## 1\. Plan([\s\S]*?)(?=\n## 3\. Execution|$)/);
const section1       = section1Match ? section1Match[1] : '';

const TASK_START = /^- \[[ x/]\] \*\*(TASK-(?:DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+)\s*—/gm;
const starts = [];
let m;
while ((m = TASK_START.exec(section1)) !== null) {
  starts.push({ id: m[1], index: m.index });
}

if (starts.length === 0) {
  issues.push({
    type: 'no tasks found',
    found: '0 tasks in section "## 1. Plan"',
    expected: 'at least 1 task block matching "- [ ] **TASK-[LAYER]-[CONCERN] — ..."',
    action: 'Add task blocks following plan-schema.md, or check the task heading format.',
  });
}

// Also flag task headings that don't match the semantic ID format at all.
const ANY_TASK_HEADING = /^- \[[ x/]\] \*\*(TASK-[^\s]+)\s*—/gm;
const validLayers = /^TASK-(DOM|APP|INFRA|MOD|TEST)-[A-Z0-9-]+$/;
const malformedIds = [];
let m2;
while ((m2 = ANY_TASK_HEADING.exec(section1)) !== null) {
  if (!validLayers.test(m2[1])) malformedIds.push(m2[1]);
}
if (malformedIds.length > 0) {
  issues.push({
    type: 'malformed task ID',
    found: malformedIds.join(', '),
    expected: 'TASK-[LAYER]-[CONCERN], LAYER in DOM/APP/INFRA/MOD/TEST',
    action: 'Rename the task ID to follow the semantic format.',
  });
}

const blocks = starts.map((s, i) => {
  const end = i + 1 < starts.length ? starts[i + 1].index : section1.length;
  return { id: s.id, text: section1.slice(s.index, end) };
});

const declaredIds = new Set(blocks.map(b => b.id));
const depsById     = {};
const waveById     = {};
const placeholders = [];

const FIELD_CHECKS = [
  { name: 'File',                 re: /\*\*File:\*\*\s*`[^`]+`/ },
  { name: 'Reason',               re: /\*\*Reason:\*\*\s*\S+/ },
  { name: 'Requirement',          re: /\*\*Requirement:\*\*\s*\S+/ },
  { name: 'depends_on',           re: /`depends_on:`\s*`\[[^\]]*\]`/ },
  { name: 'parallel_group',       re: /`parallel_group:`\s*`\[[^\]]*\]`/ },
  { name: 'wave',                 re: /`wave:`\s*`?\d+`?/ },
  { name: 'Sub-tasks',            re: /\*\*Sub-tasks:\*\*[\s\S]*?- \[[ x/]\]/ },
  { name: 'Completion criterion', re: /\*\*Completion criterion:\*\*\s*\S+/ },
];
const COMPLEXITY_DIMS = ['scope', 'ambiguity', 'coupling', 'novelty', 'reversibility'];

for (const block of blocks) {
  const missingFields = FIELD_CHECKS.filter(f => !f.re.test(block.text)).map(f => f.name);

  const missingDims = COMPLEXITY_DIMS.filter(dim => !new RegExp(`\\|\\s*${dim}\\s*\\|`, 'i').test(block.text));
  if (missingDims.length > 0) missingFields.push(`complexity dimension(s): ${missingDims.join(', ')}`);
  if (!/\|\s*\*\*score\*\*\s*\|/i.test(block.text)) missingFields.push('complexity score row');
  if (!/\|\s*\*\*tier\*\*\s*\|/i.test(block.text)) missingFields.push('complexity tier row');

  if (missingFields.length > 0) {
    issues.push({
      type: 'task missing required field(s)',
      taskId: block.id,
      missing: missingFields,
      action: `Add the missing field(s) to ${block.id}'s task block, following plan-schema.md.`,
    });
  }

  const requirementMatch = block.text.match(/\*\*Requirement:\*\*\s*([^\n]+)/);
  if (requirementMatch) {
    const reqIds = requirementMatch[1].split(',').map(s => s.trim()).filter(Boolean);
    const malformedReqIds = reqIds.filter(id => !/^(RF|RNF|RN)-\d+$/.test(id));
    if (malformedReqIds.length > 0) {
      issues.push({
        type: 'malformed requirement id',
        taskId: block.id,
        found: malformedReqIds.join(', '),
        expected: 'RF-0X, RNF-0X, or RN-0X id(s) as declared in the PRD, comma-separated',
        action: `Correct ${block.id}'s Requirement field to reference valid PRD requirement id(s).`,
      });
    }
  }

  const depsMatch = block.text.match(/`depends_on:`\s*`\[([^\]]*)\]`/);
  const deps = depsMatch && depsMatch[1].trim()
    ? depsMatch[1].split(',').map(s => s.trim()).filter(Boolean)
    : [];
  depsById[block.id] = deps;

  const waveMatch = block.text.match(/`wave:`\s*`?(\d+)`?/);
  waveById[block.id] = waveMatch ? parseInt(waveMatch[1], 10) : null;

  const unresolved = deps.filter(d => !declaredIds.has(d));
  if (unresolved.length > 0) {
    issues.push({
      type: 'task with invalid depends_on reference',
      found: `${block.id} depends on ${unresolved.join(', ')}`,
      expected: `${unresolved.join(', ')} must exist in the plan as (a) declared task(s)`,
      action: `Add the missing task(s) as declared tasks, or correct the depends_on reference in ${block.id}.`,
    });
  }

  // Detect unresolved template placeholders inside the block, excluding valid
  // empty lists ("[]"), valid depends_on/parallel_group/wave values, and markdown
  // checkbox syntax ("[ ]"/"[x]"/"[/]") on sub-task lines — a task or sub-task
  // already marked done during execution is not a leftover template placeholder.
  const bracketMatches = block.text.match(/\[[^\]]*\]/g) || [];
  for (const raw of bracketMatches) {
    const inner = raw.slice(1, -1).trim();
    if (inner === '') continue;
    if (/^[ x/]$/.test(inner)) continue;
    if (/^(TASK-[A-Z0-9-]+)(\s*,\s*TASK-[A-Z0-9-]+)*$/.test(inner)) continue;
    if (/^(domain|application|infrastructure|module|tests|config)$/.test(inner)) continue;
    placeholders.push(`${block.id}: ${raw}`);
  }
}

if (placeholders.length > 0) {
  issues.push({
    type: 'unfilled placeholders',
    found: `${placeholders.length} placeholder(s) remaining`,
    sample: placeholders.slice(0, 5),
    action: 'Replace every [placeholder] with actual plan content.',
  });
}

// Wave consistency: a task can never share a wave with — or precede — a task it
// depends on. This is what makes `implement-wave=N` in sdd-workflow trustworthy:
// wave numbers are supposed to be written by compute-waves.cjs, but a hand-edited
// depends_on or a stale wave value (compute-waves.cjs not re-run after an edit) would
// silently let the orchestrator schedule a task before its dependency. Catch it here.
const waveInconsistencies = [];
for (const block of blocks) {
  const taskWave = waveById[block.id];
  if (taskWave === null) continue; // already reported as a missing required field above
  for (const dep of depsById[block.id] || []) {
    const depWave = waveById[dep];
    if (depWave === null || depWave === undefined) continue;
    if (taskWave <= depWave) {
      waveInconsistencies.push(`${block.id} (wave ${taskWave}) depends on ${dep} (wave ${depWave})`);
    }
  }
}
if (waveInconsistencies.length > 0) {
  issues.push({
    type: 'wave inconsistency',
    found: waveInconsistencies.join('; '),
    expected: 'a dependent task always has a strictly higher wave than every task in its depends_on',
    action: 'Re-run compute-waves.cjs to recompute waves from the current depends_on/parallel_group, or fix the depends_on that no longer matches the wave field.',
  });
}

// Circular dependency detection (DFS with recursion stack).
const circular = [];
const visited = new Set();
function dfs(id, stack) {
  if (stack.includes(id)) {
    circular.push([...stack, id].join(' -> '));
    return;
  }
  if (visited.has(id)) return;
  visited.add(id);
  for (const dep of depsById[id] || []) {
    if (declaredIds.has(dep)) dfs(dep, [...stack, id]);
  }
}
for (const id of declaredIds) dfs(id, []);

if (circular.length > 0) {
  issues.push({
    type: 'circular dependency',
    found: circular.join('; '),
    expected: 'a dependency graph with no cycles',
    action: 'Break the cycle by removing or reordering one of the depends_on references.',
  });
}

if (issues.length === 0) {
  process.stdout.write(`STATUS: valid\n`);
  process.stdout.write(`FILE: ${planPath}\n`);
  process.stdout.write(`TASKS: ${blocks.length}\n`);
  process.stdout.write(`DEPENDENCIES_OK: true\n`);
  process.stdout.write(`CIRCULAR_DEPS: 0\n`);
  process.stdout.write(`PLACEHOLDERS: 0\n`);
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
  if (issue.missing)  (Array.isArray(issue.missing) ? issue.missing : [issue.missing]).forEach(x => process.stderr.write(`    Missing:  ${x}\n`));
  if (issue.sample)   issue.sample.forEach(s => process.stderr.write(`    Sample:   ${s}\n`));
  if (issue.action)   process.stderr.write(`    Action:   ${issue.action}\n`);
  process.stderr.write('\n');
});
process.exit(1);
