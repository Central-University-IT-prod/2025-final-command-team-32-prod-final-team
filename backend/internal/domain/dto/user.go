package dto

import (
	"github.com/google/uuid"
)

type UserAuth struct {
	Login            string `json:"username" validate:"required,username" example:"msa"`
	PasswordUnhashed string `json:"password" validate:"required,password" example:"qwerty123_AOISROKT(:#*L(*))"`
	PasswordHashed   string `json:"-"`
	Provider         string `json:"-"`
	YandexId         string `json:"-"`
	AccessToken      string `json:"-"`
}

type UserView struct {
	Id             uuid.UUID `json:"id" example:"some-uuid-v4"`
	Login          string    `json:"username" example:"vanya228"`
	PasswordHashed string    `json:"password" example:"SOME_HASHED_PASSWD"`
}

type UserAuthResponse struct {
	Token string `json:"token" validate:"required"`
}

type Rate struct {
	Rate int32 `json:"rate" validate:"required"`
}

type CreateView struct {
	UserId   uuid.UUID `json:"user_id" validate:"required"`
	CinemaId uuid.UUID `json:"cinema_id" validate:"required"`
}

type AddGenresRequest struct {
	Genres []string `json:"geners" example:"['Action', 'Comedy']"`
}
