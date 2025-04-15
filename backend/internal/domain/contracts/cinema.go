package contracts

import (
	"context"
	"solution/internal/database/storage"
	"solution/internal/domain/dto"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type CinemaReposiroty interface {
	CheckUserIsPrivileged(ctx context.Context, username string) bool

	Create(ctx context.Context, c *dto.CinemaCreate, vector pgvector.Vector, login string) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.CinemaView, error)
	GetRecommended(ctx context.Context, v pgvector.Vector, userId uuid.UUID, limit int) ([]*dto.CinemaView, error)
	GetTopRated(ctx context.Context, offset, limit int) ([]*dto.CinemaView, error)
	GetUserModel(ctx context.Context, login string) (*storage.User, error)
	UpdateFilm(ctx context.Context, params dto.CinemaUpdate, username string) error
	DeleteFilm(ctx context.Context, filmId uuid.UUID, username string) error

	SetCinemaPic(ctx context.Context, cinemaId uuid.UUID, picName string) error
	GetCinemaPic(ctx context.Context, cinemaId uuid.UUID) (*string, error)
	SearchFilm(ctx context.Context, titleLike string, tags []string) ([]*dto.CinemaView, error)
}

type CinemaService interface {
	CheckUserIsPrivileged(ctx context.Context, username string) bool

	GetById(ctx context.Context, id uuid.UUID) (*dto.CinemaView, *dto.HttpErr)

	UpdateFilm(ctx context.Context, params dto.CinemaUpdate, username string) *dto.HttpErr
	DeleteFilm(ctx context.Context, filmId uuid.UUID, username string) *dto.HttpErr

	GetGenres(ctx context.Context) []string
	GetRecommended(ctx context.Context, login string, limit int) ([]*dto.CinemaView, *dto.HttpErr)
	GetTopRated(ctx context.Context, offset, limit int) ([]*dto.CinemaView, *dto.HttpErr)
	Create(ctx context.Context, c *dto.CinemaCreate, login string) (uuid.UUID, *dto.HttpErr) // locally or globally depends on user auth
	SearchFilm(ctx context.Context, titleLike string, tags []string) ([]*dto.CinemaView, *dto.HttpErr)
}
