# xml-feed-gen

Go RSS generator for `shaynemcgregor.dev`.

## What This Repo Does

- Fetches blog JSON from `notion2github-BE`.
- Converts backend post objects into RSS 2.0 XML with validation and deterministic output support.
- Writes `rss.xml` atomically from the command-line generator.
- Tests RSS XML generation and fetch behavior.

## Commands

- `go test ./...`: run tests.
- `go run ./cmd/rssgen`: fetch backend JSON and write `rss.xml`.
- `go run .`: fetch backend JSON and print it.

## Local Development Notes

- This repo is not Netlify-configured.
- It consumes the backend Netlify-hosted JSON endpoint.
- This repo has been promoted into production-ready RSS generation tooling.
- The generator validates required post fields, uses `/blog/{slug}` canonical links, produces stable GUIDs, sorts newest-first, and writes output atomically.
- `go run ./cmd/rssgen` writes `cmd/rssgen/rss.xml` by default when run from `cmd/rssgen`; that tracked artifact remains prototype output and is not authoritative for the public site.
- The public live RSS URL is served from the backend generated feed through the frontend `/rss.xml` route/proxy.
- `notion2github-FE/public/rss.xml` remains the tracked rollback fallback artifact.
- Non-secret configuration names used by RSS tooling include `RSS_OUTPUT_PATH`, `RSS_SITE_BASE_URL`, and `RSS_PUBLIC_URL`.

## Artifact And Generated File Cautions

- `cmd/rssgen/rss.xml` is tracked prototype output and should be changed intentionally.
- It is not the authoritative live RSS source.
- Avoid committing generated output unless you are intentionally working on the prototype tooling.

## Relation To The Other Repos

- Consumes blog JSON produced by `notion2github-BE`.
- Shares the blog post contract with `notion2github-BE` and `notion2github-FE`.
- Produces RSS XML for production-ready generation workflows.
- Production publishing is coordinated by `notion2github-BE` through the RSS auto-update path when enabled by approved environment configuration.
