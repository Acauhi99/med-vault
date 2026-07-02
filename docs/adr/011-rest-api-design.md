# ADR-011: REST API Design

## Status

Accepted

## Context

MedVault needs an API design for the frontend-backend communication. The API should be simple, consistent, and demonstrate modern API design.

## Decision

Use REST API design with JSON payloads.

## Consequences

### Positive
- Simple and well-understood
- HTTP methods map to CRUD operations
- Stateless
- Cacheable
- Easy to document (OpenAPI/Swagger)
- Works well with SPA architecture

### Negative
- Multiple round-trips for complex operations (acceptable for PoC)
- Over-fetching/under-fetching (but acceptable)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| GraphQL | More complex, overkill for PoC |
| gRPC | Not browser-friendly, requires protobuf |
| WebSockets | Not needed for this use case |

## References

- [REST API Design](https://restfulapi.net/)
- [AWS API Gateway Best Practices](https://docs.aws.amazon.com/apigateway/latest/developerguide/getting-started.html)
