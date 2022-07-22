package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HTTP) DefaultHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (h *HTTP) DefaultErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	h.log.Error(err)

	err = c.Render(code, "error.html", nil)
	if err != nil {
		h.log.Error(err)
	}
}
