package sqs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsClient struct {
	client *sqs.Client
}

// NewSqsClient creates a new SQS client with AWS SDK v2
func NewSqsClient(ctx context.Context) (*SqsClient, error) {
	// Load AWS configuration with credentials
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
			},
		},
	))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create SQS client
	client := sqs.NewFromConfig(cfg)

	return &SqsClient{
		client: client,
	}, nil
}

// SendMessage sends a message to an SQS queue
func (s *SqsClient) SendMessage(ctx context.Context, queueURL string, messageBody string) (*types.Message, error) {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	}

	result, err := s.client.SendMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to queue %s: %w", queueURL, err)
	}

	// Create a Message struct to return
	message := &types.Message{
		MessageId: result.MessageId,
		Body:      aws.String(messageBody),
	}

	return message, nil
}

// ReceiveMessages receives messages from an SQS queue
func (s *SqsClient) ReceiveMessages(ctx context.Context, queueURL string, maxMessages int, waitTimeSeconds int) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: int32(maxMessages),
		WaitTimeSeconds:     int32(waitTimeSeconds),
	}

	result, err := s.client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from queue %s: %w", queueURL, err)
	}

	return result.Messages, nil
}

// DeleteMessage deletes a message from the SQS queue
func (s *SqsClient) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := s.client.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete message from queue %s: %w", queueURL, err)
	}

	return nil
}

// GetClient returns the underlying SQS client
func (s *SqsClient) GetClient() *sqs.Client {
	return s.client
}
