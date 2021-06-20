package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"os"
)

var editableResources = []string{
	"settings",
}

type Settings struct {
	Persistent *persistent `json:"persistent"`
}

type persistent struct {
	Cluster *cluster `json:"cluster,omitempty"`
	Script  *script  `json:"script,omitempty"`
	Indices *indices `json:"indices,omitempty"`
}

type cluster struct {
	Routing          *routing `json:"routing,omitempty"`
	MaxShardsPerNode string   `json:"max_shards_per_node,omitempty"`
}

type script struct {
	MaxCompilationsRate string `json:"max_compilations_rate,omitempty"`
}

type indices struct {
	Breaker  *breaker  `json:"breaker,omitempty"`
	Recovery *recovery `json:"recovery,omitempty"`
}

type recovery struct {
	MaxBytesPerSec string `json:"max_bytes_per_sec,omitempty"`
}

type routing struct {
	Allocation *allocation `json:"allocation,omitempty"`
}

type allocation struct {
	ClusterConcurrentRebalanced    string `json:"cluster_concurrent_rebalance,omitempty"`
	NodeConcurrentRecoveries       string `json:"node_concurrent_recoveries,omitempty"`
	NodeInitialPrimariesRecoveries string `json:"node_initial_primaries_recoveries,omitempty"`
	Disk                           *disk  `json:"disk,omitempty"`
	Enable                         string `json:"enable,omitempty"`
}

type disk struct {
	Watermark *watermark `json:"watermark,omitempty"`
}

type watermark struct {
	High string `json:"high,omitempty"`
	Low  string `json:"low,omitempty"`
}

type breaker struct {
	Fielddata *fielddata `json:"fielddata,omitempty"`
	Request   *request   `json:"request,omitempty"`
	Total     *total     `json:"total,omitempty"`
}

type fielddata struct {
	Limit string `json:"limit,omitempty"`
}

type request struct {
	Limit string `json:"limit,omitempty"`
}

type total struct {
	Limit string `json:"limit,omitempty"`
}

func (S *Settings) WithAllocationSettings(ClusterConcurrentRebalanced, NodeConcurrentRecoveries, NodeInitialPrimariesRecoveries, Enable string) *Settings {
	if ClusterConcurrentRebalanced == "" && NodeConcurrentRecoveries == "" && NodeInitialPrimariesRecoveries == "" && Enable == "" {
		return S
	}
	if S.Persistent.Cluster != nil && S.Persistent.Cluster.Routing != nil {
		if S.Persistent.Cluster.Routing.Allocation != nil {
			S.Persistent.Cluster.Routing.Allocation.NodeInitialPrimariesRecoveries = NodeInitialPrimariesRecoveries
			S.Persistent.Cluster.Routing.Allocation.ClusterConcurrentRebalanced = ClusterConcurrentRebalanced
			S.Persistent.Cluster.Routing.Allocation.NodeConcurrentRecoveries = NodeConcurrentRecoveries
			S.Persistent.Cluster.Routing.Allocation.Enable = Enable
		} else {
			S.Persistent.Cluster.Routing.Allocation = &allocation{
				ClusterConcurrentRebalanced:    ClusterConcurrentRebalanced,
				NodeConcurrentRecoveries:       NodeConcurrentRecoveries,
				NodeInitialPrimariesRecoveries: NodeInitialPrimariesRecoveries,
				Enable:                         Enable,
			}
		}
	} else {
		S.Persistent.Cluster = &cluster{
			Routing: &routing{
				Allocation: &allocation{
					ClusterConcurrentRebalanced:    ClusterConcurrentRebalanced,
					NodeConcurrentRecoveries:       NodeConcurrentRecoveries,
					NodeInitialPrimariesRecoveries: NodeInitialPrimariesRecoveries,
					Enable:                         Enable,
				},
			},
		}
	}
	return S
}

func (S *Settings) WithBreakerFielddata(BreakerFielddata string) *Settings {
	if BreakerFielddata == "" {
		return S
	}

	if S.Persistent.Indices != nil {
		S.Persistent.Indices.Breaker.Fielddata = &fielddata{
			Limit: BreakerFielddata,
		}
	} else {
		S.Persistent.Indices = &indices{
			Breaker: &breaker{
				Fielddata: &fielddata{
					Limit: BreakerFielddata,
				},
			},
		}
	}

	return S
}

func (S *Settings) WithBreakerRequest(BreakerRequest string) *Settings {
	if BreakerRequest == "" {
		return S
	}

	if S.Persistent.Indices != nil {
		S.Persistent.Indices.Breaker.Request = &request{
			Limit: BreakerRequest,
		}
	} else {
		S.Persistent.Indices = &indices{
			Breaker: &breaker{
				Request: &request{
					Limit: BreakerRequest,
				},
			},
		}
	}

	return S
}

func (S *Settings) WithBreakerTotal(BreakerTotal string) *Settings {
	if BreakerTotal == "" {
		return S
	}

	if S.Persistent.Indices != nil {
		S.Persistent.Indices.Breaker.Total = &total{
			Limit: BreakerTotal,
		}
	} else {
		S.Persistent.Indices = &indices{
			Breaker: &breaker{
				Total: &total{
					Limit: BreakerTotal,
				},
			},
		}
	}

	return S
}

func (S *Settings) WithWatermark(High, Low string) *Settings {
	if High == "" && Low == "" {
		return S
	}
	if S.Persistent.Cluster != nil && S.Persistent.Cluster.Routing != nil {
		S.Persistent.Cluster.Routing.Allocation.Disk = &disk{
			Watermark: &watermark{
				Low:  Low,
				High: High,
			},
		}
	} else {
		S.Persistent.Cluster = &cluster{
			Routing: &routing{
				Allocation: &allocation{
					Disk: &disk{
						Watermark: &watermark{
							Low:  Low,
							High: High,
						},
					},
				},
			},
		}
	}
	return S
}

func (S *Settings) WithRecovery(MaxBytesPerSec string) *Settings {
	if MaxBytesPerSec == "" {
		return S
	}
	if S.Persistent.Indices != nil {
		S.Persistent.Indices.Recovery = &recovery{
			MaxBytesPerSec: MaxBytesPerSec,
		}
	} else {
		S.Persistent.Indices = &indices{
			Recovery: &recovery{
				MaxBytesPerSec: MaxBytesPerSec,
			},
		}
	}

	return S
}

func (S *Settings) WithMaxShardsPerNode(MaxShardsPerNode string) *Settings {
	if MaxShardsPerNode == "" {
		return S
	}
	if S.Persistent.Cluster != nil {
		S.Persistent.Cluster.MaxShardsPerNode = MaxShardsPerNode
	} else {
		S.Persistent.Cluster = &cluster{
			MaxShardsPerNode: MaxShardsPerNode,
		}
	}

	return S
}

func (S *Settings) WithMaxCompilationsRate(MaxCompilationsRate string) *Settings {
	if MaxCompilationsRate == "" {
		return S
	}

	S.Persistent.Script = &script{MaxCompilationsRate: MaxCompilationsRate}

	return S
}

func NewSettings() *Settings {
	return &Settings{Persistent: &persistent{}}
}

func applyClusterSettings(cli *elasticsearch.Client) *cobra.Command {
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
	var command = &cobra.Command{
		Use:   "apply [resource]",
		Short: "apply cluster settings.",
		Long:  "apply cluster settings ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return editableResources, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "settings" {
				return es.NoResourcesError(args[0])
			}
			po := PutObject{Client: cli}
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

			res, err := po.putSettings(settings)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	}
	f := command.Flags()
	f.StringVarP(&clusterConcurrentRebalanced, "cluster_concurrent_rebalanced", "a", "", "to set cluster_concurrent_rebalanced value, such as 10")
	f.StringVarP(&nodeConcurrentRecoveries, "node_concurrent_recoveries", "n", "", "to set node_concurrent_recoveries value, such as 10")
	f.StringVarP(&nodeInitialPrimariesRecoveries, "node_initial_primaries_recoveries", "i", "", "to set node_initial_primaries_recoveries value, such as 10")
	f.StringVarP(&breakerFielddata, "breaker_fielddata", "f", "", "to set breaker_fielddata value, such as 10")
	f.StringVarP(&breakerRequest, "breaker_request", "r", "", "to set breaker_request value, such as 10")
	f.StringVarP(&breakerTotal, "breaker_total", "t", "", "to set breaker_total value, such as 10")
	f.StringVarP(&watermarkHigh, "watermark_high", "w", "", "to set watermark_high value, such as 10")
	f.StringVarP(&watermarkLow, "watermark_low", "l", "", "to set watermark_low value, such as 10")
	f.StringVarP(&maxCompilationsRate, "max_compilations_rate", "m", "", "to set max_compilations_rate value, such as 75/5")
	f.StringVarP(&maxShardsPerNode, "max_shards_per_node", "s", "", "to set max_shards_per_node value, such as 1000")
	f.StringVarP(&allocationEnable, "allocation_enable", "e", "", "to set allocation enable value, primaries or null")
	f.StringVarP(&maxBytesPerSec, "max_bytes_per_sec", "b", "", "to set indices recovery max_bytes_per_sec, default 40 ")

	err := command.RegisterFlagCompletionFunc("allocation_enable", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return []string{"primaries", "null"}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		os.Exit(-1)
	}
	return command
}

type PutObject struct {
	Client *elasticsearch.Client
}

func (o *PutObject) putSettings(settings *Settings) (*esapi.Response, error) {
	data, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	putSettings, err := o.Client.Cluster.PutSettings(reader, o.Client.Cluster.PutSettings.WithPretty())
	if err != nil {
		return nil, err
	}
	return putSettings, err
}
