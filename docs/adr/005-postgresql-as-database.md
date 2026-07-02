# ADR-005: PostgreSQL as Database

## Status

Accepted

## Context

MedVault needs a relational database for structured data storage. The database should support multi-tenancy, have strong consistency, and align with the project's security requirements.

## Decision

Use PostgreSQL (Amazon RDS) as the primary database.

## Consequences

### Positive
- ACID compliance
- Strong relational model
- Row-Level Security (RLS) for tenant isolation
- JSON support for flexible data
- Mature ecosystem and tooling
- AWS RDS managed service

### Negative
- Requires schema design
- Connection pooling needed for scale

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| MySQL | Less feature-rich than PostgreSQL |
| DynamoDB | NoSQL, less suitable for relational data |
| MongoDB | Document model less suitable for this domain |

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Amazon RDS Documentation](https://docs.aws.amazon.com/rds/)
