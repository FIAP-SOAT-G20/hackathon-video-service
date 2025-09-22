package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// migrationsFS is a filesystem that embeds the migrations folder
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

type Database struct {
	*gorm.DB
	dbDSN string
}

type GormLogger struct {
	*logger.Logger
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.InfoContext(ctx, msg, "data", data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.WarnContext(ctx, msg, "data", data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.ErrorContext(ctx, msg, "data", data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.ErrorContext(ctx, "database query failed",
			"error", err,
			"elapsed", elapsed,
			"sql", sql,
			"rows", rows,
		)
		return
	}

	l.DebugContext(ctx, "database query",
		"elapsed", elapsed,
		"sql", sql,
		"rows", rows,
	)
}

func NewPostgresConnection(cfg *config.Config, logger *logger.Logger) (*Database, error) {
	// Configure GORM with slog logger
	gormConfig := &gorm.Config{
		Logger:      &GormLogger{Logger: logger},
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), gormConfig)
	if err != nil {
		fmt.Println(cfg.DBDSN, err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	return &Database{db, cfg.DBDSN}, nil
}

// convertDSNToURL converts a GORM DSN format to a URL format for golang-migrate
func convertDSNToURL(dsn string) string {
	// If the DSN already has a scheme, return it as is
	if len(dsn) > 8 && (dsn[:8] == "postgres:" || dsn[:8] == "postgresql:") {
		return dsn
	}

	// Parse the DSN format: host=localhost port=5432 user=postgres password=postgres dbname=video_service sslmode=disable
	// Convert to: postgres://user:password@host:port/dbname?sslmode=disable

	params := make(map[string]string)
	pairs := strings.Fields(dsn)

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}

	// Build the URL
	user := params["user"]
	password := params["password"]
	host := params["host"]
	port := params["port"]
	dbname := params["dbname"]
	sslmode := params["sslmode"]

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	if sslmode != "" {
		url += fmt.Sprintf("?sslmode=%s", sslmode)
	}

	return url
}

// Migrate runs database migrations
func (db *Database) Migrate() error {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	// Convert GORM DSN format to URL format for golang-migrate
	migrationURL := convertDSNToURL(db.dbDSN)

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, migrationURL)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
