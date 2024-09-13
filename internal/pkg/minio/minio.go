package minio

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"video-service/internal/config"
)

func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Minio client: %w", err)
	}

	bucketExists, errBucketExists := minioClient.BucketExists(context.Background(), cfg.MinioBucketName)
	if errBucketExists != nil {
		log.Printf("Error checking if bucket exists: %v", errBucketExists)
		return nil, errBucketExists
	}
	if !bucketExists {
		err = minioClient.MakeBucket(context.Background(), cfg.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("error creating bucket: %w", err)
		}
		log.Printf("Bucket '%s' created successfully", cfg.MinioBucketName)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.MinioBucketName)

	err = minioClient.SetBucketPolicy(context.Background(), cfg.MinioBucketName, policy)
	if err != nil {
		log.Println("error while setting bucket policy : ", err)
		return nil, err
	}

	return minioClient, nil
}
