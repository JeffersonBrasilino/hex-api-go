# Write Operations — Scripts Reference

Schemas e comandos para os scripts de escrita. Carregue este arquivo apenas na Fase 2.4 e na Fase 3.

---

## plan-write.js — escrita completa

```json
{
  "module": "string",
  "feature": "string",
  "title": "string",
  "status": "Draft",
  "ticket": "string?",
  "strategy": {
    "approach": "string?",
    "layers": "string?",
    "skills": "string?"
  },
  "tasks": [
    {
      "id": "TASK-[LAYER]-[CONCERN]",
      "description": "string",
      "file": "relative/path/to/file.go",
      "reason": "string",
      "depends_on": ["TASK-ID"],
      "parallel_group": "domain | application | infrastructure | module | tests | config",
      "complexity": {
        "scope": 1,
        "ambiguity": 1,
        "coupling": 1,
        "novelty": 1,
        "reversibility": 1
      },
      "risk_note": "string? (apenas quando tier=high)",
      "subtasks": ["string"],
      "completion_criterion": "string"
    }
  ]
}
```

```bash
node .agentic/skills/sdd-plan/scripts/plan-write.js --input /tmp/plan-data.json
```

Saída: `{ folder, plan, author, date, tasks_count }`

**Computed pelo script (não inclua no JSON):**
- `score` — média dos 5 valores de complexity
- `tier` — `low` ≤ 2.0 · `medium` ≤ 3.5 · `high` > 3.5
- `overall_complexity` — derivado do tier mais alto entre todas as tasks

---

## plan-write.js — patch de task

Forneça o objeto task **completo** (não apenas os campos alterados).

```json
{
  "module": "string",
  "feature": "string",
  "patch_task": {
    "id": "TASK-[LAYER]-[CONCERN]",
    "description": "string",
    "file": "relative/path/to/file.go",
    "reason": "string",
    "depends_on": ["TASK-ID"],
    "parallel_group": "string",
    "complexity": { "scope": 1, "ambiguity": 1, "coupling": 1, "novelty": 1, "reversibility": 1 },
    "risk_note": "string?",
    "subtasks": ["string"],
    "completion_criterion": "string"
  }
}
```

```bash
node .agentic/skills/sdd-plan/scripts/plan-write.js --input /tmp/plan-patch.json
```

Saída: `{ folder, plan, author, date, tasks_count }`

**Erros:**
- `"Task X não encontrada no PLAN.md"` → verifique se o `id` é exatamente o mesmo do arquivo.
- `"PLAN.md não encontrado"` → use escrita completa primeiro.
