# Relatório de Avaliação: `sdd-plan-script` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.950 |
| without_skill | 0.650 |
| **delta** | **+0.300** |
| **value_tier** | **moderado** |

(N=1, então não há comparação com iteração anterior — omitido)

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.157170 | $0.031434 | claude-sonnet-5 |
| without_skill | $0.132240 | $0.026448 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.024930** | — | — |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-wave-computed-by-compute-waves-cjs-never-assigned-by-hand | +0.50 | with_skill 0.75 vs without_skill 0.25 — a skill reduz drasticamente a tendência do agente de calcular waves manualmente em vez de delegar ao script |
| eval-gate-2-handoff-table-sourced-verbatim-from-generate-task-ind | +0.50 | with_skill 1.0 vs without_skill 0.5 — a skill garante que a tabela de handoff do Gate 2 seja extraída literalmente do `generate-task-index`, em vez de reconstruída pelo agente |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-complexity-delegated-to-calculate-complexity-cjs-never-hand`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-wave-computed-by-compute-waves-cjs-never-assigned-by-hand | A saída não deve se apresentar como tendo calculado manualmente os números de wave — isso deve ser inteiramente delegado à regra mecânica do script (piso da camada + max(wave dos depends_on) + 1). Na evidência, o agente narra as atribuições específicas de wave ("TASK-DOM-SESSION-REVOKE ... cai na wave 1" etc.) explicando a regra, o que dá a impressão de que o cálculo foi compreendido/verificado pelo próprio agente, e não puramente reportado a partir da saída do script. | Ajustar as instruções da skill para, após rodar `compute-waves.cjs`, o agente apenas reportar o resultado tabular do script (task → wave) sem reexplicar ou justificar a regra de cálculo em prosa; adicionar um exemplo explícito de formato de saída permitido ("apresente a tabela de waves gerada pelo script, sem comentar a lógica de atribuição") para reforçar que a narrativa de "por que" cada task caiu em determinada wave é proibida. |

## Recomendação

Precisa de iteração — o delta geral é moderado e positivo, mas a lacuna recorrente de "explicar" o cálculo de waves em vez de apenas reportar a saída do script deve ser corrigida no texto da skill antes de considerá-la pronta.
