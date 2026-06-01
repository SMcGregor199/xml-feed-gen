# Agent Instructions for xml-feed-gen

This repo is the Go RSS generator for `shaynemcgregor.dev`.

Read the parent workspace `system_context.md` before making cross-repo or data-contract changes.

## Repo Purpose

- Fetch blog post JSON from `notion2github-BE`.
- Convert backend post objects into RSS 2.0 XML.
- Write `rss.xml` from the command-line generator.
- Test RSS XML generation.

## Important Directories And Files

- `go.mod`: Go module definition.
- `index.go`: small debug-style entry point that prints fetched backend JSON.
- `cmd/rssgen/main.go`: command-line generator that writes `rss.xml`.
- `feed/fetch.go`: backend JSON fetch logic.
- `feed/rss.go`: RSS types, ordering, URL handling, date formatting, and XML writing.
- `feed/rss_test.go`: smoke test for RSS XML output.
- `.gitignore`: ignore rules. Do not add `AGENTS.md`.

## Relevant Commands

- `go test ./...`: run tests.
- `go run ./cmd/rssgen`: fetch backend JSON and write `rss.xml`. Confirm before running because it writes an artifact.
- `go run .`: fetch backend JSON and print it. Avoid if output may expose sensitive or unexpected production content.

## External Services Used

- Fetches `https://shaynemcgregordev-be.netlify.app/.netlify/functions/blog-posts-json`.
- Emits links for `https://shaynemcgregor.dev`.
- Uses backend `notion-image` URLs as RSS enclosures when thumbnails are present.
- Does not directly use Notion, Airtable, Netlify APIs, or GitHub APIs in inspected source.

## Connections To Other Repos

- Consumes blog JSON produced by `notion2github-BE`.
- Shares the blog post data contract with `notion2github-BE` and `notion2github-FE`.
- Produces RSS XML for the public site rendered by the broader personal site system.

## Inspect Before Changing

- `git status --short`
- `.gitignore`
- `.git/info/exclude`
- `go.mod`
- `feed/fetch.go`
- `feed/rss.go`
- `feed/rss_test.go`
- `cmd/rssgen/main.go`
- parent `../system_context.md` for cross-repo contract context

Do not inspect `.env` files or secret files.

## Do Not Do Without Confirmation

- Do not deploy.
- Do not regenerate committed RSS artifacts unless requested.
- Do not change the expected backend JSON contract casually.
- Do not modify `.gitignore` to ignore `AGENTS.md`.
- Do not restore hidden duplicate context files such as `.AGENTS.md`.
- Do not create commits, branches, PRs, releases, or tags.
- Do not overwrite user-owned dirty changes.

## Testing, Build, And Formatting Guidance

- Use `go test ./...` for verification when Go source changes.
- Use standard Go formatting expectations for Go source changes.
- For documentation-only changes, read back changed docs and check `git status --short`.
- If skipping tests, report why.

## Data Contract Warnings

`feed.Post` mirrors the backend blog post contract:

- `id`
- `tag`
- `title`
- `summary`
- `link`
- `thumbnail`
- `publishedDate`
- `updatedDate`
- `body` as sections with `heading` and `paras`

RSS generation assumes:

- `link` is absolute or a slug relative to `https://shaynemcgregor.dev`.
- `publishedDate` is RFC3339-compatible when possible.
- non-empty `thumbnail` values are WebP image enclosures.

Coordinate contract changes with `notion2github-BE` and `notion2github-FE`.

## Secrets And Environment Warnings

- Do not read `.env` or secret files.
- Do not print secrets or environment values.
- No direct environment variable reads were observed in inspected Go source.
- Fetching backend JSON may contact a production endpoint; confirm intent before commands that fetch or write artifacts.

## PR Workflow Expectations

- Do not create a branch, commit, or PR unless Shayne asks.
- Before a PR, summarize dirty files and identify user-owned changes.
- Include verification commands and skipped commands in any PR summary.
- Mention RSS output changes and backend contract assumptions if relevant.
