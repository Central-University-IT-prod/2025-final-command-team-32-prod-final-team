package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/logger"
	"solution/pkg/metric"
	"solution/pkg/utils"
)

type userService struct {
	metrics     *metric.PromMetrics
	userRepo    contracts.UserRepository
	authService contracts.AuthService
}

func NewUserService(metrics *metric.PromMetrics, ur contracts.UserRepository, as contracts.AuthService) *userService {
	return &userService{
		metrics:     metrics,
		userRepo:    ur,
		authService: as,
	}
}

func (s *userService) Register(ctx context.Context, u *dto.UserAuth) (*dto.UserAuthResponse, *dto.HttpErr) {
	if s.userRepo.Exists(ctx, u.Login) {
		return nil, wrapper.ConflictErr(dto.MsgUserAlreadyExists)
	}

	_, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}

	token, httpErr := s.authService.GenerateToken(ctx, u.Login)
	if httpErr != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("generate token failed with ERR: %s", httpErr.Error()))
		return nil, httpErr
	}

	s.metrics.TotalRegistered.Add(1)

	return &dto.UserAuthResponse{Token: token}, nil
}

func (s *userService) Login(ctx context.Context, uLogin *dto.UserAuth) (*dto.UserAuthResponse, *dto.HttpErr) {
	user, err := s.userRepo.GetByLogin(ctx, uLogin.Login)
	if err != nil {
		fmt.Println(err.Error())
		return nil, wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	ok := utils.CompareHashAndPassword(user.PasswordHashed, uLogin.PasswordUnhashed)
	if !ok {
		return nil, wrapper.BadRequestErr(dto.MsgInvalidPassword)
	}

	token, httpErr := s.authService.GenerateToken(ctx, uLogin.Login)
	if httpErr != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("generate token failed with ERR: %s", httpErr.Error()))
		return nil, httpErr
	}

	return &dto.UserAuthResponse{Token: token}, nil
}

func (s *userService) GetProfile(ctx context.Context, email string) (*dto.UserView, *dto.HttpErr) {
	if !s.userRepo.Exists(ctx, email) {
		return nil, wrapper.NotFoundErr(dto.MsgUserNotFound)
	}
	user, err := s.userRepo.GetByLogin(ctx, email)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to get user from db with ERR: %s", err.Error()))
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return user, nil
}

func (s *userService) SaveFilmToFavorites(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	if !s.userRepo.FilmExists(ctx, filmId) {
		return wrapper.NotFoundErr(dto.MsgFilmNotFound)
	}

	err := s.userRepo.SaveFilmToFavorites(ctx, login, filmId)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to save film to favorites with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}

	s.metrics.TotalLikedFilms.Add(1)

	return nil
}

func (s *userService) DeleteFilmFromFavorites(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	if !s.userRepo.FilmExists(ctx, filmId) {
		return wrapper.NotFoundErr(dto.MsgFilmNotFound)
	}

	err := s.userRepo.DeleteFilmFromFavorites(ctx, login, filmId)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to delete film from favorites with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *userService) SaveFilmToBlacklist(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	if !s.userRepo.FilmExists(ctx, filmId) {
		return wrapper.NotFoundErr(dto.MsgFilmNotFound)
	}

	err := s.userRepo.SaveFilmToBlacklist(ctx, login, filmId)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to save film to blacklist with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}

	s.metrics.TotalDislikedFilms.Add(1)

	return nil
}

func (s *userService) DeleteFilmFromBlacklist(ctx context.Context, login string, filmId uuid.UUID) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	if !s.userRepo.FilmExists(ctx, filmId) {
		return wrapper.NotFoundErr(dto.MsgFilmNotFound)
	}

	err := s.userRepo.DeleteFilmFromBlacklist(ctx, login, filmId)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to delete film from blacklist with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *userService) SaveRate(ctx context.Context, login string, filmId uuid.UUID, rate int32) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	if !s.userRepo.FilmExists(ctx, filmId) {
		return wrapper.NotFoundErr(dto.MsgFilmNotFound)
	}

	err := s.userRepo.SaveRate(ctx, login, filmId, rate)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to save film rate with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *userService) GetFavorites(ctx context.Context, login string, limit int64) ([]dto.CinemaView, *dto.HttpErr) {
	if !s.userRepo.Exists(ctx, login) {
		return []dto.CinemaView{}, wrapper.NotFoundErr(dto.MsgUserNotFound)
	}

	favorites, err := s.userRepo.GetFavorites(ctx, login, limit)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to save film rate with ERR: %s", err.Error()))
		return []dto.CinemaView{}, wrapper.InternalServerErr(err.Error())
	}
	return favorites, nil
}

func (s *userService) MarkFilmsAsSeen(ctx context.Context, userId uuid.UUID, ids []uuid.UUID) *dto.HttpErr {
	err := s.userRepo.MarkFilmsAsSeen(ctx, userId, ids)
	if err != nil {
		return wrapper.InternalServerErr(err.Error())
	}

	s.metrics.TotalSeen.Add(float64(len(ids)))

	return nil
}

func (s *userService) GetVector(ctx context.Context, login string) (*pgvector.Vector, *dto.HttpErr) {
	vect, err := s.userRepo.GetVectorByLogin(ctx, login)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return vect, nil
}

func (s *userService) SetVector(ctx context.Context, login string, genres []string) *dto.HttpErr {
	if !s.userRepo.Exists(ctx, login) {
		return wrapper.NotFoundErr(dto.MsgUserNotFound)
	}
	u, _ := s.userRepo.GetByLogin(ctx, login)
	newV := utils.VectorFromTags(genres)
	adjV := utils.MultiplyVector(&newV, 5)
	err := s.userRepo.AdjustVector(ctx, u.Id, &adjV)
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to adjust user vector with ERR: %s", err.Error()))
		return wrapper.InternalServerErr(err.Error())
	}
	return nil
}

func (s *userService) SearchUser(ctx context.Context, userLike string) ([]*dto.UserView, *dto.HttpErr) {
	users, err := s.userRepo.SearchUser(ctx, userLike)
	if err != nil {
		return nil, wrapper.InternalServerErr(err.Error())
	}
	return users, nil
}
