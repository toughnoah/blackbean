package cmd

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("cat resources test", func() {
	Context("test no filecompletion", func() {
		It("test noCompletions", func() {
			noCompletions(nil, nil, "")
		})
	})
})
