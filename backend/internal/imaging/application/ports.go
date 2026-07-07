package application

import (
	"context"
	"time"
)

type Storage interface {
	GenerateUploadURL(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error)
	GenerateDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
