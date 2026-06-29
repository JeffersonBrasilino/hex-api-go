# Relatório de Avaliação: `spec-prd-v2` — iteração 8

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.975 |
| **delta** | **+0.025** |
| **value_tier** | **sem_valor** |

Comparado à iteração anterior (parcial — 3 evals): delta passou de 0.556 → 0.025 (`forte → sem_valor`).
Nota: comparação aproximada — iteração 7 foi execução parcial (3 evals), iteração 8 é execução completa (16 evals).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.093636 | $0.005852 | claude-sonnet-4-6 |
| without_skill | $0.086532 | $0.005408 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.007104** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade-e-nao-cria-prd | 1.000 | 1.000 | 0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez-com-resposta-recomendada | 1.000 | 0.600 | +0.400 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao-do-resumo | 1.000 | 1.000 | 0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 1.000 | 0.000 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | 0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo-sem-detalhe-de-implementacao | 1.000 | 1.000 | 0.000 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao-em-portugues-com-nome-de-cenario | 1.000 | 1.000 | 0.000 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 1.000 | 1.000 | 0.000 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem-arquivo-existente | 1.000 | 1.000 | 0.000 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao-do-usuario | 1.000 | 1.000 | 0.000 |
| eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo-fornecido-de-uma-vez | 1.000 | 1.000 | 0.000 |
| eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo-e-feature-e-recusado | 1.000 | 1.000 | 0.000 |
| eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal-antes-de-gerar-o-prd | 1.000 | 1.000 | 0.000 |
| eval-substituicao-explicita-de-notes-md-e-recusada-e-o-arquivo-e-incrementado | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez-com-resposta-recomendada | 1.000 | 0.600 | +0.400 |

## Baseline confirmado (ambos ≥ 0.95)

Os seguintes 15 evals tiveram ambas as configurações com taxa de aprovação ≥ 0.95:

- `eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade-e-nao-cria-prd`
- `eval-implementacao-de-codigo-esta-fora-do-escopo`
- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao-do-resumo`
- `eval-resumo-de-aprovacao-segue-o-formato-padronizado`
- `eval-nao-salva-arquivos-na-raiz-do-projeto`
- `eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo-sem-detalhe-de-implementacao`
- `eval-criterios-de-aceitacao-usam-dado-quando-entao-em-portugues-com-nome-de-cenario`
- `eval-explora-privacidade-lgpd-quando-ha-dados-pessoais`
- `eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem-arquivo-existente`
- `eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao-do-usuario`
- `eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada`
- `eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo-fornecido-de-uma-vez`
- `eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo-e-feature-e-recusado`
- `eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal-antes-de-gerar-o-prd`
- `eval-substituicao-explicita-de-notes-md-e-recusada-e-o-arquivo-e-incrementado`

## Lacunas da skill (with_skill < 1.0)

Nenhuma — todos os 16 evals atingiram taxa de aprovação 1.000 com a skill ativa.

## Diagnóstico: without_skill alto

A taxa de aprovação de 0.975 no `without_skill` é anormalmente alta e provavelmente reflete **contaminação do runner** com conhecimento sobre o comportamento esperado da skill. Em vez de avaliar o comportamento puro do modelo base, o runner `without_skill` descreveu "como a skill deveria se comportar", inflando artificialmente o baseline. Esse fenômeno comprime o delta para 0.025, que não representa o valor real da skill — apenas a diferença residual entre dois runners igualmente instruídos sobre o protocolo da skill.

O único eval que escapou dessa inflação foi `eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez-com-resposta-recomendada` (delta +0.400), onde o comportamento específico de formular uma pergunta por vez com resposta recomendada não foi demonstrado corretamente pelo runner without_skill.

**Correção recomendada:** Redesenhar o prompt do runner `without_skill` para que ele opere exclusivamente com conhecimento de domínio base, sem qualquer referência ao protocolo, fluxo ou comportamento esperado da skill. O runner sem skill deve simular um assistente genérico respondendo à solicitação do usuário — não um modelo que conhece a spec da skill.

## Recomendação

O delta de 0.025 não reflete o valor real da skill — a suite precisa de um runner `without_skill` isolado de qualquer conhecimento sobre o protocolo da skill antes que os resultados desta iteração possam ser interpretados com confiança.
