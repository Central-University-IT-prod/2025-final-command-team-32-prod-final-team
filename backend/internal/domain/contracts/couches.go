package contracts

import (
	"context"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"solution/internal/domain/dto"
)

type CouchRepository interface {
	Exists(ctx context.Context, id uuid.UUID) bool
	Create(ctx context.Context, param dto.CreateCouch) (uuid.UUID, error)
	GetOne(ctx context.Context, id uuid.UUID) (*dto.CouchView, error)

	Update(ctx context.Context, id uuid.UUID, param dto.UpdateCouch) error
	GetMany(ctx context.Context, username string) ([]dto.CouchView, error)

	GetRecommended(ctx context.Context, couchId uuid.UUID, v pgvector.Vector, limit int) ([]*dto.CinemaView, error)

	GetFavorites(ctx context.Context, couchId uuid.UUID, limit int64) ([]dto.CinemaView, error)
	SaveFilmToFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error
	DeleteFilmFromFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error
	SaveFilmToBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error
	DeleteFilmFromBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error
	AdjustVector(ctx context.Context, id uuid.UUID, v *pgvector.Vector) error

	MarkFilmsAsSeen(ctx context.Context, couchId uuid.UUID, films []uuid.UUID) error
}

type CouchService interface {
	Create(ctx context.Context, param dto.CreateCouch) (uuid.UUID, *dto.HttpErr)
	GetOne(ctx context.Context, id uuid.UUID) (*dto.CouchView, *dto.HttpErr)

	Update(ctx context.Context, id uuid.UUID, param dto.UpdateCouch) *dto.HttpErr
	GetMany(ctx context.Context, username string) ([]dto.CouchView, *dto.HttpErr)

	GetFavorites(ctx context.Context, couchId uuid.UUID, limit int64) ([]dto.CinemaView, *dto.HttpErr)
	SaveFilmToFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr
	DeleteFilmFromFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr
	SaveFilmToBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr
	DeleteFilmFromBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr

	GetRecommended(ctx context.Context, couchId uuid.UUID, limit int) ([]*dto.CinemaView, *dto.HttpErr)

	MarkFilmsAsSeen(ctx context.Context, couchId uuid.UUID, films []uuid.UUID) *dto.HttpErr
}
