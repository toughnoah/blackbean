package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUse(t *testing.T) {
	fs := afero.NewOsFs()
	filename := ".blackbean.yaml"
	home, err := homedir.Dir()
	file := filepath.Join(home, filename)
	err = ioutil.WriteFile(file, yamlExample, 0755)
	require.NoError(t, err)
	InitConfig()
	out, err := executeCommand("use backup", nil)
	require.NoError(t, err)
	require.Equal(t, out, "change to  cluster: backup\n\n")
	defer fs.Remove(file)
}
