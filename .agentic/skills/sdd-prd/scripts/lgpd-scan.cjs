#!/usr/bin/env node
'use strict';

// Reads intake text from stdin.
// Detects personal data signals that trigger LGPD compliance questions.
//
// Exit codes:
//   0 — no personal data detected (stdout: structured result)
//   1 — personal data signals found (stderr: detailed findings + required questions)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const readline = require('readline');

const HELP = `
Usage: node scripts/lgpd-scan.cjs < intake.txt

Scans intake form text (via stdin) for personal data signals that require
LGPD compliance questions in the PRD interview.

Exit codes:
  0   No personal data detected
  1   Personal data signals found — agent must ask the 3 LGPD dimensions
  2   Usage or runtime error

Examples:
  echo "$INTAKE" | node scripts/lgpd-scan.cjs
  node scripts/lgpd-scan.cjs < intake.txt
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const SIGNALS = [
  { name: 'CPF',             pattern: /\b\d{3}\.?\d{3}\.?\d{3}-?\d{2}\b|\bcpf\b/gi },
  { name: 'CNPJ',            pattern: /\b\d{2}\.?\d{3}\.?\d{3}\/?\d{4}-?\d{2}\b|\bcnpj\b/gi },
  { name: 'e-mail',          pattern: /\b[\w.+-]+@[\w-]+\.[a-z]{2,}\b|\be-?mail\b/gi },
  { name: 'endereço',        pattern: /\bendere[cç]o\b|\bcep\b|\blogradouro\b|\bbairro\b/gi },
  { name: 'data nascimento', pattern: /\bnascimento\b|\bdata de nascimento\b/gi },
  { name: 'telefone',        pattern: /\b\(?0?\d{2}\)?\s?\d{4,5}[-.]?\d{4}\b|\btelefone\b|\bcelular\b/gi },
  { name: 'nome completo',   pattern: /\bnome completo\b|\bnome do usu[aá]rio\b|\bnome do cliente\b/gi },
];

async function main() {
  const lines = [];
  const rl = readline.createInterface({ input: process.stdin, crlfDelay: Infinity });
  for await (const line of rl) lines.push(line);
  const text = lines.join('\n');

  if (!text.trim()) {
    process.stderr.write([
      'Error: no input received via stdin.',
      'Expected: filled PRD intake form text piped via stdin.',
      'Usage: echo "$INTAKE" | node scripts/lgpd-scan.cjs',
    ].join('\n') + '\n');
    process.exit(2);
  }

  const found = SIGNALS
    .filter(({ pattern }) => { pattern.lastIndex = 0; return pattern.test(text); })
    .map(({ name }) => name);

  if (found.length === 0) {
    process.stdout.write('LGPD: nao_detectado\n');
    process.stdout.write(`SINAIS_VERIFICADOS: ${SIGNALS.map(s => s.name).join(', ')}\n`);
    process.stdout.write('DEFAULT: sem dados pessoais — declarar no resumo de aprovação\n');
    process.exit(0);
  }

  process.stderr.write('LGPD: detectado\n');
  process.stderr.write(`SINAIS: ${found.join(', ')}\n`);
  process.stderr.write(`TOTAL_SINAIS: ${found.length} de ${SIGNALS.length} categorias verificadas\n`);
  process.stderr.write('ACAO: adicionar 3 perguntas ao gap list, uma por turno:\n');
  process.stderr.write('  [1] consentimento — o processamento exige consentimento explícito?\n');
  process.stderr.write('  [2] acesso_controle — quem pode ler/editar esses dados?\n');
  process.stderr.write('  [3] origem_retencao — de onde vêm; há obrigação de anonimização ou prazo de retenção?\n');
  process.exit(1);
}

main().catch(err => {
  process.stderr.write(`Error: ${err.message}\n`);
  process.exit(2);
});
