# ADR-001: Go as Backend Language

## Status

Accepted

## Context

MedVault needs a backend language that demonstrates secure, cloud-native development. The language should be suitable for REST API development, have strong standard library support, and align with the project's engineering principles. We prioritize the Go stdlib over third-party libraries when stdlib provides a production-ready solution.

## Decision

Use Go (Golang) as the backend programming language with a stdlib-first philosophy.

**Stack:**

| Layer | Tool | Reason |
|-------|------|--------|
| HTTP | `net/http` | Stdlib, production-ready, no external dependency |
| Router | `http.ServeMux` | Stdlib, sufficient for modular monolith routing |
| Config | `envconfig` | Minimal, 12-factor app aligned, no reflection-heavy config libs |
| Database | `pgx` | Best PostgreSQL driver, connection pool, binary protocol |
| Queries | `sqlc` | Type-safe SQL code generation, no ORM overhead |
| Migrations | `golang-migrate` | SQL-file based, versioned, works with pgx |
| Logging | `log/slog` | Stdlib structured logging (Go 1.21+), leveled, JSON output |
| Tests | `testing` + `httptest` | Stdlib, sufficient for HTTP handler and integration tests |
| Container | Docker | Minimal images, stdlib compiles to static binary |
| Deploy | ECS Fargate | Serverless containers, no node management |

## Consequences

### Positive
- Strong standard library covers HTTP, crypto, encoding, logging, testing
- Compiled binary with no runtime dependencies
- Excellent concurrency model (goroutines, channels)
- Static typing catches errors at compile time
- Fast compilation and startup time
- Small Docker images (scratch or distroless)
- Minimal third-party dependencies = smaller attack surface
- `log/slog` provides structured logging without external packages
- `net/http` + `http.ServeMux` sufficient for modular monolith routing

### Negative
- More verbose than some alternatives
- Requires manual error handling
- No built-in ORM (but this aligns with Clean Architecture and sqlc approach)
- `http.ServeMux` lacks middleware chaining (solved with wrapper pattern)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Python | Slower runtime, dynamic typing, larger container images |
| Node.js | Single-threaded event loop, dynamic typing |
| Java | Heavier runtime, more boilerplate, larger images |
| Rust | Steeper learning curve, slower development speed for PoC |
| chi/gorilla/mux | Unnecessary when `http.ServeMux` covers our routing needs |
| zerolog/zap | `log/slog` now covers structured logging in stdlib |
| gorm/sqlx | sqlc provides type safety without runtime reflection overhead |

## References

- [Go Documentation](https://go.dev/doc/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx)
- [sqlc Documentation](https://sqlc.dev/)
- [log/slog Documentation](https://pkg.go.dev/log/slog)
