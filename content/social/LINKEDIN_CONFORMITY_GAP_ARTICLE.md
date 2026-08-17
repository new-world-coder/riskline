# LinkedIn article — paste into “Write article”

**Title:** The Conformity Gap: why AI risk classification fails without continuous verification

**Cover line (subtitle):** EU AI Act Art. 72 and NIST Measure/Manage expect more than a one-time tier label. Here is the Classify → Assure model — and five questions I want practitioners to answer.

---

Most AI governance programmes look the same on paper.

Inventory your AI systems. Classify risk tier. File evidence in a GRC tool. Pass audit.

That workflow works beautifully — until day two.

After deployment, models get updated. Agents gain new tools. Infrastructure drifts. Logging gets turned down to save cost. The system you certified is not the system that is running.

I call the space between those two states **the Conformity Gap**.

It is not a marketing term. EU AI Act Article 43 gives you a pre-market conformity snapshot. Article 72 requires post-market monitoring for the operational lifetime of high-risk systems. Article 9 mandates a continuous, iterative risk management process. NIST AI RMF Measure and Manage functions say the same thing in different language: design-time classification is necessary but not sufficient.

---

## What classification actually does

Design-time classification answers one question: *Given what this system is intended to do, which obligations apply?*

A good classifier is deterministic — same input and ruleset version, same output. It cites articles and annexes. It runs offline in CI. It does not ask an LLM to guess your tier.

That is the layer we ship today in RiskLine: open-source EU AI Act tiers plus NIST AI RMF mapping with technical control hooks.

But a tier label is a building permit, not a building inspection.

---

## What assurance must do

Assurance answers a different question: *Are the controls tied to that tier actually present and operating?*

That requires probes — deterministic checks bound to the obligations classification produced:

- **Static probes** — CI configs, IaC, repository files
- **Runtime validation** — agent middleware, API event checks (metadata only, customer-deployed)
- **Scheduled re-probes** — cadence aligned with post-market monitoring

The ruleset that produces the tier must be the same ruleset that defines what gets verified. Otherwise classification and verification decouple, and the gap reopens.

We call this **Classify → Assure**. Layer one is public and open source. Layer two is in alpha with design partners.

---

## Five questions I am asking the community

I published a white paper and opened GitHub Discussions because I do not have all the answers:

1. **Material change** — What events should force re-classification vs re-assurance only?
2. **Agentic systems** — Which runtime signals prove human oversight without capturing prompt content?
3. **Multi-regime plans** — How should EU hard-law tiers coexist with NIST mapping-only output in one assurance plan?
4. **GRC evidence** — What JSON shape do Vanta, Drata, and auditors need to treat assure output as first-class evidence?
5. **Art. 72 cadence** — How should probe frequency scale with risk tier?

If you operate in this space, I would value your perspective in the discussions linked below.

---

## What we are not claiming

- Not legal advice or conformity assessment
- Not a universal compliance bot for every industry
- Not central traffic sniffing from a vendor SaaS — sensors run in your environment
- Not an LLM on the verdict path

---

## Links

- White paper: https://github.com/new-world-coder/riskline/blob/main/docs/whitepaper-the-conformity-gap.md
- Discussions: https://github.com/new-world-coder/riskline/discussions
- Classify demo: https://riskline-cloud-web.vercel.app/classify
- Design-partner pilots: https://riskline-cloud-web.vercel.app/pilots

*Advisory tool only — not legal advice.*

#EUAIAct #AIGovernance #RegTech #PostMarketMonitoring #NIST
