# curl examples

```bash
# start the API
go run ./cmd/riskline-api -addr :8080

curl -s localhost:8080/v1/classify \
  -H 'content-type: application/json' \
  -d @hiring-assist.json | jq .
```
