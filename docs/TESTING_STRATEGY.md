# Testing Strategy

> **This is the source of truth for testing philosophy, stack, and patterns.** Other documents (ARCHITECTURE.md, ADRs, AGENTS.md) reference this document rather than duplicating its content.

## Objective

This project adopts a pragmatic testing strategy focused on delivering confidence with the lowest possible maintenance cost.

The goal is **not to maximize the number of tests or achieve 100% code coverage**, but to ensure that every business capability is validated through meaningful, maintainable, and deterministic tests.

Testing should validate the observable behavior of the system rather than its internal implementation details.

This strategy applies equally to both the Go backend and the Next.js frontend.

---

## Guiding Principles

- Prioritize confidence over coverage percentage.
- Test behavior, not implementation.
- Keep tests deterministic, isolated, and easy to understand.
- Avoid redundant tests across different layers.
- Every test should provide unique value.
- Prefer fewer high-quality tests over many low-value tests.
- A failing test should clearly indicate what behavior was broken.

---

## Testing Pyramid

This project intentionally adopts a simplified testing pyramid composed of only two layers:

```
Integration Tests
        ▲
Unit Tests
```

There are **no End-to-End (E2E) tests**.

Instead, confidence is achieved through a combination of:

- Unit Tests
- Integration Tests

This approach provides fast feedback, lower maintenance costs, and excellent reliability while avoiding the complexity and flakiness commonly associated with E2E testing.

---

## Development Workflow

Every new feature should follow this workflow:

1. Define or update the OpenAPI contract.
2. Implement the business logic.
3. Create unit tests for the business rules.
4. Create integration tests covering the complete application flow.
5. Validate that all expected behaviors are covered.
6. Merge only after both unit and integration tests pass.

Testing is considered part of the implementation, not an afterthought.

---

## Unit Tests

Unit tests validate isolated business logic.

They should focus on:

- Business rules
- Validations
- Domain logic
- Data transformations
- Decision making

Unit tests should execute quickly and should not depend on external infrastructure.

---

## Integration Tests

Integration tests validate the complete application flow.

They should exercise the full request lifecycle, ensuring that all layers work together correctly.

For the backend, this includes the complete HTTP request flow through handlers, services, repositories, database access, and persistence.

For the frontend, integration tests should validate the interaction between components, application logic, generated API clients, and mocked or controlled backend responses where appropriate.

The objective is to verify that independently tested components behave correctly when combined.

---

## Coverage Strategy

Each business capability should be covered by a small but meaningful set of tests.

As a general guideline, every feature should include:

- One Happy Path.
- Up to three Critical Edge Cases.

The selected edge cases should represent the highest-risk scenarios capable of breaking the expected behavior.

Examples include:

- Invalid input.
- Authorization failures.
- Business rule violations.
- Resource conflicts.
- Infrastructure failures.

The exact scenarios depend on the feature being implemented.

---

## Avoid Redundancy

Tests should not repeat validations already guaranteed by another layer.

Each test should exist for a specific purpose.

Avoid creating multiple tests that verify the same behavior through different paths without adding additional confidence.

A smaller, cohesive test suite is preferred over a large suite containing duplicated assertions.

---

## Behavior Over Implementation

Tests should validate externally observable behavior.

They should remain stable even when the internal implementation changes.

Refactoring internal code without changing behavior should not require rewriting the test suite.

This encourages better software design and reduces long-term maintenance costs.

---

## Quality Goals

Success is not measured by code coverage alone.

Instead, quality is measured by:

- Confidence when refactoring.
- Reliability of critical business flows.
- Fast execution.
- Low maintenance cost.
- Clear failure diagnostics.
- Deterministic and reproducible execution.

Coverage percentage is considered a secondary metric and should never drive test creation by itself.

---

## AI-Assisted Development

Because this project is heavily AI-assisted, every implementation should include its corresponding tests as part of the same development cycle.

AI should generate:

- Business implementation.
- Unit tests.
- Integration tests.

Developers are responsible for reviewing the generated tests to ensure they validate meaningful behaviors instead of simply increasing coverage numbers.

---

## Backend Testing Stack

| Tool | Purpose |
|------|---------|
| `testing` | Test runner (stdlib) |
| `testify/assert` | Assertions |
| `net/http/httptest` | HTTP handler testing |
| `go-cmp` | Struct comparison |
| `testcontainers-go` | Real PostgreSQL in integration tests |
| `go test -cover` | Coverage reporting |
| `testing.B` | Benchmarks |
| `testing.F` | Fuzz testing |

**Unit tests:** Domain logic, value objects, aggregates — isolated, fast, no I/O.

**Integration tests:** Repositories, HTTP handlers — real PostgreSQL via testcontainers, real HTTP via httptest.

### Testing Domain Events and Outbox

Domain events are delivered via the Transactional Outbox pattern (see [ADR-017](adr/017-transactional-outbox.md)). Testing should cover:

| What to Test | How |
|-------------|-----|
| Aggregate emits correct events | Unit test: call command, inspect events on aggregate |
| Event persisted in outbox | Integration test: call command, query `domain_outbox` table |
| Projection handler processes event | Unit test: call handler with event, assert read model updated |
| Projection handler is idempotent | Unit test: call handler twice with same event, assert no duplicate |
| Outbox poller dispatches events | Integration test: insert outbox row, run poller, assert handler called |
| Failed event marked correctly | Integration test: handler returns error, assert `attempts` incremented |

**Key rules:**
- Projection handlers must be tested for idempotency (same event twice = same result)
- Outbox integration tests use real PostgreSQL via testcontainers
- Event payload serialization/deserialization tested round-trip

---

## Frontend Testing Stack

| Tool | Purpose |
|------|---------|
| Vitest | Test runner |
| `@testing-library/react` | Component testing utilities |
| `@testing-library/user-event` | User interaction simulation |
| MSW (Mock Service Worker) | API mocking |
| `@vitest/coverage-v8` | Coverage reporting |

**Unit tests:** Pure functions, Zod schemas, data transformations — isolated, fast.

**Integration tests:** Component rendering, user interactions, API calls via MSW — validate component + hook + service integration.

---

## Philosophy

The testing strategy follows the same architectural principles adopted throughout the project:

- Explicit over implicit.
- Simplicity over unnecessary complexity.
- Confidence over quantity.
- Behavior over implementation.
- Maintainability over coverage metrics.
- High-value tests over exhaustive testing.
- Cohesive test suites with minimal redundancy.

The objective is to build a test suite that remains fast, reliable, easy to evolve, and capable of providing high confidence throughout the lifetime of the project.
