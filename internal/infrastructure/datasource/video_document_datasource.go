package datasource

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/entity"
	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/database"
)

type videoDocumentDataSource struct {
	db         *database.DocumentDatabase
	collection *mongo.Collection
}

// VideoDocument represents the MongoDB document structure for video
type VideoDocument struct {
	ID          primitive.ObjectID      `bson:"_id,omitempty"`
	VideoID     uint64                  `bson:"video_id"`
	UserID      uint64                  `bson:"user_id"`
	Name        string                  `bson:"name"`
	Description string                  `bson:"description"`
	Status      valueobject.VideoStatus `bson:"status"`
	CreatedAt   time.Time               `bson:"created_at"`
	UpdatedAt   time.Time               `bson:"updated_at"`
}

func NewVideoDocumentDataSource(db *database.DocumentDatabase) port.VideoDataSource {
	collection := db.GetCollection("videos")
	return &videoDocumentDataSource{
		db:         db,
		collection: collection,
	}
}

func (ds *videoDocumentDataSource) FindByID(ctx context.Context, id uint64) (*entity.Video, error) {
	filter := bson.M{"video_id": id}

	var doc VideoDocument
	err := ds.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding video: %w", err)
	}

	return ds.documentToEntity(&doc), nil
}

func (ds *videoDocumentDataSource) FindAll(ctx context.Context, filters map[string]any, sort string, page, limit int) ([]*entity.Video, int64, error) {
	// Build MongoDB filter
	filter := bson.M{}

	for key, value := range filters {
		switch key {
		case "statuses":
			if statuses, ok := value.([]valueobject.VideoStatus); ok && len(statuses) > 0 {
				filter["status"] = bson.M{"$in": statuses}
			}
		case "statuses_exclude":
			if statuses, ok := value.([]valueobject.VideoStatus); ok && len(statuses) > 0 {
				filter["status"] = bson.M{"$nin": statuses}
			}
		case "user_id":
			if customerID, ok := value.(uint64); ok && customerID != 0 {
				filter["user_id"] = customerID
			}
		}
	}

	// Count total documents
	total, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting videos: %w", err)
	}

	// Build sort options
	sortOptions := options.Find()
	if sort != "" {
		// Convert SQL-style sort to MongoDB sort
		sortField := "created_at"
		sortOrder := -1 // Default descending

		switch sort {
		case "id asc", "video_id asc":
			sortField = "video_id"
			sortOrder = 1
		case "id desc", "video_id desc":
			sortField = "video_id"
			sortOrder = -1
		case "created_at asc":
			sortField = "created_at"
			sortOrder = 1
		case "created_at desc":
			sortField = "created_at"
			sortOrder = -1
		case "updated_at asc":
			sortField = "updated_at"
			sortOrder = 1
		case "updated_at desc":
			sortField = "updated_at"
			sortOrder = -1
		}

		sortOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	} else {
		// Default sort by created_at desc
		sortOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	// Apply pagination
	skip := (page - 1) * limit
	sortOptions.SetSkip(int64(skip))
	sortOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := ds.collection.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding videos: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx) // Ignore error for cleanup operation
	}()

	var documents []VideoDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, 0, fmt.Errorf("error decoding videos: %w", err)
	}

	// Convert documents to entities
	videos := make([]*entity.Video, len(documents))
	for i, doc := range documents {
		videos[i] = ds.documentToEntity(&doc)
	}

	return videos, total, nil
}

func (ds *videoDocumentDataSource) Create(ctx context.Context, video *entity.Video) error {
	doc := ds.entityToDocument(video)
	doc.ID = ds.uint64ToObjectID(video.ID)
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	result, err := ds.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("error creating video: %w", err)
	}

	// If no ID was set, generate one from the inserted ObjectID
	if video.ID == 0 {
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			// Convert ObjectID to uint64 (simplified approach)
			video.ID = ds.objectIDToUint64(oid)
			// Update the document with the generated video_id
			filter := bson.M{"_id": oid}
			update := bson.M{"$set": bson.M{"video_id": video.ID}}
			_, err := ds.collection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("error updating video with generated ID: %w", err)
			}
		}
	}

	return nil
}

func (ds *videoDocumentDataSource) Update(ctx context.Context, video *entity.Video) error {
	filter := bson.M{"video_id": video.ID}

	update := bson.M{
		"$set": bson.M{
			"video_id":    video.ID,
			"user_id":     video.UserID,
			"status":      video.Status,
			"name":        video.Name,
			"description": video.Description,
			"hash":        video.Hash,
			"updated_at":  time.Now(),
		},
	}

	result, err := ds.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating video: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil // No document found, but no error
	}

	return nil
}

func (ds *videoDocumentDataSource) Delete(ctx context.Context, id uint64) error {
	filter := bson.M{"video_id": id}

	result, err := ds.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting video: %w", err)
	}

	if result.DeletedCount == 0 {
		return nil // No document found, but no error
	}

	return nil
}

func (ds *videoDocumentDataSource) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// MongoDB transactions require a session
	session, err := ds.db.StartSession()
	if err != nil {
		return fmt.Errorf("error starting session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute the transaction
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

// Helper methods for conversion between entity and document

func (ds *videoDocumentDataSource) entityToDocument(video *entity.Video) *VideoDocument {
	return &VideoDocument{
		ID:          ds.uint64ToObjectID(video.ID),
		VideoID:     video.ID,
		UserID:      video.UserID,
		Name:        video.Name,
		Description: video.Description,
		Status:      video.Status,
		CreatedAt:   video.CreatedAt,
		UpdatedAt:   video.UpdatedAt,
	}
}

func (ds *videoDocumentDataSource) documentToEntity(doc *VideoDocument) *entity.Video {
	// Use VideoID field directly, fall back to converting ObjectID if VideoID is 0
	id := doc.VideoID
	if id == 0 {
		id = ds.objectIDToUint64(doc.ID)
	}

	return &entity.Video{
		ID:          id,
		UserID:      doc.UserID,
		Name:        doc.Name,
		Description: doc.Description,
		Status:      doc.Status,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}

// objectIDToUint64 converts MongoDB ObjectID to uint64
// This is a simplified approach - in production you might want a more sophisticated ID strategy
func (ds *videoDocumentDataSource) objectIDToUint64(oid primitive.ObjectID) uint64 {
	// Take the last 8 bytes of the ObjectID and convert to uint64
	bytes := [8]byte{}
	copy(bytes[:], oid[4:])

	// Convert to uint64
	var result uint64
	for i, b := range bytes {
		result |= uint64(b) << (8 * (7 - i))
	}

	// Ensure it's a positive number and within reasonable bounds
	return result & 0x7FFFFFFFFFFFFFFF
}

// uint64ToObjectID converts uint64 back to ObjectID (if needed)
func (ds *videoDocumentDataSource) uint64ToObjectID(id uint64) primitive.ObjectID {
	// This is a simplified reverse conversion
	bytes := make([]byte, 12)

	// Use current timestamp for first 4 bytes
	timestamp := uint32(time.Now().Unix())
	bytes[0] = byte(timestamp >> 24)
	bytes[1] = byte(timestamp >> 16)
	bytes[2] = byte(timestamp >> 8)
	bytes[3] = byte(timestamp)

	// Use id for last 8 bytes
	for i := 0; i < 8; i++ {
		bytes[4+i] = byte(id >> (8 * (7 - i)))
	}

	return primitive.ObjectID(bytes)
}
