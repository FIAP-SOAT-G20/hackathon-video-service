# Download Processed Video Endpoint

## Overview
A new endpoint has been added to generate presigned URLs for downloading processed videos from S3.

## Endpoint Details

### GET `/api/v1/videos/:id/processed`

Generates a presigned URL to download a processed video from S3.

#### Request Parameters
- `id` (path parameter, required): The video ID

#### Response Format
```json
{
  "download_url": "https://s3.amazonaws.com/bucket/processed/path/to/video?X-Amz-Algorithm=..."
}
```

#### Status Codes
- `200`: Success - Returns presigned download URL
- `400`: Bad Request - Invalid video ID
- `404`: Not Found - Video doesn't exist
- `422`: Unprocessable Entity - Video is not processed yet (status != FINISHED)
- `500`: Internal Server Error

#### Business Rules
1. Only videos with status `FINISHED` can be downloaded
2. The presigned URL is valid for 1 hour
3. The download URL points to the processed video in the S3 bucket under the `processed/` folder

#### Example Usage

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/videos/123/processed" \
  -H "Content-Type: application/json"
```

**Success Response (200):**
```json
{
  "download_url": "https://s3.amazonaws.com/my-bucket/processed/123/encoded_video_name?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=..."
}
```

**Error Response (422 - Video not processed):**
```json
{
  "error": "video is not processed yet"
}
```

## Implementation Details

### New Components Added:

1. **DTO**: `DownloadVideoInput` in `internal/core/dto/video_dto.go`
2. **Port Methods**: 
   - `VideoController.Download()`
   - `VideoUseCase.Download()` 
   - `S3Service.GeneratePresignedDownloadURL()`
3. **Request Struct**: `DownloadVideoUriRequest` in `internal/infrastructure/handler/request/video_request.go`
4. **Handler Method**: `VideoHandler.Download()` in `internal/infrastructure/handler/video_handler.go`
5. **Use Case Logic**: Download validation and URL generation in `internal/core/usecase/video_usecase.go`
6. **S3 Service Method**: Presigned download URL generation in `internal/infrastructure/service/s3_service.go`

### Key Features:

- **Security**: Only allows downloads for videos with `FINISHED` status
- **Temporary Access**: Presigned URLs expire after 1 hour
- **Error Handling**: Proper error responses for different scenarios
- **Documentation**: Full Swagger/OpenAPI documentation
- **Testing**: Mocks updated and test script provided

### S3 Folder Structure:
- Raw uploads: `raw/{user_id}/{video_hash}`
- Processed videos: `processed/{user_id}/{video_hash}`
