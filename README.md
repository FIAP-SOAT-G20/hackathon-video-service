<a name="readme-top"></a>

![GitHub CI - Tests](https://github.com/FIAP-SOAT-G20/hackathon-video-service/actions/workflows/ci-unit-tests.yml/badge.svg)
![GitHub CI - Security](https://github.com/FIAP-SOAT-G20/hackathon-video-service/actions/workflows/ci-govulncheck.yml/badge.svg)
![GitHub CI - Lint](https://github.com/FIAP-SOAT-G20/hackathon-video-service/actions/workflows/ci-golangci-lint.yml/badge.svg)
![GitHub CI - Coverage](https://github.com/FIAP-SOAT-G20/hackathon-video-service/actions/workflows/ci-go-test-coverage.yaml/badge.svg)
![GitHub CI - BDD](https://github.com/FIAP-SOAT-G20/hackathon-video-service/actions/workflows/ci-bdd-tests.yml/badge.svg)


# <p align="center"><b>Hackathon FIAP-X</b> <small>Video Service - G20</small></p>

<p align="center">
    <img src="https://img.shields.io/badge/Code-Go-informational?style=flat-square&logo=go&color=00ADD8" alt="Go" />
    <img src="https://img.shields.io/badge/Tools-Gin-informational?style=flat-square&logo=go&color=00ADD8" alt="Gin" />
    <img src="https://img.shields.io/badge/Tools-PostgreSQL-informational?style=flat-square&logo=postgresql&color=4169E1" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/Tools-Redis-informational?style=flat-square&logo=redis&color=DC382D" alt="Redis" />
    <img src="https://img.shields.io/badge/Tools-Docker-informational?style=flat-square&logo=docker&color=2496ED" alt="Docker" />
    <img src="https://img.shields.io/badge/Tools-Swagger-informational?style=flat-square&logo=swagger&color=85EA2D" alt="Swagger" />
    <br>
    <img src="https://img.shields.io/badge/Tools-Make-informational?style=flat-square&logo=make&color=6D00CC" alt="Make" />
    <img src="https://img.shields.io/badge/Tools-GitHub_Actions-informational?style=flat-square&logo=githubactions&color=222222" alt="GitHub Actions" />
    <img src="https://img.shields.io/badge/Tools-Postman-informational?style=flat-square&logo=postman&color=FF6C37" alt="Postman" />
</p>

## 💬 About

A scalable, cloud-native video service application built with Go, featuring clean architecture and asynchronous message processing.

> [Hackathon FIAP-X Video Service Specification](docs/specification.pdf)

## 🔗 Related Projects

This project is part of a larger system that includes:

- [Video Processor Job](https://github.com/FIAP-SOAT-G20/hackathon-video-processor-job)
- [Job Starter Lambda](https://github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda)
- [Notification Lambda](https://github.com/FIAP-SOAT-G20/hackathon-notification-lambda)
- [User Lambda](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda)
- [Infrastructure (Terraform)](https://github.com/FIAP-SOAT-G20/hackaton-infrastructure-tf)
- [Infrastructure (Kubernetes)](https://github.com/FIAP-SOAT-G20/hackaton-infrastructure-deploy)

## 📚 Dictionary - Ubiquitous Language
- **Video**: An entity representing a video file uploaded by a user, including metadata such as name, description, status, and storage information.
- **User**: An entity representing a user of the video service, including metadata such as username, email, and authentication information.
- **Admin**: An entity representing an administrator user with elevated privileges for managing the video service.
- **Logged User**: An entity representing a user who uploads and manages their own video content.
- **Video Status**: The current state of a video, which can be one of the following: CREATED, PROCESSING, FINISHED, FAILED.

## ✨ Features

- **🏗️ Clean Architecture**: Implements Clean Architecture principles with clear separation of concerns across domain, use case, and infrastructure layers
- **📊 PostgreSQL Database**: High-performance relational database with ACID compliance, migrations, and connection pooling
- **🌐 RESTful API**: Complete CRUD operations for video management with OpenAPI/Swagger documentation and standardized HTTP responses
- **☁️ AWS S3 Integration**: Secure video storage with automatic presigned URL generation for uploads and downloads with configurable expiration times
- **⚡ Asynchronous Processing**: AWS SQS integration with dedicated worker consumer for scalable, non-blocking video processing workflows
- **🔐 JWT Authentication**: Stateless, secure token-based authentication with configurable expiration and refresh capabilities
- **⚡ Redis Caching**: High-performance in-memory caching with Redis and AWS ElastiCache support for improved response times and reduced database load
- **🏥 Health Monitoring**: Comprehensive health checks with database connectivity validation and service status reporting
- **🐳 Containerized Deployment**: Multi-stage Docker builds optimized for production with Docker Compose for local development
- **🧪 Comprehensive Testing**: Unit tests (80%+ coverage), integration tests, BDD tests with Gherkin scenarios, and automated CI/CD pipelines
- **📝 Interactive Documentation**: Auto-generated OpenAPI/Swagger UI with live API testing and comprehensive endpoint documentation
- **🔧 Enhanced Developer Experience**: Hot reload with Air, extensive Makefile automation, code generation, linting, and Git workflow tools
- **📈 Production Ready**: Kubernetes deployment manifests, monitoring hooks, graceful shutdowns, and horizontal scaling support

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🚀 Getting Started

### Prerequisites

- **Go 1.25+**
- **Docker and Docker Compose**
- **Redis** (for caching - can run via Docker)
- **Make** (for development commands)
- **AWS Account** (for S3, SQS, and ElastiCache services)

### 🚀 Development Setup

```bash
# Clone and setup
git clone <repository-url>
cd hackathon-video-service
cp .env.example .env  # Configure your environment variables

# Start PostgreSQL database and Redis cache
make run-db
make run-cache

# Build and run the API server (includes database and cache)
make build
make run-api

# In another terminal, run the worker (if SQS is configured)
make run-worker
```

### 🌐 Access Points

- **API Server**: [http://localhost:8080](http://localhost:8080)
- **API Documentation**: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)
- **Health Check**: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
- **PgAdmin**: [http://localhost:5050](http://localhost:5050) (PostgreSQL admin interface)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

##  🏗️ Architecture

The application follows **Clean Architecture** principles with a microservices approach, featuring:

### System Components

1. **API Server** (`cmd/server/main.go`) - RESTful API for video management
2. **Worker Consumer** (`cmd/worker/consumer/main.go`) - SQS message processor for asynchronous video updates
3. **Database Layer** - PostgreSQL with GORM ORM and migrations
4. **Caching Layer** - Redis/AWS ElastiCache for high-performance data caching and reduced database load
5. **Message Queue** - AWS SQS for decoupled, scalable message processing
6. **Object Storage** - AWS S3 for video file storage with presigned URLs

### System Design

![System Design Diagram](docs/system-design.png)

### Project Structure

```
cmd/
├── server/                 # API application entry point
└── worker/                 # Worker applications
    └── consumer/           # SQS message consumer
internal/
├── core/                   # Business logic (Clean Architecture)
│   ├── domain/            # Entities, value objects, and domain errors
│   ├── dto/               # Data transfer objects
│   ├── port/              # Interface definitions (dependency inversion)
│   └── usecase/           # Business use cases
├── adapter/               # Interface adapters
│   ├── controller/        # HTTP controllers and request handling
│   ├── gateway/           # Repository implementations
│   └── presenter/         # Response formatting (JSON, XML, pagination)
└── infrastructure/        # Infrastructure layer
    ├── config/            # Application configuration
    ├── database/          # Database connections and migrations
    ├── datasource/        # Data access implementations
    ├── handler/           # HTTP route handlers
    ├── middleware/        # HTTP middleware (auth, logging, etc.)
    ├── pkg/               # Infrastructure packages
    │   └── aws/           # AWS integrations (S3, SQS)
    └── server/            # HTTP server setup
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## ⚙️ Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and configure according to your environment.

### Database Configuration

```bash
# PostgreSQL Configuration
DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=fiapx sslmode=disable"
DB_ENGINE="postgres"
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=300s
```

### AWS Services Configuration

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
AWS_SESSION_TOKEN=your_session_token  # Optional for temporary credentials

# S3 Configuration
AWS_S3_BUCKET_NAME=your-video-bucket
AWS_S3_BUCKET_RAW_FOLDER=raw           # Folder for uploaded videos
AWS_S3_BUCKET_PROCESSED_FOLDER=processed  # Folder for processed videos
AWS_S3_PRESIGNED_URL_EXPIRATION=15m    # URL expiration time

# SQS Configuration (for worker)
AWS_SQS_VIDEO_UPDATED_URL=https://sqs.us-east-1.amazonaws.com/123456789012/video-updated
AWS_SQS_VIDEO_UPDATED_MAX_MESSAGES=10
AWS_SQS_VIDEO_UPDATED_WAIT_TIME_SECONDS=20
```

### Application Configuration

```bash
# Environment
ENVIRONMENT=development  # development, staging, production

# Server Settings
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# JWT Settings
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION=24h
```

### Cache Configuration

```bash
# Redis/ElastiCache Configuration
CACHE_ENABLED=true                    # Enable/disable caching
CACHE_ENDPOINT=localhost              # Redis endpoint (local or AWS ElastiCache)
CACHE_PORT=6379                       # Redis port

# For AWS ElastiCache (production)
# CACHE_ENDPOINT=your-elasticache-endpoint.cache.amazonaws.com
```

> **Local Development**: Use Docker to run Redis locally:
> ```bash
> docker run -d --name redis -p 6379:6379 redis:alpine
> ```

> **Production**: Configure AWS ElastiCache endpoint for high availability and managed Redis instances.

### Worker Message Format

The SQS worker processes messages in the following format:

```json
{
  "video_id": 123,
  "status": "FINISHED",
  "hash": "abc123def456"  // Optional: base64 hash of processed video
}
```

> **Valid Status Values**: `CREATED`, `PROCESSING`, `FINISHED`, `FAILED`

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 🌐 API Reference

### Core Video Management

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/api/v1/health` | Health check endpoint | Health status |
| `GET` | `/api/v1/videos` | List videos with pagination | Video list |
| `POST` | `/api/v1/videos` | Create video and get upload URL | Video with presigned upload URL |
| `GET` | `/api/v1/videos/{id}` | Get video by ID | Video details |
| `PUT` | `/api/v1/videos/{id}` | Update video metadata | Updated video |
| `DELETE` | `/api/v1/videos/{id}` | Delete video | Deletion confirmation |
| `GET` | `/api/v1/videos/{id}/processed` | Get download URL | Presigned download URL |

> [!TIP]
> **Cache Performance**: Video list queries are cached in Redis for improved response times. Cache duration is 5 minutes for list results.

### API Examples

#### Create Video (with Upload URL)
```bash
curl -X POST http://localhost:8080/api/v1/videos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "user_id": 1,
    "name": "My Awesome Video",
    "description": "A sample video upload"
  }'
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "My Awesome Video",
  "description": "A sample video upload",
  "status": "CREATED",
  "presigned_url": "https://your-bucket.s3.amazonaws.com/raw/video-uuid?AWSAccessKeyId=...",
  "created_at": "2024-02-09T10:00:00Z",
  "updated_at": "2024-02-09T10:00:00Z"
}
```

#### Get Video with Download URL
```bash
curl -X GET http://localhost:8080/api/v1/videos/1/processed \
  -H "Authorization: Bearer <jwt-token>"
```

**Response:**
```json
{
  "presigned_url": "https://your-bucket.s3.amazonaws.com/processed/video-hash?AWSAccessKeyId=..."
}
```

#### List Videos with Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/videos?page=1&limit=10" \
  -H "Authorization: Bearer <jwt-token>"
```

**Response:**
```json
{
  "videos": [
    {
      "id": 1,
      "user_id": 1,
      "name": "My Awesome Video",
      "description": "A sample video upload",
      "status": "CREATED",
      "created_at": "2024-02-09T10:00:00Z",
      "updated_at": "2024-02-09T10:00:00Z"
    }
    // More videos...
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_pages": 5,
    "total_items": 50
  }
}
```

> [!NOTE]
> You can filter videos by `hash` and `status` using query parameters. 

### Presigned URLs

The service provides secure, time-limited URLs for:

- **Upload URLs**: Generated when creating a video (expires in 15 minutes)
- **Download URLs**: Generated when retrieving processed videos (expires in 15 minutes)
- **Security**: URLs are signed with AWS credentials and include expiration timestamps

### Authentication

Most endpoints require JWT authentication:

```bash
# Include in request headers
Authorization: Bearer <your-jwt-token>
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🛠️ Development

### Development Commands

```bash
# 📋 Display all available commands
make help

# 🔨 Build and Run
make build                      # Build the application
make run-api                    # Run API server with PostgreSQL
make run-worker                 # Run SQS worker consumer
make run-api-air                # Run API with hot reload (Air)

# 🗄️ Database Management
make run-db                     # Start PostgreSQL with PgAdmin
make stop-db                    # Stop database services
make migrate-up                 # Run database migrations
make migrate-down               # Roll back migrations
make migrate-create name=<name> # Create new migration

# ⚡ Cache Management
make run-cache                  # Start Redis cache

# 🐳 Docker Development
make compose-up                 # Start full environment
make compose-down               # Stop all services
make compose-clean              # Clean volumes and images

# 🧪 Testing
make test                       # Run test with linting
make coverage                   # Generate test coverage report
make bdd-tests                  # Run BDD/Gherkin tests
make test-integration           # Run integration tests

# 🔧 Code Quality
make lint                       # Run linter
make fmt                        # Format code
make mock                       # Generate mocks
make swagger                    # Generate Swagger docs
make scan                       # Run security scan

# 🚀 CI/CD
make docker-build               # Build Docker image
make docker-push                # Push to registry
make new-branch                 # Create new feature branch
make pull-request               # Create pull request
```

### Testing Strategy

#### Unit Tests
```bash
# Run unit tests with linting, race detection and coverage
make test

# Generate HTML coverage report
make coverage
```

> [!NOTE]
> Unit tests cover core business logic with ~80%+ coverage, including use cases and domain logic.

#### BDD Tests
```bash
# Run behavior-driven development tests
make bdd-tests
```

### Hot Reload Development

Use Air for automatic rebuilds during development:

```bash
make run-api-air
```

Air configuration is in `air.toml` and watches for changes in `.go` files.

### Code Generation

```bash
# Generate mocks for all port interfaces
make mock

# Generate/update Swagger documentation
make swagger
```

> [!TIP]
> Apply those commands after modifying interfaces or API routes.

### Database Migrations

```bash
# Create a new migration
make migrate-create name=add_video_status_field

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

> [!NOTE]
> You don't need to run migrations manually in development,  
> they run automatically on server start with `golang-migrate` and stored in `db/migrations`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🚢 Deployment

### Docker Compose Environments

#### Development with PostgreSQL
```bash
# Start API server, worker, and PostgreSQL
docker-compose up -d

# Access services
# API: http://localhost:8080
# PgAdmin: http://localhost:5050
```

#### Production Deployment

The application provides multi-stage Docker builds optimized for production:

```dockerfile
# API Server
FROM golang:1.25-alpine AS builder
# ... build stage
FROM alpine:latest
COPY --from=builder /app/bin/app /app/app
EXPOSE 8080
CMD ["/app/app"]

# Worker Consumer
FROM golang:1.25-alpine AS builder
# ... build stage  
FROM alpine:latest
COPY --from=builder /app/bin/worker /app/worker
CMD ["/app/worker"]
```

### Container Registry

Images are published to GitHub Container Registry:

```bash
# Build and push images
make docker-build
make docker-push

# Available images
ghcr.io/fiap-soat-g20/fiapx-video-service:latest
ghcr.io/fiap-soat-g20/fiapx-video-service:<version>
```

> [!NOTE]
> This images are automatically built and pushed via GitHub Actions on new releases.

### Kubernetes Deployment

For Kubernetes deployment, see the [K8s documentation](docs/k8s.jpg) which includes:

- Deployment manifests for API server and worker
- Service definitions and ingress configuration
- ConfigMaps and Secrets for environment variables
- Horizontal Pod Autoscaler (HPA) configuration
- Health check and readiness probes

### AWS Production Setup

#### Required AWS Services

1. **RDS PostgreSQL** for data persistence
2. **S3 Bucket** for video storage with appropriate IAM policies
3. **SQS Queue** for message processing
4. **ECS/EKS** or **EC2** for container hosting
5. **Application Load Balancer** for traffic distribution

#### Environment Configuration

```bash
# Production environment variables
ENVIRONMENT=production
SERVER_PORT=8080

# Database Configuration
DB_DSN="postgres://user:pass@rds-endpoint:5432/dbname?sslmode=require"

# AWS Services
AWS_S3_BUCKET_NAME=production-video-bucket
AWS_SQS_VIDEO_UPDATED_URL=https://sqs.region.amazonaws.com/account/video-updated

# Security
JWT_SECRET=production-secret-key-256-bits
```

### Health Checks

The application provides comprehensive health endpoints:

- **Liveness Probe**: `GET /api/v1/health`
- **Readiness Probe**: `GET /api/v1/health` (includes database connectivity)

Configure in Kubernetes:
```yaml
livenessProbe:
  httpGet:
    path: /api/v1/health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /api/v1/health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Scaling Considerations

- **API Server**: Stateless, can be horizontally scaled
- **Worker Consumer**: Can run multiple instances for parallel message processing
- **Database**: Use read replicas for read-heavy workloads
- **S3**: Automatically scales, consider CloudFront for global distribution

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**: `make new-branch`
3. **Make your changes** with tests
4. **Run quality checks**: `make tests lint`
5. **Commit your changes** with conventional commits
6. **Push to your branch**
7. **Create a Pull Request**: `make pull-request`

### Code Standards

- **Go Style**: Follow effective Go practices and `gofmt` formatting
- **Clean Architecture**: Maintain separation of concerns
- **Test Coverage**: Maintain >80% coverage for core business logic
- **Documentation**: Update documentation for API changes
- **Conventional Commits**: Use conventional commit format

### Pull Request Checklist

- [ ] Tests pass (`make tests`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation updated
- [ ] Integration tests included for new features
- [ ] Environment variables documented

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 📄 Documentation

### Internal Documentation

- [📋 Agreements](docs/agreements.md) - Project agreements and standards
- [🏗️ Implementation Summary](docs/implementation-summary.md) - Technical implementation details
- [🗄️ Database Schema](docs/db-schema.dbml) - Database design and relationships
- [☸️ Kubernetes Deployment](docs/k8s.jpg) - K8s architecture diagram
- [🔄 Flow Diagram](docs/flow-diagram.png) - System flow visualization
- [🧪 Validation Testing](docs/validation-testing.md) - Testing strategies

### API Documentation

- **Interactive Swagger UI**: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)
- **ReDoc**: [http://localhost:8080/redoc](http://localhost:8080/redoc)
- **OpenAPI Spec**: [docs/swagger.yaml](docs/swagger.yaml)
- **Postman Collection**: [docs/postman_collection.json](docs/postman_collection.json)
- **HTTP Requests**: [docs/requests.http](docs/requests.http)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 📚 References

### Architecture and Design Patterns

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Video Sharing Platform Architecture Design](https://medium.com/@dominikzygarski_88070/project-video-sharing-platform-architecture-design-400e718f77b3)
- [Distributed High Performance Video Processing](https://www-di.inf.puc-rio.br/~endler/paperlinks/CLOUD-2010.pdf)
- [System Design — Video Processing System](https://medium.com/@qingedaig/system-design-video-processing-system-3742af267ba5)

### AWS Services Integration

- [AWS SDK for Go Documentation](https://docs.aws.amazon.com/sdk-for-go/)
- [S3 Presigned URLs with Go](https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57)
- [Go S3 Package Documentation](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3)
- [Download and Upload with Presigned URLs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html)
- [Working with Object Metadata](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingMetadata.html)

### Message Queue and SQS

- [AWS SQS with Golang - Publishing and Consuming](https://medium.com/@weberasantos/publicando-e-consumindo-mensagem-no-sqs-aws-com-golang-6970a0e7581e)
- [Process AWS SQS Messages with Goroutines](https://medium.com/@wiraizkandar/process-aws-sqs-message-with-goroutines-98ff4799ea69)
- [Simplifying Message Queueing with Golang and Amazon SQS](https://dev.to/thanhphuchuynh/simplifying-message-queueing-with-golang-and-amazon-sqs-3gpl)

### Development Tools and Practices

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Docker Multi-stage Builds](https://docs.docker.com/develop/dev-best-practices/)
- [GitHub Actions for Go](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go)
- [How to Connect Your GitHub Project to Sonar](https://dev.to/olsido/how-to-connect-your-github-project-to-sonar-9ic)
- [How to Enable SonarCloud for Your Project](https://dev.to/olsido/how-to-enable-sonarcloud-for-your-project-aoi)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---
<p>

<p align="center">
    Made with&nbsp;&nbsp;♥️&nbsp;&nbsp;by FIAP 10SOAT G21 Team 
</p>
