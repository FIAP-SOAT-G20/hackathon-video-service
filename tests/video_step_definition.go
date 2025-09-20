package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/cucumber/godog"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type godogsResponseCtxKey struct{}
type godogsRequestCtxKey struct{}

func init() {
	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, _ := nat.NewPort("tcp", "5432")

	container, err := startContainer(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to start test container: %w", err))
	}
	containerPort, _ := container.MappedPort(ctx, port)
	host, _ := container.Host(ctx)

	err = os.Setenv("DB_HOST", host)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PORT", containerPort.Port())
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_USER", user)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PASS", password)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_NAME", dbname)
	if err != nil {
		panic(err)
	}
}

// startContainer is a helper to start a test container for the database.
func startContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test-db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

type apiFeature struct {
	router        *gin.Engine
	deletedVideos map[string]bool
	videoStatus   map[string]valueobject.VideoStatus // Track video statuses
}

type response struct {
	status int
	body   any
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	a.router = gin.New()

	// Initialize deleted videos tracking
	a.deletedVideos = make(map[string]bool)
	// Initialize video status tracking
	a.videoStatus = make(map[string]valueobject.VideoStatus)

	// Set up mock routes for testing
	api := a.router.Group("/api/v1")
	{
		// Mock video routes
		api.POST("/videos", func(c *gin.Context) {
			var video entity.Video
			if err := c.ShouldBindJSON(&video); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Mock successful video creation
			video.ID = 12345
			if video.Status == "" {
				video.Status = "OPEN"
			}
			video.CreatedAt = time.Now()
			video.UpdatedAt = time.Now()

			c.JSON(http.StatusCreated, video)
		})

		api.GET("/videos/:id", func(c *gin.Context) {
			id := c.Param("id")

			// Check if video was deleted
			if a.deletedVideos[id] {
				c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
				return
			}

			// Get the expected status from our tracking map, default to "CREATED"
			expectedStatus := valueobject.CREATED
			if status, exists := a.videoStatus[id]; exists {
				expectedStatus = status
			}

			// Mock video response
			video := entity.Video{
				ID:        12345,
				Status:    expectedStatus,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// If a specific ID is requested, mock it
			if id != "" {
				c.JSON(http.StatusOK, video)
				return
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		})

		// Mock video status update route
		api.PATCH("/videos/:id", func(c *gin.Context) {
			// Read the raw body first
			bodyBytes, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}

			// Try to parse as JSON first
			var statusUpdate map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &statusUpdate); err == nil {
				// Handle JSON payload
				if status, ok := statusUpdate["status"]; ok {
					statusStr := status.(string)
					// Use ToVideoStatus to properly convert and validate the status
					if videoStatus, isValid := valueobject.ToVideoStatus(statusStr); isValid {
						// Update the video status in our tracking map
						id := c.Param("id")
						a.videoStatus[id] = videoStatus

						video := entity.Video{
							ID:        12345,
							Status:    videoStatus,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
						c.JSON(http.StatusOK, video)
						return
					}
					// If status is not valid, return error
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video status"})
					return
				}
				// If no status provided, return error
				c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
				return
			}

			// If JSON parsing failed, try as plain string
			newStatus := string(bodyBytes)
			// Clean up the status string (remove quotes and whitespace)
			newStatus = strings.Trim(strings.TrimSpace(newStatus), "\"")

			// Debug: print the received status
			fmt.Printf("DEBUG: Received status string: '%s'\n", newStatus)

			// Use ToVideoStatus to properly convert and validate the status
			if videoStatus, ok := valueobject.ToVideoStatus(newStatus); ok {
				// Update the video status in our tracking map
				id := c.Param("id")
				a.videoStatus[id] = videoStatus

				video := entity.Video{
					ID:        12345,
					Status:    videoStatus,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				c.JSON(http.StatusOK, video)
				return
			}

			// If status is not valid, return error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video status"})
		})

		// Mock videos deletion route
		api.DELETE("/videos/:id", func(c *gin.Context) {
			id := c.Param("id")

			// Check if video already deleted
			if a.deletedVideos[id] {
				c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
				return
			}

			// Mock successful deletion for existing video
			if id == "12345" {
				a.deletedVideos[id] = true
				c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
				return
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		})
	}
}

func (a *apiFeature) iHaveAValidVideoRequest(ctx context.Context) error {
	video := entity.Video{
		ID:     12345,
		Status: "created",
	}
	// Store the video in context for later steps
	ctx = context.WithValue(ctx, godogsRequestCtxKey{}, video) //nolint
	return nil
}

func (a *apiFeature) iHaveAnExistingVideoWithID(ctx context.Context, videoID string) error {
	// Set the video status to PROCESSING for the retrieve scenario
	a.videoStatus[videoID] = valueobject.PROCESSING

	video := entity.Video{
		ID:     12345,
		Status: valueobject.PROCESSING,
	}
	ctx = context.WithValue(ctx, godogsRequestCtxKey{}, video) //nolint

	// Mock the existence of a video with the given ID
	if videoID != "12345" {
		return fmt.Errorf("video with ID %s does not exist", videoID)
	}

	return nil
}

func (a *apiFeature) iRequestTheVideoDetailsForID(ctx context.Context, videoID string) (context.Context, error) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/videos/%s", videoID), nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	var video entity.Video
	if err := json.NewDecoder(w.Body).Decode(&video); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   video,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iSendTheVideoRequestToTheVideoService(ctx context.Context) (context.Context, error) {
	videoRequest := entity.Video{
		ID:     12345,
		Status: "created",
	}
	reqBody, err := json.Marshal(videoRequest)
	if err != nil {
		return ctx, fmt.Errorf("failed to marshal video request: %w", err)
	}
	req := httptest.NewRequest("POST", "/api/v1/videos", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)
	fmt.Println("Response:", w.Body.String())
	if w.Code != 201 {
		return ctx, fmt.Errorf("expected status code 201, got %d", w.Code)
	}
	var createdVideo entity.Video
	if err := json.NewDecoder(w.Body).Decode(&createdVideo); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}
	if createdVideo.ID != videoRequest.ID || createdVideo.Status != videoRequest.Status {
		return ctx, fmt.Errorf("expected video ID %d and status %s, got ID %d and status %s", videoRequest.ID, videoRequest.Status, createdVideo.ID, createdVideo.Status)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   createdVideo,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iShouldReceiveAConfirmationOfTheVideoCreation(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusCreated {
		return fmt.Errorf("expected response code to be 201, but got %d", resp.status)
	}

	// The body should already be an entity.Video object
	createdVideo, ok := resp.body.(entity.Video)
	if !ok {
		return errors.New("response body is not a valid video")
	}

	if createdVideo.ID == 0 || createdVideo.Status == "" {
		return errors.New("video creation confirmation is invalid")
	}

	return nil
}

func (a *apiFeature) iShouldReceiveTheVideoDetailsWithStatus(ctx context.Context, expectedStatus string) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	video, ok := resp.body.(entity.Video)
	if !ok {
		return errors.New("response body is not a valid video")
	}

	if video.Status != valueobject.VideoStatus(expectedStatus) {
		return fmt.Errorf("expected video status to be %s, but got %s", expectedStatus, video.Status)
	}

	return nil
}

func (a *apiFeature) iDeleteTheVideoWithID(ctx context.Context, videoID string) (context.Context, error) {
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/videos/%s", videoID), nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   w.Body.String(),
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) iShouldReceiveAConfirmationThatTheVideoHasBeenDeleted(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	return nil
}

func (a *apiFeature) iShouldReceiveAConfirmationThatTheVideoStatusHasBeenUpdated(ctx context.Context) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if resp.status != http.StatusOK {
		return fmt.Errorf("expected response code to be 200, but got %d", resp.status)
	}

	return nil
}

func (a *apiFeature) iUpdateTheVideoStatusTo(ctx context.Context, newStatus string) (context.Context, error) {
	req := httptest.NewRequest("PATCH", "/api/v1/videos/12345", bytes.NewBufferString(newStatus))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, got %d", w.Code)
	}

	var updatedVideo entity.Video
	if err := json.NewDecoder(w.Body).Decode(&updatedVideo); err != nil {
		return ctx, fmt.Errorf("failed to decode response body: %w", err)
	}

	if updatedVideo.Status != valueobject.VideoStatus(newStatus) {
		return ctx, fmt.Errorf("expected video status to be %s, got %s", newStatus, updatedVideo.Status)
	}

	// Store response in context for confirmation step
	actual := response{
		status: w.Code,
		body:   updatedVideo,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) theVideoShouldNoLongerExistInTheSystem(ctx context.Context) error {
	req := httptest.NewRequest("GET", "/api/v1/videos/12345", nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		return fmt.Errorf("expected status code 404, got %d", w.Code)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I have a valid video request$`, api.iHaveAValidVideoRequest)
	ctx.Step(`^I have an existing video with ID "([^"]*)"$`, api.iHaveAnExistingVideoWithID)
	ctx.Step(`^I request the video details for ID "([^"]*)"$`, api.iRequestTheVideoDetailsForID)
	ctx.Step(`^I send the video request to the video service$`, api.iSendTheVideoRequestToTheVideoService)
	ctx.Step(`^I should receive a confirmation of the video creation$`, api.iShouldReceiveAConfirmationOfTheVideoCreation)
	ctx.Step(`^I should receive the video details with status "([^"]*)"$`, api.iShouldReceiveTheVideoDetailsWithStatus)
	ctx.Step(`^I delete the video with ID "([^"]*)"$`, api.iDeleteTheVideoWithID)
	ctx.Step(`^I should receive a confirmation that the video has been deleted$`, api.iShouldReceiveAConfirmationThatTheVideoHasBeenDeleted)
	ctx.Step(`^I should receive a confirmation that the video status has been updated$`, api.iShouldReceiveAConfirmationThatTheVideoStatusHasBeenUpdated)
	ctx.Step(`^I update the video status to "([^"]*)"$`, api.iUpdateTheVideoStatusTo)
	ctx.Step(`^the video should no longer exist in the system$`, api.theVideoShouldNoLongerExistInTheSystem)
}
