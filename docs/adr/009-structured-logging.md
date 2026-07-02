# ADR-009: Structured Logging

## Status

Accepted

## Context

MedVault needs observability through logging. Logs should be structured, searchable, and support audit requirements. The logging solution should align with the stdlib-first philosophy.

## Decision

Use `log/slog` (Go stdlib) for structured JSON logging with CloudWatch Logs integration.

## Consequences

### Positive
- `log/slog` is stdlib (Go 1.21+), no external dependency
- Structured JSON output by default
- Leveled logging (ERROR, WARN, INFO, DEBUG)
- Easy to query with CloudWatch Insights
- Consistent log format across components
- Supports audit logging requirements
- Integration with CloudWatch Logs

### Negative
- `log/slog` is newer, less battle-tested than zerolog/zap
- No built-in log rotation (handled by CloudWatch agent)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| zerolog | `log/slog` covers our needs without external dependency |
| zap | `log/slog` covers our needs without external dependency |
| Plain text logs | Harder to query and analyze |
| ELK Stack | Overkill for PoC, higher operational cost |
| Datadog | Third-party dependency, cost |

## References

- [log/slog Documentation](https://pkg.go.dev/log/slog)
- [CloudWatch Logs Documentation](https://docs.aws.amazon.com/cloudwatch/)
