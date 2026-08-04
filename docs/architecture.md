# Architecture notes

```
CLI ──┐
      ├──► pkg/engine ──► pkg/ruleset (embedded JSON)
API ──┘
```

- `api/openapi.yaml` is the public contract.
- `pkg/schema` mirrors that contract in Go.
- `pkg/ruleset` loads versioned rule documents; matching conditions are data, evaluation is in `pkg/engine`.
- `internal/httpserver` and both `cmd/*` entrypoints are thin.

v1 does not use OPA/Rego. The rule surface is small and we wanted zero runtime policy-engine dependency for an offline CLI. Revisit if the ruleset grows into something Rego expresses more clearly than JSON conditions.
