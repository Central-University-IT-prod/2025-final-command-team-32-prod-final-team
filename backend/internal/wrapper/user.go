package wrapper

import (
	"github.com/google/uuid"
	"solution/internal/database/storage"
	"solution/internal/domain/dto"
	"solution/pkg/utils"

	"github.com/pgvector/pgvector-go"
)

func UserAuthWithCreateParams(u *dto.UserAuth, v pgvector.Vector) *storage.CreateUserParams {
	u.PasswordHashed = utils.HashPassword(u.PasswordUnhashed)
	return &storage.CreateUserParams{
		Login:     u.Login,
		Password:  &u.PasswordHashed,
		Embedding: v,
	}
}

func UserWithView(u *storage.User) *dto.UserView {
	var pass string
	if u.Password == nil {
		pass = ""
	} else {
		pass = *u.Password
	}
	return &dto.UserView{
		Id:             u.ID,
		Login:          u.Login,
		PasswordHashed: pass,
	}
}

func ViewForDatabase(userId uuid.UUID, ids []uuid.UUID) []storage.MarkAsViewedBulkParams {
	res := make([]storage.MarkAsViewedBulkParams, 0)
	for _, id := range ids {
		res = append(res, storage.MarkAsViewedBulkParams{
			SubjectID: userId,
			CinemaID:  id,
		})
	}
	return res
}
