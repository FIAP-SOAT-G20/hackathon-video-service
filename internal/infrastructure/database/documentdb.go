package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
)

type DocumentDatabase struct {
	*mongo.Client
	database string
	logger   *logger.Logger
}

func NewDocumentDBConnection(cfg *config.Config, logger *logger.Logger) (*DocumentDatabase, error) {
	// Set up client options
	clientOptions := options.Client().ApplyURI(cfg.DocumentDBURI)

	// Set connection pool settings
	clientOptions.SetMaxPoolSize(uint64(cfg.DBMaxOpenConns))
	clientOptions.SetMinPoolSize(uint64(cfg.DBMaxIdleConns))
	clientOptions.SetMaxConnIdleTime(cfg.DBMaxLifetime)

	// Set timeouts
	clientOptions.SetConnectTimeout(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(30 * time.Second)

	// For AWS DocumentDB, we might need TLS configuration
	if cfg.DocumentDBTLSEnabled {
		tlsConfig := &tls.Config{}
		if cfg.DocumentDBTLSInsecure {
			tlsConfig.InsecureSkipVerify = true
		}
		clientOptions.SetTLSConfig(tlsConfig)
	}

	// Create the MongoDB client
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DocumentDB: %w", err)
	}

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping DocumentDB: %w", err)
	}

	logger.Info("Successfully connected to DocumentDB")

	return &DocumentDatabase{
		Client:   client,
		database: cfg.DocumentDBName,
		logger:   logger,
	}, nil
}

// GetDatabase returns the MongoDB database instance
func (d *DocumentDatabase) GetDatabase() *mongo.Database {
	return d.Client.Database(d.database)
}

// GetCollection returns a specific collection from the database
func (d *DocumentDatabase) GetCollection(name string) *mongo.Collection {
	return d.GetDatabase().Collection(name)
}

// Close closes the DocumentDB connection
func (d *DocumentDatabase) Close(ctx context.Context) error {
	if err := d.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from DocumentDB: %w", err)
	}
	d.logger.Info("DocumentDB connection closed")
	return nil
}

// CreateIndexes creates necessary indexes for the video collection
func (d *DocumentDatabase) CreateIndexes(ctx context.Context) error {
	videoCollection := d.GetCollection("videos")

	// Create unique index on video_id
	_, err := videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "video_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create video_id unique index: %w", err)
	}

	// Create index on customer_id
	_, err = videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "customer_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create customer_id index: %w", err)
	}

	// Create index on status
	_, err = videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	// Create compound index on status and customer_id
	_, err = videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}, {Key: "customer_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create compound index: %w", err)
	}

	// Create index on created_at for sorting
	_, err = videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create created_at index: %w", err)
	}

	// Create index on updated_at for sorting
	_, err = videoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "updated_at", Value: -1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create updated_at index: %w", err)
	}

	d.logger.Info("DocumentDB indexes created successfully")
	return nil
}
