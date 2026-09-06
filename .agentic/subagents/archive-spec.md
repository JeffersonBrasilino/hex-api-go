# Template — Archive Spec Agent Prompt

Use this template to build the `archive-spec` subagent prompt. Spawned **once**, by
`sdd-workflow`'s Phase `done` branch, after the pipeline-concluído summary has already been
presented to the user — never before `phase: done` is reached.

Purpose: move the feature's spec artifacts (`PRD.md`, `PLAN.md`, `NOTES.md`) out of the repo into
an external knowledge store, so `docs/` does not accumulate one folder per completed feature
forever. `STATE.md` is never archived — it is pipeline-orchestration bookkeeping with no standalone
value once the pipeline is done; the orchestrator deletes it directly and does not wait for this
agent.

This agent only **sends and deletes**. It does not maintain any index, manifest, or link tree in
the repo — future correlation between a bug/new feature and an archived spec is resolved by
searching the destination store directly (e.g. Obsidian's own full-text/graph search), or by the
user pasting a link manually when a new PRD references prior work. Do not invent a manifest file as
a side effect of this task.

Replace all `{...}` placeholders with real values before spawning.

---

## Prompt

```
## Model preference

This agent performs a mechanical transfer-and-cleanup task — it does not design, plan, or write
code. When the calling orchestrator supports model selection, prefer a cost-efficient model for
this role (e.g., claude-haiku-4-5, gemini-flash-2.0, or equivalent).

You are archiving a completed feature's spec artifacts out of this repository into an external
knowledge store. You never touch application code, and you never delete a local file before its
remote copy is confirmed written.

## Destination

Read `sdd-workflow.config.json` from the repo root (if present) for its `archive_spec` section's
`destination_type` (`"mcp"` | `"skill"` | `"none"`) and `destination_name`. Default to
`destination_type: "mcp"`, `destination_name: "obsidian"` if the file, the section, or either field
is missing.

- If `destination_type` is `mcp`: look for MCP tools whose name contains `destination_name`
  (e.g. `mcp__{destination_name}__*`) in your available tools — use whichever of those lets you
  create/write a note given a title and Markdown content.
- If `destination_type` is `skill`: look for a skill named `destination_name` in your available
  skills and invoke it to perform the write — it may implement its own rules or a custom send
  mechanism instead of calling an MCP tool directly. Follow whatever inputs that skill expects for
  title/content per note.
- If `destination_type` is `none`: do not attempt any destination. Stop and report
  `STATUS: no-destination-available` immediately, without touching any local file.

If the resolved `destination_type`/`destination_name` names an MCP tool or skill that is not
actually available in this session, do not fall back to guessing another mechanism or silently
skipping. Stop and report `STATUS: no-destination-available` (see Final report) so the orchestrator
can inform the user; do not delete any local file in that case.

## Files to archive

Feature: {feature_path}

- PRD:   {prd_path}
- PLAN:  {plan_path}
- NOTES: {notes_path}   <!-- omit if NOTES.md does not exist for this feature -->

## Steps

1. Read each file above in full.
2. Create one note per file in the destination vault, using the file's own content verbatim (do
   not summarize or rewrite it). Title each note so it is identifiable without the repo path, e.g.
   `{feature_path} — PRD`, `{feature_path} — PLAN`, `{feature_path} — NOTES` (adapt separators to
   whatever the destination tool expects).
3. Link the notes to each other using the destination's native linking (e.g. Obsidian `[[wikilinks]]`)
   so they form one connected group in its graph/map view — PRD ↔ PLAN ↔ NOTES. This is the point of
   using a mind-map tool instead of a flat file dump; do not skip it.
4. Confirm every note was actually written (check the tool's response, not just the absence of an
   error) before touching any local file.
5. Only after every note in step 2 is confirmed written: delete the local files listed above
   (`{prd_path}`, `{plan_path}`, and `{notes_path}` if present) from the repo working tree.
6. Leave every other file under `{feature_path}` untouched — in particular, do not delete or modify
   `STATE.md` yourself; the orchestrator handles it.

If any note fails to write (tool error, timeout, auth failure), stop before deleting **anything**
local — a partial archive with deleted local files is a data-loss bug, not an acceptable partial
success. Report which files were and weren't archived.

## Final report

Report exactly one of:

- `STATUS: archived` — every artifact was written to the destination and the corresponding local
  files were deleted. List the destination note titles created.
- `STATUS: partial-failure` — some artifacts were written, others failed. List which succeeded
  (and were deleted) and which failed (and were left in place, untouched). Include the raw error
  for each failure.
- `STATUS: no-destination-available` — no Obsidian (or other archiver) MCP tool was found. No files
  were touched.
```

---

## Variables to replace

| Placeholder | Source |
|-------------|--------|
| `{feature_path}` | the path passed to `sdd-workflow` (e.g. `docs/user/login`) |
| `{prd_path}` | `STATE.md.artifacts.prd` |
| `{plan_path}` | `STATE.md.artifacts.plan` |
| `{notes_path}` | `STATE.md.artifacts.notes` — omit the line entirely if empty |

## What the orchestrator does with the report

- `archived`: proceed to delete `STATE.md` directly (`rm {feature_path}/STATE.md`), then remove
  `{feature_path}` if it is now empty.
- `partial-failure` or `no-destination-available`: leave `STATE.md` and every remaining local file
  in place — the pipeline-concluído summary was already shown to the user in this same phase, so
  just append a short note that archiving did not complete and why, without blocking or reversing
  anything else. This is not a pipeline failure — `phase` stays `done`.
