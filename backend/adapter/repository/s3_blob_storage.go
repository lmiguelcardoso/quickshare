package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)


type S3BlobStorage struct {
	s3Client  *s3.S3
	bucket    string
	baseURL   string
	uploadDir string
}


func NewS3BlobStorage(s3Client *s3.S3, bucket string, baseURL string, uploadDir string) *S3BlobStorage {
	sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2"),
	})

	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	return &S3BlobStorage{
		s3Client: s3.New(sess),
		bucket: bucket,
		baseURL: baseURL,
		uploadDir: uploadDir,
	}
}

func (s *S3BlobStorage) GeneratePresignedUploadURL(objectKey string, expiresIn time.Duration) (string, error){
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(objectKey),
	})

	log.Printf("generating presigned URL for object: %s", objectKey)

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.Printf("presigned URL generated successfully: %s", urlStr)

	return urlStr, nil
}