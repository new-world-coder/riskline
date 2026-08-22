# riskline

[![CI](https://github.com/new-world-coder/riskline/actions/workflows/ci.yml/badge.svg)](https://github.com/new-world-coder/riskline/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/new-world-coder/riskline)](LICENSE)
[![Release](https://img.shields.io/github/v/release/new-world-coder/riskline?include_prereleases)](https://github.com/new-world-coder/riskline/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/new-world-coder/riskline.svg)](https://pkg.go.dev/github.com/new-world-coder/riskline)
[![Go Report Card](https://goreportcard.com/badge/github.com/new-world-coder/riskline)](https://goreportcard.com/report/github.com/new-world-coder/riskline)

Open-source AI risk classification — CLI and API, deterministic and auditable, not an LLM guessing.

**Status:** Pre-1.0 (`0.2.0-alpha`). APIs and rulesets may change. Feedback via [Issues](https://github.com/new-world-coder/riskline/issues) is welcome — especially [misclassification reports](https://github.com/new-world-coder/riskline/issues/new?template=misclassification_report.md).

You describe an AI system (purpose, data types, deployment context, autonomy, who it affects). You get back a risk tier — prohibited, high-risk, limited-risk, or minimal-risk — with clause references, a plain-language rationale, and recommended controls. The ruleset is versioned JSON, not a model making something up.

**Shipped regime packs:** `eu-ai-act` (`eu-ai-act-2024-v0.1.0`, hard law) and `nist-ai-rmf` (`nist-ai-rmf-2023-v0.1.0`, mapping — not US legal tiers). Multi-regime via `regimes[]`, `.riskline.yaml`, `RISKLINE_REGIMES`. MAS packs are next — not shipping claims yet. `geographic_scope` is not a regime selector.

Distribution model borrows from Open Policy Agent and Trivy — single binary, no runtime, CI-friendly.

**Website:** [new-world-coder.github.io/riskline](https://new-world-coder.github.io/riskline/) · [Printable one-pager](https://new-world-coder.github.io/riskline/one-pager.html) · [Console demo](https://riskline-cloud-web.vercel.app/console) (mock GRC dashboards)

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
./bin/riskline-cli -regimes eu-ai-act --json examples/curl/hiring-assist.yaml
./bin/riskline-cli -list-regimes
```

Optional project defaults: copy [`examples/riskline.yaml`](examples/riskline.yaml) to `.riskline.yaml`, or set `RISKLINE_REGIMES=eu-ai-act`.

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



## Write-ups

- [How we classify AI systems under the EU AI Act](docs/blog/01-how-we-classify-ai-systems-under-the-eu-ai-act.md)
- [An open JSON schema for AI system risk classification](docs/blog/02-open-json-schema-for-ai-system-risk-classification.md)
- [Why I built the compliance engine in Go](docs/blog/03-why-go-for-a-compliance-engine.md)
- [EU AI Act classification belongs in CI — not another portal](docs/blog/04-eu-ai-act-in-ci-not-a-portal.md)
- [Singapore AI governance meets deterministic EU classification](docs/blog/05-singapore-ai-governance-meets-deterministic-classification.md)
- [Introducing the Riskline Console](docs/blog/06-introducing-riskline-console.md) — role-based trust dashboards (demo; mock data)

See also [ROADMAP.md](ROADMAP.md) and ready-to-post copy in [`content/social/LAUNCH_COPY.md`](content/social/LAUNCH_COPY.md) / [`content/social/OUTREACH_EU_SG.md`](content/social/OUTREACH_EU_SG.md).

## Disclaimer

Every response includes this, and we keep a test that checks it stays put:

> This classification is an advisory tool based on a versioned ruleset. It is not legal advice and is not a substitute for qualified counsel or a formal conformity assessment.

## License

Apache 2.0 — see [LICENSE](LICENSE).
