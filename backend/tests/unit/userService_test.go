package unit

import (
	"context"
	"fmt"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/service"
	"solution/internal/wrapper"
	"solution/pkg/logger"
	"solution/pkg/metric"
	"solution/pkg/utils"
	mock_contracts "solution/tests/unit/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/require"
)

var (
	internalErr      = fmt.Errorf("internal Error")
	_                = internalErr
	ctx              = logger.CtxWithLogger(context.Background())
	goodUsername     = "good-user"
	alreadyUsername  = "already-user"
	already1Username = "already1"
	badUsername      = "bad-user"
	token            = "token"
	badToken         = "bad-token"
	uid              = uuid.New()
	badid            = uuid.New()
	password         = "123"
	password1        = "1234"
	passwordHash     = utils.HashPassword(password)
	goodUser         = dto.UserAuth{
		Login:            goodUsername,
		PasswordUnhashed: password,
	}
	badUser = dto.UserAuth{
		Login:            badUsername,
		PasswordUnhashed: password,
	}
	existsUser = dto.UserAuth{
		Login:            alreadyUsername,
		PasswordUnhashed: password,
	}
	exists1User = dto.UserAuth{
		Login:            already1Username,
		PasswordUnhashed: password1,
	}

	dbUser = dto.UserView{
		Login:          alreadyUsername,
		PasswordHashed: passwordHash,
	}

	Metrics = &metric.PromMetrics{
		TotalRegistered:    promauto.NewCounter(prometheus.CounterOpts{Name: "users_total_registered"}),
		TotalSeen:          promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_seen"}),
		TotalLikedFilms:    promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_liked"}),
		TotalDislikedFilms: promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_disliked"}),
		TotalCouches:       promauto.NewCounter(prometheus.CounterOpts{Name: "couches_total_created"}),
	}
)

func TestUserService(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockUserRepo := mock_contracts.NewMockUserRepository(ctrl)
	mockAuthService := mock_contracts.NewMockAuthService(ctrl)

	mockUserRepo.EXPECT().Exists(ctx, alreadyUsername).Return(true).AnyTimes()
	mockUserRepo.EXPECT().Exists(ctx, already1Username).Return(true).AnyTimes()
	mockUserRepo.EXPECT().Exists(ctx, goodUsername).Return(false).AnyTimes()
	mockUserRepo.EXPECT().Exists(ctx, badUsername).Return(false).AnyTimes()

	mockUserRepo.EXPECT().Create(ctx, &badUser).Return(uuid.UUID{}, internalErr).AnyTimes()
	mockUserRepo.EXPECT().Create(ctx, &goodUser).Return(uid, nil).AnyTimes()

	mockUserRepo.EXPECT().GetByLogin(ctx, alreadyUsername).Return(&dbUser, nil).AnyTimes()
	mockUserRepo.EXPECT().GetByLogin(ctx, already1Username).Return(&dbUser, nil).AnyTimes()
	mockUserRepo.EXPECT().GetByLogin(ctx, badUsername).Return(nil, internalErr).AnyTimes()

	mockUserRepo.EXPECT().SaveFilmToFavorites(ctx, alreadyUsername, uid).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().SaveFilmToFavorites(ctx, already1Username, uid).Return(internalErr).AnyTimes()

	mockUserRepo.EXPECT().DeleteFilmFromFavorites(ctx, alreadyUsername, uid).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().DeleteFilmFromFavorites(ctx, already1Username, uid).Return(internalErr).AnyTimes()

	mockUserRepo.EXPECT().SaveFilmToBlacklist(ctx, alreadyUsername, uid).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().SaveFilmToBlacklist(ctx, already1Username, uid).Return(internalErr).AnyTimes()

	mockUserRepo.EXPECT().DeleteFilmFromBlacklist(ctx, alreadyUsername, uid).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().DeleteFilmFromBlacklist(ctx, already1Username, uid).Return(internalErr).AnyTimes()

	mockAuthService.EXPECT().GenerateToken(ctx, alreadyUsername).Return(token, (*dto.HttpErr)(nil)).AnyTimes()
	mockAuthService.EXPECT().GenerateToken(ctx, goodUsername).Return(token, (*dto.HttpErr)(nil)).AnyTimes()
	mockAuthService.EXPECT().GenerateToken(ctx, already1Username).Return(token, (*dto.HttpErr)(nil)).AnyTimes()

	mockUserRepo.EXPECT().FilmExists(ctx, uid).Return(true).AnyTimes()
	mockUserRepo.EXPECT().FilmExists(ctx, badid).Return(false).AnyTimes()

	mockUserRepo.EXPECT().SearchUser(ctx, goodUsername).Return(nil, nil).AnyTimes()
	mockUserRepo.EXPECT().SearchUser(ctx, badUsername).Return(nil, internalErr).AnyTimes()

	userService := service.NewUserService(Metrics, mockUserRepo, mockAuthService)
	t.Run("Register", testRegister(ctx, userService))
	t.Run("Login", testLogin(ctx, userService))
	t.Run("Save-To-Favorites", testSaveLike(ctx, userService))
	t.Run("Delete-From-Favorites", testDeleteLike(ctx, userService))
	t.Run("Save-To-Blacklist", testSaveBlacklist(ctx, userService))
	t.Run("Delete-From-Blacklist", testDeleteBlacklist(ctx, userService))
	t.Run("Search-User", testSearchUser(ctx, userService))

}

func testRegister(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := service.Register(ctx, &existsUser)
		require.Equal(t, wrapper.ConflictErr(dto.MsgUserAlreadyExists), err)

		_, err = service.Register(ctx, &badUser)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)

		resp, err := service.Register(ctx, &goodUser)
		_ = resp
		require.Equal(t, (*dto.HttpErr)(nil), err)
		require.Equal(t, token, resp.Token)

	}
}

func testLogin(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		resp, err := service.Login(ctx, &existsUser)
		require.Equal(t, (*dto.HttpErr)(nil), err)
		require.Equal(t, token, resp.Token)

		_, err = service.Login(ctx, &badUser)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		_, err = service.Login(ctx, &exists1User)
		require.Equal(t, wrapper.BadRequestErr(dto.MsgInvalidPassword), err)
	}
}

func testSaveLike(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		// user not found
		err := service.SaveFilmToFavorites(ctx, goodUsername, uid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		// film not found
		err = service.SaveFilmToFavorites(ctx, alreadyUsername, badid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgFilmNotFound), err)

		// film liked
		err = service.SaveFilmToFavorites(ctx, alreadyUsername, uid)
		require.Equal(t, (*dto.HttpErr)(nil), err)
		
		err = service.SaveFilmToFavorites(ctx, already1Username, uid)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}

func testDeleteLike(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		// user not found
		err := service.DeleteFilmFromFavorites(ctx, goodUsername, uid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		// film not found
		err = service.DeleteFilmFromFavorites(ctx, alreadyUsername, badid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgFilmNotFound), err)

		// film disliked
		err = service.DeleteFilmFromFavorites(ctx, alreadyUsername, uid)
		require.Equal(t, (*dto.HttpErr)(nil), err)

		err = service.DeleteFilmFromFavorites(ctx, already1Username, uid)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}


func testSaveBlacklist(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T){
		// user not found
		err := service.SaveFilmToBlacklist(ctx, goodUsername, uid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		// film not found
		err = service.SaveFilmToBlacklist(ctx, alreadyUsername, badid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgFilmNotFound), err)

		// film blacklisted
		err = service.SaveFilmToBlacklist(ctx, alreadyUsername, uid)
		require.Equal(t, (*dto.HttpErr)(nil), err)

		err = service.SaveFilmToBlacklist(ctx, already1Username, uid)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)

	}
}
func testDeleteBlacklist(ctx context.Context, service contracts.UserService) func(t *testing.T) {
	return func(t *testing.T){
		// user not found
		err := service.DeleteFilmFromFavorites(ctx, goodUsername, uid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		// film not found
		err = service.DeleteFilmFromFavorites(ctx, alreadyUsername, badid)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgFilmNotFound), err)

		// film unblacklisted
		err = service.DeleteFilmFromFavorites(ctx, alreadyUsername, uid)
		require.Equal(t, (*dto.HttpErr)(nil), err)

		err = service.DeleteFilmFromFavorites(ctx, already1Username, uid)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)

	}
}

func testSearchUser(ctx context.Context, service contracts.UserService) func(t *testing.T){
	return func(t *testing.T){
		_, err := service.SearchUser(ctx, goodUsername)
		require.Equal(t, (*dto.HttpErr)(nil), err)
		
		_, err = service.SearchUser(ctx, badUsername)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}


func saveRate(ctx context.Context, service contracts.UserService) func(t *testing.T){
	return func(t *testing.T){
		// user not found
		err := service.SaveRate(ctx, goodUsername, uid, 10)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgUserNotFound), err)

		// film not found
		err = service.SaveRate(ctx, alreadyUsername, badid, 10)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgFilmNotFound), err)

	}

}
