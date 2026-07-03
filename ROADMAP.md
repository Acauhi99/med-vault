# Roadmap

This document outlines the phased delivery plan for MedVault. Each phase builds on the previous one, following the project's documentation-first principle.

**Testing philosophy:** Unit tests and integration tests only. See [TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md).

**Quality gates:** Pre-commit, pre-push, unified Taskfile. See [QUALITY_GATES.md](docs/QUALITY_GATES.md).

**Acceptance criteria:** Detailed checklist per phase. See [CHECKLIST.md](docs/CHECKLIST.md).

---

## Phase 1: Foundation

**Goal:** Project structure, documentation, tooling setup, and API contract.

| Task | Status |
|------|--------|
| Complete documentation (ARCHITECTURE, DOMAIN, REQUIREMENTS, SECURITY) | ✅ |
| Create ADRs for key technology decisions | ✅ |
| Define OpenAPI 3.1.3 contract (`spec/openapi.yaml`) | ✅ |
| Setup `oapi-codegen` for Go backend code generation | ⬜ |
| Setup `openapi-typescript` + `openapi-fetch` for frontend | ⬜ |
| Initialize Go module | ⬜ |
| Initialize Next.js project (App Router) | ⬜ |
| Initialize Terraform project | ⬜ |
| Setup Taskfile with `format`, `lint`, `validate`, `test` tasks | ✅ |
| Configure `gofumpt` + `golangci-lint` (backend) | ⬜ |
| Configure Biome (frontend) | ⬜ |
| Configure `tflint` + Checkov (infrastructure) | ⬜ |
| Configure Git pre-commit and pre-push hooks | ⬜ |

---

## Phase 2: Infrastructure

**Goal:** Deploy core AWS infrastructure via Terraform.

**Stack:** Terraform, modules: `network`, `application`, `database`, `storage`, `security`, `observability` (see [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md) for philosophy and module design)

| Task | Status |
|------|--------|
| `network` module (VPC, public/private subnets, NAT, routes) | ⬜ |
| `database` module (RDS PostgreSQL, private subnet, encryption) | ⬜ |
| `storage` module (S3 buckets, lifecycle rules, encryption) | ⬜ |
| `application` module (ECS Fargate, ALB, task definition) | ⬜ |
| `security` module (IAM roles, policies, KMS, Secrets Manager) | ⬜ |
| `observability` module (CloudWatch, CloudTrail, VPC Flow Logs) | ⬜ |
| Production environment composition | ⬜ |
| Remote state in S3 with versioning and encryption | ⬜ |
| WAF associated with ALB | ⬜ |
| Route 53 + CloudFront (optional for PoC) | ⬜ |

---

## Phase 3: Backend Foundation

**Goal:** Go backend with DDD structure, authentication, and tenant isolation.

**Stack:** `net/http`, `http.ServeMux`, `envconfig`, `pgx`, `sqlc`, `golang-migrate`, `log/slog`, `testing` + `httptest` + `testify/assert` + `go-cmp` + `testcontainers-go`

| Task | Status |
|------|--------|
| Project structure (domain, application, infrastructure) | ⬜ |
| Domain layer (aggregates, entities, value objects) | ⬜ |
| `envconfig` configuration loading | ⬜ |
| `pgx` connection pool setup | ⬜ |
| `golang-migrate` schema migrations | ⬜ |
| `sqlc` query code generation | ⬜ |
| `oapi-codegen` server interface generation | ⬜ |
| `net/http` server with `http.ServeMux` routing | ⬜ |
| JWT authentication middleware | ⬜ |
| Tenant context middleware | ⬜ |
| RBAC middleware | ⬜ |
| Rate limiting middleware (auth endpoints) | ⬜ |
| Repository interfaces and implementations | ⬜ |
| Error handling and response format | ⬜ |
| `log/slog` structured logging | ⬜ |
| Health check endpoint | ⬜ |
| Unit tests with `testing` + `testify/assert` | ⬜ |
| HTTP handler tests with `httptest` | ⬜ |
| Struct comparison tests with `go-cmp` | ⬜ |
| Integration tests with `testcontainers-go` | ⬜ |

---

## Phase 4: Identity & Access

**Goal:** User registration, authentication, and tenant management.

| Task | Status |
|------|--------|
| Tenant aggregate and repository | ⬜ |
| User aggregate and repository | ⬜ |
| Register user command | ⬜ |
| Authenticate user command | ⬜ |
| Refresh token command | ⬜ |
| Get current user query | ⬜ |
| Add user to tenant command | ⬜ |
| Remove user from tenant command | ⬜ |
| List tenant members query | ⬜ |
| Reactivate tenant command | ⬜ |
| Audit logging for auth events | ⬜ |

---

## Phase 5: Clinical Core

**Goal:** Medical case management with symptoms and diagnoses.

| Task | Status |
|------|--------|
| Case aggregate and repository | ⬜ |
| Symptom entity | ⬜ |
| Diagnosis value object | ⬜ |
| Create case command | ⬜ |
| Add symptom command | ⬜ |
| Assign doctor command | ⬜ |
| Write diagnosis command | ⬜ |
| Close case command | ⬜ |
| List cases queries (by patient, doctor, admin) | ⬜ |
| Get case query | ⬜ |
| Domain events and projections | ⬜ |
| Audit logging for clinical events | ⬜ |

---

## Phase 6: Imaging

**Goal:** Medical image upload and retrieval.

| Task | Status |
|------|--------|
| Image aggregate and repository | ⬜ |
| S3 pre-signed URL generation | ⬜ |
| Request upload URL command | ⬜ |
| Confirm upload command | ⬜ |
| List images query | ⬜ |
| Get download URL query | ⬜ |
| Audit logging for imaging events | ⬜ |

---

## Phase 7: Frontend

**Goal:** Next.js App Router SPA with feature-based architecture, authentication, and core workflows.

**Stack:** Next.js App Router, TypeScript, pnpm, TanStack Query, openapi-fetch, openapi-typescript, React Hook Form, Zod, Tailwind CSS, shadcn/ui, Vitest, `@testing-library/react`, `@testing-library/user-event`, MSW, `@vitest/coverage-v8`

| Task | Status |
|------|--------|
| Project setup (Next.js App Router, TypeScript, pnpm, static export) | ⬜ |
| `openapi-typescript` generation from `spec/openapi.yaml` | ⬜ |
| Feature-based directory structure (features/, infrastructure/, shared/) | ⬜ |
| Infrastructure layer (openapi-fetch instance, TanStack Query client, auth helpers) | ⬜ |
| Shared components (layouts, navigation, base UI) | ⬜ |
| Authentication feature (login, register — components, hooks, services, schemas) | ⬜ |
| Patients feature (dashboard, case list — components, hooks, services) | ⬜ |
| Doctors feature (assigned cases, diagnosis — components, hooks, services) | ⬜ |
| Admin feature (case management, audit logs — components, hooks, services) | ⬜ |
| Case creation form (React Hook Form + Zod) | ⬜ |
| Symptom submission form | ⬜ |
| Image upload component | ⬜ |
| Diagnosis view | ⬜ |
| Audit log viewer | ⬜ |
| Tenant switcher component | ⬜ |
| Vitest configured with `@testing-library/react` and MSW | ⬜ |
| Component tests with `@testing-library/react` + `@testing-library/user-event` | ⬜ |
| API mocking with MSW | ⬜ |
| Coverage reporting with `@vitest/coverage-v8` | ⬜ |

---

## Phase 8: Polish

**Goal:** Security hardening, observability, and documentation.

| Task | Status |
|------|--------|
| Security review | ⬜ |
| Input validation | ⬜ |
| CloudWatch dashboards | ⬜ |
| CloudTrail integration | ⬜ |
| Updated README with deployment instructions | ⬜ |
| Architecture diagrams (PNG/SVG) | ⬜ |

---

## Phase 9: CI/CD (Post-PoC)

**Goal:** Implement the three-pipeline deployment architecture. See [CI_CD_STRATEGY.md](docs/CI_CD_STRATEGY.md) for full philosophy and boundaries.

| Task | Status |
|------|--------|
| Infrastructure pipeline (Terraform init/validate/plan/apply) | ⬜ |
| Backend pipeline (lint, test, build, migrate, deploy, health check) | ⬜ |
| Frontend pipeline (lint, typecheck, test, build, export, S3 upload, CF invalidation) | ⬜ |
| GitHub OIDC for AWS authentication | ⬜ |
| Path-based triggers (infrastructure/, backend/, frontend/) | ⬜ |
| Concurrency groups per pipeline | ⬜ |
| Deployment validation (health checks, smoke tests) | ⬜ |
| Notifications (email/SMS) | ⬜ |

---

## Phase 10: HIPAA Compliance

**Goal:** Document and implement HIPAA Privacy Rule, Security Rule, and Breach Notification Rule controls. See [SECURITY.md](docs/SECURITY.md) for the full security architecture.

| Task | Status |
|------|--------|
| Notice of Privacy Practices (NPP) defined | ⬜ |
| Patient rights implemented (access, amendment, accounting, restrictions, confidential communications) | ⬜ |
| Minimum Necessary Standard enforced per role | ⬜ |
| Business Associate Agreement (BAA) with AWS | ⬜ |
| Uses and Disclosures policy documented | ⬜ |
| De-identification methods documented | ⬜ |
| Breach assessment process defined | ⬜ |
| Notification requirements documented (individual, HHS, media) | ⬜ |
| Breach response team identified | ⬜ |
| Breach documentation template defined | ⬜ |
| Security Officer and Privacy Officer designated | ⬜ |
| Risk analysis and management documented | ⬜ |
| Sanction policy documented | ⬜ |
| Workforce security procedures documented | ⬜ |
| Security awareness training program | ⬜ |
| Incident response procedures documented | ⬜ |
| Contingency plan documented (backup, DR, emergency mode) | ⬜ |
| Annual security evaluation scheduled | ⬜ |
| Physical safeguards documented (facility, workstation, device controls) | ⬜ |
| Automatic logoff configured (15 min inactivity) | ⬜ |
| Documentation retention enforced (6 years via S3 lifecycle) | ⬜ |
