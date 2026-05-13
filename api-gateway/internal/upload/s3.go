package upload

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Uploader(accessKey, secretKey, region, bucket string) (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return &S3Uploader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		region: region,
	}, nil
}

func (u *S3Uploader) UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("uploads/%d-%s%s", time.Now().Unix(), uuid.New().String(), ext)

	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   file,
		// ACL:    "public-read", // Uncomment if you want the file to be public
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	return url, nil
}
