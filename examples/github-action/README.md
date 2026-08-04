# Example: classify in CI

```yaml
# .github/workflows/ai-act.yml
name: EU AI Act classify
on:
  pull_request:
    paths:
      - "ai-system.yaml"
jobs:
  classify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: new-world-coder/riskline/examples/github-action@main
        with:
          file: ai-system.yaml
```

Until the action is published from a real org, point `uses:` at a local path or install the CLI with `go install` as in `action.yml`.
