# ADR-007: JWT for Authentication

## Status

Accepted

## Context

MedVault needs a stateless authentication mechanism for the SPA + API architecture. The mechanism should support multi-tenancy and role-based access control.

## Decision

Use JWT (JSON Web Tokens) for authentication.

## Consequences

### Positive
- Stateless (no server-side session storage)
- Contains tenant_id and role claims
- Short-lived access tokens (15 minutes)
- Refresh token rotation
- Easy to validate in middleware
- Works well with SPA architecture

### Negative
- Cannot invalidate tokens before expiration (but short-lived)
- Token size can grow (but acceptable)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Session cookies | Requires server-side storage |
| OAuth 2.0 | Overkill for PoC, no external identity providers |
| API keys | Not suitable for user authentication |

## References

- [JWT Documentation](https://jwt.io/introduction)
- [AWS Best Practices for JWT](https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-authentication-access-token-verification/)
