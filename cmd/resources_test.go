package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/toughnoah/blackbean/pkg/fake"
)

var _ = Describe("cat resources test", func() {

	Context("test execute es cat command.", func() {
		It("test cat command with valid resource", func() {
			testCases := []struct {
				cmd  string
				mock *fake.MockEsResponse
			}{
				{
					cmd: "get health",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get health"}`,
					},
				},
				{
					cmd: "get allocations",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get health"}`,
					},
				},
				{
					cmd: "get nodes",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get nodes"}`,
					},
				},
				{
					cmd: "get threadpool",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get threadpool"}`,
					},
				},
				{
					cmd: "get cachemem",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get cachemem"}`,
					},
				},
				{
					cmd: "get segmem",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get segmem"}`,
					},
				},
				{
					cmd: "get largeindices",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get largeindices"}`,
					},
				},
			}

			for _, tc := range testCases {
				_, err := executeCommand(tc.cmd, tc.mock)
				Expect(err).Should(BeNil())
			}
		})
		It("test cat command with invalid resource", func() {
			testCases := []struct {
				cmd  string
				mock *fake.MockEsResponse
			}{
				{
					cmd: "get test",
					mock: &fake.MockEsResponse{
						ResponseString: `{"test":"get health"}`,
					},
				},
			}
			for _, tc := range testCases {
				_, err := executeCommand(tc.cmd, tc.mock)
				Expect(err).ShouldNot(BeNil())
			}
		})
	})
})
