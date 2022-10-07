// Copyright © 2022 Ory Corp

package integration

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
)

// Modify this var to provide different function for logging retires.
var logRetry = func(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

// Function to be executed - and possibly retried.
type workerFunc func() (interface{}, error)

// Helper function introduced to inject maxRetries into retry logging - so that users provide maxRetries only once, in `withRetries` function
type getOnRetryLoggingFunc func(maxRetries uint) retry.OnRetryFunc

// Generic function to handle retries.
// getRetryLogging can be nil if you don't want to log retries. To log retries, use the result of `onRetryLogMsg` as an argument.
func withRetries(maxRetries uint, delay time.Duration, getRetryLogging getOnRetryLoggingFunc, worker workerFunc) (interface{}, error) {

	var response interface{} = nil

	err := retry.Do(func() error {
		var err error
		response, err = worker()

		if err != nil {
			return err
		}
		return nil
	},
		retry.Attempts(maxRetries),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(getRetryLogging(maxRetries)),
	)

	return response, err
}

// Returns `getOnRetryLoggingFunc` instance for given message.
func onRetryLogMsg(msg string) getOnRetryLoggingFunc {
	return func(maxRetries uint) retry.OnRetryFunc {
		return func(retryNo uint, err error) {
			logRetry("[%d / %d] %s: %v\n", retryNo+1, maxRetries, msg, err)
		}
	}
}
