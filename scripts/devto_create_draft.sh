#!/usr/bin/env bash
# Create a *draft* on dev.to. Does not publish.
# Usage: DEVTO_API_KEY=... ./scripts/devto_create_draft.sh content/devto/01-....md
set -euo pipefail
file="${1:?markdown file}"
key="${DEVTO_API_KEY:?set DEVTO_API_KEY}"

python3 - "$file" "$key" <<'PY'
import json, sys, urllib.request
path, key = sys.argv[1], sys.argv[2]
raw = open(path, encoding="utf-8").read()
if not raw.startswith("---"):
    raise SystemExit("expected YAML front matter")
parts = raw.split("---", 2)
fm_lines = [l for l in parts[1].strip().splitlines() if l.strip()]
meta = {}
for line in fm_lines:
    k, _, v = line.partition(":")
    meta[k.strip()] = v.strip().strip('"')
body = parts[2].lstrip("\n")
tags = [t.strip() for t in meta.get("tags", "").split(",") if t.strip()]
payload = {
    "article": {
        "title": meta["title"],
        "body_markdown": body,
        "published": False,
        "tags": tags[:4],
        "canonical_url": meta.get("canonical_url") or None,
    }
}
req = urllib.request.Request(
    "https://dev.to/api/articles",
    data=json.dumps(payload).encode(),
    headers={"Content-Type": "application/json", "api-key": key},
    method="POST",
)
with urllib.request.urlopen(req) as resp:
    out = json.load(resp)
print("draft id:", out.get("id"))
print("url:", out.get("url") or out.get("path"))
PY
