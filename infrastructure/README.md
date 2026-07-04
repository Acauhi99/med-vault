# Infrastructure

Terraform infrastructure for MedVault.

## Official Docs

- [Context](../CONTEXT.md)
- [Architecture](../docs/ARCHITECTURE.md)
- [Security](../docs/SECURITY.md)
- [Infrastructure](../docs/INFRASTRUCTURE.md)
- [CI/CD Strategy](../docs/CI_CD_STRATEGY.md)
- [Quality Gates](../docs/QUALITY_GATES.md)
- [Infrastructure ADRs](../docs/adr/)

## Commands

```bash
terraform -chdir=terraform/environments/production init
terraform -chdir=terraform/environments/production validate
terraform -chdir=terraform/environments/production plan
terraform -chdir=terraform/environments/production apply
```

## Structure

- `terraform/modules/` reusable platform capabilities
- `terraform/environments/production/` root environment composition
- `terraform/environments/production/backend.tf` remote state backend
