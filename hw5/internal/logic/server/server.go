package server

import (
	"time"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewAPIServer(api *api.API) (e *echo.Echo) {
	e = echo.New()
	// восстанавливается после паники и передача управления HTTPErrorHandler.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(5) * time.Second,
	}))

	// endpoints
	e.POST("/api/users/create", api.UserCreate)
	e.GET("/api/users/:userId", api.UserRead)
	e.PUT("/api/users/:userId", api.UserUpdate)
	e.DELETE("/api/users/:userId", api.UserDelete)
	e.GET("/", func(c echo.Context) error {
		return nil
	})

	return
}
