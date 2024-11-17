package delivery

import (
	"net/http"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/labstack/echo/v4"
)

type Server struct {
	server *http.Server
}

type Comic struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Html        string `json:"html"`
}

func NewServer(port string, comicUseCase usecase.ComicUseCase) *Server {
	e := echo.New()
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
	ctx.Logger().Info("start server")
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx *app_context.AppContext) error {
	ctx.Logger().Info("stop server")
	return s.server.Shutdown(ctx)
}
