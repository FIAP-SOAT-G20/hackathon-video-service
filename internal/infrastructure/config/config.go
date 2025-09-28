package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// AWS SQS settings
	AWS_SQS_VideoUpdatedURL             string
	AWS_SQS_VideoUpdatedMaxMessages     int
	AWS_SQS_VideoUpdatedWaitTimeSeconds int

	// AWS S3 settings
	AWSS3BucketName             string
	AWSS3BucketRawFolder        string
	AWSS3BucketProcessedFolder  string
	AWSS3Region                 string
	AWSS3PresignedURLExpiration time.Duration

	// Cache settings
	CacheEndpoint string
	CachePort     int
	CacheEnabled  bool
	CacheDuration time.Duration

	// Database settings
	DBDSN          string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  time.Duration

	DBEngine string

	// DocumentDB settings
	DocumentDBURI         string
	DocumentDBName        string
	DocumentDBTLSEnabled  bool
	DocumentDBTLSCertPath string
	DocumentDBTLSInsecure bool

	// Server settings
	ServerPort                    string
	ServerReadTimeout             time.Duration
	ServerWriteTimeout            time.Duration
	ServerIdleTimeout             time.Duration
	ServerGracefulShutdownTimeout time.Duration

	// Environment
	Environment string

	// JWT Settings
	JWTSecret     string
	JWTExpiration time.Duration

	// Auth0 Settings
	Auth0Domain   string
	Auth0Audience string
	Auth0JWKSURL  string

	// Metrics settings
	MetricsPort string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	AWS_SQS_VideoUpdatedMaxMessages, _ := strconv.Atoi(getEnv("AWS_SQS_VIDEO_UPDATED_MAX_MESSAGES", "10"))
	AWS_SQS_VideoUpdatedWaitTimeSeconds, _ := strconv.Atoi(getEnv("AWS_SQS_VIDEO_UPDATED_WAIT_TIME_SECONDS", "20"))

	dbMaxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	dbMaxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	dbMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	// Parse DocumentDB TLS settings
	documentDBTLSEnabled, _ := strconv.ParseBool(getEnv("DOCUMENTDB_TLS_ENABLED", "false"))
	documentDBTLSInsecure, _ := strconv.ParseBool(getEnv("DOCUMENTDB_TLS_INSECURE", "false"))

	serverReadTimeout, _ := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", "10s"))
	serverWriteTimeout, _ := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", "10s"))
	serverIdleTimeout, _ := time.ParseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s"))
	serverGracefulShutdownTimeout, _ := time.ParseDuration(getEnv("SERVER_GRACEFUL_SHUTDOWN_SEC_TIMEOUT", "5s"))

	jwtExpirationStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExpiration, err := time.ParseDuration(jwtExpirationStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION value %q: %v. Using default value 24h.", jwtExpirationStr, err)
		jwtExpiration = 24 * time.Hour
	}

	s3PresignedURLExpirationStr := getEnv("AWS_S3_PRESIGNED_URL_EXPIRATION", "15m")
	s3PresignedURLExpiration, err := time.ParseDuration(s3PresignedURLExpirationStr)
	if err != nil {
		log.Printf("Warning: invalid AWS_S3_PRESIGNED_URL_EXPIRATION value %q: %v. Using default value 15m.", s3PresignedURLExpirationStr, err)
		s3PresignedURLExpiration = 15 * time.Minute
	}

	// Parse ElastiCache settings
	cachePort, _ := strconv.Atoi(getEnv("CACHE_PORT", "6379"))
	cacheEnabled, _ := strconv.ParseBool(getEnv("CACHE_ENABLED", "false"))
	cacheDurationStr := getEnv("CACHE_DURATION", "1m")
	cacheDuration, err := time.ParseDuration(cacheDurationStr)
	if err != nil {
		log.Printf("Warning: invalid CACHE_DURATION value %q: %v. Using default value 5m.", cacheDurationStr, err)
		cacheDuration = 5 * time.Minute
	}

	return &Config{
		// AWS SQS settings
		AWS_SQS_VideoUpdatedURL:             getEnv("AWS_SQS_VIDEO_UPDATED_URL", ""),
		AWS_SQS_VideoUpdatedMaxMessages:     AWS_SQS_VideoUpdatedMaxMessages,
		AWS_SQS_VideoUpdatedWaitTimeSeconds: AWS_SQS_VideoUpdatedWaitTimeSeconds,

		// AWS S3 settings
		AWSS3BucketName:             getEnv("AWS_S3_BUCKET_NAME", "hackathon-video-service"),
		AWSS3Region:                 getEnv("AWS_REGION", "us-east-1"),
		AWSS3PresignedURLExpiration: s3PresignedURLExpiration,
		AWSS3BucketRawFolder:        getEnv("AWS_S3_BUCKET_RAW_FOLDER", "raw"),
		AWSS3BucketProcessedFolder:  getEnv("AWS_S3_BUCKET_PROCESSED_FOLDER", "processed"),

		// Cache settings
		CacheEndpoint: getEnv("CACHE_ENDPOINT", ""),
		CachePort:     cachePort,
		CacheEnabled:  cacheEnabled,
		CacheDuration: cacheDuration,

		// Database settings
		DBEngine: getEnv("DB_ENGINE", "postgresql"),

		DBDSN:          getEnv("DB_DSN", "host=localhost port=5432 user=postgres password=postgres dbname=fiapx sslmode=disable"),
		DBMaxOpenConns: dbMaxOpenConns,
		DBMaxIdleConns: dbMaxIdleConns,
		DBMaxLifetime:  dbMaxLifetime,

		// DocumentDB settings
		DocumentDBURI:         getEnv("DOCUMENTDB_URI", "mongodb://localhost:27017"),
		DocumentDBName:        getEnv("DOCUMENTDB_NAME", "video_service"),
		DocumentDBTLSEnabled:  documentDBTLSEnabled,
		DocumentDBTLSCertPath: getEnv("DOCUMENTDB_TLS_CERT_PATH", ""),
		DocumentDBTLSInsecure: documentDBTLSInsecure,

		// Server settings
		ServerPort:                    getEnv("SERVER_PORT", "8080"),
		ServerReadTimeout:             serverReadTimeout,
		ServerWriteTimeout:            serverWriteTimeout,
		ServerIdleTimeout:             serverIdleTimeout,
		ServerGracefulShutdownTimeout: serverGracefulShutdownTimeout,

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),

		// JWT Settings
		JWTSecret:     getEnv("JWT_SECRET", "SUPER_SECRET_KEY_DONT_TELL_ANYONE"),
		JWTExpiration: jwtExpiration,

		Auth0Domain:   getEnv("AUTH0_DOMAIN", "atomaz.us.auth0.com"),
		Auth0Audience: getEnv("AUTH0_AUDIENCE", "https://video-service.fiapx.com.br"),
		Auth0JWKSURL:  getEnv("AUTH0_JWKS_URL", "https://atomaz.us.auth0.com/.well-known/jwks.json"),

		MetricsPort: getEnv("METRICS_PORT", "8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
