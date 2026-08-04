# Generated clients

These directories hold **stubs** until OpenAPI Generator is wired in CI.

Regenerate (once `openapi-generator-cli` is installed):

```bash
npx @openapitools/openapi-generator-cli generate \
  -i api/openapi.yaml -g typescript-fetch -o clients/typescript

openapi-generator-cli generate \
  -i api/openapi.yaml -g python -o clients/python
```

Do not hand-edit generated sources — change `api/openapi.yaml` and regenerate.
