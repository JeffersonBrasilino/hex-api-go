# Relatório de Avaliação: `spec-plan` — iteração 3

## Pontuação Geral

| Configuração    | Taxa de Aprovação Média |
|-----------------|-------------------------|
| with_skill      | 0.969                   |
| without_skill   | 0.896                   |
| **delta**       | **+0.073**              |
| **value_tier**  | **fraco**               |

Comparado à iteração anterior: delta passou de -0.094 → +0.073 (`negativo → fraco` — melhora de +0.167 — primeira iteração positiva da skill).

## Custo estimado

| Configuração              | Total (USD)  | Média por avaliação (USD) | Modelo             |
|---------------------------|--------------|---------------------------|--------------------|
| with_skill                | $0.110880    | $0.013860                 | claude-sonnet-4-6  |
| without_skill             | $0.093045    | $0.011631                 | claude-sonnet-4-6  |
| **custo adicional da skill** | **$0.017835** | —                      | —                  |

_Nota de precificação: token split estimated 75/25 — actual split unavailable. Runs executados em batch (8 evals por subagente)._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|------|------------|---------------|-------|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | 0.750 | 0.750 | 0.000 |
| eval-language-compliance-full-response-in-pt-br-when-user-p | 1.000 | 1.000 | 0.000 |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | 1.000 | 0.750 | +0.250 |
| eval-no-code-generation-skill-must-not-produce-go-source-cod | 1.000 | 0.667 | +0.333 |
| eval-multi-layer-decomposition-plan-tasks-span-at-least-two | 1.000 | 1.000 | 0.000 |
| eval-prd-file-path-skill-reads-the-file-rather-than-inventin | 1.000 | 1.000 | 0.000 |
| eval-constraint-extraction-no-redundant-questions-when-all-i | 1.000 | 1.000 | 0.000 |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_(Nenhuma avaliação — nenhum delta ≥ 0.40 registrado nesta iteração)_

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-language-compliance-full-response-in-pt-br-when-user-p` (4/4 ambos)
- `eval-multi-layer-decomposition-plan-tasks-span-at-least-two` (4/4 ambos)
- `eval-prd-file-path-skill-reads-the-file-rather-than-inventin` (4/4 ambos)
- `eval-constraint-extraction-no-redundant-questions-when-all-i` (4/4 ambos)
- `eval-phase-gate-enforcement-phase-2-must-not-start-without-e` (4/4 ambos — **novo** em relação à iteração 2, onde with_skill era 0.500)

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|------|----------------|-------------------|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | "output explicitly asks the user for confirmation or approval (Gate 1)" | A instrução de Gate 1 já existe na SKILL.md, mas o agente termina com perguntas de elicitação numeradas sem fechar com a frase canônica. Adicionar instrução mais precisa: "Mesmo quando há perguntas de elicitação, o Gate 1 SEMPRE deve encerrar a mensagem com frase separada de confirmação após as perguntas, ex: 'Assim que responder, confirmo o escopo e avançamos para a Fase 2.'" |

## Wins da skill (nova seção — iteração positiva)

| Slug | Delta | Motivo |
|------|-------|--------|
| eval-no-code-generation-skill-must-not-produce-go-source-cod | +0.333 | without_skill pula Fase 1 e gera plano direto; with_skill recusa e redireciona ao processo correto da skill. |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | +0.250 | without_skill usa headings '### TASK-...' sem '- [ ]' e tabela em Seção 3; with_skill usa formato correto de checkbox em todos os níveis e separa command.go e handler.go em tasks distintas. |

## Recomendação

A skill **está progredindo** — `value_tier` passou de `negativo` para `fraco` pela primeira vez, com delta +0.073. Os dois problemas críticos da iteração 2 (violação de um-arquivo-por-task e ausência de resumo parcial antes de perguntas) foram resolvidos. A única lacuna remanescente é o Gate 1 explícito em `eval-scope-elicitation`: o agente produz o resumo de escopo corretamente mas não fecha com a frase de confirmação quando há perguntas de elicitação pendentes. Uma instrução precisa adicional nesse ponto pode elevar a pontuação para `moderado` na próxima iteração.
