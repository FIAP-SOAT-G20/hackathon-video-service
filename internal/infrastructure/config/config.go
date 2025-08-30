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
	AWS_SQS_OrderStatusUpdatedURL             string
	AWS_SQS_OrderStatusUpdatedMaxMessages     int
	AWS_SQS_OrderStatusUpdatedWaitTimeSeconds int

	// Database settings
	DBDSN          string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  time.Duration

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
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	AWS_SQS_OrderStatusUpdatedMaxMessages, _ := strconv.Atoi(getEnv("AWS_SQS_ORDER_STATUS_UPDATED_MAX_MESSAGES", "10"))
	AWS_SQS_OrderStatusUpdatedWaitTimeSeconds, _ := strconv.Atoi(getEnv("AWS_SQS_ORDER_STATUS_UPDATED_WAIT_TIME_SECONDS", "20"))

	dbMaxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	dbMaxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	dbMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

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

	return &Config{
		// AWS SQS settings
		AWS_SQS_OrderStatusUpdatedURL:             getEnv("AWS_SQS_ORDER_STATUS_UPDATED_URL", ""),
		AWS_SQS_OrderStatusUpdatedMaxMessages:     AWS_SQS_OrderStatusUpdatedMaxMessages,
		AWS_SQS_OrderStatusUpdatedWaitTimeSeconds: AWS_SQS_OrderStatusUpdatedWaitTimeSeconds,

		// Database settings
		DBDSN:          getEnv("DB_DSN", "host=localhost port=5432 user=postgres password=postgres dbname=fastfood_10soat_g19_tc4_order sslmode=disable"),
		DBMaxOpenConns: dbMaxOpenConns,
		DBMaxIdleConns: dbMaxIdleConns,
		DBMaxLifetime:  dbMaxLifetime,

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
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
