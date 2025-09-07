package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

type S3Service struct {
	client              *s3.Client
	bucketName          string
	RawFolderName       string
	ProcessedFolderName string
}

// NewS3Service creates a new S3Service instance
func NewS3Service(cfg *config.Config) (port.S3Service, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSS3Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Service{
		client:              client,
		bucketName:          cfg.AWSS3BucketName,
		RawFolderName:       cfg.AWSS3BucketRawFolder,
		ProcessedFolderName: cfg.AWSS3BucketProcessedFolder,
	}, nil
}

// GeneratePresignedURL generates a presigned URL for uploading a file to S3
func (s *S3Service) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	keyWithFolder := fmt.Sprintf("%s/%s", s.RawFolderName, key)
	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(keyWithFolder),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}
