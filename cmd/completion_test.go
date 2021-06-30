package cmd

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/toughnoah/blackbean/pkg/fake"
	"strings"
)

var yamlExample = []byte(`cluster:
  default: 
    url: https://a.es.com:9200
    username: Noah
    password: abc
  backup: 
    url: https://a.es.com:9200
    username: Noah
    password: abc
current: default`)

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
					cmd: "completion bash ",
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
				_, err := executeCommand(tc.cmd, nil)
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
				mock     *fake.MockEsResponse
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
					cmd:      "__complete repo get ''",
					checkOut: "repoA\nrepoB\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"repoA":"a","repoB":"b"}`,
					},
				},
				{
					cmd:      "__complete repo get test -s ''",
					checkOut: "snapshot01\nsnapshot01\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
					},
				},
				{
					cmd:      "__complete repo delete ''",
					checkOut: "repoA\nrepoB\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"repoA":"a","repoB":"b"}`,
					},
				},
				{
					cmd:      "__complete snapshot get -r ''",
					checkOut: "repoA\nrepoB\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"repoA":"a","repoB":"b"}`,
					},
				},
				{
					cmd:      "__complete snapshot get -r repoA ''",
					checkOut: "snapshot01\nsnapshot01\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
					},
				},
				{
					cmd:      "__complete snapshot delete -r repoA ''",
					checkOut: "snapshot01\nsnapshot01\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
					},
				},
				{
					cmd:      "__complete snapshot delete -r ''",
					checkOut: "repoA\nrepoB\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"repoA":"a","repoB":"b"}`,
					},
				},
				{
					cmd:      "__complete index search ''",
					checkOut: "index1\nindex2\n",
					mock: &fake.MockEsResponse{
						ResponseString: `{"index1":"a","index2":"b"}`,
					},
				},
			}
			for _, tc := range testCases {
				out, err := executeCommand(tc.cmd, tc.mock)
				Expect(err).To(BeNil())
				Expect(strings.Contains(out, tc.checkOut))
			}
		})
	})
})
