package v1

import (
	"net/http"
	"solution/internal/api/middleware"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type adminHandler struct {
	cinemaService    contracts.CinemaService
	authService      contracts.AuthService
	validatorService contracts.ValidatorService
}

func NewAdminHandler(cs contracts.CinemaService, vs contracts.ValidatorService, authService contracts.AuthService) *adminHandler {
	return &adminHandler{
		cinemaService:    cs,
		validatorService: vs,
		authService:      authService,
	}
}

func (ch *adminHandler) Setup(r fiber.Router, secretKey string) {
	f := r.Group("/admin") // for rebuild
	f.Get("/", middleware.Auth(secretKey), ch.CheckPrivileges)
	f.Put("/films/:id", middleware.Auth(secretKey), ch.UpdateFilm)
	f.Delete("/films/:id", middleware.Auth(secretKey), ch.DeleteFilm)
	f.Post("/films", middleware.Auth(secretKey), ch.CreateFilm)
}

// UpdateFilm godoc
//
//	@Tags			admin
//	@Summary		update film
//	@Description	Обновляет фильм
//	@Security		Bearer
//	@Param			Authorization	header	string				true	"access token 'Bearer {token}'"
//	@Param			id				path	string				true	"ID of film"
//	@Param			RequestBody		body	dto.CinemaUpdate	true	"Request Body"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	dto.CinemaUpdate
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/admin/films/{id} [put]
func (ch *adminHandler) UpdateFilm(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, httperr := ch.authService.GetSubject(c.UserContext(), token)
	if httperr != nil {
		return c.Status(httperr.HttpCode).JSON(httperr)
	}

	filmsId := c.Params("id")
	fid, err := uuid.Parse(filmsId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	upd := dto.CinemaUpdate{}
	if err := c.BodyParser(&upd); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	upd.ID = fid
	httpErr := ch.cinemaService.UpdateFilm(c.UserContext(), upd, email)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	return c.Status(http.StatusOK).JSON(upd)
}

// DeleteFilm godoc
//
//	@Tags			admin
//	@Summary		delete film
//	@Description	Удаляет фильм
//	@Security		Bearer
//	@Param			Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Param			id				path	string	true	"ID of film"
//	@Produce		json
//	@Success		200	{object}	dto.OkResponse
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/admin/films/{id} [delete]
func (ch *adminHandler) DeleteFilm(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, httperr := ch.authService.GetSubject(c.UserContext(), token)
	if httperr != nil {
		return c.Status(httperr.HttpCode).JSON(httperr)
	}
	filmsId := c.Params("id")
	fid, err := uuid.Parse(filmsId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	httpErr := ch.cinemaService.DeleteFilm(c.UserContext(), fid, username)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(wrapper.OkResponse())
}

// CreateFilm godoc
//
//	@Tags			admin
//	@Summary		create film
//	@Description	Создает фильм
//	@Security		Bearer
//	@Param			Authorization	header	string				true	"access token 'Bearer {token}'"
//	@Param			RequestBody		body	dto.CinemaUpdate	true	"Request Body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.OkResponse
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/admin/films [post]
func (ch *adminHandler) CreateFilm(ctx *fiber.Ctx) error {
	token := ctx.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(ctx.UserContext(), token)
	if err != nil {
		return ctx.Status(err.HttpCode).JSON(err)
	}
	filmCreate := dto.CinemaCreate{}
	if err := ctx.BodyParser(&filmCreate); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return ctx.Status(httpErr.HttpCode).JSON(httpErr)
	}
	filmCreate.Private = false
	ch.cinemaService.Create(ctx.UserContext(), &filmCreate, username)
	return ctx.SendStatus(fiber.StatusCreated)
}

func (ch *adminHandler) CheckPrivileges(ctx *fiber.Ctx) error {
	token := ctx.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(ctx.UserContext(), token)
	if err != nil {
		return ctx.Status(err.HttpCode).JSON(err)
	}
	if !ch.cinemaService.CheckUserIsPrivileged(ctx.UserContext(), username) {
		return ctx.SendStatus(403)
	}
	return ctx.SendStatus(200)
}
