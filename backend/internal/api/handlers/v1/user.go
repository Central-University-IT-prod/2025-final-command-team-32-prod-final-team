package v1

import (
	"net/http"
	"solution/internal/api/middleware"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userHandler struct {
	userService      contracts.UserService
	authService      contracts.AuthService
	validatorService contracts.ValidatorService
}

func NewUserHandler(us contracts.UserService, as contracts.AuthService, vs contracts.ValidatorService) *userHandler {
	return &userHandler{
		userService:      us,
		authService:      as,
		validatorService: vs,
	}
}

func (uh *userHandler) Setup(r fiber.Router, secretKey string) {
	u := r.Group("/users")
	u.Post("/sign-up", uh.Register)
	u.Post("/sign-in", uh.Login)
	u.Post("/genres", uh.AddGenres)
	u.Get("/search", uh.SearchUser)

	films := r.Group("/films")
	films.Post("/:id/like", middleware.Auth(secretKey), uh.SaveFilmToFavorites)
	films.Delete("/:id/like", middleware.Auth(secretKey), uh.DeleteFilmFromFavorites)

	films.Post("/:id/dislike", middleware.Auth(secretKey), uh.SaveFilmToBlacklist)
	films.Delete("/:id/dislike", middleware.Auth(secretKey), uh.DeleteFilmFromBlacklist)

	films.Post("/:id/rate", middleware.Auth(secretKey), uh.SaveRate)

	r.Get("/plans", middleware.Auth(secretKey), uh.GetFavorites)

	films.Post("/views/bulk", middleware.Auth(secretKey), uh.BulkViews)

}

// RegisterUser godoc
//
//	@Tags		users
//	@Summary	register new user
//	@Param		RequestBody	body	dto.UserAuth	true	"Registers new user and returns access token"
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	dto.UserAuthResponse
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/users/sign-up [post]
func (uh *userHandler) Register(c *fiber.Ctx) error {
	userRegister := new(dto.UserAuth)
	if err := c.BodyParser(userRegister); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	if err := uh.validatorService.ValidateRequestData(userRegister); err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	resp, err := uh.userService.Register(c.UserContext(), userRegister)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(201).JSON(resp)
}

// LoginUser godoc
//
//	@Tags		users
//	@Summary	login existed user
//	@Param		RequestBody	body	dto.UserAuth	true	"Logins existed user and returns access token"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	dto.UserAuthResponse
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	404	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/users/sign-in [post]
func (uh *userHandler) Login(c *fiber.Ctx) error {
	userLogin := new(dto.UserAuth)
	if err := c.BodyParser(userLogin); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	if err := uh.validatorService.ValidateRequestData(userLogin); err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	resp, err := uh.userService.Login(c.UserContext(), userLogin)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(resp)

}

// SaveFilmToFavorites godoc
//
//	@Tags		users
//	@Summary	like film
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string	true	"ID of film"
//	@Success	201				{object}	dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/films/{FilmID}/like [post]
func (uh *userHandler) SaveFilmToFavorites(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err = uh.userService.SaveFilmToFavorites(c.UserContext(), email, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// DeleteFilmFromFavorites godoc
//
//	@Tags		users
//	@Summary	delete like of the film
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string	true	"ID of film"
//	@Success	204				{object}	dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/films/{id}/like [delete]
func (uh *userHandler) DeleteFilmFromFavorites(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err = uh.userService.DeleteFilmFromFavorites(c.UserContext(), email, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

// SaveFilmToBlacklist godoc
//
//	@Tags		users
//	@Summary	dislike film
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string	true	"ID of film"
//	@Success	201				{object}	dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/films/{id}/dislike [post]
func (uh *userHandler) SaveFilmToBlacklist(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err = uh.userService.SaveFilmToBlacklist(c.UserContext(), email, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// DeleteFilmFromBlacklist godoc
//
//	@Tags		users
//	@Summary	delete dislike of the film
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string	true	"ID of film"
//	@Success	204				{object}	dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/films/{id}/dislike [delete]
func (uh *userHandler) DeleteFilmFromBlacklist(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	err = uh.userService.DeleteFilmFromBlacklist(c.UserContext(), email, filmId)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

// SaveRate godoc
//
//	@Tags		users
//	@Summary	rate the film
//	@Security	Bearer
//	@Param		Authorization	header		string		true	"access token 'Bearer {token}'"
//	@Param		FilmID			path		string		true	"ID of film"
//	@Param		RequestBody		body		dto.Rate	true	"Request body"
//	@Success	201				{object}	dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/films/{id}/rate [post]
func (uh *userHandler) SaveRate(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	filmId, goerr := uuid.Parse(c.Params("id"))
	if goerr != nil {
		httpErr := wrapper.BadRequestErr(goerr.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}

	rate := new(dto.Rate)

	if err := c.BodyParser(rate); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	if err := uh.validatorService.ValidateRequestData(rate); err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	err = uh.userService.SaveRate(c.UserContext(), email, filmId, rate.Rate)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// GetFavorites godoc
//
//	@Tags		users
//	@Summary	get liked films
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Success	200				{object}	[]dto.CinemaView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/plans [get]
func (uh *userHandler) GetFavorites(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	email, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}

	var limit int

	limitStr := c.Query("limit")
	if limitStr == "" {
		limit = 30
	} else {
		var goerr error
		limit, goerr = strconv.Atoi(limitStr)
		if goerr != nil {
			httpErr := wrapper.BadRequestErr(goerr.Error())
			return c.Status(httpErr.HttpCode).JSON(httpErr)
		}
	}

	favorites, err := uh.userService.GetFavorites(c.UserContext(), email, int64(limit))
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.JSON(favorites)
}

// BulkViews godoc
//
//	@Tags		users
//	@Summary	bulk views
//	@Security	Bearer
//	@Param		Authorization	header	string		true	"access token 'Bearer {token}'"
//	@Param		RequestBody		body	[]uuid.UUID	true	"Request body"
//	@Success	201
//	@Failure	400	{object}	dto.HttpErr
//	@Failure	401	{object}	dto.HttpErr
//	@Failure	500	{object}	dto.HttpErr
//	@Router		/films/views/bulk [post]
func (uh *userHandler) BulkViews(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	usr, err := uh.userService.GetProfile(c.UserContext(), username)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	cinemas := make([]uuid.UUID, 0)
	if err := c.BodyParser(&cinemas); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	err = uh.userService.MarkFilmsAsSeen(c.UserContext(), usr.Id, cinemas)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.SendStatus(http.StatusCreated)
}

// AddGenres godoc
//
//	@Tags		users
//	@Summary	add genres
//	@Security	Bearer
//	@Param		Authorization	header		string					true	"access token 'Bearer {token}'"
//	@Param		RequestBody		body		dto.AddGenresRequest	true	"Request body"
//	@Success	200				{object}	dto.OkResponse
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/users/genres [post]
func (uh *userHandler) AddGenres(c *fiber.Ctx) error {
	token := c.Locals("subject").(*jwt.Token)
	username, err := uh.authService.GetSubject(c.UserContext(), token)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	req := dto.AddGenresRequest{}
	if err := c.BodyParser(&req); err != nil {
		httpErr := wrapper.BadRequestErr(err.Error())
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	httpErr := uh.userService.SetVector(c.UserContext(), username, req.Genres)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	return c.Status(200).JSON(wrapper.OkResponse())
}

// SearchUser godoc
//
//	@Tags		users
//	@Summary	search user
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Param		query			query		string	true "Search query"
//	@Success	200				{object}	[]dto.UserView
//	@Failure	400				{object}	dto.HttpErr
//	@Failure	401				{object}	dto.HttpErr
//	@Failure	500				{object}	dto.HttpErr
//	@Router		/users/search [post]
func (uh *userHandler) SearchUser(c *fiber.Ctx) error {
	userLike := c.Query("query")
	users, err := uh.userService.SearchUser(c.UserContext(), userLike)
	if err != nil {
		return c.Status(err.HttpCode).JSON(err)
	}
	return c.Status(200).JSON(users)
}
