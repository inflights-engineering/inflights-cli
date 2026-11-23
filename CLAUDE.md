# inflights-cli

Go CLI for the Inflights platform, built with Cobra.

## Build & Test

```bash
go build ./...
go test ./... -v
go vet ./...
```

## Conventions

- Commit messages: terse, lowercase, no period
- Use "service" in all CLI-facing text — never "product" (backend uses "product" internally but the CLI must not expose that term)
- Standard Go project layout: commands in `cmd/`, internal logic in `internal/`

## On Elias's Mac: builds and tests run on the Beelink

On Elias's Mac, never install dependencies, build, test or start servers locally. Run them on
the Beelink build box with the wrappers in `~/.local/bin`, which sync this checkout or
worktree to its own copy there and return the exit code:

- `rbuild`: `go build ./...`. `rbuild "<cmd>"` runs any command there.
- `rtest`: `go vet ./...`, then `go test ./...` (4 packages at a time). A failing `rtest` means the tests failed.
- One package or test: `rtest "go test ./internal/<pkg> -run <TestName> -v"`.
- **Commands that rewrite files run locally**, because the next sync overwrites the remote
  copy: `gofmt -w` / `go fmt`, `go mod tidy`, `go get`.
- If the Beelink is unreachable, stop and say so. Never fall back to local builds.
- Several agents work at once, each in its own worktree with one branch and one PR. Stay in
  yours, never touch another worktree's files, and merge only as this file allows.

How the Beelink and the wrappers are set up: `~/42istheanswer/beelink/BEELINK.md`.
