package domain

import (
	"context"
)

type Comic struct {
	ID          int32
	Url         string
	Name        string
	Description string
	Html        string
}

type ComicRepository interface {
	Save(context.Context, *Comic) error
	Get(context.Context, int32) (*Comic, error)
	Delete(context.Context, int32) error
	List(context.Context) ([]Comic, error)
}
