#!/bin/bash

# MongoDB Authentication Test Script
# This script helps debug MongoDB authentication issues in tests

echo "🔍 MongoDB Authentication Diagnostic"
echo "===================================="

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running"
    exit 1
fi

echo "✅ Docker is running"

# Test 1: Start MongoDB without authentication
echo ""
echo "📋 Test 1: MongoDB without authentication"
echo "----------------------------------------"

docker run --rm -d --name mongo-test-noauth -p 27018:27017 mongo:7.0 >/dev/null 2>&1
sleep 3

if docker exec mongo-test-noauth mongosh --quiet --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
    echo "✅ No-auth MongoDB connection: SUCCESS"
else
    echo "❌ No-auth MongoDB connection: FAILED"
fi

docker stop mongo-test-noauth >/dev/null 2>&1

# Test 2: Start MongoDB with authentication
echo ""
echo "📋 Test 2: MongoDB with authentication"
echo "-------------------------------------"

docker run --rm -d --name mongo-test-auth \
    -p 27019:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=admin \
    -e MONGO_INITDB_ROOT_PASSWORD=password \
    mongo:7.0 >/dev/null 2>&1

echo "⏳ Waiting for authenticated MongoDB to initialize..."
sleep 8

# Test unauthenticated connection (should fail)
if docker exec mongo-test-auth mongosh --quiet --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
    echo "⚠️  Unauthenticated connection succeeded (unexpected)"
else
    echo "✅ Unauthenticated connection failed (expected)"
fi

# Test authenticated connection
if docker exec mongo-test-auth mongosh --quiet -u admin -p password --authenticationDatabase admin --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
    echo "✅ Authenticated connection: SUCCESS"
    
    # Test database operations
    if docker exec mongo-test-auth mongosh --quiet -u admin -p password --authenticationDatabase admin test --eval "db.testcoll.insertOne({test: 1})" >/dev/null 2>&1; then
        echo "✅ Database insert operation: SUCCESS"
    else
        echo "❌ Database insert operation: FAILED"
    fi
else
    echo "❌ Authenticated connection: FAILED"
fi

docker stop mongo-test-auth >/dev/null 2>&1

echo ""
echo "📋 Recommendations:"
echo "==================="
echo "1. If Test 1 passes but Test 2 fails → Authentication initialization issue"
echo "2. If both tests fail → Docker/MongoDB setup issue"
echo "3. If Test 2 passes → Use authenticated MongoDB in tests"
echo "4. If only Test 1 passes → Use non-authenticated MongoDB for simpler testing"

echo ""
echo "🔧 Connection strings to try:"
echo "No auth:  mongodb://localhost:27018/test"
echo "With auth: mongodb://admin:password@localhost:27019/test?authSource=admin"
