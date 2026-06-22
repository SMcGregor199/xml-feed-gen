# Agent Instructions for xml-feed-gen

Go RSS generator for `shaynemcgregor.dev`.

## Read Before Work

- `../system_context.md`
- `../WORKFLOW.md`
- `../DECISIONS.md`
- `README.md`

## Repo-Specific Cautions

- Consumes backend blog JSON from `notion2github-BE`.
- Generates production-ready RSS XML used by the production RSS workflow.
- `notion2github-FE/public/rss.xml` is the tracked rollback fallback RSS artifact.
- `cmd/rssgen/rss.xml` is tracked prototype output and is not authoritative.
- Not Netlify-configured.
- Do not inspect `.env` or secret files.
- Do not run `go run ./cmd/rssgen` unless you intend to refresh the tracked RSS artifact.

## Verification Expectations

- For source changes, prefer `go test ./...`.
- For documentation-only changes, read back the edited docs and check `git status --short`.
- Report any skipped verification command and why it was skipped.

## Cleanup Expectations

- Leave user-owned dirty files untouched.
- Do not commit generated output unless the task explicitly targets the RSS artifact.
- Do not modify `.gitignore` to ignore `AGENTS.md`.

## Forward-Look Note

- `xml-feed-gen` has been promoted into production-ready generation tooling. Production publishing and Netlify routing changes still require explicit approval.

## Final Handoff Expectations

- State which files changed.
- State which commands ran and which were skipped.
- Call out any RSS sync implications with `notion2github-FE`.
- Do not create commits, branches, PRs, releases, or tags unless Shayne asks.
