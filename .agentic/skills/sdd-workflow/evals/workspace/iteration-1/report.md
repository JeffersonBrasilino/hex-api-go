# Relatório de Avaliação: `sdd-workflow` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.563 |
| **delta** | **+0.438** |
| **value_tier** | **forte** |

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.092688 | $0.023172 | claude-sonnet-4-6 |
| without_skill | $0.076476 | $0.019119 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.016212** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-bootstrap-initialize-pipeline-when-state-md-does-not-ex | 1.000 | 1.000 | 0.000 |
| eval-phase-routing-prd-plan-advance-when-prd-md-is-detected | 1.000 | 0.250 | +0.750 |
| eval-wave-confirmation-gate-must-wait-for-user-approval-befo | 1.000 | 0.250 | +0.750 |
| eval-done-phase-show-completion-summary-with-artifact-paths | 1.000 | 0.750 | +0.250 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-phase-routing-prd-plan-advance-when-prd-md-is-detected | +0.750 | Sem a skill, o modelo não produziu a saída explícita de avanço de fase, não exibiu o caminho do PRD.md detectado e não instruiu o usuário a executar `/spec-plan` com o caminho correto — apenas descreveu o comportamento esperado em vez de executá-lo. |
| eval-wave-confirmation-gate-must-wait-for-user-approval-befo | +0.750 | Sem a skill, o modelo não exibiu as tasks agrupadas por `parallel_group`, não apresentou o cabeçalho de wave no formato esperado e não solicitou confirmação explícita ao usuário com o termo "sim" — descreveu o fluxo conceitualmente sem produzir a saída formatada. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-bootstrap-initialize-pipeline-when-state-md-does-not-ex`

## Lacunas da skill (with_skill < 1.0)

Nenhuma — todas as avaliações atingiram 1.000 com a skill ativa.

## Recomendação

A skill `sdd-workflow` está pronta para uso: atingiu taxa de aprovação perfeita (1.000) em todas as avaliações, com delta forte de +0.438 em relação ao baseline sem skill, demonstrando que o orquestrador produz saídas formatadas, roteamento de fases e gates de confirmação corretamente sem depender de inferência livre do modelo.
