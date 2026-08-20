# OSS competitive landscape (Aug 2026)

Research for ultra-strong positioning vs TrustModel and GitHub EU AI Act tooling.
See also [trustmodel-partner-strategy canvas](/Users/stardust/.cursor/projects/Users-stardust-git-riskline/canvases/trustmodel-partner-strategy.canvas.tsx).

## The wedge nobody ships end-to-end

**Classify → Assure from one versioned ruleset.** Competitors stop at classification (EuConform, Lucairn) or CI/repo scans (Regula, systima/comply, Legalithm) or runtime agent governance (Microsoft AGT, aulite). None bind article-level tier obligations to continuous probe specs and material-change policy in one deterministic engine.

## Tier A — Classification (Layer 1 competitors)

| Repo | Stars | Threat |
|------|-------|--------|
| [Hiepler/EuConform](https://github.com/Hiepler/EuConform) | ~123 | Deepest OSS classification mindshare; could add CI |
| [Declade/lucairn-ai-act-classifier](https://github.com/Declade/lucairn-ai-act-classifier) | ~1 | Same offline/deterministic niche; stronger Annex III today |
| [AbdelStark/eu-ai-act-toolkit](https://github.com/AbdelStark/eu-ai-act-toolkit) | ~6 | Classification + checklists + dashboard |

**Lead with:** multi-regime (EU + NIST + MAS path), embeddable Go engine, OpenAPI, Classify→Assure (now shipping v0.4).

## Tier B — CI continuous (Layer 2 competitors)

| Repo | Stars | Threat |
|------|-------|--------|
| [microsoft/Agent-Governance-Toolkit](https://github.com/microsoft/Agent-Governance-Toolkit) | ~6052 | Default agent governance kernel; EU checklist not tier classifier |
| [legalithm-org/legalithm](https://github.com/legalithm-org/legalithm) | ~0 | Drift on committed compliance records; MCP + Action |
| [systima-ai/comply](https://github.com/systima-ai/comply) | ~0 | PR baseline diff, SARIF, obligation scoring |
| [kuzivaai/getregula](https://github.com/kuzivaai/getregula) | ~4 | Code pattern scan → tier |

**Lead with:** “They scan code or repo patterns; we classify **system intent** from inventory YAML, then same ruleset drives assure probes.” Compose, don’t fight AGT.

## Tier C — Commercial overlap

| Actor | Threat |
|-------|--------|
| [TrustModel](https://trustmodel.ai) | TrustScore + eval + monitor + AGP; EU mapping on website |
| [VerifyWise](https://github.com/verifywise-ai/verifywise) | Self-hosted GRC portal (~338★) |

**Lead with:** article citations + hashed evidence + offline CLI; complement TrustModel eval with conformity verification.

## Positioning lines (use in public copy)

1. **TrustModel scores outputs. RiskLine cites articles — and proves controls tied to those articles still exist.**
2. **No TrustScore. No LLM judge. Same YAML in, same tier out — pinned by `ruleset_version`.**
3. **One ruleset, two layers:** classify obligations → assure probes → material change → conformity state.
4. **Evidence is SHA-256 hash chain locally — tamper-evident audit log, not a public blockchain.**

## Integration story (10 miles ahead)

Publish compose guides: RiskLine tier → AGT/comply/Legalithm as probe implementations. Own the **seam** between legal tier and operational proof.
