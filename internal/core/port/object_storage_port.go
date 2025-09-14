package port

import (
	"context"
	"time"
)

// ObjectStorageDatasource defines the interface for S3 operations
type ObjectStorageDatasource interface {
	// GeneratePresignedURL generates a presigned URL for uploading a file to S3
	GeneratePresignedURL(ctx context.Context, bucketName, key string, metadata map[string]string, expiration time.Duration) (string, error)
	// GeneratePresignedDownloadURL generates a presigned URL for downloading a processed file from S3
	GeneratePresignedDownloadURL(ctx context.Context, bucketName, key string, expiration time.Duration) (string, error)
}
