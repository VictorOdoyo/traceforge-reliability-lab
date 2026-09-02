# TraceForge Reliability Lab

[![CI](https://github.com/VictorOdoyo/traceforge-reliability-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/VictorOdoyo/traceforge-reliability-lab/actions/workflows/ci.yml)

TraceForge is a production-oriented reliability investigation platform for distributed traces, service dependencies, SLOs, error budgets, alerts, incidents, and deployment correlations. It pairs a Go API and ClickHouse trace store with a dense React operations console and a complete local observability stack.

The included dataset is synthetic and generated at startup. No cloud account, private credential, customer data, or paid service is required.

## Capabilities

- Search and filter traces by service, status, duration, time window, and result limit
- Inspect span waterfalls and cross-service request paths
- Aggregate service health, p95 latency, error rate, and request volume
- Visualize service dependencies and edge-level latency
- Track SLO attainment, error-budget consumption, and burn risk
- Correlate deployments with reliability regressions
- Acknowledge alerts and assign operational ownership
- Advance incidents through controlled lifecycle transitions
- Export Prometheus API metrics and OpenTelemetry traces
- Apply error-first, latency-first, and baseline tail-sampling policies
- Run with in-memory seed data or persistent ClickHouse storage

## Architecture

```text
Browser
  | same-origin HTTP
  v
React console -> Go API -> domain services -> repository port
                              |                 |-- in-memory demo adapter
                              |                 `-- ClickHouse adapter
                              v
                     OpenTelemetry Collector -> ClickHouse
                              |
                              `-> Prometheus -> Grafana
```

The domain and service layers do not depend on HTTP or ClickHouse. Repository ports keep query and incident workflows testable without infrastructure, while adapters provide local and persistent runtime modes.

## Repository Layout

```text
apps/console/                 React and TypeScript operations console
cmd/api/                      Go API entry point
internal/config/              Validated runtime configuration
internal/domain/              Traces, SLOs, alerts, and incident rules
internal/httpapi/             Versioned REST handlers and middleware
internal/platform/            ClickHouse, metrics, and telemetry adapters
internal/service/             Reliability query and command workflows
internal/seed/                Connected synthetic telemetry dataset
deploy/grafana/               Provisioned dashboards and data sources
deploy/k8s/base/              Kubernetes resources and autoscaling
deploy/otel-collector/        OTLP receivers and tail-sampling pipeline
deploy/prometheus/            Metrics scrape configuration
```

## Requirements

- Go 1.27+
- Node.js 20.20+
- pnpm 9.15+
- Docker and Docker Compose for the full stack
- kubectl for rendering Kubernetes manifests

## Quick Start

Start the API with the synthetic in-memory repository:

```bash
go run ./cmd/api
```

In another terminal, install and run the console:

```bash
pnpm install --frozen-lockfile
pnpm dev
```

Open `http://localhost:5173`. Vite proxies `/api` requests to the API on port `8080`.

## Full Observability Stack

```bash
docker compose up --build
```

| Service | URL |
| --- | --- |
| TraceForge console | http://localhost:4173 |
| Go API | http://localhost:8080 |
| API metrics | http://localhost:8080/metrics |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| ClickHouse HTTP | http://localhost:8123 |

Grafana is provisioned with the `TraceForge API Overview` dashboard. The local admin password is `traceforge`; change it outside disposable development environments.

## API Examples

```bash
curl http://localhost:8080/health/ready
curl "http://localhost:8080/api/v1/traces?service=payment-router&status=error&limit=10"
curl http://localhost:8080/api/v1/deployments/correlations
curl -X POST http://localhost:8080/api/v1/alerts/alert-301/acknowledge \
  -H "Content-Type: application/json" \
  -d '{"owner":"Victor"}'
```

## Quality Checks

```bash
go test -race -cover ./...
go vet ./...
pnpm lint
pnpm typecheck
pnpm test
pnpm build
docker compose config --quiet
kubectl kustomize deploy/k8s/base
```

Vitest enforces coverage thresholds for shared console logic. CI also enforces an overall Go statement coverage floor and builds both production containers.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | API listen port |
| `READ_TIMEOUT` | `5s` | HTTP read-header timeout |
| `WRITE_TIMEOUT` | `15s` | HTTP response timeout |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown budget |
| `CLICKHOUSE_DSN` | empty | Enables persistent trace storage |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | empty | Enables OTLP trace export |
| `VITE_API_URL` | empty | Optional console API origin override |

See [architecture](docs/ARCHITECTURE.md), [incident response](docs/INCIDENT_RESPONSE.md), [SLO operations](docs/SLO_OPERATIONS.md), and [security](SECURITY.md) for operational detail.

## License

MIT
