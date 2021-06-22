package cmd

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCurrent(t *testing.T) {
	r := bytes.NewReader(yamlExample)
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(r)
	require.NoError(t, err)
	_, err = executeCommand("current", nil)
	require.NoError(t, err)
}
