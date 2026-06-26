# Relatório de Avaliação: `spec-prd-v2` — iteração 4

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.9875 |
| without_skill | 0.9104 |
| **delta** | **+0.0771** |
| **value_tier** | **fraco** |

Comparado à iteração anterior: delta passou de 0.0614 → 0.0771 (`fraco → fraco`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.090804 | $0.005675 | claude-sonnet-4-6 |
| without_skill | $0.064800 | $0.004050 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.026004** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 1.000 | 0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 0.800 | 0.800 | 0.000 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | 0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 0.800 | +0.200 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | 0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 1.000 | 0.667 | +0.333 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 1.000 | 1.000 | 0.000 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 1.000 | 0.500 | +0.500 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 1.000 | 0.000 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao | 1.000 | 1.000 | 0.000 |
| eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo | 1.000 | 0.800 | +0.200 |
| eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo | 1.000 | 1.000 | 0.000 |
| eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal | 1.000 | 1.000 | 0.000 |
| eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | +0.500 | A skill guia o modelo a explorar ativamente requisitos de privacidade e LGPD quando dados pessoais estão presentes, comportamento ausente sem a skill. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade`
- `eval-implementacao-de-codigo-esta-fora-do-escopo`
- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao`
- `eval-nao-salva-arquivos-na-raiz-do-projeto`
- `eval-criterios-de-aceitacao-usam-dado-quando-entao`
- `eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem`
- `eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao`
- `eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada`
- `eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo`
- `eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal`
- `eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | "A resposta contém exatamente um sinal de interrogação '?'" — resposta encerrada com ponto final, zero perguntas encontradas | Reforçar no SKILL.md que ao receber contexto inicial, a resposta deve obrigatoriamente terminar com uma pergunta explícita (com `?`) dirigida ao usuário, mesmo quando uma recomendação já foi sugerida. |

## Recomendação

A skill está estável e madura em quase todos os cenários; a única lacuna restante — garantir que a entrevista sempre inicie com uma pergunta explícita (com `?`) mesmo quando o contexto já foi fornecido — deve ser corrigida antes de considerar a skill pronta para produção.
