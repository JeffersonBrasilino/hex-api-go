# Relatório de Avaliação: `spec-prd` — iteração 2

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 1.000 |
| **delta** | **0.000** |
| **value_tier** | **sem_valor** |

Comparado à iteração anterior: delta passou de **-0.130 → 0.000** (`negativo → sem_valor`). A skill deixou de ser prejudicial — melhora de +0.130 no delta de aprovação.

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.073760 | $0.014752 | claude-sonnet-4-6 |
| without_skill | $0.077225 | $0.015445 | claude-sonnet-4-6 |
| **custo adicional da skill** | **-$0.003465** | — | — |

_Nota de precificação: token split estimado 75/25 — split real indisponível (batch run). Custo negativo: with_skill foi marginalmente mais barato que o baseline._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-bootstrap-correto-solicita-modulo-e-feature-quando-nao | 1.000 | 1.000 | 0.000 |
| eval-elicitacao-faz-uma-pergunta-por-vez-em-pt-br-estilo-cav | 1.000 | 1.000 | 0.000 |
| eval-nao-gera-codigo-de-implementacao-pois-esta-fora-do-esco | 1.000 | 1.000 | 0.000 |
| eval-gate-2-impede-criacao-do-prd-sem-aprovacao-dos-criterio | 1.000 | 1.000 | 0.000 |
| eval-nao-cria-arquivos-fora-da-estrutura-docs-modulo-feature | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_Nenhuma avaliação com delta positivo ≥ 0.40 nesta iteração._

## Baseline confirmado (ambos ≥ 0.95)

Todos os 5 evals passam em ambas as configurações — guardas de regressão válidas:
- `eval-bootstrap-correto-solicita-modulo-e-feature-quando-nao`
- `eval-elicitacao-faz-uma-pergunta-por-vez-em-pt-br-estilo-cav`
- `eval-nao-gera-codigo-de-implementacao-pois-esta-fora-do-esco`
- `eval-gate-2-impede-criacao-do-prd-sem-aprovacao-dos-criterio`
- `eval-nao-cria-arquivos-fora-da-estrutura-docs-modulo-feature`

## Lacunas da skill (with_skill < 1.0)

_Nenhuma — pontuação perfeita._

## Diagnóstico: por que o delta é zero mesmo com pontuação perfeita

O delta zero não significa que a skill é inútil — significa que os evals atuais testam comportamentos que o **modelo base já conhece** sem instrução específica (recusar código, usar pt-BR, respeitar estrutura de pastas). O valor real da skill está em comportamentos que o modelo **não executaria corretamente sem ela**:

- Usar o template PRD exato com todas as 11 seções
- Estrutura específica do NOTES.md (`# Objetivos e Contexto`, `# Regras de Negócio`, etc.)
- Formato BDD em português com os termos `Dado/Quando/Então`
- Sequência de steps com gates obrigatórios em conversas longas

Estes aspectos não estão cobertos pelos evals atuais — são os candidatos para a próxima iteração de evals.

## Recomendação

A skill **está pronta para uso** — as correções da iteração 1 eliminaram o comportamento meta-descritivo e todos os evals passam. Para elevar o `value_tier` acima de `sem_valor`, adicionar evals que testem conhecimento exclusivo da skill: estrutura exata do NOTES.md, as 11 seções obrigatórias do PRD, e a sequência precisa de gates em fluxos multi-turno.
