# DocumentDB Integration

This document describes the implementation and usage of AWS DocumentDB integration in the video service.

## Overview

The video service now supports multiple database backends:
- **PostgreSQL** (default)
- **AWS DocumentDB** (MongoDB-compatible)
- **MongoDB** (for local development)

The application automatically determines which database to use based on configuration.

## Database Selection Logic

The application selects the database type based on the following priority:

1. **DocumentDB URI**: If `DOCUMENTDB_URI` contains "docdb" or "documentdb", it uses DocumentDB
2. **MongoDB URI**: If `DOCUMENTDB_URI` is set but doesn't contain DocumentDB keywords, it uses MongoDB
3. **Default**: PostgreSQL

## Configuration

### Environment Variables

```bash
# DocumentDB Configuration
DOCUMENTDB_URI=mongodb://username:password@cluster-endpoint:27017/database?tls=true&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false
DOCUMENTDB_NAME=video_service
DOCUMENTDB_TLS_ENABLED=true
DOCUMENTDB_TLS_CERT_PATH=/path/to/rds-combined-ca-bundle.pem
DOCUMENTDB_TLS_INSECURE=false

# Set environment to explicitly use DocumentDB
ENVIRONMENT=documentdb
```

### AWS DocumentDB Setup

1. **Create DocumentDB Cluster**:
   ```bash
   aws docdb create-db-cluster \
     --db-cluster-identifier video-service-cluster \
     --engine docdb \
     --master-username admin \
     --master-user-password your-password \
     --vpc-security-group-ids sg-xxxxxxxx \
     --db-subnet-group-name your-subnet-group
   ```

2. **Create DocumentDB Instance**:
   ```bash
   aws docdb create-db-instance \
     --db-instance-identifier video-service-instance \
     --db-instance-class db.t3.medium \
     --engine docdb \
     --db-cluster-identifier video-service-cluster
   ```

3. **Download TLS Certificate**:
   ```bash
   wget https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem
   ```

### Local MongoDB Setup

For local development, you can use Docker:

```bash
# Start MongoDB with Docker
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password \
  mongo:7.0

# Set environment variables
export DOCUMENTDB_URI="mongodb://admin:password@localhost:27017/video_service?authSource=admin"
export DOCUMENTDB_NAME="video_service"
export ENVIRONMENT="mongodb"
```

## Data Model

### Video Document Structure

```json
{
  "_id": ObjectId("..."),
  "user_id": 456,
  "status": "OPEN",
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

### Indexes

The following indexes are automatically created:

- `user_id`
- `status`
- `status + user_id` (compound)
- `created_at` (descending)
- `updated_at` (descending)

## Features

### Supported Operations

- ✅ Create video
- ✅ Find video by ID
- ✅ Find all videos with filtering
- ✅ Update video
- ✅ Delete video
- ✅ Transactions
- ✅ Pagination
- ✅ Sorting

### Filtering

The DocumentDB implementation supports the same filtering options as PostgreSQL:

```go
filters := map[string]any{
    "user_id": uint64(123),
    "statuses": []valueobject.VideoStatus{valueobject.OPEN, valueobject.PENDING},
    "statuses_exclude": []valueobject.VideoStatus{valueobject.CANCELLED},
}
```

### Sorting

Supported sort options:
- `id asc` / `id desc`
- `created_at asc` / `created_at desc`
- `updated_at asc` / `updated_at desc`

## Deployment

### Docker Compose

Use the provided Docker Compose file for testing:

```bash
# Start both PostgreSQL and MongoDB services
docker-compose -f docker-compose.documentdb.yml up -d

# Test PostgreSQL version (port 8080)
curl http://localhost:8080/api/v1/health

# Test MongoDB version (port 8081)
curl http://localhost:8081/api/v1/health
```

### Kubernetes

For production deployment, use the existing Kubernetes manifests and add DocumentDB configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: video-service-config
data:
  ENVIRONMENT: "documentdb"
  DOCUMENTDB_URI: "mongodb://username:password@cluster-endpoint:27017/video_service?tls=true&replicaSet=rs0"
  DOCUMENTDB_NAME: "video_service"
  DOCUMENTDB_TLS_ENABLED: "true"
---
apiVersion: v1
kind: Secret
metadata:
  name: video-service-secrets
type: Opaque
stringData:
  DOCUMENTDB_URI: "mongodb://username:password@cluster-endpoint:27017/video_service?tls=true&replicaSet=rs0"
```

## Testing

### Unit Tests

Run the DocumentDB integration tests:

```bash
# Run all tests
go test ./internal/infrastructure/datasource/...

# Run only DocumentDB tests
go test -run TestVideoDocumentDataSource ./internal/infrastructure/datasource/
```

### Performance Testing

The DocumentDB implementation includes the same performance optimizations as PostgreSQL:

- Connection pooling
- Index optimization
- Efficient pagination
- Bulk operations support

## Migration

### From PostgreSQL to DocumentDB

1. **Export data from PostgreSQL**:
   ```bash
   pg_dump -h localhost -U postgres -d video_service --data-only --column-inserts > data.sql
   ```

2. **Convert to MongoDB format** (manual process or use migration scripts)

3. **Import to DocumentDB**:
   ```bash
   mongoimport --uri="mongodb://..." --collection=videos --file=videos.json
   ```

### Schema Differences

| PostgreSQL | DocumentDB/MongoDB |
|------------|-------------------|
| `id` (serial) | `video_id` (uint64) + `_id` (ObjectID) |
| Table rows | Documents |
| Foreign keys | Embedded documents |
| ACID transactions | ACID transactions (4.0+) |

## Troubleshooting

### Common Issues

1. **TLS Connection Errors**:
   - Ensure `DOCUMENTDB_TLS_ENABLED=true`
   - Download and specify the correct CA bundle
   - Check security group settings

2. **Authentication Errors**:
   - Verify username/password in the URI
   - Check IAM roles for DocumentDB access
   - Ensure proper `authSource` in connection string

3. **Performance Issues**:
   - Monitor connection pool settings
   - Check index usage with `explain()`
   - Adjust read preferences for replicas

### Monitoring

Use the application logs to monitor database operations:

```bash
# Enable debug logging
export LOG_LEVEL=debug

# View database operations
docker logs video-service | grep "database"
```

## Security Considerations

1. **Encryption in Transit**: Always use TLS for DocumentDB connections
2. **Encryption at Rest**: Enable encryption for DocumentDB clusters
3. **Access Control**: Use IAM roles and security groups
4. **Connection String Security**: Store sensitive data in secrets management
5. **Network Security**: Use VPC and private subnets

## Performance Tuning

### Connection Pool Settings

```bash
# Adjust based on your workload
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=5m
```

### DocumentDB Instance Types

- **Development**: db.t3.medium
- **Production**: db.r5.large or higher
- **High Performance**: db.r5.xlarge with provisioned IOPS

### Read Replicas

For read-heavy workloads, configure read replicas:

```bash
DOCUMENTDB_URI="mongodb://username:password@cluster-endpoint:27017/database?readPreference=secondaryPreferred"
```
