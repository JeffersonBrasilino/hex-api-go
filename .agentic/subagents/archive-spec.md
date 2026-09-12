# Template — Archive Spec Agent Prompt

Use this template to build the `archive-spec` subagent prompt. Spawned **once**, by
`sdd-workflow`'s Phase `done` branch, after the pipeline-concluído summary has already been
presented to the user — never before `phase: done` is reached.

Purpose: move the feature's spec artifacts (`PRD.md`, `PLAN.md`, `NOTES.md`) out of the repo into
an external knowledge store, so `docs/` does not accumulate one folder per completed feature
forever. `STATE.md` is never archived — it is pipeline-orchestration bookkeeping with no standalone
value once the pipeline is done, and it must never reach the archive destination. The orchestrator
deletes it directly, *before* spawning this agent — by the time this agent runs, `STATE.md` is
already gone from `{feature_path}`. This agent never needs to special-case it, but if `STATE.md` is
ever unexpectedly still present, do not send it or move it — it is not one of the "Files to archive"
below, full stop.

This agent never touches the repo's own index or manifest — it only sends and deletes local files.
Cross-file correlation (PRD ↔ PLAN ↔ NOTES) is authored upstream, by the `sdd-prd`/`sdd-plan`
templates themselves, at creation time — not here.

The exception is `path` and `skill`: both maintain one `{module}.md` file per module at the
destination (named after the module itself, e.g. `users.md` for module `users`), listing every
archived feature grouped by type, linking only to that feature's PRD — the PRD's own `Correlations`
section is what leads from there to PLAN/NOTES, so the module file stays a one-line-per-feature
index rather than duplicating every link. `path` writes it directly, as a plain file, alongside the
mirrored folder structure. `skill` delegates it to the invoked skill, since that skill is a full
agent and can implement its own read-modify-write — the request just needs to be part of what's
asked of it.

`mcp` does **not** get index maintenance. `destination_name` there is a single, exact tool call
(e.g. `mcp__obsidian__create_note`) with no read/search tool guaranteed alongside it — this agent
cannot safely merge a new entry into an index note it cannot read back, and blindly recreating the
note risks overwriting every previous feature's entry. This is a known limitation of the single-tool
`mcp` destination, not something to work around by guessing at other tools under the same server
prefix.

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

Read `.agentic/skills/sdd-workflow/assets/sdd-workflow.config.json` (if present) for its `archive_spec` section's
`destination_type` (`"mcp"` | `"skill"` | `"path"` | `"none"`) and `destination_name`. Default to
`destination_type: "none"`, `destination_name: ""` if the file, the section, or either field is
missing.

- If `destination_type` is `mcp`: `destination_name` is the **exact** tool name to call, including
  the full `mcp__{server}__{tool}` prefix (e.g. `mcp__obsidian__create_note`) — it is resolved by
  the user's config, not guessed here. Do not search for other tools under the same server prefix
  and do not substitute a different tool if this exact one is missing.
- If `destination_type` is `skill`: look for a skill named `destination_name` in your available
  skills and invoke it to perform the write — it may implement its own rules or a custom send
  mechanism instead of calling an MCP tool directly. Follow whatever inputs that skill expects for
  title/content per note. Also ask it to maintain the per-module index (see step 2 under "Steps —
  `destination_type` is `skill`" below) — unlike `mcp`, a skill is a full agent and can read-modify-
  write its own destination.
- If `destination_type` is `path`: `destination_name` is a local filesystem folder (absolute or
  relative to the repo root), outside `docs/` — e.g. a personal notes directory that isn't part of
  the codebase. No MCP tool or skill is involved; you write files directly. Create the folder (and
  any missing parent directories) if it doesn't exist yet.
- If `destination_type` is `none`: do not attempt any destination. Stop and report
  `STATUS: no-destination-available` immediately, without touching any local file.

If the resolved `destination_type`/`destination_name` names an MCP tool or skill that is not
actually available in this session, do not fall back to guessing another mechanism or silently
skipping. Stop and report `STATUS: no-destination-available` (see Final report) so the orchestrator
can inform the user; do not delete any local file in that case. This does not apply to `path` — a
missing folder is created, never treated as unavailable.

## Files to archive

Feature: {feature_path}

- PRD:   {prd_path}
- PLAN:  {plan_path}
- NOTES: {notes_path}   <!-- omit if NOTES.md does not exist for this feature -->

`STATE.md` was already deleted by the orchestrator before you were spawned and will not be in
`{feature_path}`. It is never one of the files to archive — if it is unexpectedly still present, do
not send it or move it.

`{prd_path}` may point to `PRD.cache.md` instead of `PRD.md` — `sdd-plan` materializes a PRD
sourced from an external tool (Jira, GitHub Issues, etc.) under that name, marking it as a
temporary, regenerable local cache rather than the authored source of truth (the external card is).
That distinction only matters inside this repo. Once archived, this **is** the permanent record of
what shipped, so treat it as a plain PRD regardless of the suffix — never propagate the `.cache.`
segment past this agent.

At the destination, the PRD is always renamed to carry the feature's own name:
`PRD-{feature_slug}.md`, where `{feature_slug}` is the feature's folder name (the last segment of
`{feature_path}`, e.g. `access-control-v1` for `docs/user/access-control-v1`) — this applies
whether the source was `PRD.md` or `PRD.cache.md`. Every reference to this PRD from elsewhere in
the archive (the module file, PLAN/NOTES `Correlations` links) must use this qualified name, never
a bare `PRD.md`. `PLAN.md` and `NOTES.md` keep their plain names — the folder they live in already
disambiguates them, and they're only ever reached by navigating from the PRD, one hop away.

## Steps — `destination_type` is `mcp`

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
6. Leave every other file under `{feature_path}` untouched (`STATE.md` is already gone — see "Files
   to archive" above). Do not attempt any module index — see the `mcp` limitation noted above.

If any note fails to write (tool error, timeout, auth failure), stop before deleting anything local
and report which files were and weren't archived.

## Steps — `destination_type` is `skill`

1. Read each file above in full.
2. Invoke the skill named `destination_name`, asking it to:
   - Create one note per file, using each file's own content verbatim (do not summarize or
     rewrite it), titled so it's identifiable without the repo path (e.g. `{feature_path} — PRD`).
   - Link the notes to each other (PRD ↔ PLAN ↔ NOTES) using whatever native linking the
     destination supports.
   - Create or update one module-level index entry for `{module}` (the first segment of
     `{feature_path}` after stripping `docs/`), grouped by type, following the same content rules
     as the `{module}.md` spec under "Steps — `destination_type` is `path`" step 4 below (type
     mapping, one line per feature linking only to its PRD note, module summary preserved on
     update) — adapted to whatever native format the skill's destination uses instead of a flat
     file.
   Follow whatever inputs that skill expects for title/content per note; it owns the actual write
   mechanism.
3. Confirm every note (including the index update) was actually written — check the skill's
   response, not just the absence of an error — before touching any local file.
4. Only after every note in step 2 is confirmed written: delete the local files listed above
   (`{prd_path}`, `{plan_path}`, and `{notes_path}` if present) from the repo working tree.
5. Leave every other file under `{feature_path}` untouched (`STATE.md` is already gone — see "Files
   to archive" above).

If any note fails to write (tool error, timeout, auth failure), including the index, stop before
deleting anything local and report which files were and weren't archived.

## Steps — `destination_type` is `path`

No vault, no linking format, no external tool — this is a plain filesystem move to
`destination_name`, outside the repo's version control.

1. Strip the leading `docs/` segment from `{feature_path}` to get `{feature_rel_path}` (e.g.
   `docs/user/login` → `user/login`) — `destination_name` is itself the new docs root, so keeping
   `docs/` would nest it redundantly (`{destination_name}/docs/...`). Split `{feature_rel_path}`
   into `{module}` (its first segment) and `{feature_slug}` (everything after it, e.g. `user/login`
   → module `user`, feature slug `login`). Create `{destination_name}/{feature_rel_path}` (mkdir -p
   style) if it doesn't already exist, so the feature's module/feature folder structure is preserved
   at the destination and multiple features never collide.
2. Copy each file above into `{destination_name}/{feature_rel_path}`:
   - `{prd_path}` → `PRD-{feature_slug}.md` (this is where any `.cache` segment gets dropped — see
     "Files to archive" above).
   - `{plan_path}` → `PLAN.md`, `{notes_path}` → `NOTES.md` (unchanged names) — but before writing
     either, rewrite its `Correlations` section's PRD link (`[PRD](./PRD.md)` or
     `[PRD](./PRD.cache.md)`) to `[PRD](./PRD-{feature_slug}.md)`. Leave every other link in that
     section untouched. This keeps the destination copies internally consistent even though the
     local originals still point at the old filename.
3. Confirm each copy landed at the destination (file exists, non-empty, same byte size as the
   source for `PLAN.md`/`NOTES.md`; non-empty for the renamed PRD) before touching the source.
4. Update the module's `{module}.md` (only after step 3 is confirmed, only for `path`):
   1. Read the copied `PRD-{feature_slug}.md` (fall back to `PLAN.md` if the PRD lacks it) for its
      `**Tipo**` / `**Type**` field and its title (the first `# ` heading). Map the type to a
      section name: `feat`→`Features`, `fix`→`Fixes`, `refactor`→`Refactors`, `perf`→`Performance`,
      `chore`→`Chores`, `docs`→`Docs`, `test`→`Tests`, `build`→`Build`, `ci`→`CI`. If the type is
      missing or unrecognized, use `Features`.
   2. Read `{destination_name}/{module}/{module}.md` if it exists.
      - **Missing:** create it with a `# {module}` heading, then a short one-paragraph module
        summary — take the first paragraph of the archived PRD's "Visão Geral"/overview section; if
        that section is empty or absent, write `_(resumo do módulo pendente)_` instead. Do not
        invent a summary from unrelated content.
      - **Existing:** keep its heading and summary paragraph untouched — you're only adding/updating
        one entry in one section below.
   3. Under a `### {SectionName}` heading (create the heading, in this fixed order if multiple exist
      — Features, Fixes, Refactors, Performance, Chores, Docs, Tests, Build, CI — skipping any
      section that has no entries), add one line for this feature, linking only to its PRD (PLAN and
      NOTES are one hop away via the PRD's own `Correlations` section, so this file doesn't
      duplicate those links):
      `- {Feature Title}: [PRD](./{feature_slug}/PRD-{feature_slug}.md)`
      — path relative to the module folder, matching the structure created in step 1. If a line for
      this exact feature already exists (re-running archival on the same feature), replace that line
      instead of duplicating it.
   4. Write the updated `{module}.md` back to `{destination_name}/{module}/{module}.md`.
5. Only after every copy in step 2 is confirmed **and** `{module}.md` is updated: delete the local
   files listed above (`{prd_path}`, `{plan_path}`, and `{notes_path}` if present) from the repo
   working tree.
6. Leave every other file under `{feature_path}` untouched (`STATE.md` is already gone — see "Files
   to archive" above).

If any copy fails (permission error, disk full, path unwritable) or the `{module}.md` update fails,
stop before deleting anything local and report which files were and weren't archived.

## Final report

Report exactly one of:

- `STATUS: archived` — every artifact was written to the destination and the corresponding local
  files were deleted. List the destination note titles created (`mcp`/`skill`) or the destination
  file paths (`path`), including the updated `{module}.md` path/note for `path` and `skill`
  destinations.
- `STATUS: partial-failure` — some artifacts were written, others failed. List which succeeded
  (and were deleted) and which failed (and were left in place, untouched). Include the raw error
  for each failure.
- `STATUS: no-destination-available` — no Obsidian (or other archiver) MCP tool/skill was found, or
  `destination_type` was `none`. No files were touched. (Not applicable to `path` — see above.)
```

---

## Variables to replace

| Placeholder | Source |
|-------------|--------|
| `{feature_path}` | the path passed to `sdd-workflow` (e.g. `docs/user/login`) |
| `{prd_path}` | `STATE.md.artifacts.prd`, read before `STATE.md` is deleted |
| `{plan_path}` | `STATE.md.artifacts.plan`, read before `STATE.md` is deleted |
| `{notes_path}` | `STATE.md.artifacts.notes`, read before `STATE.md` is deleted — omit the line entirely if empty |

## What the orchestrator does with the report

`STATE.md` is already deleted by the time this agent is spawned (see phase-done.md) — there is
nothing left to clean up on that front regardless of outcome.

- `archived`: remove `{feature_path}` if it is now empty.
- `partial-failure` or `no-destination-available`: leave every remaining local file (`PRD.md`,
  `PLAN.md`, `NOTES.md`) in place — the pipeline-concluído summary was already shown to the user in
  this same phase, so just append a short note that archiving did not complete and why. Because
  `STATE.md` is gone, this feature can no longer be resumed automatically by a later
  `/sdd-workflow {feature_path}` — say so explicitly, and point the user at re-running the
  archival manually once the blocker is resolved. This is not a pipeline failure — `phase` stays
  `done`.
