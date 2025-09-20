#!/bin/bash

# DocumentDB Integration Test Script with Transaction Support
# This script uses a MongoDB replica set to support transactions

set -e

echo "🚀 DocumentDB Integration Test Script (with Transactions)"
echo "========================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

print_status "Starting DocumentDB integration test with replica set..."

# Stop any existing containers
print_status "Cleaning up existing containers..."
docker-compose -f compose-replica-set.yml down >/dev/null 2>&1 || true

# Start MongoDB replica set for testing
print_status "Starting MongoDB replica set..."
docker-compose -f compose-replica-set.yml up -d

# Wait for MongoDB replica set to be ready
print_status "Waiting for MongoDB replica set to be ready..."
for i in {1..60}; do
    if docker-compose -f compose-replica-set.yml exec -T documentdb mongosh --quiet -u admin -p password --authenticationDatabase admin --eval "rs.status().ok" >/dev/null 2>&1; then
        print_success "MongoDB replica set is ready!"
        break
    fi
    echo -n "."
    sleep 2
    if [ $i -eq 60 ]; then
        print_error "MongoDB replica set failed to start within 120 seconds"
        docker-compose -f compose-replica-set.yml logs documentdb
        exit 1
    fi
done

# Build the application
print_status "Building the application..."
go build -o bin/server ./cmd/server

# Test with MongoDB configuration (replica set supports transactions)
print_status "Testing with MongoDB replica set configuration..."
export DOCUMENTDB_URI="mongodb://admin:password@localhost:27017/video_service?authSource=admin&replicaSet=rs0"
export DOCUMENTDB_NAME="video_service"
export ENVIRONMENT="development"
export DB_ENGINE="documentdb"
export SERVER_PORT="8083"
export JWT_SECRET="test-secret-key"

# Start the server in background
print_status "Starting video service with DocumentDB replica set..."
./bin/server &
SERVER_PID=$!

# Function to cleanup on exit
cleanup() {
    print_status "Cleaning up..."
    kill $SERVER_PID 2>/dev/null || true
    docker-compose -f compose-replica-set.yml down >/dev/null 2>&1 || true
}
trap cleanup EXIT

# Wait for server to start
print_status "Waiting for server to start..."
for i in {1..15}; do
    if curl -s http://localhost:8083/api/v1/health >/dev/null 2>&1; then
        print_success "Server is running!"
        break
    fi
    echo -n "."
    sleep 2
    if [ $i -eq 15 ]; then
        print_error "Server failed to start within 30 seconds"
        exit 1
    fi
done

# Test health check
print_status "Testing health check..."
HEALTH_RESPONSE=$(curl -s http://localhost:8083/api/v1/health)
if echo "$HEALTH_RESPONSE" | grep -q '"status":"pass"' && echo "$HEALTH_RESPONSE" | grep -q '"mongodb:status"' && echo "$HEALTH_RESPONSE" | grep -q '"componentId":"db:mongodb"'; then
    print_success "Health check passed - MongoDB connection verified"
    print_status "Health response: $HEALTH_RESPONSE"
else
    print_error "Health check failed: $HEALTH_RESPONSE"
    exit 1
fi

# Run unit tests (should now support transactions)
print_status "Running unit tests with transaction support..."
TEST_OUTPUT=$(go test ./internal/infrastructure/datasource/ -v 2>&1)
TEST_EXIT_CODE=$?

echo "$TEST_OUTPUT"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_success "All unit tests passed (including transactions)"
else
    print_error "Some unit tests failed"
    exit 1
fi

print_success "DocumentDB integration test with transactions completed successfully!"

echo ""
echo "📋 Summary:"
echo "✅ MongoDB replica set started"
echo "✅ Application built successfully"
echo "✅ Server started with DocumentDB replica set"
echo "✅ Health check endpoint working"
echo "✅ All unit tests passed (including transactions)"
echo ""
echo "🔧 Manual Testing:"
echo "   Server URL: http://localhost:8083"
echo "   Health Check: curl http://localhost:8083/api/v1/health"
echo "   MongoDB Shell: docker-compose -f compose-replica-set.yml exec documentdb mongosh -u admin -p password --authenticationDatabase admin"
