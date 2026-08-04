---
title: How we classify AI systems under the EU AI Act
published: false
tags: go, opensource, aiact, compliance, devops
canonical_url: https://github.com/new-world-coder/riskline
---


I built [RiskLine](https://github.com/new-world-coder/riskline) because most "AI governance" tooling wants a portal login before it will tell you anything useful. I needed something that runs in CI, prints a tier, and shows *which article fired* — without calling a model to invent a legal opinion.

RiskLine is a small Go CLI/API over a versioned JSON ruleset. You describe a system; it returns prohibited / high_risk / limited_risk / minimal_risk plus clause refs, a rationale, recommended controls, and a disclaimer. Advisory only. Not legal advice.

Below are three fixtures from the repo (`testdata/scenarios/`) and the actual CLI output on my machine against ruleset `eu-ai-act-2024-v0.1.0`.

## 1. Prohibited — social scoring

```yaml
name: Citizen Score
purpose: Score citizens based on social behaviour and personal characteristics for access to public services
data_types: [personal_data]
deployment_context: public_authority
autonomy_level: automated_decision
affected_population: general_public
social_scoring: true
```

```
Risk tier:        prohibited
Matched rules
prohibited-social-scoring  prohibited  Article 5(1)(c)
```

This one is boring on purpose. If the input asserts social scoring, Article 5(1)(c) wins. The interesting work is making sure higher-severity tiers always beat lower ones when multiple rules match — a law-enforcement biometric ID case should not get "downgraded" to high-risk because Annex III also fires.

## 2. High-risk — recruiting screener

```yaml
name: Hiring Assist
purpose: Screen job applicants and rank candidates for interview
data_types: [personal_data, employment_data]
deployment_context: saas_b2b
autonomy_level: decision_support
affected_population: job_applicants
```

```
Risk tier:        high_risk
Matched rules
high-risk-recruitment  high_risk  Annex III (4)(a)
```

Annex III (4)(a) is the recruitment/selection use-case list. Matching is partly keyword/population heuristics. That is a judgment call, and the response says so:

> Matching on purpose keywords or job_applicants population is a pragmatic proxy; borderline HR analytics that only aggregate headcount may be over-classified…

If you have a real HR analytics system that got over-classified, that is exactly what the [misclassification report](https://github.com/new-world-coder/riskline/issues/new?template=misclassification_report.md) template is for. I am using those issues as the validation loop instead of pretending I finished twenty customer interviews first.

## 3. Limited-risk — support chatbot

```yaml
name: Support Bot
purpose: Customer support chatbot for product questions
…
autonomy_level: content_generation
```

```
Risk tier:        limited_risk
Matched rules
limited-risk-chatbot-transparency  limited_risk  Article 50(1)
```

Art. 50 is wider than chatbots. The v0.1 rule requires both an interactive autonomy level *and* chatbot-ish purpose keywords, which probably under-detects some interactive systems. That uncertainty is intentional and logged on the rule. I would rather ship a narrow, tested heuristic than a confident wrong net.

## The boundary case I am least sure about

Creditworthiness vs fraud detection.

Annex III (5)(b) treats credit scoring of natural persons as high-risk, with a fraud-detection carve-out in the Act. RiskLine does **not** auto-infer that carve-out from free text. A system whose purpose says "fraud" but still carries `financial_credit_data` can still land high-risk and emit a judgment-call note.

I am not confident we have the right product behaviour here yet. If you have shipped lending + fraud models under the Act and think the heuristic is wrong, open a misclassification issue with the input you used. That report is more useful to me than another star.

## Try it

```bash
go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha
riskline-cli testdata/scenarios/high-risk-hiring.yaml
```

Schema/contract: [`api/openapi.yaml`](https://github.com/new-world-coder/riskline/blob/main/api/openapi.yaml). Repo: [new-world-coder/riskline](https://github.com/new-world-coder/riskline).
