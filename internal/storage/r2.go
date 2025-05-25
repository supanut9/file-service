package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanut9/file-service/internal/config"
)

func UploadToR2(file multipart.File, filename, bucketName, folderPath string, isPublic bool) (string, error) {
	if bucketName == "" {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		if bucketName == "" {
			return "", fmt.Errorf("R2_BUCKET_NAME is not set and no bucketName provided")
		}
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// inside UploadToR2(...)
	mimeType := http.DetectContentType(buf.Bytes())

	key := filepath.Join(folderPath, fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(filename)))

	_, err := config.R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(mimeType), //
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %v", err)
	}

	var fileURL string
	if isPublic {
		publicHost := os.Getenv("R2_PUBLIC_ENDPOINT") // e.g. https://my-bucket.r2.dev
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
