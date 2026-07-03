# ADR-019: Docker Image Strategy

## Status

Accepted

## Context

MedVault's backend runs on ECS Fargate and requires a Docker image strategy. The image must be minimal, secure, reproducible, and production-oriented. The Docker image is a deployment artifact, not an environment provisioning tool.

## Decision

Use multi-stage Docker builds with distroless runtime, static Go compilation, and BuildKit caching.

### Multi-Stage Build

The build stage handles Go toolchain, module download, dependency resolution, compilation, and build-time validation. The runtime stage executes only the compiled binary.

The final image never contains: Go compiler, source code, Git metadata, build cache, development tools, package managers, or test files.

### Layering Strategy

```
Copy go.mod
    ↓
Copy go.sum
    ↓
Download modules
    ↓
Copy application source
    ↓
Compile application
    ↓
Create runtime image
```

Dependency resolution occurs before copying application source to maximize Docker layer cache efficiency.

### Static Compilation

- CGO disabled
- Static linking
- Minimal runtime dependencies

### Runtime Image

Distroless base. Contains only: application binary, required CA certificates, required runtime assets.

### Security

- Non-root user
- Read-only application binaries
- Minimal Linux capabilities
- No privileged execution
- No bash, curl, wget, vim, apt, apk, debugging utilities, or build tools

### Image Reproducibility

- Pin base image versions (no `latest` tags)
- Rebuild regularly for upstream security fixes
- Deterministic builds (identical source → identical artifacts)

### Runtime Configuration

No environment-specific configuration embedded. Configuration injected externally via environment variables, AWS Secrets Manager, SSM Parameter Store, or ECS Task Definitions. Same image deploys to dev, staging, and production.

### CI/CD Integration

```
Checkout → Restore Cache → Build → Validation → Production Image → Security Scan → SBOM → Sign → Push → Deploy
```

Images are immutable after publication.

## Consequences

### Positive
- Minimal attack surface (distroless, no shell, no package manager)
- Small image size (~10-20MB for Go static binary)
- Fast builds (BuildKit cache mounts for Go modules and compiler)
- Deterministic and reproducible
- Low CVE exposure
- Same image across all environments
- Fast incremental builds (layer caching)

### Negative
- No debugging tools in production (by design; use ECS Exec or sidecar for debugging)
- Distroless base has no shell (requires adjustment for debugging workflows)

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| Alpine | Larger attack surface, musl libc compatibility issues, unnecessary packages |
| Debian/Ubuntu | Larger image, more CVEs, unnecessary for static Go binary |
| Scratch | Too minimal (no CA certificates, no timezone data); distroless provides these |
| Single-stage build | Larger image, includes Go compiler and source in production |

## References

- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Docker BuildKit](https://docs.docker.com/build/buildkit/)
- [ADR-004: ECS Fargate for Compute](004-ecs-fargate-for-compute.md)
- [ADR-018: ECR for Container Registry](018-ecr-for-container-registry.md)
