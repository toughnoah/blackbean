package es

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"log"
	"net/http"
)


func NewEsClient(env string) *elasticsearch.Client {
	conf := viper.Get("es").(map[string]interface{})
	info := conf[env].(map[string]string)
	username := info["username"]
	password := info["password"]
	url := info["url"]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	config := elasticsearch.Config{
		Transport: tr,
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	}
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}