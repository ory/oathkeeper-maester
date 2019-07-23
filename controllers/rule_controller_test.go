package controllers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const data = "[]"

func TestCreateConfigmap(t *testing.T) {

	var cnt int
	var err error
	var createMapFunc func(string) error

	t.Run("should retry on error", func(t *testing.T) {

		//Given
		cnt = 0

		createMapFunc = func(data string) error {
			if cnt == 0 {
				cnt++
				return errors.New("error only on first invocation")
			}
			return nil
		}

		//when
		err = createConfigMap(data, createMapFunc)

		//then
		require.Nil(t, err)
	})

	t.Run("should retry 4 times", func(t *testing.T) {

		//Given
		cnt = 0

		createMapFunc = func(data string) error {
			if cnt < 4 {
				cnt++
				return errors.New(fmt.Sprintf("error only on first four invocations (current: %d)", cnt))
			}
			return nil
		}

		//when
		err = createConfigMap(data, createMapFunc)

		//then
		require.Nil(t, err)
	})

	t.Run("should give up after five failed attempts", func(t *testing.T) {

		//Given
		cnt = 0

		createMapFunc = func(data string) error {
			cnt++
			return errors.New(fmt.Sprintf("error on every invocation (current: %d)", cnt))
		}

		//when
		err = createConfigMap(data, createMapFunc)

		//then
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "error on every invocation")
	})
}
