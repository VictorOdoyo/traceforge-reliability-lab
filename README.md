# TraceForge Reliability Lab

TraceForge is an observability and incident-investigation platform for exploring distributed traces, service dependencies, SLOs, error budgets, alerts, and deployment correlations.

## Architecture

```text
apps/console/        React and TypeScript reliability console
cmd/api/             Go HTTP API entry point
internal/domain/     Reliability and telemetry domain model
internal/platform/   Storage, ingestion, and observability adapters
internal/service/    Application workflows and query services
deploy/              Local and Kubernetes platform definitions
```

TraceForge is designed to run locally with synthetic telemetry. It contains no private production traces, credentials, or customer data.
