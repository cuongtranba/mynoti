package usecase

import (
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
		cron:         cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
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
	v := w.cron.Stop()
	<-v.Done()
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

// Register implements domain.WatcherComic.
func (w *WatcherComic) Register(ctx *app_context.AppContext, j domain.Comic) error {
	if err := w.comicUseCase.Subscribe(ctx, &j); err != nil {
		return err
	}
	return w.watcher.Register(ctx, domain.Job{
		ID:      j.ID,
		Url:     j.Url,
		JobSpec: j.CronSpec,
	})
}

func NewWatcherComic(watcher domain.Watcher, comicUseCase domain.ComicUseCase) domain.WatcherComic {
	return &WatcherComic{
		watcher:      watcher,
		comicUseCase: comicUseCase,
	}
}
