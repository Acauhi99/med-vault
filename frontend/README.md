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
