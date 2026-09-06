# Relatório de Avaliação: `sdd-workflow` — iteração 2

## Pontuação Geral

| Configuração | Taxa de Aprovação Média |
|---|---|
| with_skill | 0.950 |
| without_skill | 0.600 |
| **delta** | **+0.350** |
| **value_tier** | **moderado** |

Comparado à iteração anterior: delta passou de 0.450 → 0.350 (`forte → moderado`).

Essa queda no `value_tier` ocorre no contexto de um refatoramento estrutural da skill: o `SKILL.md` monolítico (645 linhas) foi dividido em um roteador enxuto (~122 linhas) mais 4 arquivos de referência carregados sob demanda por fase (`references/phase-prd.md`, `phase-plan.md`, `phase-implement.md`, `phase-done.md`), além da remoção de prosa redundante de "Gotchas". A separação em arquivos carregados sob demanda parece ter deixado uma instrução crítica (recarregar o cache de waves) menos visível no fluxo de "marcar tarefas como concluídas", o que provavelmente explica a regressão pontual observada em `eval-plan-checkboxes` sem afetar as demais avaliações.

## Custo estimado

| Configuração | Total (USD) | Média por avaliação (USD) | Modelo |
|---|---|---|---|
| with_skill | $0.183000 | $0.036600 | claude-sonnet-5 |
| without_skill | $0.132000 | $0.026400 | claude-sonnet-5 |
| **custo adicional da skill** | **$0.051000** | — | — |

O custo adicional da skill nesta iteração é $0.051000, valor $0.6837 menor do que na iteração anterior (`cost_delta_change_usd`), refletindo o carregamento sob demanda dos arquivos de referência em vez do `SKILL.md` completo em toda invocação.

## Skill agrega valor claro (delta ≥ 0.40, ordem decrescente)

| Slug | Delta | Motivo |
|---|---|---|
| eval-wave-computed-via-compute-wave-cjs-from-the-persisted-wave-f | +0.75 | with_skill 1.00 (4/4) vs without_skill 0.25 (1/4) — sem a skill, o modelo raramente calcula a wave via `compute-wave.cjs` a partir do estado persistido |
| eval-retry-counter-updated-via-write-state-cjs-increment-retry-ne | +0.50 | with_skill 1.00 (4/4) vs without_skill 0.50 (2/4) — a skill reforça o incremento do contador de retry via `write-state.cjs` |
| eval-plan-md-checkboxes-marked-via-update-plan-cjs-never-hand-edi | +0.50 | with_skill 0.75 (3/4) vs without_skill 0.25 (1/4) — mesmo com a lacuna identificada abaixo, a skill ainda supera fortemente a ausência dela |

## Baseline confirmado (ambos ≥ 0.95)

- `eval-bootstrap-detect-state-cjs-write-state-cjs-initialize-state`
- `eval-wave-confirmation-gate-still-enforced-before-spawning-subage`

## Lacunas da skill (with_skill < 1.0)

| Slug | Asserção falha | Correção sugerida |
|---|---|---|
| eval-plan-md-checkboxes-marked-via-update-plan-cjs-never-hand-edi | Após atualizar o plano, a resposta deve re-executar `compute-wave.cjs` para atualizar o cache de waves, em vez de assumir que a próxima wave é `current_wave + 1` por aritmética simples | Em `references/phase-implement.md` (Step 4.3), tornar explícito e destacado o passo de re-executar `compute-wave.cjs` logo após `update-plan.cjs` marcar tarefas concluídas, deixando claro que a wave seguinte não deve ser inferida por incremento aritmético |

## Recomendação

Precisa de iteração: a skill segue trazendo valor moderado e consistente, mas o passo de recalcular a wave via `compute-wave.cjs` após atualizar o plano precisa ser reforçado na referência de fase de implementação para eliminar a única lacuna identificada.
