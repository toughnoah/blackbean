package cmd

import (
	"bytes"
	"github.com/bouk/monkey"
	"github.com/mattn/go-shellwords"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/toughnoah/blackbean/pkg/es"
	"net/http"
	"testing"
)

const (
	TestUrl      = "https://test.es.com"
	TestUsername = "test"
	TestPassword = "password"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var _ = BeforeSuite(func() {
	// block all HTTP requests

})

var _ = BeforeEach(func() {
	// remove any mocks

})

var _ = AfterSuite(func() {

})

func executeCommandForTesting(cmdToExecute string, MockTransport http.RoundTripper) (string, error) {
	defer monkey.Unpatch(es.GetEnv)
	args, err := shellwords.Parse(cmdToExecute)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	monkey.Patch(es.GetEnv, func(env string) (url, username, password string, err error) {
		return TestUrl, TestUsername, TestPassword, nil
	})
	root := NewRootCmd(MockTransport)
	root.SetErr(buf)
	root.SetOut(buf)
	root.SetArgs(args)
	if err = root.Execute(); err != nil {
		return "", err
	}
	res := buf.String()
	return res, nil
}
