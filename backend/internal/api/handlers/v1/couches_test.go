package v1

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"testing"
)

func TestCouchHandler_CreateCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}
	e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name)
}

func TestCouchHandler_CreateCouch_Invalid(t *testing.T) {
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

	// name := gofakeit.Name()
	req = fiber.Map{
		"name": "",
		"users": []string{
			"myfriendtony",
		},
	}
	e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusBadRequest)
}

func TestCouchHandler_UpdateCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}
	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()

	req = fiber.Map{
		"name": "new name" + name,
	}
	e.PUT(fmt.Sprintf("/couches/%s", id)).WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusOK).JSON().Object()
}

func TestCouchHandler_UpdateCouch_Invalid(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}
	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()

	req = fiber.Map{
		"name": "",
	}
	e.PUT(fmt.Sprintf("/couches/%s", id)).WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusBadRequest)
}

func TestCouchHandler_GetCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}
	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()

	e.GET(fmt.Sprintf("/couches/%s", id)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusOK).JSON().Object().HasValue("name", name)
}

func TestCouchHandler_GetAllCouches(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}
	e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()

	e.GET("/couches").WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(1)
}

func TestCouchHandler_SaveFilmToFavoritesCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}

	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()
	filmID := "925246d5-cfa5-4500-b94b-54b28b95c482"
	path := fmt.Sprintf("/couches/%s/films/%s/like", id, filmID)
	e.POST(path).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
}

func TestCouchHandler_DeleteFilmFromFavoritesCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}

	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()
	filmID := "919fd8fd-6acb-43ed-8920-4ed90528c4a4"
	e.POST(fmt.Sprintf("/couches/%s/films/%s/like", id, filmID)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
	e.DELETE(fmt.Sprintf("/couches/%s/films/%s/like", id, filmID)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusNoContent)
}

func TestCouchHandler_SaveFilmToBlacklistCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}

	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()
	filmID := "919fd8fd-6acb-43ed-8920-4ed90528c4a4"
	e.POST(fmt.Sprintf("/couches/%s/films/%s/dislike", id, filmID)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
}

func TestCouchHandler_DeleteFilmFromBlacklistCouch(t *testing.T) {
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

	name := gofakeit.Name()
	req = fiber.Map{
		"name": name,
		"users": []string{
			"myfriendtony",
		},
	}

	id := e.POST("/couches").WithJSON(req).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect().Status(http.StatusCreated).JSON().Object().HasValue("name", name).Value("id").String().Raw()
	filmID := "925246d5-cfa5-4500-b94b-54b28b95c482"
	e.POST(fmt.Sprintf("/couches/%s/films/%s/dislike", id, filmID)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusCreated)
	e.DELETE(fmt.Sprintf("/couches/%s/films/%s/dislike", id, filmID)).WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)).Expect().
		Status(http.StatusNoContent)
}
