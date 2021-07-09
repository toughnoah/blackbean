package cmd

import (
	"bytes"
	"github.com/bouk/monkey"
	"github.com/mattn/go-shellwords"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/toughnoah/blackbean/pkg/es"
	"net/http"
	"os"
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

func executeCommand(cmdToExecute string, MockTransport http.RoundTripper) (string, error) {
	defer monkey.Unpatch(es.GetProfile)
	defer monkey.Unpatch(InitConfig)
	args, err := shellwords.Parse(cmdToExecute)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	monkey.Patch(es.GetProfile, func() (profile *es.Profile, err error) {
		p := &es.Profile{
			Info: make(map[string]string),
		}
		p.Info["url"] = TestUrl
		p.Info["username"] = TestUsername
		p.Info["password"] = TestPassword
		return p, nil
	})
	monkey.Patch(InitConfig, func() {
		return
	})
	file, _ := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	fd := int(file.Fd())
	fakeTerminal := &MockTerminal{
		toSend:       []byte("password\rpassword\r\x1b[A\r"),
		bytesPerRead: 1,
	}
	root := NewRootCmd(MockTransport, buf, fakeTerminal, fd, args)
	root.SetErr(buf)
	root.SetOut(buf)
	root.SetArgs(args)
	if err = root.Execute(); err != nil {
		return "", err
	}
	out := buf.String()
	return out, nil
}
