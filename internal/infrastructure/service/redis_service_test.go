package service

import (
	"context"
	"fmt"
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

// Test Redis service connection failures and error handling
func TestRedisService_ConnectionErrors(t *testing.T) {
	t.Run("should handle connection failure gracefully", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "nonexistent-redis-host",
			CachePort:     6379,
		}

		// This should return a noopCacheService due to connection failure
		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)

		// Verify it's a noop service
		_, ok := service.(*noopCacheService)
		assert.True(t, ok, "should return noopCacheService when connection fails")
	})

	t.Run("should handle localhost fallback when endpoint not provided", func(t *testing.T) {
		cfg := &config.Config{
			CacheEnabled:  true,
			CacheEndpoint: "", // Empty endpoint should fallback to localhost
			CachePort:     6379,
		}

		service, err := NewRedisService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

// Test Redis service operations with error scenarios
func TestRedisService_Operations_WithErrors(t *testing.T) {
	// Start a mini redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create redis client directly for testing
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	service := &redisService{client: rdb}
	ctx := context.Background()

	t.Run("Set operation with marshalling error", func(t *testing.T) {
		// Create a value that can't be marshalled
		unmarshalable := make(chan int)

		err := service.Set(ctx, "test-key", unmarshalable, time.Minute)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal")
	})

	t.Run("Get operation with connection error", func(t *testing.T) {
		// Stop the mini redis server to simulate connection error
		mr.Close()

		_, err := service.Get(ctx, "test-key")
		assert.Error(t, err)
	})

	t.Run("GetString operation", func(t *testing.T) {
		// Restart mini redis for this test
		mr2, err := miniredis.Run()
		require.NoError(t, err)
		defer mr2.Close()

		rdb2 := redis.NewClient(&redis.Options{
			Addr: mr2.Addr(),
		})
		service2 := &redisService{client: rdb2}

		// Test successful GetString
		err = service2.Set(ctx, "string-key", "test-value", time.Minute)
		require.NoError(t, err)

		value, err := service2.GetString(ctx, "string-key")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", value)

		// Test GetString with non-existent key
		_, err = service2.GetString(ctx, "non-existent-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete operation", func(t *testing.T) {
		mr3, err := miniredis.Run()
		require.NoError(t, err)
		defer mr3.Close()

		rdb3 := redis.NewClient(&redis.Options{
			Addr: mr3.Addr(),
		})
		service3 := &redisService{client: rdb3}

		// Set a key first
		err = service3.Set(ctx, "delete-key", "value", time.Minute)
		require.NoError(t, err)

		// Delete the key
		err = service3.Delete(ctx, "delete-key")
		assert.NoError(t, err)
	})

	t.Run("Exists operation", func(t *testing.T) {
		mr4, err := miniredis.Run()
		require.NoError(t, err)
		defer mr4.Close()

		rdb4 := redis.NewClient(&redis.Options{
			Addr: mr4.Addr(),
		})
		service4 := &redisService{client: rdb4}

		// Test non-existent key
		exists, err := service4.Exists(ctx, "non-existent")
		assert.NoError(t, err)
		assert.False(t, exists)

		// Set a key and test existence
		err = service4.Set(ctx, "exists-key", "value", time.Minute)
		require.NoError(t, err)

		exists, err = service4.Exists(ctx, "exists-key")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("SetExpiration operation", func(t *testing.T) {
		mr5, err := miniredis.Run()
		require.NoError(t, err)
		defer mr5.Close()

		rdb5 := redis.NewClient(&redis.Options{
			Addr: mr5.Addr(),
		})
		service5 := &redisService{client: rdb5}

		// Set a key first
		err = service5.Set(ctx, "expire-key", "value", time.Hour)
		require.NoError(t, err)

		// Set expiration
		err = service5.SetExpiration(ctx, "expire-key", time.Minute)
		assert.NoError(t, err)
	})

	t.Run("GetTTL operation", func(t *testing.T) {
		mr6, err := miniredis.Run()
		require.NoError(t, err)
		defer mr6.Close()

		rdb6 := redis.NewClient(&redis.Options{
			Addr: mr6.Addr(),
		})
		service6 := &redisService{client: rdb6}

		// Set a key with expiration
		err = service6.Set(ctx, "ttl-key", "value", time.Hour)
		require.NoError(t, err)

		// Get TTL
		ttl, err := service6.GetTTL(ctx, "ttl-key")
		assert.NoError(t, err)
		assert.Greater(t, ttl, time.Duration(0))
	})

	t.Run("Flush operation", func(t *testing.T) {
		mr7, err := miniredis.Run()
		require.NoError(t, err)
		defer mr7.Close()

		rdb7 := redis.NewClient(&redis.Options{
			Addr: mr7.Addr(),
		})
		service7 := &redisService{client: rdb7}

		// Set some keys
		err = service7.Set(ctx, "flush-key1", "value1", time.Hour)
		require.NoError(t, err)
		err = service7.Set(ctx, "flush-key2", "value2", time.Hour)
		require.NoError(t, err)

		// Flush all keys
		err = service7.Flush(ctx)
		assert.NoError(t, err)

		// Verify keys are gone
		exists, err := service7.Exists(ctx, "flush-key1")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Ping operation", func(t *testing.T) {
		mr8, err := miniredis.Run()
		require.NoError(t, err)
		defer mr8.Close()

		rdb8 := redis.NewClient(&redis.Options{
			Addr: mr8.Addr(),
		})
		service8 := &redisService{client: rdb8}

		// Test ping
		err = service8.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Close operation", func(t *testing.T) {
		mr9, err := miniredis.Run()
		require.NoError(t, err)
		defer mr9.Close()

		rdb9 := redis.NewClient(&redis.Options{
			Addr: mr9.Addr(),
		})
		service9 := &redisService{client: rdb9}

		// Test close
		err = service9.Close()
		assert.NoError(t, err)
	})
}

// Test Redis data source list operations separately
func TestRedisDataSource_ListOperations(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cfg := &config.Config{
		CacheEnabled:  true,
		CacheEndpoint: mr.Host(),
		CachePort:     6379,
	}

	// Create Redis data source
	ds, err := NewRedisDataSource(cfg)
	require.NoError(t, err)
	require.NotNil(t, ds)

	ctx := context.Background()
	listKey := fmt.Sprintf("list-operations-key-%d", time.Now().UnixNano())

	// SetList (append to list)
	err = ds.SetList(ctx, listKey, "item1", "item2", "item3")
	assert.NoError(t, err)

	// GetList
	items, err := ds.GetList(ctx, listKey)
	assert.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "item1", items[0])

	// GetListRange
	rangeItems, err := ds.GetListRange(ctx, listKey, 0, 1)
	assert.NoError(t, err)
	assert.Len(t, rangeItems, 2)
	assert.Equal(t, "item1", rangeItems[0])
	assert.Equal(t, "item2", rangeItems[1])
}

// Test Redis data source advanced operations
func TestRedisDataSource_AdvancedOperations(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cfg := &config.Config{
		CacheEnabled:  true,
		CacheEndpoint: mr.Host(),
		CachePort:     6379,
	}

	// Create Redis data source
	ds, err := NewRedisDataSource(cfg)
	require.NoError(t, err)
	require.NotNil(t, ds)

	ctx := context.Background()

	t.Run("Hash operations", func(t *testing.T) {
		// SetHash with different data types
		err := ds.SetHash(ctx, "hash-key", "field1", "string-value")
		assert.NoError(t, err)

		err = ds.SetHash(ctx, "hash-key", "field2", 123)
		assert.NoError(t, err)

		err = ds.SetHash(ctx, "hash-key", "field3", map[string]interface{}{"nested": "value"})
		assert.NoError(t, err)

		// GetHash
		value, err := ds.GetHash(ctx, "hash-key", "field1")
		assert.NoError(t, err)
		assert.Equal(t, "string-value", value)

		// GetHash non-existent field
		_, err = ds.GetHash(ctx, "hash-key", "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		// GetAllHash
		allFields, err := ds.GetAllHash(ctx, "hash-key")
		assert.NoError(t, err)
		assert.Len(t, allFields, 3)

		// DeleteHash
		err = ds.DeleteHash(ctx, "hash-key", "field1")
		assert.NoError(t, err)

		// Verify deletion
		_, err = ds.GetHash(ctx, "hash-key", "field1")
		assert.Error(t, err)
	})

	t.Run("Set operations", func(t *testing.T) {
		// SetAdd
		err := ds.SetAdd(ctx, "set-key", "member1", "member2", "member3")
		assert.NoError(t, err)

		// SetMembers
		members, err := ds.SetMembers(ctx, "set-key")
		assert.NoError(t, err)
		assert.Len(t, members, 3)

		// SetIsMember
		isMember, err := ds.SetIsMember(ctx, "set-key", "member1")
		assert.NoError(t, err)
		assert.True(t, isMember)

		isMember, err = ds.SetIsMember(ctx, "set-key", "non-member")
		assert.NoError(t, err)
		assert.False(t, isMember)
	})
}

// Test no-op implementations
func TestNoopCacheImplementations(t *testing.T) {
	ctx := context.Background()

	t.Run("noopCacheService operations", func(t *testing.T) {
		noop := &noopCacheService{}

		// All operations should not return errors except where specified
		err := noop.Set(ctx, "key", "value", time.Minute)
		assert.NoError(t, err)

		_, err = noop.Get(ctx, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		_, err = noop.GetString(ctx, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		err = noop.Delete(ctx, "key")
		assert.NoError(t, err)

		exists, err := noop.Exists(ctx, "key")
		assert.NoError(t, err)
		assert.False(t, exists)

		err = noop.SetExpiration(ctx, "key", time.Minute)
		assert.NoError(t, err)

		ttl, err := noop.GetTTL(ctx, "key")
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), ttl)

		err = noop.Flush(ctx)
		assert.NoError(t, err)

		err = noop.Ping(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		err = noop.Close()
		assert.NoError(t, err)
	})

	t.Run("noopCacheDataSource operations", func(t *testing.T) {
		noop := &noopCacheDataSource{noopCacheService: &noopCacheService{}}

		err := noop.SetHash(ctx, "key", "field", "value")
		assert.NoError(t, err)

		_, err = noop.GetHash(ctx, "key", "field")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		_, err = noop.GetAllHash(ctx, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		err = noop.DeleteHash(ctx, "key", "field")
		assert.NoError(t, err)

		err = noop.SetList(ctx, "key", "item1", "item2")
		assert.NoError(t, err)

		_, err = noop.GetList(ctx, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		_, err = noop.GetListRange(ctx, "key", 0, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		err = noop.SetAdd(ctx, "key", "member1", "member2")
		assert.NoError(t, err)

		_, err = noop.SetMembers(ctx, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache disabled")

		isMember, err := noop.SetIsMember(ctx, "key", "member")
		assert.NoError(t, err)
		assert.False(t, isMember)
	})
}

// Test connection error detection
func TestIsConnectionError(t *testing.T) {
	t.Run("should detect connection errors", func(t *testing.T) {
		connectionErrors := []string{
			"connection refused",
			"i/o timeout",
			"no route to host",
			"network is unreachable",
			"dial tcp failed",
			"connection reset by peer",
		}

		for _, errMsg := range connectionErrors {
			err := fmt.Errorf("%s", errMsg)
			assert.True(t, isConnectionError(err), "should detect '%s' as connection error", errMsg)
		}
	})

	t.Run("should not detect non-connection errors", func(t *testing.T) {
		nonConnectionErrors := []string{
			"key not found",
			"invalid command",
			"permission denied",
		}

		for _, errMsg := range nonConnectionErrors {
			err := fmt.Errorf("%s", errMsg)
			assert.False(t, isConnectionError(err), "should not detect '%s' as connection error", errMsg)
		}
	})

	t.Run("should handle nil error", func(t *testing.T) {
		assert.False(t, isConnectionError(nil))
	})
}
