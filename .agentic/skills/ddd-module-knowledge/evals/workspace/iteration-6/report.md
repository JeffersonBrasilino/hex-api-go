# Relatório de Avaliação: `ddd-module-knowledge` — iteração 6

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.900 |
| **delta** | **+0.100** |
| **value_tier** | **fraco** |

_(Comparado à iteração anterior: delta passou de 0.211 → 0.100 (moderado → fraco).)_

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $2.088150 | $0.417630 | claude-sonnet-5 |
| without_skill | $1.561125 | $0.312225 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.527025** | — | — |

> ⚠️ **Nota sobre artefato de medição de custo**: nesta iteração as 5 perguntas foram enviadas em uma
> única chamada de subagente (batch), diferente de iterações anteriores, que fizeram uma chamada por
> eval. O script `write-timing.cjs` grava o mesmo `timing.json` (27.842 tokens / $0.417630 para
> with_skill; 20.815 tokens / $0.312225 para without_skill) em cada um dos 5 diretórios de eval, pois
> o custo real da chamada em lote é compartilhado entre todas as perguntas respondidas nela.
> `compute-benchmark.cjs` soma esse valor por eval, contando o custo da mesma chamada 5 vezes:
> $0.417630 × 5 = $2.088150 (with_skill) e $0.312225 × 5 = $1.561125 (without_skill). O custo real
> pago nesta iteração foi de ~$0.42 (with_skill) + ~$0.31 (without_skill) ≈ $0.73 no total, não
> $3.65. Consequentemente, o "custo adicional da skill" de $0.527025 também está inflado em ~5x — o
> valor real fica próximo de $0.105. Esse artefato é do modo de execução em lote desta iteração
> específica, não uma regressão real de custo da skill. Iterações futuras devem rodar uma chamada de
> subagente por eval para manter os números de custo comparáveis entre iterações.

## Skill agrega valor claro (delta >= 0.40, ordem decrescente)

Nenhum eval individual atingiu delta >= 0.40 nesta iteração.

## Baseline confirmado (ambos >= 0.95)

- eval-error-handling-ddgo-error-selection-across-three-scenarios
- eval-gotchas-aggregate-root-naming-child-entity-construction-and
- eval-layer-boundaries-identify-violations-across-http-application
- eval-reference-resolution-by-file-path-and-skill-scope-boundary

(exceto o eval de module-scaffold, onde without_skill ficou em 4/6 = 0.667)

## Lacunas da skill (with_skill < 1.0)

Nenhuma — with_skill atingiu 100% em todos os 5 evals.

## Recomendação

Precisa de iteração: o value_tier caiu de "moderado" para "fraco" e o delta de pass_rate despencou de
0.211 para 0.100 em relação à iteração anterior — não por regressão da skill (with_skill manteve
100%), mas porque o baseline sem a skill (without_skill) ficou mais forte nesta rodada, cobrindo por
conhecimento geral a maior parte das convenções DDD/hexagonal do projeto. O ganho real da skill hoje
está concentrado em detalhes específicos e estreitos: nomes exatos de struct/constructor
(`GormInvoiceRepository`), o prefixo `hex-api-go.` em `TableName()`, e o redirecionamento explícito
para a skill `make-unit-tests`. A comparação de custo desta iteração não é confiável (ver nota acima)
— refazer a medição com uma chamada por eval antes de decidir se o custo real ainda compensa o valor
entregue.
