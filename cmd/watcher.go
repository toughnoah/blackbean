package cmd

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
)

func watcher(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "watcher [subcommand]",
			Short:             "operate watcher",
			Long:              "operate watcherr ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
		}
	)
	command.AddCommand(watcherStart(cli, out))
	command.AddCommand(watcherStop(cli, out))
	command.AddCommand(watcherStats(cli, out))
	return command
}

func watcherStart(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		w       = Watcher{Client: cli}
		command = &cobra.Command{
			Use:               "start",
			Short:             "start watcher",
			Long:              "start watcher ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := w.start()
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return nil
			},
		}
	)
	return command
}

func watcherStop(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		w       = Watcher{Client: cli}
		command = &cobra.Command{
			Use:               "stop",
			Short:             "stop watcher",
			Long:              "stop watcher ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := w.stop()
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	return command
}
func watcherStats(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		w       = Watcher{Client: cli}
		command = &cobra.Command{
			Use:               "stats",
			Short:             "get watcher stats",
			Long:              "get watcher stats ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := w.stats()
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	return command
}

type Watcher struct {
	Client *elasticsearch.Client
}

func (w *Watcher) start() (*esapi.Response, error) {
	return w.Client.Watcher.Start()
}

func (w *Watcher) stop() (*esapi.Response, error) {
	return w.Client.Watcher.Stop()
}
func (w *Watcher) stats() (*esapi.Response, error) {
	return w.Client.Watcher.Stats(w.Client.Watcher.Stats.WithPretty())
}
