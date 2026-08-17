# The Conformity Gap: Why AI Risk Classification Without Continuous Verification Fails Post-Market Obligations

**Author:** Sachin Chitre, RiskLine  
**Version:** 1.0 (public draft)  
**Date:** August 2026  
**Status:** Thought leadership — not legal advice  

---

## Abstract

Most AI governance tooling today stops at design-time classification: a system is described, a risk tier is assigned, and evidence is filed for audit. That satisfies pre-market conformity narratives for many organisations — but it leaves a structural gap between **what was certified** and **what is running**.

This paper introduces the **Classify → Assure** model: a two-layer, deterministic compliance architecture in which classification produces versioned control obligations, and a separate assurance layer continuously verifies that those obligations remain present across infrastructure, code, and runtime behaviour. We argue this gap is not a product marketing construct; it is implied by EU AI Act Articles 9, 43, and 72, and by NIST AI RMF Measure and Manage functions — yet remains largely unaddressed by enterprise GRC suites focused on SOC 2 evidence automation.

We describe the problem, the regulatory basis, a reference architecture at the conceptual level, and a research agenda for practitioners. Implementation details, probe specifications, and integration patterns are intentionally out of scope for this document.

**Keywords:** EU AI Act, post-market monitoring, AI governance, deterministic compliance, continuous assurance, NIST AI RMF

---

## 1. Introduction

The EU AI Act creates a lifecycle obligation for providers of high-risk AI systems. Article 43 requires conformity assessment before market placement. Article 72 requires post-market monitoring for the operational lifetime of the system. Article 9 mandates a continuous, iterative risk management process.

In practice, vendor tooling and internal programmes converge on a familiar pattern:

1. Inventory AI systems  
2. Classify risk tier (prohibited / high / limited / minimal)  
3. Document controls in a GRC platform  
4. Pass audit  

Steps 1–3 are increasingly automated. Step 4 produces a point-in-time certificate. What happens on day 2 — after deployment, after a model update, after an agent gains new tools, after infrastructure drifts — is where most programmes become manual again.

We call the space between certified design and verified operation **the Conformity Gap**.

---

## 2. The Conformity Gap defined

**Definition.** The Conformity Gap is the measurable divergence between (a) the control obligations associated with a system's classified risk tier at time T₀, and (b) the observable state of controls in code, infrastructure, and runtime at time T₁.

The gap widens when:

- System descriptions change without re-classification  
- Infrastructure (IaC, VPC, Kubernetes) drifts from documented posture  
- Agentic systems gain autonomy or tool access not reflected in technical documentation  
- Logging, human oversight, or transparency mechanisms degrade after launch  
- Rulesets update but deployed systems are not re-evaluated  

The gap is not hypothetical. Post-market monitoring under Article 72 explicitly requires providers to collect and analyse performance data across the system's lifetime — not merely retain pre-market paperwork.

---

## 3. Regulatory basis (high level)

| Obligation | Article / Function | Design-time? | Continuous? |
|------------|-------------------|--------------|-------------|
| Risk tier determination | EU AI Act Art. 6, Annex III | Yes | On material change |
| Conformity assessment | Art. 43 | Yes (pre-market) | On substantial modification |
| Risk management system | Art. 9 | Documented at design | Iterative, lifecycle-long |
| Record-keeping | Art. 12 | Designed | Operated |
| Transparency | Art. 13 | Designed | Deployed UI/API must match |
| Human oversight | Art. 14 | Designed | Runtime must enforce |
| Post-market monitoring | Art. 72 | Planned pre-market | **Operational lifetime** |
| NIST AI RMF Measure | Measure 2.x | Baseline metrics | Drift, re-evaluation |
| NIST AI RMF Manage | Manage 2.x | Incident procedures | Active response |

**Implication.** Classification alone addresses the left column. Continuous assurance addresses the right. Conflating the two produces compliance theatre: a correct tier with no mechanism to detect when reality diverges.

*This paper is an engineering and governance perspective, not legal advice. Organisations should consult qualified counsel for conformity assessment and regulatory interpretation.*

---

## 4. Why current tooling leaves the gap open

Enterprise GRC platforms (SOC 2, ISO 27001 automation) excel at:

- Integrating cloud accounts, identity providers, and SaaS tools  
- Collecting configuration evidence on schedules  
- Mapping controls to frameworks for audit  

They are weaker at:

- **Regime-specific AI tier logic** with article-level citations  
- **Linking a classified obligation to a concrete technical hook** in the customer's stack  
- **Verifying AI-specific controls** (human override on agent actions, transparency banners on AI outputs, classify-on-change in CI)  
- **Deterministic, reproducible verdicts** auditable without model stochasticity  

LLM-based "compliance copilots" address narrative generation but introduce non-reproducibility — problematic when evidence must be re-run by auditors or regulators.

The gap persists because classification and verification are treated as separate products, separate vendors, and separate audit artefacts — with no shared, versioned ruleset binding them.

---

## 5. The Classify → Assure reference model

We propose a two-layer architecture bound by a **single versioned ruleset**:

```
┌─────────────────────────────────────────────────────────────┐
│  LAYER 1 — CLASSIFY (design-time)                           │
│  Input:  system description (metadata, not source code)     │
│  Output: risk tier + control obligations + evidence plan    │
│  Engine: deterministic rules (no LLM on verdict path)       │
└──────────────────────────┬──────────────────────────────────┘
                           │ shared ruleset version
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 2 — ASSURE (continuous)                              │
│  Input:  infrastructure state, CI events, runtime signals   │
│  Output: pass/fail per obligation + signed evidence bundle  │
│  Engine: deterministic probes (no LLM on verdict path)      │
└─────────────────────────────────────────────────────────────┘
```

### 5.1 Layer 1 — Classify

Classification answers: *Given what this system is intended to do, which regulatory obligations apply?*

Properties we consider non-negotiable:

- **Deterministic** — same input + ruleset version → same output  
- **Versioned** — ruleset changes are explicit and diffable  
- **Cited** — outputs reference articles, annexes, or framework subcategories  
- **Offline-capable** — runnable in CI without sending source code to a vendor  

### 5.2 Layer 2 — Assure

Assurance answers: *Given what was classified, are the associated controls present and operating?*

Assurance operates through **probes** — deterministic checks bound to obligations produced by classification. Probes may be:

- **Static** — evaluate IaC, repository configuration, CI pipeline definitions  
- **Event-driven** — validate runtime signals submitted by customer-deployed agents  
- **Scheduled** — re-run on cadence to support post-market monitoring narratives  

### 5.3 The binding principle

The ruleset that produces a tier must be the same ruleset that defines what assurance checks. Without this binding, classification and verification decouple — the gap reopens.

*This document does not specify probe formats, hook registries, evidence signing schemes, or API contracts. Those are implementation concerns under active development.*

---

## 6. Deployment and trust boundaries

Continuous assurance requires explicit trust boundaries:

1. **Customer-deployed sensors.** Probes that touch production systems, agent middleware, or network paths run inside the customer's environment — not on a vendor's multi-tenant sniffing plane.  
2. **Metadata-first APIs.** Hosted services receive structured metadata (system descriptions, event summaries) — not arbitrary source code or raw traffic payloads — unless the customer explicitly opts in.  
3. **Reproducible evidence.** Every assurance verdict references ruleset version, probe identity, and timestamp — enabling third-party re-run.  
4. **Advisory positioning.** Tools produce evidence for human and counsel review; they do not replace conformity assessment or legal determination.  

---

## 7. Research agenda

We invite practitioners to engage on the following open questions:

1. **Material change detection.** What events should trigger mandatory re-classification vs re-assurance only?  
2. **Agentic systems.** Which runtime signals are sufficient to verify human oversight without capturing prompt content?  
3. **Multi-regime orchestration.** How should EU AI Act hard-law tiers coexist with NIST RMF mapping-only outputs in a single assurance plan?  
4. **GRC integration.** What evidence shape do Vanta, Drata, and audit firms need to accept assurance output as first-class audit artefacts?  
5. **Post-market proportionality.** How should probe cadence scale with risk tier under Article 72?  

Discussion: https://github.com/new-world-coder/riskline/discussions/19  
Design partners: https://riskline-cloud-web.vercel.app/pilots  

---

## 8. Conclusion

The Conformity Gap is the operational space between certification and reality. Regulatory text already anticipates continuous verification — particularly EU AI Act Article 72 and NIST AI RMF Measure/Manage functions — but tooling has not caught up.

The Classify → Assure model addresses this gap with a deliberately narrow claim: **one versioned ruleset, two deterministic layers, customer-deployed assurance, reproducible evidence.** It is not a universal compliance bot, not an LLM guessing risk tiers, and not a replacement for counsel or conformity assessment.

Authority in this space will accrue to teams that publish the problem framing clearly, ship reproducible artefacts, and close the gap with evidence — not slide decks.

---

## About RiskLine

RiskLine is an open-core AI governance engine. The classification layer (EU AI Act, NIST AI RMF mapping) is open source under Apache 2.0. The assurance layer and hosted services are under active development with design partners.

- Classify demo: https://riskline-cloud-web.vercel.app/classify  
- Open source: https://github.com/new-world-coder/riskline  
- Design-partner pilots: https://riskline-cloud-web.vercel.app/pilots  

**Disclaimer:** RiskLine is an advisory tool based on versioned rulesets. It is not legal advice and is not a substitute for qualified counsel or a formal conformity assessment.

---

## Citation

If referencing this paper:

> Chitre, S. (2026). *The Conformity Gap: Why AI Risk Classification Without Continuous Verification Fails Post-Market Obligations* (v1.0). RiskLine. https://github.com/new-world-coder/riskline/blob/main/docs/whitepaper-the-conformity-gap.md

---

*© 2026 RiskLine. This white paper may be shared freely with attribution. Implementation details, probe specifications, and partner integration patterns are not licensed for commercial reproduction.*
