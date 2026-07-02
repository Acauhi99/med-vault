# Checklist

Acceptance criteria for each phase of MedVault. Mark each item as complete when verified.

---

## Phase 1: Foundation

- [ ] Go module initialized (`go.mod` exists)
- [ ] Next.js App Router project initialized (static export)
- [ ] Terraform project initialized (`main.tf` exists)
- [ ] Makefile or justfile with common commands
- [ ] golangci-lint configured
- [ ] ESLint configured
- [ ] goimports configured
- [ ] Prettier configured
- [ ] All documentation files present and non-empty
- [ ] All ADRs present and non-empty

---

## Phase 2: Infrastructure

- [ ] VPC created with public and private subnets
- [ ] NAT Gateway configured for private subnet internet access
- [ ] RDS PostgreSQL in private subnet
- [ ] RDS encryption enabled
- [ ] S3 bucket for medical images created
- [ ] S3 bucket encryption enabled
- [ ] S3 bucket public access blocked
- [ ] ECS Fargate cluster created
- [ ] ALB created in public subnet
- [ ] ALB TLS certificate configured
- [ ] Security groups configured (ALB → ECS → RDS)
- [ ] IAM roles created (ECS task role, execution role)
- [ ] Secrets Manager secret created for DB credentials
- [ ] CloudWatch log group created
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
- [ ] `net/http` server starts with `http.ServeMux` routing
- [ ] JWT middleware validates tokens
- [ ] Tenant context middleware extracts tenant_id
- [ ] RBAC middleware enforces role-based access
- [ ] Repository interfaces defined in domain layer
- [ ] Repository implementations in infrastructure layer
- [ ] Error responses follow standard format
- [ ] `log/slog` structured JSON logging works
- [ ] Health check endpoint returns 200
- [ ] Unit tests pass with `testing` + `httptest`

---

## Phase 4: Identity & Access

- [ ] Tenant can be created
- [ ] User can be registered with valid email
- [ ] Duplicate email within tenant is rejected
- [ ] Password is hashed with bcrypt
- [ ] Authentication returns JWT with correct claims
- [ ] Invalid credentials are rejected
- [ ] Refresh token issues new access token
- [ ] Invalid/expired refresh token is rejected
- [ ] Current user endpoint returns correct profile
- [ ] Audit log created for registration
- [ ] Audit log created for login success
- [ ] Audit log created for login failure

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
- [ ] Feature-based directory structure created (`features/`, `infrastructure/`, `shared/`)
- [ ] Infrastructure layer configured (Axios instance, TanStack Query client, auth helpers)
- [ ] Shared components created (layouts, navigation, base UI)
- [ ] TanStack Query installed and configured
- [ ] Axios installed with typed API client
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
- [ ] Unauthorized access redirects to login
- [ ] JWT tokens are stored in httpOnly cookies
- [ ] No business logic in frontend components
- [ ] Each feature is self-contained (components, hooks, services, schemas, types)
- [ ] No unnecessary coupling between features

---

## Phase 8: Polish

- [ ] No hardcoded secrets in codebase
- [ ] All inputs validated
- [ ] Rate limiting configured on ALB
- [ ] CloudWatch dashboard created
- [ ] CloudTrail enabled
- [ ] README includes deployment instructions
- [ ] Architecture diagrams generated
- [ ] All TODO comments resolved or documented
- [ ] No PHI in any log output
- [ ] HTTPS enforced on all connections
