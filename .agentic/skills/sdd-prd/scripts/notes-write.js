#!/usr/bin/env node

'use strict';

const fs = require('fs');
const path = require('path');

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

// --- Field renderers ---

function renderField(key, value, date) {
  switch (key) {
    case 'context':
      return value;
    case 'functionalRequirements':
    case 'nonFunctionalRequirements':
    case 'businessRules':
    case 'acceptanceCriteria':
      return value.map(r => `- ${r}`).join('\n');
    case 'dataRequirements':
      return `- **Entidades:** ${value.entities || '-'}\n- **Origem:** ${value.origin || '-'}\n- **Privacidade / LGPD:** ${value.privacyLgpd || 'sem dados pessoais'}`;
    case 'glossary':
      return value.map(g => `- **${g.term}:** ${g.definition}`).join('\n');
    case 'openQuestions':
      return value.map(q => `- [ ] ${q}`).join('\n');
    case 'decisions':
      return `- [${date}] ${value}`;
    default:
      return String(value);
  }
}

// --- Full write: template + tag replacement ---

function validateTags(rendered) {
  const missing = [...rendered.matchAll(/\{\{(\w+)\}\}/g)].map(m => m[1]);
  if (missing.length > 0) throw new Error(`Tags não substituídas: ${missing.join(', ')}`);
}

function writeNOTES(templatePath, data, date) {
  const n = data.notes || {};
  let tpl = fs.readFileSync(templatePath, 'utf8');

  tpl = tpl.replace('{{title}}', data.title);
  tpl = tpl.replace('{{date}}',  date);

  const fields = {
    context:                  n.context,
    functionalRequirements:   n.functionalRequirements,
    nonFunctionalRequirements:n.nonFunctionalRequirements,
    businessRules:            n.businessRules,
    dataRequirements:         n.dataRequirements,
    glossary:                 n.glossary && n.glossary.length > 0 ? n.glossary : undefined,
    acceptanceCriteria:       n.acceptanceCriteria,
    openQuestions:            n.openQuestions,
    decisions:                n.newDecision || 'NOTES inicializado.',
  };

  for (const [key, value] of Object.entries(fields)) {
    if (value !== undefined) {
      tpl = tpl.replace(`{{${key}}}`, renderField(key, value, date));
    }
  }

  // remove entire section blocks where tag was not replaced (optional sections without data)
  tpl = tpl.replace(/\n## [^\n]+\n\n\{\{\w+\}\}\n/g, '');

  validateTags(tpl);
  return tpl;
}

// --- Merge: parse existing NOTES, update sections ---

function parseSections(content) {
  const sections = {};
  const lines = content.split('\n');
  let currentSection = null;
  let buffer = [];

  for (const line of lines) {
    if (line.startsWith('## ')) {
      if (currentSection !== null) sections[currentSection] = buffer.join('\n').trimEnd();
      currentSection = line.slice(3).trim();
      buffer = [];
    } else if (currentSection !== null) {
      buffer.push(line);
    }
  }

  if (currentSection !== null) sections[currentSection] = buffer.join('\n').trimEnd();
  return sections;
}

function renderNOTES(sections, header) {
  const order = [
    'Objetivos e Contexto',
    'Requisitos Funcionais',
    'Requisitos Não Funcionais',
    'Regras de Negócio',
    'Requisitos de Dados',
    'Linguagem Ubíqua (Glossário)',
    'Critérios de Aceitação (rascunho BDD)',
    'Dúvidas Pendentes',
    'Registro de Decisões (append-only)',
  ];

  let out = header + '\n';
  for (const sec of order) {
    if (sections[sec] !== undefined) out += `\n## ${sec}\n${sections[sec]}\n`;
  }
  for (const key of Object.keys(sections)) {
    if (!order.includes(key)) out += `\n## ${key}\n${sections[key]}\n`;
  }
  return out;
}

function mergeNOTES(notesPath, data, date) {
  const n = data.notes || {};
  const raw = fs.readFileSync(notesPath, 'utf8');

  const firstSection = raw.indexOf('\n## ');
  const header = (firstSection > 0 ? raw.slice(0, firstSection).trim() : '')
    .replace(/Última atualização: \d{4}-\d{2}-\d{2}/, `Última atualização: ${date}`);

  const existing = parseSections(raw);

  if (n.context)                   existing['Objetivos e Contexto']                  = `\n${n.context}`;
  if (n.functionalRequirements)    existing['Requisitos Funcionais']                 = '\n' + n.functionalRequirements.map(r => `- ${r}`).join('\n');
  if (n.nonFunctionalRequirements) existing['Requisitos Não Funcionais']             = '\n' + n.nonFunctionalRequirements.map(r => `- ${r}`).join('\n');
  if (n.businessRules)             existing['Regras de Negócio']                     = '\n' + n.businessRules.map(r => `- ${r}`).join('\n');
  if (n.dataRequirements) {
    const dr = n.dataRequirements;
    existing['Requisitos de Dados'] = `\n- **Entidades:** ${dr.entities || '-'}\n- **Origem:** ${dr.origin || '-'}\n- **Privacidade / LGPD:** ${dr.privacyLgpd || 'sem dados pessoais'}`;
  }
  if (n.glossary && n.glossary.length > 0) existing['Linguagem Ubíqua (Glossário)'] = '\n' + n.glossary.map(g => `- **${g.term}:** ${g.definition}`).join('\n');
  if (n.acceptanceCriteria)        existing['Critérios de Aceitação (rascunho BDD)'] = '\n' + n.acceptanceCriteria.map(ac => `- ${ac}`).join('\n');
  if (n.openQuestions)             existing['Dúvidas Pendentes']                     = '\n' + n.openQuestions.map(q => `- [ ] ${q}`).join('\n');

  const newDecision = n.newDecision || 'NOTES atualizado';
  const prev = existing['Registro de Decisões (append-only)'] || '';
  existing['Registro de Decisões (append-only)'] = prev
    ? `${prev}\n- [${date}] ${newDecision}`
    : `\n- [${date}] ${newDecision}`;

  return renderNOTES(existing, header);
}

// --- Schema ---

const SCHEMA = {
  description: 'Cria NOTES.md via template (escrita completa) ou faz merge em arquivo existente. Seções não-append substituem o estado anterior. "Registro de Decisões" sempre acumula.',
  required: ['module', 'feature', 'title'],
  fields: {
    module:  'string — módulo DDD (normalizado para kebab-case)',
    feature: 'string — nome da feature (normalizado para kebab-case)',
    title:   'string — título legível',
    notes: {
      context:                  'string?',
      functionalRequirements:   'string[]?',
      nonFunctionalRequirements:'string[]?',
      businessRules:            'string[]?',
      dataRequirements:         '{ entities: string, origin: string, privacyLgpd: string }?',
      glossary:                 '{ term: string, definition: string }[]?',
      acceptanceCriteria:       'string[]?',
      openQuestions:            'string[]?',
      newDecision:              'string? — append-only no Registro de Decisões',
    },
  },
  behavior: 'NOTES.md inexistente → escrita completa via template. NOTES.md existente → merge por seção.',
  output: '{ notes, date }',
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
    console.error('Usage: node notes-write.js --input <json-file>');
    console.error('       node notes-write.js --help');
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

  const module     = toKebabCase(data.module);
  const feature    = toKebabCase(data.feature);
  const folderPath = path.join('docs', module, feature);
  const notesPath  = path.join(folderPath, 'NOTES.md');
  const templatePath = path.join('.agentic/skills/sdd-prd/references/notes-template.md');

  fs.mkdirSync(folderPath, { recursive: true });

  const date = getDate();
  const content = fs.existsSync(notesPath)
    ? mergeNOTES(notesPath, data, date)
    : writeNOTES(templatePath, data, date);

  fs.writeFileSync(notesPath, content, 'utf8');

  console.log(JSON.stringify({ notes: notesPath, date }));
}

main();
