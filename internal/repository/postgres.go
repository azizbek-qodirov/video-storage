package repository

import (
	"context"
	"fmt"

	"video-service/internal/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type VideoRepository struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *model.Video) error {
	query := `INSERT INTO videos (id, name, size, url) VALUES ($1, $2, $3, $4)`

	if _, err := r.db.ExecContext(ctx, query, video.ID, video.Name, video.Size, video.URL); err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	return nil
}

func (r *VideoRepository) GetByID(ctx context.Context, videoID string) (*model.Video, error) {
	var video model.Video

	query := `SELECT * FROM videos WHERE id = $1`

	err := r.db.GetContext(ctx, &video, query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video by ID: %w", err)
	}

	return &video, nil
}

func (r *VideoRepository) Delete(ctx context.Context, videoID string) error {
	query := `DELETE FROM videos WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, query, videoID); err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	return nil
}

func (r *VideoRepository) GetAll(ctx context.Context) ([]*model.Video, error) {
	var videos []*model.Video

	query := `SELECT * FROM videos`

	if err := r.db.SelectContext(ctx, &videos, query); err != nil {
		return nil, fmt.Errorf("failed to get all videos: %w", err)
	}

	return videos, nil
}
