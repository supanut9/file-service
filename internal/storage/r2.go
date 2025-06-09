package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanut9/file-service/internal/config"
)

func UploadToR2(r2Client *s3.Client, r2Config config.R2Config, file multipart.File, fileHeader *multipart.FileHeader, bucketName, folderPath string, isPublic bool) (string, string, error) {
	if bucketName == "" {
		return "", "", fmt.Errorf("bucketName cannot be empty")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	mimeType := http.DetectContentType(buf.Bytes())
	key := filepath.Join(folderPath, fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename)))

	_, err := r2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to R2: %v", err)
	}

	var fileURL string
	if isPublic {
		publicHost := r2Config.PublicEndpoint
		fileURL = fmt.Sprintf("%s/%s", publicHost, key)
	} else {
		endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2Config.AccountID)
		fileURL = fmt.Sprintf("%s/%s/%s", endpoint, bucketName, key)
	}

	return fileURL, key, nil
}

func DeleteFromR2(r2Client *s3.Client, r2Config config.R2Config, bucketName, key string) {
	if bucketName == "" {
		log.Println("Error deleting orphaned file: bucketName is not set.")
		return
	}

	log.Printf("Attempting to delete orphaned file from R2: %s", key)
	_, err := r2Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to delete orphaned file %s from R2: %v", key, err)
	} else {
		log.Printf("Successfully deleted orphaned file: %s", key)
	}
}
