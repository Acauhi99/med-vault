# ADR-010: DDD with CQRS Architecture

## Status

Accepted

## Context

MedVault needs an architecture that demonstrates clean separation of concerns, testability, and alignment with the domain. The architecture should be easy to navigate while following DDD tactical and strategic patterns with CQRS.

## Decision

Use Domain-Driven Design (DDD) with CQRS (Command Query Responsibility Segregation).

## Consequences

### Positive
- Explicit domain model (Aggregates, Entities, Value Objects)
- Clear separation of reads (queries) and writes (commands)
- Domain Events for decoupling
- Testable at every layer
- Easy to navigate (bounded contexts map to business capabilities)
- Ubiquitous Language aligns code with domain

### Negative
- More concepts to learn (but well-documented in DOMAIN.md)
- Requires discipline to maintain boundaries

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Clean Architecture | Less explicit domain model |
| Hexagonal Architecture | Similar but less focused on domain |
| CRUD-only | Too simple, doesn't demonstrate architectural thinking |

## References

- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
