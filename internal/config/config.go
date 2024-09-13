package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	GATEWAY_HTTP_PORT string
	REDIS_URL         string
	PostgresUser      string
	PostgresPassword  string
	PostgresHost      string
	PostgresPort      int
	PostgresDatabase  string
	MinioEndpoint     string
	MinioAccessKey    string
	MinioSecretKey    string
	MinioBucketName   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := &Config{}

	config.REDIS_URL = cast.ToString(coalesce("REDIS_URL", "localhost:6379"))
	config.GATEWAY_HTTP_PORT = cast.ToString(coalesce("GATEWAY_HTTP_PORT", ":8088"))

	config.PostgresUser = cast.ToString(coalesce("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(coalesce("POSTGRES_PASSWORD", "postgres"))
	config.PostgresHost = cast.ToString(coalesce("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToInt(coalesce("POSTGRES_PORT", 5432))
	config.PostgresDatabase = cast.ToString(coalesce("POSTGRES_DB", "video_service"))

	config.MinioEndpoint = cast.ToString(coalesce("MINIO_ENDPOINT", "localhost:9000"))
	config.MinioAccessKey = cast.ToString(coalesce("MINIO_ACCESS_KEY", "minioadmin"))
	config.MinioSecretKey = cast.ToString(coalesce("MINIO_SECRET_KEY", "minioadmin"))
	config.MinioBucketName = cast.ToString(coalesce("MINIO_BUCKET_NAME", "videos"))

	return config
}

func coalesce(key string, defaultValues ...interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists && val != "" {
		return val
	}
	for _, defaultValue := range defaultValues {
		if defaultValue != nil && fmt.Sprintf("%v", defaultValue) != "" {
			return defaultValue
		}
	}
	return nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDatabase)
}
