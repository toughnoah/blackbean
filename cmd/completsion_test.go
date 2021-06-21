package cmd

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cat resources test", func() {
	Context("test no filecompletion", func() {
		It("test noCompletions", func() {
			noCompletions(nil, nil, "")
		})
		It("test noCompletions", func() {
			testCases := []struct {
				cmd string
			}{
				{
					cmd: "completion bash",
				},
				{
					cmd: "completion zsh",
				},
				{
					cmd: "completion fish",
				},
				{
					cmd: "completion powershell",
				},
			}
			for _, tc := range testCases {
				_, err := executeCommandForTesting(tc.cmd, nil)
				Expect(err).To(BeNil())
			}
		})
		It("test noCompletions", func() {
			res, err := executeCommandForTesting("__complete get ''", nil)
			Expect(err).To(BeNil())
			fmt.Println(res)
		})
	})
})
