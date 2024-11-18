package main

import (
	"context"

	"github.com/cuongtranba/mynoti/internal/appfx"
)

func main() {
	if err := appfx.CLIApp.Start(context.Background()); err != nil {
		panic(err)
	}
}
