# TrustModel outreach — Aarushi Jaitly (Phase 0)

**Status:** Draft for your review before posting. Do not paste product name in the public comment.

---

## Public LinkedIn reply (post this on the Conformity Gap article)

Thank you, Aarushi — I think you've put your finger on the part of AI governance that becomes harder after deployment: **proving that the control remains effective once the system is in the real world.**

There's a distinction I've been thinking about:

**Monitoring tells us what happened.**  
**Continuous verification should tell us whether what happened was still within the authorised control boundary.**

Those aren't necessarily the same thing.

A model update, new tool/API permission, changed data source, modified system prompt, or a new agent-to-agent interaction can alter that boundary without anyone consciously deciding that the AI's risk classification has changed.

That makes me wonder whether the next evolution of AI assurance is less about producing another point-in-time assessment and more about creating **machine-verifiable evidence of compliance at runtime** — evidence an auditor can re-run, not reconstruct from screenshots.

I'd be genuinely interested in your perspective, especially given your work in AI safety and agent evaluation: **where do you see the bigger unsolved problem today — detecting that an AI system has drifted, or proving that its controls remained effective when it actually made a consequential decision?**

---

## DM (send only if she replies or engages)

Your comment got me thinking about something I've been building on the evidence side rather than the assessment side.

The question I'm trying to answer is: **can an application generate machine-verifiable proof that a particular AI action was still permissible under its policy and regulatory boundary at the moment it occurred?**

I've been working on an open-source engine called **RiskLine** around deterministic EU AI Act classification plus a continuous verification layer (material-change detection, conformity state, hashed evidence — local integrity chain, not blockchain).

Given your work at TrustModel, I'd genuinely value your technical perspective on where this overlaps with your architecture and where it might fill a gap. Would you be open to a 20-minute walkthrough where you can challenge the design?

---

## What not to say (yet)

- "I want to become your Europe head"
- "TrustModel doesn't do EU"
- "RiskLine replaces TrustModel"
- Any mention of acquisition in the first conversation

---

## Meeting demo script (when she says yes)

1. **Classify** — `riskline-cli --json examples/curl/hiring-assist.yaml` → high-risk + `technical_controls`
2. **Diff** — `riskline-cli diff --json examples/curl/hiring-assist.yaml examples/curl/hiring-assist-automated.yaml` → reclassify on autonomy change
3. **Assure** — `riskline-cli assure --probes examples/curl/hiring-assist-probes.json --json <classify.json>` → GREEN; fail a probe → RED
4. **Evidence** — show `content_hash` + `previous_hash` chain; explain SHA-256 local integrity, **not** a public blockchain ledger
5. **Complement slide** — TrustModel = TrustScore / eval / monitor; RiskLine = article-cited tier + obligation-bound verification from the same ruleset

---

## Follow-up if product team is interested

> TrustModel provides trust assessment and assurance scoring. RiskLine provides an independent, developer-native **conformity verification** layer that continuously generates evidence from the running system's declared boundary — especially for EU AI Act article-level obligations. Happy to explore an evidence API feed into your Control Plane if useful.
