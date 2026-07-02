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
- AWS infrastructure via Terraform (see [INFRASTRUCTURE.md](INFRASTRUCTURE.md))
- DDD with CQRS backend architecture
- Feature-Based frontend architecture (see [ADR-015](adr/015-frontend-feature-based-architecture.md))
- Design-First API with OpenAPI (see [ADR-016](adr/016-design-first-api-documentation.md))
- Pragmatic testing strategy (see [TESTING_STRATEGY.md](TESTING_STRATEGY.md))
- Local quality gates: pre-commit, pre-push, unified task runner (see [QUALITY_GATES.md](QUALITY_GATES.md))
- Security by default: encryption, private networking, least privilege (see [SECURITY.md](SECURITY.md))

### Out of Scope

- Real patient data
- Production deployment
- CI/CD pipeline (planned for future)
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

- [ARCHITECTURE.md](ARCHITECTURE.md) — system design and component details
- [DOMAIN.md](DOMAIN.md) — domain model and business rules
- [REQUIREMENTS.md](REQUIREMENTS.md) — functional and non-functional requirements
- [SECURITY.md](SECURITY.md) — security controls and threat model
- [INFRASTRUCTURE.md](INFRASTRUCTURE.md) — Terraform philosophy, modules, state, security
- [TESTING_STRATEGY.md](TESTING_STRATEGY.md) — testing philosophy and approach
- [QUALITY_GATES.md](QUALITY_GATES.md) — validation layers, tooling, execution strategy
- [PROJECT_PRINCIPLES.md](PROJECT_PRINCIPLES.md) — engineering principles
