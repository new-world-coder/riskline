# An open JSON schema for AI system risk classification

The useful part of RiskLine is not the binary. It is the contract.

Everything — CLI table output, `POST /v1/classify`, future generated clients — hangs off [`api/openapi.yaml`](https://github.com/new-world-coder/riskline/blob/main/api/openapi.yaml). If the handler and the doc disagree, that is a bug, not a documentation problem.

## Request: describe the system, don't upload your codebase

```json
{
  "name": "Hiring Assist",
  "purpose": "Screen job applicants and rank candidates for interview",
  "data_types": ["personal_data", "employment_data"],
  "deployment_context": "saas_b2b",
  "autonomy_level": "decision_support",
  "affected_population": "job_applicants",
  "geographic_scope": "eu"
}
```

Enums over free-text for the structural fields. Purpose stays free text because real systems do not fit a taxonomy on day one — and because keyword heuristics need something to chew on. Flags like `social_scoring` and `real_time_remote_biometric_id` exist for Article 5 cases where "infer from prose" is a bad idea.

## Response: tier is not enough

A label without a clause reference is how tools become shelfware. Required fields on every success response:

| Field | Why it is mandatory |
|-------|---------------------|
| `risk_tier` | The thing you came for |
| `ruleset_version` | Regulations move; "we said high-risk in 2024" is useless without the ruleset id |
| `last_updated` | Same — perishable mappings |
| `matched_rules[]` | Which article/annex fired, paraphrased |
| `rationale` | Plain-language why |
| `recommended_controls` | Next engineering actions, not a policy essay |
| `disclaimer` | Advisory tool, not legal advice — tested so refactors cannot drop it |

`judgment_calls` is optional but important. When the ruleset made an interpretation the Act text does not fully decide, we say so in the payload instead of burying it in a blog post.

## What I want feedback on

1. Are the request enums too coarse for your inventory format?
2. Should `geographic_scope: non_eu` short-circuit differently, or stay "still classify, human decides applicability"?
3. Is `recommended_controls` useful, or should it move to a separate optional profile?

Open a GitHub issue. Schema feedback is a contribution, not noise — especially [ruleset updates](https://github.com/new-world-coder/riskline/issues/new?template=ruleset_update.md) and [misclassifications](https://github.com/new-world-coder/riskline/issues/new?template=misclassification_report.md).

Repo: [github.com/new-world-coder/riskline](https://github.com/new-world-coder/riskline).
