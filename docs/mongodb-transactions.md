# MongoDB Transactions Support

## Issue
The MongoDB transaction tests were failing with the error:
```
Transaction numbers are only allowed on a replica set member or mongos
```

## Root Cause
MongoDB transactions require either:
1. **Replica Set**: A cluster of MongoDB instances for high availability
2. **Sharded Cluster**: A horizontally scaled MongoDB deployment

The default `compose.yml` uses a **standalone MongoDB instance** which does not support transactions.

## Solutions Implemented

### 1. Graceful Test Handling
Updated `video_document_datasource_test.go` to:
- Detect when transactions are not supported
- Skip transaction tests with a clear message
- Prevent nil pointer dereference panics
- Continue with other tests

```go
// Check if this is a standalone MongoDB (transactions not supported)
if err != nil && (strings.Contains(err.Error(), "replica set member") || 
    strings.Contains(err.Error(), "Transaction numbers are only allowed")) {
    t.Skip("Skipping transaction test: MongoDB transactions require a replica set or sharded cluster")
    return
}
```

### 2. Enhanced Test Script
Updated `test-documentdb.sh` to:
- Handle skipped tests gracefully
- Provide informative output about transaction limitations
- Show test results summary

### 3. Replica Set Support (Optional)
Created `compose-replica-set.yml` and `test-documentdb-transactions.sh` for full transaction testing:
- Uses MongoDB replica set configuration
- Supports all transaction operations
- Separate test script for comprehensive testing

## Usage

### Standard Testing (Standalone MongoDB)
```bash
# Uses standalone MongoDB - transactions will be skipped
./scripts/test-documentdb.sh
```

### Full Transaction Testing (Replica Set)
```bash
# Uses MongoDB replica set - supports all features including transactions
./scripts/test-documentdb-transactions.sh
```

## Files Modified/Created

1. **Modified**:
   - `internal/infrastructure/datasource/video_document_datasource_test.go` - Added transaction error handling
   - `scripts/test-documentdb.sh` - Enhanced test output and error handling

2. **Created**:
   - `compose-replica-set.yml` - MongoDB replica set configuration
   - `scripts/setup-replica-set.js` - Replica set initialization script
   - `scripts/test-documentdb-transactions.sh` - Full transaction testing script

## Benefits

1. **Robust Testing**: Tests no longer fail due to transaction limitations
2. **Clear Feedback**: Users understand why transaction tests are skipped
3. **Flexibility**: Choice between simple setup and full feature testing
4. **Production Ready**: Replica set configuration mirrors production setups

## Recommendations

- Use **standalone MongoDB** (`compose.yml`) for development and basic testing
- Use **replica set** (`compose-replica-set.yml`) for production-like testing and transaction validation
- CI/CD pipelines should use replica set configuration for comprehensive testing
