package infrastructure

import (
	"context"
	"fmt"
	"time"
)

type StubStorage struct {
	Bucket string
	Region string
}

func NewStubStorage(bucket, region string) *StubStorage {
	return &StubStorage{Bucket: bucket, Region: region}
}

func (s *StubStorage) GenerateUploadURL(_ context.Context, key, _ string, ttl time.Duration) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s?X-Amz-Algorithm=stub&Expires=%d",
		s.Bucket, s.Region, key, time.Now().Add(ttl).Unix()), nil
}

func (s *StubStorage) GenerateDownloadURL(_ context.Context, key string, ttl time.Duration) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s?X-Amz-Algorithm=stub&Expires=%d",
		s.Bucket, s.Region, key, time.Now().Add(ttl).Unix()), nil
}
