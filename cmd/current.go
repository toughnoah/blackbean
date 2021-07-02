package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
)

func current(out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:               "current-es",
		Short:             "show current cluster context",
		Long:              "show current cluster context ... wordless",
		Args:              cobra.NoArgs,
		ValidArgsFunction: noCompletions,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "current using cluster: %s\n\n", viper.Get("current"))
		},
	}
	return command
}
