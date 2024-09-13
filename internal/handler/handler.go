package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"video-service/internal/service"
)

type Handler struct {
	videoService *service.VideoService
}

func NewHandler(videoService *service.VideoService) *Handler {
	return &Handler{videoService: videoService}
}

// UploadVideo godoc
// @Summary Upload a video
// @Description Uploads a new video file.
// @Tags Videos
// @Accept mpfd
// @Produce json
// @Param file formData file true "Video file"
// @Success 201 {object} model.Video "Video uploaded successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /video/upload [post]
func (h *Handler) UploadVideo(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100*1024*1024) // 100 MB

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), getExt(header.Filename))

	tmpFile, err := os.Create(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	video, err := h.videoService.Upload(c, tmpFile.Name(), fileName, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}

// GetVideo godoc
// @Summary Get a video by ID
// @Description Retrieves a video by its unique ID.
// @Tags Videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} model.Video "Video retrieved successfully"
// @Failure 404 {object} string "Video not found"
// @Failure 500 {object} string "Server error"
// @Router /video/{id} [get]
func (h *Handler) GetVideo(c *gin.Context) {
	videoID := c.Param("id")

	video, err := h.videoService.Get(c, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if video == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	c.JSON(http.StatusOK, video)
}

// DeleteVideo godoc
// @Summary Delete a video by ID
// @Description Deletes a video by its unique ID.
// @Tags Videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 204 {object} string "Video deleted successfully"
// @Failure 404 {object} string "Video not found"
// @Failure 500 {object} string "Server error"
// @Router /video/{id} [delete]
func (h *Handler) DeleteVideo(c *gin.Context) {
	videoID := c.Param("id")

	err := h.videoService.Delete(c, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllVideos godoc
// @Summary Get all videos
// @Description Retrieves a list of all videos.
// @Tags Videos
// @Produce json
// @Success 200 {array} model.Video "Videos retrieved successfully"
// @Failure 500 {object} string "Server error"
// @Router /videos [get]
func (h *Handler) GetAllVideos(c *gin.Context) {
	videos, err := h.videoService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func getExt(filename string) string {
	for i := len(filename) - 1; i >= 0 && !os.IsPathSeparator(filename[i]); i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}
