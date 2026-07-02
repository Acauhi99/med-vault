# ADR-014: golang-migrate for Database Migrations

## Status

Accepted

## Context

MedVault needs a tool to manage PostgreSQL schema migrations. Migrations must be versioned, repeatable, and work with the pgx driver. The tool should support SQL-file based migrations that can be embedded in the Go binary or run from the CLI.

## Decision

Use `golang-migrate/migrate` for database schema management.

## Consequences

### Positive
- SQL-file based migrations (plain SQL, no DSL)
- Versioned with numeric timestamps
- Up/down migration support
- Works with pgx via the `pgx` database driver
- Can embed migrations in Go binary or run from CLI
- Integrates with `sqlc` (migrations define the schema sqlc generates from)
- Simple, focused tool — does one thing well

### Negative
- CLI tool requires separate installation
- No automatic migration generation (manual SQL writing)
- Less feature-rich than some alternatives (no rollback protection)

## Usage

```
migrations/
├── 000001_create_tenants.up.sql
├── 000001_create_tenants.down.sql
├── 000002_create_users.up.sql
├── 000002_create_users.down.sql
└── ...
```

```bash
# Apply migrations
migrate -path migrations -database "postgres://..." up

# Rollback
migrate -path migrations -database "postgres://..." down 1
```

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| goose | Similar but golang-migrate has better pgx integration |
| atlas | Newer, less mature ecosystem |
| dbmate | Ruby-inspired, less Go-native |
| Custom migration code | Reimplements what golang-migrate does well |

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
