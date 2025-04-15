package v1

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"testing"
)

func TestCinemaHandler_CreateCinema(t *testing.T) {
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
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	cinemareq := fiber.Map{
		"name":        gofakeit.Name(),
		"poster_url":  gofakeit.URL(),
		"genres":      []string{"Action"},
		"description": gofakeit.Sentence(10),
	}
	e.POST("/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)
}

func TestCinemaHandler_CreateCinema_Admin(t *testing.T) {
	url := "http://prod-team-32-n26k57br.REDACTED:8080/api/v1/"
	e := httpexpect.Default(t, url)

	username := "admin"
	password := "adminadmin"
	req := fiber.Map{
		"username": username,
		"password": password,
	}
	res := e.POST("/users/sign-in").WithJSON(req).
		Expect().
		Status(http.StatusOK).JSON()
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	cinemareq := fiber.Map{
		"name":        gofakeit.Name(),
		"poster_url":  gofakeit.URL(),
		"genres":      []string{"Action"},
		"description": gofakeit.Sentence(10),
	}
	e.POST("/admin/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)
}

func TestCinemaHandler_CreateCinema_Invalid(t *testing.T) {
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
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	cinemareq := fiber.Map{
		"name":        "",
		"poster_url":  gofakeit.URL(),
		"genres":      []string{"Action"},
		"description": gofakeit.Sentence(10),
	}
	e.POST("/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusBadRequest)
}

func TestCinemaHandler_GetById(t *testing.T) {
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
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	name := gofakeit.Name()
	poster := gofakeit.URL()
	genres := []string{"Action"}
	description := gofakeit.Sentence(10)
	cinemareq := fiber.Map{
		"name":        name,
		"poster_url":  poster,
		"genres":      genres,
		"description": description,
	}
	resf := e.POST("/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)
	id := resf.JSON().Object().Value("id").String().Raw()

	e.GET(fmt.Sprintf("/films/%s", id)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusOK).JSON().Object().
		HasValue("name", name).
		HasValue("poster_url", poster).
		HasValue("genres", genres).
		HasValue("description", description)
}

func TestCinemaHandler_GetById_Invalid(t *testing.T) {
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
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	name := gofakeit.Name()
	poster := gofakeit.URL()
	genres := []string{"Action"}
	description := gofakeit.Sentence(10)
	cinemareq := fiber.Map{
		"name":        name,
		"poster_url":  poster,
		"genres":      genres,
		"description": description,
	}
	resf := e.POST("/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)
	resf.JSON().Object().Value("id").String().Raw()
	id := "shitid"
	e.GET(fmt.Sprintf("/films/%s", id)).
		Expect().Status(http.StatusBadRequest)
}

func TestCinemaHandler_GetTopRated(t *testing.T) {
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
	t.Logf("%v", res.Object().Raw())
	token := res.Object().Value("token").String().Raw()
	name := gofakeit.Name()
	poster := gofakeit.URL()
	genres := []string{"Action"}
	description := gofakeit.Sentence(10)
	cinemareq := fiber.Map{
		"name":        name,
		"poster_url":  poster,
		"genres":      genres,
		"description": description,
		"rating":      5.5,
	}
	cinemareq1 := fiber.Map{
		"name":        gofakeit.Name(),
		"poster_url":  poster,
		"genres":      genres,
		"description": description,
		"rating":      7.5,
	}
	e.POST("/films").WithJSON(cinemareq).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)
	e.POST("/films").WithJSON(cinemareq1).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated)

	e.GET("/films/popular").Expect().Status(http.StatusOK).
		JSON().Array()
}
