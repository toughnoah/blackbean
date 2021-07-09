package es

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Es Suite")
}

var (
	yamlExample = []byte(`cluster:
  a:
    url: https://a.es.com
    username: blackbean
    password: bulldog
current: a
`)
	yamlExampleWithNoCluster = []byte(`c:
  a:
    url: https://a.es.com
    username: blackbean
    password: bulldog
current: a
`)
	yamlExampleWithNoCluster2 = []byte(`cluster:

current: a
`)
	yamlExampleWithNoUser = []byte(`cluster:
  a:
    url: https://a.es.com
    password: bulldog
current: a
`)
	yamlExampleWithNoPwd = []byte(`cluster:
  a:
    url: https://a.es.com
    username: blackbean
current: a
`)

	yamlExampleWithNoEnv = []byte(`cluster:
  a:
    url: https://a.es.com
    username: blackbean
    password: bulldog
current: b
`)
	yamlExampleWithNoUrl = []byte(`cluster:
  a:
    username: blackbean
    password: bulldog
current: a
`)
	yamlExampleWithBadClusterFormat = []byte(`cluster:
  1
current: a
`)
	yamlExampleWithBadEnvFormat = []byte(`cluster:
  a: 1
current: a
`)
	yamlExampleWithBadUserFormat = []byte(`cluster:
  a:
    username: 1
    password: bulldog
current: a
`)
	yamlExampleWithBadPwdFormat = []byte(`cluster:
  a:
    username: blackbean
    password: 1
current: a
`)
	yamlExampleWithBadUrlFormat = []byte(`cluster:
  a:
    url: 1
    username: blackbean
    password: bulldog
current: a
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
current: prod
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
			}
			for _, tc := range testCases {
				r := bytes.NewReader(tc.yamlExample)
				viper.SetConfigType("yaml")
				err := viper.ReadConfig(r)
				Expect(err).To(BeNil())
				fmt.Println(string(tc.yamlExample))
				profile, err := GetProfile()
				Expect(err).To(tc.err)
				Expect(profile.Info[ConfigUrl]).To(Equal(tc.url))
				Expect(profile.Info[ConfigUsername]).To(Equal(tc.user))
				Expect(profile.Info[ConfigPassword]).To(Equal(tc.pwd))
			}
		})
	})
	Context("test no resource error", func() {
		It("test NoResourcesError", func() {
			err := NoResourcesError("test")
			Expect(err).ShouldNot(BeNil())
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
					res = Check(i, env)
					if !res {
						return res
					}
				}
				return res
			}).Should(Equal(true))
		})
	})
})

func TestValidate(t *testing.T) {

	testCase := []struct {
		name  string
		valid []string
		noun  string
		err   error
		pass  bool
	}{
		{
			name: "test correct noun",
			valid: []string{
				"a",
				"b",
			},
			noun: "a",
			pass: true,
		},
		{
			name: "test correct noun",
			valid: []string{
				"a",
				"b",
			},
			noun: "c",
			err:  fmt.Errorf("no valid resources exists with the name: \"c\""),
		},
		{
			name:  "test correct noun",
			valid: []string{},
			noun:  "c",
			err:   errors.New("empty valid resources are not allowed"),
		},
	}

	for _, tc := range testCase {
		err := Validate(tc.noun, tc.valid)
		if tc.pass {
			require.Equal(t, err, nil)
		} else {
			assert.Error(t, err)
		}

	}

}

func TestCheck(t *testing.T) {
	testCases := []struct {
		name        string
		checkString string
		checkSlice  []string
		want        bool
	}{
		{
			name:        "check success",
			checkString: "a",
			checkSlice:  []string{"a", "b", "c"},
			want:        true,
		},
		{
			name:        "check failed",
			checkString: "a",
			checkSlice:  []string{"b", "c"},
			want:        false,
		},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.want, Check(tc.checkString, tc.checkSlice))
	}
}

func Test_decodeFromFile(t *testing.T) {
	testCase := []struct {
		name string
		file string
		want string
	}{
		{
			name: "test json",
			file: "../pkg/testdata/query.json",
			want: `{"query":{"match_all": {}}}`,
		},
		{
			name: "test yaml",
			file: "../pkg/testdata/query.yaml",
			want: `{"query":{"match":{"name":"test"}}}`,
		},
	}
	for _, tc := range testCase {
		file, err := DecodeFromFile(tc.file)
		if err != nil {
			return
		}
		require.Equal(t, tc.want, string(file))
	}
}

func TestAddRequestBodyFlag(t *testing.T) {
	AddRequestBodyFlag(&cobra.Command{}, new(RequestBody))
}

func TestGetRawRequestBody(t *testing.T) {
	req := new(RequestBody)
	req.Filename = "../testdata/query.yaml"
	req.Data = `{"abc":"1"}`
	body, _ := GetRawRequestBody(req)
	require.Equal(t, "{\"query\":{\"match\":{\"name\":\"test\"}}}", string(body))
	req = new(RequestBody)
	req.Data = `{"abc":"1"}`
	body, _ = GetRawRequestBody(req)
	require.Equal(t, "{\"abc\":\"1\"}", string(body))
}

func TestGetFlagValue(t *testing.T) {
	cmd := &cobra.Command{}
	var test string
	cmd.Flags().StringVar(&test, "test", "1", "for test")
	require.Equal(t, GetFlagValue(cmd, "test"), "1")
}
