# AGENTS.md

## Cursor Cloud specific instructions

`riskline` is a zero-config, pure Go product (EU AI Act risk classifier) with three surfaces sharing one embedded JSON ruleset: `riskline-cli` (one-shot CLI), `riskline-api` (HTTP server), and placeholder client stubs. There is no database, cache, queue, `.env`, or external service — the ruleset is embedded via `//go:embed`, so nothing external needs to be running.

Standard build/lint/test/run commands are already documented; do not duplicate them. See `CONTRIBUTING.md` ("Dev loop"), `README.md` ("Quick start" / "API"), and `.github/workflows/ci.yml` (authoritative CI: build, `go vet`, `go test -race`, ≥85% coverage gate on `pkg/engine` and `pkg/ruleset`, `staticcheck`, `gosec`, OpenAPI smoke test).

Non-obvious caveats for future agents:

- Go toolchain: `go.mod` pins `go 1.24.0`. `go` auto-downloads the required toolchain on first use (default `GOTOOLCHAIN=auto`), so the first command can pause to fetch. Installing `staticcheck@latest` pulls a module requiring Go >= 1.25, which triggers an automatic switch to a `go1.25.x` toolchain — expected, not an error.
- Run the API: `go run ./cmd/riskline-api -addr :8080` (only long-running service; port `8080`, override via `-addr`). Endpoints: `GET /healthz` (returns `ok`) and `POST /v1/classify`. Sample payloads live in `examples/curl/` (`hiring-assist.json` for the API, `hiring-assist.yaml` for the CLI).
- `staticcheck` and `gosec` are installed on demand via `go install` (not vendored); they are not part of the update script.
