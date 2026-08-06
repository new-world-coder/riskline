# RiskLine demo video — shot list & narration

Target length: **~75 seconds**. Source of truth for on-screen output: the real CLI in this repo (`bin/riskline-cli` or `go run ./cmd/riskline-cli`). A matching asciinema cast lives at [`riskline-demo.cast`](riskline-demo.cast).

Convert the cast to gif/mp4 with [agg](https://github.com/asciinema/agg) when you need a Product Hunt / social upload:

```bash
agg --font-size 18 --theme monokai demo/riskline-demo.cast demo/riskline-demo.gif
# or render mp4 via your preferred pipeline from the cast / gif
```

Upload the cast to [asciinema.org](https://asciinema.org) for an embeddable player.

---

## Voice

Direct, engineer-to-engineer, anti-hype. Mirror [`content/social/LAUNCH_COPY.md`](../content/social/LAUNCH_COPY.md). Do not say “enterprise-ready,” “AI-powered compliance,” or claim full Act coverage.

---

## Shot list

| Time | Visual | Narration |
|------|--------|-----------|
| 0:00–0:08 | Dark terminal. Brand line or `$ # RiskLine` then install command typing. | “RiskLine — offline EU AI Act risk classification. CLI and API. Deterministic ruleset, not an LLM.” |
| 0:08–0:22 | Run Citizen Score fixture. Highlight `Risk tier: prohibited` and `Article 5(1)(c)`. | “Describe an AI system. Here’s a public-authority social scoring tool — classified prohibited under Article 5.” |
| 0:22–0:40 | Run Hiring Assist. Highlight `high_risk` and `Annex III (4)(a)`. | “Recruitment ranking? High-risk — Annex III — with recommended controls and judgment calls you can review.” |
| 0:40–0:55 | Run Support Bot. Highlight `limited_risk` / `Article 50(1)`. | “A customer support chatbot lands as limited-risk with a transparency obligation.” |
| 0:55–1:08 | `riskline-cli --json …` peek, or brief `curl` to `/v1/classify`. | “Same contract as JSON for CI, or as an HTTP API. Offline by default — the CLI doesn’t phone home.” |
| 1:08–1:15 | Freeze on GitHub URL + `go install …@v0.1.0-alpha`. | “Open source, Apache 2.0, pre-1.0. Misclassification reports welcome. github.com/new-world-coder/riskline” |

---

## Commands to film (exact)

```bash
go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha
# or from a local build:
# go build -o bin/riskline-cli ./cmd/riskline-cli

./bin/riskline-cli testdata/scenarios/prohibited-social-scoring.yaml
./bin/riskline-cli testdata/scenarios/high-risk-hiring.yaml
./bin/riskline-cli testdata/scenarios/limited-risk-chatbot.yaml
./bin/riskline-cli --json testdata/scenarios/high-risk-hiring.yaml | head -n 40
```

Optional API beat (second terminal or after starting the API):

```bash
go run ./cmd/riskline-api -addr :8080
curl -s localhost:8080/v1/classify \
  -H 'content-type: application/json' \
  -d @examples/curl/hiring-assist.json
```

---

## Recording tips

- Font: IBM Plex Mono or JetBrains Mono, **18–22pt**, high contrast dark theme (charcoal + cyan, not purple neon).
- Window: ~1280×720 or 1920×1080; crop chrome; no desktop clutter.
- Pause ~1s after each `Risk tier:` line so viewers can read the clause ref.
- Do not skip the disclaimer if you show full human output for more than a beat — or cut after matched rules and say “advisory, not legal advice” in VO.
- Prefer the checked-in cast for consistency with the landing page demo strings.

## Re-record the cast

```bash
go build -o bin/riskline-cli ./cmd/riskline-cli
asciinema rec demo/riskline-demo.cast --overwrite \
  --command 'bash demo/record-demo.sh'
```

`demo/record-demo.sh` drives the same three scenarios plus a JSON peek.
