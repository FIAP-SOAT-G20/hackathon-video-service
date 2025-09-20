package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
)

// DatabaseType represents the type of database
type DatabaseType string

const (
	PostgresDB DatabaseType = "postgres"
	DocumentDB DatabaseType = "documentdb"
	MongoDB    DatabaseType = "mongodb"
)

// DatabaseConfig holds configuration for different database types
type DatabaseConfig struct {
	Type       DatabaseType
	PostgresDB *Database
	DocumentDB *DocumentDatabase
}

// NewDatabase creates a new database connection based on configuration
func NewDatabase(cfg *config.Config, logger *logger.Logger) (*DatabaseConfig, error) {
	dbType := determineDatabaseType(cfg)

	dbConfig := &DatabaseConfig{
		Type: dbType,
	}

	switch dbType {
	case PostgresDB:
		db, err := NewPostgresConnection(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		dbConfig.PostgresDB = db

		// Run migrations for PostgreSQL
		if err := db.Migrate(); err != nil {
			return nil, fmt.Errorf("failed to run PostgreSQL migrations: %w", err)
		}

		logger.Info("Connected to PostgreSQL database")

	case DocumentDB, MongoDB:
		db, err := NewDocumentDBConnection(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DocumentDB/MongoDB: %w", err)
		}
		dbConfig.DocumentDB = db

		// Create indexes for DocumentDB/MongoDB
		if err := db.CreateIndexes(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to create DocumentDB indexes: %w", err)
		}

		logger.Info("Connected to DocumentDB/MongoDB database")

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	return dbConfig, nil
}

// determineDatabaseType determines which database to use based on configuration
func determineDatabaseType(cfg *config.Config) DatabaseType {
	// Check environment variable for explicit database type
	if dbType := cfg.DBEngine; dbType != "" {
		switch strings.ToLower(dbType) {
		case "documentdb":
			return DocumentDB
		case "mongodb":
			return MongoDB
		case "postgres", "postgresql":
			return PostgresDB
		}
	}

	// Default to PostgreSQL
	return PostgresDB
}

// Close closes the database connection
func (dc *DatabaseConfig) Close() error {
	switch dc.Type {
	case PostgresDB:
		if dc.PostgresDB != nil {
			if sqlDB, err := dc.PostgresDB.DB.DB(); err == nil {
				return sqlDB.Close()
			}
		}
	case DocumentDB, MongoDB:
		if dc.DocumentDB != nil {
			return dc.DocumentDB.Close(context.Background())
		}
	}
	return nil
}
