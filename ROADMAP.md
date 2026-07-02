# Roadmap

This document outlines the phased delivery plan for MedVault. Each phase builds on the previous one, following the project's documentation-first principle.

---

## Phase 1: Foundation

**Goal:** Project structure, documentation, and tooling setup.

| Task | Status |
|------|--------|
| Complete documentation (ARCHITECTURE, DOMAIN, REQUIREMENTS, SECURITY) | ✅ |
| Create ADRs for key technology decisions | ✅ |
| Initialize Go module | ⬜ |
| Initialize Next.js project (App Router) | ⬜ |
| Initialize Terraform project | ⬜ |
| Setup Makefile / justfile for common commands | ⬜ |
| Configure linter (golangci-lint, ESLint) | ⬜ |
| Configure formatter (goimports, Prettier) | ⬜ |

---

## Phase 2: Infrastructure

**Goal:** Deploy core AWS infrastructure via Terraform.

| Task | Status |
|------|--------|
| VPC module (public/private subnets, NAT, routes) | ⬜ |
| RDS PostgreSQL module (private subnet, encryption) | ⬜ |
| S3 module (medical images, audit logs) | ⬜ |
| ECS Fargate module (cluster, service, task definition) | ⬜ |
| ALB module (public subnet, TLS termination) | ⬜ |
| IAM module (roles, policies) | ⬜ |
| Secrets Manager module | ⬜ |
| CloudWatch Logs module | ⬜ |
| WAF module | ⬜ |
| Route 53 + CloudFront (optional for PoC) | ⬜ |

---

## Phase 3: Backend Foundation

**Goal:** Go backend with DDD structure, authentication, and tenant isolation.

**Stack:** `net/http`, `http.ServeMux`, `envconfig`, `pgx`, `sqlc`, `golang-migrate`, `log/slog`, `testing` + `httptest`

| Task | Status |
|------|--------|
| Project structure (domain, application, infrastructure) | ⬜ |
| Domain layer (aggregates, entities, value objects) | ⬜ |
| `envconfig` configuration loading | ⬜ |
| `pgx` connection pool setup | ⬜ |
| `golang-migrate` schema migrations | ⬜ |
| `sqlc` query code generation | ⬜ |
| `net/http` server with `http.ServeMux` routing | ⬜ |
| JWT authentication middleware | ⬜ |
| Tenant context middleware | ⬜ |
| RBAC middleware | ⬜ |
| Repository interfaces and implementations | ⬜ |
| Error handling and response format | ⬜ |
| `log/slog` structured logging | ⬜ |
| Health check endpoint | ⬜ |
| Unit tests with `testing` + `httptest` | ⬜ |

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

**Stack:** Next.js App Router, TypeScript, pnpm, TanStack Query, Axios, React Hook Form, Zod, Tailwind CSS, shadcn/ui

| Task | Status |
|------|--------|
| Project setup (Next.js App Router, TypeScript, pnpm, static export) | ⬜ |
| Feature-based directory structure (features/, infrastructure/, shared/) | ⬜ |
| Infrastructure layer (Axios instance, TanStack Query client, auth helpers) | ⬜ |
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

---

## Phase 8: Polish

**Goal:** Security hardening, observability, and documentation.

| Task | Status |
|------|--------|
| Security review | ⬜ |
| Input validation | ⬜ |
| Rate limiting | ⬜ |
| CloudWatch dashboards | ⬜ |
| CloudTrail integration | ⬜ |
| Updated README with deployment instructions | ⬜ |
| Architecture diagrams (PNG/SVG) | ⬜ |

---

## Phase 9: Future (Post-PoC)

| Task | Status |
|------|--------|
| CI/CD pipeline (GitHub Actions) | ⬜ |
| Container registry (ECR) | ⬜ |
| Blue/green deployments | ⬜ |
| Notifications (email/SMS) | ⬜ |
