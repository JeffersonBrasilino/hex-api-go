# Write Operations — Scripts Reference

Schemas e comandos para os scripts de escrita. Carregue este arquivo apenas nos Steps 3 e 4.

---

## prd-write.js — escrita completa

```json
{
  "module": "string",
  "feature": "string",
  "title": "string",
  "status": "Draft",
  "sections": {
    "overview": "string",
    "objectives": ["string"],
    "userStories": [{ "epic": "string", "stories": ["string"] }],
    "functionalRequirements": [{ "id": "RF-01", "description": "string" }],
    "nonFunctionalRequirements": [{ "id": "RNF-01", "qualifier": "string", "description": "string" }],
    "businessRules": [{ "id": "RN-01", "description": "string" }],
    "dataRequirements": { "entities": "string", "origin": "string", "privacyLgpd": "string" },
    "acceptanceCriteria": [{ "name": "string", "scenario": "Dado [...], quando [...], então [...]." }],
    "outOfScope": ["string"],
    "openQuestions": ["string"],
    "futureConsiderations": ["string"]
  }
}
```

```bash
node .agentic/skills/sdd-prd/scripts/prd-write.js --input /tmp/prd-data.json
```

Saída: `{ folder, prd, author, date, version }`

---

## prd-write.js — patch parcial

Inclua apenas os campos que mudaram.

```json
{
  "module": "string",
  "feature": "string",
  "title": "string",
  "patch": {
    "acceptanceCriteria": [{ "name": "string", "scenario": "string" }],
    "outOfScope": ["string"]
  }
}
```

```bash
node .agentic/skills/sdd-prd/scripts/prd-write.js --input /tmp/prd-patch.json
```

Saída: `{ folder, prd, author, date, version }` — versão incrementada automaticamente.

**Erros:**
- `"Tags não substituídas: X"` → preencha os campos ausentes no JSON e reexecute.
- `"PRD não encontrado"` → use escrita completa primeiro.

---

## notes-write.js

```json
{
  "module": "string",
  "feature": "string",
  "title": "string",
  "notes": {
    "context": "string?",
    "functionalRequirements": ["string?"],
    "nonFunctionalRequirements": ["string?"],
    "businessRules": ["string?"],
    "dataRequirements": { "entities": "string", "origin": "string", "privacyLgpd": "string" },
    "glossary": [{ "term": "string", "definition": "string" }],
    "acceptanceCriteria": ["string?"],
    "openQuestions": ["string?"],
    "newDecision": "string?"
  }
}
```

```bash
node .agentic/skills/sdd-prd/scripts/notes-write.js --input /tmp/notes-data.json
```

Saída: `{ notes, date }`
Comportamento: arquivo ausente → escrita completa via template. Arquivo existente → merge por seção.
