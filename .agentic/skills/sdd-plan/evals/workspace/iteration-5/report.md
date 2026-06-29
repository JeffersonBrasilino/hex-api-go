# Relatório de Avaliação: `spec-plan` — iteração 5

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.653 |
| **delta** | **+0.347** |
| **value_tier** | **moderado** |

Comparado à iteração anterior: delta passou de 0.208 → 0.347 (`sem mudança`). (Nota: iteração 4 rodou 4 evals; iteração 5 rodou 6 evals — conjunto diferente, comparação direta é apenas indicativa.)

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.163944 | $0.027324 | claude-opus-4-8 |
| without_skill | $0.107670 | $0.017945 | claude-opus-4-8 |
| **custo adicional da skill** | **$0.056274** | — | — |

(Nota: claude-opus-4-8 não consta na tabela de preços; preços de sonnet aplicados como fallback, split de tokens estimado 75/25.)

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-scope-elicitation-produce-scope-summary-from-text | 1.000 | 0.750 | +0.250 |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file | 1.000 | 0.500 | +0.500 |
| eval-no-code-generation-skill-must-not-produce-go-source | 1.000 | 0.667 | +0.333 |
| eval-phase-gate-enforcement-phase-2-must-not-start | 1.000 | 1.000 | 0.000 |
| eval-dependency-analysis-tasks-declare-depends-on-and | 1.000 | 0.800 | +0.200 |
| eval-complexity-evaluation-tasks-scored-on-five-dimensions | 1.000 | 0.200 | +0.800 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-complexity-evaluation-tasks-scored-on-five-dimensions | +0.800 | Sem a skill o modelo usa dimensões erradas (Domain/Integration/Risk/Testing/Effort), trata 'score' como soma em vez de média e omite os campos 'tier' e 'risk_note'. A skill impõe as cinco dimensões corretas (scope/ambiguity/coupling/novelty/reversibility), a média e o threshold de tier. |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file | +0.500 | Sem a skill os IDs de task usam camadas por extenso (TASK-DOMAIN-, TASK-APPLICATION-, TASK-INFRASTRUCTURE-) em vez de DOM/APP/INFRA, e os títulos de task ficam em heading '### TASK-...' sem checkbox. A skill garante o padrão de ID semântico e o formato de checkbox `[ ]`. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-phase-gate-enforcement-phase-2-must-not-start` (with_skill 1.000 / without_skill 1.000)

## Lacunas da skill (with_skill < 1.0)

Nenhuma — with_skill atingiu 1.000 em todos os evals.

## Recomendação

Pronta — a skill atinge 1.000 em todos os 6 evals e entrega delta moderado (+0.347) sobre o baseline.
