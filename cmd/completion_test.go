package cmd

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"strings"
)

var yamlExample = []byte(`cluster:
  default: 
    url: https://a.es.com:9200
    username: Noah
    password: abc
  testCompletion:
    url: https://b.es.com:9200
    username: blackbean
    password: abc`)

var _ = Describe("cat resources test", func() {
	Context("test no FileCompletion", func() {
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
			r := bytes.NewReader(yamlExample)
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(r)
			Expect(err).To(BeNil())
			testCases := []struct {
				cmd      string
				checkOut string
			}{
				{
					cmd:      "__complete get ''",
					checkOut: "health\nnodes\nallocations\nthreadpool\ncachemem\nsegmem\nlargeindices",
				},
				{
					cmd:      "__complete apply settings --allocation_enable ''",
					checkOut: "primaries\nnull\n",
				},
				{
					cmd:      "__complete apply settings --cluster ''",
					checkOut: "testCompletion\n",
				},
			}
			for _, tc := range testCases {
				out, err := executeCommandForTesting(tc.cmd, nil)
				Expect(err).To(BeNil())
				Expect(strings.Contains(out, tc.checkOut))
			}
		})
	})
})
