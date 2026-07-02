# ADR-007: JWT for Authentication

## Status

Accepted

## Context

MedVault needs a stateless authentication mechanism for the SPA + API architecture. The mechanism should support multi-tenancy (users can belong to multiple tenants with different roles) and role-based access control.

## Decision

Use JWT (JSON Web Tokens) for authentication with a two-step login flow.

### Login Flow

```
Step 1 — Authenticate:
POST /auth/login { email, password }
→ Temporary JWT (no tenant) + list of available tenants with roles

Step 2 — Select tenant:
POST /auth/select-tenant { tenant_id }  (with temporary JWT)
→ Final JWT { user_id, tenant_id, role }
```

### JWT Claims (Final Token)

```json
{
  "sub": "user_id",
  "tenant_id": "selected_tenant_id",
  "role": "doctor",
  "iat": 1704067200,
  "exp": 1704068100
}
```

The `role` comes from the `user_tenants` table, not from the user record.

## Consequences

### Positive
- Stateless (no server-side session storage)
- Contains tenant_id and role claims for the selected tenant
- Users can switch tenants without re-authenticating (just call select-tenant again)
- Short-lived access tokens (15 minutes)
- Refresh token rotation
- Easy to validate in middleware
- Works well with SPA architecture

### Negative
- Cannot invalidate tokens before expiration (but short-lived)
- Token size can grow (but acceptable)
- Login flow is two-step (minor UX friction)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Session cookies | Requires server-side storage |
| OAuth 2.0 | Overkill for PoC, no external identity providers |
| API keys | Not suitable for user authentication |
| Single JWT with all tenant_ids | Token bloat, exposes tenant list to client |

## References

- [JWT Documentation](https://jwt.io/introduction)
- [AWS Best Practices for JWT](https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-authentication-access-token-verification/)
- [DOMAIN.md](../DOMAIN.md) — User and UserTenant aggregates
