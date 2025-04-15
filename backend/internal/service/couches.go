package service

import (
	"context"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/metric"

	"github.com/google/uuid"
)

var _ contracts.CouchService = (*couchService)(nil)

type couchService struct {
	metrics   *metric.PromMetrics
	couchRepo contracts.CouchRepository
}

func NewCouchService(metrics *metric.PromMetrics, repository contracts.CouchRepository) *couchService {
	return &couchService{metrics: metrics, couchRepo: repository}
}

func (c couchService) MarkFilmsAsSeen(ctx context.Context, couchId uuid.UUID, films []uuid.UUID) *dto.HttpErr {
	err := c.couchRepo.MarkFilmsAsSeen(ctx, couchId, films)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}

	c.metrics.TotalSeen.Add(float64(len(films)))

	return nil
}

func (c couchService) GetOne(ctx context.Context, id uuid.UUID) (*dto.CouchView, *dto.HttpErr) {
	couch, err := c.couchRepo.GetOne(ctx, id)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return couch, nil
}

func (c couchService) Create(ctx context.Context, param dto.CreateCouch) (uuid.UUID, *dto.HttpErr) {
	if param.Name == "" || param.AuthorName == "" {
		return uuid.Nil, wrapper.BadRequestErr("name and authorName are required")
	}
	id, err := c.couchRepo.Create(ctx, param)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerErr(err.Error())
	}

	c.metrics.TotalCouches.Add(1)

	return id, nil
}

func (c couchService) Update(ctx context.Context, id uuid.UUID, param dto.UpdateCouch) *dto.HttpErr {
	err := c.couchRepo.Update(ctx, id, param)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (c couchService) GetMany(ctx context.Context, username string) ([]dto.CouchView, *dto.HttpErr) {
	res, err := c.couchRepo.GetMany(ctx, username)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return res, nil
}

func (c couchService) GetFavorites(ctx context.Context, couchId uuid.UUID, limit int64) ([]dto.CinemaView, *dto.HttpErr) {
	res, err := c.couchRepo.GetFavorites(ctx, couchId, limit)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return res, nil
}

func (c couchService) SaveFilmToFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr {
	err := c.couchRepo.SaveFilmToFavorites(ctx, couchId, filmId)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (c couchService) DeleteFilmFromFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr {
	err := c.couchRepo.DeleteFilmFromFavorites(ctx, couchId, filmId)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (c couchService) SaveFilmToBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr {
	err := c.couchRepo.SaveFilmToBlacklist(ctx, couchId, filmId)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (c couchService) DeleteFilmFromBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) *dto.HttpErr {
	err := c.couchRepo.DeleteFilmFromBlacklist(ctx, couchId, filmId)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *couchService) GetRecommended(ctx context.Context, couchId uuid.UUID, limit int) ([]*dto.CinemaView, *dto.HttpErr) {
	if !s.couchRepo.Exists(ctx, couchId) {
		return nil, wrapper.NotFoundErr(dto.MsgCouchNotFound)
	}
	couch, err := s.couchRepo.GetOne(ctx, couchId)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	cinemas, err := s.couchRepo.GetRecommended(ctx, couchId, couch.Embedding, limit)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return cinemas, nil
}
