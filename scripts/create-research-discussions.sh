#!/usr/bin/env bash
# One-time: create Classify → Assure research discussions (requires gh auth).
set -euo pipefail

REPO_ID="R_kgDOTto-0w"
CAT_QA="DIC_kwDOTto-084DCp4L"
CAT_ANN="DIC_kwDOTto-084DCp4J"
WP="https://github.com/new-world-coder/riskline/blob/main/docs/whitepaper-the-conformity-gap.md"

create_one() {
  local category="$1"
  local title="$2"
  local body="$3"
  jq -n \
    --arg query 'mutation($input: CreateDiscussionInput!) { createDiscussion(input: $input) { discussion { number url } } }' \
    --arg repositoryId "$REPO_ID" \
    --arg categoryId "$category" \
    --arg title "$title" \
    --arg body "$body" \
    '{query: $query, variables: {input: {repositoryId: $repositoryId, categoryId: $categoryId, title: $title, body: $body}}}' \
    | gh api graphql --input -
}

create_one "$CAT_QA" "Research Q1: Material change — re-classify or re-assure only?" \
"What events in your organisation should trigger a **mandatory re-classification** (new system description → new tier) versus a **re-assurance only** pass (same tier, verify controls still hold)?

Examples we hear about:
- Model version bump without purpose change
- New tool/API attached to an agent
- IaC change affecting logging or human-review gates
- Geographic scope expansion (EU users added)

Context: [The Conformity Gap white paper](${WP}) frames **Classify → Assure** — design-time tiering plus continuous verification.

**We are collecting practitioner signal** to prioritise probe design. No implementation details here — just your operational rules of thumb.

Advisory discussion only — not legal advice."

create_one "$CAT_QA" "Research Q2: Agent runtime signals for human oversight (without prompt capture)" \
"For **agentic / autonomous AI systems**, which **runtime signals** are sufficient to evidence human oversight under your programme — **without** capturing full prompt content?

Candidates (opinion welcome):
- Human override events in audit log
- Kill-switch endpoint health
- Allowlisted action set violations blocked
- Escalation queue depth / SLA

What would **you** accept as audit evidence? What is impractical or privacy-toxic?

Related: EU AI Act Art. 14 (human oversight), NIST AI RMF Manage functions.

White paper: ${WP}

Advisory only — not legal advice."

create_one "$CAT_QA" "Research Q3: Multi-regime assurance plans (EU AI Act + NIST mapping)" \
"RiskLine can classify under **EU AI Act** (hard-law tiers) and **NIST AI RMF** (mapping-only, \`mapping_only: true\`) in one request.

**Question:** How should a single **assurance plan** treat both outputs?

- One probe list merging all \`technical_controls[]\`?
- Separate evidence bundles per regime?
- EU probes mandatory; NIST probes optional for US-facing buyers?

If you run dual-regime programmes today, how do you avoid contradicting yourself in audit?

White paper: ${WP}

Advisory only — not legal advice."

create_one "$CAT_QA" "Research Q4: GRC evidence shape for Vanta, Drata, and auditors" \
"If a deterministic **assure** layer produced pass/fail probe results (ruleset version, probe id, timestamp, signed JSON), **what shape** would your GRC stack or auditor accept as first-class evidence?

- Vanta custom resource / document upload?
- Drata Custom Connection JSON schema?
- Plain Markdown + hash?
- Something else entirely?

We are **not** asking for vendor endorsement — asking what **you** have seen work in practice.

White paper: ${WP}

Advisory only — not legal advice."

create_one "$CAT_QA" "Research Q5: Art. 72 probe cadence vs risk tier" \
"EU AI Act **Article 72** post-market monitoring must be proportionate to the technology and risks.

**How should automated re-probe cadence scale with tier?**

| Tier | Your cadence instinct |
|------|----------------------|
| High-risk | ? |
| Limited-risk | ? |
| Minimal-risk | ? |

Triggers we consider: deploy events, ruleset version bumps, scheduled cron, material config drift.

What is **too much** (noise) vs **too little** (Art. 72 gap)?

White paper: ${WP}

Advisory only — not legal advice."

create_one "$CAT_ANN" "White paper: The Conformity Gap + Classify → Assure research threads" \
"We published a public white paper on the gap between **design-time classification** and **continuous verification**:

📄 ${WP}

**Classify** (open source today): deterministic EU AI Act + NIST AI RMF mapping — offline CLI, citations, no LLM on the verdict path.

**Assure** (alpha, design partners): planned probes bound to the same ruleset — static, runtime, scheduled. Implementation private while we validate with partners.

**Join the research agenda** — five open Q&A threads in this repo's Discussions (material change, agent signals, multi-regime, GRC evidence shape, Art. 72 cadence).

Design-partner pilots: https://riskline-cloud-web.vercel.app/pilots

Advisory only — not legal advice."

echo "Done — check https://github.com/new-world-coder/riskline/discussions"
