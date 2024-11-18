package usecase

import (
	"context"

	"github.com/cuongtranba/mynoti/internal/domain"
)

type ComicUseCase interface {
	Subscribe(ctx context.Context, comic *domain.Comic) error
}

type comicUseCase struct {
	repo domain.ComicRepository
}

func (w *comicUseCase) Subscribe(ctx context.Context, comic *domain.Comic) error {
	if err := validate.Struct(comic); err != nil {
		return err
	}
	return w.repo.Save(ctx, comic)
}

func NewComicUseCase(repo domain.ComicRepository) ComicUseCase {
	return &comicUseCase{
		repo: repo,
	}
}
