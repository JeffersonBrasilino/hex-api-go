# Relatório de Avaliação: `spec-prd-v2` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.917 |
| without_skill | 0.602 |
| **delta** | **+0.315** |
| **value_tier** | **moderado** |

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.080004 | $0.008000 | claude-sonnet-4-6 |
| without_skill | $0.063108 | $0.006311 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.016896** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 0.800 | +0.200 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 1.000 | 0.600 | +0.400 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 0.750 | +0.250 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | +0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 0.200 | +0.800 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 0.667 | +0.333 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 0.667 | 0.333 | +0.334 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 0.750 | 0.500 | +0.250 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 0.750 | 0.500 | +0.250 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 0.667 | +0.333 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | +0.800 | Sem skill o modelo descreve o resumo em vez de apresentá-lo; com skill segue o formato padronizado exato |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | +0.400 | Sem skill o modelo não oferece resposta recomendada e responde em inglês; com skill ambos são atendidos |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao` (1.000 / 1.000) — o comportamento de bloqueio do gate de aprovação é robusto independentemente da skill.

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | RNFs listados como texto plano, não no formato `RNF-0X` | Adicionar instrução explícita no SKILL.md exigindo identificadores `RNF-0X` na seção 5 do PRD |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | Caminho `docs/auth/reset-password/PRD.md` não aparece na resposta | Instruir o modelo a confirmar o caminho completo do arquivo na mensagem ao usuário (Step 3) |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | Critérios apresentados sem mencionar gate de aprovação obrigatório antes do PRD | Reforçar no Step 1 que o gate (Step 2) é obrigatório mesmo durante iteração de critérios |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | Três dimensões LGPD descritas de uma vez em vez de uma pergunta por turno | Adicionar regra explícita: cada dimensão LGPD deve ser explorada em mensagens separadas, uma por vez |

## Recomendação

A skill agrega valor moderado (+0.315) e precisa de iteração focada em 4 lacunas pontuais: formato `RNF-0X`, confirmação do caminho de saída, disciplina do gate durante iteração de critérios BDD, e sequenciamento das perguntas LGPD.
