#!/usr/bin/env node

'use strict';

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// --- Utilities ---

function toKebabCase(str) {
  return str
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, '')
    .trim()
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-');
}

function getDate() {
  return new Date().toISOString().slice(0, 10);
}

function getGitAuthor() {
  try {
    return execSync('git config user.name', { encoding: 'utf8' }).trim();
  } catch {
    return 'Desconhecido';
  }
}

function escapeRegex(str) {
  return str.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, '\\$&');
}

function validateTags(rendered) {
  const missing = [...rendered.matchAll(/\{\{(\w+)\}\}/g)].map(m => m[1]);
  if (missing.length > 0) throw new Error(`Tags não substituídas: ${missing.join(', ')}`);
}

// --- Complexity ---

function calcComplexity(c) {
  const dims = ['scope', 'ambiguity', 'coupling', 'novelty', 'reversibility'];
  const score = dims.reduce((sum, d) => sum + (Number(c[d]) || 0), 0) / dims.length;
  const rounded = Math.round(score * 10) / 10;
  const tier = rounded <= 2.0 ? 'low' : rounded <= 3.5 ? 'medium' : 'high';
  return { score: rounded, tier };
}

function overallTier(tasks) {
  const tiers = tasks.map(t => calcComplexity(t.complexity || {}).tier);
  if (tiers.includes('high')) return 'High';
  if (tiers.includes('medium')) return 'Medium';
  return 'Low';
}

// --- Renderers ---

function renderTaskBlock(task) {
  const c = task.complexity || {};
  const { score, tier } = calcComplexity(c);
  const dependsOn = (task.depends_on || []).length > 0
    ? task.depends_on.join(', ')
    : '';

  const subtasksLines = (task.subtasks || [])
    .map((s, i) => `    - [ ] ${task.id}.${i + 1} — ${s}`)
    .join('\n');

  let block = `- [ ] **${task.id} — ${task.description}**\n`;
  block += `  - **File:** \`${task.file}\`\n`;
  block += `  - **Reason:** ${task.reason}\n`;
  block += `  - **Dependencies:**\n`;
  block += `    - \`depends_on:\` \`[${dependsOn}]\`\n`;
  block += `    - \`parallel_group:\` \`${task.parallel_group || 'config'}\`\n`;
  block += `  - **Complexity:**\n`;
  block += `    | Dimension       | Score |\n`;
  block += `    |-----------------|-------|\n`;
  block += `    | scope           | ${c.scope ?? 0}     |\n`;
  block += `    | ambiguity       | ${c.ambiguity ?? 0}     |\n`;
  block += `    | coupling        | ${c.coupling ?? 0}     |\n`;
  block += `    | novelty         | ${c.novelty ?? 0}     |\n`;
  block += `    | reversibility   | ${c.reversibility ?? 0}     |\n`;
  block += `    | **score**       | ${score}   |\n`;
  block += `    | **tier**        | ${tier} |\n`;

  if (tier === 'high' && task.risk_note) {
    block += `    - \`risk_note:\` ${task.risk_note}\n`;
  }

  if (subtasksLines) {
    block += `  - **Sub-tasks:**\n${subtasksLines}\n`;
  }

  block += `  - **Completion criterion:** ${task.completion_criterion || '[To be defined]'}`;

  return block;
}

function renderExecutionBlock(task) {
  let block = `- [ ] **Execution — ${task.id}: ${task.description}**\n`;
  block += `  - *Agent Notes:* [Filled by agent: what exactly was changed and why.]\n`;
  block += `  - *Files Modified:*\n`;
  block += `    - \`${task.file}\`\n`;
  block += `  - *Validation Evidence:* [Test output / log / diff snippet goes here.]\n`;
  block += `  - *Validation Status:* \`✅ Validated\` | \`❌ Failed\` | \`⚠️ Blocked: [reason]\``;
  return block;
}

// --- Full write ---

function writePLAN(templatePath, data, meta) {
  const tasks = data.tasks || [];
  const strategy = data.strategy || {};
  const derivedLayers = [...new Set(tasks.map(t => t.parallel_group).filter(Boolean))].join(' | ');

  const modifiedFiles = tasks
    .map(t => `| \`${t.file}\` | \`Created\` / \`Modified\` | [What changed and why.] |`)
    .join('\n');

  let tpl = fs.readFileSync(templatePath, 'utf8');

  tpl = tpl.replace('{{title}}',            data.title);
  tpl = tpl.replace('{{status}}',           data.status || 'Draft');
  tpl = tpl.replace('{{date}}',             meta.date);
  tpl = tpl.replace('{{author}}',           meta.author);
  tpl = tpl.replace('{{ticket}}',           data.ticket || 'N/A');
  tpl = tpl.replace('{{overall_complexity}}', overallTier(tasks));
  tpl = tpl.replace('{{approach}}',         strategy.approach || '[To be defined]');
  tpl = tpl.replace('{{layers}}',           strategy.layers || derivedLayers || '[To be defined]');
  tpl = tpl.replace('{{skills}}',           strategy.skills || '`adjust-go-code`, `make-unit-tests`, `ddd-module-knowledge`');
  tpl = tpl.replace('{{task_blocks}}',      tasks.map(renderTaskBlock).join('\n\n---\n\n'));
  tpl = tpl.replace('{{execution_blocks}}', tasks.map(renderExecutionBlock).join('\n\n---\n\n'));
  tpl = tpl.replace('{{modified_files}}',   modifiedFiles);

  validateTags(tpl);
  return tpl;
}

// --- Patch: replace a single task block ---

function patchPLAN(planPath, patchTask) {
  const content = fs.readFileSync(planPath, 'utf8');
  const eid = escapeRegex(patchTask.id);

  // Replace task block in Section 1 (anchored to the next separator)
  const taskBlockRegex = new RegExp(
    `(- \\[[ x/]\\] \\*\\*${eid} — .*?)(?=\n\n---)`,
    's'
  );

  const newBlock = renderTaskBlock(patchTask);
  let updated = content.replace(taskBlockRegex, newBlock);

  if (updated === content) {
    throw new Error(`Task ${patchTask.id} não encontrada no PLAN.md`);
  }

  // Update execution block description header if description changed
  const execHeaderRegex = new RegExp(
    `(- \\[[ x/]\\] \\*\\*Execution — ${eid}: )([^*]*)\\*\\*`
  );
  updated = updated.replace(execHeaderRegex, `$1${patchTask.description}**`);

  return updated;
}

// --- Schema ---

const SCHEMA = {
  description: 'Escreve ou atualiza PLAN.md. Dois modos detectados pelo campo presente.',
  required: ['module', 'feature'],
  modes: {
    'escrita completa': {
      trigger: 'campos "title" + "tasks" presentes',
      fields: {
        module: 'string — módulo DDD (normalizado para kebab-case)',
        feature: 'string — nome da feature (normalizado para kebab-case)',
        title: 'string — título legível',
        status: 'string? — "Draft" (padrão)',
        ticket: 'string? — referência do ticket (padrão: N/A)',
        strategy: {
          approach: 'string? — abordagem arquitetural',
          layers: 'string? — camadas afetadas (derivado dos parallel_groups se omitido)',
          skills: 'string? — skills/ferramentas mapeadas',
        },
        tasks: [
          {
            id: 'string — TASK-[LAYER]-[CONCERN]',
            description: 'string — descrição curta',
            file: 'string — caminho relativo do arquivo',
            reason: 'string — por que este arquivo precisa mudar',
            depends_on: 'string[]? — IDs semânticos ([] se nenhum)',
            parallel_group: 'string — domain | application | infrastructure | module | tests | config',
            complexity: {
              scope: 'number 1-5',
              ambiguity: 'number 1-5',
              coupling: 'number 1-5',
              novelty: 'number 1-5',
              reversibility: 'number 1-5',
            },
            risk_note: 'string? — apenas quando tier calculado for high',
            subtasks: 'string[]? — mudanças específicas dentro do arquivo',
            completion_criterion: 'string — como saber que a task está concluída',
          },
        ],
      },
    },
    'patch de task': {
      trigger: 'campo "patch_task" presente',
      note: 'Forneça o objeto task completo, não apenas os campos alterados.',
      fields: {
        module: 'string',
        feature: 'string',
        patch_task: 'Task — mesmos campos do modo escrita completa',
      },
    },
  },
  computed: {
    score: 'avg(scope, ambiguity, coupling, novelty, reversibility) — calculado pelo script',
    tier: 'low ≤ 2.0 | medium ≤ 3.5 | high > 3.5 — calculado pelo script',
    overall_complexity: 'High se qualquer task for high, Medium se qualquer for medium, Low caso contrário',
  },
  output: '{ folder, plan, author, date, tasks_count }',
};

// --- Main ---

function main() {
  const args = process.argv.slice(2);

  if (args.includes('--help') || args.includes('-h')) {
    console.log(JSON.stringify(SCHEMA, null, 2));
    process.exit(0);
  }

  const inputIdx = args.indexOf('--input');
  if (inputIdx === -1 || !args[inputIdx + 1]) {
    console.error('Usage: node plan-write.js --input <json-file>');
    console.error('       node plan-write.js --help');
    process.exit(1);
  }

  let data;
  try {
    data = JSON.parse(fs.readFileSync(args[inputIdx + 1], 'utf8'));
  } catch (err) {
    console.error(`Failed to parse input: ${err.message}`);
    process.exit(1);
  }

  if (!data.module || !data.feature) {
    console.error('Input must have: module, feature');
    process.exit(1);
  }

  const module       = toKebabCase(data.module);
  const feature      = toKebabCase(data.feature);
  const folderPath   = path.join('docs', module, feature);
  const planPath     = path.join(folderPath, 'PLAN.md');
  const templatePath = path.join('.agentic/skills/sdd-plan/references/plan-schema.md');

  fs.mkdirSync(folderPath, { recursive: true });

  const date   = getDate();
  const author = getGitAuthor();
  let content;
  let tasksCount;

  if (data.patch_task) {
    if (!fs.existsSync(planPath)) {
      console.error(`PLAN.md não encontrado em ${planPath}. Patch requer arquivo existente.`);
      process.exit(1);
    }
    content = patchPLAN(planPath, data.patch_task);
    tasksCount = (data.patch_task.id ? 1 : 0);
  } else if (data.tasks) {
    if (!data.title) {
      console.error('Input must have: title (required for escrita completa)');
      process.exit(1);
    }
    content = writePLAN(templatePath, data, { author, date });
    tasksCount = data.tasks.length;
  } else {
    console.error('Input must have either "tasks" (escrita completa) or "patch_task" (patch de task)');
    process.exit(1);
  }

  fs.writeFileSync(planPath, content, 'utf8');
  console.log(JSON.stringify({ folder: folderPath, plan: planPath, author, date, tasks_count: tasksCount }));
}

main();
