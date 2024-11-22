package usecase

import (
	"encoding/json"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
)

type echoNotifier[T any] struct {
}

// Notify implements domain.Notifier.
func (e *echoNotifier[T]) Notify(ctx *app_context.AppContext, data T) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.Logger().Info("notifier", logger.String("data", string(b)))
	return nil
}

func NewEchoNotifier() domain.Notifier[domain.Comic] {
	return &echoNotifier[domain.Comic]{}
}
