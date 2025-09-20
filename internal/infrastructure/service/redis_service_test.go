package service

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

func TestNewRedisService(t *testing.T) {
	t.Run("should return noop cache service when cache is disabled", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled: false,
		}

		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)

		// Verify it's a noop service
		_, ok := service.(*noopCacheService)
		assert.True(t, ok, "should return noopCacheService when cache is disabled")
	})

	t.Run("should create redis service with endpoint", func(t *testing.T) {
		// Start a mini redis server for testing
		mr, err := miniredis.Run()
		require.NoError(t, err)
		defer mr.Close()

		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: mr.Host(),
			CachePort:     6379, // Use integer
		}

		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("should fallback to noop service when redis connection fails", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "nonexistent-host",
			CachePort:     6379,
		}

		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)

		// Should fallback to noop service
		_, ok := service.(*noopCacheService)
		assert.True(t, ok, "should fallback to noopCacheService when connection fails")
	})

	t.Run("should use localhost when endpoint is empty", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "", // Empty endpoint
			CachePort:     6379,
		}

		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestNewRedisDataSource(t *testing.T) {
	t.Run("should return noop cache data source when cache is disabled", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled: false,
		}

		dataSource, err := NewRedisDataSource(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, dataSource)

		// Verify it's a noop data source
		_, ok := dataSource.(*noopCacheDataSource)
		assert.True(t, ok, "should return noopCacheDataSource when cache is disabled")
	})

	t.Run("should create redis data source with endpoint", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "localhost",
			CachePort:     6379,
		}

		dataSource, err := NewRedisDataSource(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, dataSource)

		// Verify it's a redis data source
		_, ok := dataSource.(*redisDataSource)
		assert.True(t, ok, "should return redisDataSource when cache is enabled")
	})

	t.Run("should use localhost when endpoint is empty", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "", // Empty endpoint
			CachePort:     6379,
		}

		dataSource, err := NewRedisDataSource(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, dataSource)
	})
}

func TestRedisService_Set(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should set string value successfully", func(t *testing.T) {
		err := service.Set(ctx, "test:string", "test-value", time.Minute)
		assert.NoError(t, err)

		// Verify value was set
		val, err := rdb.Get(ctx, "test:string").Result()
		assert.NoError(t, err)
		assert.Equal(t, "test-value", val)
	})

	t.Run("should set byte slice value successfully", func(t *testing.T) {
		data := []byte("test-bytes")
		err := service.Set(ctx, "test:bytes", data, time.Minute)
		assert.NoError(t, err)

		// Verify value was set
		val, err := rdb.Get(ctx, "test:bytes").Result()
		assert.NoError(t, err)
		assert.Equal(t, string(data), val)
	})

	t.Run("should marshal and set complex object successfully", func(t *testing.T) {
		obj := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		}
		err := service.Set(ctx, "test:object", obj, time.Minute)
		assert.NoError(t, err)

		// Verify JSON was stored
		val, err := rdb.Get(ctx, "test:object").Result()
		assert.NoError(t, err)
		assert.Contains(t, val, "key1")
		assert.Contains(t, val, "value1")
	})

	t.Run("should handle marshal error gracefully", func(t *testing.T) {
		// Use a channel which cannot be marshaled to JSON
		invalidObj := make(chan int)
		err := service.Set(ctx, "test:invalid", invalidObj, time.Minute)
		assert.Error(t, err)
	})

	t.Run("should set with no expiration", func(t *testing.T) {
		err := service.Set(ctx, "test:no-expiry", "test-value", 0)
		assert.NoError(t, err)

		// Verify TTL is -1 (no expiration)
		ttl, err := rdb.TTL(ctx, "test:no-expiry").Result()
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(-1), ttl)
	})
}

func TestRedisService_Get(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should get existing value successfully", func(t *testing.T) {
		// Set a value first
		err := rdb.Set(ctx, "test:get", "test-value", 0).Err()
		require.NoError(t, err)

		val, err := service.Get(ctx, "test:get")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", val)
	})

	t.Run("should return nil for non-existent key", func(t *testing.T) {
		val, err := service.Get(ctx, "test:nonexistent")
		assert.Error(t, err) // miniredis returns an error for non-existent keys
		assert.Nil(t, val)
	})
}

func TestRedisService_GetString(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should get string value successfully", func(t *testing.T) {
		// Set a value first
		err := rdb.Set(ctx, "test:string", "test-value", 0).Err()
		require.NoError(t, err)

		val, err := service.GetString(ctx, "test:string")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", val)
	})

	t.Run("should return empty string for non-existent key", func(t *testing.T) {
		val, err := service.GetString(ctx, "test:nonexistent")
		assert.Error(t, err)
		assert.Equal(t, "", val)
		// Note: miniredis behavior may differ from redis.Nil
	})
}

func TestRedisService_Delete(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should delete existing key successfully", func(t *testing.T) {
		// Set a value first
		err := rdb.Set(ctx, "test:delete", "test-value", 0).Err()
		require.NoError(t, err)

		err = service.Delete(ctx, "test:delete")
		assert.NoError(t, err)

		// Verify key was deleted
		exists, err := rdb.Exists(ctx, "test:delete").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("should handle deleting non-existent key", func(t *testing.T) {
		err := service.Delete(ctx, "test:nonexistent")
		assert.NoError(t, err) // Redis doesn't error when deleting non-existent keys
	})
}

func TestRedisService_Exists(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should return true for existing key", func(t *testing.T) {
		// Set a value first
		err := rdb.Set(ctx, "test:exists", "test-value", 0).Err()
		require.NoError(t, err)

		exists, err := service.Exists(ctx, "test:exists")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false for non-existent key", func(t *testing.T) {
		exists, err := service.Exists(ctx, "test:nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestRedisService_SetExpiration(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should set expiration for existing key", func(t *testing.T) {
		// Set a value first without expiration
		err := rdb.Set(ctx, "test:expire", "test-value", 0).Err()
		require.NoError(t, err)

		err = service.SetExpiration(ctx, "test:expire", time.Minute)
		assert.NoError(t, err)

		// Verify TTL was set
		ttl, err := rdb.TTL(ctx, "test:expire").Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= time.Minute)
	})

	t.Run("should handle setting expiration for non-existent key", func(t *testing.T) {
		err := service.SetExpiration(ctx, "test:nonexistent", time.Minute)
		assert.NoError(t, err) // Redis doesn't error, but returns 0
	})
}

func TestRedisService_GetTTL(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Logf("Failed to close Redis client: %v", err)
		}
	}()

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("should get TTL for key with expiration", func(t *testing.T) {
		// Set a value with expiration
		err := rdb.Set(ctx, "test:ttl", "test-value", time.Minute).Err()
		require.NoError(t, err)

		ttl, err := service.GetTTL(ctx, "test:ttl")
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= time.Minute)
	})

	t.Run("should return -1 for key with no expiration", func(t *testing.T) {
		// Set a value without expiration
		err := rdb.Set(ctx, "test:no-ttl", "test-value", 0).Err()
		require.NoError(t, err)

		ttl, err := service.GetTTL(ctx, "test:no-ttl")
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(-1), ttl)
	})

	t.Run("should return -2 for non-existent key", func(t *testing.T) {
		ttl, err := service.GetTTL(ctx, "test:nonexistent")
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(-2), ttl)
	})
}

// Test noop implementations
func TestNoopCacheService(t *testing.T) {
	service := &noopCacheService{}
	ctx := context.Background()

	t.Run("all operations should return nil/false/errors", func(t *testing.T) {
		err := service.Set(ctx, "key", "value", time.Minute)
		assert.NoError(t, err)

		val, err := service.Get(ctx, "key")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Nil(t, val)

		str, err := service.GetString(ctx, "key")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Equal(t, "", str)

		err = service.Delete(ctx, "key")
		assert.NoError(t, err)

		exists, err := service.Exists(ctx, "key")
		assert.NoError(t, err)
		assert.False(t, exists)

		err = service.SetExpiration(ctx, "key", time.Minute)
		assert.NoError(t, err)

		ttl, err := service.GetTTL(ctx, "key")
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), ttl)
	})
}

func TestNoopCacheDataSource(t *testing.T) {
	dataSource := &noopCacheDataSource{}
	ctx := context.Background()

	t.Run("all basic cache operations should return appropriate defaults/errors", func(t *testing.T) {
		err := dataSource.Set(ctx, "key", "value", time.Minute)
		assert.NoError(t, err)

		val, err := dataSource.Get(ctx, "key")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Nil(t, val)

		err = dataSource.Delete(ctx, "key")
		assert.NoError(t, err)

		exists, err := dataSource.Exists(ctx, "key")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("all hash operations should return appropriate defaults/errors", func(t *testing.T) {
		err := dataSource.SetHash(ctx, "hash", "field", "value")
		assert.NoError(t, err)

		hashVal, err := dataSource.GetHash(ctx, "hash", "field")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Nil(t, hashVal)

		hashVals, err := dataSource.GetAllHash(ctx, "hash")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Empty(t, hashVals)

		err = dataSource.DeleteHash(ctx, "hash", "field")
		assert.NoError(t, err)
	})

	t.Run("all list operations should return appropriate defaults/errors", func(t *testing.T) {
		err := dataSource.SetList(ctx, "list", "value")
		assert.NoError(t, err)

		listVals, err := dataSource.GetList(ctx, "list")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Empty(t, listVals)

		listRange, err := dataSource.GetListRange(ctx, "list", 0, -1)
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Empty(t, listRange)
	})

	t.Run("all set operations should return appropriate defaults/errors", func(t *testing.T) {
		err := dataSource.SetAdd(ctx, "set", "member")
		assert.NoError(t, err)

		setMembers, err := dataSource.SetMembers(ctx, "set")
		assert.Error(t, err) // noop returns "cache disabled" error
		assert.Empty(t, setMembers)

		isMember, err := dataSource.SetIsMember(ctx, "set", "member")
		assert.NoError(t, err) // This doesn't error in noop
		assert.False(t, isMember)
	})
}

// Verify interface implementations
func TestInterfaceImplementations(t *testing.T) {
	t.Run("redisService implements CacheService", func(t *testing.T) {
		var _ port.CacheService = (*redisService)(nil)
	})

	t.Run("noopCacheService implements CacheService", func(t *testing.T) {
		var _ port.CacheService = (*noopCacheService)(nil)
	})

	t.Run("redisDataSource implements CacheDataSource", func(t *testing.T) {
		var _ port.CacheDataSource = (*redisDataSource)(nil)
	})

	t.Run("noopCacheDataSource implements CacheDataSource", func(t *testing.T) {
		var _ port.CacheDataSource = (*noopCacheDataSource)(nil)
	})
}
