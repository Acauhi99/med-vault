MedVault frontend.

## Commands

```bash
pnpm dev
pnpm build
pnpm lint
pnpm test
```

## Environment

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## Structure

- `app/` route shell and layout
- `features/authentication/` auth UI, schema, and service layer
- `infrastructure/` API client, query provider, and session store

## Validation

- Zod v4 helpers in use: `z.email()`, `z.uuid()`, `z.iso.datetime()`
- Check the official Zod docs before adding validators; deprecated helpers still exist in typings
