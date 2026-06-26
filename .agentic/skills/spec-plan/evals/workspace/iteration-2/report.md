# Relatório de Avaliação: `spec-plan` — iteração 2

## Pontuação Geral

| Configuração    | Taxa de Aprovação Média |
|-----------------|-------------------------|
| with_skill      | 0.833                   |
| without_skill   | 0.927                   |
| **delta**       | **-0.094**              |
| **value_tier**  | **negativo**            |

Comparado à iteração anterior: delta passou de -0.063 → -0.094 (`sem mudança` — tier mantido em `negativo`, mas delta piorou em -0.032).

## Custo estimado

| Configuração              | Total (USD)  | Média por avaliação (USD) | Modelo             |
|---------------------------|--------------|---------------------------|--------------------|
| with_skill                | $0.131298    | $0.016412                 | claude-sonnet-4-6  |
| without_skill             | N/D          | N/D                       | claude-sonnet-4-6  |
| **custo adicional da skill** | N/D       | —                         | —                  |

_Nota de precificação: token split estimated 75/25 — actual split unavailable. Custo without_skill indisponível — outputs eram pré-existentes de execução incompleta anterior._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|------|------------|---------------|-------|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | 0.750 | 0.750 | 0.000 |
| eval-language-compliance-full-response-in-pt-br-when-user-p | 1.000 | 1.000 | 0.000 |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | 0.750 | 1.000 | -0.250 |
| eval-no-code-generation-skill-must-not-produce-go-source-cod | 0.667 | 0.667 | 0.000 |
| eval-multi-layer-decomposition-plan-tasks-span-at-least-two | 1.000 | 1.000 | 0.000 |
| eval-prd-file-path-skill-reads-the-file-rather-than-inventin | 1.000 | 1.000 | 0.000 |
| eval-constraint-extraction-no-redundant-questions-when-all-i | 1.000 | 1.000 | 0.000 |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | 0.500 | 1.000 | -0.500 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_(Nenhuma avaliação — nenhum delta ≥ 0.40 registrado nesta iteração)_

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-language-compliance-full-response-in-pt-br-when-user-p` (4/4 ambos)
- `eval-multi-layer-decomposition-plan-tasks-span-at-least-two` (4/4 ambos)
- `eval-prd-file-path-skill-reads-the-file-rather-than-inventin` (4/4 ambos)
- `eval-constraint-extraction-no-redundant-questions-when-all-i` (4/4 ambos)

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|------|----------------|-------------------|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | "output explicitly asks the user for confirmation or approval (Gate 1)" | Adicionar instrução explícita na skill: após apresentar o resumo de escopo, sempre incluir frase de Gate 1 como "Você confirma este escopo para avançarmos para a Fase 2?" |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | "Each task references exactly one file in the 'File:' field" | Adicionar restrição explícita na skill: command.go e handler.go de um mesmo caso de uso devem ser tasks separadas (TASK-APP-[CONCERN]-COMMAND e TASK-APP-[CONCERN]-HANDLER). |
| eval-no-code-generation-skill-must-not-produce-go-source-cod | "The response addresses Phase 1 scope elicitation before jumping to create a plan" | A skill deve explicitar que mesmo para prompts que pedem plano diretamente, a Fase 1 é obrigatória antes de qualquer geração. |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | "response produces a Phase 1 scope summary and explicitly waits for user confirmation" | Esclarecer que quando há informações insuficientes, o agente deve apresentar um resumo parcial + perguntas, e não apenas as perguntas. O resumo de escopo (mesmo incompleto) sempre deve aparecer antes das perguntas. |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | "scope summary identifies at least one invariant or constraint relevant to the notification domain" | Mesmo em resumos parciais, o agente deve listar invariantes que puder inferir do contexto (ex.: "integração com evento UserCreated — mecanismo a confirmar"). |

## Recomendação

A skill **precisa de iteração** — o `value_tier` é `negativo` pela segunda iteração consecutiva, com o delta piorando de -0.063 para -0.094. Os dois problemas críticos são: (1) o agente with_skill omite o resumo de escopo formatado quando há perguntas de elicitação, e (2) viola a regra de um-arquivo-por-task ao agrupar command.go e handler.go na mesma task. Ambos são corrigíveis com instruções mais precisas na SKILL.md.
