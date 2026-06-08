package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Storage(
	accessKey string,
	secretKey string,
	region string,
	bucket string,
) (*S3Storage, error) {

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

func (s *S3Storage) UploadFile(
	ctx context.Context,
	key string,
	data []byte,
	contentType string,
) (string, error) {

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		s.bucket,
		s.region,
		key,
	)

	return url, nil
}

func (s *S3Storage) GetBaseURL() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucket, s.region)
}
