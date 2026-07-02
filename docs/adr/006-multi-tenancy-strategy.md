# ADR-006: Multi-Tenancy Strategy

## Status

Accepted

## Context

MedVault is a multi-tenant platform. Each healthcare organization represents one tenant. Tenant isolation must be enforced across all layers. Users (especially doctors) may belong to multiple tenants with different roles in each.

## Decision

Use shared database, shared schema with `tenant_id` column isolation. Users exist independently of tenants — tenant membership is modeled via a `user_tenants` join table with a role per tenant.

**Schema design:**
- `users` — tenant-independent (no `tenant_id` column)
- `user_tenants` — join table linking users to tenants with a role
- All other tables — include `tenant_id` column for isolation

**Login flow:**
1. User authenticates (email + password) → receives list of available tenants
2. User selects a tenant → receives JWT with `tenant_id` + `role`

## Consequences

### Positive
- Users can belong to multiple tenants with different roles
- Doctor can work at Clinic A (role: doctor) and Hospital B (role: doctor) with one login
- Single database to manage
- Cost-effective for PoC
- Easy to demonstrate tenant isolation
- Row-Level Security as defense-in-depth
- Role is per-tenant, not global — correct for healthcare

### Negative
- Requires careful query design (always include tenant_id)
- Noisy neighbor risk (acceptable for PoC)
- Schema changes affect all tenants
- Login flow is two-step (authenticate + select tenant)
- JWT must carry tenant context

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| User-per-tenant (user has one tenant_id) | Cannot model doctors working at multiple clinics |
| Database per tenant | Higher operational cost, overkill for PoC |
| Schema per tenant | More complex, harder to demonstrate |
| Home tenant + cross-tenant grants | Over-engineering for PoC; many-to-many covers the need |

## References

- [AWS Multi-Tenancy Best Practices](https://docs.aws.amazon.com/whitepapers/latest/multi-tenant-saas-applications/multi-tenant-saas-applications.html)
- [DOMAIN.md](../DOMAIN.md) — User and UserTenant aggregates
