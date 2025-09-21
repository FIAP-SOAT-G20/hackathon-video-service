package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler/request"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/middleware"
)

type VideoHandler struct {
	controller      port.VideoController
	jwtService      port.JWTService
	cacheMiddleware port.CacheMiddleware
	cacheDuration   time.Duration
}

func NewVideoHandler(controller port.VideoController, jwtService port.JWTService, cacheStore port.CacheMiddleware, cacheDuration time.Duration) *VideoHandler {
	return &VideoHandler{
		controller:      controller,
		jwtService:      jwtService,
		cacheMiddleware: cacheStore,
		cacheDuration:   cacheDuration,
	}
}

func (h *VideoHandler) Register(router *gin.RouterGroup) {
	// Apply cache middleware to the List endpoint
	router.GET("", middleware.CachePage(h.cacheMiddleware, h.cacheDuration, h.List))
	router.POST("", h.Create)
	router.GET("/:id", h.Get)
	router.PUT("/:id", h.Update)
	router.PATCH("/:id", h.UpdatePartial)
	router.DELETE("/:id", h.Delete)
	router.GET("/:id/processed", h.Download)
}

// List godoc
//
//	@Summary		List videos
//	@Description	List all videos
//	@Description	## Video list is sorted by:
//	@Description	- **Status** in **descending** order (CREATED > PROCESSING > FINISHED)
//	@Description	- **Created date** (CreatedAt) in **ascending** order (oldest first)
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			user_id			query		int										false	"Filter by user ID"
//	@Param			status			query		string									false	"Filter by status (Accept many), options: <sub>CREATED, PROCESSING, FINISHED</sub>, ex: <sub>CREATED</sub> or <sub>CREATED,PROCESSING</sub>"
//	@Param			status_exclude	query		string									false	"Exclude by status (Accept many), options: <sub>NONE, CREATED, PROCESSING, FINISHED, FAILED</sub>, ex: <sub>FAILED</sub> (default)"	default(FAILED)
//	@Param			hash			query		string									false	"Filter by hash"
//	@Param			sort			query		string									false	"Sort by field (Accept many). Use `<field_name>:d` for descending, and the default order is ascending"	default(status:d,created_at)
//	@Param			page			query		int										false	"Page number"																							default(1)
//	@Param			limit			query		int										false	"Items per page"																						default(10)
//	@Success		200				{object}	presenter.VideoJsonPaginatedResponse	"OK"
//	@Failure		400				{object}	middleware.ErrorJsonResponse			"Bad Request"
//	@Failure		500				{object}	middleware.ErrorJsonResponse			"Internal Server Error"
//	@Router			/videos [get]
func (h *VideoHandler) List(c *gin.Context) {
	var query request.ListVideosQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	// Default sort
	if query.Sort == "" {
		query.Sort = "status:d,created_at"
	}

	// Default status_exclude
	var statusExclude []valueobject.VideoStatus
	if query.StatusExclude == "" {
		query.StatusExclude = "FAILED"
	}

	// Convert status_exclude
	if strings.ToUpper(query.StatusExclude) != "NONE" {
		for _, s := range strings.Split(query.StatusExclude, ",") {
			statusExclude = append(statusExclude, valueobject.VideoStatus(strings.TrimSpace(s)))
		}
	}

	// Convert status
	var status []valueobject.VideoStatus
	if query.Status != "" {
		for _, s := range strings.Split(query.Status, ",") {
			videoStatus, ok := valueobject.ToVideoStatus(strings.TrimSpace(s))
			if !ok {
				_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
				return
			}
			status = append(status, videoStatus)
		}
	}

	input := dto.ListVideosInput{
		UserID:        query.UserID,
		Status:        status,
		StatusExclude: statusExclude,
		Hash:          query.Hash,
		Page:          query.Page,
		Limit:         query.Limit,
		Sort:          query.Sort,
	}

	output, err := h.controller.List(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Create godoc
//
//	@Summary		Create video
//	@Description	Creates a new video
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			video	body		request.CreateVideoBodyRequest	true	"Video data"
//	@Success		201		{object}	presenter.VideoJsonResponse		"Created"
//	@Failure		400		{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		500		{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/videos [post]
func (h *VideoHandler) Create(c *gin.Context) {
	var body request.CreateVideoBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.CreateVideoInput{
		UserID:      body.UserID,
		Name:        body.Name,
		Description: body.Description,
	}

	output, err := h.controller.Create(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusCreated, "application/json", output)
}

// Get godoc
//
//	@Summary		Get video
//	@Description	Search for a video by ID
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"Video ID"
//	@Success		200	{object}	presenter.VideoJsonResponse		"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/videos/{id} [get]
func (h *VideoHandler) Get(c *gin.Context) {
	var uri request.GetVideoUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.GetVideoInput{
		ID: uri.ID,
	}

	output, err := h.controller.Get(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Update godoc
//
//	@Summary		Update video
//	@Description	Update an existing video
//	@Description	The status are: **CREATED**, **FAILED**, **PROCESSING**, **FINISHED**
//	@Description	## Transition of status:
//	@Description	- CREATED      -> FAILED || PROCESSING
//	@Description	- FAILED       -> {},
//	@Description	- PROCESSING   -> FINISHED
//	@Description	- FINISHED     -> {}
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Video ID"
//	@Param			video	body		request.UpdateVideoBodyRequest	true	"Video data"
//	@Success		200		{object}	presenter.VideoJsonResponse		"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404		{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500		{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/videos/{id} [put]
func (h *VideoHandler) Update(c *gin.Context) {
	var uri request.UpdateVideoUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateVideoBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateVideoInput{
		ID:     uri.ID,
		Status: body.Status,
		Hash:   body.Hash,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// UpdatePartial godoc
//
//	@Summary		Partial update video
//	@Description	Partially updates an existing video
//	@Description	The status are: **CREATED**, **FAILED**, **PROCESSING**, **FINISHED**
//	@Description	## Transition of status:
//	@Description	- CREATED      -> FAILED || PROCESSING
//	@Description	- FAILED       -> {},
//	@Description	- PROCESSING   -> FINISHED
//	@Description	- FINISHED     -> {}
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Video ID"
//	@Param			video	body		request.UpdateVideoPartilRequest	true	"Video data"
//	@Success		200		{object}	presenter.VideoJsonResponse			"OK"
//	@Failure		400		{object}	middleware.ErrorJsonResponse		"Bad Request"
//	@Failure		404		{object}	middleware.ErrorJsonResponse		"Not Found"
//	@Failure		500		{object}	middleware.ErrorJsonResponse		"Internal Server Error"
//	@Router			/videos/{id} [patch]
func (h *VideoHandler) UpdatePartial(c *gin.Context) {
	var uri request.UpdateVideoUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	var body request.UpdateVideoPartilBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidBody))
		return
	}

	input := dto.UpdateVideoInput{
		ID:     uri.ID,
		Status: body.Status,
		Hash:   body.Hash,
	}

	output, err := h.controller.Update(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Delete godoc
//
//	@Summary		Delete video
//	@Description	Deletes a video by ID
//	@Tags			videos
//	@Produce		json
//	@Param			id	path		int								true	"Video ID"
//	@Success		200	{object}	presenter.VideoJsonResponse		"OK"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/videos/{id} [delete]
func (h *VideoHandler) Delete(c *gin.Context) {
	var uri request.DeleteVideoUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DeleteVideoInput{
		ID: uri.ID,
	}

	output, err := h.controller.Delete(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Download godoc
//
//	@Summary		Download processed video
//	@Description	Generate a presigned URL to download a processed video from S3
//	@Description	Only available for videos with status FINISHED
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"Video ID"
//	@Success		200	{object}	map[string]string				"OK - Returns download_url"
//	@Failure		400	{object}	middleware.ErrorJsonResponse	"Bad Request"
//	@Failure		404	{object}	middleware.ErrorJsonResponse	"Not Found"
//	@Failure		422	{object}	middleware.ErrorJsonResponse	"Video not processed yet"
//	@Failure		500	{object}	middleware.ErrorJsonResponse	"Internal Server Error"
//	@Router			/videos/{id}/processed [get]
func (h *VideoHandler) Download(c *gin.Context) {
	var uri request.DownloadVideoUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(domain.NewInvalidInputError(domain.ErrInvalidParam))
		return
	}

	input := dto.DownloadVideoInput{
		ID: uri.ID,
	}

	output, err := h.controller.Download(
		c.Request.Context(),
		presenter.NewVideoJsonPresenter(),
		input,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}
