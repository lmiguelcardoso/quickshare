package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)


type S3BlobStorage struct {
	s3Client *s3.S3
	bucket   string
	region   string
}

func NewS3BlobStorage(region, bucket, accessKeyID, secretAccessKey string) (*S3BlobStorage, error) {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}

	if accessKeyID != "" && secretAccessKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	log.Printf("S3 connection established: region=%s, bucket=%s", region, bucket)

	return &S3BlobStorage{
		s3Client: s3.New(sess),
		bucket:   bucket,
		region:   region,
	}, nil
}

func (s *S3BlobStorage) GeneratePresignedUploadURL(objectKey string, expiresIn time.Duration) (string, error){
	req, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
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

func (s *S3BlobStorage) GetPublicURL(objectKey string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, objectKey)
}

func (s *S3BlobStorage) ObjectExists(objectKey string) (bool, error) {
	_, err := s.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return false, nil
			default:
				return false, fmt.Errorf("failed to check if object exists: %w", err)
			}
		}
		return false, fmt.Errorf("failed to check if object exists: %w", err)
	}

	return true, nil
}

func (s *S3BlobStorage) GetObjectMetadata(objectKey string) (map[string]string, error) {
	resp, err := s.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	metadata := make(map[string]string)
	for k, v := range resp.Metadata {
		if v != nil {
			metadata[k] = *v
		}
	}

	return metadata, nil
}

func (s *S3BlobStorage) Delete(objectKey string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

