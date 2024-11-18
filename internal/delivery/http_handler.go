package delivery

import (
	"net/http"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/cuongtranba/mynoti/pkg/middleware"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server *http.Server
}

func NewServer(
	port string,
	comicUseCase usecase.ComicUseCase,
	log *logger.Logger,
) *Server {
	e := echo.New()

	e.Use(echo_middleware.Recover())
	e.Use(echo.WrapMiddleware(middleware.ContextMiddleware()))
	e.Use(echo.WrapMiddleware(middleware.LoggerMiddleware(log)))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/subscribe", func(c echo.Context) error {
		var comic Comic
		if err := c.Bind(&comic); err != nil {
			return err
		}
		err := comicUseCase.Subscribe(c.Request().Context(), &domain.Comic{
			Url:         comic.Url,
			Name:        comic.Name,
			Description: comic.Description,
			Html:        comic.Html,
			CronSpec:    comic.CronSpec,
		})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	})

	return &Server{
		server: &http.Server{
			Addr:    port,
			Handler: e,
		},
	}
}

func (s *Server) Start(ctx *app_context.AppContext) error {
	ctx.Logger().Info("start server", logger.String("port", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx *app_context.AppContext) error {
	ctx.Logger().Info("stop server")
	return s.server.Shutdown(ctx)
}
