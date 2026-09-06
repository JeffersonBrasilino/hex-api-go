# Relatório de Avaliação: `sdd-workflow` — iteração 3 (re-run direcionado)

## Contexto

Iteração de verificação pontual (`--ids 4`), rodando apenas `eval-plan-md-checkboxes-marked-via-update-plan-cjs-never-hand-edi` — o único eval que apresentou lacuna na iteração 2 (with_skill 0.75, faltava mencionar o passo de re-executar `compute-wave.cjs` após marcar tasks como `Done`).

Correção aplicada entre as iterações: em `references/phase-implement.md` (Step 4.3), o passo "refresh the wave cache" foi promovido a item obrigatório em negrito, no mesmo padrão visual dos outros dois passos da seção ("Update PLAN.md" / "Update STATE.md").

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 1.000 |
| without_skill | 0.250 |
| **delta** | **+0.750** |
| **value_tier** | **forte** |

Comparado à iteração 2 (mesmo eval, antes da correção): with_skill passou de 0.75 (3/4) → 1.00 (4/4). A asserção que antes falhava — "After updating, the output re-runs compute-wave.cjs to refresh the wave cache instead of assuming the next wave is current_wave + 1 by arithmetic" — agora passa: a resposta inclui explicitamente o comando `compute-wave.cjs` e o tratamento condicional de `STATUS: ok` / `STATUS: complete`.

## Lacunas da skill (with_skill < 1.0)

Nenhuma. Todas as 4 asserções passaram.

## Recomendação

Pronta: a lacuna identificada na iteração 2 foi fechada. Combinada com o resultado geral da iteração 2 (0.950 de pass rate médio nos 5 evals, `value_tier: moderado`), o split em progressive disclosure não introduziu regressões — a única perda pontual já foi corrigida.
