package sqs

import (
	"context"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	awsclient "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestSqsClient creates a new SQS client for testing
func createTestSqsClient(t *testing.T, ctx context.Context) *SqsClient {
	awsClientFactory, err := awsclient.NewClientFactory(ctx, "us-east-1")
	require.NoError(t, err)

	sqsClient, err := NewSqsClient(awsClientFactory)
	require.NoError(t, err)
	require.NotNil(t, sqsClient)

	return sqsClient
}

func TestNewSqsHandler(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Create logger
	logger := logger.NewLogger("test")

	// Create SQS handler
	handler := NewSqsHandler(sqsClient, "https://test-queue.com", 10, 10, logger)
	require.NotNil(t, handler)
	assert.Equal(t, sqsClient, handler.sqsClient)
	assert.Equal(t, "https://test-queue.com", handler.queueURL)
	assert.Equal(t, 10, handler.maxMessages)
}

func TestSqsHandler_ReceiveMessages(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Create logger
	logger := logger.NewLogger("test")

	// Create SQS handler
	handler := NewSqsHandler(sqsClient, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", 10, 10, logger)

	// Test processing messages with invalid queue URL
	err := handler.ReceiveMessages(ctx, func(message types.Message) (bool, error) {
		return false, nil
	})

	// We expect an error since the queue doesn't exist, but the handler should handle it gracefully
	assert.Error(t, err)
}

func TestSqsHandler_SendMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Create logger
	logger := logger.NewLogger("test")

	// Create SQS handler
	handler := NewSqsHandler(sqsClient, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", 10, 10, logger)

	// Test sending message to invalid queue
	message, err := handler.SendMessage(ctx, "test message")

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
	assert.Nil(t, message)
}

func TestSqsClient_ReceiveMessages(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Test with invalid queue URL
	messages, err := sqsClient.ReceiveMessages(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", 1, 1)

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
	assert.Nil(t, messages)
}

func TestSqsClient_SendMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Test sending message to invalid queue
	message, err := sqsClient.SendMessage(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", "test message")

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
	assert.Nil(t, message)
}

func TestSqsClient_DeleteMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient := createTestSqsClient(t, ctx)

	// Test with invalid queue URL and receipt handle
	err := sqsClient.DeleteMessage(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", "invalid-receipt-handle")

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
}
