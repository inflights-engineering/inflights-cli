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
