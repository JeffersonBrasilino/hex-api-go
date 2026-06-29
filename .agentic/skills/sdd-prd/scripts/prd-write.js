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

function validateTags(rendered) {
  const missing = [...rendered.matchAll(/\{\{(\w+)\}\}/g)].map(m => m[1]);
  if (missing.length > 0) throw new Error(`Tags não substituídas: ${missing.join(', ')}`);
}

// --- Section map: field key → markdown heading ---

const SECTION_MAP = {
  overview:                  '## 1. Visão Geral',
  objectives:                '## 2. Objetivos',
  userStories:               '## 3. Histórias de Usuário',
  functionalRequirements:    '## 4. Requisitos Funcionais',
  nonFunctionalRequirements: '## 5. Requisitos Não Funcionais',
  businessRules:             '## 6. Regras de Negócio',
  dataRequirements:          '## 7. Requisitos de Dados',
  acceptanceCriteria:        '## 8. Critérios de Aceitação',
  outOfScope:                '## 9. Fora de Escopo',
  openQuestions:             '## 10. Perguntas Abertas',
  futureConsiderations:      '## 11. Considerações Futuras',
};

const HEADING_TO_FIELD = Object.fromEntries(Object.entries(SECTION_MAP).map(([k, v]) => [v, k]));

// --- Field renderers ---

function renderField(key, value) {
  switch (key) {
    case 'overview':
      return value;
    case 'objectives':
      return value.map(o => `- ${o}`).join('\n');
    case 'userStories':
      return value.map(e => `- **${e.epic}**\n${(e.stories || []).map(s => `  - ${s}`).join('\n')}`).join('\n');
    case 'functionalRequirements':
      return value.map(r => `- **${r.id}:** ${r.description}`).join('\n');
    case 'nonFunctionalRequirements':
      return value.map(r => `- **${r.id} (${r.qualifier}):** ${r.description}`).join('\n');
    case 'businessRules':
      return value.map(r => `- **${r.id}:** ${r.description}`).join('\n');
    case 'dataRequirements':
      return `- **Entidades envolvidas:** ${value.entities}\n- **Origem dos dados:** ${value.origin}\n- **Privacidade / LGPD:** ${value.privacyLgpd}`;
    case 'acceptanceCriteria':
      return value.map((ac, i) => `- **Cenário ${i + 1} — ${ac.name}:** ${ac.scenario}`).join('\n');
    case 'outOfScope':
      return value.map(i => `- ${i}`).join('\n');
    case 'openQuestions':
      return value.length ? value.map(i => `- ${i}`).join('\n') : 'Nenhuma.';
    case 'futureConsiderations':
      return value.map(i => `- ${i}`).join('\n');
    default:
      return String(value);
  }
}

// --- Full write: template + tag replacement ---

function writePRD(templatePath, data, meta) {
  let tpl = fs.readFileSync(templatePath, 'utf8');

  tpl = tpl.replace('{{title}}',   data.title);
  tpl = tpl.replace('{{author}}',  meta.author);
  tpl = tpl.replace('{{date}}',    meta.date);
  tpl = tpl.replace('{{status}}',  data.status || 'Draft');
  tpl = tpl.replace('{{version}}', meta.version);

  for (const key of Object.keys(SECTION_MAP)) {
    const value = (data.sections || {})[key];
    if (value !== undefined) {
      tpl = tpl.replace(`{{${key}}}`, renderField(key, value));
    }
  }

  validateTags(tpl);
  return tpl;
}

// --- Patch: parse existing PRD, replace only changed sections ---

function parsePRD(content) {
  const lines = content.split('\n');
  let headerLines = [];
  let currentHeading = null;
  let currentBody = [];
  const sections = {};
  const order = [];

  for (const line of lines) {
    if (/^## \d+\./.test(line)) {
      if (currentHeading !== null) {
        sections[currentHeading] = currentBody.join('\n').trimEnd();
      } else {
        headerLines = currentBody.slice();
      }
      currentHeading = line;
      order.push(line);
      currentBody = [];
    } else {
      currentBody.push(line);
    }
  }

  if (currentHeading !== null) {
    sections[currentHeading] = currentBody.join('\n').trimEnd();
  }

  return { header: headerLines.join('\n').trimEnd(), sections, order };
}

function assemblePRD(parsed) {
  let out = parsed.header + '\n';
  for (const heading of parsed.order) {
    out += `\n${heading}\n\n${parsed.sections[heading]}\n`;
  }
  return out;
}

function extractVersion(header) {
  const match = header.match(/\|\s*\*\*Versão\*\*\s*\|\s*(\d+)\.(\d+)/);
  if (!match) return '1.0';
  return `${match[1]}.${parseInt(match[2]) + 1}`;
}

function updateHeaderMeta(header, version, date) {
  return header
    .replace(/(\|\s*\*\*Versão\*\*\s*\|)[^|]+(\|)/, `$1 ${version} $2`)
    .replace(/(\|\s*\*\*Última Atualização\*\*\s*\|)[^|]+(\|)/, `$1 ${date} $2`);
}

function patchPRD(prdPath, data, date) {
  const existing = fs.readFileSync(prdPath, 'utf8');
  const parsed = parsePRD(existing);

  const newVersion = extractVersion(parsed.header);
  parsed.header = updateHeaderMeta(parsed.header, newVersion, date);

  for (const heading of parsed.order) {
    const fieldKey = HEADING_TO_FIELD[heading];
    if (fieldKey && (data.patch || {})[fieldKey] !== undefined) {
      parsed.sections[heading] = '\n' + renderField(fieldKey, data.patch[fieldKey]);
    }
  }

  const result = assemblePRD(parsed);
  validateTags(result);
  return { content: result, version: newVersion };
}

// --- Schema ---

const SCHEMA = {
  description: 'Escreve ou atualiza PRD.md. Dois modos detectados pelo campo presente.',
  required: ['module', 'feature', 'title'],
  modes: {
    'escrita completa': {
      trigger: 'campo "sections" presente',
      fields: {
        module: 'string — módulo DDD (normalizado para kebab-case)',
        feature: 'string — nome da feature (normalizado para kebab-case)',
        title: 'string — título legível',
        status: 'string? — "Draft" (padrão), "Em Revisão" ou "Aprovado"',
        sections: {
          overview: 'string',
          objectives: 'string[]',
          userStories: '{ epic: string, stories: string[] }[]',
          functionalRequirements: '{ id: string, description: string }[]',
          nonFunctionalRequirements: '{ id: string, qualifier: string, description: string }[]',
          businessRules: '{ id: string, description: string }[]',
          dataRequirements: '{ entities: string, origin: string, privacyLgpd: string }',
          acceptanceCriteria: '{ name: string, scenario: string }[]',
          outOfScope: 'string[]',
          openQuestions: 'string[]',
          futureConsiderations: 'string[]',
        },
      },
    },
    'patch parcial': {
      trigger: 'campo "patch" presente — inclua só os campos que mudaram',
      fields: {
        module: 'string',
        feature: 'string',
        title: 'string',
        patch: 'Partial<sections> — mesmos tipos de "escrita completa"',
      },
    },
  },
  output: '{ folder, prd, author, date, version }',
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
    console.error('Usage: node prd-write.js --input <json-file>');
    console.error('       node prd-write.js --help');
    process.exit(1);
  }

  let data;
  try {
    data = JSON.parse(fs.readFileSync(args[inputIdx + 1], 'utf8'));
  } catch (err) {
    console.error(`Failed to parse input: ${err.message}`);
    process.exit(1);
  }

  if (!data.module || !data.feature || !data.title) {
    console.error('Input must have: module, feature, title');
    process.exit(1);
  }

  const module      = toKebabCase(data.module);
  const feature     = toKebabCase(data.feature);
  const folderPath  = path.join('docs', module, feature);
  const prdPath     = path.join(folderPath, 'PRD.md');
  const templatePath = path.join('.agentic/skills/sdd-prd/references/prd-template.md');

  fs.mkdirSync(folderPath, { recursive: true });

  const date    = getDate();
  const author  = getGitAuthor();
  let version   = '1.0';

  if (data.patch) {
    if (!fs.existsSync(prdPath)) {
      console.error(`PRD não encontrado em ${prdPath}. Patch requer arquivo existente.`);
      process.exit(1);
    }
    const result = patchPRD(prdPath, data, date);
    version = result.version;
    fs.writeFileSync(prdPath, result.content, 'utf8');
  } else if (data.sections) {
    fs.writeFileSync(prdPath, writePRD(templatePath, data, { author, date, version }), 'utf8');
  }

  console.log(JSON.stringify({ folder: folderPath, prd: prdPath, author, date, version }));
}

main();
