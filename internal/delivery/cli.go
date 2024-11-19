package delivery

import (
	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/urfave/cli/v2"
)

type Cli struct {
	app  *cli.App
	args []string
}

func NewCli(logger *logger.Logger, useCase domain.ComicUseCase, args []string) *Cli {
	app := &cli.App{
		Name:        "mynoti",
		Usage:       `A tool to manage your comic subscription.`,
		Description: `mynoti is a tool to manage your comic subscription. Use the subscribe command to start tracking your comic.`,
		UsageText:   `mynoti [global options] command [command options] [arguments...]`,
		Commands: []*cli.Command{
			newSubscribeCommand(logger, useCase),
		},
	}
	return &Cli{
		app:  app,
		args: args,
	}
}

func (c *Cli) Run(ctx *app_context.AppContext) error {
	return c.app.RunContext(ctx, c.args)
}

func (c *Cli) Stop(ctx *app_context.AppContext) error {
	return nil
}

func newSubscribeCommand(logger *logger.Logger, useCase domain.ComicUseCase) *cli.Command {
	return &cli.Command{
		Name:    "subscribe",
		Aliases: []string{"s"},
		Usage:   "subscribe to a comic",
		Flags:   newComicFlags(),
		Action: func(cCtx *cli.Context) error {
			return handleSubscribeAction(cCtx, logger, useCase)
		},
	}
}

func newComicFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Usage:    "the URL of the comic",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "the name of the comic",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"d"},
			Usage:    "the description of the comic",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "cron-spec",
			Aliases:  []string{"c"},
			Usage:    "the cron spec of the comic",
			Required: true,
		},
	}
}

func handleSubscribeAction(cCtx *cli.Context, logger *logger.Logger, useCase domain.ComicUseCase) error {
	comic := &domain.Comic{
		Url:         cCtx.String("url"),
		Name:        cCtx.String("name"),
		Description: cCtx.String("description"),
		CronSpec:    cCtx.String("cron-spec"),
	}
	ctx := app_context.New(cCtx.Context).WithLogger(logger)
	return useCase.Subscribe(ctx, comic)
}
