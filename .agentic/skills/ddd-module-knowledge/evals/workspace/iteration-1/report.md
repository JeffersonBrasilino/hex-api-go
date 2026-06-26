# Relatório de Avaliação: `ddd-module-knowledge` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.978 |
| without_skill | 0.756 |
| **delta** | **+0.222** |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-module-structure-verify-correct-directory-layout | 1.000 | 0.571 | +0.429 |
| eval-naming-convention-command-directory-name | 1.000 | 0.500 | +0.500 |
| eval-naming-convention-gorm-repository-struct-constructor | 1.000 | 0.667 | +0.333 |
| eval-naming-convention-tablename-return-value | 1.000 | 0.400 | +0.600 |
| eval-error-handling-notfounderror-usage | 1.000 | 1.000 | 0.000 |
| eval-error-handling-validationerror-vs-invaliddataerror | 1.000 | 0.200 | +0.800 |
| eval-layer-boundary-http-handler-calling-repository-direct | 1.000 | 0.667 | +0.333 |
| eval-layer-boundary-application-handler-containing-sql-quer | 1.000 | 1.000 | 0.000 |
| eval-gotcha-domain-event-emission-by-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-builder-pattern-usage-for-child-entity | 1.000 | 0.250 | +0.750 |
| eval-gotcha-aggregate-root-name-must-equal-module-name | 0.667 | 1.000 | -0.333 |
| eval-naming-convention-http-handler-function-and-request-st | 1.000 | 0.667 | +0.333 |
| eval-scope-boundary-skill-must-not-generate-unit-tests | 1.000 | 0.667 | +0.333 |
| eval-domain-contract-where-interfaces-must-be-defined | 1.000 | 0.750 | +0.250 |
| eval-mapper-visibility-rule | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-error-handling-validationerror-vs-invaliddataerror | +0.800 | Sem a skill, o modelo confunde camadas (usa "application ou domain boundary" para ValidationError) e mapeia InvalidDataError para "400/422" de forma ambígua. A skill fornece mapeamento preciso: 400 vs 422, application/infrastructure vs domain. |
| eval-gotcha-builder-pattern-usage-for-child-entity | +0.750 | Sem a skill, o modelo abre exceção ("unless the child entity has many optional fields") e não responde "No" explicitamente. A skill impõe a regra sem hedging. |
| eval-naming-convention-tablename-return-value | +0.600 | Sem a skill, o modelo retorna apenas "users" sem o prefixo do projeto. A skill fornece o padrão exato `hex-api-go.users`. |
| eval-naming-convention-command-directory-name | +0.500 | Sem a skill, o modelo sugere "create_product" (underscore). A skill impõe "createproduct" (sem separador). |
| eval-module-structure-verify-correct-directory-layout | +0.429 | Sem a skill, faltam: subdiretório `domain/contract/`, `domain/event/` (plural errado), e `product.go` na raiz do módulo. |

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-error-handling-notfounderror-usage`
- `eval-layer-boundary-application-handler-containing-sql-quer`
- `eval-gotcha-domain-event-emission-by-child-entity`
- `eval-mapper-visibility-rule`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-gotcha-aggregate-root-name-must-equal-module-name | "Output specifies PascalCase as the required casing" | A resposta com skill mostra o exemplo `account → Account` mas não nomeia explicitamente "PascalCase". Adicionar ao SKILL.md que a entidade agregada deve usar "PascalCase" como convenção de nomenclatura explícita. Curiosamente, a resposta sem skill passou nesta asserção por mencionar "PascalCase" diretamente. |

## Recomendação

A skill está **pronta para uso** com delta médio de +0.222; a única lacuna (eval-11, delta -0.333) é um gap de redação na resposta — a instrução existe no SKILL.md mas não foi mencionada explicitamente na resposta, sugerindo adicionar a palavra "PascalCase" no gotcha de nomenclatura do aggregate root para garantir que a resposta a inclua.
