# Relatório de Avaliação: `sdd-prd-script` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.700 |
| **delta** | **+0.300** |
| **value_tier** | **moderado** |

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.846390 | $0.169278 | claude-sonnet-5 |
| without_skill | $0.624615 | $0.124923 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.221775** | — | — |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-parse-intake-cjs-lgpd-scan-cjs-run-on-intake-but-gaps-nenhum | 0.50 (1.0 vs 0.5) | Sem a skill, o modelo falhou em mencionar explicitamente a execução de parse-intake.cjs e não deixou claro que a Fase 2 (entrevista do PM) deve rodar mesmo quando GAPS: nenhum. |
| eval-approval-gate-step-2-still-required-before-prd-md-is-written | 0.50 (1.0 vs 0.5) | Sem a skill, faltaram os passos explícitos de parse-intake.cjs/lgpd-scan.cjs/entrevista do PM e o formato padronizado exato do resumo de aprovação. |

## Baseline confirmado (ambos ≥ 0.95)

- eval-lgpd-dimensions-asked-one-at-a-time-each-with-a-recommended (with_skill 1.0, without_skill 1.0)

## Lacunas da skill (with_skill < 1.0)

Nenhuma. A configuração with_skill obteve pontuação máxima (1.0) em todas as 5 avaliações, sem falhas registradas nas asserções.

Observação: dois casos ficaram abaixo do limiar de valor claro (delta ≥ 0.40) por apresentarem delta de 0.25 — eval-validate-prd-cjs-and-detect-implementation-detail-cjs-both-g (1.0 vs 0.75, sem_skill falhou em capturar o detalhe de implementação removido na seção Contexto Técnico do NOTES.md) e eval-technical-detail-volunteered-mid-interview-triggers-merge-no (1.0 vs 0.75, sem_skill foi vago sobre capturar o detalhe imediatamente via merge-notes.cjs em vez de adiar). Ambos reforçam, ainda que abaixo do limiar formal, a disciplina adicional trazida pela skill.

## Recomendação

Skill pronta: com with_skill atingindo pontuação perfeita (1.0) e um delta moderado (0.300) explicado por lacunas genuínas do baseline em ordenação de execução de scripts, rigor no tratamento de LGPD e disciplina de captura no NOTES.md, as falhas observadas em without_skill validam exatamente que o fluxo explícito e orientado por scripts da skill é o que previne essas omissões.
