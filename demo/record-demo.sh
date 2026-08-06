#!/usr/bin/env bash
# Driven by asciinema for demo/riskline-demo.cast
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLI="${ROOT}/bin/riskline-cli"
if [[ ! -x "$CLI" ]]; then
  (cd "$ROOT" && go build -o bin/riskline-cli ./cmd/riskline-cli)
fi

export PS1='$ '
slow() { sleep "${1:-0.6}"; }

echo "# RiskLine — offline EU AI Act risk classification"
slow 0.8
echo "$ go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha"
slow 0.9

echo "$ ./bin/riskline-cli testdata/scenarios/prohibited-social-scoring.yaml"
slow 0.4
"$CLI" "$ROOT/testdata/scenarios/prohibited-social-scoring.yaml"
slow 1.4

echo "$ ./bin/riskline-cli testdata/scenarios/high-risk-hiring.yaml"
slow 0.4
"$CLI" "$ROOT/testdata/scenarios/high-risk-hiring.yaml"
slow 1.4

echo "$ ./bin/riskline-cli testdata/scenarios/limited-risk-chatbot.yaml"
slow 0.4
"$CLI" "$ROOT/testdata/scenarios/limited-risk-chatbot.yaml"
slow 1.2

echo "$ ./bin/riskline-cli --json testdata/scenarios/high-risk-hiring.yaml"
slow 0.4
"$CLI" --json "$ROOT/testdata/scenarios/high-risk-hiring.yaml" | head -n 24
slow 1.0

echo "# github.com/new-world-coder/riskline"
slow 0.8
