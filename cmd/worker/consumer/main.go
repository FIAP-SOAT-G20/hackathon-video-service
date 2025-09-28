package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/usecase"
	appConfig "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	awsclient "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws/s3"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SNSNotification represents the AWS SNS notification structure
type SNSNotification struct {
	Type             string `json:"Type"`
	MessageID        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

// VideoUpdated represents the video update data nested in SNS Message
type VideoUpdated struct {
	VideoID    uint64                  `json:"video_id"`
	UserID     uint64                  `json:"user_id"`
	Status     valueobject.VideoStatus `json:"status"`
	OccurredAt string                  `json:"occurred_at"`
}

func main() {

	ctx := context.Background()

	// Load AWS Config
	appCfg := appConfig.LoadConfig()

	loggerInstance := logger.NewLogger(appCfg.Environment)

	db, err := database.NewPostgresConnection(appCfg, loggerInstance)
	if err != nil {
		loggerInstance.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	videoDS := datasource.NewVideoDataSource(db.DB)
	videoGateway := gateway.NewVideoGateway(videoDS)

	// Create AWS client factory
	awsClientFactory, err := awsclient.NewClientFactory(ctx, appCfg.AWSS3Region)
	if err != nil {
		loggerInstance.Error("Failed to create AWS client factory", "error", err.Error())
		os.Exit(1)
	}

	s3Client, err := s3.NewS3Client(awsClientFactory)
	if err != nil {
		loggerInstance.Error("Failed to create S3 client", "error", err.Error())
		os.Exit(1)
	}

	// Since this is a worker, we don't need cache service for the use case
	videoUC := usecase.NewVideoUseCase(videoGateway, s3Client, appCfg, loggerInstance)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		port := appCfg.MetricsPort
		loggerInstance.Info("Prometheus metrics server running", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			loggerInstance.Error("Metrics server failed", "error", err.Error())
		}
	}()

	if appCfg.AWS_SQS_VideoUpdatedURL == "" {
		loggerInstance.Error("AWS SQS Order Status Updated URL is not configured")
		os.Exit(1)
	}

	sqsClient, err := sqs.NewSqsClient(awsClientFactory)
	if err != nil {
		loggerInstance.Error("Failed to create SQS client", "error", err.Error())
		os.Exit(1)
	}

	sqsHandler := sqs.NewSqsHandler(
		sqsClient,
		appCfg.AWS_SQS_VideoUpdatedURL,
		appCfg.AWS_SQS_VideoUpdatedMaxMessages,
		appCfg.AWS_SQS_VideoUpdatedWaitTimeSeconds,
		loggerInstance,
	)

	loggerInstance.Info("Starting SQS consumer", "queueURL", appCfg.AWS_SQS_VideoUpdatedURL)

	// Receive messages from SQS
	for {
		err = sqsHandler.ReceiveMessages(ctx, func(message types.Message) (bool, error) {
			loggerInstance.Info("Processing message", "message", message)

			reprocess, err := processedMessage(ctx, message, loggerInstance, videoUC)
			if err != nil {
				loggerInstance.Error("Failed to process message", "error", err.Error(), "messageID", *message.MessageId)
				return reprocess, err
			}

			return false, nil
		})
		if err != nil {
			loggerInstance.Error("Failed to receive messages", "error", err.Error())
		}
	}
}

func processedMessage(ctx context.Context, message types.Message, logger *logger.Logger, uc port.VideoUseCase) (reprocess bool, err error) {
	// Here you can implement the logic to process the message
	// For example, you can unmarshal the message body and update the order status in your database
	//{\n  \"Type\" : \"Notification\",\n  \"MessageId\" : \"2a33bf6b-bb93-57c3-afd0-de38c7b6f234\",\n  \"TopicArn\" : \"arn:aws:sns:us-east-1:905417995957:video-status-updated\",\n  \"Message\" : \"{\\\"video_id\\\":25,\\\"user_id\\\":5,\\\"status\\\":\\\"FINISHED\\\",\\\"occurred_at\\\":\\\"2025-09-28T18:09:41.110973622Z\\\"}\",\n  \"Timestamp\" : \"2025-09-28T18:09:41.115Z\",\n  \"SignatureVersion\" : \"1\",\n  \"Signature\" : \"n2sK9472MGBlYH6D58MSJjo64pxWlpevdXgJxqmPLhkKf2Aox+90cADrCmycfQaHpRVqFCJwbMvKl2JSofOBjtpdw33LQzyJi9KsQQ6IbjYiiIsgf2SVTqJZdeZeJbBAZ533iFyfOhK5lVM//nLiRSVrz5zHYHQfmzKLYfY/B6KxvE8S3X5nxYG3sAg7bk3gnp92kpLAVRojwNif+XUnDYrliCyBNmEPQg/z9Y7hR+LT+K5OPiwKjZ/u6wLB7ht0E4c+uRU6+l7WONIAshM95HFh4tpO8g7UuVKYXPQ8C9XnLNsTAxminr8vTnHiD4Mewfh3N9WgA2eAXF/N1bh8Ww==\",\n  \"SigningCertURL\" : \"https://sns.us-east-1.amazonaws.com/SimpleNotificationService-6209c161c6221fdf56ec1eb5c821d112.pem\",\n  \"UnsubscribeURL\" : \"https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe\u0026SubscriptionArn=arn:aws:sns:us-east-1:905417995957:video-status-updated:c9aa674a-a2e1-4d99-9871-327425c720d1\"\n}
	logger.Info("Processing message", "messageID", *message.MessageId, "body", *message.Body)

	// First unmarshal the SNS notification structure
	var snsNotification SNSNotification
	err = json.Unmarshal([]byte(*message.Body), &snsNotification)
	if err != nil {
		return false, err
	}

	// Then unmarshal the nested Message field to get the video update data
	var updatedVideo VideoUpdated
	err = json.Unmarshal([]byte(snsNotification.Message), &updatedVideo)
	if err != nil {
		return false, err
	}

	if updatedVideo.VideoID == 0 {
		return false, domain.NewValidationError(errors.New(domain.ErrVideoIsMandatory))
	}

	if updatedVideo.Status == "" {
		return false, domain.NewValidationError(errors.New(domain.ErrStatusIsMandatory))
	}

	// Get Video by ID - use UserID from the message if provided, otherwise use VideoID as fallback
	userID := updatedVideo.UserID
	if userID == 0 {
		userID = updatedVideo.VideoID // fallback for backward compatibility
	}
	_, err = uc.Get(ctx, dto.GetVideoInput{ID: updatedVideo.VideoID, UserID: userID})
	if err != nil {
		if err.Error() == domain.ErrInternalError {
			return true, err
		}
		return false, err
	}

	// Update the video status in the database
	uoi := dto.UpdateVideoInput{
		ID:     updatedVideo.VideoID,
		UserID: updatedVideo.UserID,
		Status: updatedVideo.Status,
	}

	_, err = uc.Update(ctx, uoi)
	if err != nil {
		if err.Error() == domain.ErrInternalError {
			return true, err
		}
		return false, err
	}

	logger.Info("Message processed successfully",
		"videoID", updatedVideo.VideoID,
		"userID", updatedVideo.UserID,
		"status", updatedVideo.Status,
		"occurredAt", updatedVideo.OccurredAt,
	)

	return false, nil
}
