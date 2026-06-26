# Relatório de Avaliação: `spec-prd-v2` — iteração 5

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.972 |
| without_skill | 0.948 |
| **delta** | **+0.024** |
| **value_tier** | **sem_valor** |

Comparado à iteração anterior: delta passou de 0.077 → 0.024 (`fraco → sem_valor`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.082134 | $0.005133 | claude-sonnet-4-6 |
| without_skill | $0.084228 | $0.005264 | claude-sonnet-4-6 |
| **custo adicional da skill** | **-$0.002094** | — | — |

> Custo adicional negativo indica que o batch com skill consumiu menos tokens totais do que sem skill nesta execução.

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 1.000 | 0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 0.800 | 1.000 | -0.200 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | 0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 1.000 | 0.000 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | 0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 1.000 | 0.167 | +0.833 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 1.000 | 1.000 | 0.000 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 0.750 | 1.000 | -0.250 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 1.000 | 0.000 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao | 1.000 | 1.000 | 0.000 |
| eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo | 1.000 | 1.000 | 0.000 |
| eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo | 1.000 | 1.000 | 0.000 |
| eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal | 1.000 | 1.000 | 0.000 |
| eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | +0.833 | Sem a skill, o modelo apenas anunciou a geração do PRD sem exibir o conteúdo — 5 de 6 seções obrigatórias (RF-0X, RNF-0X, Fora de Escopo, Critérios de Aceitação com BDD, ausência de implementação) falharam. A skill garante PRD completo e estruturado. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade`
- `eval-implementacao-de-codigo-esta-fora-do-escopo`
- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao`
- `eval-resumo-de-aprovacao-segue-o-formato-padronizado`
- `eval-nao-salva-arquivos-na-raiz-do-projeto`
- `eval-criterios-de-aceitacao-usam-dado-quando-entao`
- `eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem`
- `eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao`
- `eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada`
- `eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo`
- `eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo`
- `eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal`
- `eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | "A resposta faz exatamente uma pergunta ao usuário, não múltiplas ao mesmo tempo" | Reforçar a instrução na seção Step 1 — Interview: a regra de uma única `?` por turno deve valer mesmo quando dois tópicos estão pendentes. A resposta com skill continha pergunta implícita embutida no texto antes da `?` final. |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | "A resposta NÃO faz mais de uma pergunta na mesma mensagem" | A skill violou a regra de pergunta única na seção de LGPD: fez duas perguntas no mesmo turno ("os dados são coletados no cadastro ou só ao editar o perfil?" + "o consentimento explícito é necessário?"). A instrução de LGPD deve reforçar explicitamente que cada dimensão é explorada em turnos separados, nunca bundled — igual à regra geral da entrevista. |

## Recomendação

A skill precisa de iteração para corrigir a violação da regra de pergunta única por turno, que ocorre de forma consistente em dois cenários distintos (entrevista padrão e exploração LGPD), causando regressão em relação ao baseline sem skill e rebaixamento do value_tier de `fraco` para `sem_valor`.
