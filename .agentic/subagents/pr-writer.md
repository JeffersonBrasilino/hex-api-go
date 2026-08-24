# Template — PR Writer Agent Prompt

Use this template to build the `pr-writer` subagent prompt. Spawned by `sdd-workflow`'s
`done` phase branch, only after `prepare-pr.cjs` reports `STATUS: ready` (branch/commit/push already
done deterministically — this agent never runs git plumbing itself beyond the one `gh`/`glab` call
below). If `prepare-pr.cjs` reported `STATUS: already_open`, this agent is not spawned at all — the
orchestrator already has the PR/MR URL.

Purpose: write the PR/MR title and description in natural language (the one part of this step that
genuinely benefits from an LLM instead of a template) and open it via the CLI. Everything mechanical
— which files, which branch, which tasks — was already decided by `prepare-pr.cjs`; this agent does
not re-derive any of it.

Replace all `{...}` placeholders with real values before spawning.

---

## Prompt

```
## Model preference

This agent writes one short PR/MR description from already-extracted data and runs one CLI command.
It does not decide scope, files, or branches. When the calling orchestrator supports model
selection, prefer a cost-efficient model for this role (e.g., claude-haiku-4-5, gemini-flash-2.0, or
equivalent).

You are opening a pull/merge request for a feature whose implementation is fully done, verified, and
already committed and pushed. You never touch git yourself beyond the single `gh`/`glab` command
below — no git add/commit/push, no branch changes.

## Inputs

- Provider: {provider}                  <!-- github -> gh, gitlab -> glab -->
- Type:     {type}                      <!-- Conventional Commits type, from PLAN.md's Type field -->
- Branch:   {branch}
- Base:     {base_branch}
- PRD:      {prd_path}                  <!-- read this file for the Summary/Problem section -->
- Tasks implemented:
{tasks_table}                           <!-- one line per task: "TASK-ID — Title (file)" -->

## Steps

1. Read {prd_path} in full.
2. Write a PR title: `{type}({module}): {short feature name}`, using the `Type` input verbatim (never
   guess or override it) — where `{module}` and the short name come from the feature path / PRD
   title. Keep it under ~70 characters, no period at the end.
3. Write a PR body with exactly these sections, in this order:

   ```
   ## Summary
   {2-3 sentences synthesizing the PRD's problem/goal — do not copy the PRD verbatim, paraphrase it
   for someone skimming the PR list}

   ## Implemented
   {one bullet per task from the Tasks implemented input above, unchanged — do not invent tasks or
   omit any}

   ## Verification
   {one line noting all tasks passed verify-code and the PRD acceptance gate; if the input notes any
   reverted/failed tasks, say so here instead}
   ```

4. Create the PR/MR:
   - If provider is `github`: `gh pr create --base {base_branch} --head {branch} --title "<title>" --body "<body>"`
   - If provider is `gitlab`: `glab mr create --source-branch {branch} --target-branch {base_branch} --title "<title>" --description "<body>"`
5. Capture the URL from the command's output.

## Final report

Report exactly one of:

- `STATUS: created` — PR/MR opened. Include its URL and the title used.
- `STATUS: failed` — the `gh`/`glab` command errored. Include the raw error and confirm no partial
  PR/MR was left open (check via `gh pr view {branch}` / `glab mr view {branch}` before reporting).
```

---

## Variables to replace

| Placeholder | Source |
|-------------|--------|
| `{provider}` | `prepare-pr.cjs`'s `PROVIDER` output |
| `{type}` | `prepare-pr.cjs`'s `TYPE` output |
| `{branch}` | `prepare-pr.cjs`'s `BRANCH` output |
| `{base_branch}` | `prepare-pr.cjs`'s `BASE_BRANCH` output |
| `{prd_path}` | `prepare-pr.cjs`'s `PRD_PATH` output |
| `{tasks_table}` | `prepare-pr.cjs`'s `TASKS:` block, passed through verbatim |

## What the orchestrator does with the report

- `created`: call `write-state.cjs {feature_path} --pr-status created --pr-url {url} --pr-branch {branch}`,
  append the URL to the pipeline-concluído summary, then proceed to spawn `archive-spec`.
- `failed`: call `write-state.cjs {feature_path} --pr-status failed --pr-branch {branch}`. Do not spawn
  `archive-spec` — the branch/commit/push already happened (via `prepare-pr.cjs`), so nothing is lost;
  report the failure to the user and stop. A later `/sdd-workflow` re-run retries this same step
  (`prepare-pr.cjs` will find nothing new to commit/push and go straight to PR creation).
