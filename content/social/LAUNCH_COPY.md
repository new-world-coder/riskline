# RiskLine launch copy (post yourself)

Do **not** auto-post these unattended. Paste after a quick read. Repo: https://github.com/new-world-coder/riskline

## Show HN (title + text)

**Title:** Show HN: RiskLine – offline EU AI Act risk classification CLI/API

**Text:**
I built a small Go tool that classifies an AI system description under the EU AI Act (prohibited / high-risk / limited-risk / minimal-risk) using a versioned JSON ruleset — not an LLM.

Why: most governance tools want a portal; I needed something that runs in CI, stays offline by default, and shows which article/annex fired.

```
go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha
```

Happy to take misclassification reports — that is how I am validating the ruleset.

https://github.com/new-world-coder/riskline

## Reddit — r/golang

**Title:** RiskLine – EU AI Act risk classification as a Go CLI (offline, versioned ruleset)

**Body:**
Built a deterministic classifier for EU AI Act risk tiers. Engine is pure Go (`pkg/engine`), CLI and API are thin wrappers. Rules live in embedded JSON with `ruleset_version` / `last_updated` on every mapping.

Not using OPA yet — rule surface is small. Distribution model borrows from OPA/Trivy (single binary, CI-friendly).

Looking for feedback from people who have put Go tools into other teams' pipelines:

https://github.com/new-world-coder/riskline

Post 3 (why Go): https://github.com/new-world-coder/riskline/blob/main/docs/blog/03-why-go-for-a-compliance-engine.md

## Reddit — r/opensource

**Title:** RiskLine – open-source EU AI Act risk classification (CLI + OpenAPI)

**Body:**
Apache-2.0. Classifies AI systems into EU AI Act risk tiers with clause references and a mandatory disclaimer field.

Contribution paths that actually help:
- misclassification reports
- ruleset update issues
- a few good-first-issues (goldens, examples, docs)

https://github.com/new-world-coder/riskline

## LinkedIn

Shipped an open-source EU AI Act risk classifier as a CLI/API — deterministic ruleset, not a model guessing compliance.

If you are inventorying AI systems for the Act, you can run it offline:

go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha

I am collecting misclassification reports as the validation loop. Methodology write-up: [link to post 1 once live on your blog/dev.to]

https://github.com/new-world-coder/riskline

## X / Twitter

Open-sourced RiskLine: offline EU AI Act risk classification (CLI + API).

Versioned JSON ruleset → prohibited / high-risk / limited-risk / minimal-risk, with clause refs + disclaimer.

go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha

https://github.com/new-world-coder/riskline

## Suggested posting order

1. Push content to GitHub (done via this PR/commit)
2. Publish Post 1 to dev.to as draft → review → publish
3. Same day: LinkedIn + X with link to Post 1 or repo
4. Next day: Show HN (once you can answer comments for a few hours)
5. Same week: r/golang, then r/opensource (don't crosspost dump)
6. Post 2 and Post 3 spaced 3–7 days apart
