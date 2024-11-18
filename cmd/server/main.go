package main

import (
	"context"

	"github.com/cuongtranba/mynoti/internal/appfx"
)

func main() {
	if err := appfx.ServerAPP.Start(context.Background()); err != nil {
		panic(err)
	}
}
