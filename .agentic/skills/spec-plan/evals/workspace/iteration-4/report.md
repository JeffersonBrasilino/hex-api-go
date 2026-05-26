# Relatório de Avaliação: `spec-plan` — iteração 4

## Pontuação Geral

| Configuração    | Taxa de Aprovação Média |
|-----------------|-------------------------|
| with_skill      | 1.000                   |
| without_skill   | 0.792                   |
| **delta**       | **+0.208**              |
| **value_tier**  | **moderado**            |

Comparado à iteração anterior: delta passou de +0.073 → +0.208 (`fraco → moderado` — melhora de +0.135 — primeira iteração com `with_skill` perfeito em todos os evals).

## Custo estimado

| Configuração              | Total (USD)  | Média por avaliação (USD) | Modelo             |
|---------------------------|--------------|---------------------------|--------------------|
| with_skill                | $0.104116    | $0.026029                 | claude-sonnet-4-6  |
| without_skill             | $0.060820    | $0.015205                 | claude-sonnet-4-6  |
| **custo adicional da skill** | **$0.043296** | —                      | —                  |

_Nota de precificação: token split estimated 75/25 — actual split unavailable. Runs executados em batch (4 evals por subagente). Custo médio por eval maior que na iteração 3 pois esta iteração rodou 4 evals (vs 8 anteriores), diluindo menos o custo fixo do batch._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|------|------------|---------------|-------|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | 1.000 | 0.750 | +0.250 |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | 1.000 | 0.750 | +0.250 |
| eval-no-code-generation-skill-must-not-produce-go-source-cod | 1.000 | 0.667 | +0.333 |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_(Nenhuma avaliação — nenhum delta ≥ 0.40 registrado nesta iteração)_

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-phase-gate-enforcement-phase-2-must-not-start-without-e` (4/4 ambos)

## Lacunas da skill (with_skill < 1.0)

_(Nenhuma — pontuação perfeita em todos os 4 evals ativos)_

## Wins da skill (delta > 0)

| Slug | Delta | Motivo |
|------|-------|--------|
| eval-no-code-generation-skill-must-not-produce-go-source-cod | +0.333 | with_skill executa Fase 1 obrigatória antes de qualquer plano; without_skill pula direto para o plano por camadas sem elicitação de escopo. |
| eval-scope-elicitation-produce-scope-summary-from-text-feat | +0.250 | Instrução CRITICAL adicionada na iteração 4: with_skill encerra com frase canônica de Gate 1 após perguntas de elicitação; without_skill termina na última pergunta sem frase de confirmação. |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | +0.250 | with_skill usa `- [ ] **TASK-...**` para todos os task entries e espelha na Seção 3 com `- [ ] **Execution — TASK-...**`; without_skill usa `### TASK-...` (heading) e tabela markdown na Seção 3. |

## Recomendação

A skill **atingiu maturidade operacional** nesta iteração — `value_tier` passou de `fraco` para `moderado` com `with_skill` perfeito (1.000) pela primeira vez. A correção do Gate 1 resolveu a última lacuna restante da iteração 3. Os próximos passos recomendados são: (1) reativar evals desabilitados com variações de prompt mais desafiadoras para aumentar o sinal diferencial, ou (2) adicionar novos evals que testem cenários de borda ainda não cobertos (ex: PRD com conflitos entre pré/pós-condições, feature com impacto em múltiplos módulos).
