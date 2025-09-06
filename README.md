# Hackathon Video Service

A scalable video service application with support for multiple database backends including PostgreSQL and AWS DocumentDB.

## Features

- **Multi-Database Support**: PostgreSQL, AWS DocumentDB, and MongoDB
- **RESTful API**: Complete CRUD operations for video management
- **Authentication**: JWT-based authentication
- **Health Checks**: Comprehensive health monitoring
- **Docker Support**: Full containerization with Docker Compose
- **Testing**: Unit tests and integration tests
- **Documentation**: OpenAPI/Swagger documentation

## Database Support

This application supports multiple database backends:

- **PostgreSQL** (default) - Traditional relational database
- **AWS DocumentDB** - Managed MongoDB-compatible service
- **MongoDB** - NoSQL document database for local development

### Quick Start with PostgreSQL

```bash
# Start PostgreSQL database
make run-db

# Build and run the application
make build
make run-postgres
```

### Quick Start with DocumentDB/MongoDB

```bash
# Start MongoDB for local development
make documentdb-up

# Start MongoDB with Mongo Express UI
make documentdb-up-with-ui

# Build and run with DocumentDB
make build
make run-documentdb
```

**MongoDB Web UI**: When using `make documentdb-up-with-ui`, you can access the Mongo Express web interface at [http://localhost:8082](http://localhost:8082) to view and manage your MongoDB data.

## Configuration

The application automatically selects the database type based on environment variables:

```bash
# For PostgreSQL (default)
export DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=video_service sslmode=disable"

# For DocumentDB/MongoDB
export DOCUMENTDB_URI="mongodb://username:password@cluster-endpoint:27017/database"
export DOCUMENTDB_NAME="video_service"
export ENVIRONMENT="documentdb"
```

See [DocumentDB Integration Guide](docs/documentdb-integration.md) for detailed configuration.

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/videos` - List videos
- `POST /api/v1/videos` - Create video
- `GET /api/v1/videos/{id}` - Get video by ID
- `PUT /api/v1/videos/{id}` - Update video
- `DELETE /api/v1/videos/{id}` - Delete video

## Development

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make

### Available Commands

```bash
make help                       # Show all available commands
make build                      # Build the application
make test                       # Run tests
make run-postgres               # Run with PostgreSQL
make run-documentdb             # Run with DocumentDB/MongoDB
make documentdb-up              # Start MongoDB only
make documentdb-up-with-ui      # Start MongoDB with Mongo Express UI
make mongo-express-up           # Start Mongo Express UI (requires MongoDB)
make mongo-express-down         # Stop Mongo Express UI
make compose-up-with-ui         # Start full environment with UI
make documentdb-test            # Test DocumentDB integration
make test-integration           # Run integration tests for both databases
```

### Testing

```bash
# Run unit tests
make test

# Run DocumentDB integration tests
make documentdb-test

# Run all integration tests
make test-integration
```

## Deployment

### Docker Compose

```bash
# Start with PostgreSQL
docker-compose up -d

# Start with both PostgreSQL and MongoDB for testing
docker-compose -f docker-compose.documentdb.yml up -d
```

### Production

For production deployment with AWS DocumentDB, see the [deployment guide](docs/documentdb-integration.md#deployment).

## Architecture

The application follows Clean Architecture principles with support for multiple database backends through a factory pattern:

```
cmd/server/                 # Application entry point
internal/
├── core/                   # Business logic
│   ├── domain/            # Entities and value objects
│   ├── dto/               # Data transfer objects
│   └── port/              # Interface definitions
├── adapter/               # Interface adapters
│   ├── controller/        # HTTP controllers
│   ├── gateway/           # Data access gateways
│   └── presenter/         # Response formatters
└── infrastructure/        # Infrastructure layer
    ├── database/          # Database connections (PostgreSQL, DocumentDB)
    ├── datasource/        # Data source implementations
    └── handler/           # HTTP handlers
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
