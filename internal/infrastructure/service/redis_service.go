package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

type redisService struct {
	client *redis.Client
}

// NewRedisService creates a new Redis service implementation
func NewRedisService(cfg *config.Config) (port.CacheService, error) {
	if !cfg.CacheEnabled {
		return &noopCacheService{}, nil
	}

	var addr string
	if cfg.CacheEndpoint != "" {
		addr = fmt.Sprintf("%s:%d", cfg.CacheEndpoint, cfg.CachePort)
	} else {
		// For local development or when endpoint is not provided
		addr = fmt.Sprintf("localhost:%d", cfg.CachePort)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        "",              // ElastiCache typically uses AUTH token, but for simplicity using empty
		DB:              0,               // default database
		ReadTimeout:     2 * time.Second, // Reduced timeout
		WriteTimeout:    2 * time.Second, // Reduced timeout
		DialTimeout:     3 * time.Second, // Reduced timeout
		MaxRetries:      2,               // Add retries
		MinRetryBackoff: 50 * time.Millisecond,
		MaxRetryBackoff: 500 * time.Millisecond,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("Warning: Redis connection failed, falling back to no-op cache: %v\n", err)
		return &noopCacheService{}, nil // Gracefully degrade to no-op cache
	}

	return &redisService{
		client: rdb,
	}, nil
}

// NewRedisDataSource creates a new Redis data source implementation with advanced operations
func NewRedisDataSource(cfg *config.Config) (port.CacheDataSource, error) {
	if !cfg.CacheEnabled {
		return &noopCacheDataSource{}, nil
	}

	var addr string
	if cfg.CacheEndpoint != "" {
		addr = fmt.Sprintf("%s:%d", cfg.CacheEndpoint, cfg.CachePort)
	} else {
		// For local development or when endpoint is not provided
		addr = fmt.Sprintf("localhost:%d", cfg.CachePort)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // ElastiCache typically uses AUTH token, but for simplicity using empty
		DB:           0,  // default database
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  5 * time.Second,
	})

	return &redisDataSource{
		redisService: &redisService{client: rdb},
	}, nil
}

// Set stores a key-value pair in the cache with optional expiration
func (r *redisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
	}

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Get retrieves a value from the cache by key
func (r *redisService) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key %s not found", key)
		}
		// Check if it's a connection error
		if isConnectionError(err) {
			return nil, fmt.Errorf("cache connection failed for key %s: %w", key, err)
		}
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

// GetString retrieves a string value from the cache by key
func (r *redisService) GetString(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s not found", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

// Delete removes a key-value pair from the cache
func (r *redisService) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

// Exists checks if a key exists in the cache
func (r *redisService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return result > 0, nil
}

// SetExpiration sets an expiration time for an existing key
func (r *redisService) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	return nil
}

// GetTTL gets the time-to-live for a key
func (r *redisService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	return ttl, nil
}

// Flush removes all keys from the cache
func (r *redisService) Flush(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}

	return nil
}

// Ping tests the connection to the cache
func (r *redisService) Ping(ctx context.Context) error {
	err := r.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping cache: %w", err)
	}

	return nil
}

// Close closes the connection to the cache
func (r *redisService) Close() error {
	return r.client.Close()
}

// redisDataSource extends redisService with advanced Redis operations
type redisDataSource struct {
	*redisService
}

// SetHash stores a hash field-value pair
func (r *redisDataSource) SetHash(ctx context.Context, key, field string, value interface{}) error {
	var data interface{}
	switch v := value.(type) {
	case string, int, int64, float64, bool:
		data = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal hash value for key %s field %s: %w", key, field, err)
		}
		data = string(jsonData)
	}

	err := r.client.HSet(ctx, key, field, data).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash field %s in key %s: %w", field, key, err)
	}

	return nil
}

// GetHash retrieves a hash field value
func (r *redisDataSource) GetHash(ctx context.Context, key, field string) (interface{}, error) {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("hash field %s in key %s not found", field, key)
		}
		return nil, fmt.Errorf("failed to get hash field %s in key %s: %w", field, key, err)
	}

	return val, nil
}

// GetAllHash retrieves all hash field-value pairs
func (r *redisDataSource) GetAllHash(ctx context.Context, key string) (map[string]interface{}, error) {
	vals, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all hash fields for key %s: %w", key, err)
	}

	result := make(map[string]interface{})
	for field, value := range vals {
		result[field] = value
	}

	return result, nil
}

// DeleteHash removes a hash field
func (r *redisDataSource) DeleteHash(ctx context.Context, key, field string) error {
	err := r.client.HDel(ctx, key, field).Err()
	if err != nil {
		return fmt.Errorf("failed to delete hash field %s in key %s: %w", field, key, err)
	}

	return nil
}

// SetList appends values to a list
func (r *redisDataSource) SetList(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.RPush(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("failed to append to list %s: %w", key, err)
	}

	return nil
}

// GetList retrieves all values from a list
func (r *redisDataSource) GetList(ctx context.Context, key string) ([]interface{}, error) {
	vals, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list %s: %w", key, err)
	}

	result := make([]interface{}, len(vals))
	for i, val := range vals {
		result[i] = val
	}

	return result, nil
}

// GetListRange retrieves a range of values from a list
func (r *redisDataSource) GetListRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error) {
	vals, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list range %s[%d:%d]: %w", key, start, stop, err)
	}

	result := make([]interface{}, len(vals))
	for i, val := range vals {
		result[i] = val
	}

	return result, nil
}

// SetAdd adds members to a set
func (r *redisDataSource) SetAdd(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add members to set %s: %w", key, err)
	}

	return nil
}

// SetMembers retrieves all members from a set
func (r *redisDataSource) SetMembers(ctx context.Context, key string) ([]interface{}, error) {
	vals, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get set members %s: %w", key, err)
	}

	result := make([]interface{}, len(vals))
	for i, val := range vals {
		result[i] = val
	}

	return result, nil
}

// SetIsMember checks if a member exists in a set
func (r *redisDataSource) SetIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	result, err := r.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check set membership %s: %w", key, err)
	}

	return result, nil
}

// noopCacheService is a no-operation implementation when caching is disabled
type noopCacheService struct{}

func (n *noopCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (n *noopCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheService) GetString(ctx context.Context, key string) (string, error) {
	return "", fmt.Errorf("cache disabled")
}

func (n *noopCacheService) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *noopCacheService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *noopCacheService) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (n *noopCacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (n *noopCacheService) Flush(ctx context.Context) error {
	return nil
}

func (n *noopCacheService) Ping(ctx context.Context) error {
	return fmt.Errorf("cache disabled")
}

func (n *noopCacheService) Close() error {
	return nil
}

// noopCacheDataSource is a no-operation implementation when caching is disabled
type noopCacheDataSource struct {
	*noopCacheService
}

func (n *noopCacheDataSource) SetHash(ctx context.Context, key, field string, value interface{}) error {
	return nil
}

func (n *noopCacheDataSource) GetHash(ctx context.Context, key, field string) (interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheDataSource) GetAllHash(ctx context.Context, key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheDataSource) DeleteHash(ctx context.Context, key, field string) error {
	return nil
}

func (n *noopCacheDataSource) SetList(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (n *noopCacheDataSource) GetList(ctx context.Context, key string) ([]interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheDataSource) GetListRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheDataSource) SetAdd(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (n *noopCacheDataSource) SetMembers(ctx context.Context, key string) ([]interface{}, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noopCacheDataSource) SetIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, nil
}

// isConnectionError checks if the error is related to connection issues
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "no route to host") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "connection reset")
}
