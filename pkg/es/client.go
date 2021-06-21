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

var (
	NoClusterErr = errors.New("error: read config file from getting cluster")
	YamFormatErr = errors.New("error: bad config file format, please check docs for guidance")
	NoEnvErr     = errors.New("error: can' not find any env specification in config file")
	NoUserErr    = errors.New("error: can' not find 'username' specification in config file")
	NoPwdErr     = errors.New("error: can' not find 'password' specification in config file")
	NoUrlErr     = errors.New("error: can' not find 'url' specification in config file")
)

func GetEnv(env string) (url, username, password string, err error) {
	if viper.Get(configSpec) == nil {
		return "", "", "", NoClusterErr
	}
	conf, ok := viper.Get(configSpec).(map[string]interface{})
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
	if info[configUsername] == nil {
		return "", "", "", NoUserErr
	}
	username, ok = info[configUsername].(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	if info[configPassword] == nil {
		return "", "", "", NoPwdErr
	}
	password, ok = info[configPassword].(string)
	if !ok {
		return "", "", "", YamFormatErr
	}
	if info[configUrl] == nil {
		return "", "", "", NoUrlErr
	}
	url, ok = info[configUrl].(string)
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
	if viper.Get(configSpec) == nil {
		log.Fatal("error reading config file for shell completion")
	}
	cfg := viper.Get(configSpec).(map[string]interface{})
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
