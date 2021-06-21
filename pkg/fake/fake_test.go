package fake

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fake Suite")
}

var _ = Describe("put settings test", func() {
	Context("test mockTransport", func() {
		It("test mockTransport", func() {
			mock := MockEsResponse{
				`{"fake":"test"}`,
			}
			_, err := mock.RoundTrip(nil)
			Expect(err).To(BeNil())
		})
	})
})
