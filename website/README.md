# RiskLine website

Static GitHub Pages site for RiskLine.

- [`index.html`](index.html) — landing page with embedded terminal demo
- [`one-pager.html`](one-pager.html) — printable overview (Print → PDF)
- [`assets/og.png`](assets/og.png) — social / Product Hunt thumbnail

Published by [`.github/workflows/pages.yml`](../.github/workflows/pages.yml) from this directory after Pages is enabled for the repo (Settings → Pages → GitHub Actions).

Local preview:

```bash
python3 -m http.server -d website 8080
```
