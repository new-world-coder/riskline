/**
 * Stub TypeScript client for riskline.
 * Replace by running the OpenAPI generator — see ../README.md.
 */
export type RiskTier = "prohibited" | "high_risk" | "limited_risk" | "minimal_risk";

export interface ClassifyRequest {
  name?: string;
  purpose: string;
  data_types: string[];
  deployment_context: string;
  autonomy_level: string;
  affected_population: string;
  geographic_scope?: string;
  /** Regime pack IDs (e.g. eu-ai-act). Not the same as geographic_scope. */
  regimes?: string[];
}

export interface TechnicalControl {
  id: string;
  paper_ref: string;
  summary: string;
  technical_hook: string;
  evidence_type: string;
}

export interface RegimeClassification {
  regime: string;
  character: string;
  risk_tier: RiskTier;
  ruleset_version: string;
  last_updated: string;
  matched_rules: unknown[];
  rationale: string;
  recommended_controls: string[];
  technical_controls?: TechnicalControl[];
  judgment_calls?: string[];
  mapping_only?: boolean;
}

export interface ClassifyResponse {
  risk_tier: RiskTier;
  ruleset_version: string;
  last_updated: string;
  disclaimer: string;
  rationale: string;
  matched_rules: unknown[];
  recommended_controls: string[];
  technical_controls?: TechnicalControl[];
  regime?: string;
  mapping_only?: boolean;
  classifications?: RegimeClassification[];
}

export async function classifySystem(
  baseUrl: string,
  body: ClassifyRequest
): Promise<ClassifyResponse> {
  const res = await fetch(`${baseUrl.replace(/\/$/, "")}/v1/classify`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`classify failed: ${res.status} ${await res.text()}`);
  }
  return res.json() as Promise<ClassifyResponse>;
}
