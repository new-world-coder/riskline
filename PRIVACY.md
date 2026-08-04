# Privacy (stub)

This is a short public statement of data handling for early users. It is not a full privacy policy and will be replaced after legal review.

## CLI (default)

The `riskline-cli` binary classifies AI system descriptions **locally**. By default it does not phone home, upload your system inventory, or call a remote model. Input files and classification output stay on the machine where you run it.

## Hosted API (optional)

If you use a hosted API endpoint instead of the CLI:

- Classification requests (system purpose, data types, deployment context, etc.) are sent to the API server so it can return a risk tier and rationale.
- v1 classification is **deterministic rule evaluation** — it does not send your payload to an LLM for the risk decision itself.
- Metadata needed to operate the service (API keys, usage counters) may be stored in Postgres when that layer is enabled (`TODO(hosted):`).

Sub-processors and residency details will be listed here before any design-partner traffic is invited onto a hosted environment.
