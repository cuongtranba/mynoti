package main

import (
	"context"
	"log"
	"syscall"
	"time"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
	"github.com/cuongtranba/mynoti/internal/repository"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/signal"
)

const timeout = 10 * time.Second

func main() {
	config := config.LoadConfig()
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	con, err := postgres.Connect(ctx, config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer con.Close(ctx)

	comicRepo := repository.NewComicRepository(comic.New(con))
	comicUseCase := usecase.NewComicUseCase(comicRepo)
	httpServer := delivery.NewServer(config.Port, comicUseCase)
	if err := signal.Run(ctx, httpServer, timeout, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}
