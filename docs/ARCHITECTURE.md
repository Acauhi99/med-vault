# Architecture

MedVault is a multi-tenant healthcare platform built on AWS. This document describes the system architecture, component responsibilities, data flow, and deployment model.

---

## System Context

```
Internet
    │
    ▼
Route 53 (DNS)
    │
    ▼
CloudFront (CDN + SSL)
    │
    ▼
S3 (Next.js Static Export)
    │
    ▼
AWS WAF (Request filtering)
    │
    ▼
ALB (Load balancing + TLS)
    │
    ▼
ECS Fargate (Go Backend)
    │
    ├──▶ Amazon RDS PostgreSQL
    └──▶ S3 (Medical Images)
```

---

## Components

### Frontend (Next.js App Router — Static Export)

**Purpose:** Presentation layer for patients, doctors, and administrators.

**Responsibilities:**
- UI rendering and navigation
- Form handling and client-side validation
- Authentication state management
- API communication with Go backend
- User experience

**The frontend is NOT responsible for:**
- Business rules
- Authorization logic
- Tenant isolation
- Data transformation beyond presentation

**Architectural Style:** Feature-Based Architecture (see [ADR-015](adr/015-frontend-feature-based-architecture.md))

The frontend does NOT follow MVC, MVP, or MVVM. Every business capability is implemented as an isolated feature, inspired by Vertical Slice Architecture and DDD principles.

**Stack:**

| Tool | Purpose |
|------|---------|
| Next.js (App Router) | Framework, routing, layouts |
| TypeScript | Type safety |
| pnpm | Package management |
| TanStack Query | Server state management, caching |
| Axios | HTTP client for REST API |
| React Hook Form | Form handling |
| Zod | Schema validation |
| Tailwind CSS | Utility-first styling |
| shadcn/ui | Component system |

**Constraints (intentionally not used):**

| Feature | Reason |
|---------|--------|
| Server Components | No business logic in frontend |
| Server Actions | Mutations go through Go REST API |
| API Routes | Backend is exclusively Go |
| SSR | Static export only |
| ISR | Static export only |
| BFF (Backend-for-Frontend) | Go backend serves all consumers |
| Node.js middleware | Static export, no runtime required |

**Directory Structure:**

```
frontend/
├── app/                        # Pages (routing and composition only)
├── features/                   # Business capabilities (one dir per feature)
│   ├── authentication/         # Example feature
│   │   ├── components/         # Presentation components
│   │   ├── hooks/              # TanStack Query hooks, UI orchestration
│   │   ├── services/           # Axios API calls
│   │   ├── schemas/            # Zod validation schemas
│   │   ├── types/              # TypeScript types
│   │   └── index.ts            # Public exports
│   ├── patients/
│   ├── doctors/
│   ├── admin/
│   └── shared/                 # Cross-feature reusable elements
├── infrastructure/             # External integrations
│   ├── api/                    # Axios instance, interceptors
│   ├── auth/                   # JWT storage, token refresh
│   ├── query/                  # TanStack Query client config
│   └── config/                 # Environment configuration
└── shared/                     # App-wide reusable code
    ├── components/             # Layout, navigation, base UI
    ├── layouts/                # Page layouts
    ├── lib/                    # Utilities, helpers
    └── types/                  # Global types
```

**Layer Responsibilities:**

| Layer | Responsibility | Rules |
|-------|---------------|-------|
| Pages | Route composition, feature wiring | No business logic |
| Components | Presentation only | Props-only data flow, no HTTP calls |
| Hooks | UI orchestration, TanStack Query | No raw HTTP requests |
| Services | API communication (Axios) | No business rules |
| Schemas | Validation (Zod) | One schema per feature |
| Infrastructure | Axios, auth, query client | External integrations only |
| Shared | Reusable UI, layouts, utilities | Never a dumping ground |

**Request Flow:**

```
Component → Hook → TanStack Query → Service → Axios → Go REST API
                                                          ↓
Component ← Hook ← TanStack Query Cache ←←←←←←←←← Response
```

**Business Logic:** All business rules live exclusively in the Go backend. The frontend validates user input (Zod) and handles UI state only.

**Multi-Tenant:** The frontend propagates tenant context from JWT. Tenant isolation is enforced by the backend. The frontend never enforces tenant security.

**Testing:** Unit tests and integration tests only (see [TESTING_STRATEGY.md](TESTING_STRATEGY.md) for philosophy, [ADR-015](adr/015-frontend-feature-based-architecture.md#testing-strategy) for frontend stack)

| Tool | Purpose |
|------|---------|
| Vitest | Test runner (fast, modern, TypeScript-native) |
| `@testing-library/react` | Component rendering and interaction |
| `@testing-library/user-event` | Real user actions (type, click, navigate) |
| MSW (Mock Service Worker) | API mocking without coupling to Axios |
| `@vitest/coverage-v8` | Coverage via V8 engine |

- **Unit tests:** Components, hooks, services — isolated, fast, mocked APIs via MSW
- **Integration tests:** Pages, feature workflows — route composition, feature wiring
- No end-to-end tests

**Deployment:**

```
next build + next export → S3 → CloudFront → Client
```

- Static HTML/CSS/JS exported at build time
- S3 serves static assets
- CloudFront provides CDN and TLS termination
- No Node.js runtime in production

**Security:**
- No PHI stored locally
- Session tokens in httpOnly cookies
- API calls authenticated via JWT

**Why this architecture:**
- Feature isolation → changes are localized, easy to reason about
- Predictable structure → AI agents infer patterns with minimal prompting
- Aligns with backend DDD philosophy (features ≈ bounded contexts)
- Components are pure presentation → easy to test
- Services are pure HTTP → easy to test
- Static export keeps AWS architecture simple
- Go backend owns authentication, authorization, multi-tenancy, and business logic

---

### Backend (Go Modular Monolith)

**Purpose:** Core business logic, API layer, tenant isolation.

**Responsibilities:**
- JWT authentication and token validation
- Role-based authorization (Patient, Doctor, Administrator)
- Tenant context extraction and propagation
- REST API endpoints
- Request validation
- Audit logging
- File upload handling (delegates to S3)
- Database access with tenant isolation

**Architecture:** Modular Monolith with DDD + CQRS per module.

**Backend Stack:**

| Layer | Tool |
|-------|------|
| HTTP | `net/http` (stdlib) |
| Router | `http.ServeMux` |
| Config | `envconfig` |
| Database | `pgx` |
| Queries | `sqlc` |
| Migrations | `golang-migrate` |
| Logging | `log/slog` (stdlib) |
| Tests | `testing` + `httptest` + `testify/assert` + `go-cmp` + `testcontainers-go` |
| Container | Docker |
| Deploy | ECS Fargate |

```
┌─────────────────────────────────────────────────────────────────┐
│                     http.ServeMux (Router)                       │
├─────────┬─────────┬──────────┬──────────┬───────────────────────┤
│  Auth   │Clinical │ Imaging  │  Audit   │  Shared Kernel        │
│ Module  │ Module  │  Module  │  Module  │  (middleware, tenant)  │
├─────────┴─────────┴──────────┴──────────┴───────────────────────┤
│                     Module Contracts                            │
│            (interfaces for inter-module communication)           │
├─────────────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                           │
│        (pgx + sqlc repos, S3, auth, log/slog)                   │
└─────────────────────────────────────────────────────────────────┘
```

**Module structure (each module):**

```
module/
├── domain/          # Aggregates, entities, value objects, events
├── application/     # Command/query handlers, services
├── infrastructure/  # Repository implementations (pgx + sqlc)
├── ports.go         # Public interfaces for other modules
└── module.go        # Module registration and wiring
```

**Modular monolith principles:**
- Modules communicate via in-process interfaces (no direct DB access across modules)
- Each module owns its database tables
- Shared Kernel provides cross-cutting concerns (middleware, tenant context, logging)
- Future decomposition: replace in-process interfaces with network calls
- No circular dependencies between modules

**CQRS principles (per module):**
- Commands mutate state through aggregates
- Queries read from optimized read models
- Domain events bridge command and query sides
- No direct reads from write-optimized aggregates

**Key invariants:**
- Every query includes `tenant_id` filter
- No business logic leaks into infrastructure
- External dependencies injected via interfaces
- Aggregates enforce consistency boundaries
- Read models are eventually consistent with write models

**Database access:** pgx + sqlc for type-safe SQL (see [ADR-013](adr/013-pgx-sqlc-for-database-access.md))

**Testing:** Unit tests and integration tests only (see [TESTING_STRATEGY.md](TESTING_STRATEGY.md) for philosophy, [ADR-001](adr/001-go-as-backend-language.md#testing-strategy) for backend stack)

| Tool | Purpose |
|------|---------|
| `testing` | Test runner (stdlib) |
| `testify/assert` | Assertions |
| `net/http/httptest` | HTTP handler testing |
| `go-cmp` | Struct comparison |
| `testcontainers-go` | Real PostgreSQL in integration tests |
| `go test -cover` | Coverage reporting |
| `testing.B` | Benchmarks |
| `testing.F` | Fuzz testing |

- **Unit tests:** Domain logic, value objects, aggregates — isolated, fast, no I/O
- **Integration tests:** Repositories, HTTP handlers — real PostgreSQL via testcontainers, real HTTP via httptest
- No end-to-end tests

---

### Database (Amazon RDS PostgreSQL)

**Purpose:** Persistent storage for all business entities.

**Responsibilities:**
- Store tenants, users, cases, diagnoses, audit logs
- Enforce referential integrity
- Support tenant isolation via `tenant_id` column

**Design:**
- Shared database, shared schema with `tenant_id` column
- Every table includes `tenant_id` as part of primary key or foreign key
- Row-Level Security (RLS) as defense-in-depth

---

### Medical Image Storage (Amazon S3)

**Purpose:** Store uploaded medical images (X-rays, scans, etc.).

**Responsibilities:**
- Encrypted storage (SSE-S3 or SSE-KMS)
- Pre-signed URLs for temporary access
- Bucket policy enforces TLS-only access

**Security:**
- No public access
- Bucket versioning enabled
- Access logged to CloudWatch

---

## Data Flow

### Authentication Flow

```
Client → ALB → Backend
  1. Client sends credentials (email + password)
  2. Backend validates against database (bcrypt)
  3. Backend returns JWT (access token + refresh token)
  4. Client stores token in httpOnly cookie
  5. Subsequent requests include JWT in Authorization header
```

### Request Lifecycle (Authenticated)

```
Client → WAF → ALB → ECS → Backend
  1. WAF filters malicious requests
  2. ALB terminates TLS, forwards to ECS
  3. Backend extracts JWT from request
  4. Backend validates JWT signature and expiration
  5. Backend extracts tenant_id and user_id from JWT claims
  6. Backend sets tenant context for request
  7. Handler executes business logic with tenant isolation
  8. Response returned to client
```

### File Upload Flow

```
Client → Backend → S3
  1. Client requests pre-signed upload URL
  2. Backend generates pre-signed URL (limited time, limited scope)
  3. Client uploads directly to S3 via pre-signed URL
  4. S3 notifies backend (optional: S3 event or callback)
  5. Backend records image metadata in database
```

---

## Multi-Tenancy Model

**Strategy:** Shared database, shared schema with `tenant_id` isolation.

**Isolation points:**
- **Authentication:** JWT contains `tenant_id` claim
- **Authorization:** Middleware extracts `tenant_id`, injects into context
- **Database:** Every query includes `WHERE tenant_id = $1`
- **Storage:** S3 paths prefixed with `/{tenant_id}/`
- **API:** Middleware rejects requests without valid tenant context

**Defense-in-depth:**
1. Application-level `tenant_id` filtering (primary)
2. PostgreSQL Row-Level Security (secondary)
3. IAM policies scoped per tenant (future: per-tenant KMS keys)

---

## Security Architecture

### Encryption

- **At rest:** RDS encryption (AWS KMS), S3 encryption (SSE-S3 or SSE-KMS)
- **In transit:** TLS 1.2+ everywhere (CloudFront, ALB, RDS SSL)
- **Application:** JWT signing (HMAC-SHA256 or RSA)

### Secrets Management

- Database credentials: AWS Secrets Manager
- JWT signing key: AWS Secrets Manager
- No hardcoded secrets in code or environment variables

### Audit Logging

- Every state-changing operation logged
- Audit log includes: user_id, tenant_id, action, resource, timestamp, IP
- Logs stored in CloudWatch Logs
- CloudTrail for AWS API audit trail

### Network Security

- Backend runs in private subnets (no public IP)
- RDS in private subnets
- ALB in public subnets only
- Security groups restrict inter-service communication
- VPC flow logs enabled

---

## API Documentation Strategy (Design-First)

**Approach:** OpenAPI 3.1.3 as single source of truth for all API contracts.

### Contract Flow

```
spec/openapi.yaml (single source of truth)
    ├── oapi-codegen → Go server interfaces + types
    └── openapi-typescript + openapi-fetch → Type-safe TypeScript client
```

### Directory Structure

```
med-vault/
├── spec/
│   └── openapi.yaml              # API contract
├── backend/
│   └── internal/generated/       # Generated Go interfaces (oapi-codegen)
└── frontend/
    └── src/generated/            # Generated TypeScript types (openapi-typescript)
```

### Code Generation

| Tool | Purpose | Output |
|------|---------|--------|
| `oapi-codegen` | Generate Go server stubs from OpenAPI | `server.go`, `types.go` |
| `openapi-typescript` | Generate TypeScript types from OpenAPI | `api.d.ts` |
| `openapi-fetch` | Type-safe HTTP client | Runtime dependency |

### Rules

- **Never** manually write HTTP contracts in Go or TypeScript
- **Never** generate OpenAPI from Go code
- **Always** define endpoints in `spec/openapi.yaml` first
- **Always** run code generation after spec changes
- Frontend consumes generated types — no manual `fetch` with `any` payloads
- Backend implements generated interfaces — no ad-hoc handler signatures

### API Evolution

1. Define or update `spec/openapi.yaml`
2. Run generation scripts
3. Backend implements generated interfaces
4. Frontend consumes generated types
5. Contract drift is impossible

See [ADR-016](adr/016-design-first-api-documentation.md) for full rationale.

---

## Deployment Model

### Infrastructure

- All resources managed via Terraform
- Modular structure representing platform capabilities (see [INFRASTRUCTURE.md](INFRASTRUCTURE.md))
- Modules: `network`, `application`, `database`, `storage`, `security`, `observability`
- Single environment (PoC) with production-like configuration
- Remote state in S3 with versioning and encryption

### Compute

- ECS Fargate (serverless containers)
- Auto-scaling based on CPU/memory
- Health checks via ALB target groups

### CI/CD (Future)

- GitHub Actions for build/test/deploy
- ECR for container images
- Blue/green or rolling deployments

---

## Observability

### Logging

- `log/slog` structured JSON logs from backend
- CloudWatch Logs for centralized log storage
- Log levels: ERROR, WARN, INFO, DEBUG

### Metrics

- CloudWatch metrics for ECS, ALB, RDS
- Custom metrics: request count, latency, error rate

### Tracing (Future)

- AWS X-Ray for distributed tracing

---

## Cost Considerations

| Service | Cost Strategy |
|---------|---------------|
| ECS Fargate | Use smallest task size that works for PoC |
| RDS | db.t3.micro or db.t4g.micro for PoC |
| S3 | Standard storage, lifecycle policies for old images |
| CloudFront | Use free tier + minimal traffic for PoC |
| CloudWatch | Logs ingestion may cost; set retention policies |

**Goal:** Demonstrate production architecture without production costs.

---

## Trade-offs and Decisions

| Decision | Rationale |
|----------|-----------|
| Shared DB + tenant_id | Simplest multi-tenant model; avoids per-tenant DB overhead |
| ECS Fargate over EC2 | No server management; aligns with "managed services first" |
| REST over GraphQL | Simpler; sufficient for PoC scope |
| JWT over session cookies | Stateless; suitable for SPA + API architecture |
| S3 pre-signed URLs | Direct upload reduces backend bandwidth; time-limited access |
| DDD + CQRS | Explicit domain model; clear separation of reads/writes; testable; aligns with project principles |
| Tactical DDD | Aggregates enforce consistency; Value Objects for validation; Domain Events for decoupling |
| Strategic DDD | Bounded Contexts map to business capabilities; Ubiquitous Language aligns code with domain |

See `docs/adr/` for detailed Architecture Decision Records.
