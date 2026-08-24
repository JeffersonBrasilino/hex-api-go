#!/usr/bin/env node
'use strict';

// Reads intake form text from stdin.
// Extracts a structured gap list and defaults for the PRD interview.
// Always exits 0 — gap detection is informational data, not a failure.
//
// Exit codes:
//   0 — analysis complete (stdout: gap list, defaults, behavior count)
//   2 — usage or runtime error (stderr: what failed + usage hint)

const readline = require('readline');

const HELP = `
Usage: node scripts/parse-intake.cjs < intake.txt

Parses a filled PRD intake form and outputs a structured gap list.
Output is plain-text key: value pairs designed for low token overhead.

Output fields (stdout):
  GAPS: <comma-separated list> | nenhum
  DEFAULTS: <key=value pairs separated by |>   (omitted when empty)
  COMPORTAMENTOS: N encontrados                 (omitted when < 2)
  CAMPOS_PREENCHIDOS: N de 6

Exit codes:
  0   Analysis complete (always — gaps are data, not errors)
  2   Usage or runtime error

Examples:
  echo "$INTAKE" | node scripts/parse-intake.cjs
  node scripts/parse-intake.cjs < intake.txt
`.trim();

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  process.stdout.write(HELP + '\n');
  process.exit(0);
}

const PLACEHOLDER = /^\s*\[.*\]\s*$/;
const HIGH_TRAFFIC = /\b(e-?commerce|pagamento|payment|autentica|auth|api p[uú]blica|login|checkout|banco|financ)/i;

function extractField(text, labelPattern) {
  const re = new RegExp(`\\*\\*${labelPattern}:\\*\\*\\s*([\\s\\S]*?)(?=\\*\\*[^:]+:\\*\\*|$)`, 'i');
  const m = text.match(re);
  if (!m) return null;
  const val = m[1].trim();
  return (!val || PLACEHOLDER.test(val)) ? null : val;
}

function countBehaviors(value) {
  if (!value) return 0;
  const bullets  = (value.match(/^[\s]*[-*•]\s+.+/gm)  || []).length;
  const numbered = (value.match(/^\s*\d+[.)]\s+.+/gm)  || []).length;
  const maxList  = Math.max(bullets, numbered);
  if (maxList >= 2) return maxList;
  return value.split('\n').filter(l => l.trim()).length;
}

async function main() {
  const lines = [];
  const rl = readline.createInterface({ input: process.stdin, crlfDelay: Infinity });
  for await (const line of rl) lines.push(line);
  const text = lines.join('\n');

  if (!text.trim()) {
    process.stderr.write([
      'Error: no input received via stdin.',
      'Expected: filled PRD intake form with fields **Problema:**, **Usuários:**, etc.',
      'Usage: echo "$INTAKE" | node scripts/parse-intake.cjs',
    ].join('\n') + '\n');
    process.exit(2);
  }

  const problema      = extractField(text, 'Problema');
  const usuarios      = extractField(text, 'Usu[aá]rios');
  const comportamento = extractField(text, 'O que deve fazer');
  const naoEntra      = extractField(text, 'O que N[ÃA]O entra');
  const criterio      = extractField(text, 'Crit[eé]rio de sucesso');
  const restricoes    = extractField(text, 'Restri[cç][õo]es');

  const gaps     = [];
  const defaults = [];
  let filled     = 0;

  if (problema)   filled++; else gaps.push('problema');
  if (usuarios)   filled++; else gaps.push('usuarios');

  const bCount = countBehaviors(comportamento);
  if (bCount >= 2) filled++;
  else             gaps.push('comportamentos (mínimo 2 comportamentos distintos)');

  if (naoEntra)  filled++; else defaults.push('nao_entra=sem restrições explícitas de escopo');
  if (criterio)  filled++; else defaults.push('criterio=derivar dos comportamentos listados');

  const highTraffic = HIGH_TRAFFIC.test(text);
  if (restricoes) {
    filled++;
  } else if (highTraffic) {
    gaps.push('restricoes_perf (domínio de alto tráfego detectado)');
  } else {
    defaults.push('perf=padrão — sem restrição especial');
  }

  process.stdout.write(`GAPS: ${gaps.length > 0 ? gaps.join(', ') : 'nenhum'}\n`);
  if (defaults.length > 0) process.stdout.write(`DEFAULTS: ${defaults.join(' | ')}\n`);
  if (bCount >= 2)          process.stdout.write(`COMPORTAMENTOS: ${bCount} encontrados\n`);
  process.stdout.write(`CAMPOS_PREENCHIDOS: ${filled} de 6\n`);
  process.exit(0);
}

main().catch(err => {
  process.stderr.write(`Error: ${err.message}\n`);
  process.exit(2);
});
