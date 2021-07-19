package es

import (
	"errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	rootHandler            = RootHandler{}
	clusterHandler         = &ClusterHandler{}
	infoHandle             = &InfoHandler{}
	_              Handler = &RootHandler{}
	_              Handler = &ClusterHandler{}
	_              Handler = &InfoHandler{}
)

type Handler interface {
	Handle(profile *Profile)
	next(profile *Profile)
}

type Profile struct {
	ClusterInfo *ClusterInfo
	env         string
	raw         interface{}
	handleErr   error
}

type ClusterInfo struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}

type ClusterHandler struct {
	handler Handler
}

type RootHandler struct {
	handler Handler
}

func (r *RootHandler) Handle(profile *Profile) {
	var ok bool
	if viper.Get(CurrentSpec) == nil || viper.Get(ConfigSpec) == nil {
		profile.handleErr = errors.New("can not read 'current/cluster' spec from .blackbean")
		return
	}
	profile.env, ok = viper.Get(CurrentSpec).(string)
	if !ok {
		profile.handleErr = errors.New("bad 'current' type from .blackbean, want string")
		return
	}
	r.next(profile)
}

func (r *RootHandler) next(profile *Profile) {
	if r.handler != nil {
		r.handler.Handle(profile)
	}
}

func (c *ClusterHandler) Handle(profile *Profile) {
	cluster, ok := viper.Get(ConfigSpec).(map[string]interface{})
	if !ok {
		profile.handleErr = errors.New("can not read 'cluster' from .blackbean")
		return
	}
	profile.raw = cluster[profile.env]
	c.next(profile)
}

func (c *ClusterHandler) next(profile *Profile) {
	if c.handler != nil {
		c.handler.Handle(profile)
	}
}

type InfoHandler struct {
	handler Handler
}

func (i *InfoHandler) Handle(profile *Profile) {
	ci := &ClusterInfo{}
	marshal, err := yaml.Marshal(profile.raw)
	if err != nil {
		profile.handleErr = err
		return
	}
	err = yaml.Unmarshal(marshal, ci)
	if err != nil {
		profile.handleErr = err
		return
	}
	profile.ClusterInfo = ci
	i.next(profile)
}

func (i *InfoHandler) next(profile *Profile) {
	if i.handler != nil {
		i.handler.Handle(profile)
	}
}
