package signal

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
)

type Runner interface {
	Start() error
	Stop(context.Context) error
}

// Run runs r.Start() and waits for either r.Start() to complete successfully, or
// ctx to be canceled. When ctx is canceled, Run calls r.Stop() with a context that
// has a timeout of timeout. If ctx is canceled before r.Start() completes, Run
// returns an error immediately.
//
// sig is a list of os.Signals that can be used to cancel ctx. If sig is empty,
// Run returns an error immediately.
//
// Run returns the first error it encounters, either from r.Start() or r.Stop().
func Run(ctx context.Context, r Runner, timeout time.Duration, sig ...os.Signal) error {
	if len(sig) == 0 {
		return errors.New("signal is empty")
	}
	ctx, stop := signal.NotifyContext(ctx, sig...)
	defer stop()
	errStartChan := make(chan error, 1)
	go func() {
		errStartChan <- r.Start()
	}()

	select {
	case err := <-errStartChan:
		if err != nil {
			return errors.WithMessage(err, "start error")
		}
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := r.Stop(ctx); err != nil {
			return errors.WithMessage(err, "stop error")
		}
	}
	return nil
}