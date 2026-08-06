/**
 * RiskLine landing terminal demo.
 * Output strings captured from the real CLI (v0.1.0-alpha / eu-ai-act-2024-v0.1.0).
 */
(function () {
  "use strict";

  const SCENARIOS = {
    prohibited: {
      tier: "prohibited",
      label: "prohibited",
      ref: "Article 5(1)(c)",
      command: "riskline-cli testdata/scenarios/prohibited-social-scoring.yaml",
      output: [
        "System:           Citizen Score",
        "Risk tier:        prohibited",
        "Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)",
        "",
        "Rationale",
        "Classified as prohibited because 1 rule(s) matched. Highest-severity matches: Article 5(1)(c) (prohibited-social-scoring).",
        "",
        "Matched rules",
        "ID                         TIER        REF",
        "prohibited-social-scoring  prohibited  Article 5(1)(c)",
        "",
        "Recommended controls",
        "  - Do not deploy; redesign the use case away from social scoring",
        "  - Document why the prior design was discontinued",
      ],
    },
    high_risk: {
      tier: "high_risk",
      label: "high_risk",
      ref: "Annex III (4)(a)",
      command: "riskline-cli testdata/scenarios/high-risk-hiring.yaml",
      output: [
        "System:           Hiring Assist",
        "Risk tier:        high_risk",
        "Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)",
        "",
        "Rationale",
        "Classified as high_risk because 1 rule(s) matched. Highest-severity matches: Annex III (4)(a) (high-risk-recruitment).",
        "",
        "Matched rules",
        "ID                     TIER       REF",
        "high-risk-recruitment  high_risk  Annex III (4)(a)",
        "",
        "Recommended controls",
        "  - Establish a risk management system and data governance for the system",
        "  - Ensure human oversight of ranking/filtering decisions",
        "  - Keep logs and technical documentation for conformity assessment",
        "  - Assess bias and discrimination risks in candidate scoring",
      ],
    },
    limited_risk: {
      tier: "limited_risk",
      label: "limited_risk",
      ref: "Article 50(1)",
      command: "riskline-cli testdata/scenarios/limited-risk-chatbot.yaml",
      output: [
        "System:           Support Bot",
        "Risk tier:        limited_risk",
        "Ruleset:          eu-ai-act-2024-v0.1.0 (2026-08-04)",
        "",
        "Rationale",
        "Classified as limited_risk because 1 rule(s) matched. Highest-severity matches: Article 50(1) (limited-risk-chatbot-transparency).",
        "",
        "Matched rules",
        "ID                                 TIER          REF",
        "limited-risk-chatbot-transparency  limited_risk  Article 50(1)",
        "",
        "Recommended controls",
        "  - Disclose AI interaction to end users at the start of the conversation",
        "  - Keep disclosure clear and accessible",
      ],
    },
  };

  const ORDER = ["prohibited", "high_risk", "limited_risk"];
  const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const body = document.getElementById("terminal-body");
  const badge = document.getElementById("tier-badge");
  const buttons = Array.from(document.querySelectorAll(".scenario-btn"));
  const heroPlane = document.getElementById("hero-plane");
  const copyBtn = document.getElementById("copy-install");
  const installCmd = document.getElementById("copy-install")
    ? document.getElementById("install-cmd")
    : null;

  let runToken = 0;
  let autoTimer = null;

  function escapeHtml(s) {
    return s
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");
  }

  function colorizeOutput(lines, tier) {
    return lines
      .map((line) => {
        if (line.startsWith("Risk tier:")) {
          const prefix = escapeHtml("Risk tier:        ");
          const value = escapeHtml(tier);
          return prefix + '<span class="tier-' + tier + '">' + value + "</span>";
        }
        return escapeHtml(line);
      })
      .join("\n");
  }

  function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  async function typeText(token, text, speed) {
    let html = body.innerHTML;
    // strip trailing cursor if present
    html = html.replace(/<span class="cursor"><\/span>$/, "");
    for (let i = 0; i < text.length; i++) {
      if (token !== runToken) return;
      html += escapeHtml(text[i]);
      body.innerHTML = html + '<span class="cursor"></span>';
      body.scrollTop = body.scrollHeight;
      if (!reduceMotion) await sleep(speed);
    }
  }

  async function playScenario(key, { autoAdvance } = { autoAdvance: false }) {
    const scenario = SCENARIOS[key];
    if (!scenario || !body) return;

    runToken += 1;
    const token = runToken;

    buttons.forEach((btn) => {
      btn.setAttribute("aria-pressed", btn.dataset.scenario === key ? "true" : "false");
    });

    if (badge) {
      badge.className = "tier-badge " + scenario.tier;
      badge.textContent = scenario.label + " · " + scenario.ref;
      badge.classList.remove("visible");
      badge.setAttribute("aria-hidden", "true");
    }

    const prompt = '<span class="prompt">$</span> <span class="cmd">';
    body.innerHTML = prompt + '<span class="cursor"></span>';

    if (reduceMotion) {
      body.innerHTML =
        prompt +
        escapeHtml(scenario.command) +
        "</span>\n\n" +
        colorizeOutput(scenario.output, scenario.tier);
      if (badge) {
        badge.classList.add("visible");
        badge.setAttribute("aria-hidden", "false");
      }
      return;
    }

    await typeText(token, scenario.command, 18);
    if (token !== runToken) return;

    body.innerHTML =
      prompt + escapeHtml(scenario.command) + "</span>\n" + '<span class="cursor"></span>';
    await sleep(280);
    if (token !== runToken) return;

    const colored = colorizeOutput(scenario.output, scenario.tier);
    const chunks = colored.split("\n");
    let built =
      prompt + escapeHtml(scenario.command) + "</span>\n\n";
    for (let i = 0; i < chunks.length; i++) {
      if (token !== runToken) return;
      built += (i === 0 ? "" : "\n") + chunks[i];
      body.innerHTML = built + '<span class="cursor"></span>';
      body.scrollTop = body.scrollHeight;
      await sleep(28);
    }

    if (badge) {
      badge.classList.add("visible");
      badge.setAttribute("aria-hidden", "false");
    }

    if (autoAdvance) {
      clearTimeout(autoTimer);
      autoTimer = setTimeout(() => {
        const idx = ORDER.indexOf(key);
        const next = ORDER[(idx + 1) % ORDER.length];
        playScenario(next, { autoAdvance: true });
      }, 4200);
    }
  }

  buttons.forEach((btn) => {
    btn.addEventListener("click", () => {
      clearTimeout(autoTimer);
      playScenario(btn.dataset.scenario, { autoAdvance: false });
    });
  });

  // Soft hero parallax
  if (heroPlane && !reduceMotion) {
    const heroOutput = document.getElementById("hero-output");
    window.addEventListener(
      "scroll",
      () => {
        const y = Math.min(window.scrollY, 400);
        const offset = y * 0.18 + "px";
        heroPlane.style.setProperty("--parallax", offset);
        if (heroOutput) heroOutput.style.setProperty("--parallax", offset);
      },
      { passive: true }
    );
  }

  if (copyBtn && installCmd) {
    copyBtn.addEventListener("click", async () => {
      const text = installCmd.textContent.trim();
      try {
        await navigator.clipboard.writeText(text);
        copyBtn.textContent = "copied";
        setTimeout(() => {
          copyBtn.textContent = "copy";
        }, 1600);
      } catch {
        copyBtn.textContent = "failed";
        setTimeout(() => {
          copyBtn.textContent = "copy";
        }, 1600);
      }
    });
  }

  // Start demo when section is near viewport
  const demo = document.getElementById("demo");
  if (demo && "IntersectionObserver" in window) {
    const io = new IntersectionObserver(
      (entries) => {
        if (entries.some((e) => e.isIntersecting)) {
          playScenario("prohibited", { autoAdvance: !reduceMotion });
          io.disconnect();
        }
      },
      { threshold: 0.25 }
    );
    io.observe(demo);
  } else {
    playScenario("prohibited", { autoAdvance: !reduceMotion });
  }
})();
