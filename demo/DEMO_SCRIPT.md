# RiskLine demo video script (~75 seconds)

Use this to record a terminal screencast or convert `riskline-demo.cast` with [agg](https://github.com/asciinema/agg) to gif/mp4 for Product Hunt / social.

**Prep**
- Terminal: dark theme, font size 16–18pt, 80×24 cols minimum
- Build: `go build -o bin/riskline-cli ./cmd/riskline-cli`
- Crop: terminal only, no desktop clutter

---

## Shot list

| Time | Visual | Narration |
|------|--------|-----------|
| 0:00–0:08 | Title card or terminal prompt | "RiskLine classifies AI systems under the EU AI Act — offline, with a versioned ruleset, not an LLM." |
| 0:08–0:15 | `go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha` | "Install is one line. Same binary runs locally and in CI." |
| 0:15–0:28 | `./bin/riskline-cli testdata/scenarios/prohibited-social-scoring.yaml` | "Citizen Score — social scoring for public services. Prohibited under Article 5." |
| 0:28–0:42 | `./bin/riskline-cli testdata/scenarios/high-risk-hiring.yaml` | "Hiring Assist screens applicants. High-risk — Annex III recruitment. You get clause refs and recommended controls." |
| 0:42–0:55 | `./bin/riskline-cli testdata/scenarios/limited-risk-chatbot.yaml` | "Support Bot — limited risk. Article 50 transparency obligations." |
| 0:55–1:05 | `./bin/riskline-cli --json testdata/scenarios/high-risk-hiring.yaml` \| head | "Machine output for pipelines — same engine, JSON contract." |
| 1:05–1:15 | Browser on github.com/new-world-coder/riskline or `curl` tease | "Open source, Apache 2.0. Misclassification reports are how we validate the ruleset." |

---

## Recording commands

```bash
cd riskline
go build -o bin/riskline-cli ./cmd/riskline-cli

# Option A: record fresh cast (requires asciinema)
asciinema rec demo/riskline-demo.cast

# Option B: use the committed cast in demo/
asciinema play demo/riskline-demo.cast

# Convert to gif/mp4 (requires agg)
agg demo/riskline-demo.cast riskline-demo.gif
```

---

## Voice notes

- Direct, anti-hype — match `content/social/LAUNCH_COPY.md`
- Say "advisory, not legal advice" once near the end
- Don't claim full Annex III coverage — "prioritised categories in v0.1"

---

## Product Hunt / LinkedIn

- Thumbnail: `website/assets/og.svg` (export to 240×240 PNG for PH logo if needed)
- Link: site at `https://new-world-coder.github.io/riskline/` once Pages is enabled
- One-pager PDF: open `website/one-pager.html` → Print → Save as PDF
