# Relatório de Avaliação: `spec-prd-v2` — iteração 2

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.882 |
| **delta** | **+0.118** |
| **value_tier** | **fraco** |

Comparado à iteração anterior: delta passou de 0.315 → 0.118 (`moderado` → `fraco`).

A queda no delta não indica regressão da skill — a skill manteve 100% de aprovação em todos os 10 evals. O que mudou foi uma melhora expressiva do baseline (without_skill subiu de 0.602 para 0.882), reduzindo o diferencial. Isso sinaliza que o modelo base internalizou parte dos comportamentos esperados, tornando a skill comparativamente menos diferenciada nesta iteração.

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.088338 | $0.008834 | claude-sonnet-4-6 |
| without_skill | $0.072570 | $0.007257 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.015768** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade | 1.000 | 1.000 | +0.000 |
| eval-com-contexto-a-skill-entrevista-uma-pergunta-por-vez | 1.000 | 0.800 | +0.200 |
| eval-implementacao-de-codigo-esta-fora-do-escopo | 1.000 | 1.000 | +0.000 |
| eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao | 1.000 | 1.000 | +0.000 |
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | 1.000 | 0.600 | +0.400 |
| eval-nao-salva-arquivos-na-raiz-do-projeto | 1.000 | 1.000 | +0.000 |
| eval-prd-gerado-contem-secoes-obrigatorias-e-fora-de-escopo | 1.000 | 0.667 | +0.333 |
| eval-criterios-de-aceitacao-usam-dado-quando-entao | 1.000 | 1.000 | +0.000 |
| eval-explora-privacidade-lgpd-quando-ha-dados-pessoais | 1.000 | 0.750 | +0.250 |
| eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem | 1.000 | 1.000 | +0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-resumo-de-aprovacao-segue-o-formato-padronizado | +0.400 | Baseline falha em duas asserções: ausência de "Critérios de aceitação" no resumo e ausência de cenários BDD (caminho feliz + erro/borda). A skill garante ambos com o formato padronizado exato. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-sem-contexto-a-skill-pergunta-o-problema-funcionalidade` (1.000 / 1.000) — comportamento de context-check consolidado no modelo base
- `eval-implementacao-de-codigo-esta-fora-do-escopo` (1.000 / 1.000) — recusa de código é robusta independentemente da skill
- `eval-gate-de-aprovacao-impede-gerar-prd-sem-aprovacao` (1.000 / 1.000) — gate de aprovação respeitado em ambos
- `eval-nao-salva-arquivos-na-raiz-do-projeto` (1.000 / 1.000) — convenção de path docs/ consolidada
- `eval-criterios-de-aceitacao-usam-dado-quando-entao` (1.000 / 1.000) — formato BDD e gate agora robustos em ambos
- `eval-notas-nao-sao-criadas-no-bootstrap-e-nao-sobrescrevem` (1.000 / 1.000) — comportamento de merge de NOTES.md consolidado

## Lacunas da skill (with_skill < 1.0)

Nenhuma. A skill atingiu 100% de aprovação em todos os 10 evals nesta iteração.

## Recomendação

A skill está funcionando sem lacunas próprias (with_skill = 1.000). O delta reduzido (+0.118, tier "fraco") é consequência do amadurecimento do modelo base, não de regressão da skill. Para iteração 3, recomenda-se substituir ou ampliar os 6 evals de baseline consolidado por cenários adversariais mais complexos onde o modelo base ainda falha, e aprofundar a cobertura de LGPD com asserções que verifiquem as três dimensões ao longo de um fluxo completo multi-turno.
