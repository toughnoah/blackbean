package es

import (
	"github.com/spf13/viper"
	"strings"
)


func GetConfigEnv(toComplete string) []string {
	var envArray  []string
	cfg :=viper.Get("es").(map[string]interface{})
	for env :=range cfg{
		if strings.HasPrefix(env,toComplete){
			envArray = append(envArray, env)
		}
	}
	return envArray
}