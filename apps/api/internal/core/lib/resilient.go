package lib

import (
	"context"
	"sync"
	"time"

	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

// ResilientWorker wraps a function that runs continuously and restarts it on failure.
type ResilientWorker struct {
	name         string
	work         func(ctx context.Context) error
	restartDelay time.Duration
	maxRetries   int
	mu           sync.Mutex
	running      bool
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewResilientWorker creates a new resilient worker.
func NewResilientWorker(
	name string,
	work func(ctx context.Context) error,
	options ...ResilientOption,
) *ResilientWorker {
	w := &ResilientWorker{
		name:         name,
		work:         work,
		restartDelay: 5 * time.Second,
		maxRetries:   0, // 0 = infinite
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

type ResilientOption func(*ResilientWorker)

func WithRestartDelay(delay time.Duration) ResilientOption {
	return func(w *ResilientWorker) {
		w.restartDelay = delay
	}
}

func WithMaxRetries(max int) ResilientOption {
	return func(w *ResilientWorker) {
		w.maxRetries = max
	}
}

// Start begins the resilient worker.
func (w *ResilientWorker) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return nil
	}
	w.running = true
	ctx, w.cancel = context.WithCancel(ctx)
	w.mu.Unlock()

	w.wg.Add(1)
	go w.run(ctx)
	return nil
}

// Stop gracefully stops the worker.
func (w *ResilientWorker) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
	}
	w.mu.Unlock()
	w.wg.Wait()
	w.mu.Lock()
	w.running = false
	w.mu.Unlock()
}

func (w *ResilientWorker) run(ctx context.Context) {
	defer w.wg.Done()

	retries := 0
	for {
		select {
		case <-ctx.Done():
			logger.InfoCtx(ctx, "worker stopped", "worker", w.name)
			return
		default:
		}

		logger.InfoCtx(ctx, "worker starting", "worker", w.name)
		err := w.work(ctx)

		w.mu.Lock()
		if !w.running {
			w.mu.Unlock()
			return
		}
		w.mu.Unlock()

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			logger.ErrorCtx(ctx, err, "worker failed, will restart", "worker", w.name)
			retries++
			if w.maxRetries > 0 && retries >= w.maxRetries {
				logger.ErrorCtx(
					ctx,
					nil,
					"worker exceeded max retries, giving up",
					"worker",
					w.name,
					"retries",
					retries,
				)
				return
			}
		} else {
			retries = 0
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(w.restartDelay):
		}
	}
}
