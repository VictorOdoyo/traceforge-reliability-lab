# Architecture

## Design Goals

TraceForge is designed to make reliability evidence easy to query without coupling operational workflows to one transport or storage engine. The system favors explicit domain transitions, bounded query inputs, observable request handling, and a local mode that remains useful when infrastructure is unavailable.

## Components

### Operations Console

The React console loads dashboard resources concurrently and presents them as investigation surfaces: service health, latency trends, dependency topology, traces, SLOs, deployment evidence, alerts, and incidents. Large chart and graph dependencies are loaded lazily so the initial route remains small. Mutations use the same versioned API as queries.

### HTTP API

The Go API owns validation, orchestration, and transport concerns. Middleware adds request identifiers, structured access logs, panic recovery, CORS rules, security headers, Prometheus instrumentation, and OpenTelemetry spans. Handlers translate JSON contracts into domain commands and query filters.

### Domain and Services

Domain types enforce trace assembly, SLO math, alert ownership, incident transitions, suppressions, and sampling decisions. Application services calculate reliability summaries, dependency edges, and deployment correlations. They depend only on the repository interface.

### Storage

The in-memory adapter provides a connected synthetic dataset for development and tests. When `CLICKHOUSE_DSN` is configured, the ClickHouse adapter creates a MergeTree trace table with time-based partitioning and TTL retention, persists ingested traces, and uses the memory adapter for control-plane records.

### Telemetry Pipeline

The API exports OTLP spans to the collector. Tail sampling retains errors and slow requests while preserving a baseline sample. ClickHouse stores accepted spans, Prometheus scrapes API and collector metrics, and Grafana provisions a reliability dashboard automatically.

## Request Flow

1. The console requests versioned resources through its same-origin `/api` path.
2. Nginx proxies requests to the Go API in container deployments.
3. Middleware records metrics, propagates trace context, and recovers panics.
4. Handlers validate query parameters or command bodies.
5. Services query the repository port and apply domain rules.
6. The response layer emits JSON or a stable problem envelope.

## Failure Model

- Invalid input returns a `400` response without reaching storage.
- Missing resources return `404` and preserve repository state.
- Repository failures return `500` with no internal details exposed.
- Health endpoints remain separate from business queries.
- Kubernetes readiness removes unhealthy API instances from service routing.
- Graceful shutdown allows in-flight requests to finish within a fixed budget.

## Deployment Boundaries

The API and console are stateless and can scale horizontally. ClickHouse owns persistent trace data. The collector is independently scalable when throughput increases. Production deployments should replace example credentials, configure TLS at ingress, and use managed secrets and storage classes.
