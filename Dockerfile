# syntax=docker/dockerfile:1.7
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/traceforge-api ./cmd/api

FROM alpine:3.24
RUN addgroup -S traceforge && adduser -S -G traceforge traceforge
COPY --from=build /out/traceforge-api /usr/local/bin/traceforge-api
USER traceforge
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 CMD wget -qO- http://127.0.0.1:8080/health/ready || exit 1
ENTRYPOINT ["traceforge-api"]
