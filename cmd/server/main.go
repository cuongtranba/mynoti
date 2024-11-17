package main

import (
	"context"
	"syscall"
	"time"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
	"github.com/cuongtranba/mynoti/internal/repository"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/cuongtranba/mynoti/pkg/signal"
)

const timeout = 10 * time.Second

func main() {
	config := config.LoadConfig()
	con, err := postgres.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer con.Close(context.Background())

	comicRepo := repository.NewComicRepository(comic.New(con))
	comicUseCase := usecase.NewComicUseCase(comicRepo)
	log := logger.NewLogger("api")
	httpServer := delivery.NewServer(config.Port, comicUseCase, log)
	if err := signal.Run(app_context.New(context.Background()), httpServer, timeout, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}
