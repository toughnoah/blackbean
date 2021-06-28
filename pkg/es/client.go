package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

const (
	CurrentSpec = "current"

	ConfigSpec = "cluster"

	ConfigUsername = "username"

	ConfigPassword = "password"

	ConfigUrl = "url"
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
