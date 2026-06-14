# Contributing to cadence

Thanks for wanting to help. Here's everything you need to know.

## Setup

You need Go 1.25+ and optionally [Ollama](https://ollama.com/download) for embedding-related work.

```bash
git clone https://github.com/ibnaleem/cadence
cd cadence
go build ./...
go vet ./...
```

If you're working on anything that touches `add`, `done`, or `edit` with name-based lookup, pull the embedding model:

```bash
ollama pull embeddinggemma
```

## Project structure

```
main.go                  Entry point, prints the banner
cmd/                     One file per command (add, done, delete, edit, list, streak, undo)
internal/util/
  habits.go              All DB operations (CRUD, streak, weekly, fuzzy match)
  setup.go               DB init and schema migrations
  embed.go               Ollama client and cosine similarity
  util.go                CheckError helper
internal/theme/          ANSI colour helpers
internal/tui/            Scrollable viewport (used for long output)
```

Commands register themselves in `cmd/root.go`'s `init()`. DB functions live in `internal/util/habits.go`. Keep that separation — commands should only call util functions, never touch the DB directly.

## Making changes

Branch off `main`:

```bash
git checkout -b feat/your-thing   # or fix/, docs/, refactor/
```

Commit after every discrete change. Not after every file, not in one big dump at the end. A new command, a new DB helper, a schema change — each gets its own commit. Present tense messages:

```
feat: add pause command to suspend a habit without losing history
fix: prevent streak reset when habit logged before midnight UTC
docs: add --date flag to done command in README
```

Run `go build ./...` after every commit. Do not move on if it fails.

Before opening a PR, run:

```bash
go build ./...
go vet ./...
go test ./...
```

## Adding a command

1. Create `cmd/<name>.go` with a `var <name>Cmd = &cobra.Command{...}`
2. Add any DB helpers it needs to `internal/util/habits.go`
3. Register it in `cmd/root.go` `init()` with `rootCmd.AddCommand(<name>Cmd)`
4. Commit the DB helper and the command separately

## Schema changes

Add new columns via `ALTER TABLE` in `SetupSchema` in `internal/util/setup.go`, after the `CREATE TABLE IF NOT EXISTS` block. SQLite has no `ADD COLUMN IF NOT EXISTS`, so ignore the error:

```go
db.Exec(`ALTER TABLE habits ADD COLUMN paused INTEGER NOT NULL DEFAULT 0`)
```

If you're adding a `UNIQUE` constraint or index that existing databases won't have, use `CREATE UNIQUE INDEX IF NOT EXISTS` rather than baking it into the original `CREATE TABLE` (which only runs once on new installs).

## Pull requests

- One PR per feature or fix. Don't bundle unrelated changes.
- Keep the PR description short: what changed and why.
- Base on `main`. If your branch conflicts, rebase or merge `main` in and resolve before asking for review.
- CI runs `go build`, `go vet`, and `go test` on push. Make sure all three pass before marking ready.

## What to work on

Check the open issues. Anything labelled `P1` or `P2` is highest priority. If you want to work on something not listed, open an issue first so we can discuss scope before you build it.

## Code style

- Follow standard Go conventions (`gofmt`, `go vet` clean)
- No unnecessary comments — if the function name explains it, don't add a docstring
- Error messages are lowercase, no trailing punctuation: `"no habit with id %d"` not `"No habit with ID %d."`
- Use `theme.*` helpers for all terminal output — never hardcode ANSI codes
- DB functions return errors, not booleans for "not found" — use `sql.ErrNoRows`

## Questions

Open an issue. Don't email.
