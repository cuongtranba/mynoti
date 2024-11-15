package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cuongtranba/mynoti/internal/config"
	"github.com/cuongtranba/mynoti/internal/db/postgres"
	"github.com/cuongtranba/mynoti/internal/delivery"
	"github.com/cuongtranba/mynoti/internal/repository"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/cuongtranba/mynoti/internal/usecase"
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

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httpServer.Start()
	}()

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal, shutting down gracefully...")
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), timeout)
	defer cancelShutdown()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}
