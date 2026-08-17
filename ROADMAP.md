# Roadmap

Public signal of direction. If it is not listed here, assume it is not being built yet.

## Now (v0.3.x)

- Multi-regime foundation: pack loader, `regimes[]`, `.riskline.yaml` / `RISKLINE_REGIMES`
- EU AI Act pack (`eu-ai-act` → `eu-ai-act-2024-v0.1.0`)
- **NIST AI RMF mapping pack** (`nist-ai-rmf` → `nist-ai-rmf-2023-v0.1.0`) with **technical controls**
- Offline CLI (`riskline-cli`) + thin API (`riskline-api`)
- Versioned embedded ruleset with judgment-call notes
- Misclassification + ruleset-update issue templates as the validation loop

## Next (when design-partner / issue signal says so)

- Deeper Annex III coverage, one well-tested category at a time ([#6](https://github.com/new-world-coder/riskline/issues/6))
- Stronger Art. 50 transparency detection beyond chatbot keywords
- Generated TypeScript / Python clients kept in sync with OpenAPI
- Real OpenAPI contract tests in CI ([#4](https://github.com/new-world-coder/riskline/issues/4))
- MAS FEAT pack (`mas-feat`) — principles/control gaps, not EU tiers
- MAS AIRG pack (`mas-airg`) after supervisory text stabilizes
- Deeper NIST subcategory coverage from design-partner feedback

## Planned (not started)

- Evidence generator
- Executive / GRC dashboard
- Third-party / vendor risk portal
- Hosted multi-tenant SaaS + billing

## Explicitly out of scope for now

- Embedding ISO/IEC clause text
- Claiming "enterprise-ready" or publishing a price list before design-partner feedback
- Collapsing non-EU regimes into prohibited/high/limited/minimal
- Inferring regimes from install locale or `geographic_scope`
