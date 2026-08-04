# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/riskline-api ./cmd/riskline-api \
 && CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/riskline-cli ./cmd/riskline-cli

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/riskline-api /usr/local/bin/riskline-api
COPY --from=build /out/riskline-cli /usr/local/bin/riskline-cli
EXPOSE 8080
USER nobody
ENTRYPOINT ["/usr/local/bin/riskline-api"]
CMD ["-addr", ":8080"]
