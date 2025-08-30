package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/usecase"
	appConfig "github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/pkg/aws/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

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

	orderDS := datasource.NewOrderDataSource(db.DB)
	orderHistoryDS := datasource.NewOrderHistoryDataSource(db.DB)
	orderGateway := gateway.NewOrderGateway(orderDS)
	orderHistoryGateway := gateway.NewOrderHistoryGateway(orderHistoryDS)
	orderHistoryUC := usecase.NewOrderHistoryUseCase(orderHistoryGateway)
	orderUC := usecase.NewOrderUseCase(orderGateway, orderHistoryUC)

	if appCfg.AWS_SQS_OrderStatusUpdatedURL == "" {
		loggerInstance.Error("AWS SQS Order Status Updated URL is not configured")
		os.Exit(1)
	}

	sqsClient, err := sqs.NewSqsClient(ctx)
	if err != nil {
		loggerInstance.Error("Failed to create SQS client", "error", err.Error())
		os.Exit(1)
	}

	sqsHandler := sqs.NewSqsHandler(
		sqsClient,
		appCfg.AWS_SQS_OrderStatusUpdatedURL,
		appCfg.AWS_SQS_OrderStatusUpdatedMaxMessages,
		appCfg.AWS_SQS_OrderStatusUpdatedWaitTimeSeconds,
		loggerInstance,
	)

	loggerInstance.Info("Starting SQS consumer", "queueURL", appCfg.AWS_SQS_OrderStatusUpdatedURL)

	// Receive messages from SQS
	for {
		err = sqsHandler.ReceiveMessages(ctx, func(message types.Message) (bool, error) {
			loggerInstance.Info("Processing message", "message", message)

			reprocess, err := processedMessage(ctx, message, loggerInstance, orderUC)
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

func processedMessage(ctx context.Context, message types.Message, logger *logger.Logger, uc port.OrderUseCase) (reprocess bool, err error) {
	// Here you can implement the logic to process the message
	// For example, you can unmarshal the message body and update the order status in your database
	logger.Info("Processing message", "messageID", *message.MessageId, "body", *message.Body)

	// Unmarshal the message body to your entity
	var updatedOrderStatus entity.OrderStatusUpdated
	err = json.Unmarshal([]byte(*message.Body), &updatedOrderStatus)
	if err != nil {
		return false, err
	}

	if updatedOrderStatus.OrderID == 0 {
		return false, domain.NewValidationError(errors.New(domain.ErrOrderIsMandatory))
	}

	if updatedOrderStatus.Status == "" {
		return false, domain.NewValidationError(errors.New(domain.ErrStatusIsMandatory))
	}

	// Get Order by ID
	_, err = uc.Get(ctx, dto.GetOrderInput{ID: updatedOrderStatus.OrderID})
	if err != nil {
		if err.Error() == domain.ErrInternalError {
			return true, err
		}
		return false, err
	}

	// Update the order status in the database
	uoi := dto.UpdateOrderInput{
		ID:     updatedOrderStatus.OrderID,
		Status: updatedOrderStatus.Status,
	}
	if updatedOrderStatus.StaffID != nil {
		uoi.StaffID = *updatedOrderStatus.StaffID
	}
	_, err = uc.Update(ctx, uoi)
	if err != nil {
		if err.Error() == domain.ErrInternalError {
			return true, err
		}
		return false, err
	}

	logger.Info("Message processed successfully", "orderID", updatedOrderStatus.OrderID, "status", updatedOrderStatus.Status, "staffID", updatedOrderStatus.StaffID)

	return false, nil
}
