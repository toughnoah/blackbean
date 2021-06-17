package cmd

import (
	"github.com/spf13/cobra"
	 "github.com/toughnoah/blackbean/pkg/es"
	"log"
)


func getClusterHealth() *cobra.Command {
	var command = &cobra.Command{
		Use:   "health [env]",
		Short: "get cluster health",
		Long:  "get cluster health ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return es.GetConfigEnv(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := es.NewEsClient(args[0])
			res, err := cli.Cat.Health(cli.Cat.Health.WithV(true))
			if err !=nil{
				return err
			}
			log.Print(res)
			return nil
		},
	}
	return command
}

func init()  {
	rootCmd.AddCommand(getClusterHealth())
}
