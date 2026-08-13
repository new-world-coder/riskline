# Architecture notes

```
CLI / API ──► pkg/config (regimes) ──► pkg/engine ──► pkg/ruleset.Loader ──► embedded pack JSON
```

- `api/openapi.yaml` is the public contract (`info.version` tracks API surface; pack versions are independent).
- `pkg/schema` mirrors that contract in Go.
- `pkg/ruleset.Loader` maps regime IDs (`eu-ai-act`) to embedded JSON packs (`eu-ai-act-2024-v0.1.0`).
- Evaluation is in `pkg/engine`. Top-level response fields always reflect the **first** resolved regime for backward compatibility.
- `classifications[]` is only populated when more than one regime is evaluated — keeps the EU golden JSON stable.
- `geographic_scope` is a request hint, **not** a regime selector.
- Regime resolution precedence: request / `-regimes` → `RISKLINE_REGIMES` → `.riskline.yaml` → `[eu-ai-act]`.
- `internal/httpserver` and both `cmd/*` entrypoints are thin.

v1 does not use OPA/Rego. The rule surface is small and we wanted zero runtime policy-engine dependency for an offline CLI. Revisit if the ruleset grows into something Rego expresses more clearly than JSON conditions.
