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

Every cloud resource is reproducible via Terraform. No manual console configuration. No snowflake servers. Infrastructure models platform capabilities, not individual AWS resources. Readability preferred over reducing duplication. See [INFRASTRUCTURE.md](INFRASTRUCTURE.md) and [ADR-003](adr/003-terraform-for-infrastructure.md).

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

---

## Design-First API

Define API contracts before writing code. OpenAPI is the single source of truth. Both backend (Go) and frontend (TypeScript) generate types from the same spec. Never manually write HTTP contracts. See [ADR-016](adr/016-design-first-api-documentation.md).

---

## Pragmatic Testing

Tests validate behavior, not implementation. Confidence over coverage. Fewer high-quality tests over many low-value tests. No E2E tests — unit + integration only. See [TESTING_STRATEGY.md](TESTING_STRATEGY.md).

---

## Quality Gates

Quality is enforced as early as possible. Every code change passes local validation before reaching CI. Pre-commit for formatting, pre-push for functional validation. Use Taskfile commands (`task lint`, `task test`, `task validate`) instead of technology-specific commands. See [QUALITY_GATES.md](QUALITY_GATES.md).
