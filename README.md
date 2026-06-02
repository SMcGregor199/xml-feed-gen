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
- This repo is experimental/prototype tooling and is not part of the live RSS workflow right now.
- `go run ./cmd/rssgen` writes `cmd/rssgen/rss.xml`, which is prototype output and not authoritative for the public site.
- The live RSS artifact is `notion2github-FE/public/rss.xml`.
- This repo could be promoted into the live workflow later, but only through an explicit decision.

## Artifact And Generated File Cautions

- `cmd/rssgen/rss.xml` is tracked prototype output and should be changed intentionally.
- It is not the authoritative live RSS source.
- Avoid committing generated output unless you are intentionally working on the prototype tooling.

## Relation To The Other Repos

- Consumes blog JSON produced by `notion2github-BE`.
- Shares the blog post contract with `notion2github-BE` and `notion2github-FE`.
- Produces prototype RSS XML used for experimentation, not the live public feed.
