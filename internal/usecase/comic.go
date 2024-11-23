package usecase

import (
	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
)

type comicUseCase struct {
	repo        domain.ComicRepository
	htmlFetcher domain.HtmlFetcher
}

// GetByID implements domain.ComicUseCase.
func (w *comicUseCase) GetByID(ctx *app_context.AppContext, id int32) (*domain.Comic, error) {
	return w.repo.Get(ctx, id)
}

func (w *comicUseCase) Subscribe(ctx *app_context.AppContext, comic *domain.Comic) (*domain.Comic, error) {
	if err := validate.Struct(comic); err != nil {
		return nil, err
	}
	htmlContent, err := w.htmlFetcher.Fetch(ctx, comic.Url)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Debug("html content", logger.String("content", htmlContent))
	comic.Html = htmlContent
	return w.repo.Save(ctx, comic)
}

func NewComicUseCase(
	repo domain.ComicRepository,
	htmlFetcher domain.HtmlFetcher,
) domain.ComicUseCase {
	return &comicUseCase{
		repo:        repo,
		htmlFetcher: htmlFetcher,
	}
}
