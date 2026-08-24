# Relatório de Avaliação: `ddd-module-knowledge` — iteração 5

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.961 |
| without_skill | 0.750 |
| **delta** | **+0.211** |
| **value_tier** | **moderado** |

Comparado à iteração anterior (iteração 4): delta passou de 0.244 → 0.211 (`moderado → moderado`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.105120 | $0.007008 | claude-sonnet-4-6 |
| without_skill | $0.094677 | $0.006312 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.010443** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-module-structure-verify-layout | 1.000 | 0.429 | +0.571 |
| eval-naming-command-directory-name | 1.000 | 0.500 | +0.500 |
| eval-naming-gorm-repository-struct | 1.000 | 0.667 | +0.333 |
| eval-naming-tablename-return-value | 1.000 | 0.400 | +0.600 |
| eval-error-not-found-error-usage | 1.000 | 0.750 | +0.250 |
| eval-error-validation-vs-invalid-data | 1.000 | 1.000 | 0.000 |
| eval-layer-boundary-http-handler-repo | 1.000 | 0.667 | +0.333 |
| eval-layer-boundary-app-handler-sql | 0.750 | 1.000 | -0.250 |
| eval-gotcha-domain-event-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-builder-pattern-child-entity | 1.000 | 0.750 | +0.250 |
| eval-gotcha-aggregate-root-name | 1.000 | 1.000 | 0.000 |
| eval-naming-http-handler-function | 0.667 | 0.333 | +0.334 |
| eval-scope-boundary-no-unit-tests | 1.000 | 1.000 | 0.000 |
| eval-domain-contract-interface-location | 1.000 | 0.750 | +0.250 |
| eval-mapper-visibility-rule | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-naming-tablename-return-value | +0.600 | Sem a skill, o modelo retorna apenas `users` sem o prefixo `hex-api-go.` — convenção exclusiva do projeto |
| eval-module-structure-verify-layout | +0.571 | Sem a skill, o modelo não separa `domain/contract/` e `domain/event/` como subdiretórios distintos |
| eval-naming-command-directory-name | +0.500 | Sem a skill, o modelo sugere `create_product` (snake_case) em vez de `createproduct` (sem separador) |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-error-validation-vs-invalid-data`
- `eval-gotcha-domain-event-child-entity`
- `eval-gotcha-aggregate-root-name`
- `eval-scope-boundary-no-unit-tests`
- `eval-mapper-visibility-rule`

Cinco avaliações atingem teto (1.000/1.000). Os padrões de emissão de evento por aggregate root, nomenclatura do aggregate root, escopo da skill e visibilidade de mapper são conhecimentos bem consolidados no modelo base — a skill confirma mas não diferencia nessas questões.

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-layer-boundary-app-handler-sql | Não menciona `domain contract (interface)` como fronteira de abstração correta | Adicionar instrução explícita de que a fronteira é o contrato de domínio em `domain/contract/` |
| eval-naming-http-handler-function | Não especifica o nome do arquivo `create_product_handler.go` | Incluir convenção de nomenclatura de arquivo: `{verb}_{entity}_handler.go` |

**eval-layer-boundary-app-handler-sql** é o único caso com delta negativo (−0.250), onde o modelo base superou a skill por mencionar explicitamente a interface de domínio como fronteira de abstração.

## Recomendação

A skill mantém o tier **moderado** pelo segundo ciclo consecutivo com leve queda de −0.033 no delta. O valor continua positivo e concentrado nas convenções de nomenclatura específicas do projeto. Duas ações prioritárias para elevar o tier a **forte** (delta ≥ 0.40): (1) adicionar instrução explícita sobre domain contract como fronteira de abstração em `eval-layer-boundary-app-handler-sql`; (2) incluir convenção de nomenclatura de arquivo do handler HTTP para fechar a lacuna em `eval-naming-http-handler-function`.
