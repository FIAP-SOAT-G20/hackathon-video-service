package datasource

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
)

// TestVideoDocumentDataSource_Simple provides a simple test without authentication
func TestVideoDocumentDataSource_Simple(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Start MongoDB container without authentication (simpler and more reliable)
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("27017/tcp"),
			wait.ForLog("Waiting for connections").WithOccurrence(1).WithStartupTimeout(30*time.Second),
		),
	}
	
	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer mongoContainer.Terminate(ctx)

	// Get container endpoint
	endpoint, err := mongoContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	// Setup test configuration without authentication
	cfg := &config.Config{
		DocumentDBURI:  "mongodb://" + endpoint + "/test_video_service",
		DocumentDBName: "test_video_service",
	}

	logger := logger.NewLogger("test")

	// Connect to DocumentDB
	db, err := database.NewDocumentDBConnection(cfg, logger)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create datasource
	videoDS := NewVideoDocumentDataSource(db)

	t.Run("Create and Find Video", func(t *testing.T) {
		// Create a test video
		video := &entity.Video{
			ID:        1,
			UserID:    100,
			Status:    valueobject.CREATED,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := videoDS.Create(ctx, video)
		assert.NoError(t, err)

		// Find the video
		found, err := videoDS.FindByID(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, uint64(1), found.ID)
		assert.Equal(t, uint64(100), found.UserID)
		assert.Equal(t, valueobject.CREATED, found.Status)
	})

	t.Run("Update Video", func(t *testing.T) {
		// Find the video created in previous test
		video, err := videoDS.FindByID(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, video)

		// Update the video
		video.Status = valueobject.PROCESSING
		video.UpdatedAt = time.Now()

		err = videoDS.Update(ctx, video)
		assert.NoError(t, err)

		// Verify the update
		updated, err := videoDS.FindByID(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, valueobject.PROCESSING, updated.Status)
	})

	t.Run("Delete Video", func(t *testing.T) {
		// Delete the video
		err := videoDS.Delete(ctx, 1)
		assert.NoError(t, err)

		// Verify deletion
		found, err := videoDS.FindByID(ctx, 1)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("Find All with Filters", func(t *testing.T) {
		// Create multiple videos
		videos := []*entity.Video{
			{ID: 10, UserID: 200, Status: valueobject.CREATED, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 11, UserID: 200, Status: valueobject.PROCESSING, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 12, UserID: 201, Status: valueobject.FINISHED, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		for _, video := range videos {
			err := videoDS.Create(ctx, video)
			require.NoError(t, err)
		}

		// Test filter by user ID
		filters := map[string]any{"user_id": uint64(200)}
		found, total, err := videoDS.FindAll(ctx, filters, "", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, found, 2)

		// Test filter by status
		filters = map[string]any{"statuses": []valueobject.VideoStatus{valueobject.CREATED}}
		found, total, err = videoDS.FindAll(ctx, filters, "", 1, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, video := range found {
			assert.Equal(t, valueobject.CREATED, video.Status)
		}
	})

	t.Run("Transaction", func(t *testing.T) {
		// Test transaction functionality
		err := videoDS.Transaction(ctx, func(txCtx context.Context) error {
			video := &entity.Video{
				ID:        20,
				UserID:    300,
				Status:    valueobject.PROCESSING,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			return videoDS.Create(txCtx, video)
		})
		
		// Check if this is a standalone MongoDB (transactions not supported)
		if err != nil && (strings.Contains(err.Error(), "replica set member") || 
			strings.Contains(err.Error(), "Transaction numbers are only allowed")) {
			t.Skip("Skipping transaction test: MongoDB transactions require a replica set or sharded cluster")
			return
		}
		
		assert.NoError(t, err)

		// Verify the video was created (only if transaction succeeded)
		if err == nil {
			found, err := videoDS.FindByID(ctx, 20)
			assert.NoError(t, err)
			assert.NotNil(t, found)
			if found != nil {
				assert.Equal(t, valueobject.PROCESSING, found.Status)
			}
		}
	})
}
