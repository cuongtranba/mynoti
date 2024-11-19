package appfx

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/internal/repository"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
)

func getDbUrl(cfg *config.Config) string {
	return cfg.DatabaseURL
}

func connectDb(lc fx.Lifecycle, databaseUrl string) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, err := postgres.Connect(ctx, databaseUrl)
	if err != nil {
		panic(err)
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return con.Close(ctx)
		},
	})
	return con
}

func newComicQuery(con *pgx.Conn) *comic.Queries {
	if con == nil {
		panic("db connection is nil")
	}
	return comic.New(con)
}

var dbModule = fx.Module(
	"DBModule",
	fx.Provide(
		config.LoadConfig,
		getDbUrl,
		connectDb,
	),
)

var UseCaseComicModule = fx.Module(
	"UseCaseComicModule",
	dbModule,
	htmlFetcherModule,
	fx.Provide(
		newComicQuery,
		repository.NewComicRepository,
		usecase.NewComicUseCase,
	),
)

func newLoggerName(appName string) func() *logger.Logger {
	return func() *logger.Logger {
		return logger.NewLogger(appName)
	}
}

var CLIModule = fx.Module(
	"CLIModule",
	UseCaseComicModule,
	fx.Provide(newLoggerName("cli")),
	fx.Invoke(func(lc fx.Lifecycle, logger *logger.Logger, useCase domain.ComicUseCase) *delivery.Cli {
		var server *delivery.Cli
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				server = delivery.NewCli(logger, useCase, os.Args)
				return server.Run(app_context.New(ctx))
			},
			OnStop: func(ctx context.Context) error {
				ctxc, done := context.WithTimeout(context.Background(), 5*time.Second)
				defer done()
				return server.Stop(app_context.New(ctxc))
			},
		})
		return server
	}),
)

var httpModule = fx.Provide(
	func() *http.Client {
		return http.DefaultClient
	},
)

var htmlFetcherModule = fx.Module(
	"HTMLFetcherModule",
	httpModule,
	fx.Provide(
		usecase.NewHtmlFetcher,
	),
)

var CLIApp = fx.New(
	fx.NopLogger,
	CLIModule,
)

var ServerAPP = fx.New(
	fx.NopLogger,
	UseCaseComicModule,
	fx.Provide(
		newLoggerName("api"),
	),
	fx.Invoke(func(lc fx.Lifecycle, config *config.Config, logger *logger.Logger, useCase domain.ComicUseCase) *delivery.Server {
		var server *delivery.Server
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				server = delivery.NewServer(config.Port, useCase, logger)
				return server.Start(app_context.New(ctx))
			},
			OnStop: func(ctx context.Context) error {
				ctxc, done := context.WithTimeout(context.Background(), 5*time.Second)
				defer done()
				return server.Stop(app_context.New(ctxc))
			},
		})
		return server
	}),
)

type fxRunner struct {
	app *fx.App
}

func NewFxRunner(fxApp *fx.App) *fxRunner {
	return &fxRunner{
		app: fxApp,
	}
}

func (r *fxRunner) Start(ctx *app_context.AppContext) error {
	return r.app.Start(ctx)
}

func (r *fxRunner) Stop(ctx *app_context.AppContext) error {
	return r.app.Stop(ctx)
}
