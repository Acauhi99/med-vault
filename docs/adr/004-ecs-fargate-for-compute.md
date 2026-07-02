# ADR-004: ECS Fargate for Compute

## Status

Accepted

## Context

MedVault needs a compute platform for the Go backend. The platform should be serverless (no server management), scalable, and cost-effective for a PoC.

## Decision

Use Amazon ECS with Fargate for compute.

## Consequences

### Positive
- No server management (serverless containers)
- Pay-per-use pricing
- Auto-scaling support
- Integration with ALB, IAM, CloudWatch
- Small Docker images (Go binary)
- Production-like architecture for PoC

### Negative
- Cold start latency (minimal for PoC)
- Less control than EC2 (but acceptable)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| EC2 | Requires server management, more operational overhead |
| Lambda | Not suitable for long-running REST API |
| EKS | Overkill for PoC, requires Kubernetes knowledge |
| App Runner | Less mature than ECS Fargate |

## References

- [ECS Fargate Documentation](https://docs.aws.amazon.com/ecs/latest/developerguide/AWS_Fargate.html)
