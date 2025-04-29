package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/service"
)

type S3Repository struct {
	client *s3.S3
	config *models.AWSConfig
}

func NewS3Repository(
	client *s3.S3,
	cfg *models.AWSConfig,
) service.S3Repository {
	return &S3Repository{
		client: client,
		config: cfg,
	}
}

func (r *S3Repository) UploadFile(ctx context.Context, file []byte, fileName string, contentType string) (string, error) {
	_, err := r.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.config.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	fileUrl := fmt.Sprintf("%s/%s/%s", r.config.Endpoint, r.config.Bucket, fileName)
	return fileUrl, nil
}
