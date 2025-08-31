package datasource

import (
	"fmt"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/database"
)

// VideoDataSourceFactory creates video datasource based on database type
func NewVideoDataSourceFactory(dbConfig *database.DatabaseConfig) (port.VideoDataSource, error) {
	switch dbConfig.Type {
	case database.PostgresDB:
		if dbConfig.PostgresDB == nil {
			return nil, fmt.Errorf("PostgreSQL database not initialized")
		}
		return NewVideoDataSource(dbConfig.PostgresDB.DB), nil

	case database.DocumentDB, database.MongoDB:
		if dbConfig.DocumentDB == nil {
			return nil, fmt.Errorf("DocumentDB/MongoDB database not initialized")
		}
		return NewVideoDocumentDataSource(dbConfig.DocumentDB), nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}
}
