# AGENTS.md

Instructions for AI agents working on this repository.

---

## Project

MedVault is a healthcare platform PoC demonstrating secure, multi-tenant architecture on AWS. See [README](README.md) for overview, [docs/](docs/) for details.

**Primary goal:** Architecture demonstration, not production software.

---

## Rules

### Never

- Hardcode secrets or credentials
- Disable encryption or TLS
- Bypass tenant isolation (`WHERE tenant_id = $1`)
- Bypass authentication or authorization
- Add unnecessary complexity or abstractions
- Use `// TODO` without a linked issue

### Always

- Follow DDD bounded context boundaries (see [DOMAIN.md](docs/DOMAIN.md))
- Include `tenant_id` in every database query
- Validate inputs at API boundaries
- Log state-changing operations (audit trail)
- Update documentation when changing architecture
- Prefer explicit code over clever abstractions
- Write tests as part of implementation (see [TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md))
- Run quality gates before pushing (see [QUALITY_GATES.md](docs/QUALITY_GATES.md))

---

## Source of Truth

| Concern | Document |
|---------|----------|
| System design | [ARCHITECTURE.md](docs/ARCHITECTURE.md) |
| Domain model | [DOMAIN.md](docs/DOMAIN.md) |
| Requirements | [REQUIREMENTS.md](docs/REQUIREMENTS.md) |
| Security controls | [SECURITY.md](docs/SECURITY.md) |
| Infrastructure | [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md) |
| Testing philosophy | [TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md) |
| Quality gates | [QUALITY_GATES.md](docs/QUALITY_GATES.md) |
| Engineering principles | [PROJECT_PRINCIPLES.md](docs/PROJECT_PRINCIPLES.md) |
| Implementation progress | [ROADMAP.md](ROADMAP.md) |
| Acceptance criteria | [CHECKLIST.md](docs/CHECKLIST.md) |
| Technology decisions | [docs/adr/](docs/adr/) |

---

## Architecture

- **Pattern:** DDD with CQRS (see [ADR-010](docs/adr/010-ddd-with-cqrs-architecture.md))
- **Event delivery:** Transactional Outbox with polling on PostgreSQL (see [ADR-017](docs/adr/017-transactional-outbox.md))
- **Bounded Contexts:** Identity & Access, Clinical, Imaging, Audit
- **Multi-tenancy:** Shared DB with `tenant_id` column; users can belong to multiple tenants with role per tenant (see [ADR-006](docs/adr/006-multi-tenancy-strategy.md))
- **Backend:** Go, Clean Architecture layers (domain → application → infrastructure)
- **Frontend:** Next.js App Router (static export, Client Components only), Feature-Based Architecture (see [ADR-015](docs/adr/015-frontend-feature-based-architecture.md)), no API Routes, no Server Actions, no SSR
- **Infrastructure:** Terraform, ECS Fargate, RDS PostgreSQL, S3

---

## Code Conventions

### Go Backend

- Modular monolith (see [ADR-012](docs/adr/012-modular-monolith-architecture.md))
- Modules: auth, clinical, imaging, audit + shared kernel
- Domain layer: no external dependencies
- Value objects: immutable, equality by attributes
- Aggregates: enforce consistency boundaries
- Repositories: interfaces in domain, implementations in infrastructure
- Database access: pgx + sqlc (see [ADR-013](docs/adr/013-pgx-sqlc-for-database-access.md))
- Migrations: golang-migrate (see [ADR-014](docs/adr/014-golang-migrate-for-migrations.md))
- Migrations run as CI/CD step before deployment, never at application startup
- HTTP: `net/http` stdlib, router: `http.ServeMux`
- Config: `envconfig` (12-factor, env vars)
- Logging: `log/slog` (stdlib structured JSON, no PHI in logs)
- Tests: `testing` (runner), `testify/assert` (assertions), `httptest` (HTTP), `go-cmp` (struct comparison), `testcontainers-go` (integration DB)
- Inter-module communication: via ports.go interfaces (no direct DB access across modules)
- Errors: return domain errors, wrap infrastructure errors
- **API:** Design-First with OpenAPI, generate interfaces via `oapi-codegen` (see [ADR-016](docs/adr/016-design-first-api-documentation.md))

### Frontend (Next.js)

- Feature-Based Architecture (see [ADR-015](docs/adr/015-frontend-feature-based-architecture.md))
- Next.js App Router with static export
- Client Components only (`'use client'`)
- No API Routes, no Server Actions, no SSR, no ISR
- TypeScript strict mode
- Package manager: pnpm
- Server state: TanStack Query
- HTTP client: `openapi-fetch` (type-safe, generated from OpenAPI)
- Types: `openapi-typescript` (generated from `spec/openapi.yaml`)
- Forms: React Hook Form + Zod validation
- Styling: Tailwind CSS + shadcn/ui
- Testing: Vitest, `@testing-library/react`, `@testing-library/user-event`, MSW, `@vitest/coverage-v8`
- No PHI stored in browser
- No business logic in frontend — rules belong in Go backend

**Directory structure:**

```
frontend/
├── app/                 # Pages (routing and composition only)
├── features/            # One dir per business capability
│   ├── authentication/  # components, hooks, services, schemas, types
│   ├── patients/
│   ├── doctors/
│   ├── admin/
│   └── shared/          # Cross-feature reusable elements
├── infrastructure/      # openapi-fetch, auth, query client, config
├── generated/           # Generated TypeScript types from OpenAPI
└── shared/              # Layouts, base UI, utilities, global types
```

**Layer rules:**

| Layer | Allowed | Forbidden |
|-------|---------|-----------|
| Pages | Route composition, feature wiring | Business logic, direct HTTP |
| Components | Presentation, props-only data flow | HTTP calls, business rules |
| Hooks | TanStack Query, UI orchestration | Raw HTTP requests |
| Services | openapi-fetch calls, endpoint definitions | Business rules |
| Schemas | Zod validation (request, response, form) | Business logic |
| Infrastructure | openapi-fetch instance, auth, query client | Business logic |
| Shared | Reusable UI, layouts, utilities | Dumping ground |

**Feature rules:**

- Each feature is self-contained (components, hooks, services, schemas, types)
- Avoid unnecessary coupling between features
- Shared code only when there is a clear cross-feature need
- Every new feature follows the same directory structure
- Consistency over clever abstractions

### Terraform

- Modular structure representing platform capabilities (see [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md))
- Modules: `network`, `application`, `database`, `storage`, `security`, `observability`
- No hardcoded values — variables for all configurable parameters
- State in S3 (remote backend) with versioning and encryption
- Security by default: private subnets, encryption at rest/in transit, least privilege IAM
- Readability preferred over reducing duplicated code
- Never duplicate infrastructure code between environments
- Major infrastructure decisions recorded as ADRs

---

## When Modifying

1. Read the relevant source-of-truth document first
2. Check if the change affects other bounded contexts
3. Update documentation if architecture changes
4. Follow existing patterns in the codebase
5. Verify tenant isolation is maintained
