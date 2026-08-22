# Roadmap

Public signal of direction. If it is not listed here, assume it is not being built yet.

## Now (v0.5.x)

- Multi-regime foundation: pack loader, `regimes[]`, `.riskline.yaml` / `RISKLINE_REGIMES`
- EU AI Act pack (`eu-ai-act` → `eu-ai-act-2024-v0.1.0`) with **technical_controls** on key high-risk rules
- **Assure layer (v0.4-alpha):** `POST /v1/assure`, `POST /v1/diff`, conformity state, material change, SHA-256 evidence chain
- **CI gate (v0.5-alpha):** GitHub Action `riskline-assure` — classify → diff → assure in PR pipelines
- **Signed evidence (v0.5-alpha):** Ed25519 bundles via CLI `--sign` / `verify` (local, not blockchain)
- **Runtime metadata (v0.5-alpha):** model_id, system_prompt_hash, tools, human_approval_required in material-change fingerprints
- NIST AI RMF mapping pack with technical controls
- Offline CLI: `classify`, `diff`, `assure`, `verify` subcommands
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

- **Runtime assurance (Phase 2–3 done):** `pkg/runtime` engine + `POST /v1/runtime/{register,verify,observe}` handlers; signed `VerificationReceipt` (Phase 4 next)
- **Integrations (Phase 1 done):** `integrations/n8n` scaffold (Classify, Assure, Compare); Verify Runtime node after runtime API ships
- Executive / GRC dashboard
- Third-party / vendor risk portal
- Hosted multi-tenant SaaS + billing

## Explicitly out of scope for now

- Embedding ISO/IEC clause text
- Claiming "enterprise-ready" or publishing a price list before design-partner feedback
- Collapsing non-EU regimes into prohibited/high/limited/minimal
- Inferring regimes from install locale or `geographic_scope`
