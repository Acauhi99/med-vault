# Local Quality Gates

Validation philosophy, quality gates, tooling, and execution strategy for MedVault.

---

## Quality Philosophy

Software quality should be enforced as early as possible. Every code change should pass a series of local quality gates before reaching the CI pipeline.

```
Developer
    ↓
IDE Formatting
    ↓
Git Pre-Commit
    ↓
Git Pre-Push
    ↓
Continuous Integration
    ↓
Deployment
```

Each stage exists to catch a different category of problems. The earlier a problem is detected, the cheaper it is to fix.

---

## Quality Principles

- **Fail Fast** — stop the pipeline at the first error
- **Shift Left** — move validation as close to the developer as possible
- **Automation over Manual Verification** — never rely on human discipline for correctness
- **Consistency over Cleverness** — predictable tooling beats clever shortcuts
- **Reproducibility** — every validation produces the same result regardless of environment
- **Deterministic Builds** — same input always produces same output
- **Local Validation before CI** — developers should rarely experience CI failures caused by locally detectable problems
- **Keep Feedback Loops Short** — fast feedback enables fast iteration

The goal is to prevent trivial issues from consuming CI resources.

---

## Validation Layers

### IDE

**Responsibility:** Automatic formatting, import organization, immediate developer feedback.

The IDE should solve formatting whenever possible. Developers should not think about formatting — the IDE handles it.

### Pre-Commit

**Responsibility:** Prevent commits containing obvious problems.

- Execute only fast validations
- Expected execution time: under 10 seconds
- No long-running tests
- No network calls
- Formatting, static analysis, quick linting

### Pre-Push

**Responsibility:** Ensure pushed code is healthy enough to reach CI.

- Functional validations
- Type checking
- Unit tests
- Integration tests
- Infrastructure validation

Expected execution time may be longer. Developers should rarely experience CI failures caused by problems that could have been detected locally.

### Continuous Integration

**Responsibility:** Final validation gate.

- Execute all local validations again inside a clean environment
- May execute heavier validations intentionally excluded from local hooks
- The CI pipeline should never rely on validations that developers cannot execute locally

---

## Frontend Quality Gates

### Tooling

| Tool | Purpose |
|------|---------|
| Biome | Formatting and linting |
| TypeScript | Type checking |
| Vitest | Unit and integration testing |
| `@testing-library/react` | Component testing utilities |
| MSW (Mock Service Worker) | API mocking |
| `@vitest/coverage-v8` | Coverage reporting |

### Pre-Commit

| Check | Tool | Why |
|-------|------|-----|
| Format | `biome format --write` | Consistent code style |
| Lint | `biome lint` | Catch common errors |

These operations should remain extremely fast.

### Pre-Push

| Check | Tool | Why |
|-------|------|-----|
| Format + Lint | `biome check` | Full formatting and linting pass |
| Type Check | `tsc --noEmit` | Catch type errors |
| Unit Tests | `vitest run` | Validate business logic |
| Integration Tests | `vitest run` (integration) | Validate component + service integration |

The objective is ensuring that every pushed frontend change is production-ready.

---

## Backend Quality Gates

### Tooling

| Tool | Purpose |
|------|---------|
| `gofumpt` | Code formatting (stricter than `gofmt`) |
| `golangci-lint` | Static analysis (50+ linters) |
| `go test` | Unit and integration testing |
| `govulncheck` | Security vulnerability scanning (recommended for CI, optional locally) |

### Pre-Commit

| Check | Tool | Why |
|-------|------|-----|
| Format | `gofumpt -w` | Consistent Go formatting |
| Lint | `golangci-lint run` | Catch common errors and style issues |

### Pre-Push

| Check | Tool | Why |
|-------|------|-----|
| Full Lint | `golangci-lint run` | Comprehensive static analysis |
| Tests | `go test ./...` | Validate all packages |

Coverage collection may be performed by CI to avoid slowing local workflow.

---

## Infrastructure Quality Gates

### Tooling

| Tool | Purpose |
|------|---------|
| `terraform fmt` | Official Terraform formatting |
| `terraform validate` | Configuration correctness |
| `tflint` | Terraform and AWS provider best-practice issues |
| Checkov | Infrastructure as Code security analysis and policy validation |

**Why Checkov:** Provides broader Infrastructure-as-Code security analysis than alternatives, aligns with the project's HIPAA-aware architecture, detects misconfigurations across Terraform resources.

### Pre-Commit

| Check | Tool | Why |
|-------|------|-----|
| Format | `terraform fmt -recursive` | Consistent Terraform formatting |
| Validate | `terraform validate` | Configuration correctness |

These validations should remain lightweight.

### Pre-Push

| Check | Tool | Why |
|-------|------|-----|
| Format Check | `terraform fmt -check -recursive` | Ensure formatting is applied |
| Validate | `terraform validate` | Configuration correctness |
| Lint | `tflint` | Best-practice compliance |
| Security | `checkov -d .` | Security and policy validation |

Infrastructure should never be pushed before passing all validation gates.

---

## Documentation Quality Gates

### Philosophy

Documentation is not a byproduct — it is a first-class deliverable. Every structural change must update the corresponding doc(s) in the same change set. Docs and code are never out of sync.

### What Triggers a Documentation Update

| Change Type | Documentation to Update |
|-------------|------------------------|
| API endpoint added/removed/changed | `spec/openapi.yaml` + `docs/REQUIREMENTS.md` |
| New aggregate, entity, or value object | `docs/DOMAIN.md` |
| New bounded context or module | `docs/DOMAIN.md` + `docs/ARCHITECTURE.md` |
| New domain event | `docs/DOMAIN.md` + `docs/diagrams/domain-events-flow.md` |
| New or changed status transition | `docs/diagrams/case-lifecycle.md` |
| New infrastructure resource | `docs/INFRASTRUCTURE.md` |
| New security control | `docs/SECURITY.md` |
| New or changed user flow | Related diagram in `docs/diagrams/` |
| New CI/CD pipeline change | `docs/CI_CD_STRATEGY.md` |
| New ADR decision | New file in `docs/adr/` |
| HIPAA-relevant control change | `docs/SECURITY.md` (HIPAA sections) |

### Pre-Push Documentation Check

Before pushing code that changes architecture or flows:

1. **Read the relevant doc(s)** — understand the current documented contract
2. **Update docs in the same commit** — docs and code are never separate commits
3. **Verify diagrams match code** — Mermaid diagrams must reflect actual behavior
4. **Verify cross-references** — links between docs must not be broken

### Documentation Anti-Patterns

- ❌ Code changes without reading the relevant doc
- ❌ Structural changes without updating the doc
- ❌ Doc describes behavior that doesn't match the code
- ❌ Diagram shows a flow that the code doesn't implement
- ❌ New API endpoint not in `spec/openapi.yaml`
- ❌ New domain event not in `docs/DOMAIN.md`
- ❌ Broken internal links between docs

### Source of Truth Hierarchy

When docs conflict:

1. `spec/openapi.yaml` — API contract (highest authority)
2. `docs/SECURITY.md` — security and compliance
3. `docs/DOMAIN.md` — domain model
4. `docs/ARCHITECTURE.md` — system architecture
5. `docs/diagrams/` — flow diagrams
6. Other docs — supporting documentation

---

## Unified Task Runner

Every validation is executed through a single **Taskfile**. The Taskfile becomes the project's official developer interface.

Developers should not memorize technology-specific commands. Every operation is exposed through standardized tasks.

### Standard Tasks

| Task | Description |
|------|-------------|
| `task format` | Format all code (Biome, gofumpt, terraform fmt) |
| `task lint` | Lint all code (Biome, golangci-lint, tflint) |
| `task validate` | Validate all configurations (TypeScript, terraform validate) |
| `task test` | Run all tests (Vitest, go test) |
| `task build` | Build all artifacts |
| `task pre-commit` | Run all pre-commit checks |
| `task pre-push` | Run all pre-push checks |

This abstraction provides a consistent developer experience across frontend, backend, and infrastructure.

### Git Hooks Integration

Git hooks should invoke Taskfile commands rather than embedding validation logic directly. This keeps hooks simple and validation logic centralized.

---

## Git Hooks Philosophy

### Pre-Commit

- Extremely fast (under 10 seconds)
- Formatting
- Static analysis
- No long-running tests
- No network calls

### Pre-Push

- Functional validation
- Type checking
- Unit tests
- Integration tests
- Infrastructure validation

Git hooks should remain simple orchestration layers. Business logic must never exist inside hook scripts.

---

## AI Development

AI agents should always execute standardized Taskfile commands instead of technology-specific commands whenever possible.

**Instead of:**

```bash
pnpm lint
go test ./...
terraform validate
```

**Prefer:**

```bash
task lint
task test
task validate
```

This guarantees consistency regardless of the underlying implementation.

---

## Future Evolution

The following are documented as future improvements, not current implementation requirements.

| Improvement | Purpose |
|-------------|---------|
| Secret scanning | Prevent secrets from reaching the repository |
| License compliance | Validate dependency licenses |
| Dependency vulnerability scanning | Detect known vulnerabilities in dependencies |
| SBOM generation | Software Bill of Materials for compliance |
| Container image scanning | Detect vulnerabilities in Docker images |
| SAST (Static Application Security Testing) | Detect security vulnerabilities in source code |
| DAST (Dynamic Application Security Testing) | Test running application for vulnerabilities |
| Performance regression tests | Detect performance degradation |
| Contract testing | Validate API contract adherence between services |
| Mutation testing | Validate test suite effectiveness |

---

## Philosophy

The quality gates strategy follows the same architectural principles adopted throughout the project:

- Explicit over implicit
- Simplicity over unnecessary complexity
- Automation over manual verification
- Confidence over coverage
- Early detection over late discovery

The objective is to build a validation pipeline that remains fast, reliable, and capable of providing high confidence throughout the lifetime of the project.
