package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

type videoJsonPresenter struct{}

// VideoJsonResponse represents the response of a video
func NewVideoJsonPresenter() port.Presenter {
	return &videoJsonPresenter{}
}

// Present write the response to the client
func (p *videoJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Video:
		output := ToVideoJsonResponse(v)
		return json.Marshal(output)
	case []*entity.Video:
		videoOutputs := make([]VideoJsonResponse, len(v))
		for i, video := range v {
			videoOutputs[i] = ToVideoJsonResponse(video)
		}

		output := &VideoJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Videos: videoOutputs,
		}
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}

// ToVideoJsonResponse convert entity.Video to VideoJsonResponse
func ToVideoJsonResponse(video *entity.Video) VideoJsonResponse {
	return VideoJsonResponse{
		ID:          video.ID,
		UserID:      video.UserID,
		Name:        video.Name,
		Description: video.Description,
		Status:      string(video.Status),
		CreatedAt:   video.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   video.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
