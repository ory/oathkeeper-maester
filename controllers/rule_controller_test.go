package controllers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const delay = time.Millisecond

func TestRetryOnError(t *testing.T) {

	var cnt, attempts int
	var err error
	var createMapFunc func() error

	t.Run("should retry on error", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 2

		createMapFunc = func() error {
			if cnt == 0 {
				cnt++
				return errors.New("error only on first invocation")
			}
			return nil
		}

		//when
		err = retryOnError(createMapFunc, attempts, delay)

		//then
		require.Nil(t, err)
	})

	t.Run("should retry 4 times", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 5

		createMapFunc = func() error {
			if cnt < 4 {
				cnt++
				return errors.New(fmt.Sprintf("error only on first four invocations (current: %d)", cnt))
			}
			return nil
		}

		//when
		err = retryOnError(createMapFunc, attempts, delay)

		//then
		require.Nil(t, err)
	})

	t.Run("should give up after five failed attempts", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 5

		createMapFunc = func() error {
			cnt++
			return errors.New(fmt.Sprintf("error on every invocation (current: %d)", cnt))
		}

		//when
		err = retryOnError(createMapFunc, attempts, delay)

		//then
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "error on every invocation")
	})
}
