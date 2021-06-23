// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toughnoah/blackbean/pkg/es"
	"log"
	"net/http"
	"os"
)

var (
	cfgFile string
)

func NewRootCmd(transport http.RoundTripper) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "blackbean",
		Short: "Basic interact with es via command line",
		Long: `blackbean command provides a set of commands to talk with es via cli.
Besides, blackbean is the name of my favorite french bulldog.`,
		ValidArgsFunction: noCompletions,
	}
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackbean.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// We can safely ignore any errors that flags.Parse encounters since
	// those errors will be caught later during the call to cmd.Execution.
	// This call is required to gather configuration information prior to
	// execution.
	flags.ParseErrorsWhitelist.UnknownFlags = true
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}
	InitConfig()
	url, username, password, err := es.GetProfile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	cli, err := es.NewEsClient(url, username, password, transport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		os.Exit(-1)
	}
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(catClusterResources(cli))
	rootCmd.AddCommand(applyClusterSettings(cli))
	rootCmd.AddCommand(getSnapshot(cli))
	rootCmd.AddCommand(useCluster())
	return rootCmd
}

func InitConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".blackbean" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".blackbean")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("ERROR: ", "can't not read config file! ", err)
	}
}
