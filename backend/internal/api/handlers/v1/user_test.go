package v1

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)
	// Выполняем запрос и проверяем результат
	req := fiber.Map{
		"username": gofakeit.Name(),
		"password": gofakeit.Password(true, true, true, true, true, 8),
	}
	e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON()
}

func TestUserHandler_Register_Invalid(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)
	// Выполняем запрос и проверяем результат
	req := fiber.Map{
		"username": "",
		"password": gofakeit.Password(true, true, true, true, true, 8),
	}
	e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusBadRequest)
}

func TestUserHandler_Register_Invalid1(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)
	// Выполняем запрос и проверяем результат
	req := fiber.Map{
		"username": gofakeit.Name(),
		"password": "",
	}
	e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusBadRequest)
}

func TestUserHandler_Login(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)
	// Выполняем запрос и проверяем результат
	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated)

	e.POST("/users/sign-in").WithJSON(req).
		Expect().
		Status(http.StatusOK)
}

func TestUserHandler_Login_Invalid(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)
	// Выполняем запрос и проверяем результат
	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated)

	req["password"] = ""
	e.POST("/users/sign-in").WithJSON(req).
		Expect().
		Status(http.StatusBadRequest)
}

func TestUserHandler_GetInfo(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	e.GET("/users/me").
		WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("username", username)
}

func TestUserHandler_SaveFilmToFavorites(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	resl := e.POST(fmt.Sprintf("/films/%s/like", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
	t.Log(resl.Body().Raw())
}

func TestUserHandler_DeleteFilmFromFavorites(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	e.POST(fmt.Sprintf("/films/%s/like", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)

	e.DELETE(fmt.Sprintf("/films/%s/like", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusNoContent)
}

func TestUserHandler_SaveFilmToBlacklist(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	e.POST(fmt.Sprintf("/films/%s/dislike", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
}

func TestUserHandler_DeleteFilmFromBlacklist(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	e.POST(fmt.Sprintf("/films/%s/dislike", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)

	e.DELETE(fmt.Sprintf("/films/%s/dislike", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusNoContent)
}

func TestUserHandler_GetFavorites(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	e.POST(fmt.Sprintf("/films/%s/like", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)

	e.GET("/plans").WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusOK).JSON().Array().Length().IsEqual(1)
}

func TestUserHandler_GetFavorites1(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 8)
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-up").WithJSON(req).
		Expect().
		Status(http.StatusCreated).JSON()
	token := res.Object().Value("token").String().Raw()

	filmId, _ := uuid.Parse("ff1c1749-5e8c-46fc-84d9-5f28a73c565b")

	e.POST(fmt.Sprintf("/films/%s/like", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)

	e.POST(fmt.Sprintf("/films/%s/dislike", filmId.String())).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)

	e.GET("/plans").WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusOK).JSON().Array().Length().IsEqual(0)
}
