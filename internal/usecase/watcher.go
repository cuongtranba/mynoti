package usecase

import (
	"context"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type watcher struct {
	cron         *cron.Cron
	notifier     domain.Notifier[domain.Comic]
	htmlFetcher  domain.HtmlFetcher
	comicUseCase domain.ComicUseCase
}

func NewWatcher(
	htmlFetcher domain.HtmlFetcher,
	notifier domain.Notifier[domain.Comic],
	comicUseCase domain.ComicUseCase,
) domain.Watcher {
	return &watcher{
		cron: cron.New(
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithSeconds(),
		),
		notifier:     notifier,
		htmlFetcher:  htmlFetcher,
		comicUseCase: comicUseCase,
	}
}
func (w *watcher) Unregister(ctx *app_context.AppContext, id int32) error {
	w.cron.Remove(cron.EntryID(id))
	return nil
}

// List implements domain.Watcher.
func (w *watcher) List(ctx *app_context.AppContext) ([]domain.Job, error) {
	return nil, nil
}

// Register implements domain.Watcher.
func (w *watcher) Register(ctx *app_context.AppContext, j domain.Job) error {
	ctx.Logger().Info("register job", "id", j.ID, "url", j.Url, "cron", j.JobSpec)
	_, err := w.cron.AddFunc(j.JobSpec, func() {
		ctx := app_context.New(context.WithoutCancel(ctx))
		ctx.Logger().Info("run job", "id", j.ID, "url", j.Url, "cron", j.JobSpec)
		htmlContent, err := w.htmlFetcher.Fetch(ctx, j.Url)
		if err != nil {
			ctx.Logger().Error("failed to fetch html", "err", err.Error())
			return
		}
		comic, err := w.comicUseCase.GetByID(ctx, j.ID)
		if err != nil {
			ctx.Logger().Error("failed to get comic", "err", err.Error())
			return
		}
		if comic == nil {
			ctx.Logger().Error("record not found", "id", j.ID)
			return
		}
		if comic.Html == htmlContent {
			return
		}

		if err := w.notifier.Notify(ctx, domain.Comic{
			ID:   j.ID,
			Url:  j.Url,
			Html: htmlContent,
		}); err != nil {
			ctx.Logger().Error("failed to notify", "err", err.Error())
		}
	})
	return errors.WithMessage(err, "failed to add cron")
}

// Stop implements domain.Watcher.
func (w *watcher) Stop(ctx *app_context.AppContext) error {
	ctx.Logger().Info("stop watcher")
	v := w.cron.Stop()
	<-v.Done()
	ctx.Logger().Info("watcher stopped")
	return nil
}

// Watch implements domain.Watcher.
func (w *watcher) Watch(ctx *app_context.AppContext) error {
	ctx.Logger().Info("start watcher")
	w.cron.Start()
	return nil
}

type WatcherComic struct {
	watcher      domain.Watcher
	comicUseCase domain.ComicUseCase
}

func (w *WatcherComic) Watch(ctx *app_context.AppContext) error {
	return w.watcher.Watch(ctx)
}

func (w *WatcherComic) Stop(ctx *app_context.AppContext) error {
	return w.watcher.Stop(ctx)
}

func (w *WatcherComic) Register(ctx *app_context.AppContext, j domain.Comic) error {
	result, err := w.comicUseCase.Subscribe(ctx, &j)
	if err != nil {
		return errors.WithMessage(err, "failed to subscribe comic")
	}
	if result == nil {
		return errors.New("failed to subscribe comic")
	}
	return w.watcher.Register(ctx, domain.Job{
		ID:      result.ID,
		Url:     result.Url,
		JobSpec: result.CronSpec,
	})
}

func NewWatcherComic(watcher domain.Watcher, comicUseCase domain.ComicUseCase) domain.WatcherComic {
	return &WatcherComic{
		watcher:      watcher,
		comicUseCase: comicUseCase,
	}
}
