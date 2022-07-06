package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/check"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/links"
)

type API struct {
	Host  string
	Links *links.Links
	Check *check.Check
}

type LinkFilter struct {
	Token *string `param:"token"`
}

func NewAPI(h string, links *links.Links, ch *check.Check) *API {
	return &API{h, links, ch}
}

// Create http handler
func (api *API) Create(c echo.Context) error {
	c.Logger().Info("Create")
	newLink := api.Links.NewModel()
	if err := c.Bind(newLink); err != nil {
		return err
	}

	if err := api.Links.Create(c.Request().Context(), *newLink); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	host := strings.TrimRight(api.Host, "/")
	link := host + c.Echo().Reverse("redirect", newLink.Token)
	stat := host + c.Echo().Reverse("stat", newLink.Token)
	return c.JSON(http.StatusOK, map[string]string{"link": link, "stat": stat})
}

// Readiness check application
func (api *API) Readiness(c echo.Context) error {
	err := api.Check.Check()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.
			Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (api *API) Redirect(c echo.Context) error {
	c.Logger().Info("Redirect")
	f := LinkFilter{}
	if err := c.Bind(&f); err != nil {
		return err
	}

	if f.Token == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "link token can't be empty")
	}

	link, err := api.Links.FindByToken(c.Request().Context(), *f.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = api.Links.SaveStat(c.Request().Context(), link.ID, c.RealIP())
	if err != nil {
		c.Echo().Logger.Error(err)
	}

	c.Response().Header().Set("Cache-Control", "no-cache")

	return c.Redirect(http.StatusMovedPermanently, link.Link)
}

func (api *API) HTML(c echo.Context) error {
	c.Logger().Info("HTML")

	token, err := api.Links.GetJWTToken()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "form.html", map[string]interface{}{"jwtToken": token})
}

func (api *API) Stat(c echo.Context) error {
	c.Logger().Info("Stat")
	f := LinkFilter{}
	if err := c.Bind(&f); err != nil {
		return err
	}

	if f.Token == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "link id can't be empty")
	}

	ls, err := api.Links.FindAllByToken(c.Request().Context(), *f.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "stat.html", ls)
}
