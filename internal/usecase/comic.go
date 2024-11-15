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
	return w.repo.Save(ctx, comic)
}

func NewComicUseCase(repo domain.ComicRepository) ComicUseCase {
	return &comicUseCase{
		repo: repo,
	}
}
