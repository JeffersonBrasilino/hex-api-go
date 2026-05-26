# Relatório de Avaliação: `spec-prd` — iteração 1

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.870 |
| without_skill | 1.000 |
| **delta** | **-0.130** |
| **value_tier** | **negativo** |

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.074710 | $0.014942 | claude-sonnet-4-6 |
| without_skill | $0.076660 | $0.015332 | claude-sonnet-4-6 |
| **custo adicional da skill** | **-$0.001950** | — | — |

_Nota de precificação: token split estimado 75/25 — split real indisponível (batch run). Custo adicional negativo indica que a configuração with_skill foi marginalmente mais barata._

## Por Avaliação

| Slug | with_skill | without_skill | delta |
|---|---|---|---|
| eval-bootstrap-correto-solicita-modulo-e-feature-quando-nao | 1.000 | 1.000 | 0.000 |
| eval-elicitacao-faz-uma-pergunta-por-vez-em-pt-br-estilo-cav | 0.600 | 1.000 | -0.400 |
| eval-nao-gera-codigo-de-implementacao-pois-esta-fora-do-esco | 0.750 | 1.000 | -0.250 |
| eval-gate-2-impede-criacao-do-prd-sem-aprovacao-dos-criterio | 1.000 | 1.000 | 0.000 |
| eval-nao-cria-arquivos-fora-da-estrutura-docs-modulo-feature | 1.000 | 1.000 | 0.000 |

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

_Nenhuma avaliação com delta positivo ≥ 0.40 nesta iteração._

## Baseline confirmado (ambos ≥ 0.95)

Estas avaliações passam sem a skill — são guardas de regressão válidas:
- `eval-bootstrap-correto-solicita-modulo-e-feature-quando-nao`
- `eval-gate-2-impede-criacao-do-prd-sem-aprovacao-dos-criterio`
- `eval-nao-cria-arquivos-fora-da-estrutura-docs-modulo-feature`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-elicitacao-faz-uma-pergunta-por-vez-em-pt-br-estilo-cav | "A resposta faz exatamente uma pergunta ao usuário, não múltiplas ao mesmo tempo" — o agente with_skill forneceu dois pontos de interrogação no mesmo exemplo de resposta | Adicionar ao SKILL.md uma instrução explícita: "Nunca inclua dois pontos de interrogação em uma mesma mensagem" |
| eval-elicitacao-faz-uma-pergunta-por-vez-em-pt-br-estilo-cav | "A resposta menciona a criação ou uso do arquivo NOTES.md em docs/user/user-login/" — o agente with_skill omitiu a notificação de criação do NOTES.md | Reforçar no Step 0/Step 1 que a notificação do caminho do NOTES.md é obrigatória a cada resposta de elicitação |
| eval-nao-gera-codigo-de-implementacao-pois-esta-fora-do-esco | "A resposta inicia ou continua o processo de especificação do PRD (faz perguntas, cria notas ou estrutura de diretório)" — o agente with_skill descreveu comportamento em vez de enactuar, não iniciou elicitação | Reformular o skill prompt para instruir respostas enativas ("responda como se fosse o product owner") em vez de descritivas |

## Diagnóstico Principal

O resultado `negativo` reflete um problema de **modo de resposta do agente with_skill**, não necessariamente de conteúdo da skill. O agente with_skill respondeu em estilo meta-descritivo ("O skill faria X", "O skill executaria Y") em vez de atuar como a própria skill. O agente without_skill, sem restrições de instrução, simulou o comportamento correto da skill com maior fidelidade ao workflow esperado.

**Impacto por eval:**
- Evals que testam conhecimento declarativo (bootstrap, gate-2, nao-cria-fora): sem diferença — ambos os agentes acertam.
- Evals que testam comportamento enativo (elicitacao, nao-gera-codigo): o with_skill falhou por não demonstrar ações concretas.

## Recomendação

A skill **precisa de iteração**: reformular o SKILL.md para orientar o agente a responder de forma enativa (agindo como product owner) e não descritiva; reforçar a obrigatoriedade da notificação do NOTES.md em todas as interações de elicitação; corrigir o exemplo de pergunta única para garantir um único ponto de interrogação por mensagem.
