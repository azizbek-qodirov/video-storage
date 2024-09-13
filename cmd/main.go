package main

import (
	"log"
	"net/http"

	_ "video-service/internal/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"video-service/internal/config"
	"video-service/internal/handler"
	"video-service/internal/pkg/minio"
	"video-service/internal/repository"
	"video-service/internal/service"
)

// @title           Video Service API
// @version         1.0
// @description     This is a simple video service API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	mc, err := minio.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("Error connecting to Minio: %v", err)
	}

	videoRepo := repository.NewVideoRepository(db)

	videoService := service.NewVideoService(videoRepo, mc, cfg)

	handler := handler.NewHandler(videoService)

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	{
		videos := api.Group("/video")
		{
			videos.POST("/upload", handler.UploadVideo)
			videos.GET("/:id", handler.GetVideo)
			videos.DELETE("/:id", handler.DeleteVideo)
		}
		api.GET("/videos", handler.GetAllVideos)
	}

	router.Run(cfg.GATEWAY_HTTP_PORT)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
