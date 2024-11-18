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

type handlerContext struct {
	ctx *app_context.AppContext
	echo.Context
}

func Warp(handler func(*handlerContext) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		customCtx := &handlerContext{ctx: c.Request().Context().(*app_context.AppContext), Context: c}
		return handler(customCtx)
	}
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
	e.POST("/subscribe", Warp(func(hc *handlerContext) error {
		var comic Comic
		if err := hc.Bind(&comic); err != nil {
			return err
		}
		err := comicUseCase.Subscribe(hc.ctx, &domain.Comic{
			Url:         comic.Url,
			Name:        comic.Name,
			Description: comic.Description,
			Html:        comic.Html,
			CronSpec:    comic.CronSpec,
		})
		if err != nil {
			return err
		}
		return hc.JSON(http.StatusOK, nil)
	}))

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
