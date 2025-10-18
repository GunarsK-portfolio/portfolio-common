package utils

import "fmt"

// BuildFileURL constructs the full URL for a file stored in MinIO/S3
// Format: {filesAPIURL}/files/{fileType}/{s3Key}
// Example: http://localhost:8085/api/v1/files/image/avatars/user123.jpg
func BuildFileURL(filesAPIURL, fileType, s3Key string) string {
	if filesAPIURL == "" || s3Key == "" {
		return ""
	}
	return fmt.Sprintf("%s/files/%s/%s", filesAPIURL, fileType, s3Key)
}
