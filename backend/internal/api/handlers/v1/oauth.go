package v1

import (
	"fmt"
	"solution/internal/domain/contracts"
	"solution/internal/wrapper"
	"solution/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type yandexCred struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_i"`
	RefreshToken string `json:"refresh_token"`
}

type yandexUserInfo struct {
	Login    string `json:"login"`
	Id       string `json:"id"`
	ClientId string `json:"client_id"`
	Uid      string `json:"ui"`
	Psuid    string `json:"psuid"`
}

type oauthHandler struct {
	userService contracts.UserService
	authService contracts.AuthService
}

func NewOuathHandler(us contracts.UserService, as contracts.AuthService) *oauthHandler {
	return &oauthHandler{
		userService: us,
		authService: as,
	}
}

func (oh *oauthHandler) Setup(r fiber.Router) {
	r.Get("/users/yaauth/callback", oh.YandexCallback)
}

func (oh *oauthHandler) YandexCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		httpErr := wrapper.BadRequestErr("No auth code provided from yandex api")
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	logger.FromCtx(c.UserContext()).Info(c.UserContext(), fmt.Sprintf("code: %s", code))
	tokenInfo, err := oh.authService.ExchangeCodeForToken(c.UserContext(), code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"TOKENINFO err": err.Error()})
	}
	userInfo, err := oh.authService.ExchangeTokenForUserInfo(c.UserContext(), tokenInfo.AccessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"USERINFO err": err.Error()})
	}
	user := wrapper.NewYandexUser(tokenInfo, userInfo)
	resp, httpErr := oh.userService.Register(c.UserContext(), user)
	if httpErr != nil {
		return c.Status(httpErr.HttpCode).JSON(httpErr)
	}
	return c.Status(200).JSON(resp)
}
