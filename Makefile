run:
	go run cmd/main.go

drun:
	docker run -p 9000:9000 -p 9001:9001 --name minio -v ~/minio/data:/data -e "MINIO_ROOT_USER=user" -e "MINIO_ROOT_PASSWORD=password" quay.io/minio/minio server /data --console-address ":9001"

swag-gen:
	~/go/bin/swag init -g ./cmd/main.go -o ./internal/docs force 1
