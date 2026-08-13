# Examples

- `curl/` — API request bodies
- `riskline.yaml` — sample project config (copy to `.riskline.yaml`)
- `github-action/` — CI classification action
- `pre-commit/` — local hook sketch
- `../testdata/scenarios/` — fixture systems used in docs/blog (prohibited, high-risk, limited-risk, credit boundary)

```bash
go run ./cmd/riskline-cli testdata/scenarios/high-risk-hiring.yaml
go run ./cmd/riskline-cli -list-regimes
```
