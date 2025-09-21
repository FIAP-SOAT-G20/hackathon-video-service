package main

import (
	"context"
	"encoding/json"
	"errors"
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
)

type VideoUpdated struct {
	VideoID uint64                  `json:"video_id"`
	Status  valueobject.VideoStatus `json:"status"`
	Hash    string                  `json:"hash"`
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
	logger.Info("Processing message", "messageID", *message.MessageId, "body", *message.Body)

	// Unmarshal the message body to your entity
	var updatedVideo VideoUpdated
	err = json.Unmarshal([]byte(*message.Body), &updatedVideo)
	if err != nil {
		return false, err
	}

	if updatedVideo.VideoID == 0 {
		return false, domain.NewValidationError(errors.New(domain.ErrVideoIsMandatory))
	}

	if updatedVideo.Status == "" {
		return false, domain.NewValidationError(errors.New(domain.ErrStatusIsMandatory))
	}

	// Get Video by ID
	_, err = uc.Get(ctx, dto.GetVideoInput{ID: updatedVideo.VideoID})
	if err != nil {
		if err.Error() == domain.ErrInternalError {
			return true, err
		}
		return false, err
	}

	// Update the video status in the database
	uoi := dto.UpdateVideoInput{
		ID:     updatedVideo.VideoID,
		Status: updatedVideo.Status,
		Hash:   updatedVideo.Hash,
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
		"status", updatedVideo.Status,
		"hash", updatedVideo.Hash,
	)

	return false, nil
}
