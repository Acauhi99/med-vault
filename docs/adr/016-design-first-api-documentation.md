# ADR-016: Design-First API Documentation with OpenAPI

## Status

Accepted

## Context

MedVault needs a clear API contract between the Go backend and Next.js frontend. The API design must be:
- Consistent across both sides
- Type-safe on both sides
- Documented and versioned
- Generated, not hand-written

Without a single source of truth, frontend and backend teams (or AI agents) drift apart — endpoints, payloads, and error formats diverge silently.

## Decision

Use **Design-First** approach with OpenAPI as the single source of truth for all API contracts.

### Contract Flow

```
OpenAPI Spec (spec/openapi.yaml)
    ├── oapi-codegen → Go server interfaces + types (backend)
    ├── openapi-typescript → TypeScript types (frontend)
    └── openapi-fetch → Type-safe HTTP client (frontend)
```

### Directory Structure

```
med-vault/
├── spec/
│   └── openapi.yaml              # Single source of truth
├── backend/
│   └── internal/generated/       # Generated from OpenAPI (oapi-codegen)
│       ├── server.go
│       └── types.go
└── frontend/
    └── generated/            # Generated from OpenAPI (openapi-typescript)
        └── api.d.ts
```

### Tools

| Tool | Purpose | Output |
|------|---------|--------|
| `oapi-codegen` | Generate Go server interfaces from OpenAPI | `server.go`, `types.go` |
| `openapi-typescript` | Generate TypeScript types from OpenAPI | `api.d.ts` |
| `openapi-fetch` | Type-safe HTTP client for frontend | Runtime dependency |

### API Evolution Workflow

1. Define or update `spec/openapi.yaml`
2. Run code generation scripts (backend + frontend)
3. Backend implements generated interfaces
4. Frontend consumes generated types
5. Both sides are always in sync

## Consequences

### Positive
- Single source of truth for API contracts
- Type-safe on both sides (Go + TypeScript)
- Frontend never manually writes HTTP contracts
- Backend never generates OpenAPI from code
- Contract drift is impossible if generation is used
- Changes to API are explicit (diff on `openapi.yaml`)
- Tests can validate contract adherence

### Negative
- Requires generation step in development workflow
- OpenAPI spec can become verbose for simple APIs
- Learning curve for OpenAPI syntax

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Backend generates OpenAPI | Loses contract-first discipline; frontend still hand-writes types |
| GraphQL | Overkill for PoC, adds complexity |
| Manual HTTP contracts | Prone to drift, not type-safe |
| Swagger Codegen | Less modern than oapi-codegen for Go |

## References

- [OpenAPI Specification](https://spec.openapis.org/oas/latest.html)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
- [openapi-typescript](https://openapi-ts.dev/)
- [openapi-fetch](https://openapi-ts.dev/openapi-fetch/)
