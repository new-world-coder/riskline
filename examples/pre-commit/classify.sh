#!/usr/bin/env bash
# pre-commit style hook: classify a declared AI system description
set -euo pipefail
FILE="${1:-ai-system.yaml}"
if [[ ! -f "$FILE" ]]; then
  echo "no $FILE — skip"
  exit 0
fi
riskline-cli --json "$FILE"
