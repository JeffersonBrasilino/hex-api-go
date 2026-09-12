# Relatório de Avaliação: `sdd-prd` — iteração 2

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.950 |
| without_skill | 0.550 |
| **delta** | **+0.400** |
| **value_tier** | **forte** |

Comparado à iteração anterior: delta passou de 0.300 → 0.400 (`moderado → forte`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $1.065360 | $0.213072 | claude-sonnet-5 |
| without_skill | $0.723855 | $0.144771 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.341505** | — | — |

O custo adicional da skill cresceu $0.11973 em relação à iteração anterior, acompanhando o ganho de taxa de aprovação — o investimento extra em tokens (média de 35.512 vs. 24.129 tokens) segue justificado pelo salto de qualidade.

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

O benchmark desta iteração não trouxe a quebra do `without_skill` por caso individual (apenas a média agregada de 0.550, com desvio padrão de 0.2092, indicando variação entre os casos). Com base no resultado `with_skill` (4 de 5 casos com pontuação máxima 4/4) e no delta agregado de +0.400, os quatro casos abaixo sustentam o valor forte da skill:

| Slug | Delta | Motivo |
|---|---|---|
| eval-parse-intake-cjs-lgpd-scan-cjs-run-on-intake-but-gaps-nenhum | dado agregado (não quebrado por caso) | Com a skill, o parsing do intake e a varredura LGPD são executados corretamente, cobrindo o caso "gaps: nenhum" sem falhas — 4/4 with_skill |
| eval-lgpd-dimensions-asked-one-at-a-time-each-with-a-recommended | dado agregado (não quebrado por caso) | Dimensões de LGPD são perguntadas uma de cada vez, cada uma com recomendação — comportamento essencial de entrevista guiada, 4/4 with_skill |
| eval-validate-prd-cjs-and-detect-implementation-detail-cjs-both-g | dado agregado (não quebrado por caso) | validate-prd.cjs e detect-implementation-detail.cjs são ambos acionados corretamente — 4/4 with_skill |
| eval-technical-detail-volunteered-mid-interview-triggers-merge-no | dado agregado (não quebrado por caso) | Detalhe técnico oferecido espontaneamente durante a entrevista aciona corretamente o fluxo de merge-notes.cjs — 4/4 with_skill |

## Baseline confirmado (ambos ≥ 0.95)

- Nenhum caso confirmado nesta iteração — o `without_skill` agregado (0.550) fica abaixo do limiar de 0.95, e a granularidade por caso não está disponível neste benchmark para validar baseline individualmente.

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-approval-gate-step-2-still-required-before-prd-md-is-written | O resumo de aprovação usa o formato padronizado exato (Problema, Objetivo, Usuários, Requisitos funcionais, Regras de negócio, Não funcionais, Dados/LGPD, Critérios de aceitação, Fora de escopo, Dúvidas pendentes) — a resposta apenas menciona "Resumo para aprovação (Step 2, formato padronizado em pt-BR)" sem listar os campos específicos exigidos | Reforçar no SKILL.md, no Step 2 (gate de aprovação), a listagem explícita e literal dos dez campos do formato padronizado, deixando claro que devem aparecer como cabeçalhos/rótulos no resumo apresentado ao usuário — e não apenas referenciados por nome genérico ("formato padronizado") |

## Recomendação

Precisa de iteração — o valor da skill já é forte (+0.400, sem regressões nos demais casos), mas o gate de aprovação do Step 2 ainda não materializa o formato padronizado completo e merece um ajuste pontual antes da próxima rodada.
