# Frontend

Next.js App Router app for MedVault.

## Official Docs

- [Context](../CONTEXT.md)
- [Architecture](../docs/ARCHITECTURE.md)
- [Security](../docs/SECURITY.md)
- [Testing Strategy](../docs/TESTING_STRATEGY.md)
- [Quality Gates](../docs/QUALITY_GATES.md)
- [Frontend ADR-015](../docs/adr/015-frontend-feature-based-architecture.md)
- [Design-First API ADR-016](../docs/adr/016-design-first-api-documentation.md)

## Commands

```bash
pnpm dev
pnpm build
pnpm lint
pnpm test
```

## Structure

- `app/` route shell and layout
- `features/` feature modules
- `infrastructure/` API client, query provider, and session store

## Current Surface

- `/` auth workspace: login, register, tenant select
- `/cases` patient, doctor, admin case views
- `/cases/new` patient case creation
- `/members` admin tenant members
- tenant reactivation dialog on `/members`
- `/audit` admin audit logs

## Not Yet In UI

- explicit tenant switch outside auth flow
