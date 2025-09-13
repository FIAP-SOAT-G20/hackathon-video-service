#!/bin/bash

# Test script for the new download endpoint
# This demonstrates how to use the /api/v1/videos/:id/processed endpoint

echo "Testing the new download endpoint..."

# Base URL (adjust as needed)
BASE_URL="http://localhost:8080/api/v1"

# Example video ID (replace with actual video ID)
VIDEO_ID=1

# Test the download endpoint
echo "Testing GET $BASE_URL/videos/$VIDEO_ID/processed"
curl -X GET "$BASE_URL/videos/$VIDEO_ID/processed" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n"

echo -e "\n\nExpected responses:"
echo "- 200: Success with download_url in response"
echo "- 404: Video not found"
echo "- 422: Video is not processed yet (status != FINISHED)"
echo "- 500: Internal server error"
