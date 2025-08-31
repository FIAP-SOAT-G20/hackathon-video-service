package gateway

import (
	"context"
	"strings"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

type videoGateway struct {
	dataSource port.VideoDataSource
}

func NewVideoGateway(dataSource port.VideoDataSource) port.VideoGateway {
	return &videoGateway{dataSource}
}

func (g *videoGateway) FindByID(ctx context.Context, id uint64) (*entity.Video, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *videoGateway) FindAll(
	ctx context.Context,
	customerId uint64,
	status []valueobject.VideoStatus,
	statusExclude []valueobject.VideoStatus,
	page,
	limit int,
	sort string,
) ([]*entity.Video, int64, error) {

	// Create filters
	filters := make(map[string]interface{})
	if customerId != 0 {
		filters["customer_id"] = customerId
	}
	if status != nil {
		filters["statuses"] = status
	}
	if statusExclude != nil {
		filters["statuses_exclude"] = statusExclude
	}

	// Create Sort "status:d,created_at" -> "status desc, created_at asc"
	sortFormatted := strings.ReplaceAll(sort, ":d", " desc")

	return g.dataSource.FindAll(ctx, filters, sortFormatted, page, limit)
}

func (g *videoGateway) Create(ctx context.Context, video *entity.Video) error {
	return g.dataSource.Create(ctx, video)
}

func (g *videoGateway) Update(ctx context.Context, video *entity.Video) error {
	return g.dataSource.Update(ctx, video)
}

func (g *videoGateway) Delete(ctx context.Context, id uint64) error {
	return g.dataSource.Delete(ctx, id)
}
