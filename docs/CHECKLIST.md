# Checklist

Acceptance criteria for each phase of MedVault. Mark each item as complete when verified.

**Testing philosophy:** Unit tests and integration tests only. See [TESTING_STRATEGY.md](TESTING_STRATEGY.md).

**Quality gates:** Pre-commit, pre-push, unified Taskfile. See [QUALITY_GATES.md](QUALITY_GATES.md).

**Phased delivery plan:** See [ROADMAP.md](../ROADMAP.md) for the overall implementation timeline.

---

## Phase 0: AWS Account Setup

- [x] AWS CLI v2 installed
- [x] Root account MFA enabled
- [x] Admin IAM user created (`medvault-admin`) with MFA
- [x] IAM groups created (`medvault-admins`, `medvault-devs`, `medvault-terraform`)
- [x] Terraform state backend (S3 bucket + versioning + encryption + public access blocked)
- [x] DynamoDB table for state locking
- [x] GitHub OIDC identity provider created
- [x] GitHub Actions IAM role (`medvault-github-actions`) with scoped policy
- [x] CloudTrail enabled (multi-region, management events)
- [x] AWS Config enabled (all resources, continuous recording)
- [x] ECR repository created (`medvault/backend`, scan on push, AES256)
- [x] Billing budget alerts configured

---

## Phase 1: Foundation

- [x] Go module initialized (`go.mod` exists)
- [x] Next.js App Router project initialized (static export)
- [x] Terraform project initialized (`main.tf` exists)
- [x] OpenAPI 3.1.3 contract defined (`spec/openapi.yaml`)
- [x] `oapi-codegen` configured for backend generation
- [x] `openapi-typescript` configured for frontend generation
- [x] Taskfile created with `format`, `lint`, `validate`, `test`, `pre-commit`, `pre-push` tasks
- [x] `gofumpt` configured (backend formatting)
- [x] `golangci-lint` configured (backend linting)
- [x] Biome configured (frontend formatting + linting)
- [x] `tflint` configured (infrastructure linting)
- [x] Checkov configured (infrastructure security)
- [x] Git pre-commit hook configured
- [x] Git pre-push hook configured
- [x] All documentation files present and non-empty
- [x] All ADRs present and non-empty

---

## Phase 2: Infrastructure

- [x] Module structure follows capability-based design (see [INFRASTRUCTURE.md](INFRASTRUCTURE.md))
- [x] `network` module: VPC, public/private subnets, optional NAT, route tables, internet gateway
- [x] `database` module: RDS PostgreSQL in private subnet, encryption enabled
- [x] `storage` module: S3 buckets with encryption, versioning, lifecycle rules
- [x] `application` module: ECS Fargate cluster, ALB, task definition, security groups
- [x] `security` module: IAM roles, policies, KMS, Secrets Manager
- [x] `observability` module: CloudWatch, CloudTrail, VPC Flow Logs, AWS Config
- [x] Production environment composes modules correctly
- [x] Remote state in S3 with versioning and encryption
- [x] Security groups configured (ALB → ECS → RDS)
- [x] S3 bucket public access blocked
- [x] WAF associated with ALB

---

## Phase 3: Backend Foundation

- [x] Project structure follows DDD layers (domain, application, infrastructure)
- [x] Domain layer has no external dependencies
- [x] Aggregates defined per bounded context
- [x] Value objects are immutable
- [x] `envconfig` loads configuration from environment variables
- [x] `pgx` connection pool configured and tested
- [x] `golang-migrate` migrations run successfully
- [ ] `sqlc` generates type-safe query code (N/A — using raw pgx)
- [x] `oapi-codegen` generates Go server interfaces from OpenAPI
- [x] `net/http` server starts with `http.ServeMux` routing
- [x] JWT middleware validates tokens
- [x] Tenant context middleware extracts tenant_id
- [x] RBAC middleware enforces role-based access
- [x] Repository interfaces defined in domain layer
- [x] Repository implementations in infrastructure layer
- [x] Error responses follow standard format
- [x] `log/slog` structured JSON logging works
- [x] Health check endpoint returns 200
- [x] Rate limiting middleware on auth endpoints (login, register, refresh)
- [ ] Unit tests pass with `testing` + `testify/assert` (N/A — using stdlib testing)
- [x] HTTP handler tests pass with `httptest`
- [ ] Struct comparison tests pass with `go-cmp` (N/A — using manual assertions)
- [ ] Integration tests pass with `testcontainers-go` (TODO — requires testcontainers)
- [x] Coverage reporting works with `go test -cover`
- [x] Dockerfile created with multi-stage build (see [ADR-019](adr/019-docker-image-strategy.md))
- [x] Build stage: Go toolchain, module download, compilation, validation
- [x] Runtime stage: distroless, binary + CA certs only
- [x] CGO disabled, static binary
- [x] Non-root user, read-only binary
- [x] `.dockerignore` excludes git, IDE, caches, docs, test artifacts
- [x] BuildKit cache mounts for Go modules and compiler
- [x] No Go compiler, source code, or build tools in production image
- [x] Image builds successfully with `docker build`

---

## Phase 4: Identity & Access

- [x] Tenant can be created
- [x] Suspended tenant can be reactivated
- [x] User can be registered with valid email
- [x] Email is unique system-wide (duplicate registration rejected)
- [x] Password is hashed with bcrypt
- [x] Authentication returns JWT with correct claims
- [x] Invalid credentials are rejected
- [x] Refresh token issues new access token
- [x] Invalid/expired refresh token is rejected
- [x] Current user endpoint returns correct profile
- [x] Admin can add user to tenant with role
- [x] Admin can remove user from tenant
- [x] Admin can list tenant members
- [x] Non-admin cannot manage tenant members
- [x] Audit log created for registration
- [x] Audit log created for login success
- [x] Login failure emitted as structured security log
- [x] Audit log created for member added/removed

---

## Phase 5: Clinical Core

- [x] Patient can create a case
- [x] Case is linked to tenant and patient
- [x] Patient can add symptoms to a case
- [x] Admin can assign a doctor to a case
- [x] Case status transitions correctly (Open → Assigned)
- [x] Doctor can write a diagnosis
- [x] Diagnosis is linked to the assigned doctor
- [x] Case status transitions to Diagnosed
- [x] Admin can close a diagnosed case
- [x] Case status transitions to Closed
- [x] Patient can list own cases
- [x] Doctor can list assigned cases
- [x] Admin can list all cases for tenant
- [x] Cases can be filtered by status
- [x] List endpoints return paginated results with pagination metadata
- [x] Cross-tenant access is blocked
- [x] Domain events are published for each mutation
- [x] Audit logs created for all clinical events

---

## Phase 6: Imaging

- [x] Pre-signed upload URL is generated
- [x] Upload URL is scoped to tenant path
- [x] Upload URL expires after defined period
- [x] Image metadata is recorded after upload
- [x] Image is linked to correct case and tenant
- [x] Images can be listed by case
- [x] Pre-signed download URL is generated
- [x] Cross-tenant image access is blocked
- [x] Audit logs created for imaging events

---

## Phase 7: Audit

- [x] Audit log listing (admin-only)
- [x] Audit log writing infrastructure
- [x] Tenant isolation on audit queries

---

## Phase 8: Frontend

- [x] Next.js App Router project initialized
- [x] pnpm configured as package manager
- [x] Static export configured (`output: 'export'` in next.config.js)
- [x] TypeScript strict mode enabled
- [x] `openapi-typescript` generates types from `spec/openapi.yaml`
- [x] Feature-based directory structure created (`features/`, `infrastructure/`, `shared/`, `generated/`)
- [x] Infrastructure layer configured (openapi-fetch instance, TanStack Query client, auth helpers)
- [x] Shared components created (layouts, navigation, base UI)
- [x] TanStack Query installed and configured
- [x] `openapi-fetch` installed and configured
- [x] React Hook Form + Zod installed and configured
- [x] Tailwind CSS installed and configured
- [ ] shadcn/ui installed with base components
- [ ] All components use `'use client'` directive
- [x] No API Routes present
- [x] No Server Actions present
- [x] No SSR (static export only)
- [x] Authentication feature complete (login, register — components, hooks, services, schemas)
- [x] Patients feature complete (dashboard, case list)
- [x] Doctors feature complete (assigned cases, diagnosis)
- [x] Admin feature complete (case management, audit logs)
- [x] Case creation form submits correctly (React Hook Form + Zod)
- [x] Symptom form adds symptoms to a case
- [x] Image upload works with pre-signed URL
- [x] Diagnosis is viewable by authorized roles
- [x] Audit log viewer is accessible to admins
- [x] Tenant switcher component allows switching between tenants
- [x] Tenant switcher calls select-tenant API and updates JWT
- [ ] Unauthorized access redirects to login
- [x] JWT tokens are stored in client session state
- [x] No business logic in frontend components
- [x] Each feature is self-contained (components, hooks, services, schemas, types)
- [x] No unnecessary coupling between features
- [x] Vitest configured with `@testing-library/react` and MSW
- [x] Component tests pass with `@testing-library/react` + `@testing-library/user-event`
- [x] API mocking works with MSW
- [ ] Coverage reporting works with `@vitest/coverage-v8`

---

## Phase 9: Polish

- [ ] No hardcoded secrets in codebase
- [ ] All inputs validated
- [ ] CloudWatch dashboard created
- [x] CloudTrail enabled
- [ ] README includes deployment instructions
- [ ] Architecture diagrams generated
- [ ] All TODO comments resolved or documented
- [ ] No PHI in any log output
- [x] HTTPS enforced on all connections

---

## Phase 10: CI/CD

- [x] Infrastructure pipeline runs Terraform init/validate/plan/apply
- [ ] Infrastructure pipeline requires manual approval for production
- [ ] Backend pipeline runs format, lint, unit tests, build
- [x] Backend pipeline builds and pushes Docker image to ECR
- [x] Backend pipeline runs database migrations before deployment
- [x] Backend pipeline deploys to ECS and validates health check
- [x] Frontend pipeline runs format, typecheck, unit tests, integration tests
- [x] Frontend pipeline builds and exports static assets
- [x] Frontend pipeline builds Docker image, pushes to ECR, and deploys to ECS
- [x] GitHub OIDC configured for AWS authentication (no long-lived credentials)
- [x] Path-based triggers configured (infrastructure/, backend/, frontend/)
- [x] Concurrency groups prevent parallel deployments of same component
- [x] Each pipeline supports independent rollback
- [x] CI/CD strategy documented in CI_CD_STRATEGY.md

---

## Phase 11: HIPAA Compliance

### Privacy Rule (45 CFR §164.500–534)
- [x] Notice of Privacy Practices (NPP) defined and available to patients
- [x] Patient right to access PHI implemented (view own data via API)
- [ ] Patient right to amend PHI implemented (request amendment workflow)
- [x] Patient right to accounting of disclosures implemented (audit log report)
- [ ] Patient right to request restrictions implemented
- [ ] Patient right to confidential communications implemented
- [x] Minimum Necessary Standard enforced per role (Patient, Doctor, Admin)
- [x] Business Associate Agreement (BAA) signed with AWS
- [x] Uses and Disclosures policy documented (TPO, required by law, authorization)
- [x] De-identification methods documented (Safe Harbor, Expert Determination)

### Breach Notification Rule (45 CFR §164.400–414)
- [x] Breach definition documented
- [x] Breach assessment process defined (contain, assess, document, notify)
- [x] Risk assessment factors documented
- [x] Individual notification process defined (within 60 days)
- [x] HHS notification process defined (≥500 individuals: 60 days; <500: annual)
- [x] Media notification process defined (≥500 in a state: 60 days)
- [x] Breach response team identified (Security Officer, Privacy Officer, Legal, IT, Communications)
- [x] Breach documentation template defined
- [x] Breach records retained for 6 years

### Administrative Safeguards (45 CFR §164.308)
- [x] Security Officer designated and documented
- [x] Privacy Officer designated and documented
- [x] Risk analysis methodology documented
- [x] Risk management process defined
- [x] Sanction policy documented (disciplinary actions for violations)
- [x] Information system activity review process defined (weekly audit log review)
- [x] Workforce security procedures documented (background checks, termination)
- [x] Access authorization process defined (RBAC, minimum necessary)
- [x] Security awareness training program documented
- [x] Security incident response procedures documented
- [x] Contingency plan documented:
  - [x] Data backup plan (RDS daily, S3 versioning)
  - [x] Disaster recovery plan (RTO: 4 hours, RPO: 1 hour)
  - [x] Emergency mode operation plan
  - [x] Testing and revision procedures (annual)
  - [x] Applications and data criticality analysis
- [x] Annual security evaluation scheduled

### Physical Safeguards (45 CFR §164.310)
- [x] Facility access controls documented (AWS managed)
- [x] Workstation security policy documented (endpoint protection, encryption)
- [x] Workstation use policy documented (screen lock, secure areas)
- [x] Device and media controls documented (disposal, re-use, accountability)
- [x] Automatic logoff configured (15 minutes inactivity)

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
- [x] Audit logs retained for 6 years
- [x] Breach documentation retained for 6 years
- [x] Security incident records retained for 6 years
- [x] Policy documentation retained for 6 years
- [x] Retention policy documented and enforced via S3 lifecycle
