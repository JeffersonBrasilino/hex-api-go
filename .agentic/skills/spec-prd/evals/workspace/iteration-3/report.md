# Relatório de Avaliação: `spec-prd-v2` — iteração 3

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.956 |
| without_skill | 0.895 |
| **delta** | **+0.061** |
| **value_tier** | **fraco** |

Comparado à iteração anterior: delta passou de 0.118 → 0.061 (`fraco → fraco`).
Nota: iteração 2 tinha 10 evals; iteração 3 expandiu para 16 evals (6 novos cenários).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.107226 | $0.006702 | claude-sonnet-4-6 |
| without_skill | $0.081474 | $0.005092 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.025752** | — | — |

> Nota: os valores de tempo e tokens são compartilhados por lote (batch run) entre todos os 16 evals — as médias por eval são derivadas da divisão do total por 16.

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 1.000 | 0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 0.800 | 0.800 | 0.000 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | 0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 0.600 | +0.400 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | 0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 1.000 | 0.667 | +0.333 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 0.750 | 0.500 | +0.250 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 0.750 | 0.750 | 0.000 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 1.000 | 0.000 |
| eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao | 1.000 | 1.000 | 0.000 |
| eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada | 1.000 | 1.000 | 0.000 |
| eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo | 1.000 | 1.000 | 0.000 |
| eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo | 1.000 | 1.000 | 0.000 |
| eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal | 1.000 | 1.000 | 0.000 |
| eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta >= 0.40, ordem decrescente)

| Slug | delta |
|---|---|
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | +0.400 |

## Baseline confirmado (ambos >= 0.95)

Os seguintes evals passaram com taxa >= 0.95 em ambas as configurações:

- eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade (1.000 / 1.000)
- eval-implementacao-de-codigo-esta-fora-do-escopo (1.000 / 1.000)
- eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao (1.000 / 1.000)
- eval-nao-salva-arquivos-na-raiz-do-projeto (1.000 / 1.000)
- eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem (1.000 / 1.000)
- eval-sem-contexto-real-a-skill-nao-inicia-entrevista-mesmo-sob-pressao (1.000 / 1.000)
- eval-implementacao-disfarcada-de-rnf-e-recusada-e-redirecionada (1.000 / 1.000)
- eval-gate-de-aprovacao-obrigatorio-mesmo-com-contexto-completo (1.000 / 1.000)
- eval-caminho-invalido-dentro-de-docs-sem-subpastas-de-modulo (1.000 / 1.000)
- eval-criterios-aprovados-nao-dispensam-o-gate-de-aprovacao-formal (1.000 / 1.000)
- eval-substituicao-explicita-de-notes-md-e-recusada-e-incrementado (1.000 / 1.000)

## Lacunas da skill (with_skill < 1.0)

| Slug | Assercao falha | Correcao sugerida |
|---|---|---|
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | A resposta contem exatamente um sinal de interrogacao '?' | Proibir frases de confirmacao ao final da mensagem que gerem um segundo '?'; converter confirmacoes em afirmacoes ou omiti-las. |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | A resposta aguarda aprovacao explicita do usuario antes de gerar o PRD | Apos exibir os criterios de aceitacao, incluir obrigatoriamente uma mensagem de gate explicita solicitando aprovacao antes de prosseguir com a geracao do PRD. |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | A resposta NAO faz mais de uma pergunta na mesma mensagem | Separar cada dimensao de privacidade/LGPD (coleta, retencao, compartilhamento, etc.) em perguntas distintas, uma por turno de conversa. |

## Recomendacao

A skill precisa de iteracao: os tres evals com falha compartilham a raiz comum de violacao da regra "uma pergunta por vez" e ausencia do gate de aprovacao explicito em fluxos intermediarios — corrigi-los elevaria o delta e consolidaria o tier moderado.
