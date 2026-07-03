# Context

Project background, scope, and goals for MedVault.

---

## Why

Healthcare applications require higher security standards than typical web applications. This project demonstrates how AWS managed services can be combined to support HIPAA-oriented architecture — without claiming HIPAA certification.

The code is secondary. Architecture is the primary deliverable.

---

## What

A multi-tenant healthcare platform where:

- **Patients** submit symptoms, upload medical images, view consultation history
- **Doctors** review assigned cases, examine images, write diagnoses
- **Administrators** manage cases, assign doctors, inspect audit logs

The business domain is intentionally simple. The complexity is in infrastructure, security, and multi-tenancy.

---

## Scope

### In Scope

- Multi-tenant architecture with tenant isolation
- JWT authentication and role-based authorization
- Medical case lifecycle (create → assign → diagnose → close)
- Medical image upload via S3 pre-signed URLs
- Audit logging for all state-changing operations
- AWS infrastructure via Terraform (see [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md))
- DDD with CQRS backend architecture
- Transactional Outbox for domain event delivery (see [ADR-017](docs/adr/017-transactional-outbox.md))
- Feature-Based frontend architecture (see [ADR-015](docs/adr/015-frontend-feature-based-architecture.md))
- Design-First API with OpenAPI (see [ADR-016](docs/adr/016-design-first-api-documentation.md))
- Pragmatic testing strategy (see [TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md))
- Local quality gates: pre-commit, pre-push, unified task runner (see [QUALITY_GATES.md](docs/QUALITY_GATES.md))
- Security by default: encryption, private networking, least privilege (see [SECURITY.md](docs/SECURITY.md))
- CI/CD pipeline — three independent pipelines: Infrastructure, Backend, Frontend (Phase 9 — see [CI_CD_STRATEGY.md](docs/CI_CD_STRATEGY.md))

### Out of Scope

- Real patient data
- Production deployment
- Real-time notifications
- Video consultations
- Payment processing
- Mobile applications

---

## Target Audience

This project is evaluated by:

- CEO
- HR Leadership
- Engineering Leadership

The audience assesses architectural thinking more than implementation details.

---

## Related Documents

- [ARCHITECTURE.md](docs/ARCHITECTURE.md) — system design and component details
- [DOMAIN.md](docs/DOMAIN.md) — domain model and business rules
- [REQUIREMENTS.md](docs/REQUIREMENTS.md) — functional and non-functional requirements
- [SECURITY.md](docs/SECURITY.md) — security controls and threat model
- [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md) — Terraform philosophy, modules, state, security
- [TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md) — testing philosophy and approach
- [QUALITY_GATES.md](docs/QUALITY_GATES.md) — validation layers, tooling, execution strategy
- [CI_CD_STRATEGY.md](docs/CI_CD_STRATEGY.md) — CI/CD philosophy, pipeline architecture, deployment boundaries
- [PROJECT_PRINCIPLES.md](docs/PROJECT_PRINCIPLES.md) — engineering principles
