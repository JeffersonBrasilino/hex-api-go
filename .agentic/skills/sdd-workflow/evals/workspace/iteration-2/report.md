# Relatório de Avaliação: `sdd-workflow` — iteração 2

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.500 |
| **delta** | **+0.500** |
| **value_tier** | **forte** |

Comparado à iteração anterior (iteração 1): delta passou de 0.4375 → 0.500 (`forte → forte`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.092307 | $0.023077 | claude-sonnet-4-6 |
| without_skill | $0.082584 | $0.020646 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.009723** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-bootstrap-initialize-pipeline | 1.000 | 0.750 | +0.250 |
| eval-phase-routing-prd-to-plan | 1.000 | 0.250 | +0.750 |
| eval-wave-confirmation-gate | 1.000 | 0.250 | +0.750 |
| eval-done-phase-completion-summary | 1.000 | 0.750 | +0.250 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-phase-routing-prd-to-plan | +0.750 | Sem a skill, o modelo descreve o processo abstratamente sem produzir saída voltada ao usuário com mensagens e formatação esperadas |
| eval-wave-confirmation-gate | +0.750 | Sem a skill, o modelo descreve o conceito de gate sem mostrar o formato wave header ou solicitar confirmação explícita ao usuário |

## Baseline confirmado (ambos ≥ 0.95)

Nenhum eval atingiu o limiar de baseline confirmado (≥ 0.95) sem a skill. O without_skill máximo foi 0.750 (dois evals).

## Lacunas da skill (with_skill < 1.0)

Nenhuma lacuna identificada. Todos os 4 evals atingiram pass_rate = 1.000 com a skill ativa nesta iteração.

## Recomendação

A skill `sdd-workflow` mantém o tier **forte** com melhora de +0.0625 no delta (iteração 1: 0.4375 → iteração 2: 0.500). Os padrões de falha sem a skill são sistemáticos: (1) idioma errado — sem a skill, o modelo responde em inglês em vez de pt-BR; (2) abstração sem execução — sem a skill, o modelo descreve comportamentos do pipeline em nível conceitual sem produzir as saídas concretas voltadas ao usuário. Recomenda-se manter a skill sem alterações estruturais.
