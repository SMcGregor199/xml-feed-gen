# xml-feed-gen

Go RSS generator for `shaynemcgregor.dev`.

## What This Repo Does

- Fetches blog JSON from `notion2github-BE`.
- Converts backend post objects into RSS 2.0 XML.
- Writes `rss.xml` from the command-line generator.
- Tests RSS XML generation.

## Commands

- `go test ./...`: run tests.
- `go run ./cmd/rssgen`: fetch backend JSON and write `rss.xml`.
- `go run .`: fetch backend JSON and print it.

## Local Development Notes

- This repo is not Netlify-configured.
- It consumes the backend Netlify-hosted JSON endpoint.
- `go run ./cmd/rssgen` writes a tracked artifact, so only run it when you intend to refresh RSS output.

## Artifact And Generated File Cautions

- `cmd/rssgen/rss.xml` is tracked and should be changed intentionally.
- Treat any regenerated RSS file as a release artifact, not a routine temp file.
- Avoid committing generated output unless you are intentionally updating the public feed.

## Relation To The Other Repos

- Consumes blog JSON produced by `notion2github-BE`.
- Shares the blog post contract with `notion2github-BE` and `notion2github-FE`.
- Produces RSS XML used by the broader personal site system.
