# ADR-015: Frontend Feature-Based Architecture

## Status

Accepted

## Context

MedVault frontend needs an internal architecture that maximizes maintainability, scalability, testability, and AI-assisted development. Traditional frontend patterns (MVC, MVP, MVVM) organize code by technical role, which fragments business capabilities across files and makes features hard to reason about.

## Decision

Use a Feature-Based Architecture inspired by Vertical Slice Architecture, DDD, and Clean Architecture principles.

### Architectural Style

This project intentionally does NOT follow traditional frontend patterns:

| Pattern | Reason for Rejection |
|---------|---------------------|
| MVC | Separates by technical role, not business capability |
| MVP | Same fragmentation as MVC |
| MVVM | Overhead of bindings, not suited for static export SPA |

Instead, every business capability is implemented as an isolated feature. Features are organized by domain, not by technical role.

### Why Feature-Based

- **Maintainability:** One feature = one directory. Changes are localized.
- **Scalability:** New features don't touch existing code.
- **Testability:** Features are independently testable.
- **AI-assisted development:** Predictable structure reduces ambiguity for AI agents.
- **Feature isolation:** No accidental coupling between unrelated capabilities.

### Directory Structure

```
frontend/
├── app/                        # Next.js App Router (pages only)
├── features/                   # Business capabilities
│   ├── authentication/
│   │   ├── components/         # Presentation components
│   │   ├── hooks/              # Feature-specific hooks
│   │   ├── services/           # API communication
│   │   ├── schemas/            # Zod validation schemas
│   │   ├── types/              # TypeScript types
│   │   └── index.ts            # Public exports
│   ├── patients/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/
│   │   ├── schemas/
│   │   ├── types/
│   │   └── index.ts
│   ├── doctors/
│   │   └── ...
│   ├── admin/
│   │   └── ...
│   └── shared/                 # Cross-feature reusable elements
│       ├── components/         # Shared UI components
│       ├── hooks/              # Shared hooks
│       ├── lib/                # Shared utilities
│       └── types/              # Shared types
├── infrastructure/             # External integrations
│   ├── api/                    # openapi-fetch instance, interceptors
│   ├── auth/                   # JWT storage, token refresh
│   ├── query/                  # TanStack Query client config
│   └── config/                 # Environment configuration
└── shared/                     # App-wide shared code
    ├── components/             # Layout, navigation, base UI
    ├── layouts/                # Page layouts
    ├── lib/                    # Utilities, helpers
    └── types/                  # Global types
```

### Feature Anatomy

Each feature may contain:

| Directory | Purpose |
|-----------|---------|
| `components/` | Presentation components (receive data via props) |
| `hooks/` | TanStack Query hooks, local state, UI orchestration |
| `services/` | openapi-fetch API calls, endpoint definitions |
| `schemas/` | Zod schemas for request/response/form validation |
| `types/` | TypeScript types specific to this feature |
| `utils/` | Feature-specific utilities (rare) |
| `index.ts` | Public API exports |

Features should be self-contained. Avoid unnecessary coupling between features. Shared code exists only when there is a clear cross-feature need.

### Layer Responsibilities

#### Pages (`app/`)

- Responsible only for page composition and routing
- Orchestrate feature components
- Contain little or no business logic
- Maps routes to feature components

#### Components (`features/*/components/`)

- Responsible only for presentation
- Receive data and callbacks through props
- Do not communicate directly with the backend
- Do not contain business rules
- Use shadcn/ui primitives from `shared/`

#### Hooks (`features/*/hooks/`)

- Orchestrate frontend behavior
- TanStack Query integration (queries, mutations)
- Local state management
- UI state (modals, drawers, form state)
- Interaction orchestration
- Do not perform raw HTTP requests (delegate to services)

#### Services (`features/*/services/`)

- Encapsulate all HTTP communication
- Use openapi-fetch (configured in `infrastructure/api/`)
- API requests, request configuration, response mapping
- Endpoint definitions
- Do not contain business rules

#### Schemas (`features/*/schemas/`)

- Zod schemas for request validation
- Zod schemas for response validation
- Form validation schemas
- Runtime type validation
- Live inside their corresponding feature

#### Infrastructure (`infrastructure/`)

- openapi-fetch instance with interceptors (auth tokens, error handling)
- TanStack Query client configuration
- JWT storage and token refresh logic
- Environment configuration

#### Shared (`shared/`)

- Reusable UI components (layouts, navigation, base elements)
- Page layouts
- Utility functions
- Global TypeScript types
- Should never become a dumping ground

### Request Flow

```
UI (Component)
    ↓ calls hook
Feature Hook
    ↓ uses
TanStack Query
    ↓ calls
Service (openapi-fetch)
    ↓ HTTP request
Go REST API
    ↓ response
TanStack Query Cache
    ↓ triggers
UI Update
```

**Why this separation:**
- Components never touch HTTP → easy to test with mock data
- Services never touch UI state → easy to test with mock HTTP
- Hooks bridge UI and services → orchestrator, testable with both mocked
- TanStack Query handles caching, refetching, optimistic updates → no manual cache management

### Business Logic Boundary

Business rules never belong in the frontend. The frontend only validates user input and provides user experience.

**Backend responsibilities (never in frontend):**

- Authentication and authorization
- Tenant isolation
- Permission checks
- Business rule validation
- Medical workflows
- Diagnosis rules
- Audit decisions

**Frontend responsibilities:**

- User input validation (Zod schemas)
- UI state management
- Navigation
- Presentation
- API communication

### Multi-Tenant Awareness

Every frontend request operates within a tenant context. The frontend is responsible only for propagating tenant context received during authentication (JWT claims). Tenant isolation is enforced by the backend. The frontend must never attempt to enforce tenant security.

### Testing Strategy

Unit tests and integration tests only. No end-to-end tests (see [TESTING_STRATEGY.md](../TESTING_STRATEGY.md) for full philosophy).

| Tool | Purpose |
|------|---------|
| Vitest | Test runner (fast, modern, TypeScript-native) |
| `@testing-library/react` | Component rendering and interaction |
| `@testing-library/user-event` | Real user actions (type, click, navigate) |
| MSW (Mock Service Worker) | API mocking without coupling to HTTP client |
| `@vitest/coverage-v8` | Coverage via V8 engine |

| Layer | Test Type | What to Test |
|-------|-----------|-------------|
| Components | Presentation tests | Renders correctly, responds to props, user interactions |
| Hooks | Behavior tests | State transitions, query integration |
| Services | API communication tests | Correct endpoints, request shapes (via MSW) |
| Pages | Integration tests | Route composition, feature wiring |

**Why this stack:** Behavior-oriented tests (testing-library), decoupled from implementation (MSW), fast execution (Vitest + V8), predictable for AI agents.

Business rules should never require frontend unit tests because they belong to the backend.

### AI Development Guidelines

This repository is designed for AI-assisted development. The frontend architecture is intentionally repetitive and predictable:

- Every new feature follows the same directory structure
- Every feature uses the same patterns (hooks, services, schemas)
- Consistency is more important than clever abstractions
- Reduces architectural ambiguity for AI agents
- Minimal prompting required to infer project patterns

## Consequences

### Positive
- Clear feature boundaries aligned with business capabilities
- Predictable structure for AI-assisted development
- Easy to add new features without touching existing code
- Components are pure presentation → easy to test
- Services are pure HTTP → easy to test
- Hooks are orchestrators → testable with mocks
- Aligns with backend DDD philosophy (features ≈ bounded contexts)

### Negative
- More directories than a flat structure
- Requires discipline to keep features self-contained
- Shared code requires judgment about what belongs where

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| MVC | Separates by technical role, fragments business capabilities |
| MVVM | Overhead of data bindings, not suited for static export |
| Flat structure | No feature isolation, hard to scale |
| Barrel files everywhere | Circular dependency risk |
| Atomic Design | UI-centric, not business-capability-centric |

## References

- [Vertical Slice Architecture](https://jimmybogard.com/vertical-slice-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/solid-cleanarchitecture.html)
- [Feature-Based Architecture](https://www.youtube.com/watch?v=Gx1MuLSEvZQ)
