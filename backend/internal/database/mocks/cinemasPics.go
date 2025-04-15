package mocks

import (
	"context"
	"log"
	"math/rand"
	"os"
	"solution/internal/database/storage"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/connections/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type picMocks struct {
	query       *storage.Queries
	pool        *pgxpool.Pool
	fileService contracts.FileService
}

func NewPicMocks(db *postgres.DB, fileService contracts.FileService) *picMocks {
	return &picMocks{
		query:       db.Queries(),
		pool:        db.Pool(),
		fileService: fileService,
	}
}

func (m *picMocks) getFilms(ctx context.Context) ([]dto.CinemaView, error) {
	cinemasDB, err := m.query.GetAllCinemas(ctx, m.pool)
	if err != nil {
		return []dto.CinemaView{}, err
	}
	cinemas := make([]dto.CinemaView, len(cinemasDB))
	for i, cinemaDB := range cinemasDB {
		cinemas[i] = *wrapper.CinemaWithView(&cinemaDB)
	}
	return cinemas, nil
}

func (m *picMocks) setPic(ctx context.Context, filmId uuid.UUID, fileKey string, fileContent []byte) *dto.HttpErr {
	return m.fileService.UploadFile(ctx, filmId, fileKey, fileContent)
}

func (m *picMocks) SetRandomPicsForCinemas() {
	ctx := context.Background()
	cinemas, err := m.getFilms(ctx)
	if err != nil {
		log.Fatalln("failed to get all cinemas")
	}

	dir := "./internal/database/mocks/images"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln("failed to read dir")
	}

	if len(files) == 0 {
		log.Fatalln("0 files in images dir")
	}

	for _, cinema := range cinemas {
		file := files[rand.Intn(len(files))]

		content, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			log.Println("failed to read file:", err)
			continue
		}

		if err := m.setPic(ctx, cinema.Id, file.Name(), content); err != nil {
			log.Println("failed to set pic")
		}
	}
}
