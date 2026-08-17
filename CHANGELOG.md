# Changelog

## 0.3.0-alpha — 2026-08-17

- **NIST AI RMF mapping pack** (`nist-ai-rmf` → `nist-ai-rmf-2023-v0.1.0`): Govern/Map/Measure/Manage subcategory mappings — not fake US legal tiers.
- **Technical controls compiler**: rules may emit structured `{ paper_ref, technical_hook, evidence_type }` on responses.
- **`mapping_only`** flag on mapping-pack outcomes so UIs do not misread `risk_tier` as US law.
- **`geographic_scope`**: added `us` and `us_and_global` request hints.
- EU-only default responses remain golden-compatible.

## 0.2.0-alpha — 2026-08-13

- Multi-regime foundation (P0): regime pack loader, `regimes[]` on classify requests, optional `classifications[]` when more than one pack is evaluated.
- Project config: `.riskline.yaml` and `RISKLINE_REGIMES` (plus CLI/API `-regimes`).
- Shipped pack remains `eu-ai-act` → ruleset `eu-ai-act-2024-v0.1.0`. Unknown regimes are rejected; no MAS/NIST packs yet.
- OpenAPI `info.version` bumped to `0.2.0-alpha`. EU-only default responses stay golden-compatible (no additive fields).

## 0.1.0-alpha — 2026-08-04

- Initial EU AI Act ruleset `eu-ai-act-2024-v0.1.0` (prohibited Art. 5 subset, prioritized Annex III high-risk categories, Art. 50 limited-risk heuristics).
- Classification engine, CLI, and `POST /v1/classify` API.
- Judgment-call notes recorded on boundary rules (biometric LE exceptions, fraud carve-out, Art. 50 breadth, etc.).
