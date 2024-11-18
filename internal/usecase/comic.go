package usecase

import (
	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
)

type ComicUseCase interface {
	Subscribe(ctx *app_context.AppContext, comic *domain.Comic) error
}

type comicUseCase struct {
	repo domain.ComicRepository
}

func (w *comicUseCase) Subscribe(ctx *app_context.AppContext, comic *domain.Comic) error {
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
