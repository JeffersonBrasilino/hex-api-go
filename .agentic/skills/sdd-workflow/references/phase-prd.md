# Phase: `prd`

Step 1 already confirmed no PRD artifact (`PRD.md` nor `PRD.cache.md`) exists yet for this feature
(that's the only way the script reports this phase). `sdd-prd` is a standalone tool for the
product team, not orchestrated by this pipeline — ask the user directly for the PRD source instead
of instructing them to run it:
```
Aguardando PRD.

Qual é a origem do PRD desta feature — caminho de um arquivo local (ex: docs/<module>/<feature>/PRD.md)
ou um link/referência de card (Jira/GitHub/Trello/etc.)?
```
- **Local path given**: confirm the file exists; if so, tell the user to invoke `/sdd-plan
  {path}` (a new session) to continue. If it doesn't exist yet, tell them to produce it first (with
  `sdd-prd` or by hand) and return once it's saved.
- **Card link given**: tell the user to invoke `/sdd-plan {link}` (a new session) — the plan
  skill fetches and normalizes the card content itself.
Stop.
