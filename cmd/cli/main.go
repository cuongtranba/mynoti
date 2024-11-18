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
	app := appfx.NewFxRunner(appfx.CLIApp)
	defer app.Stop(app_context.New(context.Background()))
	err := signal.Run(
		app_context.New(context.Background()),
		app,
		3*time.Second,
		syscall.SIGINT,
		os.Interrupt,
		os.Kill,
	)
	if err != nil {
		panic(err)
	}
}
