package repository

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"solution/internal/domain/contracts"
)

var _ contracts.FileRepository = (*fileRepo)(nil)

type fileRepo struct {
	minioClient *minio.Client
	bucketName  string
}

func NewFileRepo(minioClient *minio.Client, bucketName string) *fileRepo {
	return &fileRepo{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r *fileRepo) UploadFile(ctx context.Context, fileKey string, fileContent []byte, minioPubHost string) (string, error) {
	reader := bytes.NewReader(fileContent)
	_, err := r.minioClient.PutObject(ctx, r.bucketName, fileKey, reader, int64(len(fileContent)), minio.PutObjectOptions{})
	url := "https" + "://" + minioPubHost + r.bucketName + "/" + fileKey
	return url, err
}
