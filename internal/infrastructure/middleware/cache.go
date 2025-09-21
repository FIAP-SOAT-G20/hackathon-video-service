package middleware

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
)

// CacheStore holds cache middleware configuration
type CacheStore struct {
	Endpoint string
	Port     int
	Duration time.Duration
	Logger   *logger.Logger
}

func NewCacheStore(cfg *config.Config, logger *logger.Logger) *CacheStore {
	// Handle nil config for testing scenarios
	if cfg == nil {
		return &CacheStore{
			Endpoint: "localhost",
			Port:     6379,
			Duration: time.Minute,
			Logger:   logger,
		}
	}
	
	return &CacheStore{
		Endpoint: cfg.CacheEndpoint,
		Port:     cfg.CachePort,
		Duration: cfg.CacheDuration,
		Logger:   logger,
	}
}

// NewMemoryCacheStore creates a new in-memory cache store
func (cs *CacheStore) NewMemoryCacheStore() port.CacheStore {
	return persistence.NewInMemoryStore(time.Minute)
}

// NewRedisCacheStore creates a new Redis cache store
func (cs *CacheStore) NewRedisCacheStore(logger *logger.Logger) port.CacheStore {

	// Create Redis store with host:port format
	host := fmt.Sprintf("%s:%d", cs.Endpoint, cs.Port)
	store := persistence.NewRedisCache(host, "", time.Minute)

	// Test Redis connection
	if err := store.Set("cache_test_key", "test", time.Second*10); err != nil {
		logger.Error("failed to connect to Redis cache, falling back to in-memory cache", "error", err.Error())
		return cs.NewMemoryCacheStore()
	}

	logger.Info("Connected to Redis cache", "host", host)

	return store
}

// CachePage creates a cache middleware specifically for video list endpoint
func CachePage(store port.CacheStore, next gin.HandlerFunc) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute*5, func(c *gin.Context) {
		// Cache miss - call the next handler to generate the response
		c.Set("cache_miss", true)
		next(c)
	})
}
