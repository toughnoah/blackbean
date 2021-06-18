package es

import (
	"crypto/tls"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func NewEsClient(env string) (*elasticsearch.Client, error) {
	conf := viper.Get(configSpec).(map[string]interface{})
	if conf == nil {
		return nil, errors.New("error: can' not find 'cluster' specification in config file")
	}
	if conf[env] == nil {
		return nil, errors.New("error: can' not find any cluster env specification in config file")
	}
	info := conf[env].(map[string]interface{})
	if info[configUsername] == nil {
		return nil, errors.New("error: can' not find 'username' specification in config file")
	}
	username := info[configUsername].(string)
	if info[configPassword] == nil {
		return nil, errors.New("error: can' not find 'password' specification in config file")
	}
	password := info[configPassword].(string)
	if info[configUrl] == nil {
		return nil, errors.New("error: can' not find 'url' specification in config file")
	}
	url := info[configUrl].(string)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	escfg := elasticsearch.Config{
		Transport: tr,
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	}
	es, err := elasticsearch.NewClient(escfg)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func GetConfigEnv(toComplete string) []string {
	var envArray []string
	cfg := viper.Get(configSpec).(map[string]interface{})
	for env := range cfg {
		if strings.HasPrefix(env, toComplete) {
			envArray = append(envArray, env)
		}
	}
	return envArray
}
