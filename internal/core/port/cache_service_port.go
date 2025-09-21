package port

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	// Set stores a key-value pair in the cache with optional expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get retrieves a value from the cache by key
	Get(ctx context.Context, key string) (interface{}, error)

	// GetString retrieves a string value from the cache by key
	GetString(ctx context.Context, key string) (string, error)

	// Delete removes a key-value pair from the cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// SetExpiration sets an expiration time for an existing key
	SetExpiration(ctx context.Context, key string, expiration time.Duration) error

	// GetTTL gets the time-to-live for a key
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	// Flush removes all keys from the cache
	Flush(ctx context.Context) error

	// Ping tests the connection to the cache
	Ping(ctx context.Context) error

	// Close closes the connection to the cache
	Close() error
}

// CacheDataSource defines the interface for cache data operations (extending basic cache with advanced operations)
type CacheDataSource interface {
	CacheService

	// SetHash stores a hash field-value pair
	SetHash(ctx context.Context, key, field string, value interface{}) error

	// GetHash retrieves a hash field value
	GetHash(ctx context.Context, key, field string) (interface{}, error)

	// GetAllHash retrieves all hash field-value pairs
	GetAllHash(ctx context.Context, key string) (map[string]interface{}, error)

	// DeleteHash removes a hash field
	DeleteHash(ctx context.Context, key, field string) error

	// SetList appends values to a list
	SetList(ctx context.Context, key string, values ...interface{}) error

	// GetList retrieves all values from a list
	GetList(ctx context.Context, key string) ([]interface{}, error)

	// GetListRange retrieves a range of values from a list
	GetListRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error)

	// SetAdd adds members to a set
	SetAdd(ctx context.Context, key string, members ...interface{}) error

	// SetMembers retrieves all members from a set
	SetMembers(ctx context.Context, key string) ([]interface{}, error)

	// SetIsMember checks if a member exists in a set
	SetIsMember(ctx context.Context, key string, member interface{}) (bool, error)
}
