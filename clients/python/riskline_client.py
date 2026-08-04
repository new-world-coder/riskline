"""Stub Python client for riskline.

Replace by running the OpenAPI generator — see ../README.md.
"""

from __future__ import annotations

import json
import urllib.request
from typing import Any


def classify_system(base_url: str, body: dict[str, Any]) -> dict[str, Any]:
    url = base_url.rstrip("/") + "/v1/classify"
    data = json.dumps(body).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        headers={"content-type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))
