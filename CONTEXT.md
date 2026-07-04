# Context

Project background, scope, and goals for MedVault.

## Agent Start

Use this order when starting a session:

1. Read `AGENTS.md`, then this file.
2. Read the source-of-truth docs for the task.
3. Read the matching flow diagram in `docs/diagrams/`.
4. Inspect the code path with `codegraph_explore` before editing.
5. Make the smallest change that matches the docs.
6. Update docs and tests in the same change set.
7. Run the repo quality gate before reporting done.

## Source Of Truth

| Concern | Document |
|---------|----------|
| System design | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) |
| Domain model | [docs/DOMAIN.md](docs/DOMAIN.md) |
| Requirements | [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md) |
| Security | [docs/SECURITY.md](docs/SECURITY.md) |
| Infrastructure | [docs/INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md) |
| Testing | [docs/TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md) |
| Quality gates | [docs/QUALITY_GATES.md](docs/QUALITY_GATES.md) |
| Engineering principles | [docs/PROJECT_PRINCIPLES.md](docs/PROJECT_PRINCIPLES.md) |
| CI/CD | [docs/CI_CD_STRATEGY.md](docs/CI_CD_STRATEGY.md) |
| Checklist | [docs/CHECKLIST.md](docs/CHECKLIST.md) |
| ADRs | [docs/adr/](docs/adr/) |

## Project Goal

MedVault is a proof-of-concept multi-tenant healthcare platform. The point is the architecture: security, tenant isolation, bounded contexts, and clear contracts.

## Scope

- In scope: multi-tenancy, JWT auth, RBAC, case lifecycle, image upload, audit logging, AWS IaC, DDD/CQRS backend, feature-based frontend, OpenAPI-first contracts.
- Out of scope: real patient data, production deployment, realtime notifications, video, payments, mobile apps.

## Rules For Agents

- Treat `docs/` as the source of truth.
- Treat `README.md` files as navigation only.
- Never add code that contradicts `docs/ARCHITECTURE.md`, `docs/DOMAIN.md`, or `docs/SECURITY.md`.
- If code and docs disagree, fix the doc or the code in the same change.
