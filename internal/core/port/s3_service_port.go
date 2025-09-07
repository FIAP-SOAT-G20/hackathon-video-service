package port

import (
	"context"
	"time"
)

// S3Service defines the interface for S3 operations
type S3Service interface {
	// GeneratePresignedURL generates a presigned URL for uploading a file to S3
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
