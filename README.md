# MedVault

Cloud-native healthcare platform reference implementation. Demonstrates secure, multi-tenant architecture on AWS using Infrastructure as Code.

> **This is a Proof of Concept.** It demonstrates architectural thinking, not production readiness. Do not deploy with real patient data.

---

## Agent Start

1. Read [`AGENTS.md`](AGENTS.md), then [`CONTEXT.md`](CONTEXT.md).
2. Read the source-of-truth docs for the task.
3. Read the matching diagram in [`docs/diagrams/`](docs/diagrams/).
4. Use the README files as navigation only.

## Official Docs

| Document | Purpose |
|----------|---------|
| [Architecture](docs/ARCHITECTURE.md) | System design, components, data flow, deployment |
| [Domain Model](docs/DOMAIN.md) | DDD bounded contexts, aggregates, CQRS mapping |
| [Requirements](docs/REQUIREMENTS.md) | Functional and non-functional requirements, API spec |
| [Security](docs/SECURITY.md) | Threat model, encryption, auth, audit, HIPAA compliance |
| [Infrastructure](docs/INFRASTRUCTURE.md) | Terraform philosophy, modules, security, state, evolution |
| [Testing Strategy](docs/TESTING_STRATEGY.md) | Testing philosophy, principles, pyramid, coverage approach |
| [Quality Gates](docs/QUALITY_GATES.md) | Validation layers, tooling, pre-commit/pre-push, task runner |
| [Principles](docs/PROJECT_PRINCIPLES.md) | Engineering principles and decision criteria |
| [Context](CONTEXT.md) | Project background, scope, and goals |
| [Roadmap](ROADMAP.md) | Phased delivery plan with status tracking |
| [Checklist](docs/CHECKLIST.md) | Acceptance criteria per phase |
| [Diagrams](docs/diagrams/) | Flow diagrams: auth, patient, doctor, admin, case lifecycle, images, events |

---

## Project Map

| Path | Purpose |
|------|---------|
| `backend/README.md` | Go backend entrypoint |
| `frontend/README.md` | Next.js frontend entrypoint |
| `infrastructure/README.md` | Terraform entrypoint |
| `docs/` | Official architecture, security, domain, and process docs |
| `spec/` | OpenAPI contract |
| `Taskfile.yml` | Validation and task entrypoint |

---

## Roles

| Role | Capabilities |
|------|-------------|
| Patient | Register, submit cases, upload images, view history |
| Doctor | View assigned cases, review images, write diagnoses |
| Administrator | Assign doctors, close cases, inspect audit logs |

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
