package domain

import (
	"context"
)

type Comic struct {
	ID          int32
	Url         string `validate:"required"`
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
