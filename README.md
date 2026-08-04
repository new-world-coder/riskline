# riskline

Open-source EU AI Act risk classification — CLI and API, deterministic and auditable, not an LLM guessing.

You describe an AI system (purpose, data types, deployment context, autonomy, who it affects). You get back a risk tier — prohibited, high-risk, limited-risk, or minimal-risk — with clause references, a plain-language rationale, and recommended controls. The ruleset is versioned JSON, not a model making something up.

**v1 does EU AI Act only.** No NIST, no ISO mappings, no dashboards, no evidence packs, no billing. If you need those, they're roadmap — not shipping claims.

## Why this exists

Most "AI governance" tools want you in a portal. Engineering teams already live in CI and the terminal. Classification that can't run offline, can't be diffed, and can't explain *which article fired* is hard to trust — especially when the output looks like legal judgment.

This is the wedge: a small, embeddable engine + CLI + HTTP API with the same contract.

## Quick start (CLI)

```bash
go install github.com/new-world-coder/riskline/cmd/riskline-cli@latest
# or from this repo:
go build -o bin/riskline-cli ./cmd/riskline-cli

./bin/riskline-cli examples/curl/hiring-assist.yaml
./bin/riskline-cli --json examples/curl/hiring-assist.yaml
```

The CLI reads local files only. It does not phone home. See [PRIVACY.md](PRIVACY.md).

## API

```bash
go run ./cmd/riskline-api -addr :8080
curl -s localhost:8080/v1/classify -H 'content-type: application/json' -d @examples/curl/hiring-assist.json
```

Contract: [`api/openapi.yaml`](api/openapi.yaml).

## What you should not expect yet

- Full Annex III coverage (v1 prioritises recruitment, employment, credit, biometrics, law-enforcement risk scoring, plus a couple of adjacent cases)
- "Enterprise-ready" anything, customer logos, or published prices
- Hosted multi-tenant SaaS (stubs marked `TODO(hosted):`)

Early access for a hosted endpoint: open a thread in [GitHub Discussions](https://github.com/new-world-coder/riskline/discussions) — no public price list on purpose.

## Disclaimer

Every response includes this, and we keep a test that checks it stays put:

> This classification is an advisory tool based on a versioned ruleset. It is not legal advice and is not a substitute for qualified counsel or a formal conformity assessment.

## License

Apache 2.0 — see [LICENSE](LICENSE).
