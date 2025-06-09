package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanut9/file-service/internal/config"
)

// Pass the file size and the multipart.File directly for streaming
func UploadToR2(file multipart.File, fileHeader *multipart.FileHeader, bucketName, folderPath string, isPublic bool) (string, error) {
	if bucketName == "" {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		if bucketName == "" {
			return "", fmt.Errorf("R2_BUCKET_NAME is not set and no bucketName provided")
		}
	}

	head := make([]byte, 512)
	_, err := file.Read(head)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to seek file: %v", err)
	}

	mimeType := http.DetectContentType(head)
	key := filepath.Join(folderPath, fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename)))

	_, err = config.R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String(mimeType),
		ContentLength: &fileHeader.Size,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %v", err)
	}

	var fileURL string
	if isPublic {
		publicHost := os.Getenv("R2_PUBLIC_ENDPOINT")
		if publicHost == "" {
			publicHost = fmt.Sprintf("https://%s.r2.dev", bucketName)
		}
		fileURL = fmt.Sprintf("%s/%s", publicHost, key)
	} else {
		endpoint := os.Getenv("R2_ENDPOINT")
		if endpoint == "" {
			accountId := os.Getenv("R2_ACCOUNT_ID")
			endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
		}
		fileURL = fmt.Sprintf("%s/%s/%s", endpoint, bucketName, key)
	}

	return fileURL, nil
}
