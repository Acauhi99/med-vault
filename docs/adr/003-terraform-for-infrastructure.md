# ADR-003: Terraform for Infrastructure as Code

## Status

Accepted

## Context

MedVault requires all infrastructure to be managed as code. The tool should support AWS resources, have a declarative approach, and align with the project's reproducibility principle.

## Decision

Use Terraform as the Infrastructure as Code tool.

## Consequences

### Positive
- Declarative configuration
- AWS provider is mature and well-documented
- State management for resource tracking
- Module system for reusable components
- Plan/apply workflow for safe changes
- Large community and ecosystem

### Negative
- State file management requires care
- HCL learning curve
- Plan output can be verbose

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| AWS CloudFormation | AWS-only, less flexible syntax |
| Pulumi | Requires programming language knowledge |
| AWS CDK | Less mature than Terraform for some resources |

## References

- [Terraform Documentation](https://developer.hashicorp.com/terraform/docs)
- [AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
