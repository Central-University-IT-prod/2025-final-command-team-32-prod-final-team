package repository

import (
	"context"
	"fmt"

	"solution/internal/database/storage"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/connections/postgres"
	"solution/pkg/logger"
	"solution/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

var _ contracts.UserRepository = (*userRepo)(nil)

type userRepo struct {
	query *storage.Queries
	pool  *pgxpool.Pool
}

func NewUserRepository(db *postgres.DB) *userRepo {
	return &userRepo{
		query: db.Queries(),
		pool:  db.Pool(),
	}
}

func (r *userRepo) Exists(ctx context.Context, email string) bool {
	ok, _ := r.query.UserExists(ctx, r.pool, email)
	return ok
}

func (r *userRepo) Create(ctx context.Context, u *dto.UserAuth) (uuid.UUID, error) {
	s := make([]float32, 46)
	params := wrapper.UserAuthWithCreateParams(u, pgvector.NewVector(s))
	err := r.query.CreateUser(ctx, r.pool, *params)
	if err != nil {
		return uuid.UUID{}, err
	}
	user, _ := r.query.GetUserByLogin(ctx, r.pool, u.Login)
	return user.ID, nil
}

func (r *userRepo) GetByLogin(ctx context.Context, email string) (*dto.UserView, error) {
	u, err := r.query.GetUserByLogin(ctx, r.pool, email)
	logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("pass: %s", *u.Password))
	return wrapper.UserWithView(&u), err
}

func (r *userRepo) GetById(ctx context.Context, id uuid.UUID) (*dto.UserView, error) {
	u, err := r.query.GetUserById(ctx, r.pool, id)
	return wrapper.UserWithView(&u), err
}

func (r *userRepo) SaveFilmToFavorites(ctx context.Context, login string, filmId uuid.UUID) error {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return err
	}
	check, _ := r.query.CheckIfFavorite(ctx, r.pool, storage.CheckIfFavoriteParams{
		SubjectID: u.ID,
		CinemaID:  filmId,
	})
	if check {
		return nil
	}
	f, err := r.query.GetCinemaByID(ctx, r.pool, filmId)
	if err != nil {
		return err
	}
	newV := utils.SumVectors(&u.Embedding, &f.Embedding)
	SubjectID := u.ID

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.AdjustVector(ctx, u.ID, &newV)
	if err != nil {
		return err
	}

	isBlacklisted, err := r.query.CheckIfBlacklisted(ctx, tx, storage.CheckIfBlacklistedParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	if isBlacklisted {
		err = r.query.DeleteFilmFromBlacklist(ctx, tx, storage.DeleteFilmFromBlacklistParams{
			SubjectID: SubjectID,
			CinemaID:  filmId,
		})
		if err != nil {
			return err
		}
	}

	err = r.query.SaveFilmToFavorites(ctx, tx, storage.SaveFilmToFavoritesParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})

	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *userRepo) DeleteFilmFromFavorites(ctx context.Context, login string, filmId uuid.UUID) error {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return err
	}
	SubjectID := u.ID
	return r.query.DeleteFilmFromFavorites(ctx, r.pool, storage.DeleteFilmFromFavoritesParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})
}

func (r *userRepo) SaveFilmToBlacklist(ctx context.Context, login string, filmId uuid.UUID) error {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return err
	}
	check, _ := r.query.CheckIfBlacklisted(ctx, r.pool, storage.CheckIfBlacklistedParams{
		SubjectID: u.ID,
		CinemaID:  filmId,
	})
	if check {
		return nil
	}

	f, err := r.query.GetCinemaByID(ctx, r.pool, filmId)
	if err != nil {
		return err
	}
	newV := utils.SubVectors(&u.Embedding, &f.Embedding)
	SubjectID := u.ID

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.AdjustVector(ctx, u.ID, &newV)
	if err != nil {
		return err
	}

	isFavorite, err := r.query.CheckIfFavorite(ctx, tx, storage.CheckIfFavoriteParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	if isFavorite {
		err = r.query.DeleteFilmFromFavorites(ctx, tx, storage.DeleteFilmFromFavoritesParams{
			SubjectID: SubjectID,
			CinemaID:  filmId,
		})
		if err != nil {
			return err
		}
	}
	err = r.query.SaveFilmToBlacklist(ctx, tx, storage.SaveFilmToBlacklistParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *userRepo) DeleteFilmFromBlacklist(ctx context.Context, login string, filmId uuid.UUID) error {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return err
	}
	SubjectID := u.ID
	return r.query.DeleteFilmFromBlacklist(ctx, r.pool, storage.DeleteFilmFromBlacklistParams{
		SubjectID: SubjectID,
		CinemaID:  filmId,
	})
}

func (r *userRepo) SaveRate(ctx context.Context, login string, filmId uuid.UUID, rate int32) error {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return err
	}
	newV := utils.MultiplyVector(&u.Embedding, utils.AdjustRatingWeigth(rate))
	userId := u.ID

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	r.AdjustVector(ctx, u.ID, &newV)
	rateExists, err := r.query.RateExists(ctx, tx, storage.RateExistsParams{
		UserID:   userId,
		CinemaID: filmId,
	})
	if err != nil {
		return err
	}
	if rateExists {
		return r.query.UpdateRate(ctx, tx, storage.UpdateRateParams{
			UserID:   userId,
			CinemaID: filmId,
			Upd:      rate,
		})
	} else {
		err = r.query.MarkAsViewed(ctx, tx, storage.MarkAsViewedParams{
			SubjectID: userId,
			CinemaID:  filmId,
		})
		if err != nil {
			return err
		}
	}
	err = r.query.SaveRate(ctx, tx, storage.SaveRateParams{
		UserID:   userId,
		CinemaID: filmId,
		Upd:      rate,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *userRepo) GetFavorites(ctx context.Context, login string, limit int64) ([]dto.CinemaView, error) {
	u, err := r.query.GetUserByLogin(ctx, r.pool, login)
	if err != nil {
		return []dto.CinemaView{}, err
	}
	userId := u.ID

	favorites, err := r.query.GetFavorites(ctx, r.pool, storage.GetFavoritesParams{
		SubjectID: userId,
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
			UserID:   userId,
			CinemaID: favorite.CinemaID,
		})
		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		}
		if err == nil {
			film.UserRating = &rating
		}

		films[i] = *film
	}

	return films, nil
}

func (r *userRepo) MarkFilmsAsSeen(ctx context.Context, SubjectID uuid.UUID, ids []uuid.UUID) error {
	params := wrapper.ViewForDatabase(SubjectID, ids)
	_, err := r.query.MarkAsViewedBulk(ctx, r.pool, params)
	return err
}

func (r *userRepo) IsFavorited(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (bool, error) {
	return r.query.IsFavorite(ctx, r.pool, storage.IsFavoriteParams{
		SubjectID: userId,
		CinemaID:  filmId,
	})
}

func (r *userRepo) IsBlacklisted(ctx context.Context, userId uuid.UUID, filmId uuid.UUID) (bool, error) {
	return r.query.IsBlacklisted(ctx, r.pool, storage.IsBlacklistedParams{
		SubjectID: userId,
		CinemaID:  filmId,
	})
}

func (r *userRepo) FilmExists(ctx context.Context, filmId uuid.UUID) bool {
	ok, _ := r.query.FilmExists(ctx, r.pool, filmId)
	return ok
}

func (r *userRepo) GetFilmById(ctx context.Context, filmId uuid.UUID) (*dto.CinemaView, error) {
	c, err := r.query.GetCinemaByID(ctx, r.pool, filmId)
	if err != nil {
		return nil, err
	}
	return wrapper.CinemaWithView(&c), nil
}

func (r *userRepo) AdjustVector(ctx context.Context, id uuid.UUID, v *pgvector.Vector) error {
	return r.query.UploadVector(ctx, r.pool, storage.UploadVectorParams{ID: id, Embedding: *v})
}

func (r *userRepo) GetVectorByLogin(ctx context.Context, login string) (*pgvector.Vector, error) {
	vect, err := r.query.GetVector(ctx, r.pool, login)
	if err != nil {
		return nil, err
	}
	return &vect, nil
}

func (r *userRepo) SearchUser(ctx context.Context, userLike string) ([]*dto.UserView, error) {
	logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("query: %s", userLike))
	users, err := r.query.UserSearch(ctx, r.pool, userLike)
	if err != nil {
		return nil, err
	}
	usersV := make([]*dto.UserView, 0)
	for _, u := range users {
		usersV = append(usersV, wrapper.UserWithView(&u))
	}
	return usersV, nil
}
