package main

import (
	"context"
	"os"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
	"github.com/cuongtranba/mynoti/internal/repository"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
)

func main() {
	config := config.LoadConfig()
	con, err := postgres.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer con.Close(context.Background())

	comicRepo := repository.NewComicRepository(comic.New(con))
	comicUseCase := usecase.NewComicUseCase(comicRepo)
	log := logger.NewLogger("cli")
	cli := delivery.NewCli(log, comicUseCase, os.Args)
	if err := cli.Run(app_context.New(context.Background())); err != nil {
		log.Fatal(err)
	}
}
