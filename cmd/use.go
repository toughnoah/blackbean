package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toughnoah/blackbean/pkg/es"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var validArgs []string

func useCluster() *cobra.Command {
	var command = &cobra.Command{
		Use:   "use [cluster]",
		Short: "change current cluster context",
		Long:  "change current cluster context ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			validArgs = es.CompleteConfigEnv(toComplete)
			return validArgs, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			m := &Modify{}
			if err := m.ModifyCurrentCluster(args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	return command
}

type Modify struct {
	err error
}

func (m *Modify) ModifyCurrentCluster(cluster string) error {
	var blackbeanConfig map[string]interface{}
	path := m.GetConfig()
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &blackbeanConfig)
	if err != nil {
		return err
	}
	checked := m.CheckClusterConfigExists(cluster)
	if !checked {
		return fmt.Errorf("no valid resources exists with the name: %q", cluster)
	}

	blackbeanConfig[es.CurrentSpec] = cluster
	m.ModifyConfigFile(path, blackbeanConfig)
	return m.err
}

func (m *Modify) ModifyConfigFile(path string, config map[string]interface{}) {
	if m.err != nil {
		return
	}
	bytesFile, err := yaml.Marshal(config)
	if err = ioutil.WriteFile(path, bytesFile, 0755); err != nil {
		m.err = err
	}
}
func (m *Modify) GetConfig() string {
	var path string
	if cfgFile != "" {
		path = cfgFile
	} else {
		home, err := homedir.Dir()
		if err != nil {
			m.err = err
			return ""
		}
		path = filepath.Join(home, ".blackbean.yaml")
	}
	return path
}
func (m *Modify) CheckClusterConfigExists(cluster string) (checked bool) {
	if m.err != nil {
		return
	}

	if viper.Get(es.ConfigSpec) == nil {
		m.err = es.NoClusterErr
		return
	}
	clusterMap, ok := viper.Get(es.ConfigSpec).(map[string]interface{})
	if !ok {
		m.err = es.YamFormatErr
		return
	}
	for k, _ := range clusterMap {
		if k == cluster {
			checked = true
		}
	}
	return
}
