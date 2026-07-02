# ADR-002: Next.js App Router with Static Export

## Status

Accepted

## Context

MedVault needs a frontend framework for the SPA. The framework should demonstrate modern web development, align with the company's frontend stack, and remain intentionally simple per project scope.

## Decision

Use Next.js App Router with static export. Client Components only. No API Routes, no Server Actions, no SSR.

### Stack

| Tool | Purpose | Rationale |
|------|---------|-----------|
| Next.js (App Router) | Framework, routing, layouts | Official direction for new projects; nested layouts; file-based routing |
| TypeScript | Type safety | Compile-time error catching; better IDE support |
| pnpm | Package management | Fast, disk-efficient, deterministic installs |
| TanStack Query | Server state management | Caching, background refetching, optimistic updates |
| openapi-fetch | Type-safe HTTP client | Generated from OpenAPI; type-safe requests and responses |
| React Hook Form | Form handling | Performance; minimal re-renders; Zod integration |
| Zod | Schema validation | Type-safe validation; shared schemas between forms and API |
| Tailwind CSS | Styling | Utility-first; no custom CSS; consistent design system |
| shadcn/ui | Component system | Production-quality; repository-owned; Tailwind-native |

### What is intentionally NOT used

| Feature | Reason |
|---------|--------|
| Server Components | No business logic in frontend; all components are interactive |
| Server Actions | Mutations go through Go REST API |
| API Routes | Backend is exclusively Go |
| SSR | Static export only; no server-side rendering |
| ISR | Static export only; no incremental static regeneration |
| BFF (Backend-for-Frontend) | Go backend serves all consumers directly |
| Node.js middleware | Static export; no runtime required in production |

## Consequences

### Positive
- Aligns with company's frontend stack
- Modern routing and layout capabilities (nested layouts, route segments)
- Excellent developer experience (file-based routing, TypeScript, Biome)
- Simplified AWS deployment (static export → S3 → CloudFront)
- Complete separation between presentation and business logic
- No Node.js runtime required in production
- All business logic remains in the Go backend
- shadcn/ui provides production-quality components with minimal implementation effort

### Negative
- Static export limits some Next.js features (no ISR, no middleware at edge)
- Client Components only means no server-side data fetching (acceptable for PoC)

## Deployment Model

```
next build → S3 → CloudFront → Client
```

- Static HTML/CSS/JS exported at build time via `output: 'export'` in `next.config.js`
- S3 serves static assets
- CloudFront provides CDN and TLS termination
- No persistent Node.js server

## Frontend Principles

The frontend is responsible for:
- UI rendering and navigation
- Form handling and client-side validation
- Authentication state management
- API communication with Go backend
- User experience

The frontend is NOT responsible for:
- Business rules
- Authorization logic
- Tenant isolation
- Data transformation beyond presentation

Business rules belong exclusively to the Go backend.

## UI Philosophy

This project intentionally avoids custom UI development. The UI stack prioritizes:
- Consistency across all pages
- Accessibility (WCAG compliance via shadcn/ui)
- Simplicity (Tailwind utilities over custom CSS)
- Maintainability (components in repository, not in node_modules)
- Development speed (shadcn/ui provides production-quality components)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| React (Vite) | Less routing/layout capability than Next.js App Router |
| Next.js Pages Router | App Router is the official direction for new projects |
| Next.js with SSR | Requires persistent Node.js runtime, adds operational complexity |
| Vue.js | Smaller ecosystem than React |
| Angular | Heavier framework, more complex for PoC |
| Svelte | Smaller ecosystem, less mature |

## References

- [Next.js App Router Documentation](https://nextjs.org/docs/app)
- [Next.js Static Export Documentation](https://nextjs.org/docs/app/building-your-application/deploying/static-exports)
- [TanStack Query Documentation](https://tanstack.com/query)
- [shadcn/ui Documentation](https://ui.shadcn.com/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
