# ADR-008: S3 for Medical Image Storage

## Status

Accepted

## Context

MedVault needs to store medical images (X-rays, scans, etc.). The storage should be encrypted, scalable, and cost-effective.

## Decision

Use Amazon S3 for medical image storage.

## Consequences

### Positive
- Encrypted at rest (SSE-S3 or SSE-KMS)
- Pre-signed URLs for temporary access
- Scalable and durable
- Cost-effective with lifecycle policies
- Accessed only through backend-issued pre-signed URLs

### Negative
- Requires pre-signed URL generation
- Bucket policy configuration needed

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| EFS | Not suitable for object storage |
| EBS | Block storage, not suitable for images |
| DynamoDB | Not suitable for binary data |

## References

- [S3 Documentation](https://docs.aws.amazon.com/s3/)
- [Pre-signed URLs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/PresignedUrl.html)
