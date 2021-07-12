// Package cmd /*
package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"log"
	"net/http"
)

var (
	cfgFile string
)

func NewRootCmd(transport http.RoundTripper, out io.Writer, in io.ReadWriter, fd int, args []string) *cobra.Command {
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
	_ = flags.Parse(args)
	InitConfig()
	profile, err := es.GetProfile()
	if err != nil {
		log.Fatalf("Get Profile error %v", err)
	}

	cli, err := es.NewEsClient(profile.Info[es.ConfigUrl], profile.Info[es.ConfigUsername], profile.Info[es.ConfigPassword], transport)
	if err != nil {
		log.Fatalf("New client error %v", err)
	}
	rootCmd.AddCommand(NewCompletionCmd(out))
	rootCmd.AddCommand(catClusterResources(cli, out))
	rootCmd.AddCommand(apply(cli, out, args))
	rootCmd.AddCommand(snapshot(cli, out))
	rootCmd.AddCommand(repo(cli, out))
	rootCmd.AddCommand(useCluster(out))
	rootCmd.AddCommand(current(out))
	rootCmd.AddCommand(index(cli, out))
	rootCmd.AddCommand(alias(cli, out))
	rootCmd.AddCommand(reroute(cli, out, args))
	rootCmd.AddCommand(watcher(cli, out))
	rootCmd.AddCommand(explain(cli, out, args))
	rootCmd.AddCommand(user(cli, out, in, fd))
	rootCmd.AddCommand(role(cli, out))
	rootCmd.AddCommand(template(cli, out))
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
