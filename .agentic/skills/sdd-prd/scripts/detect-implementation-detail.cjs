#!/usr/bin/env node
'use strict';

// Scans a PRD.md for implementation-detail keywords that do not belong in a product document.
// A PRD must describe WHAT and WHY, never HOW (architecture, files, libs, layers).
//
// Exit codes:
//   0 — no implementation details detected (stdout: structured clean result)
//   1 — implementation details found (stderr: detailed violation list + action per type)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const fs = require('fs');

const HELP = `
Usage: node scripts/detect-implementation-detail.cjs <path/to/PRD.md>

Scans a PRD.md for technical implementation keywords that violate the
product-only rule: no architecture, file paths, libraries, or layer/contract
decisions in a PRD.

Exit codes:
  0   No implementation details detected
  1   Implementation details found — detailed violations on stderr
  2   Usage or runtime error (file not found, missing argument)

Examples:
  node scripts/detect-implementation-detail.cjs docs/auth/login/PRD.md
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const [,, prdPath] = process.argv;

if (!prdPath) {
  process.stderr.write([
    'Error: <path/to/PRD.md> is required.',
    'Usage: node scripts/detect-implementation-detail.cjs <path/to/PRD.md>',
  ].join('\n') + '\n');
  process.exit(2);
}

if (!fs.existsSync(prdPath)) {
  process.stderr.write([
    `Error: file not found: "${prdPath}"`,
    `Action: confirm the path is correct and the PRD was written before running this check.`,
  ].join('\n') + '\n');
  process.exit(2);
}

const FORBIDDEN = [
  { label: 'package (Go)',                   pattern: /\bpackage\b/gi,                             suggestion: 'Descreva o módulo de negócio, não o pacote de código.' },
  { label: 'struct (Go)',                    pattern: /\bstruct\b/gi,                              suggestion: 'Use "entidade" ou "dado" em vez de struct.' },
  { label: 'interface (technical)',          pattern: /\binterface\b/gi,                           suggestion: 'Use "contrato de negócio" ou "integração" se for produto.' },
  { label: 'repository pattern',             pattern: /\brepository\b/gi,                          suggestion: 'Descreva a necessidade de persistência, não o padrão técnico.' },
  { label: 'handler',                        pattern: /\bhandler\b/gi,                             suggestion: 'Use "endpoint de negócio" ou descreva a operação pelo nome.' },
  { label: 'controller',                     pattern: /\bcontroller\b/gi,                          suggestion: 'Descreva o fluxo do usuário, não a camada de controle.' },
  { label: 'middleware',                     pattern: /\bmiddleware\b/gi,                          suggestion: 'Descreva o comportamento esperado (ex: "autenticação obrigatória").' },
  { label: 'database technology',            pattern: /\b(sql|mysql|postgres|mongodb|redis)\b/gi,  suggestion: 'Use "armazenamento" ou "banco de dados" sem nomear a tecnologia.' },
  { label: 'JWT',                            pattern: /\bjwt\b/gi,                                 suggestion: 'Descreva a necessidade (ex: "sessão autenticada") sem citar o token.' },
  { label: 'HTTP status code',               pattern: /\bhttp\s+[245]\d{2}\b/gi,                  suggestion: 'Descreva o comportamento de negócio (ex: "retorna erro ao usuário").' },
  { label: 'endpoint (impl)',                pattern: /\bendpoint\b/gi,                            suggestion: 'Descreva a operação pelo nome de negócio, não pela rota técnica.' },
  { label: 'file extension',                pattern: /\.[a-z]{2,4}\b(?=\s|$)/gi,                  suggestion: 'Remova referências a arquivos de código — não pertencem ao PRD.' },
  { label: 'technical layer (port/adapter)', pattern: /\b(port|adapter|usecase layer|domain layer)\b/gi, suggestion: 'Use terminologia de produto — essas são decisões de arquitetura.' },
];

const content    = fs.readFileSync(prdPath, 'utf8');
const violations = [];

for (const { label, pattern, suggestion } of FORBIDDEN) {
  pattern.lastIndex = 0;
  const matches = content.match(pattern);
  if (matches) {
    violations.push({ label, count: matches.length, occurrences: [...new Set(matches)].slice(0, 3), suggestion });
  }
}

if (violations.length === 0) {
  process.stdout.write(`STATUS: clean\n`);
  process.stdout.write(`FILE: ${prdPath}\n`);
  process.stdout.write(`RULES_CHECKED: ${FORBIDDEN.length}\n`);
  process.stdout.write(`VIOLATIONS: 0\n`);
  process.exit(0);
}

process.stderr.write(`STATUS: violations_found\n`);
process.stderr.write(`FILE: ${prdPath}\n`);
process.stderr.write(`RULES_CHECKED: ${FORBIDDEN.length}\n`);
process.stderr.write(`VIOLATIONS: ${violations.length}\n\n`);
violations.forEach((v, i) => {
  process.stderr.write(`[${i + 1}] ${v.label}\n`);
  process.stderr.write(`    Occurrences: ${v.count}x — sample: ${v.occurrences.join(', ')}\n`);
  process.stderr.write(`    Action: ${v.suggestion}\n\n`);
});
process.stderr.write(`Rule: PRD descreve comportamento e intenção de produto — nunca decisões técnicas.\n`);
process.exit(1);
