package contracts

import (
	"context"
	"solution/internal/domain/dto"

	"github.com/google/uuid"
)

type FileRepository interface {
	UploadFile(ctx context.Context, fileKey string, fileContent []byte, minioPubHost string) (string, error)
}

type FileService interface {
	UploadFile(ctx context.Context, filmdId uuid.UUID, fileKey string, fileContent []byte) *dto.HttpErr
	GetFileLink(ctx context.Context, cinemaId uuid.UUID) (string, *dto.HttpErr)
}
