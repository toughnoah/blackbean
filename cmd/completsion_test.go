package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
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
			out, err := executeCommandForTesting("__complete get ''", nil)
			Expect(err).To(BeNil())
			Expect(strings.Contains(out, "health\nnodes\nallocations\nthreadpool\ncachemem\nsegmem\nlargeindices"))
		})
	})
})
