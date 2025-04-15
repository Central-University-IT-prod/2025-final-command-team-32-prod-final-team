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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)
var(
	limit = 10
	existCouchId = uuid.New()
	badCouchId = uuid.New()
	filmId = uuid.New()
)


func TestCouchService(t *testing.T){
	ctrl := gomock.NewController(t)
	couchRepo := mock_contracts.NewMockCouchRepository(ctrl)
	
	couchRepo.EXPECT().Exists(ctx, existCouchId).Return(true).AnyTimes()
	couchRepo.EXPECT().Exists(ctx, badCouchId).Return(false).AnyTimes()

	couchRepo.EXPECT().SaveFilmToFavorites(ctx, existCouchId, filmId).Return(nil).AnyTimes()
	couchRepo.EXPECT().SaveFilmToFavorites(ctx, badCouchId, filmId).Return(internalErr).AnyTimes()
	couchRepo.EXPECT().DeleteFilmFromFavorites(ctx, existCouchId, filmId).Return(nil).AnyTimes()
	couchRepo.EXPECT().DeleteFilmFromFavorites(ctx, badCouchId, filmId).Return(internalErr).AnyTimes()
	couchRepo.EXPECT().SaveFilmToBlacklist(ctx, existCouchId, filmId).Return(nil).AnyTimes()
	couchRepo.EXPECT().SaveFilmToBlacklist(ctx, badCouchId, filmId).Return(internalErr).AnyTimes()
	couchRepo.EXPECT().DeleteFilmFromBlacklist(ctx, existCouchId, filmId).Return(nil).AnyTimes()
	couchRepo.EXPECT().DeleteFilmFromBlacklist(ctx, badCouchId, filmId).Return(internalErr).AnyTimes()




	couchService := service.NewCouchService(Metrics, couchRepo)
	t.Run("Get-Recommended", getRecommended(ctx, couchService))
	t.Run("Save-Favorite", saveFav(ctx, couchService))
	t.Run("Delete-Favorite", delFav(ctx, couchService))
	t.Run("Save-Blakclist", saveBlack(ctx, couchService))
	t.Run("Delete-Blacklist", delBlack(ctx, couchService))
}


func getRecommended(ctx context.Context, service contracts.CouchService) func(t *testing.T){
	return func(t *testing.T){
		_, err := service.GetRecommended(ctx, badCouchId, limit)
		require.Equal(t, wrapper.NotFoundErr(dto.MsgCouchNotFound), err)
	}
}


func saveFav(ctx context.Context, service contracts.CouchService) func(t *testing.T){
	return func(t *testing.T){
	err := service.SaveFilmToFavorites(ctx, existCouchId, filmId)
	require.Equal(t, (*dto.HttpErr)(nil), err)

	err = service.SaveFilmToFavorites(ctx, badCouchId, filmId)
	require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}


func delFav(ctx context.Context, service contracts.CouchService) func(t *testing.T){
	return func(t *testing.T){
	err := service.DeleteFilmFromFavorites(ctx, existCouchId, filmId)
	require.Equal(t, (*dto.HttpErr)(nil), err)

	err = service.DeleteFilmFromFavorites(ctx, badCouchId, filmId)
	require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}

func saveBlack(ctx context.Context, service contracts.CouchService) func(t *testing.T){
	return func(t *testing.T){
	err := service.SaveFilmToBlacklist(ctx, existCouchId, filmId)
	require.Equal(t, (*dto.HttpErr)(nil), err)

	err = service.SaveFilmToBlacklist(ctx, badCouchId, filmId)
	require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}
func delBlack(ctx context.Context, service contracts.CouchService) func(t *testing.T){
	return func(t *testing.T){
	err := service.DeleteFilmFromBlacklist(ctx, existCouchId, filmId)
	require.Equal(t, (*dto.HttpErr)(nil), err)

	err = service.DeleteFilmFromBlacklist(ctx, badCouchId, filmId)
	require.Equal(t, wrapper.InternalServerErr(internalErr.Error()), err)
	}
}
