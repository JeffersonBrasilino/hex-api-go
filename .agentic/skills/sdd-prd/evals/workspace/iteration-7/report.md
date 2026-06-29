# Relatório de Avaliação: `spec-prd-v2` — iteração 7 (parcial: IDs 2, 5, 7)

> **Nota:** iteração parcial — apenas 3 dos 16 evals foram executados, correspondendo às lacunas corrigidas na iteração 6.
> Comparado à iteração anterior: delta de `pass_rate` passou de **-0.027 → +0.556** (`negativo → forte`).
> Comparação aproximada — escopos diferentes (3 vs 16 evals).

---

## Pontuação Geral

| Métrica | with_skill | without_skill | delta | value_tier |
|---|---|---|---|---|
| pass_rate (média) | 1.000 | 0.444 | **+0.556** | **forte** |
| time_seconds (média) | 22.236 s | 7.141 s | +15.095 s | — |
| tokens (média) | 5 188 | 3 234 | +1 954 | — |

---

## Custo Estimado

| | with_skill | without_skill | delta |
|---|---|---|---|
| Total (USD) | $0.093390 | $0.058212 | +$0.035178 |
| Média por eval (USD) | $0.031130 | $0.019404 | +$0.011726 |

---

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| `eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez` | 1.000 (5/5) | 0.600 (3/5) | **+0.400** |
| `eval-resumo-de-aprovacao-segue-o-formato-padronizado` | 1.000 (5/5) | 0.400 (2/5) | **+0.600** |
| `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo` | 1.000 (6/6) | 0.333 (2/6) | **+0.667** |

---

## Skill agrega valor claro (delta >= 0.40, ordem decrescente)

| Slug | delta | Motivo |
|---|---|---|
| `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo` | +0.667 | Skill instrui seções RF-0X/RNF-0X com hífen, PRD inline e confirmação de caminho |
| `eval-resumo-de-aprovacao-segue-o-formato-padronizado` | +0.600 | Skill instrui preenchimento do bloco completo; sem skill o formato padronizado é ignorado |
| `eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez` | +0.400 | Skill instrui recomendação antes da pergunta e um único `?`; sem skill não há valor padrão |

---

## Baseline confirmado (ambos >= 0.95)

Nenhum eval atingiu baseline confirmado: `with_skill` alcançou 1.000 nos três casos, mas `without_skill` permaneceu abaixo de 0.95 nos três.

---

## Lacunas da skill (with_skill < 1.0)

Nenhuma — todos os evals passaram com 100% de aprovação no cenário `with_skill`.

---

## Recomendação

A skill demonstrou valor forte e consistente nos 3 evals parciais, revertendo o resultado negativo da iteração 6; recomenda-se executar o conjunto completo dos 16 evals para confirmar a melhoria em toda a suíte antes de considerar a skill estável.
