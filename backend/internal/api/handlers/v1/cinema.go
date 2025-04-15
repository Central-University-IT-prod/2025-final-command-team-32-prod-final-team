package v1

import (
	"fmt"
	"io"
	"net/http"
	"solution/internal/api/middleware"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/logger"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type cinemaHandler struct {
	admin            string
	cinemaService    contracts.CinemaService
	authService      contracts.AuthService
	validatorService contracts.ValidatorService
	fileService      contracts.FileService
}

func NewCinemaHandler(cs contracts.CinemaService, vs contracts.ValidatorService, authService contracts.AuthService, adminLogin string, fileService contracts.FileService) *cinemaHandler {
	return &cinemaHandler{
		cinemaService:    cs,
		validatorService: vs,
		authService:      authService,
		admin:            adminLogin,
		fileService:      fileService,
	}
}

func (ch *cinemaHandler) Setup(r fiber.Router, secretKey string) {
	f := r.Group("/films")
	f.Post("/", middleware.Auth(secretKey), ch.CreateCinema)
	f.Get("/genres", ch.GetGenres)
	f.Get("/feed", middleware.Auth(secretKey), ch.GetRecommended)
	f.Get("/popular", ch.GetTopRated)
	f.Get("/search", ch.SearchFilm)
	f.Get("/:id", ch.GetById)

	f.Post("/:id/picture", middleware.Auth(secretKey), ch.SetCinemaPic)

}

// CreateCinema godoc
//
//	@Tags			films
//	@Summary		upload new film
//	@Description	Добавляет новый фильм в подборку пользователя или добавляет глобальный фильм в базу админом
//	@Security		Bearer
//	@Param			Authorization	header	string				true	"access token 'Bearer {token}'"
//	@Param			RequestBody		body	dto.CinemaCreate	true	"Request Body"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	dto.CinemaCreateResponse
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/films/ [post]
func (h cinemaHandler) CreateCinema(c *fiber.Ctx) error {
	cinemacreate := new(dto.CinemaCreate)
	if err := c.BodyParser(cinemacreate); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	if err := h.validatorService.ValidateRequestData(cinemacreate); err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	token, ok := c.Locals("subject").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(nil)
	}
	username, err := h.authService.GetSubject(c.Context(), token)

	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	// add local film
	cinemacreate.Private = true
	id, err := h.cinemaService.Create(c.UserContext(), cinemacreate, username)
	if err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	return c.Status(http.StatusCreated).JSON(dto.CinemaCreateResponse{Id: id})
}

// RegisterUser godoc
//
//	@Tags			films
//	@Summary		get recommended films
//	@Description	Список рекомендованных фильмов
//	@Param			limit		query	int		true	"limit"
//	@Param			username	query	string	true	"usernames"
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.CinemaView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/films/feed [get]
func (ch *cinemaHandler) GetRecommended(c *fiber.Ctx) error {
	token, ok := c.Locals("subject").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(nil)
	}
	username, _ := ch.authService.GetSubject(c.Context(), token)
	limit := c.QueryInt("limit", 10)
	cinemas, err := ch.cinemaService.GetRecommended(c.UserContext(), username, limit)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(cinemas)

}

// RegisterUser godoc
//
//	@Tags			films
//	@Summary		get film
//	@Description	Получить фильм по id
//	@Param			filmId	path	string	true	"film uuid"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.CinemaView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/films/{filmId} [get]
func (ch *cinemaHandler) GetById(c *fiber.Ctx) error {
	id, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	cinema, err := ch.cinemaService.GetById(c.UserContext(), id)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	return c.Status(200).JSON(cinema)
}

// RegisterUser godoc
//
//	@Tags			films
//	@Summary		get top rated films
//	@Description	Получить список Самых Высоко оценённых фильмов
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.CinemaView
//	@Failure		400	{object}	dto.HttpErr
//	@Failure		500	{object}	dto.HttpErr
//	@Router			/films/popular [get]
func (ch *cinemaHandler) GetTopRated(c *fiber.Ctx) error {
	logger.FromCtx(c.UserContext()).Info(c.UserContext(), "IN TOP RATED")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset")
	cinemas, err := ch.cinemaService.GetTopRated(c.UserContext(), offset, limit)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(cinemas)
}

// SetCinemaPic godoc
//
//	@Summary	Set picture of the cinema
//	@Tags		films
//	@Accept		multipart/form-data
//	@Produce	json
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string	true	"UUID of the film"
//	@Param		uploadfile		formData	file	true	"File of the pic"
//	@Success	200
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	401	{object}	dto.HttpErr
//	@Failure	404	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/films/{FilmID}/picture [post]
func (ch *cinemaHandler) SetCinemaPic(c *fiber.Ctx) error {
	cinemaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.FromCtx(c.UserContext()).Info(c.UserContext(), fmt.Sprintf("cinemaId: %s", cinemaId))
		httpErr := wrapper.BadRequestErr("bad cinema id")
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	token, ok := c.Locals("subject").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(nil)
	}
	username, httpErr := ch.authService.GetSubject(c.Context(), token)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	if !ch.cinemaService.CheckUserIsPrivileged(c.UserContext(), username) {
		httpErr := wrapper.AccessForbiddenErr(dto.MsgUserAccessForbidden)
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	file, err := c.FormFile("uploadfile")
	if err != nil {
		httpErr := wrapper.BadRequestErr("failed to get file")
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	src, err := file.Open()
	if err != nil {
		httpErr := wrapper.InternalServerErr("failed to produce file")
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		httpErr := wrapper.BadRequestErr("failed to get file")
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	httpErr = ch.fileService.UploadFile(c.UserContext(), cinemaId, file.Filename, fileBytes)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	return c.SendStatus(http.StatusOK)
}

// SearchFilm godoc
//
//	@Tags		cinemas
//	@Summary	search film
//	@Param		query	query		string	true	"Search query"
//	@Success	200		{object}	dto.OkResponse
//	@Failure	400		{object}	dto.HttpErr
//	@Failure	500		{object}	dto.HttpErr
//	@Router		/films/search [post]
func (ch *cinemaHandler) SearchFilm(c *fiber.Ctx) error {
	titleLike := c.Query("query")
	tags := c.Query("tags")
	tagsSlice := strings.Split(tags, ",")
	logger.FromCtx(c.UserContext()).Info(c.UserContext(), fmt.Sprintf("slice : %#v", len(tagsSlice)))
	if tagsSlice[0] == "" {
		tagsSlice = []string{}
	}

	cinemas, err := ch.cinemaService.SearchFilm(c.UserContext(), titleLike, tagsSlice)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(cinemas)
}

// GetGenres godoc
//
//	@Summary	Get genres
//	@Tags		films
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	404	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/films/genres [post]
func (ch *cinemaHandler) GetGenres(c *fiber.Ctx) error {
	genres := ch.cinemaService.GetGenres(c.UserContext())
	return c.Status(200).JSON(genres)
}
