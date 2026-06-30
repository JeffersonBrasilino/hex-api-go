# Relatório de Avaliação: `sdd-prd` — iteração 9

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.788 |
| **delta** | **+0.213** |
| **value_tier** | **moderado** |

Comparado à iteração anterior (8): delta passou de 0.025 → 0.213 (`sem_valor → moderado`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.098088 | $0.024522 | claude-sonnet-4-6 |
| without_skill | $0.072462 | $0.018115 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.025626** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade-e-na | 1.000 | 0.800 | +0.200 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 1.000 | 0.750 | +0.250 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pr | 1.000 | 0.600 | +0.400 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo-fo | 1.000 | 1.000 | +0.000 |

## Skill agrega valor claro (delta >= 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pr | +0.400 | Sem a skill, o modelo pede nome de modulo/dominio antes de entender o problema e faz multiplas perguntas simultaneas; com a skill, faz exatamente uma unica pergunta de contexto. |

## Baseline confirmado (ambos >= 0.95)

- `eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo-fo`

## Lacunas da skill (with_skill < 1.0)

Nenhuma lacuna identificada — todas as avaliacoes atingiram pass_rate 1.0 com a skill ativa.

## Recomendacao

A skill esta pronta para uso em producao: atingiu 100% nas quatro avaliacoes e recuperou a classificacao "moderado" apos a regressao da iteracao anterior; recomenda-se uma proxima iteracao focada em elevar o baseline do modelo sem a skill nos cenarios de pergunta unica para consolidar a margem de valor.
