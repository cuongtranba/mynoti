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

type Job struct {
	ID      int32
	Url     string
	JobSpec string
}

type ComicRepository interface {
	Save(context.Context, *Comic) error
	Get(context.Context, int32) (*Comic, error)
	Delete(context.Context, int32) error
	List(context.Context) ([]Comic, error)
}

type ComicUseCase interface {
	Subscribe(ctx *app_context.AppContext, comic *Comic) error
	GetByID(ctx *app_context.AppContext, id int32) (*Comic, error)
}

type HtmlFetcher interface {
	Fetch(ctx *app_context.AppContext, url string) (string, error)
}

type Watcher interface {
	Watch(ctx *app_context.AppContext) error
	Stop(ctx *app_context.AppContext) error
	Register(ctx *app_context.AppContext, j Job) error
	Unregister(ctx *app_context.AppContext, id int32) error
	List(ctx *app_context.AppContext) ([]Job, error)
}

type Notifier[T any] interface {
	Notify(ctx *app_context.AppContext, data T) error
}

type WatcherComic interface {
	Register(ctx *app_context.AppContext, j Comic) error
}
