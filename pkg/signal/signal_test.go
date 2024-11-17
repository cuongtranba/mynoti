package signal

import (
	"context"
	"syscall"
	"testing"
	"time"
)

type mockRunner struct {
	startErr error
	stopErr  error
	started  bool
	stopped  bool
}

func (m *mockRunner) Start() error {
	if m.startErr != nil {
		return m.startErr
	}
	m.started = true
	return nil
}

func (m *mockRunner) Stop(ctx context.Context) error {
	if m.stopErr != nil {
		return m.stopErr
	}
	m.stopped = true
	return nil
}

func TestRun_Start_Stop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := &mockRunner{}
	done := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		err := Run(ctx, r, 3*time.Second, syscall.SIGINT)
		errCh <- err
		close(errCh)
		close(done)
	}()
	cancel()
	<-done

	err := <-errCh
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !r.started {
		t.Errorf("expected runner to be started")
	}
	if !r.stopped {
		t.Errorf("expected runner to be stopped")
	}
}
