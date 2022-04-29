package syncretry

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	RetryError         = errors.New("retry error")
	RetryAttemptsError = errors.New("retry attempts error")
)

type Options struct {
	Delays []int
}

func NewSyncRetry(options *Options) *SyncRetry {
	return &SyncRetry{
		m:       &sync.Mutex{},
		options: options,
	}
}

type SyncRetry struct {
	m       *sync.Mutex
	options *Options
}

func (sq *SyncRetry) Do(ctx context.Context, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}

	if !errors.Is(err, RetryError) {
		return err
	}

	sq.m.Lock()
	defer sq.m.Unlock()

	for _, delay := range sq.options.Delays {
		select {
		case <-time.After(time.Duration(delay) * time.Second):
			break
		case <-ctx.Done():
			return ctx.Err()
		}

		err = f()
		if err == nil {
			return nil
		}

		if errors.Is(err, RetryError) {
			continue
		} else {
			return err
		}
	}

	return RetryAttemptsError
}
