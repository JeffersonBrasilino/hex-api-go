# Relatório de Avaliação: `sdd-workflow-script` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.550 |
| **delta** | **+0.450** |
| **value_tier** | **forte** |

Esta é a iteração 1 do skill, portanto não há iteração anterior para comparação (`delta_vs_prev_iteration` indisponível).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $2.293725 | $0.458745 | claude-sonnet-5 |
| without_skill | $1.559025 | $0.311805 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.734700** | — | — |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-wave-computed-via-compute-wave-cjs-from-the-persisted-wave-f | 1.00 | Sem a skill (0/4), o modelo usou "lógica equivalente" em vez de invocar `compute-wave.cjs`, fez inferência manual sobre os campos `depends_on`/`wave`, não tratou deadlock e omitiu `TOTAL_WAVES`. Com a skill (4/4), o script foi invocado corretamente em todos os casos. |
| eval-bootstrap-detect-state-cjs-write-state-cjs-initialize-state | 0.50 | Sem a skill (2/4), o conteúdo do STATE.md foi narrado em prosa em vez de mostrar a chamada exata ao script `write-state.cjs`. Com a skill (4/4), a sintaxe de invocação explícita foi usada de forma consistente. |
| eval-retry-counter-updated-via-write-state-cjs-increment-retry-ne | 0.50 | Sem a skill (2/4), foi descrita uma edição manual de estrutura YAML em vez da sintaxe exata `write-state.cjs --increment-retry`, e não houve menção explícita de aguardar ambos os agentes antes da próxima tarefa. Com a skill (4/4), o comportamento correto foi seguido integralmente. |

**Nota:** eval-plan-md-checkboxes-marked-via-update-plan-cjs-never-hand-edi teve delta de 0.25 (with_skill 4/4=1.00 vs without_skill 3/4=0.75), abaixo do limiar de 0.40 e por isso não incluído na tabela acima. Sem a skill, o modelo assumiu que a próxima wave seria `current_wave + 1` por aritmética, em vez de reexecutar `compute-wave.cjs` para determiná-la — um erro discreto que a skill corrige.

## Baseline confirmado (ambos ≥ 0.95)

| Slug | with_skill | without_skill | Delta |
|---|---|---|---|
| eval-wave-confirmation-gate-still-enforced-before-spawning-subage | 1.00 | 1.00 | 0.00 |

Ambas as configurações já sabem pedir confirmação antes de disparar subagentes — trata-se de conhecimento de base amplamente presente no modelo, não de um ganho atribuível à skill.

## Lacunas da skill (with_skill < 1.0)

Nenhuma. A configuração with_skill obteve 1.00 (4/4) em todas as 5 avaliações desta iteração — não foram identificadas lacunas.

## Recomendação

A skill está pronta, dado o índice de aprovação de 100% em with_skill e o delta médio forte de +0.45 frente ao baseline; o custo incremental (~$0.73 nas 5 avaliações, cerca de 2x o baseline puro) é esperado pela verbosidade das invocações de script e plenamente justificado pelos ganhos de correção, em especial pela eliminação total das falhas no cenário de `compute-wave.cjs`.
