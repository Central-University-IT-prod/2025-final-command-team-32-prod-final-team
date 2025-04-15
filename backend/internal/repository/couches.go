package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"slices"
	"solution/internal/database/storage"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/connections/postgres"
	"solution/pkg/utils"
)

var _ contracts.CouchRepository = (*couchRepo)(nil)

type couchRepo struct {
	query *storage.Queries
	pool  *pgxpool.Pool
}

func NewCouchRepository(db *postgres.DB) *couchRepo {
	return &couchRepo{
		query: db.Queries(),
		pool:  db.Pool(),
	}
}

func (c *couchRepo) MarkFilmsAsSeen(ctx context.Context, couchId uuid.UUID, films []uuid.UUID) error {
	params := wrapper.CouchViewForDatabase(couchId, films)
	_, err := c.query.MarkCouchAsViewedBulk(ctx, c.pool, params)
	return err
}

func (c *couchRepo) Exists(ctx context.Context, id uuid.UUID) bool {
	ok, _ := c.query.CouchExists(ctx, c.pool, id)
	return ok
}

func (c *couchRepo) Create(ctx context.Context, param dto.CreateCouch) (uuid.UUID, error) {
	vectors := make([]pgvector.Vector, 0)
	for _, sitter := range param.Sitters {
		vector, err := c.query.GetVector(ctx, c.pool, sitter)
		if err != nil {
			continue
		}
		vectors = append(vectors, vector)
	}
	param.Embedding = utils.AvgVector(vectors)

	id, err := c.query.CreateCouch(ctx, c.pool, wrapper.ToCreateCouch(param))
	if err != nil {
		return uuid.Nil, err
	}

	create := wrapper.SittersConvert(param.Sitters, id)
	_, err = c.query.CreateSitter(ctx, c.pool, create)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (c *couchRepo) GetOne(ctx context.Context, id uuid.UUID) (*dto.CouchView, error) {
	couch, err := c.query.GetCouch(ctx, c.pool, id)
	if err != nil {
		return nil, err
	}
	sitters, err := c.query.GetSitters(ctx, c.pool, couch.ID)
	if err != nil {
		return nil, err
	}
	return wrapper.ToCouchView(couch, sitters), nil
}

func (c *couchRepo) Update(ctx context.Context, id uuid.UUID, param dto.UpdateCouch) error {
	if param.Name != nil {
		err := c.query.UpdateCouch(ctx, c.pool, storage.UpdateCouchParams{
			Name: *param.Name,
			ID:   id,
		})
		if err != nil {
			return err
		}
	}
	if param.Sitters != nil {
		err := c.query.ClearSitters(ctx, c.pool, id)
		if err != nil {
			return err
		}
		rawSitters := *param.Sitters
		if len(rawSitters) == 0 || !slices.Contains(rawSitters, param.AuthorName) {
			rawSitters = append(rawSitters, param.AuthorName)
			param.Sitters = &rawSitters
		}
		sitters := wrapper.SittersConvert(*param.Sitters, id)
		_, err = c.query.CreateSitter(ctx, c.pool, sitters)
		if err != nil {
			return err
		}
		vectors := make([]pgvector.Vector, 0)
		for _, sitter := range *param.Sitters {
			vector, err := c.query.GetVector(ctx, c.pool, sitter)
			if err != nil {
				continue
			}
			vectors = append(vectors, vector)
		}
		err = c.query.SetCouchVector(ctx, c.pool, storage.SetCouchVectorParams{
			Embedding: utils.AvgVector(vectors),
			ID:        id,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *couchRepo) GetMany(ctx context.Context, username string) ([]dto.CouchView, error) {
	couches, err := c.query.GetCouches(ctx, c.pool, username)
	if err != nil {
		return nil, err
	}
	res := make([]dto.CouchView, 0)
	for _, couch := range couches {
		sitters, err := c.query.GetSitters(ctx, c.pool, couch.CouchID)
		if err != nil {
			return nil, err
		}

		res = append(res, wrapper.ToCouchesView(couch, sitters))
	}
	return res, nil
}

func (r *couchRepo) GetFavorites(ctx context.Context, couchId uuid.UUID, limit int64) ([]dto.CinemaView, error) {
	favorites, err := r.query.GetFavorites(ctx, r.pool, storage.GetFavoritesParams{
		SubjectID: couchId,
		Limit:     int32(limit),
	})
	if err != nil {
		return []dto.CinemaView{}, nil
	}

	films := make([]dto.CinemaView, len(favorites))
	for i, favorite := range favorites {
		filmDB, err := r.query.GetCinemaByID(ctx, r.pool, favorite.CinemaID)
		if err != nil {
			return nil, err
		}
		film := wrapper.CinemaWithView(&filmDB)
		rating, err := r.query.GetRate(ctx, r.pool, storage.GetRateParams{
			UserID:   couchId,
			CinemaID: favorite.CinemaID,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		if err == nil {
			film.UserRating = &rating
		}

		films[i] = *film
	}

	return films, nil
}

func (r *couchRepo) SaveFilmToFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error {
	check, _ := r.query.CheckIfFavorite(ctx, r.pool, storage.CheckIfFavoriteParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
	if check {
		return nil
	}
	f, err := r.query.GetCinemaByID(ctx, r.pool, filmId)
	if err != nil {
		return err
	}
	vector, err := r.query.GetCouchVector(ctx, r.pool, couchId)
	if err != nil {
		return err
	}
	newV := utils.SumVectors(&vector, &f.Embedding)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.AdjustVector(ctx, couchId, &newV)
	if err != nil {
		return err
	}

	isBlacklisted, err := r.query.CheckIfBlacklisted(ctx, tx, storage.CheckIfBlacklistedParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	if isBlacklisted {
		err = r.query.DeleteFilmFromBlacklist(ctx, tx, storage.DeleteFilmFromBlacklistParams{
			SubjectID: couchId,
			CinemaID:  filmId,
		})
		if err != nil {
			return err
		}
	}

	err = r.query.SaveFilmToFavorites(ctx, tx, storage.SaveFilmToFavoritesParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})

	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *couchRepo) DeleteFilmFromFavorites(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error {
	return r.query.DeleteFilmFromFavorites(ctx, r.pool, storage.DeleteFilmFromFavoritesParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
}

func (r *couchRepo) SaveFilmToBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error {
	check, _ := r.query.CheckIfBlacklisted(ctx, r.pool, storage.CheckIfBlacklistedParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
	if check {
		return nil
	}

	f, err := r.query.GetCinemaByID(ctx, r.pool, filmId)
	if err != nil {
		return err
	}
	vector, err := r.query.GetCouchVector(ctx, r.pool, couchId)
	if err != nil {
		return err
	}
	newV := utils.SubVectors(&vector, &f.Embedding)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.AdjustVector(ctx, couchId, &newV)
	if err != nil {
		return err
	}

	isFavorite, err := r.query.CheckIfFavorite(ctx, tx, storage.CheckIfFavoriteParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	if isFavorite {
		err = r.query.DeleteFilmFromFavorites(ctx, tx, storage.DeleteFilmFromFavoritesParams{
			SubjectID: couchId,
			CinemaID:  filmId,
		})
		if err != nil {
			return err
		}
	}
	err = r.query.SaveFilmToBlacklist(ctx, tx, storage.SaveFilmToBlacklistParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *couchRepo) DeleteFilmFromBlacklist(ctx context.Context, couchId uuid.UUID, filmId uuid.UUID) error {
	return r.query.DeleteFilmFromBlacklist(ctx, r.pool, storage.DeleteFilmFromBlacklistParams{
		SubjectID: couchId,
		CinemaID:  filmId,
	})
}

func (r *couchRepo) AdjustVector(ctx context.Context, id uuid.UUID, v *pgvector.Vector) error {
	return r.query.UploadVector(ctx, r.pool, storage.UploadVectorParams{ID: id, Embedding: *v})
}

func (r *couchRepo) GetRecommended(ctx context.Context, couchId uuid.UUID, v pgvector.Vector, limit int) ([]*dto.CinemaView, error) {
	cinemas, err := r.query.GetCouchFeed(ctx, r.pool, storage.GetCouchFeedParams{Limit: int32(limit), Embedding: v, SubjectID: couchId})
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
