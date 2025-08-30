package sqs

import (
	"context"
	"testing"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSqsHandler(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)
	require.NotNil(t, sqsClient)

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
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)

	// Create logger
	logger := logger.NewLogger("test")

	// Create SQS handler
	handler := NewSqsHandler(sqsClient, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", 10, 10, logger)

	// Test processing messages with invalid queue URL
	err = handler.ReceiveMessages(ctx, func(message types.Message) (bool, error) {
		return false, nil
	})

	// We expect an error since the queue doesn't exist, but the handler should handle it gracefully
	assert.Error(t, err)
}

func TestSqsHandler_SendMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)

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
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)

	// Test with invalid queue URL
	messages, err := sqsClient.ReceiveMessages(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", 1, 1)

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
	assert.Nil(t, messages)
}

func TestSqsClient_SendMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)

	// Test sending message to invalid queue
	message, err := sqsClient.SendMessage(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", "test message")

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
	assert.Nil(t, message)
}

func TestSqsClient_DeleteMessage(t *testing.T) {
	ctx := context.Background()

	// Create SQS client
	sqsClient, err := NewSqsClient(ctx)
	require.NoError(t, err)

	// Test with invalid queue URL and receipt handle
	err = sqsClient.DeleteMessage(ctx, "https://sqs.invalid-region.amazonaws.com/123456789012/test-queue", "invalid-receipt-handle")

	// We expect an error since the queue doesn't exist
	assert.Error(t, err)
}
