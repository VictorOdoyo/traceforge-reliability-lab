# Security Policy

## Supported Version

Security fixes are applied to the latest commit on `main`.

## Reporting a Vulnerability

Do not open a public issue for a suspected vulnerability. Use GitHub private vulnerability reporting for this repository and include the affected component, reproduction steps, impact, and any proposed mitigation. Expect an acknowledgement within five business days.

## Security Posture

- Runtime containers use non-root users where supported.
- Kubernetes workloads drop Linux capabilities and disable privilege escalation.
- API responses do not expose internal storage errors.
- HTTP middleware applies content-type, frame, sniffing, and referrer protections.
- Query limits and request-size validation bound common resource-exhaustion paths.
- CI performs formatting, static analysis, tests, manifest rendering, and container builds.
- Dependabot monitors Go, npm, Docker, and GitHub Actions dependencies.

## Deployment Responsibilities

Credentials in Compose and the Kubernetes base are development examples. Before deployment, use an external secret manager, rotate all passwords, terminate TLS at a trusted ingress, restrict ClickHouse network access, configure image digests, apply namespace network policies, and review retention settings for sensitive telemetry attributes.

Trace attributes can contain customer or authentication data. Production instrumentation must redact secrets and minimize personally identifiable information before export.
