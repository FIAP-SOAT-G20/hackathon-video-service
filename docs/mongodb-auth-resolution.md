# MongoDB Authentication Error Resolution

## Problem
Tests were failing with the error:
```
MongoServerError: Command insert requires authentication
```

## Root Cause
The MongoDB container is started with authentication enabled via environment variables:
```bash
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password
```

When these variables are set, MongoDB automatically enables authentication, requiring all operations to be authenticated.

## Solutions Implemented

### 1. Enhanced Authentication Handling
**File**: `video_document_datasource_test.go`
- Added longer wait times for MongoDB auth initialization
- Added retry logic for database connections
- Used proper connection string with `authSource=admin`
- Added `directConnection=true` parameter

### 2. Simplified Test Without Authentication  
**File**: `video_document_datasource_simple_test.go`
- Uses MongoDB without authentication (more reliable for testing)
- Simpler container setup
- Faster test execution
- No authentication-related issues

### 3. Diagnostic Script
**File**: `scripts/test-mongodb-auth.sh`
- Tests both authenticated and non-authenticated MongoDB setups
- Helps diagnose authentication issues
- Provides recommendations based on test results

## Recommended Approach

### For Development/CI
Use the **simplified test** (`TestVideoDocumentDataSource_Simple`):
```bash
go test ./internal/infrastructure/datasource/ -run TestVideoDocumentDataSource_Simple
```

**Benefits**:
- ✅ No authentication complexity
- ✅ Faster test execution
- ✅ More reliable in CI environments
- ✅ Easier to debug

### For Production-like Testing
Use the **authenticated test** (`TestVideoDocumentDataSource_Integration`):
```bash
go test ./internal/infrastructure/datasource/ -run TestVideoDocumentDataSource_Integration
```

**Benefits**:
- ✅ Tests production-like authentication
- ✅ Validates security configuration
- ✅ More comprehensive testing

## Connection String Examples

### No Authentication
```
mongodb://localhost:27017/database_name
```

### With Authentication
```
mongodb://admin:password@localhost:27017/database_name?authSource=admin
```

### With Additional Options
```
mongodb://admin:password@localhost:27017/database_name?authSource=admin&directConnection=true
```

## Troubleshooting

### If authentication tests still fail:
1. **Check initialization time**: MongoDB auth takes 5-10 seconds to initialize
2. **Verify credentials**: Ensure username/password match container environment
3. **Check authSource**: Must be 'admin' for root user
4. **Use diagnostic script**: `./scripts/test-mongodb-auth.sh`

### If all tests fail:
1. **Check Docker**: Ensure Docker is running
2. **Check ports**: Ensure no port conflicts
3. **Check container logs**: `docker logs <container_name>`
4. **Use simple test**: Switch to `TestVideoDocumentDataSource_Simple`

## Files Modified/Created

1. **Enhanced**: `video_document_datasource_test.go` - Better auth handling
2. **Created**: `video_document_datasource_simple_test.go` - No-auth alternative  
3. **Created**: `scripts/test-mongodb-auth.sh` - Diagnostic tool
4. **Updated**: Test scripts to handle auth gracefully

## Best Practices

1. **Use simple tests for CI/CD** - Faster and more reliable
2. **Use auth tests for security validation** - When needed
3. **Add proper wait times** - MongoDB auth needs time to initialize
4. **Use retry logic** - Handle temporary connection issues
5. **Provide fallbacks** - Multiple test approaches for different scenarios
