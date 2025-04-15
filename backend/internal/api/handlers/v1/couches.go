package v1

import (
	"errors"
	"net/http"
	"slices"
	"solution/internal/api/middleware"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type couchHandler struct {
	authService      contracts.AuthService
	validatorService contracts.ValidatorService
	couchService     contracts.CouchService
	userService      contracts.UserService
}

func NewCouchHandler(vs contracts.ValidatorService, authService contracts.AuthService, cs contracts.CouchService, us contracts.UserService) *couchHandler {
	return &couchHandler{
		authService:      authService,
		validatorService: vs,
		couchService:     cs,
		userService:      us,
	}
}

func (ch *couchHandler) Setup(r fiber.Router, secretKey string) {
	couch := r.Group("/couches")
	couch.Post("/", middleware.Auth(secretKey), ch.CreateCouch)
	couch.Get("/:id", middleware.Auth(secretKey), ch.GetCouch)
	couch.Put("/:id", middleware.Auth(secretKey), ch.UpdateCouch)
	couch.Get("/", middleware.Auth(secretKey), ch.GetAllCouches)

	couch.Get("/:couch_id/feed", ch.GetRecommended)
	couch.Post("/:couch_id/films/:id/like", middleware.Auth(secretKey), ch.SaveFilmToFavoritesCouch)
	couch.Delete("/:couch_id/films/:id/like", middleware.Auth(secretKey), ch.DeleteFilmFromFavoritesCouch)

	couch.Post("/:couch_id/films/:id/dislike", middleware.Auth(secretKey), ch.SaveFilmToBlacklistCouch)
	couch.Delete("/:couch_id/films/:id/dislike", middleware.Auth(secretKey), ch.DeleteFilmFromBlacklistCouch)

	couch.Get("/:couch_id/plans", middleware.Auth(secretKey), ch.GetFavorites)
	couch.Post("/:couch_id/views/bulk", ch.MarkFilmsAsSeen)
}

// CreateCouch godoc
//
//	@Tags			couches
//	@Summary		create new couch
//	@Description	Создает новый диван
//	@Param			RequestBody	body	dto.CreateCouch	true	"Request Body"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	dto.CouchView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/ [post]
func (ch *couchHandler) CreateCouch(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	usr, err := ch.userService.GetProfile(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	couch := dto.CreateCouch{}
	if err := c.BodyParser(&couch); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	couch.AuthorName = username
	couch.AuthorID = usr.Id
	if !slices.Contains(couch.Sitters, username) {
		couch.Sitters = append(couch.Sitters, username)
	}

	id, err := ch.couchService.Create(c.UserContext(), couch)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.CouchView{
		Id:      id,
		Name:    couch.Name,
		Sitters: couch.Sitters,
		Author:  username,
	})
}

// UpdateCouch godoc
//
//	@Tags			couches
//	@Summary		update couch
//	@Description	Обновляет диван
//	@Param			id			path	string			true	"Couch ID"
//	@Param			RequestBody	body	dto.UpdateCouch	true	"Request Body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.CouchView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id} [put]
func (ch *couchHandler) UpdateCouch(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	_, err = ch.userService.GetProfile(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	req := dto.UpdateCouch{}
	if err := c.BodyParser(&req); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	cid, httperr := uuid.Parse(c.Params("id"))
	if httperr != nil {
		return c.Status(http.StatusBadRequest).JSON(httperr)
	}
	req.AuthorName = username
	if req.Name != nil && *req.Name == "" {
		httpErr := wrapper.BadRequestErr(errors.New("name is required").Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	err = ch.couchService.Update(c.UserContext(), cid, req)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	couch, err := ch.couchService.GetOne(c.UserContext(), cid)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	if couch.Sitters == nil {
		couch.Sitters = []string{}
	}
	return c.Status(fiber.StatusOK).JSON(couch)
}

// GetAllCouches godoc
//
//	@Tags			couches
//	@Summary		get couches
//	@Description	Получает все диваны
//	@Produce		json
//	@Success		200	{object}	[]dto.CouchView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/ [get]
func (ch *couchHandler) GetAllCouches(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	_, err = ch.userService.GetProfile(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	res, err := ch.couchService.GetMany(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// GetCouch godoc
//
//	@Tags			couches
//	@Summary		get couch
//	@Description	Получает диван по его ID
//	@Param			id	path	string	true	"Couch ID"
//	@Produce		json
//	@Success		200	{object}	dto.CouchView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id} [get]
func (ch *couchHandler) GetCouch(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := ch.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	_, err = ch.userService.GetProfile(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	cid, httperr := uuid.Parse(c.Params("id"))
	if httperr != nil {
		return c.Status(http.StatusBadRequest).JSON(httperr)
	}
	couch, err := ch.couchService.GetOne(c.UserContext(), cid)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(couch)
}

// GetFavorites godoc
//
//	@Tags			couches
//	@Summary		get couch likes
//	@Description	Получает любимые фильмы дивана
//	@Param			id	path	string	true	"Couch ID"
//	@Produce		json
//	@Success		200	{object}	dto.CouchView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/plans [get]
func (ch *couchHandler) GetFavorites(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	var limit int

	limitStr := c.Query("limit")
	limit, _ = strconv.Atoi(limitStr)
	if limitStr == "" {
		limit = 30
	}
	// for rerun
	favorites, err := ch.couchService.GetFavorites(c.UserContext(), cid, int64(limit))
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.JSON(favorites)
}

// SaveFilmToFavoritesCouch godoc
//
//	@Tags			couches
//	@Summary		like film in couch
//	@Description	Лайкает фильм в диване
//	@Param			id		path	string	true	"Couch ID"
//	@Param			filmId	path	string	true	"Film ID"
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/films/{filmId}/like [post]
func (ch *couchHandler) SaveFilmToFavoritesCouch(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	err := ch.couchService.SaveFilmToFavorites(c.UserContext(), cid, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// DeleteFilmFromFavoritesCouch godoc
//
//	@Tags			couches
//	@Summary		delete like of the film in couch
//	@Description	Удаляет лайк фильма в диване
//	@Param			id		path	string	true	"Couch ID"
//	@Param			filmId	path	string	true	"Film ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/films/{filmId}/like [delete]
func (ch *couchHandler) DeleteFilmFromFavoritesCouch(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err := ch.couchService.DeleteFilmFromFavorites(c.UserContext(), cid, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

// SaveFilmToBlacklistCouch godoc
//
//	@Tags			couches
//	@Summary		dislike film in couch
//	@Description	Дизлайкает фильм в диване
//	@Param			id		path	string	true	"Couch ID"
//	@Param			filmId	path	string	true	"Film ID"
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/films/{filmId}/dislike [post]
func (ch *couchHandler) SaveFilmToBlacklistCouch(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err := ch.couchService.SaveFilmToBlacklist(c.UserContext(), cid, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// DeleteFilmFromBlacklistCouch godoc
//
//	@Tags			couches
//	@Summary		delete dislike of the film in couch
//	@Description	Удаляет дизлайк фильма в диване
//	@Param			id		path	string	true	"Couch ID"
//	@Param			filmId	path	string	true	"Film ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/films/{filmId}/dislike [delete]
func (ch *couchHandler) DeleteFilmFromBlacklistCouch(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err := ch.couchService.DeleteFilmFromBlacklist(c.UserContext(), cid, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetRecommended godoc
//
//	@Tags			couches
//	@Summary		get recommended films in couch
//	@Description	Показывает рекомендации в диване
//	@Param			id	path	string	true	"Couch ID"
//	@Produce		json
//	@Success		200	{object}	[]dto.CinemaView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		401	{object}	dto.HttpErr
//	@Failure		403	{object}	dto.HttpErr
//	@Failure		404	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/couches/{id}/feed [get]
func (ch *couchHandler) GetRecommended(c *fiber.Ctx) error {
	couchId := c.Params("couch_id")
	limit := c.QueryInt("limit", 10)
	cid, goerr := uuid.Parse(couchId)
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	cinemas, err := ch.couchService.GetRecommended(c.UserContext(), cid, limit)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(cinemas)
}

// MarkFilmAsSeen godoc
//
//	@Tags		couches
//	@Summary	mark film as seen in couch
//	@Param		id			path	string		true	"Couch ID"
//	@Param		RequestBody	body	[]uuid.UUID	true	"Request Body"
//	@Produce	json
//	@Success	201
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	401	{object}	dto.HttpErr
//	@Failure	403	{object}	dto.HttpErr
//	@Failure	404	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/couches/{id}/views/bulk [post]
func (ch *couchHandler) MarkFilmsAsSeen(c *fiber.Ctx) error {
	couchid := c.Params("couch_id")
	cid, goerr := uuid.Parse(couchid)
	if goerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	cinemas := make([]uuid.UUID, 0)
	if err := c.BodyParser(&cinemas); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	err := ch.couchService.MarkFilmsAsSeen(c.UserContext(), cid, cinemas)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}
