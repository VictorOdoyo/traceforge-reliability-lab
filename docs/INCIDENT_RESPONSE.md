# Incident Response

## Purpose

This runbook describes how an operator uses TraceForge to investigate a reliability regression and maintain a useful incident record.

## Triage

1. Confirm the alert is firing and record an owner by acknowledging it.
2. Check service health for elevated error rate, p95 latency, and request volume.
3. Open the trace explorer with the affected service and `error` status selected.
4. Inspect representative span waterfalls for the first failing dependency.
5. Compare deployment evidence inside the incident window.
6. Review dependency edges for a shared upstream or downstream service.

## Incident Lifecycle

Use transitions in one direction:

```text
investigating -> identified -> monitoring -> resolved
```

An incident may move directly from `investigating` or `identified` to `resolved` when the cause and recovery are immediate. Resolved incidents cannot accept new timeline events. Each non-resolution transition appends an actor-attributed timeline event.

## Evidence Checklist

- Alert start time and owner
- Affected service and customer-facing operation
- Representative trace identifiers
- Error and latency change from baseline
- Suspected deployment and commit
- Dependency or region concentration
- Mitigation action and verification query
- SLO budget impact

## Mitigation

Prefer reversible changes. Roll back a correlated deployment, disable a risky path, reduce load, or isolate an unhealthy dependency. Re-run the same trace query after mitigation so the comparison uses equivalent filters.

## Resolution

Move to `monitoring` after the primary signal recovers. Resolve only after error rate and latency remain healthy for an agreed observation window and no new correlated alerts fire. Preserve trace IDs and deployment evidence for the post-incident review.
