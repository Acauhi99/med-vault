# MedVault

Cloud-native healthcare platform reference implementation. Demonstrates secure, multi-tenant architecture on AWS using Infrastructure as Code.

> **This is a Proof of Concept.** It demonstrates architectural thinking, not production readiness. Do not deploy with real patient data.

---

## Documentation

| Document | Purpose |
|----------|---------|
| [Architecture](docs/ARCHITECTURE.md) | System design, components, data flow, deployment |
| [Domain Model](docs/DOMAIN.md) | DDD bounded contexts, aggregates, CQRS mapping |
| [Requirements](docs/REQUIREMENTS.md) | Functional and non-functional requirements, API spec |
| [Security](docs/SECURITY.md) | Threat model, encryption, auth, audit, compliance |
| [Infrastructure](docs/INFRASTRUCTURE.md) | Terraform philosophy, modules, security, state, evolution |
| [Testing Strategy](docs/TESTING_STRATEGY.md) | Testing philosophy, principles, pyramid, coverage approach |
| [Quality Gates](docs/QUALITY_GATES.md) | Validation layers, tooling, pre-commit/pre-push, task runner |
| [Principles](docs/PROJECT_PRINCIPLES.md) | Engineering principles and decision criteria |
| [Context](docs/CONTEXT.md) | Project background, scope, and goals |
| [Roadmap](ROADMAP.md) | Phased delivery plan with status tracking |
| [Checklist](docs/CHECKLIST.md) | Acceptance criteria per phase |

### Architecture Decision Records

| ADR | Decision |
|-----|----------|
| [ADR-001](docs/adr/001-go-as-backend-language.md) | Go as backend language |
| [ADR-002](docs/adr/002-react-as-frontend-framework.md) | Next.js App Router with static export |
| [ADR-003](docs/adr/003-terraform-for-infrastructure.md) | Terraform for IaC |
| [ADR-004](docs/adr/004-ecs-fargate-for-compute.md) | ECS Fargate for compute |
| [ADR-005](docs/adr/005-postgresql-as-database.md) | PostgreSQL as database |
| [ADR-006](docs/adr/006-multi-tenancy-strategy.md) | Multi-tenancy via tenant_id isolation |
| [ADR-007](docs/adr/007-jwt-for-authentication.md) | JWT for authentication |
| [ADR-008](docs/adr/008-s3-for-medical-image-storage.md) | S3 for medical image storage |
| [ADR-009](docs/adr/009-structured-logging.md) | Structured JSON logging |
| [ADR-010](docs/adr/010-ddd-with-cqrs-architecture.md) | DDD with CQRS |
| [ADR-011](docs/adr/011-rest-api-design.md) | REST API design |
| [ADR-012](docs/adr/012-modular-monolith-architecture.md) | Modular monolith architecture |
| [ADR-013](docs/adr/013-pgx-sqlc-for-database-access.md) | pgx + sqlc for database access |
| [ADR-014](docs/adr/014-golang-migrate-for-migrations.md) | golang-migrate for migrations |
| [ADR-015](docs/adr/015-frontend-feature-based-architecture.md) | Frontend feature-based architecture |
| [ADR-016](docs/adr/016-design-first-api-documentation.md) | Design-First API documentation |

---

## Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js App Router, TypeScript, pnpm, TanStack Query, openapi-fetch, React Hook Form, Zod, Tailwind CSS, shadcn/ui |
| Frontend Types | openapi-typescript (generated from OpenAPI) |
| Frontend Testing | Vitest, `@testing-library/react`, `@testing-library/user-event`, MSW, `@vitest/coverage-v8` |
| Backend | Go, `net/http`, `http.ServeMux`, `envconfig` |
| Backend API | oapi-codegen (generated from OpenAPI) |
| Database | PostgreSQL (Amazon RDS) |
| DB Access | pgx + sqlc (type-safe SQL) |
| Migrations | golang-migrate |
| Logging | `log/slog` (stdlib) |
| Backend Testing | `testing`, `testify/assert`, `httptest`, `go-cmp`, `testcontainers-go` |
| Storage | Amazon S3 |
| Compute | Amazon ECS Fargate |
| Infrastructure | Terraform |
| Auth | JWT (HMAC-SHA256) |
| API Contract | OpenAPI 3.1.3 (`spec/openapi.yaml`) |

---

## Repository Structure

```
med-vault/
├── backend/           # Go backend (DDD layers)
├── frontend/          # Next.js App Router (feature-based architecture)
│   ├── app/           # Pages (routing and composition only)
│   ├── features/      # Business capabilities (authentication, patients, doctors, admin)
│   ├── infrastructure:# openapi-fetch, auth, query client, config
│   ├── generated/     # Generated TypeScript types from OpenAPI
│   └── shared/        # Layouts, base UI, utilities, global types
├── infrastructure/    # Terraform (modules + environments)
│   └── terraform/
│       ├── modules/       # Reusable platform capabilities
│       └── environments/  # Environment-specific configs
├── spec/              # OpenAPI contract (single source of truth) [Phase 1]
│   └── openapi.yaml   # API contract (to be created in Phase 1)
└── docs/              # Project documentation
    ├── adr/           # Architecture Decision Records
    └── diagrams/      # Architecture diagrams
```

---

## Quick Start

> Prerequisites: Go, Node.js (pnpm), Terraform, AWS CLI configured

```bash
# Backend
cd backend
go mod init github.com/med-vault/backend
go run ./cmd/server

# Frontend
cd frontend
pnpm create next-app@latest . --typescript --tailwind --eslint --app --import-alias "@/*"
pnpm build    # exports static files to out/
pnpm dev

# Infrastructure
cd infrastructure
terraform init
terraform plan
terraform apply
```

---

## Roles

| Role | Capabilities |
|------|-------------|
| Patient | Register, submit cases, upload images, view history |
| Doctor | View assigned cases, review images, write diagnoses |
| Administrator | Assign doctors, close cases, inspect audit logs |

---

## License

Internal use only. Not for production deployment.
