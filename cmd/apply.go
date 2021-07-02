package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"log"
	"os"
)

var (
	clusterConcurrentRebalanced    string
	nodeConcurrentRecoveries       string
	nodeInitialPrimariesRecoveries string
	breakerFielddata               string
	breakerRequest                 string
	breakerTotal                   string
	watermarkHigh                  string
	watermarkLow                   string
	maxCompilationsRate            string
	maxShardsPerNode               string
	allocationEnable               string
	maxBytesPerSec                 string
)

func apply(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var command = &cobra.Command{
		Use:   "apply [subcommand]",
		Short: "apply cluster changes",
		Long:  "apply cluster changes ... wordless",
	}
	command.AddCommand(applySettings(cli, out, osArgs))
	command.AddCommand(applyFlushed(cli, out))
	command.AddCommand(applyClearCache(cli, out))
	return command
}

func applySettings(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var (
		req     = new(es.RequestBody)
		command = &cobra.Command{
			Use:               "settings --[flag] ",
			Short:             "apply cluster settings change",
			Long:              "apply cluster settings change ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				if AllFlagNotSet(osArgs) {
					return errors.New("At lease one flag should be specified to change cluster settings")
				}
				o := applyObject{Client: cli}
				res, err := o.putSettings(req)
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVarP(&clusterConcurrentRebalanced, "cluster_concurrent_rebalanced", "a", "", "to set cluster_concurrent_rebalanced value, such as 10 (Only for 'settings' resource)")
	f.StringVarP(&nodeConcurrentRecoveries, "node_concurrent_recoveries", "n", "", "to set node_concurrent_recoveries value, such as 10 (Only for 'settings' resource)")
	f.StringVarP(&nodeInitialPrimariesRecoveries, "node_initial_primaries_recoveries", "i", "", "to set node_initial_primaries_recoveries value, such as 10 (Only for 'settings' resource)")
	f.StringVarP(&breakerFielddata, "breaker_fielddata", "k", "", "to set breaker_fielddata value, such as 60% (Only for 'settings' resource)")
	f.StringVarP(&breakerRequest, "breaker_request", "r", "", "to set breaker_request value, such as 60% (Only for 'settings' resource)")
	f.StringVarP(&breakerTotal, "breaker_total", "t", "", "to set breaker_total value, such as 60% (Only for 'settings' resource)")
	f.StringVarP(&watermarkHigh, "watermark_high", "w", "", "to set watermark_high value, such as 85% (Only for 'settings' resource)")
	f.StringVarP(&watermarkLow, "watermark_low", "l", "", "to set watermark_low value, such as 85% (Only for 'settings' resource)")
	f.StringVarP(&maxCompilationsRate, "max_compilations_rate", "m", "", "to set max_compilations_rate value, such as 75/5m (Only for 'settings' resource)")
	f.StringVarP(&maxShardsPerNode, "max_shards_per_node", "s", "", "to set max_shards_per_node value, such as 1000 (Only for 'settings' resource)")
	f.StringVarP(&allocationEnable, "allocation_enable", "e", "", "to set allocation enable value, primaries or null (Only for 'settings' resource)")
	f.StringVarP(&maxBytesPerSec, "max_bytes_per_sec", "b", "", "to set indices recovery max_bytes_per_sec, default 40 (Only for 'settings' resource)")

	err := command.RegisterFlagCompletionFunc("allocation_enable", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"primaries", "null"}, cobra.ShellCompDirectiveNoFileComp
	})

	if err != nil {
		log.Fatal(err)
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func applyFlushed(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "flush",
			Short:             "apply indices to flush",
			Long:              "apply indices to flush ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				o := applyObject{Client: cli}
				res, err := o.flush()
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	return command
}

func applyClearCache(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "clearCache",
			Short:             "apply indices to clear cache",
			Long:              "apply indices to clear cache ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				o := applyObject{Client: cli}
				res, err := o.clearCache()
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	return command
}

type applyObject struct {
	Client *elasticsearch.Client
}

func (o *applyObject) putSettings(req *es.RequestBody) (*esapi.Response, error) {
	data, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if data == nil {
		settings := NewSettings()
		settings.WithAllocationSettings(
			clusterConcurrentRebalanced,
			nodeConcurrentRecoveries,
			nodeInitialPrimariesRecoveries,
			allocationEnable).
			WithBreakerRequest(breakerRequest).
			WithBreakerTotal(breakerTotal).
			WithBreakerRequest(breakerRequest).
			WithWatermark(watermarkHigh, watermarkLow).
			WithRecovery(maxBytesPerSec).
			WithMaxShardsPerNode(maxShardsPerNode).
			WithMaxCompilationsRate(maxCompilationsRate)
		data, err = json.Marshal(settings)
		if err != nil {
			return nil, errors.Wrap(err, "failed to Marshal settings")
		}
	}
	putSettings, err := o.Client.Cluster.PutSettings(bytes.NewReader(data), o.Client.Cluster.PutSettings.WithPretty())
	if err != nil {
		return nil, errors.Wrap(err, "failed when sending put request")
	}
	return putSettings, err
}

func (o *applyObject) flush() (*esapi.Response, error) {
	return o.Client.Indices.FlushSynced(
		o.Client.Indices.FlushSynced.WithIgnoreUnavailable(true),
		o.Client.Indices.FlushSynced.WithPretty())
}

func (o *applyObject) clearCache() (*esapi.Response, error) {
	return o.Client.Indices.ClearCache()
}

func AllFlagNotSet(args []string) bool {
	return len(args) == 4 && args[len(os.Args)-1] == "settings"
}
