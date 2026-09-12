# Phase: `done`

Read `STATE.md.verify.failed_tasks`. Compute successful tasks as
`STATE.md.implement.completed_tasks` (tasks that passed verification).

```
Pipeline concluído para {feature_path}.

Artefatos:
  PRD:  {artifacts.prd}
  PLAN: {artifacts.plan}

--- Tasks implementadas com sucesso ({N}) ---
  {TASK-ID}
  {TASK-ID}

--- Tasks com falha ({N}) ---
  {TASK-ID} — {reason}
  {TASK-ID} — {reason}

{Se failed_tasks estiver vazio: "Todas as tasks foram implementadas e verificadas com sucesso."}
{Se houver falhas: "Tasks com falha foram revertidas e marcadas no PLAN.md. Revise manualmente antes de re-executar o pipeline."}
```

After presenting the summary above, run the PR gate (Step 3.1 below). Only once it resolves to
`created` or `skipped` does the archive step run — this happens on every arrival at `done`, even
with failed tasks:

If `sdd-workflow.config.json`'s `archive_spec.destination_type` is `mcp`, check that the exact tool
named in `destination_name` is actually reachable in this session (e.g. via `ToolSearch` with
`select:{destination_name}`) **before** spawning `archive-spec` — no point paying for a subagent
spawn just to have it report back that the tool is missing. If it is not found, skip the spawn
entirely, append `Spec não arquivada: tool MCP "{destination_name}" indisponível nesta sessão.` to
the summary, and leave every local file untouched (same outcome as the subagent's own
`no-destination-available`). Otherwise, proceed:

Before spawning the subagent, read whatever fields you still need from `STATE.md` (`artifacts.prd`,
`artifacts.plan`, `artifacts.notes`), then delete `{feature_path}/STATE.md` — it is pipeline-only
bookkeeping with no value once the PR gate has resolved, and must never reach the archive
destination. Deleting it here, before `archive-spec` runs, means a failure past this point (a
partial archive, an unavailable MCP tool) can no longer be resumed by a later `/sdd-workflow
{feature_path}` invocation — `detect-state.cjs` has nothing left to detect. This is accepted: the
summary already shown to the user covers the pipeline outcome, and the "What the orchestrator does
with the report" contract below tells the user explicitly when a manual follow-up is needed instead
of relying on resume.

Load [archive-spec](.agentic/subagents/archive-spec.md) and inject:
- `{feature_path}` → the feature path for this pipeline
- `{prd_path}` → the `artifacts.prd` value read above
- `{plan_path}` → the `artifacts.plan` value read above
- `{notes_path}` → the `artifacts.notes` value read above (omit the line if empty)

Wait for its report and follow the contract in archive-spec.md ("What the orchestrator does with
the report"): on `archived`, remove `{feature_path}` if now empty, then append `Spec arquivada e
removida de docs/.`; otherwise append the agent's report verbatim, note that `STATE.md` has already
been removed so this feature can't be resumed automatically, and tell the user to re-run archiving
manually once the blocker is resolved (`phase` stays `done` either way).

## Step 3.1 — PR gate (runs before archive-spec)

Skip this entirely if `STATE.md.pr.status` is already `created` — proceed straight to the archive
step above.

1. Run:
   ```
   node .agentic/skills/sdd-workflow/scripts/prepare-pr.cjs {feature_path}
   ```
2. **`STATUS: blocked_missing_tool`**: run
   `write-state.cjs {feature_path} --pr-status blocked_missing_tool`, relay the script's `INSTALL`
   line to the user verbatim, and **stop the whole invocation** — do not spawn `archive-spec`. The
   next `/sdd-workflow {feature_path}` re-enters `phase: done`, re-runs this same step, and picks up
   automatically once the tool is installed.
3. **`STATUS: blocked_protected_branch`**: run
   `write-state.cjs {feature_path} --pr-status blocked_protected_branch`, relay the script's
   `MESSAGE` to the user, and stop the whole invocation — same resume behavior as above once the dev
   has switched to a feature branch.
4. **`STATUS: already_open`**: an open PR/MR already exists for this branch (resume case — a prior
   run got this far but failed before `archive-spec`). Run
   `write-state.cjs {feature_path} --pr-status created --pr-url {PR_URL} --pr-branch {BRANCH}`,
   append the URL to the summary already shown to the user, and proceed to `archive-spec` above.
5. **`STATUS: ready`**: branch/commit/push are done. Before spawning the writer, respect `{yolo}`
   (parsed in Step 0) exactly like the wave gate:
   - **`{yolo}` is `false`**: present `Branch: {BRANCH} -> {BASE_BRANCH}` and the task list, then wait
     for explicit confirmation (`sim`) before proceeding — pushing/opening a PR is visible to others
     and harder to reverse than a local wave.
   - **`{yolo}` is `true`**: skip confirmation, proceed immediately.

   Then load [pr-writer](.agentic/subagents/pr-writer.md) and inject:
   - `{provider}` → the script's `PROVIDER`
   - `{type}` → the script's `TYPE`
   - `{branch}` → the script's `BRANCH`
   - `{base_branch}` → the script's `BASE_BRANCH`
   - `{prd_path}` → the script's `PRD_PATH`
   - `{tasks_table}` → the script's `TASKS:` block, verbatim

   Wait for its report and follow the contract in pr-writer.md ("What the orchestrator does with
   the report"): on `created`, append the URL to the summary and proceed to `archive-spec` above;
   on `failed`, report the failure and stop — do not spawn `archive-spec`.
