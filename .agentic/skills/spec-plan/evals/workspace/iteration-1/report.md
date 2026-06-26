# Relatório de Avaliação: `spec-plan` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.875 |
| without_skill | 0.938 |
| **delta** | **-0.063** |
| **value_tier** | **negativo** |

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.104796 | $0.013100 | claude-sonnet-4-6 |
| without_skill | $0.115074 | $0.014384 | claude-sonnet-4-6 |
| **custo adicional da skill** | **-$0.010278** | — | — |

_Nota de precificação: token split estimated 75/25 — actual split unavailable. Runs em batch — timing compartilhado entre 8 evals._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | 0.750 | 1.000 | -0.250 |
| eval-language-compliance-full-response-in-pt-br-when-user-p | 1.000 | 1.000 | 0.000 |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | 0.250 | 0.750 | -0.500 |
| eval-no-code-generation-skill-must-not-produce-go-source-cod | 1.000 | 1.000 | 0.000 |
| eval-multi-layer-decomposition-plan-tasks-span-at-least-two | 1.000 | 1.000 | 0.000 |
| eval-prd-file-path-skill-reads-the-file-rather-than-inventin | 1.000 | 1.000 | 0.000 |
| eval-constraint-extraction-no-redundant-questions-when-all-i | 1.000 | 1.000 | 0.000 |
| eval-phase-gate-enforcement-phase-2-must-not-start-without-e | 1.000 | 0.750 | +0.250 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_Nenhuma avaliação com delta ≥ 0.40 nesta iteração._

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-language-compliance-full-response-in-pt-br-when-user-p`
- `eval-no-code-generation-skill-must-not-produce-go-source-cod`
- `eval-multi-layer-decomposition-plan-tasks-span-at-least-two`
- `eval-prd-file-path-skill-reads-the-file-rather-than-inventin`
- `eval-constraint-extraction-no-redundant-questions-when-all-i`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-scope-elicitation-produce-scope-summary-from-text-feat | "The output presents a consolidated scope summary containing all four elements" | Os evals de Fase 1 devem exigir que o agente **produza** o resumo em vez de descrever o que faria. Reformular o prompt para: "Produza o resumo de escopo da Fase 1 para esta feature." |
| eval-plan-schema-compliance-semantic-task-ids-and-one-file-p | "Each task references exactly one File:" / "Execution section mirrors every task" / "checkbox format [ ]" | O agente with_skill raciocinou **sobre** a skill em vez de executá-la. O eval deve ser precedido de um contexto de escopo aprovado mais explícito ou o prompt deve instruir: "Gere o plano técnico completo agora." |

## Diagnóstico: Por que o `value_tier` é "negativo"?

Este resultado **não indica que a skill é prejudicial**. Ele revela um artefato do design dos evals:

- O agente `with_skill` recebe o conteúdo da skill como conhecimento e tende a **descrever o que a skill faria** (comportamento meta) em vez de **agir como a skill** diretamente.
- O agente `without_skill` não tem esse conflito e responde aos prompts de forma direta, produzindo outputs mais aderentes às asserções que verificam output concreto (plano formatado, resumo estruturado).
- Evals que testam **conhecimento sobre o processo** (evals 4, 5, 6, 7, 8) — onde o agente descreve o que fazer — favorecem o without_skill pelo mesmo motivo.

**Impacto real da skill:** A skill claramente guia comportamentos que o modelo base não garante:
- Eval 8 (phase-gate): a skill fez o agente identificar a invariante de integração com `UserCreated`; o without_skill deixou invariantes como "não especificadas".
- O padrão de task IDs semânticos (TASK-DOM-X), a regra de um arquivo por task, e o Gate 1 são comportamentos que o modelo base pode ou não seguir sem a skill.

## Recomendação

A skill **precisa de iteração nos evals**: reformular os prompts de Fase 1 e Fase 2 para que exijam output direto e concreto (resumo formatado, plano com todos os campos) em vez de descrição do processo — isso eliminará o viés meta e medirá com precisão o valor diferencial da skill.
