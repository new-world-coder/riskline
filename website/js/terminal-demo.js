const INSTALL_LINE = "$ go install github.com/new-world-coder/riskline/cmd/riskline-cli@v0.1.0-alpha";

const SCENARIOS = {
  prohibited: {
    label: "Citizen Score",
    tier: "prohibited",
    tierClass: "tier-prohibited",
    file: "testdata/scenarios/prohibited-social-scoring.yaml",
    output: `System:           Citizen Score
Risk tier:        prohibited
Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)

Rationale
Classified as prohibited because 1 rule(s) matched. Highest-severity matches: Article 5(1)(c) (prohibited-social-scoring).

Matched rules
ID                         TIER        REF
prohibited-social-scoring  prohibited  Article 5(1)(c)

Recommended controls
  - Do not deploy; redesign the use case away from social scoring
  - Document why the prior design was discontinued`,
  },
  high: {
    label: "Hiring Assist",
    tier: "high_risk",
    tierClass: "tier-high",
    file: "testdata/scenarios/high-risk-hiring.yaml",
    output: `System:           Hiring Assist
Risk tier:        high_risk
Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)

Rationale
Classified as high_risk because 1 rule(s) matched. Highest-severity matches: Annex III (4)(a) (high-risk-recruitment).

Matched rules
ID                     TIER       REF
high-risk-recruitment  high_risk  Annex III (4)(a)

Recommended controls
  - Establish a risk management system and data governance for the system
  - Ensure human oversight of ranking/filtering decisions
  - Keep logs and technical documentation for conformity assessment
  - Assess bias and discrimination risks in candidate scoring`,
  },
  limited: {
    label: "Support Bot",
    tier: "limited_risk",
    tierClass: "tier-limited",
    file: "testdata/scenarios/limited-risk-chatbot.yaml",
    output: `System:           Support Bot
Risk tier:        limited_risk
Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)

Rationale
Classified as limited_risk because 1 rule(s) matched. Highest-severity matches: Article 50(1) (limited-risk-chatbot-transparency).

Matched rules
ID                                 TIER          REF
limited-risk-chatbot-transparency  limited_risk  Article 50(1)

Recommended controls
  - Disclose AI interaction to end users at the start of the conversation
  - Keep disclosure clear and accessible`,
  },
};

function escapeHtml(text) {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

function colorizeOutput(text, tierClass) {
  return escapeHtml(text).replace(
    /Risk tier:\s+(\S+)/,
    (_, tier) => `Risk tier:        <span class="${tierClass}">${tier}</span>`
  );
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function typeLine(el, html, charDelay = 8) {
  el.insertAdjacentHTML("beforeend", '<span class="line"></span>');
  const line = el.lastElementChild;
  for (let i = 0; i < html.length; i++) {
    line.innerHTML = html.slice(0, i + 1);
    await sleep(charDelay);
  }
  el.insertAdjacentHTML("beforeend", "\n");
}

async function runScenario(key) {
  const body = document.getElementById("terminal-output");
  const scenario = SCENARIOS[key];
  if (!body || !scenario) return;

  body.innerHTML = "";
  document.querySelectorAll(".scenario-tab").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.scenario === key);
  });

  await typeLine(body, `<span class="dim">${escapeHtml(INSTALL_LINE)}</span>`, 4);
  await sleep(300);
  await typeLine(
    body,
    `<span class="prompt">$</span> riskline-cli ${escapeHtml(scenario.file)}`,
    10
  );
  await sleep(200);
  await typeLine(
    body,
    colorizeOutput(scenario.output, scenario.tierClass),
    3
  );
}

function initTerminalDemo() {
  document.querySelectorAll(".scenario-tab").forEach((btn) => {
    btn.addEventListener("click", () => runScenario(btn.dataset.scenario));
  });
  runScenario("prohibited");
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initTerminalDemo);
} else {
  initTerminalDemo();
}
