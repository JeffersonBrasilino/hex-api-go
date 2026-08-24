#!/usr/bin/env node
'use strict';

// Validates a PRD.md file against structural requirements.
// Checks: all 11 sections present, no unfilled placeholders, ≥2 acceptance criteria.
//
// Exit codes:
//   0 — PRD is valid (stdout: structured validation summary)
//   1 — validation failed (stderr: detailed issue list with expected vs found + actions)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/validate-prd.cjs <path/to/PRD.md>

Validates a PRD.md file against the sdd-prd structural requirements:
  - All 11 sections present (Visão Geral through Considerações Futuras)
  - No unfilled template placeholders remaining
  - At least 2 acceptance criteria (1 happy-path + 1 error/edge)

Exit codes:
  0   PRD is valid
  1   Validation failed — detailed issues on stderr
  2   Usage or runtime error (file not found, missing argument)

Examples:
  node scripts/validate-prd.cjs docs/auth/login/PRD.md
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [,, prdPath] = process.argv;

if (!prdPath) {
  process.stderr.write([
    'Error: <path/to/PRD.md> is required.',
    'Usage: node scripts/validate-prd.cjs <path/to/PRD.md>',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(prdPath)) {
  process.stderr.write([
    `Error: file not found: "${prdPath}"`,
    `Action: run scaffold.cjs to initialise PRD.md, then write its content before validating.`,
  ].join('\n') + '\n');
  process.exit(2);
}

const REQUIRED_SECTIONS = [
  '## 1. Visão Geral',
  '## 2. Objetivos',
  '## 3. Histórias de Usuário',
  '## 4. Requisitos Funcionais',
  '## 5. Requisitos Não Funcionais',
  '## 6. Regras de Negócio',
  '## 7. Requisitos de Dados',
  '## 8. Critérios de Aceitação',
  '## 9. Fora de Escopo',
  '## 10. Perguntas Abertas',
  '## 11. Considerações Futuras',
];

const content = fs.readFileSync(prdPath, 'utf8');
const issues  = [];

const missingSections = REQUIRED_SECTIONS.filter(s => !content.includes(s));

const placeholders = (content.match(/\[[^\]]+\]/g) || [])
  .filter(p => !/não se aplica|n\/a/i.test(p));

const acSection  = content.match(/## 8\. Critérios de Aceitação([\s\S]*?)(?=## 9\.|$)/);
const acText     = acSection ? acSection[1] : '';
const scenarios  = (acText.match(/\*\*Cenário/gi) || []).length;
const hasError   = /erro|borda|inv[aá]lid|falha/i.test(acText);

if (missingSections.length > 0) {
  issues.push({
    type: 'missing sections',
    found: `${REQUIRED_SECTIONS.length - missingSections.length} of ${REQUIRED_SECTIONS.length}`,
    missing: missingSections,
    action: 'Add each missing section to the PRD in the correct order.',
  });
}

if (placeholders.length > 0) {
  issues.push({
    type: 'unfilled placeholders',
    found: `${placeholders.length} placeholder(s) remaining`,
    sample: placeholders.slice(0, 3),
    action: 'Replace every [placeholder] with actual PRD content.',
  });
}

if (scenarios < 2) {
  issues.push({
    type: 'insufficient acceptance criteria',
    found: `${scenarios} scenario(s)`,
    expected: 'at least 2 (1 happy-path + 1 error/edge)',
    action: 'Add a **Cenário** with an error, edge case, or invalid input context.',
  });
} else if (!hasError) {
  issues.push({
    type: 'missing error/edge scenario',
    found: `${scenarios} scenario(s), none flagged as error/edge`,
    expected: 'at least 1 scenario covering error, borda, inválido, or falha',
    action: 'Add or rename a scenario to explicitly cover an error or edge case.',
  });
}

if (issues.length === 0) {
  process.stdout.write(`STATUS: valid\n`);
  process.stdout.write(`FILE: ${prdPath}\n`);
  process.stdout.write(`SECTIONS: ${REQUIRED_SECTIONS.length}/${REQUIRED_SECTIONS.length}\n`);
  process.stdout.write(`PLACEHOLDERS: 0\n`);
  process.stdout.write(`CRITERIOS: ${scenarios} (cenário(s) de aceitação encontrados)\n`);
  process.exit(0);
}

process.stderr.write(`STATUS: invalid\n`);
process.stderr.write(`FILE: ${prdPath}\n`);
process.stderr.write(`ISSUES: ${issues.length}\n\n`);
issues.forEach((issue, i) => {
  process.stderr.write(`[${i + 1}] ${issue.type}\n`);
  if (issue.found)    process.stderr.write(`    Found:    ${issue.found}\n`);
  if (issue.expected) process.stderr.write(`    Expected: ${issue.expected}\n`);
  if (issue.missing)  issue.missing.forEach(m => process.stderr.write(`    Missing:  ${m}\n`));
  if (issue.sample)   issue.sample.forEach(s => process.stderr.write(`    Sample:   ${s}\n`));
  if (issue.action)   process.stderr.write(`    Action:   ${issue.action}\n`);
  process.stderr.write('\n');
});
process.exit(1);
