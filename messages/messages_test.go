package messages

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
)

func TestMessagesPersistTogetherAndAreConsumedOnce(t *testing.T) {
	session.SetupSession()
	app := fiber.New()
	app.Use(session.Middleware)
	app.Get("/set", func(c fiber.Ctx) error {
		SetErrorMessage(c, "erro")
		SetSuccessMessage(c, "sucesso")
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Get("/read", func(c fiber.Ctx) error { return c.JSON(LoadMessages(c)) })
	res, err := app.Test(httptest.NewRequest("GET", "/set", nil))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	cookies := res.Cookies()
	for _, want := range []Message{{Error: "erro", Success: "sucesso"}, {}} {
		req := httptest.NewRequest("GET", "/read", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		var got Message
		err = json.NewDecoder(res.Body).Decode(&got)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("messages = %+v, want %+v", got, want)
		}
		cookies = res.Cookies()
	}
}
