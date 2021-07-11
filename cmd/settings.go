package cmd

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
	if len(ClusterConcurrentRebalanced) == 0 && len(NodeConcurrentRecoveries) == 0 && len(NodeInitialPrimariesRecoveries) == 0 && len(Enable) == 0 {
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
	if len(BreakerFielddata) == 0 {
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
	if len(BreakerRequest) == 0 {
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
	if len(BreakerTotal) == 0 {
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
	if len(High) == 0 && len(Low) == 0 {
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
	if len(MaxBytesPerSec) == 0 {
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
	if len(MaxShardsPerNode) == 0 {
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
	if len(MaxCompilationsRate) == 0 {
		return S
	}

	S.Persistent.Script = &script{MaxCompilationsRate: MaxCompilationsRate}

	return S
}

func NewSettings() *Settings {
	return &Settings{Persistent: &persistent{}}
}
