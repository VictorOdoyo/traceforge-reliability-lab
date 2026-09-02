# SLO Operations

## Model

TraceForge models a service-level objective with a target percentage, rolling window, good-event count, and total-event count. Attainment is the ratio of good events to total events. The error budget is the allowed bad-event share implied by the target.

For a `99.9%` target, the allowed failure ratio is `0.1%`. Consuming half of that allowance leaves 50% of the error budget.

## Investigation Sequence

1. Sort SLOs by remaining budget and identify the fastest burn.
2. Compare error rate and p95 latency in service health.
3. Filter traces to the SLO service and error status.
4. Inspect deployment correlations within the degradation window.
5. Review alerts and active incidents for the same service.
6. Decide whether to halt releases, mitigate, or continue monitoring.

## Operational Policy

- Above 50% remaining: normal delivery policy.
- Between 20% and 50%: review risky releases and increase monitoring.
- Below 20%: prioritize reliability work and require explicit release approval.
- Exhausted: freeze nonessential changes until the service returns within policy.

These thresholds are examples for the lab. Production policy must reflect business impact, traffic patterns, and the SLO window.

## Sampling Considerations

Tail sampling retains errors and slow traces before applying baseline probabilistic sampling. This improves diagnostic coverage but sampled traces are not a replacement for unsampled metric counters when calculating formal SLO compliance.
