package controllers

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Oathkeeper controller", func() {

	Context("createConfigmap function", func() {
		It("should retry on error", func() {

			//Given
			cnt := 0
			data := "should be valid json"

			createMapFunc := func(data string) error {
				if cnt == 0 {
					cnt = 1
					return errors.New("error only on first invocation")
				}
				return nil
			}

			err := createConfigMap(data, createMapFunc)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should retry four times", func() {

			//Given
			cnt := 0
			data := "should be valid json"

			createMapFunc := func(data string) error {
				if cnt < 4 {
					cnt += 1
					return errors.New(fmt.Sprintf("error only on first four invocations (current: %d)", cnt))
				}
				return nil
			}

			err := createConfigMap(data, createMapFunc)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should give up after five failed attempts", func() {

			//Given
			cnt := 0
			data := "should be valid json"

			createMapFunc := func(data string) error {
				cnt++
				return errors.New(fmt.Sprintf("error on every invocation (current: %d)", cnt))
			}

			err := createConfigMap(data, createMapFunc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error on every invocation"))
		})

	})
})
