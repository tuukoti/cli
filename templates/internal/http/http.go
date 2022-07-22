package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type HTTP struct {
	log *logrus.Logger
}

func RegisterRoutes(e *echo.Echo, log *logrus.Logger) {
	h := HTTP{log: log}

	e.HTTPErrorHandler = h.DefaultErrorHandler

	e.GET("/", h.DefaultHandler)
}
