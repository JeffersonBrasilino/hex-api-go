# Relatório de Avaliação: `spec-prd-v2` — iteração 6

## Pontuação Geral

| Métrica | with_skill | without_skill | delta | value_tier |
|---|---|---|---|---|
| pass_rate (média) | 0.908 | 0.935 | **-0.027** | **negativo** |
| time_seconds (média) | 3.159 s | 5.476 s | -2.317 s | — |
| tokens (média) | 953.25 | 678.375 | +274.875 | — |

Comparado à iteração anterior: o delta de pass_rate passou de **0.024 → -0.027** (value_tier: **sem_valor → negativo**). A skill agora prejudica levemente o desempenho médio em relação ao baseline sem skill.

## Custo Estimado

| | with_skill | without_skill | delta |
|---|---|---|---|
| Total (USD) | $0.091512 | $0.065124 | +$0.026388 |
| Média por eval (USD) | $0.005720 | $0.004070 | +$0.001650 |

A skill consome mais tokens e gera custo adicional de ~$0.026 por batch sem entregar ganho de qualidade nesta iteração.

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 1.000 | 0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 0.800 | 0.800 | 0.000 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | 0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 0.400 | 1.000 | **-0.600** |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | 0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 0.333 | 0.667 | **-0.334** |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 1.000 | 1.000 | 0.000 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 1.000 | 0.500 | **+0.500** |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 1.000 | 0.000 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao | 1.000 | 1.000 | 0.000 |
| eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo | 1.000 | 1.000 | 0.000 |
| eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo | 1.000 | 1.000 | 0.000 |
| eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal | 1.000 | 1.000 | 0.000 |
| eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta >= 0.40, ordem decrescente)

| Slug | delta | Motivo |
|---|---|---|
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | +0.500 | Skill instrui exploração LGPD uma dimensão por turno; sem skill faz 5 perguntas simultâneas |

Apenas um eval apresentou delta >= 0.40. Os dois evals com maior regressão (`eval-resumo-de-aprovacao-segue-o-formato-padronizado` e `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo`) cancelam este ganho e ainda geram resultado líquido negativo.

## Baseline confirmado (ambos >= 0.95)

- `eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade`
- `eval-implementacao-de-codigo-esta-fora-do-escopo`
- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao`
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

| Slug | Assertion que falhou | Sugestão de correção |
|---|---|---|
| `eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez` | A resposta contém exatamente um sinal de interrogação `?` | Garantir que cada turno de entrevista produza uma única sentença interrogativa terminada com `?`, sem listas ou sub-perguntas adicionais |
| `eval-resumo-de-aprovacao-segue-o-formato-padronizado` | (1) Resumo para aprovação não é apresentado; (2) Campos obrigatórios ausentes; (3) Cenários de caminho feliz e erro ausentes | O SKILL.md deve instruir a skill a gerar o bloco de resumo completo (Problema, Objetivo, Regras de negócio, Critérios de aceitação, cenários) — não apenas descrevê-lo |
| `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo` | (1) Seção `## 4. Requisitos Funcionais` com `RF-0X` ausente; (2) `## 5. Requisitos Não Funcionais` com `RNF-0X` ausente; (3) Seção `Fora de Escopo` ausente; (4) `## 8. Critérios de Aceitação` ausente | A skill está descrevendo o PRD em vez de escrevê-lo. O SKILL.md deve deixar explícito que após a aprovação a skill ESCREVE o arquivo com todas as seções; identificadores devem usar hífen: `RF-01`, `RNF-01` |

## Recomendação

Corrigir com prioridade máxima os dois evals com delta negativo expressivo (`eval-resumo-de-aprovacao-segue-o-formato-padronizado` e `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo`), ajustando o SKILL.md para que a skill produza conteúdo concreto e formatado em vez de descrever comportamentos esperados, e para que o formato dos identificadores de requisitos use hífen (`RF-0X`, `RNF-0X`).
