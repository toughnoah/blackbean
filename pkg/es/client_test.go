package es

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/spf13/viper"
)

var (
	yamlExample = []byte(`cluster:
  a:
    url: https://a.es.com
    username: blackbean
    password: bulldog
`)
	yamlExampleWithNoCluster = []byte(`c:
  a:
    url: https://a.es.com
    username: blackbean
    password: bulldog
`)
	yamlExampleWithNoCluster2 = []byte(`cluster:

`)
	yamlExampleWithNoUser = []byte(`cluster:
  a:
    url: https://a.es.com
    password: bulldog
`)
	yamlExampleWithNoPwd = []byte(`cluster:
  a:
    url: https://a.es.com
    username: blackbean
`)
	yamlExampleWithNoUrl = []byte(`cluster:
  a:
    username: blackbean
    password: bulldog
`)
	yamlExampleWithBadClusterFormat = []byte(`cluster:
  1
`)
	yamlExampleWithBadEnvFormat = []byte(`cluster:
  a: 1
`)
	yamlExampleWithBadUserFormat = []byte(`cluster:
  a:
    username: 1
    password: bulldog
`)
	yamlExampleWithBadPwdFormat = []byte(`cluster:
  a:
    username: blackbean
    password: 1
`)
	yamlExampleWithBadUrlFormat = []byte(`cluster:
  a:
    url: 1
    username: blackbean
    password: bulldog
`)
	yamlExampleForShellCompletion = []byte(`cluster:
  prod:
    url: https://a.es.com
    username: blackbean
    password: bulldog
  prd:
    url: https://a.es.com
    username: blackbean
    password: bulldog
  pr:
    url: https://a.es.com
    username: blackbean
    password: bulldog
`)
)

var _ = Describe("put settings test", func() {
	Context("test get env", func() {
		It("test mock env", func() {
			testCases := []struct {
				yamlExample []byte
				url         string
				user        string
				pwd         string
				err         types.GomegaMatcher
				env         string
			}{
				{
					yamlExample: yamlExample,
					url:         "https://a.es.com",
					user:        "blackbean",
					pwd:         "bulldog",
					err:         BeNil(),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithNoCluster,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoClusterErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithNoCluster2,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoClusterErr),
					env:         "a",
				},
				{
					yamlExample: yamlExample,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoEnvErr),
					env:         "c",
				},
				{
					yamlExample: yamlExampleWithNoUser,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoUserErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithNoPwd,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoPwdErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithNoUrl,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(NoUrlErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithBadClusterFormat,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(YamFormatErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithBadEnvFormat,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(YamFormatErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithBadUserFormat,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(YamFormatErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithBadPwdFormat,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(YamFormatErr),
					env:         "a",
				},
				{
					yamlExample: yamlExampleWithBadUrlFormat,
					url:         "",
					user:        "",
					pwd:         "",
					err:         Equal(YamFormatErr),
					env:         "a",
				},
			}
			for _, tc := range testCases {
				r := bytes.NewReader(tc.yamlExample)
				viper.SetConfigType("yaml")
				err := viper.ReadConfig(r)
				Expect(err).To(BeNil())
				url, user, pwd, err := GetEnv(tc.env)
				Expect(err).To(tc.err)
				Expect(url).To(Equal(tc.url))
				Expect(user).To(Equal(tc.user))
				Expect(pwd).To(Equal(pwd))
			}
		})
		Context("test no resource error", func() {
			It("test NoResourcesError", func() {
				err := NoResourcesError("test")
				Expect(err).To(Equal(errors.New("no such resources [test]")))
			})
		})
		Context("test new es client", func() {
			It("test NewEsClient", func() {
				_, err := NewEsClient("https://test.es.com", "test", "test", nil)
				Expect(err).Should(BeNil())
			})
		})
		Context("test get env for shell completion", func() {
			It("test NewEsClient", func() {
				r := bytes.NewReader(yamlExampleForShellCompletion)
				viper.SetConfigType("yaml")
				err := viper.ReadConfig(r)
				Expect(err).To(BeNil())
				env := CompleteConfigEnv("pr")
				Eventually(func() bool {
					var res bool
					for _, i := range []string{"prod", "prd", "pr"} {
						res = check(i, env)
						if !res {
							return res
						}
					}
					return res
				}).Should(Equal(true))
			})
		})
	})
})

func check(i string, env []string) bool {
	for _, e := range env {
		if i == e {
			return true
		}
	}
	return false
}
