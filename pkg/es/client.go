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

func GetEnv(env string) (url, username, password string, err error) {
	if viper.Get(configSpec) == nil {
		err = errors.New("error: error reading config file from get env")
		return "", "", "", err
	}
	conf := viper.Get(configSpec).(map[string]interface{})
	if conf == nil {
		err = errors.New("error: can' not find 'cluster' specification in config file")
		return "", "", "", err
	}
	if conf[env] == nil {
		err = errors.New("error: can' not find any cluster env specification in config file")
		return "", "", "", err
	}
	info := conf[env].(map[string]interface{})
	if info[configUsername] == nil {
		err = errors.New("error: can' not find 'username' specification in config file")
		return "", "", "", err
	}
	username = info[configUsername].(string)
	if info[configPassword] == nil {
		err = errors.New("error: can' not find 'password' specification in config file")
		return "", "", "", err
	}
	password = info[configPassword].(string)
	if info[configUrl] == nil {
		err = errors.New("error: can' not find 'url' specification in config file")
		return "", "", "", err
	}
	url = info[configUrl].(string)
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
		log.Fatal("error reading config file from CompleteConfigEnv")
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
