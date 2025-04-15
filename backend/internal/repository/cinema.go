package repository

import (
	"context"
	"solution/internal/database/storage"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/connections/postgres"

	"github.com/pgvector/pgvector-go"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ contracts.CinemaReposiroty = (*cinemaRepo)(nil)

type cinemaRepo struct {
	query *storage.Queries
	pool  *pgxpool.Pool
}

func NewCinemaRepository(db *postgres.DB) *cinemaRepo {
	return &cinemaRepo{
		query: db.Queries(),
		pool:  db.Pool(),
	}
}

func (r cinemaRepo) Create(ctx context.Context, c *dto.CinemaCreate, vector pgvector.Vector, login string) (uuid.UUID, error) {
	var posterURL *string
	if c.PosterURL != "" {
		posterURL = &c.PosterURL
	}
	var description *string
	if c.Description != "" {
		description = &c.Description
	}
	var rating *float64
	if c.Rating != 0 {
		v := float64(c.Rating)
		rating = &v
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback(ctx)

	params := wrapper.ToCreateCinema(c, posterURL, description, rating, vector)
	res, err := r.query.CreateCinema(ctx, tx, params)
	if err != nil {
		return uuid.Nil, err
	}

	if c.Private {
		user, err := r.query.GetUserByLogin(ctx, tx, login)
		if err != nil {
			return uuid.UUID{}, err
		}
		err = r.query.SaveFilmToFavorites(ctx, tx, storage.SaveFilmToFavoritesParams{
			SubjectID: user.ID,
			CinemaID:  res,
		})
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return res, nil
}

func (r cinemaRepo) GetById(ctx context.Context, id uuid.UUID) (*dto.CinemaView, error) {
	cinemaDB, err := r.query.GetCinemaByID(ctx, r.pool, id)
	if err != nil {
		return nil, err
	}
	return wrapper.CinemaWithView(&cinemaDB), nil
}

func (r cinemaRepo) GetRecommended(ctx context.Context, v pgvector.Vector, userId uuid.UUID, limit int) ([]*dto.CinemaView, error) {

	cinemas, err := r.query.GetFeed(ctx, r.pool, storage.GetFeedParams{Limit: int32(limit), Embedding: v, SubjectID: userId})
	if err != nil {
		return nil, err
	}
	var cinemasV []*dto.CinemaView
	for _, c := range cinemas {
		cv := wrapper.FeedWitView(&c)
		cinemasV = append(cinemasV, cv)
	}
	return cinemasV, nil
}

func (r *cinemaRepo) GetUserModel(ctx context.Context, login string) (*storage.User, error) {
	u, err := r.query.GetUserModelByLogin(ctx, r.pool, login)
	return &u, err
}

func (r *cinemaRepo) GetTopRated(ctx context.Context, offset, limit int) ([]*dto.CinemaView, error) {
	cinemas, err := r.query.GetTopRated(ctx, r.pool, storage.GetTopRatedParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	cinemasV := []*dto.CinemaView{}
	for _, c := range cinemas {
		cv := wrapper.CinemaWithView(&c)
		cinemasV = append(cinemasV, cv)
	}
	return cinemasV, nil
}

func (r *cinemaRepo) SetCinemaPic(ctx context.Context, cinemaId uuid.UUID, picName string) error {
	return r.query.SetCinemaPic(ctx, r.pool, storage.SetCinemaPicParams{
		ID:        cinemaId,
		PosterUrl: &picName,
	})
}

func (r *cinemaRepo) GetCinemaPic(ctx context.Context, cinemaId uuid.UUID) (*string, error) {
	return r.query.GetCinemaPic(ctx, r.pool, cinemaId)
}

func (r *cinemaRepo) SearchFilm(ctx context.Context, titleLike string, tags []string) ([]*dto.CinemaView, error) {
	var err error
	var cinemas []storage.Cinema
	switch {
	case len(tags) == 0:
		cinemas, err = r.query.SearchFilm(ctx, r.pool, titleLike)
		if err != nil {
			return nil, err
		}
	case titleLike == "":
		cinemas, err = r.query.SearchTagsOnly(ctx, r.pool, tags)
		if err != nil {
			return nil, err
		}
	case titleLike != "":
		cinemas, err = r.query.SearchFilmWitTags(ctx, r.pool, storage.SearchFilmWitTagsParams{Genres: tags, Lower: titleLike})
		if err != nil {
			return nil, err
		}
	}

	cinemasV := make([]*dto.CinemaView, 0)
	for _, c := range cinemas {
		cv := wrapper.CinemaWithView(&c)
		cinemasV = append(cinemasV, cv)
	}
	return cinemasV, nil
}

func (r *cinemaRepo) CheckUserIsPrivileged(ctx context.Context, username string) bool {
	ok, _ := r.query.UserExists(ctx, r.pool, username)
	if !ok {
		return false
	}
	ok, _ = r.query.UserIsPrivileged(ctx, r.pool, username)
	return ok
}

func (r *cinemaRepo) UpdateFilm(ctx context.Context, params dto.CinemaUpdate, username string) error {
	err := r.query.UpdateCinema(ctx, r.pool, wrapper.ToUpdateCinema(params))
	if err != nil {
		return err
	}
	return nil
}

func (r *cinemaRepo) DeleteFilm(ctx context.Context, filmId uuid.UUID, username string) error {
	err := r.query.DeleteCinema(ctx, r.pool, filmId)
	if err != nil {
		return err
	}
	return nil
}
