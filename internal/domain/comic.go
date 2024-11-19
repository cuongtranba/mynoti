package domain

import (
	"context"

	"github.com/cuongtranba/mynoti/pkg/app_context"
)

type Comic struct {
	ID          int32
	Url         string `validate:"required,url"`
	Name        string `validate:"required"`
	Description string `validate:"required"`
	Html        string
	CronSpec    string `validate:"required"`
}

type ComicRepository interface {
	Save(context.Context, *Comic) error
	Get(context.Context, int32) (*Comic, error)
	Delete(context.Context, int32) error
	List(context.Context) ([]Comic, error)
}

type ComicUseCase interface {
	Subscribe(ctx *app_context.AppContext, comic *Comic) error
}

type HtmlFetcher interface {
	Fetch(ctx *app_context.AppContext, url string) (string, error)
}
