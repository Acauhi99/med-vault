# Checklist

Acceptance criteria for each phase of MedVault. Mark each item as complete when verified.

**Testing philosophy:** Unit tests and integration tests only. See [TESTING_STRATEGY.md](TESTING_STRATEGY.md).

**Quality gates:** Pre-commit, pre-push, unified Taskfile. See [QUALITY_GATES.md](QUALITY_GATES.md).

**Phased delivery plan:** See [ROADMAP.md](../ROADMAP.md) for the overall implementation timeline.

---

## Phase 1: Foundation

- [ ] Go module initialized (`go.mod` exists)
- [ ] Next.js App Router project initialized (static export)
- [ ] Terraform project initialized (`main.tf` exists)
- [ ] OpenAPI 3.1.3 contract defined (`spec/openapi.yaml`)
- [ ] `oapi-codegen` configured for backend generation
- [ ] `openapi-typescript` configured for frontend generation
- [ ] Taskfile created with `format`, `lint`, `validate`, `test`, `pre-commit`, `pre-push` tasks
- [ ] `gofumpt` configured (backend formatting)
- [ ] `golangci-lint` configured (backend linting)
- [ ] Biome configured (frontend formatting + linting)
- [ ] `tflint` configured (infrastructure linting)
- [ ] Checkov configured (infrastructure security)
- [ ] Git pre-commit hook configured
- [ ] Git pre-push hook configured
- [ ] All documentation files present and non-empty
- [ ] All ADRs present and non-empty

---

## Phase 2: Infrastructure

- [ ] Module structure follows capability-based design (see [INFRASTRUCTURE.md](INFRASTRUCTURE.md))
- [ ] `network` module: VPC, public/private subnets, NAT, route tables, internet gateway
- [ ] `database` module: RDS PostgreSQL in private subnet, encryption enabled
- [ ] `storage` module: S3 buckets with encryption, versioning, lifecycle rules
- [ ] `application` module: ECS Fargate cluster, ALB, task definition, security groups
- [ ] `security` module: IAM roles, policies, KMS, Secrets Manager
- [ ] `observability` module: CloudWatch, CloudTrail, VPC Flow Logs
- [ ] Production environment composes modules correctly
- [ ] Remote state in S3 with versioning and encryption
- [ ] Security groups configured (ALB → ECS → RDS)
- [ ] S3 bucket public access blocked
- [ ] WAF associated with ALB

---

## Phase 3: Backend Foundation

- [ ] Project structure follows DDD layers (domain, application, infrastructure)
- [ ] Domain layer has no external dependencies
- [ ] Aggregates defined per bounded context
- [ ] Value objects are immutable
- [ ] `envconfig` loads configuration from environment variables
- [ ] `pgx` connection pool configured and tested
- [ ] `golang-migrate` migrations run successfully
- [ ] `sqlc` generates type-safe query code
- [ ] `oapi-codegen` generates Go server interfaces from OpenAPI
- [ ] `net/http` server starts with `http.ServeMux` routing
- [ ] JWT middleware validates tokens
- [ ] Tenant context middleware extracts tenant_id
- [ ] RBAC middleware enforces role-based access
- [ ] Repository interfaces defined in domain layer
- [ ] Repository implementations in infrastructure layer
- [ ] Error responses follow standard format
- [ ] `log/slog` structured JSON logging works
- [ ] Health check endpoint returns 200
- [ ] Rate limiting middleware on auth endpoints (login, register, refresh)
- [ ] Unit tests pass with `testing` + `testify/assert`
- [ ] HTTP handler tests pass with `httptest`
- [ ] Struct comparison tests pass with `go-cmp`
- [ ] Integration tests pass with `testcontainers-go`
- [ ] Coverage reporting works with `go test -cover`
- [ ] Dockerfile created with multi-stage build (see [ADR-019](adr/019-docker-image-strategy.md))
- [ ] Build stage: Go toolchain, module download, compilation, validation
- [ ] Runtime stage: distroless, binary + CA certs only
- [ ] CGO disabled, static binary
- [ ] Non-root user, read-only binary
- [ ] `.dockerignore` excludes git, IDE, caches, docs, test artifacts
- [ ] BuildKit cache mounts for Go modules and compiler
- [ ] No Go compiler, source code, or build tools in production image
- [ ] Image builds successfully with `docker build`

---

## Phase 4: Identity & Access

- [ ] Tenant can be created
- [ ] Suspended tenant can be reactivated
- [ ] User can be registered with valid email
- [ ] Email is unique system-wide (duplicate registration rejected)
- [ ] Password is hashed with bcrypt
- [ ] Authentication returns JWT with correct claims
- [ ] Invalid credentials are rejected
- [ ] Refresh token issues new access token
- [ ] Invalid/expired refresh token is rejected
- [ ] Current user endpoint returns correct profile
- [ ] Admin can add user to tenant with role
- [ ] Admin can remove user from tenant
- [ ] Admin can list tenant members
- [ ] Non-admin cannot manage tenant members
- [ ] Audit log created for registration
- [ ] Audit log created for login success
- [ ] Audit log created for login failure
- [ ] Audit log created for member added/removed

---

## Phase 5: Clinical Core

- [ ] Patient can create a case
- [ ] Case is linked to tenant and patient
- [ ] Patient can add symptoms to a case
- [ ] Admin can assign a doctor to a case
- [ ] Case status transitions correctly (Open → Assigned)
- [ ] Doctor can write a diagnosis
- [ ] Diagnosis is linked to the assigned doctor
- [ ] Case status transitions to Diagnosed
- [ ] Admin can close a diagnosed case
- [ ] Case status transitions to Closed
- [ ] Patient can list own cases
- [ ] Doctor can list assigned cases
- [ ] Admin can list all cases for tenant
- [ ] Cases can be filtered by status
- [ ] List endpoints return paginated results with pagination metadata
- [ ] Cross-tenant access is blocked
- [ ] Domain events are published for each mutation
- [ ] Audit logs created for all clinical events

---

## Phase 6: Imaging

- [ ] Pre-signed upload URL is generated
- [ ] Upload URL is scoped to tenant path
- [ ] Upload URL expires after defined period
- [ ] Image metadata is recorded after upload
- [ ] Image is linked to correct case and tenant
- [ ] Images can be listed by case
- [ ] Pre-signed download URL is generated
- [ ] Cross-tenant image access is blocked
- [ ] Audit logs created for imaging events

---

## Phase 7: Frontend

- [ ] Next.js App Router project initialized
- [ ] pnpm configured as package manager
- [ ] Static export configured (`output: 'export'` in next.config.js)
- [ ] TypeScript strict mode enabled
- [ ] `openapi-typescript` generates types from `spec/openapi.yaml`
- [ ] Feature-based directory structure created (`features/`, `infrastructure/`, `shared/`, `generated/`)
- [ ] Infrastructure layer configured (openapi-fetch instance, TanStack Query client, auth helpers)
- [ ] Shared components created (layouts, navigation, base UI)
- [ ] TanStack Query installed and configured
- [ ] `openapi-fetch` installed and configured
- [ ] React Hook Form + Zod installed and configured
- [ ] Tailwind CSS installed and configured
- [ ] shadcn/ui installed with base components
- [ ] All components use `'use client'` directive
- [ ] No API Routes present
- [ ] No Server Actions present
- [ ] No SSR (static export only)
- [ ] Authentication feature complete (login, register — components, hooks, services, schemas)
- [ ] Patients feature complete (dashboard, case list)
- [ ] Doctors feature complete (assigned cases, diagnosis)
- [ ] Admin feature complete (case management, audit logs)
- [ ] Case creation form submits correctly (React Hook Form + Zod)
- [ ] Symptom form adds symptoms to a case
- [ ] Image upload works with pre-signed URL
- [ ] Diagnosis is viewable by authorized roles
- [ ] Audit log viewer is accessible to admins
- [ ] Tenant switcher component allows switching between tenants
- [ ] Tenant switcher calls select-tenant API and updates JWT
- [ ] Unauthorized access redirects to login
- [ ] JWT tokens are stored in httpOnly cookies
- [ ] No business logic in frontend components
- [ ] Each feature is self-contained (components, hooks, services, schemas, types)
- [ ] No unnecessary coupling between features
- [ ] Vitest configured with `@testing-library/react` and MSW
- [ ] Component tests pass with `@testing-library/react` + `@testing-library/user-event`
- [ ] API mocking works with MSW
- [ ] Coverage reporting works with `@vitest/coverage-v8`

---

## Phase 8: Polish

- [ ] No hardcoded secrets in codebase
- [ ] All inputs validated
- [ ] CloudWatch dashboard created
- [ ] CloudTrail enabled
- [ ] README includes deployment instructions
- [ ] Architecture diagrams generated
- [ ] All TODO comments resolved or documented
- [ ] No PHI in any log output
- [ ] HTTPS enforced on all connections

---

## Phase 9: CI/CD

- [ ] Infrastructure pipeline runs Terraform init/validate/plan/apply
- [ ] Infrastructure pipeline requires manual approval for production
- [ ] Backend pipeline runs format, lint, unit tests, build
- [ ] Backend pipeline builds and pushes Docker image to ECR
- [ ] Backend pipeline runs database migrations before deployment
- [ ] Backend pipeline deploys to ECS and validates health check
- [ ] Frontend pipeline runs format, typecheck, unit tests, integration tests
- [ ] Frontend pipeline builds and exports static assets
- [ ] Frontend pipeline uploads to S3 and invalidates CloudFront
- [ ] GitHub OIDC configured for AWS authentication (no long-lived credentials)
- [ ] Path-based triggers configured (infrastructure/, backend/, frontend/)
- [ ] Concurrency groups prevent parallel deployments of same component
- [ ] Each pipeline supports independent rollback
- [ ] CI/CD strategy documented in CI_CD_STRATEGY.md

---

## Phase 10: HIPAA Compliance

### Privacy Rule (45 CFR §164.500–534)
- [ ] Notice of Privacy Practices (NPP) defined and available to patients
- [ ] Patient right to access PHI implemented (view own data via API)
- [ ] Patient right to amend PHI implemented (request amendment workflow)
- [ ] Patient right to accounting of disclosures implemented (audit log report)
- [ ] Patient right to request restrictions implemented
- [ ] Patient right to confidential communications implemented
- [ ] Minimum Necessary Standard enforced per role (Patient, Doctor, Admin)
- [ ] Business Associate Agreement (BAA) signed with AWS
- [ ] Uses and Disclosures policy documented (TPO, required by law, authorization)
- [ ] De-identification methods documented (Safe Harbor, Expert Determination)

### Breach Notification Rule (45 CFR §164.400–414)
- [ ] Breach definition documented
- [ ] Breach assessment process defined (contain, assess, document, notify)
- [ ] Risk assessment factors documented
- [ ] Individual notification process defined (within 60 days)
- [ ] HHS notification process defined (≥500 individuals: 60 days; <500: annual)
- [ ] Media notification process defined (≥500 in a state: 60 days)
- [ ] Breach response team identified (Security Officer, Privacy Officer, Legal, IT, Communications)
- [ ] Breach documentation template defined
- [ ] Breach records retained for 6 years

### Administrative Safeguards (45 CFR §164.308)
- [ ] Security Officer designated and documented
- [ ] Privacy Officer designated and documented
- [ ] Risk analysis methodology documented
- [ ] Risk management process defined
- [ ] Sanction policy documented (disciplinary actions for violations)
- [ ] Information system activity review process defined (weekly audit log review)
- [ ] Workforce security procedures documented (background checks, termination)
- [ ] Access authorization process defined (RBAC, minimum necessary)
- [ ] Security awareness training program documented
- [ ] Security incident response procedures documented
- [ ] Contingency plan documented:
  - [ ] Data backup plan (RDS daily, S3 versioning)
  - [ ] Disaster recovery plan (RTO: 4 hours, RPO: 1 hour)
  - [ ] Emergency mode operation plan
  - [ ] Testing and revision procedures (annual)
  - [ ] Applications and data criticality analysis
- [ ] Annual security evaluation scheduled

### Physical Safeguards (45 CFR §164.310)
- [ ] Facility access controls documented (AWS managed)
- [ ] Workstation security policy documented (endpoint protection, encryption)
- [ ] Workstation use policy documented (screen lock, secure areas)
- [ ] Device and media controls documented (disposal, re-use, accountability)
- [ ] Automatic logoff configured (15 minutes inactivity)

### Technical Safeguards (45 CFR §164.312)
- [x] Unique user identification (JWT user_id claim)
- [x] Emergency access procedure documented
- [x] Automatic logoff implemented (15 minutes)
- [x] Audit controls implemented (structured logging)
- [x] Integrity controls implemented (referential integrity, validation)
- [x] Person/entity authentication implemented (JWT)
- [x] Transmission security implemented (TLS 1.2+)
- [x] Encryption at rest implemented (AES-256)

### Documentation Retention (45 CFR §164.530(j))
- [ ] Audit logs retained for 6 years
- [ ] Breach documentation retained for 6 years
- [ ] Security incident records retained for 6 years
- [ ] Policy documentation retained for 6 years
- [ ] Retention policy documented and enforced via S3 lifecycle
