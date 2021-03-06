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
		req     = &es.RequestBody{}
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
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&clusterConcurrentRebalanced, "cluster_concurrent_rebalanced", "", "to set cluster_concurrent_rebalanced value, such as 10 (Only for 'settings' resource)")
	f.StringVar(&nodeConcurrentRecoveries, "node_concurrent_recoveries", "", "to set node_concurrent_recoveries value, such as 10 (Only for 'settings' resource)")
	f.StringVar(&nodeInitialPrimariesRecoveries, "node_initial_primaries_recoveries", "", "to set node_initial_primaries_recoveries value, such as 10 (Only for 'settings' resource)")
	f.StringVar(&breakerFielddata, "breaker_fielddata", "", "to set breaker_fielddata value, such as 60% (Only for 'settings' resource)")
	f.StringVar(&breakerRequest, "breaker_request", "", "to set breaker_request value, such as 60% (Only for 'settings' resource)")
	f.StringVar(&breakerTotal, "breaker_total", "", "to set breaker_total value, such as 60% (Only for 'settings' resource)")
	f.StringVar(&watermarkHigh, "watermark_high", "", "to set watermark_high value, such as 85% (Only for 'settings' resource)")
	f.StringVar(&watermarkLow, "watermark_low", "", "to set watermark_low value, such as 85% (Only for 'settings' resource)")
	f.StringVar(&maxCompilationsRate, "max_compilations_rate", "", "to set max_compilations_rate value, such as 75/5m (Only for 'settings' resource)")
	f.StringVar(&maxShardsPerNode, "max_shards_per_node", "", "to set max_shards_per_node value, such as 1000 (Only for 'settings' resource)")
	f.StringVar(&allocationEnable, "allocation_enable", "", "to set allocation enable value, primaries or null (Only for 'settings' resource)")
	f.StringVar(&maxBytesPerSec, "max_bytes_per_sec", "", "to set indices recovery max_bytes_per_sec, default 40 (Only for 'settings' resource)")

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
					fmt.Fprintln(out, res)
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
					fmt.Fprintln(out, res)
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
	return len(args) == 4 && args[len(args)-1] == "settings"
}
