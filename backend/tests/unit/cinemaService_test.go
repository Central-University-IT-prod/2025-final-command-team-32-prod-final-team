package unit

import (
	"context"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/service"
	"solution/internal/wrapper"
	mock_contracts "solution/tests/unit/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)
var(
	cc = dto.CinemaCreate{}
	cu = dto.CinemaUpdate{}
	cn = (*dto.CinemaCreate)(nil)
	admin = "admin"
	user = "user"
	aux = "aux"
)


func TestCinemaService(t *testing.T){
	ctrl := gomock.NewController(t)
	cinemaRepo := mock_contracts.NewMockCinemaReposiroty(ctrl)

	cinemaRepo.EXPECT().CheckUserIsPrivileged(ctx, admin).Return(true).AnyTimes()
	cinemaRepo.EXPECT().CheckUserIsPrivileged(ctx,	user).Return(false).AnyTimes()
	cinemaRepo.EXPECT().CheckUserIsPrivileged(ctx,	aux).Return(true).AnyTimes()
	cinemaRepo.EXPECT().UpdateFilm(ctx, cu, admin).Return(nil).AnyTimes()
	cinemaRepo.EXPECT().UpdateFilm(ctx, cu, aux).Return(internalErr).AnyTimes()
	cinemaRepo.EXPECT().DeleteFilm(ctx, cu, admin).Return(nil).AnyTimes()
	cinemaRepo.EXPECT().DeleteFilm(ctx, cu, aux).Return(internalErr).AnyTimes()

	
	

	cinemaService := service.NewCinemaService(cinemaRepo)
	t.Run("Create-F", createF(ctx, cinemaService))
	t.Run("Update-F", updateF(ctx, cinemaService))
	t.Run("Delete-F", deleteF(ctx, cinemaService))
}



func createF(ctx context.Context, service contracts.CinemaService)func(t *testing.T){
	return func(t *testing.T){
		_, err := service.Create(ctx, cn, admin)
		require.Equal(t, wrapper.BadRequestErr("not enough data"), err)

		_, err = service.Create(ctx, &cc, user)
		require.Equal(t, wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden), err)
	}
}
func updateF(ctx context.Context, service contracts.CinemaService)func(t *testing.T){

	return func(t *testing.T){
		err := service.UpdateFilm(ctx, cu, user)
		require.Equal(t, wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden), err)

		err = service.UpdateFilm(ctx, cu, aux)
		require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
		
		err = service.UpdateFilm(ctx, cu, admin)
		require.Equal(t, (*dto.HttpErr)(nil), err)
		
	}
}
func deleteF(ctx context.Context, service contracts.CinemaService)func(t *testing.T){

	return func(t *testing.T){

		err := service.DeleteFilm(ctx, filmId, user)
		require.Equal(t, wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden), err)

	}
}
