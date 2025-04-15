package service

import (
	"context"
	"errors"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"

	"github.com/google/uuid"
)

var _ contracts.CinemaService = (*cinemaService)(nil)

type cinemaService struct {
	cinemaRepo contracts.CinemaReposiroty
}

func NewCinemaService(cr contracts.CinemaReposiroty) *cinemaService {
	return &cinemaService{
		cinemaRepo: cr,
	}
}

func (s *cinemaService) GetById(ctx context.Context, id uuid.UUID) (*dto.CinemaView, *dto.HttpErr) {
	cinemaDB, err := s.cinemaRepo.GetById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, wrapper.NotFoundErr("cinema not found")
	} else if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}

	return cinemaDB, nil
}

func (s *cinemaService) GetRecommended(ctx context.Context, login string, limit int) ([]*dto.CinemaView, *dto.HttpErr) {
	u, err := s.cinemaRepo.GetUserModel(ctx, login)
	if err != nil {
		return nil, wrapper.NotFoundErr(dto.MsgUserNotFound)
	}
	cinemas, err := s.cinemaRepo.GetRecommended(ctx, u.Embedding, u.ID, limit)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	if cinemas == nil {
		cinemas = []*dto.CinemaView{}
	}
	return cinemas, nil
}

func (s *cinemaService) Create(ctx context.Context, c *dto.CinemaCreate, login string) (uuid.UUID, *dto.HttpErr) { // locally or globally depends on user auth{
	if c == nil {
		return uuid.Nil, wrapper.BadRequestErr("not enough data")
	}
	priv := s.cinemaRepo.CheckUserIsPrivileged(ctx, login)
	if !priv && !c.Private {
		return uuid.UUID{}, wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden)
	}
	vector := pgvector.NewVector(make([]float32, 46))
	if !c.Private {
		vector = utils.VectorFromTags(c.Genres)
	}
	// check for admin
	if c.Private {
		c = &dto.CinemaCreate{
			Name:        c.Name,
			PosterURL:   c.PosterURL,
			Description: c.Description,
			Genres:      c.Genres,
			Private:     c.Private,
		}
	}
	id, err := s.cinemaRepo.Create(ctx, c, vector, login)
	if err != nil {
		return id, wrapper.InternalServerErr(err.Error())
	}
	return id, nil
}

func (s *cinemaService) GetTopRated(ctx context.Context, offset, limit int) ([]*dto.CinemaView, *dto.HttpErr) {
	cinemas, err := s.cinemaRepo.GetTopRated(ctx, offset, limit)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return cinemas, nil
}

func (s *cinemaService) SearchFilm(ctx context.Context, titleLike string, tags []string) ([]*dto.CinemaView, *dto.HttpErr) {
	cinemas, err := s.cinemaRepo.SearchFilm(ctx, titleLike, tags)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return cinemas, nil
}

func (s *cinemaService) UpdateFilm(ctx context.Context, params dto.CinemaUpdate, username string) *dto.HttpErr {
	priv := s.cinemaRepo.CheckUserIsPrivileged(ctx, username)
	if !priv {
		return wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden)
	}
	err := s.cinemaRepo.UpdateFilm(ctx, params, username)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *cinemaService) DeleteFilm(ctx context.Context, filmId uuid.UUID, username string) *dto.HttpErr {
	priv := s.cinemaRepo.CheckUserIsPrivileged(ctx, username)
	if !priv {
		return wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden)
	}
	err := s.cinemaRepo.DeleteFilm(ctx, filmId, username)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *cinemaService) GetGenres(ctx context.Context) []string {
	return dto.VectorTag[:]
}

func (s *cinemaService) CheckUserIsPrivileged(ctx context.Context, username string) bool {
	return s.cinemaRepo.CheckUserIsPrivileged(ctx, username)
}
