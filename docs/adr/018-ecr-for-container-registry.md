# ADR-018: ECR for Container Registry

## Status

Accepted

## Context

MedVault's backend runs on ECS Fargate and requires a container image registry to store Docker images. The registry must integrate with ECS, support CI/CD pipelines, and align with the project's managed-services-first philosophy.

## Decision

Use Amazon Elastic Container Registry (ECR) for storing Docker images.

## Consequences

### Positive
- Native AWS integration (ECS pulls images without authentication configuration)
- IAM-based access control (same permission model as other AWS resources)
- Image scanning for vulnerabilities (optional, integrates with CI/CD)
- Image versioning and immutability support
- No additional infrastructure to manage
- Pay-per-storage pricing (cost-effective for PoC)

### Negative
- AWS-only (no multi-cloud portability)
- ECR repository must be provisioned via Terraform

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Docker Hub | External dependency, rate limits on free tier, no native ECS integration |
| GitHub Container Registry (ghcr.io) | External dependency, requires separate credentials management |
| Self-hosted registry | Operational overhead, violates managed-services-first principle |

## References

- [ECR Documentation](https://docs.aws.amazon.com/ecr/)
- [ECS + ECR Integration](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_ECS.html)
- [ADR-004: ECS Fargate for Compute](004-ecs-fargate-for-compute.md)
