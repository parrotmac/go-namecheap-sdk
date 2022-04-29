package syncretry

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testRetryDelays = []int{1, 2, 3}

func TestNewSyncRetry(t *testing.T) {
	t.Run("instance", func(t *testing.T) {
		sr := NewSyncRetry(&Options{testRetryDelays})
		assert.NotNil(t, sr)
	})
}

func TestSyncRetry_Do(t *testing.T) {
	t.Run("one_func_success", func(t *testing.T) {
		sr := NewSyncRetry(&Options{testRetryDelays})
		done := false

		err := sr.Do(context.TODO(), func() error {
			done = true
			return nil
		})

		assert.Nil(t, err)
		assert.True(t, done)
	})

	t.Run("two_funcs_sync_success", func(t *testing.T) {
		sr := NewSyncRetry(&Options{testRetryDelays})
		done := 0

		err1 := sr.Do(context.TODO(), func() error {
			done++
			return nil
		})

		err2 := sr.Do(context.TODO(), func() error {
			done++
			return nil
		})

		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, 2, done)
	})

	t.Run("should_forward_error", func(t *testing.T) {
		testError := errors.New("test error")
		sr := NewSyncRetry(&Options{testRetryDelays})

		err := sr.Do(context.TODO(), func() error {
			return testError
		})

		assert.NotNil(t, err)
		assert.ErrorIs(t, testError, err)
	})

	t.Run("two_funcs_parallel_success", func(t *testing.T) {
		sr := NewSyncRetry(&Options{testRetryDelays})
		done := int32(0)

		firstDone := make(chan error)
		secondDone := make(chan error)

		go func() {
			err := sr.Do(context.TODO(), func() error {
				atomic.AddInt32(&done, 1)
				time.Sleep(time.Millisecond * time.Duration(200))
				return nil
			})

			firstDone <- err
		}()

		go func() {
			err := sr.Do(context.TODO(), func() error {
				atomic.AddInt32(&done, 1)
				time.Sleep(time.Millisecond * time.Duration(200))
				return nil
			})

			secondDone <- err
		}()

		err1 := <-firstDone
		err2 := <-secondDone

		assert.Equal(t, int32(2), done)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
	})

	t.Run("one_func_retry_last_success", func(t *testing.T) {
		delays := []int{1, 1, 1}
		sr := NewSyncRetry(&Options{delays})
		count := 0

		err := sr.Do(context.TODO(), func() error {
			if count == len(testRetryDelays) {
				return nil
			}
			count++
			return RetryError
		})

		assert.Nil(t, err)
		assert.Equal(t, len(delays), count)
	})

	t.Run("one_func_exceed_error", func(t *testing.T) {
		delays := []int{1, 1}
		sr := NewSyncRetry(&Options{delays})
		count := 0

		err := sr.Do(context.TODO(), func() error {
			count++
			return RetryError
		})

		assert.ErrorIs(t, RetryAttemptsError, err)
	})

	t.Run("two_func_retry_success", func(t *testing.T) {
		delays := []int{1, 1, 1}
		sr := NewSyncRetry(&Options{delays})

		firstFuncCalls := int32(0)
		secondFuncCalls := int32(0)

		firstDone := make(chan error)
		secondDone := make(chan error)

		go func() {
			count := 0
			err := sr.Do(context.TODO(), func() error {
				count++
				atomic.AddInt32(&firstFuncCalls, 1)
				if count != 2 {
					return RetryError
				}
				return nil
			})

			firstDone <- err
		}()

		go func() {
			count := 0
			err := sr.Do(context.TODO(), func() error {
				count++
				atomic.AddInt32(&secondFuncCalls, 1)
				if count != 2 {
					return RetryError
				}
				return nil
			})

			secondDone <- err
		}()

		err1 := <-firstDone
		err2 := <-secondDone

		assert.Equal(t, int32(2), firstFuncCalls)
		assert.Equal(t, int32(2), secondFuncCalls)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
	})

	t.Run("parallel_funcs_exceeded_error", func(t *testing.T) {
		delays := []int{1, 1}
		sr := NewSyncRetry(&Options{delays})

		firstFuncCalls := int32(0)
		secondFuncCalls := int32(0)

		firstDone := make(chan error)
		secondDone := make(chan error)

		go func() {
			count := 0
			err := sr.Do(context.TODO(), func() error {
				count++
				atomic.AddInt32(&firstFuncCalls, 1)
				return RetryError
			})

			firstDone <- err
		}()

		go func() {
			count := 0
			err := sr.Do(context.TODO(), func() error {
				count++
				atomic.AddInt32(&secondFuncCalls, 1)
				return RetryError
			})

			secondDone <- err
		}()

		err1 := <-firstDone
		err2 := <-secondDone

		assert.Equal(t, int32(3), firstFuncCalls)
		assert.Equal(t, int32(3), secondFuncCalls)
		assert.ErrorIs(t, RetryAttemptsError, err1)
		assert.ErrorIs(t, RetryAttemptsError, err2)
	})
}
