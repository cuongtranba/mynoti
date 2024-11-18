package appfx

import (
	"context"
	"os"
	"time"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
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
	var (
		con *pgx.Conn
		err error
	)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			con, err = postgres.Connect(ctx, databaseUrl)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return con.Close(ctx)
		},
	})
	return con
}

func newComicQuery(con *pgx.Conn) *comic.Queries {
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
	fx.Invoke(func(logger *logger.Logger, useCase usecase.ComicUseCase) *delivery.Cli {
		return delivery.NewCli(logger, useCase, os.Args)
	}),
)

var CLIApp = fx.New(
	CLIModule,
)

var ServerAPP = fx.New(
	UseCaseComicModule,
	fx.Provide(newLoggerName("api")),
	fx.Invoke(func(lc fx.Lifecycle, config *config.Config, logger *logger.Logger, useCase usecase.ComicUseCase) *delivery.Server {
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
