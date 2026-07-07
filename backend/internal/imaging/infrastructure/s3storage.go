package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	bucket string
	client *s3.Client
	presign *s3.PresignClient
}

func NewS3Storage(ctx context.Context, bucket, region string) (*S3Storage, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	return newS3Storage(bucket, s3.NewFromConfig(cfg)), nil
}

func newS3Storage(bucket string, client *s3.Client) *S3Storage {
	return &S3Storage{
		bucket:  bucket,
		client:  client,
		presign: s3.NewPresignClient(client),
	}
}

func (s *S3Storage) GenerateUploadURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	out, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func (s *S3Storage) GenerateDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	out, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
