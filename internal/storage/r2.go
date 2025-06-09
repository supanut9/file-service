package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Pass the file size and the multipart.File directly for streaming
func UploadToR2(r2Client *s3.Client, file multipart.File, fileHeader *multipart.FileHeader, bucketName, folderPath string, isPublic bool) (string, string, error) {
	if bucketName == "" {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		if bucketName == "" {
			return "", "", fmt.Errorf("R2_BUCKET_NAME is not set and no bucketName provided")
		}
	}

	head := make([]byte, 512)
	_, err := file.Read(head)
	if err != nil && err != io.EOF {
		return "", "", fmt.Errorf("failed to read file header: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("failed to seek file: %v", err)
	}

	mimeType := http.DetectContentType(head)
	key := filepath.Join(folderPath, fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename)))

	_, err = r2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String(mimeType),
		ContentLength: &fileHeader.Size,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to R2: %v", err)
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

	return fileURL, key, nil
}

func DeleteFromR2(r2Client *s3.Client, bucketName, key string) {
	if bucketName == "" {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		if bucketName == "" {
			log.Println("Error deleting orphaned file: R2_BUCKET_NAME is not set.")
			return
		}
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
