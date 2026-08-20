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

## Example: assure gate (classify + diff + assure)

See `examples/github-action/assure-gate.yml` for a minimal workflow that:
- classifies a checked-in system description (`ai-system.yaml`)
- diffs against a checked-in baseline (`ai-system-baseline.yaml`)
- runs assurance using probe results (`probes.json`)
- fails the workflow if `conformity_state` is `red` or if the material-change impact is `reclassify`
