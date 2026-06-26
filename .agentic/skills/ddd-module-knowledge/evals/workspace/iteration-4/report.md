# Relatório de Avaliação: `ddd-module-knowledge` — iteração 4

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.983 |
| without_skill | 0.739 |
| **delta** | **+0.244** |
| **value_tier** | **moderado** |

Comparado à iteração anterior: delta passou de 0.322 → 0.244 (`sem mudança de tier: moderado → moderado`).

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.099288 | $0.006619 | claude-sonnet-4-6 |
| without_skill | $0.077658 | $0.005177 | claude-sonnet-4-6 |
| **custo adicional da skill** | **$0.021630** | — | — |

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-module-structure-verify-correct-directory-layout-for-a | 1.000 | 0.571 | +0.429 |
| eval-naming-convention-command-directory-name | 1.000 | 0.500 | +0.500 |
| eval-naming-convention-gorm-repository-struct-and-constructo | 1.000 | 0.667 | +0.333 |
| eval-naming-convention-tablename-return-value | 1.000 | 0.400 | +0.600 |
| eval-error-handling-notfounderror-usage | 1.000 | 1.000 | 0.000 |
| eval-error-handling-validationerror-vs-invaliddataerror | 1.000 | 0.200 | +0.800 |
| eval-layer-boundary-http-handler-calling-repository-directly | 1.000 | 0.667 | +0.333 |
| eval-layer-boundary-application-handler-containing-sql-query | 0.750 | 1.000 | -0.250 |
| eval-gotcha-domain-event-emission-by-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-builder-pattern-usage-for-child-entity | 1.000 | 1.000 | 0.000 |
| eval-gotcha-aggregate-root-name-must-equal-module-name | 1.000 | 1.000 | 0.000 |
| eval-naming-convention-http-handler-function-and-request-str | 1.000 | 0.333 | +0.667 |
| eval-scope-boundary-skill-must-not-generate-unit-tests | 1.000 | 1.000 | 0.000 |
| eval-domain-contract-where-interfaces-must-be-defined | 1.000 | 0.750 | +0.250 |
| eval-mapper-visibility-rule | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-error-handling-validationerror-vs-invaliddataerror | +0.800 | Sem a skill, o modelo inverte completamente as camadas: atribui `ValidationError` ao domínio e `InvalidDataError` à camada de aplicação/infraestrutura, além de omitir os códigos HTTP 400/422. A skill impõe a distinção correta por camada e código. |
| eval-naming-convention-http-handler-function-and-request-str | +0.667 | Sem a skill, o modelo sugere `CreateProduct` (sem sufixo `Handler`) e omite o nome do arquivo. A skill impõe `CreateProductHandler`, `CreateProductRequest` e `create_product_handler.go`. |
| eval-naming-convention-tablename-return-value | +0.600 | Sem a skill, o modelo retorna apenas `users` ignorando o prefixo do projeto. A skill impõe o padrão `hex-api-go.users`. |
| eval-naming-convention-command-directory-name | +0.500 | Sem a skill, o modelo sugere `create_product` (snake_case). A skill impõe `createproduct` (sem separador, tudo minúsculo). |
| eval-module-structure-verify-correct-directory-layout-for-a | +0.429 | Sem a skill, o modelo omite `domain/contract/`, `domain/event/` como subdiretórios distintos e não menciona `product.go` na raiz do módulo. |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-error-handling-notfounderror-usage`
- `eval-gotcha-domain-event-emission-by-child-entity`
- `eval-gotcha-builder-pattern-usage-for-child-entity`
- `eval-gotcha-aggregate-root-name-must-equal-module-name`
- `eval-scope-boundary-skill-must-not-generate-unit-tests`
- `eval-mapper-visibility-rule`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-layer-boundary-application-handler-containing-sql-query | A asserção "Output mentions domain contract (interface) as the correct abstraction boundary" falhou (0.750). O with_skill não citou explicitamente o contrato de domínio como fronteira de abstração. | Adicionar no SKILL.md, nas Application Layer Rules: "Handlers devem acessar persistência exclusivamente via interface definida em `domain/contract/` — nunca via SQL direto, ORM direto ou qualquer detalhe de infraestrutura." |

## Regressões vs iteração anterior

O delta caiu de **0.322** (iteração 3) para **0.244** (iteração 4), redução de 0.078 pontos. O tier permanece `moderado`. A principal causa é a regressão em `eval-layer-boundary-application-handler-containing-sql-query` (delta = −0.250): com a skill carregada, o modelo omite a menção ao contrato de domínio como fronteira de abstração, enquanto o baseline sem skill responde corretamente. Essa regressão localizada está relacionada à instrução da Application Layer não reforçar explicitamente o uso de `domain/contract/` como mecanismo de isolamento.

## Recomendação

A skill precisa de iteração pontual: incluir nas Application Layer Rules uma instrução explícita sobre uso obrigatório de `domain/contract/` como única fronteira de acesso a persistência, corrigindo a regressão em `eval-layer-boundary-application-handler-containing-sql-query` e recuperando o delta para próximo de 0.322.
