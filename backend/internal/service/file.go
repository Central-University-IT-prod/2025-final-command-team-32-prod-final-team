package service

import (
	"context"
	"path/filepath"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ contracts.FileService = (*fileService)(nil)

type fileService struct {
	fileRepo        contracts.FileRepository
	cinemaRepo      contracts.CinemaReposiroty
	minioPublicHost string
}

func NewFileService(fr contracts.FileRepository, cr contracts.CinemaReposiroty, minioPubHost string) *fileService {
	return &fileService{
		fileRepo:        fr,
		cinemaRepo:      cr,
		minioPublicHost: minioPubHost,
	}
}

func (s *fileService) UploadFile(ctx context.Context,
	filmId uuid.UUID,
	fileName string,
	fileContent []byte) *dto.HttpErr {
	_, err := s.cinemaRepo.GetById(ctx, filmId)
	if err == pgx.ErrNoRows {
		return wrapper.NotFoundErr("cinema not found")
	}

	// Generate file key
	id := uuid.New().String()
	ext := filepath.Ext(fileName)
	fileKey := id + ext

	url, err := s.fileRepo.UploadFile(ctx, fileKey, fileContent, s.minioPublicHost)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}

	err = s.cinemaRepo.SetCinemaPic(ctx, filmId, url)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}

	return nil
}

func (s *fileService) GetFileLink(ctx context.Context, cinemaId uuid.UUID) (string, *dto.HttpErr) {
	picURL, err := s.getPicURL(ctx, cinemaId)
	if err != nil {
		return "", wrapper.InternalServerErr(err.Error())
	}

	return picURL, nil
}

func (s *fileService) getPicURL(ctx context.Context, cinemaId uuid.UUID) (string, error) {
	// Get picture id from db
	picURL, err := s.cinemaRepo.GetCinemaPic(ctx, cinemaId)
	if err != nil || picURL == nil {
		return "", nil
	}

	// Return pic URL
	return *picURL, nil
}
