package service

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	awsclient "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new S3Service instance using AWS client factory
func NewS3Service(ctx context.Context, bucketName string, awsClientFactory *awsclient.ClientFactory) (port.S3Service, error) {
	client := s3.NewFromConfig(awsClientFactory.GetConfig())

	return &S3Service{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// GeneratePresignedURL generates a presigned URL for uploading a file to S3
func (s *S3Service) GeneratePresignedURL(ctx context.Context, key string, metadata map[string]string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      awssdk.String(s.bucketName),
		Key:         awssdk.String(key),
		ContentType: awssdk.String("video/mp4"),
		Metadata:    metadata,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// GeneratePresignedDownloadURL generates a presigned URL for downloading a processed file from S3
func (s *S3Service) GeneratePresignedDownloadURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	desiredFilename := "fiapx-video-images.zip"

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     awssdk.String(s.bucketName),
		Key:                        awssdk.String(key),
		ResponseContentDisposition: awssdk.String(fmt.Sprintf("attachment; filename=\"%s\"", desiredFilename)),
		ResponseContentType:        awssdk.String("application/zip"),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return request.URL, nil
}
