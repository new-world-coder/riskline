# Introducing the Riskline Console

We're building trust infrastructure for AI systems — starting with deterministic EU AI Act classification and expanding toward continuous verification. Today we're sharing a preview of the **Riskline Console**: a role-based operating view for the people who need to *see* AI compliance posture, not just run a CLI command.

Advisory only. Not legal advice. Not a substitute for counsel or a formal conformity assessment.

## The problem

AI governance isn't a developer-only problem. When a board member asks "what's our AI risk exposure?", the answer shouldn't require reading JSON from a terminal. Compliance officers need obligation tracking. Security teams need control coverage. Auditors need evidence trails. Developers need CI integration. Each persona asks a different question — but they're looking at the same underlying inventory of AI systems.

Most tools pick one audience and ignore the rest. We think the classify step is shared infrastructure; the views should be persona-specific.

## What the Console is

The Riskline Console is a **trust & compliance operating console** with dashboards tailored to six roles:

- **Executive** — board-ready posture, risk exposure, systems requiring attention
- **Compliance** — EU AI Act obligation coverage, open findings, evidence gaps
- **Operations** — system inventory, verification cadence, incident queue
- **Security** — control coverage, probe failures, threat signals
- **Developer** — API usage, CI integrations, classify workflow
- **Auditor** — evidence trails, findings, attestation exports

Plus a **regulation update feed** showing incoming regulatory changes and a **try-it flow** for interactive demos.

## What the Console is not

The Console demo you can explore today is a **mango-skin preview**. It shows what the operating experience looks like — with mock data — while the proprietary engine stays backend-only:

- **Continuous verification** (probe executors, drift detection, post-market monitoring) — private, alpha
- **Policy engine** (auto-remediation, enforcement rules) — private
- **Regulation update pipeline** (impact scoring, auto-reclassification) — private

The public demo shows *outcomes* (pass/fail, compliance %, findings). It never exposes *how* verification runs. That's deliberate — and it's how we ship a credible demo without giving away the IP.

What's real today:

- **Classify** — deterministic EU AI Act tier with clause citations: https://riskline-cloud-web.vercel.app/classify
- **Hosted API** — same OpenAPI contract as the open-source CLI: https://riskline-cloud-web.vercel.app/developers
- **Pilot waitlist** — human-led classification engagements: https://riskline-cloud-web.vercel.app/pilots

## Why now

Buyers need to see the vision — a CLI install command doesn't close a design-partner pilot. Persona coverage wins deals: the compliance officer and the CISO are different buyers. A demo-grade console de-risks the pitch: we show the UX without shipping proprietary logic.

The Console is a sales artifact that sits alongside our existing classify report and evidence export — not a replacement for them.

## Architecture in one paragraph

Developers classify via CLI or API (open-source engine, versioned rules, no LLM guessing). Results flow into persona dashboards. Continuous verification — the Assure layer — runs backend-only and surfaces outcomes to Operations and Security views. Regulation updates arrive through a gateway pipeline and appear in the Compliance and Executive views.

## Try it

- **Console demo (mock metrics):** https://riskline-cloud-web.vercel.app/console
- **Live classify:** https://riskline-cloud-web.vercel.app/classify
- **Pilot waitlist:** https://riskline-cloud-web.vercel.app/pilots

## What's next

- Connect console widgets to real APIs as Assure alpha opens
- SSO-gated console for paying customers
- Regulation feed from live pipeline (when impact scoring ships)

---

*Advisory tool only. Not legal advice. Not a conformity assessment.*

**Published:** August 2026 · **Engine:** [github.com/new-world-coder/riskline](https://github.com/new-world-coder/riskline)
