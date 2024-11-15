package domain

import (
	"context"

	"github.com/google/uuid"
)

type Comic struct {
	ID          int64
	Url         string
	Name        string
	Description string
	Html        string
}

type ComicRepository interface {
	Save(context.Context, *Comic) error
	Get(context.Context, uuid.UUID) (*Comic, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context) ([]Comic, error)
}
