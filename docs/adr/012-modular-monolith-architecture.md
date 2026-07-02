# ADR-012: Modular Monolith Architecture

## Status

Accepted

## Context

MedVault backend needs an architecture that demonstrates clean separation of concerns while remaining easy to evolve. The architecture should support future decomposition into microservices if the service becomes highly demanded.

## Decision

Use a modular monolith architecture. Each bounded context is a self-contained module with its own domain, application, and infrastructure layers. Modules communicate via in-process interfaces.

## Module Structure

Module names align with the DDD Bounded Context names defined in [DOMAIN.md](../DOMAIN.md):

| Module | Bounded Context |
|--------|----------------|
| `auth` | Identity & Access |
| `clinical` | Clinical |
| `imaging` | Imaging |
| `audit` | Audit |
| `shared` | Shared Kernel (cross-cutting) |

```
backend/
├── cmd/server/           # Entry point
├── internal/
│   ├── auth/             # Identity & Access module
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   ├── ports.go
│   │   └── module.go
│   ├── clinical/         # Clinical module
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   ├── ports.go
│   │   └── module.go
│   ├── imaging/          # Imaging module
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   ├── ports.go
│   │   └── module.go
│   ├── audit/            # Audit module
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   ├── ports.go
│   │   └── module.go
│   └── shared/           # Shared Kernel
│       ├── middleware/    # Tenant, auth, logging
│       ├── domain/       # Shared value objects
│       └── infrastructure/ # Shared DB connection, S3 client
├── migrations/           # SQL migrations
├── sqlc.yaml            # sqlc configuration
└── go.mod
```

## Consequences

### Positive
- Clear module boundaries aligned with bounded contexts
- Each module is independently testable
- Future decomposition: replace in-process interfaces with network calls
- Shared kernel reduces duplication (middleware, tenant context)
- Single deployment unit (simpler than microservices for PoC)
- Database tables owned by specific modules

### Negative
- Requires discipline to maintain module boundaries
- In-process communication can lead to tight coupling if not careful
- Shared database requires clear table ownership

## Module Communication

```
Module A ──calls──▶ Module B (via ports.go interface)
```

- Modules expose public interfaces in `ports.go`
- Other modules depend on interfaces, not implementations
- No direct database access across modules
- No circular dependencies between modules

## Future Decomposition

When the service requires independent scaling:

1. Extract module into separate service
2. Replace in-process interface with HTTP/gRPC call
3. Separate database per service
4. Add API gateway for routing

The modular monolith makes this a mechanical transformation, not a rewrite.

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Traditional monolith | No clear boundaries, harder to decompose later |
| Microservices | Overkill for PoC, adds operational complexity |
| Hexagonal architecture | Similar but less focused on module boundaries |

## References

- [Modular Monolith by Spotify](https://engineering.atspotify.com/2020/04/when-should-i-build-a-monolith/)
- [Modular Monolith by Shopify](https://shopify.engineering/building-a-modular-monolith)
