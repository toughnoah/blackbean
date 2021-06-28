package es

import (
	"errors"
	"github.com/spf13/viper"
)

var (
	rootHandler    = RootHandler{}
	clusterHandler = &ClusterHandler{}
	infoHandle     = &InfoHandler{}
)

type Handler interface {
	Handle(profile *Profile)
	next(profile *Profile)
}

type Profile struct {
	Info      map[string]string
	env       string
	raw       interface{}
	handleErr error
}

type ClusterHandler struct {
	handler Handler
}

type RootHandler struct {
	handler Handler
}

func (r *RootHandler) Handle(profile *Profile) {
	var ok bool
	if viper.Get(CurrentSpec) == nil {
		profile.handleErr = errors.New("can not read 'current' from config")
		return
	}
	profile.env, ok = viper.Get(CurrentSpec).(string)
	if !ok {
		profile.handleErr = errors.New("can not read 'env' from config")
		return
	}
	if viper.Get(ConfigSpec) == nil {
		profile.handleErr = errors.New("can not read 'cluster' from config")
		return
	}
	r.next(profile)
}

func (r *RootHandler) next(profile *Profile) {
	if profile.handleErr != nil {
		return
	}
	if r.handler != nil {
		r.handler.Handle(profile)
	}
}

func (c *ClusterHandler) Handle(profile *Profile) {
	cluster, ok := viper.Get(ConfigSpec).(map[string]interface{})
	if !ok {
		profile.handleErr = errors.New("can not read 'cluster' from config")
		return
	}
	profile.raw = cluster[profile.env].(map[string]interface{})
	c.next(profile)
}

func (c *ClusterHandler) next(profile *Profile) {
	if profile.handleErr != nil {
		return
	}
	if c.handler != nil {
		c.handler.Handle(profile)
	}
}

type InfoHandler struct {
	handler Handler
}

func (i *InfoHandler) Handle(profile *Profile) {
	rawMap, ok := profile.raw.(map[string]interface{})
	profile.Info = make(map[string]string)
	profile.Info["url"], ok = rawMap["url"].(string)
	profile.Info["username"], ok = rawMap["username"].(string)
	profile.Info["password"], ok = rawMap["password"].(string)
	if !ok {
		profile.handleErr = errors.New("can not get info from 'cluster' map, url, username,password must to be string type")
		return
	}
	i.next(profile)
}

func (i *InfoHandler) next(profile *Profile) {
	if profile.handleErr != nil {
		return
	}
	if i.handler != nil {
		i.handler.Handle(profile)
	}
}
