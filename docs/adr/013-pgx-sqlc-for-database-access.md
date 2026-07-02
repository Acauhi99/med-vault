# ADR-013: pgx + sqlc for Database Access

## Status

Accepted

## Context

MedVault needs a database access layer for PostgreSQL. The layer should be type-safe, SQL-first, and avoid ORM overhead while maintaining developer productivity.

## Decision

Use pgx as the PostgreSQL driver and sqlc for generating type-safe Go code from SQL queries.

## pgx

**Purpose:** Pure Go PostgreSQL driver and client library.

**Features:**
- Pure Go implementation (no CGO)
- Connection pooling (pgxpool)
- Prepared statement support
- LISTEN/NOTIFY support
- COPY protocol support
- Binary protocol (faster than text)
- Transaction support

## sqlc

**Purpose:** Generate type-safe Go code from SQL queries.

**Workflow:**
1. Write SQL queries in `.sql` files
2. Define schema in `migrations/`
3. Run `sqlc generate` to produce Go code
4. Use generated code in repository implementations

**Example:**

```sql
-- query.sql
-- name: GetCaseByID :one
SELECT * FROM cases
WHERE id = $1 AND tenant_id = $2;

-- name: ListCasesByPatient :many
SELECT * FROM cases
WHERE patient_id = $1 AND tenant_id = $2
ORDER BY created_at DESC;
```

```go
// Generated code (sqlc)
func (q *Queries) GetCaseByID(ctx context.Context, arg GetCaseByIDParams) (Case, error) {
    // Type-safe, compile-time checked
}
```

## Consequences

### Positive
- Type-safe database access (compile-time checks)
- SQL-first approach (no ORM abstractions)
- Generated code is readable and debuggable
- No runtime reflection
- Excellent performance (pgx binary protocol)
- Easy to understand for developers who know SQL
- Generated code follows Go conventions

### Negative
- Requires sqlc CLI in build pipeline
- SQL queries must be maintained separately from Go code
- Schema changes require regeneration

## Project Structure

```
backend/
├── migrations/           # SQL migrations (schema)
├── queries/              # SQL query files (per module)
│   ├── auth.sql
│   ├── clinical.sql
│   ├── imaging.sql
│   └── audit.sql
├── sqlc.yaml            # sqlc configuration
└── internal/
    └── [module]/
        └── infrastructure/
            └── queries/  # Generated Go code
                └── queries.go
```

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| database/sql + pgx driver | No code generation, manual scanning |
| GORM | Too much magic, hard to optimize |
| sqlx | Less type safety than sqlc |
| sqlx + sqlx | More manual work than sqlc |

## References

- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [sqlc Documentation](https://sqlc.dev/)
- [sqlc GitHub](https://github.com/sqlc-dev/sqlc)
