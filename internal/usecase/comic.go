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

func (w *comicUseCase) Subscribe(ctx *app_context.AppContext, comic *domain.Comic) error {
	if err := validate.Struct(comic); err != nil {
		return err
	}
	htmlContent, err := w.htmlFetcher.Fetch(ctx, comic.Url)
	if err != nil {
		return err
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
