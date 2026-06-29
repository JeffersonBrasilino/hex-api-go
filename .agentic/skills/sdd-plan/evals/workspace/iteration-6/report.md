# Relatório de Avaliação: `sdd-plan` — iteração 6

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.503 |
| **delta** | **+0.497** |
| **value_tier** | **forte** |

Comparado à iteração anterior (iteração 5): delta passou de 0.347 → 0.497 (`moderado → forte`). A troca de modelo (claude-opus-4-8 → claude-sonnet-4-6) não prejudicou o desempenho da skill; pelo contrário, a diferença entre com e sem skill aumentou 0.150 pontos, confirmando que a skill agrega valor independente de modelo.

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.288237 | $0.048040 | claude-sonnet-4-6 |
| without_skill | $0.223581 | $0.037264 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.064656** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-scope-elicitation-scope-summary | 1.000 | 0.750 | +0.250 |
| eval-plan-schema-semantic-task-ids | 1.000 | 0.750 | +0.250 |
| eval-no-code-generation | 1.000 | 0.667 | +0.333 |
| eval-phase-gate-enforcement | 1.000 | 0.250 | +0.750 |
| eval-dependency-analysis-depends-on | 1.000 | 0.600 | +0.400 |
| eval-complexity-evaluation-five-dimensions | 1.000 | 0.000 | +1.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-complexity-evaluation-five-dimensions | +1.000 | Sem a skill, o modelo usou dimensões incorretas sem pontuações numéricas. Com a skill, todas as 5 dimensões obrigatórias foram aplicadas corretamente |
| eval-phase-gate-enforcement | +0.750 | Sem a skill, o modelo desviava para `/sdd-prd` em vez de executar a Fase 1 de elicitação. Com a skill, o gate é respeitado em 100% dos casos |
| eval-dependency-analysis-depends-on | +0.400 | Sem a skill, os campos `depends_on` e `parallel_group` eram omitidos. Com a skill, a análise de dependências está completa e consistente |

## Baseline confirmado (ambos ≥ 0.95)

Nenhum eval atingiu o limiar de baseline confirmado (≥ 0.95) sem a skill. O without_skill máximo foi 0.750 (dois evals). Em todos os casos o baseline sem skill ficou abaixo de 0.80.

## Lacunas da skill (with_skill < 1.0)

Nenhuma lacuna identificada. Todos os 6 evals atingiram pass_rate = 1.000 com a skill ativa nesta iteração.

## Recomendação

A skill `sdd-plan` atingiu o tier **forte** nesta iteração, avançando de `moderado` (iteração 5, delta 0.347) para `forte` (iteração 6, delta 0.497). Os três maiores ganhos estão em áreas estruturais críticas: avaliação de complexidade com dimensões padronizadas, enforcement de gates de fase e análise de dependências com campos obrigatórios. Recomenda-se manter a skill sem alterações e monitorar se o tier `forte` se sustenta em iterações futuras com variação de contexto de entrada.
