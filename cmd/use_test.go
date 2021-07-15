package cmd

import (
	"bytes"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestUse(t *testing.T) {
	fs := afero.NewOsFs()
	filename := ".blackbean.yaml"
	home, err := homedir.Dir()

	require.NoError(t, err)
	file := filepath.Join(home, filename)
	err = ioutil.WriteFile(file, yamlExample, 0755)
	require.NoError(t, err)

	InitConfig()
	out, err := executeCommand("use backup", nil)
	require.NoError(t, err)
	require.Equal(t, out, "change to cluster: backup\n\n")
	readFile, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	require.Equal(t, true, strings.Contains(string(readFile), "current: backup"))
	defer func(fs afero.Fs, name string) {
		_ = fs.Remove(name)
	}(fs, file)
}

func TestModify_CheckClusterConfigExists(t *testing.T) {
	m := new(Modify)
	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(bytes.NewReader([]byte("")))
	checked := m.CheckClusterConfigExists("a")
	require.Equal(t, false, checked)
	require.Equal(t, errors.New("can not read 'cluster' from .blackbean").Error(), m.err.Error())
}
