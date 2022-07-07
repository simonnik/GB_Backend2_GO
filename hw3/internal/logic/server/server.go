package server

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/api"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/config"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/logic/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewAPIServer(cfg *config.Config, api *api.API) (e *echo.Echo) {
	e = echo.New()
	fs := os.DirFS("web/template")
	p, err := template.ParseFS(fs, "*.html", "*/*.html")
	t := &Template{
		templates: template.Must(p, err),
	}
	e.Renderer = t
	// восстанавливается после паники и передача управления HTTPErrorHandler.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	initLogLevel(e, cfg)

	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(cfg.ServerTimeout) * time.Second,
	}))
	authMiddlewares := middlewares.JWTAuthMiddleware(cfg)

	// endpoints
	e.POST("/api/create", api.Create, authMiddlewares...)
	e.GET("/:token", api.Redirect).Name = "redirect"
	e.GET("/html/form", api.HTML)
	e.GET("/html/stat/:token", api.Stat).Name = "stat"
	e.GET("/", func(c echo.Context) error {
		return nil
	})

	e.GET("/check/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/check/readiness", api.Readiness)

	return
}

func initLogLevel(e *echo.Echo, cfg *config.Config) {
	var loglevelMap = map[string]echoLog.Lvl{
		"debug": echoLog.DEBUG,
		"info":  echoLog.INFO,
		"error": echoLog.ERROR,
		"warn":  echoLog.WARN,
		"off":   echoLog.OFF,
	}

	logLevel, ok := loglevelMap[cfg.Log.Level]
	if !ok {
		e.Logger.Errorf("Undefined log level %#v", cfg.Log.Level)
		logLevel = echoLog.ERROR
	}
	e.Logger.SetLevel(logLevel)
	e.Logger.Infof("%#v", cfg)
}
