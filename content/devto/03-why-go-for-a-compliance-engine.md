---
title: Why I built the compliance engine in Go
published: false
tags: go, opensource, aiact, compliance, devops
canonical_url: https://github.com/new-world-coder/riskline
---


I did not pick Go because it is fashionable. I picked it because the product has to land in someone else's CI without a negotiation about runtimes.

RiskLine's job is: read a YAML/JSON system description, evaluate a versioned ruleset, print a tier. That has to work as a pre-commit hook, a GitHub Action step, and a one-liner `go install`. Node and Python are fine languages. They are worse at "here is a single static binary, trust it in your pipeline."

## What "CI-friendly" actually means

- No interpreter version matrix for adopters
- Cross-compile for linux/mac/windows from one tree
- Cold start fast enough that nobody disables the hook
- Offline by default — a compliance CLI that phones home is a trust problem

That distribution model is the same reason people reach for Open Policy Agent, Trivy, and Terraform-shaped tools. RiskLine is not built on OPA/Rego in v0.1 — the rule surface is small enough that embedded JSON + Go matching was enough — but the packaging lesson is the same. If the ruleset grows into something Rego expresses more clearly, that is a later decision, not a day-one dependency.

## Embeddability

`pkg/engine` does not import net/http or CLI packages. The API server and the CLI are thin wrappers. That is what makes "SDK" real instead of "HTTP client with feelings." If you want classification inside another Go service, you call `engine.Classify`, not a localhost port.

## Tradeoffs I am accepting

Go JSON tags and OpenAPI can drift if you are careless — hence the contract being a checked-in YAML file and a (still too weak) CI smoke test. Help wanted on making that a real schema validation step: [issue #4](https://github.com/new-world-coder/riskline/issues/4).

Also: writing rules in JSON conditions is less expressive than Rego. For prohibited flags and a handful of Annex III heuristics, it has been fine. I will revisit when the ruleset stops being fine.

## Try the binary

```bash
go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha
riskline-cli --json path/to/system.yaml
```

Source: [new-world-coder/riskline](https://github.com/new-world-coder/riskline).
