package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/activities"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/users"
)

type API struct {
	Users      *users.Users
	Activities *activities.Activities
}

type UserFilter struct {
	UserId *string `param:"userId"`
}

func NewAPI(users *users.Users, acts *activities.Activities) *API {
	return &API{users, acts}
}

// UserCreate http handler
func (api *API) UserCreate(c echo.Context) error {
	c.Logger().Info("User create")
	newModel := api.Users.NewModel()
	if err := c.Bind(newModel); err != nil {
		return err
	}

	if err := api.Users.Create(c.Request().Context(), newModel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a := &activities.Activity{
		newModel.UserId,
		"2020-12-12",
		"Registered",
	}
	if err := api.Activities.Create(c.Request().Context(), a); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, newModel)
}

// UserRead http handler
func (api *API) UserRead(c echo.Context) error {
	c.Logger().Info("User read")
	u := &users.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	err := api.Users.Read(c.Request().Context(), u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

// UserUpdate http handler
func (api *API) UserUpdate(c echo.Context) error {
	c.Logger().Info("User update")
	u := &users.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	err := api.Users.Update(c.Request().Context(), u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

// UserDelete http handler
func (api *API) UserDelete(c echo.Context) error {
	c.Logger().Info("User delete")
	u := &users.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	err := api.Users.Delete(c.Request().Context(), u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}
