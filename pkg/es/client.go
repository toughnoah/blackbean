package es

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toughnoah/blackbean/pkg/util"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	CurrentSpec = "current"

	ConfigSpec = "cluster"

	ConfigUsername = "username"

	ConfigPassword = "password"

	ConfigUrl = "url"

	EmptyData = "{}"

	EmptyFile = ""
)

func GetProfile() (*Profile, error) {
	profile := &Profile{}
	rootHandler.handler = clusterHandler
	clusterHandler.handler = infoHandle
	rootHandler.Handle(profile)
	if profile.handleErr != nil {
		return nil, profile.handleErr
	}
	return profile, nil
}

func NewEsClient(url, username, password string, transport http.RoundTripper) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Transport: transport,
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new client error")
	}
	return es, nil
}

func CompleteConfigEnv(toComplete string) []string {
	var envArray []string
	if viper.Get(ConfigSpec) == nil {
		log.Fatal("error reading config file for shell completion")
	}
	cfg := viper.Get(ConfigSpec).(map[string]interface{})
	for env := range cfg {
		if strings.HasPrefix(env, toComplete) {
			envArray = append(envArray, env)
		}
	}
	return envArray
}

func NoResourcesError(Resources string) error {
	return errors.New(fmt.Sprintf("no such resources [%s]", Resources))
}

func Validate(noun string, valid []string) error {
	if len(valid) == 0 {
		return errors.New("empty valid resources are not allowed")
	}
	for _, n := range valid {
		if n == noun {
			return nil
		}
	}
	return errors.Errorf("no valid resources exists with the name: %q", noun)
}
func Check(i string, env []string) bool {
	for _, e := range env {
		if i == e {
			return true
		}
	}
	return false
}

func DecodeFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	utf16bom := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	reader := transform.NewReader(f, utf16bom)
	raw := new(json.RawMessage)
	d := util.NewYAMLOrJSONDecoder(reader, 4096)
	if err = d.Decode(raw); err != nil {
		if err == io.EOF {
			return []byte(`{}`), nil
		}
		return nil, fmt.Errorf("error parsing %s: %v", filename, err)
	}
	return *raw, nil
}

type RequestBody struct {
	Filename string
	Data     string
}

func AddRequestBodyFlag(cmd *cobra.Command, body *RequestBody) {
	f := cmd.Flags()
	f.StringVarP(&body.Filename, "filename", "f", "", "get request body from specific file.")
	f.StringVarP(&body.Data, "data", "d", "{}", "specify request body")
}

func GetRawRequestBody(req *RequestBody) (raw []byte, err error) {
	if req.Filename == EmptyFile && req.Data == EmptyData {
		return
	}
	if req.Filename != EmptyFile {
		return DecodeFromFile(req.Filename)
	} else {
		raw = []byte(req.Data)
		return
	}
}

func GetFlagValue(cmd *cobra.Command, flag string) string {
	return cmd.Flags().Lookup(flag).Value.String()
}
