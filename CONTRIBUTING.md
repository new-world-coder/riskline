# Contributing

Thanks for wanting to help. This project stays useful by staying small and honest about what the ruleset covers.

## Ground rules

1. **EU AI Act only for now.** Don't open PRs that paste NIST/ISO control text into the repo — ISO/IEC clause text is copyrighted and explicitly out of scope. Reference clause numbers and paraphrase in your own words if you ever extend beyond the AI Act.
2. **Rules live in JSON**, not buried in Go `if` trees. Add or edit files under `pkg/ruleset/data/`, bump `version` / `last_updated`, and note judgment calls in the rule's `judgment_call_note`.
3. **Tests travel with the change.** Table rows in `pkg/engine` for new scenarios; golden files under `testdata/golden` when the public JSON shape matters.
4. **Engine stays pure.** `pkg/engine` must not import HTTP or CLI packages.

## Dev loop

```bash
go test ./... -race -cover
go vet ./...
go build -o bin/riskline-cli ./cmd/riskline-cli
go build -o bin/riskline-api ./cmd/riskline-api
```

## Issues we actually want

Use the templates:

- **Misclassification report** — you think the tier is wrong for a real system. This is our validation signal; we read these weekly.
- **Ruleset update** — the Act / guidance moved, or a mapping looks stale.

No CLA for now. By contributing you agree your work is Apache 2.0 licensed.

## Code of conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
