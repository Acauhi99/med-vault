# Infrastructure

Infrastructure architecture for MedVault. Infrastructure is a first-class citizen in this project — treated with the same engineering discipline as application code.

> **Disclaimer:** MedVault is a Proof of Concept. Infrastructure demonstrates architectural thinking, not production readiness.

---

## Philosophy

Infrastructure should be:

- **Simple** — avoid unnecessary abstractions
- **Secure** — security by default, not as an afterthought
- **Reproducible** — every resource managed via Terraform
- **Maintainable** — readable over clever
- **Explainable** — understandable before becoming reusable

Avoid premature modularization. Infrastructure should model platform capabilities rather than individual AWS resources.

---

## Goals

The infrastructure demonstrates:

- AWS architecture with managed services
- Infrastructure as Code (Terraform)
- HIPAA-aware cloud design
- Security by default
- Cost awareness
- Reproducibility
- Cloud-native architecture

This project is intentionally designed as a production-inspired Proof of Concept. Although only one AWS environment will be deployed, the repository clearly demonstrates how it would evolve into a multi-environment infrastructure.

---

## Terraform Philosophy

Terraform is the single source of truth for all AWS resources. No infrastructure changes should exist outside Terraform unless explicitly documented.

- Declarative infrastructure over clever abstractions
- Readability preferred over reducing duplicated code
- Infrastructure should remain approachable for engineers unfamiliar with the project

---

## Repository Structure

```
infrastructure/
└── terraform/
    ├── modules/
    │   ├── network/
    │   ├── application/
    │   ├── database/
    │   ├── storage/
    │   ├── security/
    │   └── observability/
    │
    ├── environments/
    │   └── production/
    │       ├── backend.tf
    │       ├── providers.tf
    │       ├── versions.tf
    │       ├── variables.tf
    │       ├── terraform.tfvars
    │       ├── main.tf
    │       └── outputs.tf
    │
    └── README.md
```

### Directories

| Directory | Purpose |
|-----------|---------|
| `modules/` | Reusable infrastructure modules representing platform capabilities |
| `modules/network/` | VPC, subnets, route tables, internet gateway |
| `modules/application/` | ECS cluster, service, ALB, task definition |
| `modules/database/` | RDS PostgreSQL, subnet groups, parameter groups |
| `modules/storage/` | S3 buckets, lifecycle rules, policies |
| `modules/security/` | IAM, KMS, Secrets Manager |
| `modules/observability/` | CloudTrail, CloudWatch, AWS Config, GuardDuty |
| `environments/` | Environment-specific configurations |
| `environments/production/` | Production environment composition |

---

## Module Philosophy

Modules represent platform capabilities, not individual AWS resources.

### Anti-patterns

Avoid modules that simply encapsulate single AWS resources:

- `security_group`
- `iam_role`
- `kms_key`
- `cloudwatch_log_group`

### Correct approach

Modules should represent complete capabilities with cohesive resources:

**network**
- VPC
- Public Subnets
- Private Subnets
- Route Tables
- Internet Gateway

**application**
- ECS Cluster
- ECS Service
- Task Definition
- Application Load Balancer
- Target Groups
- Security Groups
- IAM Roles
- CloudWatch Log Groups
- ECR Repository (Docker images)

**database**
- PostgreSQL (Amazon RDS)
- DB Subnet Group
- Parameter Groups
- Security Groups
- Encryption

**storage**
- S3 Buckets
- Lifecycle Rules
- Bucket Policies
- Encryption
- Versioning

**security**
- IAM
- KMS
- Secrets Manager
- IAM Policies
- IAM Roles

**observability**
- CloudTrail
- CloudWatch
- AWS Config
- GuardDuty
- Security Hub

---

## Module Design Guidelines

A module should have:

- A single responsibility
- Cohesive resources
- Explicit inputs
- Explicit outputs

Rules:

- Avoid modules with dozens of variables
- Prefer a small number of meaningful inputs
- Avoid generic modules attempting to solve every possible scenario

---

## Root Environment

The Production environment should remain intentionally small. The environment primarily composes modules.

The root module should avoid directly managing AWS resources whenever possible.

**Responsibilities:**

- Provider configuration
- Backend configuration
- Module composition
- Shared variables
- Outputs

Business infrastructure belongs inside modules.

---

## Environment Strategy

The project initially deploys only Production. The architecture is designed to evolve.

**Future environments:**

- development
- staging
- production

Every environment consumes the same reusable modules. The only differences between environments should be:

- Variable values
- Sizing
- Scaling
- Networking
- Secrets

Never duplicate infrastructure code between environments.

---

## Terraform State

| Concern | Decision |
|---------|----------|
| Backend | Amazon S3 |
| Versioning | Enabled |
| Encryption | Server-side encryption (SSE-S3 or SSE-KMS) |
| Locking | DynamoDB table (for real production) |

**Why not local state:** Local Terraform state is not production-ready. It has no concurrency protection, no backup, no encryption, and no versioning. Always use remote state for any environment that matters.

---

## Security Principles

Infrastructure enforces security by default.

### Least Privilege

Every IAM role grants only the permissions required for its function. No wildcard permissions in production.

### Private Networking

- Backend runs in private subnets (no public IP)
- RDS in private subnets
- ALB in public subnets only
- Security groups restrict inter-service communication

### Encryption at Rest

| Resource | Method |
|----------|--------|
| RDS PostgreSQL | AES-256 (AWS KMS) |
| S3 medical images | AES-256 (SSE-S3) |
| S3 audit logs | AES-256 (SSE-S3) |
| Secrets | AES-256 (Secrets Manager) |

### Encryption in Transit

| Connection | Protocol |
|------------|----------|
| Client → CloudFront | TLS 1.2+ |
| CloudFront → ALB | TLS 1.2+ |
| ALB → ECS | HTTP (internal VPC) |
| ECS → RDS | TLS 1.2+ |
| ECS → S3 | HTTPS |

### Secrets Management

- Database credentials: AWS Secrets Manager
- JWT signing key: AWS Secrets Manager
- No hardcoded secrets in code or environment variables
- Automatic rotation for database credentials

### Audit Logging

- CloudTrail for AWS API audit trail
- VPC Flow Logs for network audit
- Application audit logs via CloudWatch

### Immutable Infrastructure

Infrastructure changes replace rather than modify. New deployments create new resources; old resources are destroyed.

---

## AWS Shared Responsibility Model

This project operates under the AWS Shared Responsibility Model. See [SECURITY.md](SECURITY.md#aws-shared-responsibility-model) for the full breakdown of AWS, Infrastructure, Application, and Developer responsibilities.

---

## Cost Philosophy

The project intentionally balances production realism with cost awareness.

- When multiple AWS services could satisfy the same requirement, prefer the simpler and more cost-effective option
- Avoid introducing unnecessary AWS services solely for completeness
- The infrastructure should be explainable during a technical interview

**Cost-conscious choices:**

- ECS Fargate over EKS (no cluster management cost)
- RDS PostgreSQL over Aurora (simpler for PoC)
- S3 standard over S3 Intelligent-Tiering (predictable access patterns)
- CloudWatch over Datadog (AWS-native, no additional cost)

---

## Future Evolution

The following are documented as future improvements, not current implementation requirements.

### CI/CD

See [CI_CD_STRATEGY.md](CI_CD_STRATEGY.md) for the full pipeline architecture. The Infrastructure pipeline owns Terraform only — it provisions the platform but never deploys application code.

### Multi-Account Strategy

- Separate AWS accounts for development, staging, production
- AWS Organizations for account management
- Cross-account IAM roles

### Multi-Region

- Multi-region deployment for disaster recovery
- Route 53 failover routing
- Cross-region replication for S3 and RDS

### Auto Scaling

- ECS Service auto-scaling based on CPU/memory
- RDS read replicas for read-heavy workloads
- S3 Intelligent-Tiering for cost optimization

### WAF Hardening

- AWS WAF rules for common attack patterns
- Rate limiting
- Geo-restriction (if needed)

### Disaster Recovery

- RDS automated backups with point-in-time recovery
- S3 cross-region replication
- Regular disaster recovery testing

### Secrets Rotation

- Automatic rotation for all secrets
- AWS Lambda functions for rotation
- Integration with application health checks

### Policy as Code

- AWS Config rules for compliance
- Service Control Policies (SCPs)
- IAM Access Analyzer

---

## Architecture Decision Records

Major infrastructure decisions should be recorded as ADRs.

| ADR | Decision |
|-----|----------|
| [ADR-003](adr/003-terraform-for-infrastructure.md) | Terraform for IaC |
| [ADR-004](adr/004-ecs-fargate-for-compute.md) | ECS Fargate for compute |
| [ADR-005](adr/005-postgresql-as-database.md) | PostgreSQL as database |
| [ADR-008](adr/008-s3-for-medical-image-storage.md) | S3 for medical image storage |

**Future ADRs to create:**

- Why ECS instead of EKS
- Why RDS instead of Aurora
- Why CloudFront + S3
- Why Fargate instead of EC2
- Why Shared PostgreSQL Multi-Tenancy
- Why Managed AWS Services over self-hosted

Each ADR should explain: Context, Decision, Alternatives, Consequences.

---

## AI Development Guidelines

This repository is designed for AI-assisted development.

- Infrastructure documentation prioritizes consistency and predictability
- Terraform modules follow repeatable conventions
- Architectural consistency is more important than aggressive abstraction
- Every infrastructure decision should be easy to explain to another engineer
- The documentation teaches future contributors how to think about the infrastructure rather than simply describing the current implementation
