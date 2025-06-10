package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]any{
		"Titulo": "PDV Home",
	})
}
