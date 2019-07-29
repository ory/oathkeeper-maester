package controllers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const delay = time.Millisecond

func TestRetryOnErrorWith(t *testing.T) {

	assert := assert.New(t)

	var cnt, attempts int
	var createMapFunc func() error

	t.Run("should retry on error and exit on first successful attempt with no error", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 3

		createMapFunc = func() error {
			cnt++
			if cnt == 1 {
				return errors.New("error only on first invocation")
			}
			return nil
		}

		//when
		err := retryOnErrorWith(createMapFunc, attempts, delay)

		//then
		assert.Nil(err)
		assert.NotEqual(attempts, cnt)
		assert.Equal(2, cnt)
	})

	t.Run("should return no error if the last attempt is successful", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 5

		createMapFunc = func() error {
			cnt++
			if cnt < 5 {
				return errors.New(fmt.Sprintf("error only on first four invocations (current: %d)", cnt))
			}
			return nil
		}

		//when
		err := retryOnErrorWith(createMapFunc, attempts, delay)

		//then
		assert.Nil(err)
		assert.Equal(attempts, cnt)
	})

	t.Run("should return an error if all attempts fail", func(t *testing.T) {

		//Given
		cnt, attempts = 0, 2

		createMapFunc = func() error {
			cnt++
			return errors.New(fmt.Sprintf("error on every invocation (current: %d)", cnt))
		}

		//when
		err := retryOnErrorWith(createMapFunc, attempts, delay)

		//then
		assert.NotNil(err)
		assert.Equal(attempts, cnt)
	})
}
