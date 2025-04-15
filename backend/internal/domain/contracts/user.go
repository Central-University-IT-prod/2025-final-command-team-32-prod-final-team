package contracts

import (
	"context"

	"solution/internal/domain/dto"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type UserRepository interface {
	Exists(ctx context.Context, login string) bool
	FilmExists(ctx context.Context, filmId uuid.UUID) bool

	Create(ctx context.Context, u *dto.UserAuth) (uuid.UUID, error)
	GetByLogin(ctx context.Context, login string) (*dto.UserView, error)
	GetVectorByLogin(ctx context.Context, login string) (*pgvector.Vector, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.UserView, error)
	GetFilmById(ctx context.Context, filmId uuid.UUID) (*dto.CinemaView, error)

	AdjustVector(ctx context.Context, id uuid.UUID, v *pgvector.Vector) error

	GetFavorites(ctx context.Context, login string, limit int64) ([]dto.CinemaView, error)
	SaveFilmToFavorites(ctx context.Context, login string, filmId uuid.UUID) error
	DeleteFilmFromFavorites(ctx context.Context, login string, filmId uuid.UUID) error
	SaveFilmToBlacklist(ctx context.Context, login string, filmId uuid.UUID) error
	DeleteFilmFromBlacklist(ctx context.Context, login string, filmId uuid.UUID) error

	SaveRate(ctx context.Context, login string, filmId uuid.UUID, rate int32) error
	MarkFilmsAsSeen(ctx context.Context, userId uuid.UUID, ids []uuid.UUID) error

	SearchUser(ctx context.Context, userLike string) ([]*dto.UserView, error)
}

type UserService interface {
	Register(ctx context.Context, u *dto.UserAuth) (*dto.UserAuthResponse, *dto.HttpErr)
	Login(ctx context.Context, uLogin *dto.UserAuth) (*dto.UserAuthResponse, *dto.HttpErr)
	GetVector(ctx context.Context, login string) (*pgvector.Vector, *dto.HttpErr)
	GetProfile(ctx context.Context, login string) (*dto.UserView, *dto.HttpErr)

	SetVector(ctx context.Context, login string, genres []string) *dto.HttpErr

	GetFavorites(ctx context.Context, login string, limit int64) ([]dto.CinemaView, *dto.HttpErr)
	SaveFilmToFavorites(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr
	DeleteFilmFromFavorites(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr
	SaveFilmToBlacklist(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr
	DeleteFilmFromBlacklist(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr

	SaveRate(ctx context.Context, login string, filmId uuid.UUID, rate int32) *dto.HttpErr
	MarkFilmsAsSeen(ctx context.Context, userId uuid.UUID, ids []uuid.UUID) *dto.HttpErr

	SearchUser(ctx context.Context, userLike string) ([]*dto.UserView, *dto.HttpErr)
}
