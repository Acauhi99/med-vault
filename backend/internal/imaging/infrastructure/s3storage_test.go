package infrastructure

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS3StorageGenerateUploadURL(t *testing.T) {
	storage := newS3Storage("med-vault-dev", s3.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("akid", "secret", "token"),
	}))

	got, err := storage.GenerateUploadURL(context.Background(), "tenant/case/image.png", "image/png", time.Minute)
	if err != nil {
		t.Fatalf("generate upload url: %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	if !strings.Contains(u.Host, "med-vault-dev") {
		t.Fatalf("host = %q, want bucket host", u.Host)
	}
	if u.Query().Get("X-Amz-Expires") != "60" {
		t.Fatalf("expires = %q, want 60", u.Query().Get("X-Amz-Expires"))
	}
	if u.Query().Get("X-Amz-Signature") == "" {
		t.Fatal("missing signature")
	}
}

func TestS3StorageGenerateDownloadURL(t *testing.T) {
	storage := newS3Storage("med-vault-dev", s3.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("akid", "secret", "token"),
	}))

	got, err := storage.GenerateDownloadURL(context.Background(), "tenant/case/image.png", time.Minute)
	if err != nil {
		t.Fatalf("generate download url: %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	if !strings.Contains(u.Host, "med-vault-dev") {
		t.Fatalf("host = %q, want bucket host", u.Host)
	}
	if u.Query().Get("X-Amz-Signature") == "" {
		t.Fatal("missing signature")
	}
}
