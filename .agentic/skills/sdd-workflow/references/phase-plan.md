# Phase: `plan`

Step 1 already confirmed either `PLAN.md` doesn't exist yet, or it exists but isn't approved
(`PLAN_STATUS` is `Draft`, `Planning`, or `Unknown`). Distinguish using the script's `PLAN` field:

**`PLAN: (none)`:**
```
PRD detectado: {PRD}

--- Avancando para fase: PLAN ---

Abra uma nova sessão e execute:

  /sdd-plan {PRD}

Quando o PLAN.md estiver aprovado e salvo, volte e execute /sdd-workflow novamente.
```

**`PLAN.md` exists but not approved:**
```
PLAN.md detectado, porém ainda não aprovado (Status: {PLAN_STATUS}).

Abra uma nova sessão e execute /sdd-plan {PRD} para revisar/aprovar o plano.
Quando o Status mudar para "Ready for Implementation", volte e execute /sdd-workflow novamente.
```
Stop.

`{PRD}` here is exactly the `PRD:` value `detect-state.cjs` reported in Step 1 — a local `PRD.md`
path, or a previously materialized `PRD.cache.md` path if this feature's plan was started from a
card. Never hardcode `{feature_path}/PRD.md`; the actual filename varies.
