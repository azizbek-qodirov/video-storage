package service

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"video-service/internal/config"
	"video-service/internal/model"
	"video-service/internal/repository"
)

type VideoService struct {
	repo        *repository.VideoRepository
	minioClient *minio.Client
	cfg         *config.Config
}

func NewVideoService(repo *repository.VideoRepository, minioClient *minio.Client, cfg *config.Config) *VideoService {
	return &VideoService{repo: repo, minioClient: minioClient, cfg: cfg}
}

func (s *VideoService) Upload(c *gin.Context, filePath string, fileName string, fileSize int64) (*model.Video, error) {
	objectName := fileName
	bucketName := s.cfg.MinioBucketName

	_, err := s.minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Minio: %w", err)
	}

	video := &model.Video{
		ID:   uuid.New(),
		Name: fileName,
		Size: fileSize,
		URL:  fmt.Sprintf("%s/%s/%s", s.cfg.MinioEndpoint, bucketName, objectName),
	}

	err = s.repo.Create(c, video)
	if err != nil {
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	os.Remove(filePath)

	return video, nil
}

func (s *VideoService) Get(c *gin.Context, videoID string) (*model.Video, error) {
	video, err := s.repo.GetByID(c, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return video, nil
}

func (s *VideoService) Delete(c *gin.Context, videoID string) error {
	video, err := s.repo.GetByID(c, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	err = s.minioClient.RemoveObject(
		context.Background(),
		s.cfg.MinioBucketName,
		video.Name,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete video from Minio: %w", err)
	}

	err = s.repo.Delete(c, videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video record: %w", err)
	}

	return nil
}

func (s *VideoService) GetAll(c *gin.Context) ([]*model.Video, error) {
	videos, err := s.repo.GetAll(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}

	return videos, nil
}
