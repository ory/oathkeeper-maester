// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oathkeeper controller")
}
