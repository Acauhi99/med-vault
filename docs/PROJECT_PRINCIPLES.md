# Engineering Principles

Principles that guide every decision in MedVault.

---

## Security First

Security is considered before functionality. If a choice trades security for convenience, choose security. See [SECURITY.md](SECURITY.md) for implementation details.

---

## Architecture Is the Product

This project demonstrates architectural thinking. Code quality matters, but architectural consistency matters more. Every change should be explainable during a technical review.

---

## Documentation Before Implementation

Understand the problem before writing code. Read the relevant docs. If docs are missing, write them first.

---

## Multi-Tenant by Default

Every request belongs to exactly one tenant. No query, no API call, no storage operation ignores tenant isolation. See [DOMAIN.md](DOMAIN.md) for enforcement points.

---

## Infrastructure as Code

Every cloud resource is reproducible via Terraform. No manual console configuration. No snowflake servers. See [ADR-003](adr/003-terraform-for-infrastructure.md).

---

## Cloud Native

Prefer managed AWS services over self-managed infrastructure. The goal is demonstrating architecture, not operations. See [ADR-012](adr/012-modular-monolith-architecture.md).

---

## Explicit Over Clever

Readable code beats clever code. Explicit error handling beats silent failures. Named constants beat magic numbers. If a colleague cannot understand it in 30 seconds, simplify it.

---

## Bounded Contexts Are Boundaries

Respect DDD bounded contexts. Do not leak domain logic across boundaries. Do not share database tables across contexts without explicit design. See [DOMAIN.md](DOMAIN.md).

---

## Simplicity Is a Feature

The simplest solution that demonstrates the concept is the right solution. Complexity must justify itself. If in doubt, choose the boring option.

---

## Decision Traceability

Every significant technology choice is documented as an ADR in [docs/adr/](adr/). If you make a decision, record it.

---

## Feature Isolation

Frontend features are self-contained business capabilities. Each feature owns its components, hooks, services, schemas, and types. Avoid unnecessary coupling between features. Shared code exists only when there is a clear cross-feature need. See [ADR-015](adr/015-frontend-feature-based-architecture.md).
