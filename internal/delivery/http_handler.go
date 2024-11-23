package delivery

import (
	"net/http"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/cuongtranba/mynoti/pkg/middleware"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server        *http.Server
	swatcherComic domain.WatcherComic
}

type handlerContext struct {
	ctx *app_context.AppContext
	echo.Context
}

func Warp(handler func(*handlerContext) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		customCtx := &handlerContext{ctx: c.Request().Context().(*app_context.AppContext), Context: c}
		if err := handler(customCtx); err != nil {
			customCtx.ctx.Logger().Error("failed to handle request", "err", err.Error())
			return err
		}
		return nil
	}
}

func NewServer(
	port string,
	swatcherComic domain.WatcherComic,
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
		err := swatcherComic.Register(hc.ctx, domain.Comic{
			Url:         comic.Url,
			Name:        comic.Name,
			Description: comic.Description,
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
		swatcherComic: swatcherComic,
	}
}

func (s *Server) Start(ctx *app_context.AppContext) error {
	ctx.Logger().Info("start server", logger.String("port", s.server.Addr))
	go func() {
		s.swatcherComic.Watch(ctx)
	}()
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx *app_context.AppContext) error {
	ctx.Logger().Info("stop server")
	go func() {
		<-ctx.Done()
		if err := s.swatcherComic.Stop(ctx); err != nil {
			ctx.Logger().Error("failed to stop watcher", "err", err.Error())
		}
	}()
	return s.server.Shutdown(ctx)
}
