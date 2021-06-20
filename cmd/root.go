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
	Cluster string
)

func NewRootCmd(transport http.RoundTripper) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "blackbean",
		Short: "Basic interact with es via command line",
		Long: `blackbean command provides a set of commands to talk with es via cli.
Besides, blackbean is the name of my favorite french bulldog.`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return editableResources, cobra.ShellCompDirectiveNoFileComp
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackbean.yaml)")

	rootCmd.PersistentFlags().StringVarP(&Cluster, "cluster", "c", "default", "to specify a es cluster")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	err := rootCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return es.CompleteConfigEnv(toComplete), cobra.ShellCompDirectiveNoFileComp
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		os.Exit(-1)
	}
	url, username, password, err := es.GetEnv(Cluster)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		os.Exit(-1)
	}

	cli, err := es.NewEsClient(url, username, password, transport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		os.Exit(-1)
	}
	rootCmd.AddCommand(catClusterResources(cli))
	rootCmd.AddCommand(applyClusterSettings(cli))
	rootCmd.AddCommand(completionCmd)
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
