# Step 2 — Ensure local config exists

`.agentic/skills/sdd-workflow/assets/sdd-workflow.config.json` holds this repo's own settings for
the pipeline (git branch/PR conventions, archive destination, PRD-platform destination). It is
gitignored on purpose — every clone/project configures it independently, it is never the framework's
concern to version it. `sdd-plan`, `sdd-prd`, `prepare-pr.cjs`, and `archive-spec` all read this
same file; this step is the only place it gets written.

## Step 2.1 — Check existence

```
test -f .agentic/skills/sdd-workflow/assets/sdd-workflow.config.json && echo EXISTS || echo MISSING
```

- **`EXISTS`** and the user's request did not explicitly ask to reconfigure: skip straight to
  Step 3. Say nothing about config — this is invisible plumbing once it exists. *(strict — never
  re-ask just because a session restarted or state looks stale.)*
- **`MISSING`**, or the user's request explicitly asked to change/reset the sdd-workflow config
  (e.g. "reconfigurar sdd-workflow", "mudar o destino do PRD", "trocar branch base"): run the
  interview below. When reconfiguring an existing file, read its current values first and use them
  as the defaults shown in the interview instead of the framework defaults.

## Step 2.2 — Interview

Ask via `AskUserQuestion`, pt-BR labels, grouped by section. Skip follow-up fields a chosen
provider/type makes irrelevant (see notes per field) — don't ask for a Jira project key when the
user picked `none`.

**Git**
1. `git.provider` — `auto` (detecta pelo remote) | `github` | `gitlab` | `none`. Default: `auto`.
2. `git.base_branch` — nome da branch base. Default: `main`.
3. `git.auto_create_branch` — criar a branch da feature automaticamente? Default: `true`.
4. `git.branch_pattern` — padrão de nome de branch. Default: `{type}/{feature-slug}`.
5. `git.default_type` — tipo conventional-commit default (`feat`|`fix`|`refactor`|`perf`|`chore`|
   `docs`|`test`|`build`|`ci`). Default: `feat`.

**Archive** (destino de PRD/PLAN/NOTES ao concluir a pipeline — ver `subagents/archive-spec.md`)
6. `archive_spec.destination_type` — `mcp` | `skill` | `path` | `none`. Default: `none`.
7. `archive_spec.destination_name` — significado depende do tipo acima; só perguntar se não for
   `none`:
   - `mcp` | `skill`: nome do MCP/skill de destino (ex: `obsidian`).
   - `path`: caminho de uma pasta local (absoluto ou relativo à raiz do repo) fora de `docs/`, para
     onde os arquivos são movidos (ex: `~/notes/arquivados`). Criada automaticamente se não existir.
   Default: `""`.

**PRD** (plataforma de gerenciamento de projetos — ver `sdd-prd/SKILL.md` Step 3.5)
8. `prd.provider` — `jira` | `github` | `trello` | `none`. Default: `none`.
9. `prd.board_url` — URL do board/projeto. Só perguntar se o provider acima não for `none`.
10. `prd.project_key` — chave do projeto (Jira). Só perguntar se `provider` for `jira`.
11. `prd.list_id` — id da lista alvo (Trello). Só perguntar se `provider` for `trello`.
12. `prd.repo` — `owner/repo` (GitHub). Só perguntar se `provider` for `github`.

## Step 2.3 — Write

```
node .agentic/skills/sdd-workflow/scripts/write-config.cjs \
  --git-provider {value} --base-branch {value} --auto-create-branch {true|false} \
  --branch-pattern {value} --default-type {value} \
  --archive-destination-type {value} --archive-destination-name {value} \
  --prd-provider {value} --prd-board-url {value} --prd-project-key {value} \
  --prd-list-id {value} --prd-repo {value}
```

- **exit 0**: config written. Confirm in one pt-BR line (e.g. `Configuração do sdd-workflow salva
  em .agentic/skills/sdd-workflow/assets/sdd-workflow.config.json.`) and continue to Step 3.
- **exit 1**: an enum value was invalid — fix just that field from the interview answer and retry;
  don't restart the whole interview for one bad field.
