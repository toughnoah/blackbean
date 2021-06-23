package es

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
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

var (
	NoCurrentError = errors.New("read config file failed from getting current")
	NoClusterErr   = errors.New("read config file failed from getting cluster")
	YamFormatErr   = errors.New("bad config file format, please check docs for guidance")
	NoEnvErr       = errors.New("can' not find any env specification in config file")
	NoUserErr      = errors.New("can' not find 'username' specification in config file")
	NoPwdErr       = errors.New("can' not find 'password' specification in config file")
	NoUrlErr       = errors.New("can' not find 'url' specification in config file")
)

func GetProfile() (url, username, password string, err error) {
	if viper.Get(CurrentSpec) == nil {
		return "", "", "", NoCurrentError
	}
	env, ok := viper.Get(CurrentSpec).(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	if viper.Get(ConfigSpec) == nil {
		return "", "", "", NoClusterErr
	}
	conf, ok := viper.Get(ConfigSpec).(map[string]interface{})
	if !ok {
		return "", "", "", YamFormatErr
	}
	if conf[env] == nil {
		return "", "", "", NoEnvErr
	}
	info, ok := conf[env].(map[string]interface{})
	if !ok {
		return "", "", "", YamFormatErr
	}
	if info[ConfigUsername] == nil {
		return "", "", "", NoUserErr
	}
	username, ok = info[ConfigUsername].(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	if info[ConfigPassword] == nil {
		return "", "", "", NoPwdErr
	}
	password, ok = info[ConfigPassword].(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	if info[ConfigUrl] == nil {
		return "", "", "", NoUrlErr
	}
	url, ok = info[ConfigUrl].(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	return url, username, password, nil
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
		return nil, err
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

func Validate(nouns string, valid []string) error {
	if len(valid) == 0 {
		return errors.New("empty valid resources are not allowed")
	}
	for _, noun := range valid {
		if noun == nouns {
			return nil
		}
	}
	return fmt.Errorf("no valid resources exists with the name: %q", nouns)
}
