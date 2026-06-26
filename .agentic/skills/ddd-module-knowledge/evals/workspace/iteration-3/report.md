# Relatório de Avaliação: `ddd-module-knowledge` — iteração 3

## Pontuação Geral

| Configuração   | Taxa de Aprovação Média |
|----------------|-------------------------|
| with_skill     | 1.000                   |
| without_skill  | 0.678                   |
| **delta**      | **+0.322**              |
| **value_tier** | **moderado**            |

Comparado à iteração anterior: delta passou de 0.289 → 0.322 (`sem mudança` de tier — ambas "moderado"). A melhoria de +0.033 no delta consolida a skill em pontuação perfeita (with_skill 1.000 pela primeira vez), eliminando a única lacuna existente na iteração 2.

## Custo estimado

| Configuração              | Total (USD)  | Média por avaliação (USD) | Modelo            |
|---------------------------|--------------|---------------------------|-------------------|
| with_skill                | $0.082476    | $0.005498                 | claude-sonnet-4-6 |
| without_skill             | $0.060612    | $0.004041                 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.021864** | —                      | —                 |

Nota de precificação: token split estimated 75/25 — actual split unavailable for both configurations. Runs em modo batch (15 evals por agente); timing compartilhado entre todos os evals do batch.

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
| eval-gotcha-builder-pattern-usage-for-child-entity | 1.000 | 0.750 | +0.250 |
| eval-gotcha-aggregate-root-name-must-equal-module-name | 1.000 | 1.000 | 0.000 |
| eval-naming-convention-http-handler-function-and-request-str | 1.000 | 0.333 | +0.667 |
| eval-scope-boundary-skill-must-not-generate-unit-tests | 1.000 | 0.667 | +0.333 |
| eval-domain-contract-where-interfaces-must-be-defined | 1.000 | 0.750 | +0.250 |
| eval-mapper-visibility-rule | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|------|-------|--------|
| eval-error-handling-validationerror-vs-invaliddataerror | +1.000 | without_skill inverte os HTTP codes (ValidationError→422) e atribui camadas incorretas; skill fornece distinção precisa 400/422 e application-infra vs domain |
| eval-naming-convention-http-handler-function-and-request-str | +0.667 | without_skill omite o sufixo 'Handler' no nome da função e não menciona o nome do arquivo |
| eval-naming-convention-tablename-return-value | +0.600 | without_skill retorna apenas 'users' sem o prefixo 'hex-api-go.' obrigatório |
| eval-module-structure-verify-correct-directory-layout-for-a | +0.571 | without_skill omite domain/contract/, usa 'events/' (plural), colapsa command/query e nomeia o arquivo raiz como module.go |
| eval-naming-convention-command-directory-name | +0.500 | without_skill sugere 'create_product' (snake_case) em vez de 'createproduct' (sem separador) |

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-error-handling-notfounderror-usage`
- `eval-layer-boundary-application-handler-containing-sql-quer`
- `eval-gotcha-domain-event-emission-by-child-entity`
- `eval-gotcha-aggregate-root-name-must-equal-module-name`
- `eval-mapper-visibility-rule`

## Lacunas da skill (with_skill < 1.0)

_(Nenhuma — pontuação perfeita)_

## Recomendação

A skill atingiu pontuação perfeita (with_skill 1.000) pela primeira vez — o ajuste no `domain-contract-pattern.md` que explicita a responsabilidade da camada de infraestrutura em implementar as interfaces resolveu a única lacuna remanescente; não são necessárias iterações adicionais salvo inclusão de novos casos de uso.
