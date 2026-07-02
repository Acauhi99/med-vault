# ADR-006: Multi-Tenancy Strategy

## Status

Accepted

## Context

MedVault is a multi-tenant platform. Each healthcare organization represents one tenant. Tenant isolation must be enforced across all layers.

## Decision

Use shared database, shared schema with `tenant_id` column isolation.

## Consequences

### Positive
- Simplest multi-tenant model
- Single database to manage
- Cost-effective for PoC
- Easy to demonstrate tenant isolation
- Row-Level Security as defense-in-depth

### Negative
- Requires careful query design (always include tenant_id)
- Noisy neighbor risk (acceptable for PoC)
- Schema changes affect all tenants

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Database per tenant | Higher operational cost, overkill for PoC |
| Schema per tenant | More complex, harder to demonstrate |
| Hybrid | More complex, not needed for PoC scope |

## References

- [AWS Multi-Tenancy Best Practices](https://docs.aws.amazon.com/whitepapers/latest/multi-tenant-saas-applications/multi-tenant-saas-applications.html)
