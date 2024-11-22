package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/cuongtranba/mynoti/internal/appfx"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/signal"
)

func main() {
	if err := signal.Run(
		app_context.New(context.Background()),
		appfx.NewFxRunner(appfx.ServerAPP),
		3*time.Second,
		syscall.SIGINT,
		os.Kill,
		os.Interrupt,
		syscall.SIGTERM,
	); err != nil {
		panic(err)
	}
}
