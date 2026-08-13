# Changelog

## 0.2.0-alpha — 2026-08-13

- Multi-regime foundation (P0): regime pack loader, `regimes[]` on classify requests, optional `classifications[]` when more than one pack is evaluated.
- Project config: `.riskline.yaml` and `RISKLINE_REGIMES` (plus CLI/API `-regimes`).
- Shipped pack remains `eu-ai-act` → ruleset `eu-ai-act-2024-v0.1.0`. Unknown regimes are rejected; no MAS/NIST packs yet.
- OpenAPI `info.version` bumped to `0.2.0-alpha`. EU-only default responses stay golden-compatible (no additive fields).

## 0.1.0-alpha — 2026-08-04

- Initial EU AI Act ruleset `eu-ai-act-2024-v0.1.0` (prohibited Art. 5 subset, prioritized Annex III high-risk categories, Art. 50 limited-risk heuristics).
- Classification engine, CLI, and `POST /v1/classify` API.
- Judgment-call notes recorded on boundary rules (biometric LE exceptions, fraud carve-out, Art. 50 breadth, etc.).
