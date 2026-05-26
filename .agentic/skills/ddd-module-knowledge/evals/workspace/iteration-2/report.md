# Relatório de Avaliação: `ddd-module-knowledge` — iteração 2

## Pontuação Geral

| Configuração   | Taxa de Aprovação Média |
|----------------|-------------------------|
| with_skill     | 0.983                   |
| without_skill  | 0.694                   |
| **delta**      | **+0.289**              |
| **value_tier** | **moderado**            |

Comparado à iteração anterior: delta passou de 0.222 → 0.289 (`sem mudança` de tier — ambas "moderado"). A melhoria de +0.067 no delta indica progresso na skill, mas sem mudar de tier.

## Custo estimado

| Configuração              | Total (USD)  | Média por avaliação (USD) | Modelo            |
|---------------------------|--------------|---------------------------|-------------------|
| with_skill                | $0.095640    | $0.006376                 | claude-sonnet-4-6 |
| without_skill             | $0.072435    | $0.004829                 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.023205** | —                      | —                 |

Nota de precificação: token split estimado 75/25 — split real não disponível para ambas as configurações. Os runs foram em modo batch (15 evals por agente); timing compartilhado entre todos os evals do batch.

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|------|------------|---------------|-------|
| eval-module-structure-verify-correct-directory-layout-for-a | 1.000 | 0.429 | +0.571 |
| eval-naming-convention-command-directory-name | 1.000 | 0.500 | +0.500 |
| eval-naming-convention-gorm-repository-struct-and-constructo | 1.000 | 0.667 | +0.333 |
| eval-naming-convention-tablename-return-value | 1.000 | 0.400 | +0.600 |
| eval-error-handling-notfounderror-usage | 1.000 | 1.000 | 0.000 |
| eval-error-handling-validationerror-vs-invaliddataerror | 1.000 | 0.000 | +1.000 |
| eval-layer-boundary-http-handler-calling-repository-directly | 1.000 | 0.667 | +0.333 |
| eval-layer-boundary-application-handler-containing-sql-quer | 1.000 | 1.000 | 0.000 |
| eval-gotcha-domain-event-emission-by-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-builder-pattern-usage-for-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-aggregate-root-name-must-equal-module-name | 1.000 | 1.000 | 0.000 |
| eval-naming-convention-http-handler-function-and-request-str | 1.000 | 0.333 | +0.667 |
| eval-scope-boundary-skill-must-not-generate-unit-tests | 1.000 | 0.667 | +0.333 |
| eval-domain-contract-where-interfaces-must-be-defined | 0.750 | 0.750 | 0.000 |
| eval-mapper-visibility-rule | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|------|-------|--------|
| eval-error-handling-validationerror-vs-invaliddataerror | +1.000 | without_skill inverte completamente HTTP codes (400↔422) e camadas (domain↔application-infra) para os dois erros |
| eval-naming-convention-http-handler-function-and-request-str | +0.667 | without_skill omite o sufixo 'Handler' no nome da função e não fornece o nome do arquivo |
| eval-naming-convention-tablename-return-value | +0.600 | without_skill retorna só 'users' sem o prefixo 'hex-api-go.' obrigatório |
| eval-module-structure-verify-correct-directory-layout-for-a | +0.571 | without_skill omite domain/contract/, usa 'events/' (plural), colapsa command/query e nomeia o arquivo raiz como module.go |
| eval-naming-convention-command-directory-name | +0.500 | without_skill sugere 'create_product' (snake_case) em vez de 'createproduct' (sem separador) |

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-error-handling-notfounderror-usage`
- `eval-layer-boundary-application-handler-containing-sql-quer`
- `eval-gotcha-domain-event-emission-by-child-entity`
- `eval-gotcha-builder-pattern-usage-for-child-entity`
- `eval-gotcha-aggregate-root-name-must-equal-module-name`
- `eval-mapper-visibility-rule`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|------|----------------|-------------------|
| eval-domain-contract-where-interfaces-must-be-defined | "Output states the infrastructure layer is responsible for implementing the interface" | Adicionar ao SKILL.md e/ou reference domain-contract-pattern.md uma frase explícita: "A camada de infraestrutura é responsável por implementar as interfaces definidas em domain/contract/" |

## Recomendação

A skill está sólida em nomenclatura e tratamento de erros (tier "moderado" com tendência de melhora: +0.067 vs iteração 1); o único ajuste necessário é explicitar no padrão de domain contract que a infraestrutura é quem implementa a interface.
