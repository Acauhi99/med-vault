# CI/CD Strategy

> **This is the source of truth for CI/CD philosophy, pipeline architecture, and deployment boundaries.** Other documents reference this document rather than duplicating its content.

This document defines the CI/CD architecture for MedVault. It covers philosophy, pipeline responsibilities, deployment lifecycle, and operational boundaries.

> **Disclaimer:** MedVault is a Proof of Concept. This document defines the target CI/CD architecture. Implementation occurs in Phase 9 (see [ROADMAP.md](../ROADMAP.md)).

---

## Philosophy

The CI/CD architecture mirrors the application architecture.

Each deployment pipeline is responsible only for the lifecycle of the artifact it produces. Pipelines never own responsibilities belonging to another layer of the system.

**Objectives:**

- Separation of concerns
- Maintainability
- Scalability
- Deployment safety
- Rollback simplicity
- Reproducibility

Every pipeline has a single, well-defined responsibility.

---

## Pipeline Architecture

MedVault has three independent deployment pipelines:

```
Infrastructure
     ↓
  Backend
     ↓
  Frontend
```

Each pipeline owns its own lifecycle. Pipelines communicate only through:

- Deployed infrastructure (Terraform outputs consumed by application pipelines)
- Published artifacts (Docker images in ECR, static assets in S3)

No pipeline directly deploys another project's code.

---

## Infrastructure Pipeline

The Infrastructure pipeline owns the AWS platform.

### Responsibilities

- Terraform Init
- Terraform Validate
- Terraform Plan
- Manual Approval (production)
- Terraform Apply

### Resources Provisioned

| Category | Resources |
|----------|-----------|
| Networking | VPC, public/private subnets, NAT gateway, route tables, internet gateway |
| Compute | ECS Cluster, ECS Service, ALB, target groups |
| Database | RDS PostgreSQL, subnet groups, parameter groups |
| Storage | S3 buckets, lifecycle rules, bucket policies |
| Security | IAM roles, policies, KMS keys, Secrets Manager |
| Observability | CloudWatch, CloudTrail, VPC Flow Logs, AWS Config, Security Hub, GuardDuty |
| CDN | CloudFront distribution |

### Boundaries

The Infrastructure pipeline must NOT:

- Deploy application code
- Build Docker images
- Upload frontend assets
- Execute database migrations

Infrastructure provisions the platform. Applications consume the platform.

---

## Backend Pipeline

The Backend pipeline owns the Go application runtime.

### Execution Flow

```
Format → Lint → Unit Tests → Build → Docker Image → Push to ECR → Database Migration → ECS Deploy → Health Check
```

### Responsibilities

| Step | Tool | Purpose |
|------|------|---------|
| Format | `gofumpt` | Consistent code formatting |
| Static Analysis | `golangci-lint` | Code quality and correctness |
| Unit Tests | `go test` | Behavior validation |
| Build | `go build` | Binary compilation |
| Docker Image | `docker build` | Container packaging |
| Push | `docker push` → Amazon ECR | Artifact publication |
| Migration | `golang-migrate` | Schema deployment |
| Deploy | ECS Service Update | Runtime deployment |
| Health Check | ALB target group | Availability validation |

### Artifact

The Backend pipeline produces a Docker image stored in Amazon ECR.

The image is the sole artifact. It is immutable, versioned, and independently deployable.

> **Docker image strategy:** See [ADR-019: Docker Image Strategy](adr/019-docker-image-strategy.md) for multi-stage build, distroless runtime, security constraints, and layering strategy.

### Boundaries

The Backend pipeline must NOT:

- Modify Terraform infrastructure
- Deploy frontend assets
- Manage CloudFront

---

## Frontend Pipeline

The Frontend pipeline owns the static web application.

### Execution Flow

```
Format → Lint → Type Check → Unit Tests → Integration Tests → Next.js Build → Static Export → Upload to S3 → CloudFront Invalidation
```

### Responsibilities

| Step | Tool | Purpose |
|------|------|---------|
| Format | Biome | Consistent code formatting |
| Lint | Biome | Code quality |
| Type Check | `tsc` | TypeScript correctness |
| Unit Tests | Vitest | Component behavior |
| Integration Tests | Vitest + MSW | API contract validation |
| Build | `next build` | Application compilation |
| Static Export | `next build` with `output: 'export'` | HTML/CSS/JS generation |
| Upload | `aws s3 sync` → S3 | Asset publication |
| Invalidation | `aws cloudfront create-invalidation` | Cache refresh |

### Artifact

The Frontend pipeline produces static HTML, CSS, and JavaScript files stored in S3.

### Boundaries

The Frontend pipeline must NOT:

- Deploy backend services
- Execute Terraform
- Manage databases
- Execute migrations

---

## Database Migration Strategy

Database schema migrations belong to the Backend deployment lifecycle.

### Ownership

| Concern | Owner |
|---------|-------|
| Database provisioning | Infrastructure (Terraform) |
| Schema migrations | Backend pipeline (`golang-migrate`) |
| Application code | Backend pipeline |

Terraform provisions the database instance. The Backend pipeline owns the schema.

### Execution Order

```
Backend Pipeline:
  1. Build application
  2. Run migrations        ← schema changes
  3. Deploy new version    ← application uses new schema
  4. Health check
```

Migrations execute **before** the new application version becomes active.

### Safety Guarantees

- Every migration has a corresponding `.down.sql` for rollback
- Migrations run as a dedicated step (ECS Run Task or CI job), not at application startup
- Multiple application containers never race on the same migration
- Rollback is explicit: `migrate down N`, then rollback application

### Why Separate from Application Startup

- Avoids race conditions in multi-container deployments
- Application startup stays clean (no schema changes during boot)
- Rollback is decoupled (rollback app without touching schema, or vice versa)
- CI/CD controls the migration lifecycle, not the application

---

## Environment Variables

### Infrastructure Variables

Owned by Terraform.

| Variable | Example |
|----------|---------|
| AWS Region | `us-east-1` |
| Project Name | `medvault` |
| Network Configuration | VPC CIDR, subnet ranges |
| Domain Name | `medvault.example.com` |

These are Terraform inputs. They define the platform.

### Backend Variables

Injected through ECS task definition via AWS Secrets Manager or AWS Systems Manager Parameter Store.

| Variable | Source |
|----------|--------|
| Database connection string | Secrets Manager |
| JWT signing key | Secrets Manager |
| AWS Region | Parameter Store |
| S3 bucket names | Parameter Store |
| KMS key ARNs | Parameter Store |
| Application config | Environment-specific |

The deployment pipeline references secrets. It never exposes secret values.

### Frontend Variables

Only public runtime configuration.

| Variable | Purpose |
|----------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API endpoint |

Because the frontend is statically exported, these variables become part of the generated application at build time.

**Sensitive information must never be embedded in the frontend bundle.**

---

## Secrets Management

### Storage

Secrets are stored in AWS:

- **AWS Secrets Manager** — database credentials, JWT signing keys
- **AWS Systems Manager Parameter Store** — non-secret configuration

### Authentication

GitHub Actions authenticates using **GitHub OIDC** (OpenID Connect).

| Practice | Status |
|----------|--------|
| Long-lived AWS credentials | Avoided |
| Static IAM users | Avoided |
| Temporary credentials (OIDC) | ✅ Configured |
| Secret values in logs | Never |

**OIDC Provider:** `arn:aws:iam::836734448013:oidc-provider/token.actions.githubusercontent.com`
**IAM Role:** `medvault-github-actions` (scoped to `repo:Acauhi99/med-vault:*`)

---

## Deployment Order

### Initial Deployment

```
Infrastructure → Backend → Frontend
```

The platform must exist before applications can deploy.

### Steady-State Independence

After the initial deployment, every pipeline is independently deployable:

- Updating the frontend does not require redeploying the backend
- Updating the backend does not require redeploying the frontend
- Infrastructure changes occur independently whenever possible

---

## Workflow Triggers

Pipelines run only when their files change.

| Pipeline | Trigger Path |
|----------|-------------|
| Infrastructure | `infrastructure/` |
| Backend | `backend/`, `spec/openapi.yaml` |
| Frontend | `frontend/`, `spec/openapi.yaml` |

The OpenAPI spec triggers both Backend and Frontend pipelines because it is the shared contract.

Unnecessary pipeline executions are avoided.

---

## Concurrency Strategy

Each pipeline owns its own concurrency group:

| Pipeline | Concurrency Group |
|----------|-------------------|
| Infrastructure | `production-infrastructure` |
| Backend | `production-backend` |
| Frontend | `production-frontend` |

**Rules:**

- A newer deployment cancels older pending deployments of the same component
- Deployments of different components remain independent

---

## Rollback Strategy

Each layer supports independent rollback.

| Layer | Rollback Method |
|-------|-----------------|
| Infrastructure | `terraform apply` with previous state |
| Backend | Redeploy previous container image (ECS task definition rollback) |
| Frontend | Redeploy previous static build (S3 version restore) |

Rollback is explicit and traceable.

---

## Deployment Validation

Deployment is not considered successful until validation completes.

### Infrastructure Validation

- Terraform state consistency
- Resource health checks

### Backend Validation

- Application health check endpoint (`/health`)
- Database connectivity
- API availability (smoke test)

### Frontend Validation

- Static asset availability
- CloudFront propagation
- Application accessibility (HTTP 200)

---

## Future Evolution

The following are documented as future improvements, not current implementation requirements.

| Capability | Purpose |
|------------|---------|
| Blue/Green Deployments | Zero-downtime backend releases |
| Canary Releases | Gradual traffic shifting |
| Progressive Delivery | Risk-based deployment |
| Feature Flags | Decouple deploy from release |
| Multi-Region Deployments | Disaster recovery, latency |
| Multi-Account AWS Strategy | Environment isolation |
| Automated Rollback | Health-check-triggered rollback |
| Disaster Recovery | Cross-region failover |
| Container Vulnerability Scanning | Image security |
| Supply Chain Security | Dependency verification |
| SBOM Generation | Artifact transparency |
| Release Automation | Tag-based releases |
| Preview Environments | PR-based previews |

---

## AI Development Guidelines

AI agents should reason about deployments using pipeline ownership:

- **Infrastructure** owns infrastructure
- **Backend** owns backend
- **Frontend** owns frontend

Agents should avoid introducing cross-pipeline responsibilities.

When modifying deployment-related code:

- Identify which pipeline owns the change
- Do not mix concerns across pipelines
- Respect the boundaries defined in this document

The deployment architecture should remain modular, predictable, and easy to evolve.
