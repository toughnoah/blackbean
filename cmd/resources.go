package cmd

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
)

var resources = []string{"health", "nodes", "allocations", "threadpool", "cachemem", "segmem", "largeindices", "allocationExp"}

func catClusterResources(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:   "cat [resource]",
		Short: "cat allocation/nodes/health/nodes/threadpool/cache memory/segments memory/large indices.",
		Long:  "cat es cluster info ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return resources, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := catResources(args[0], cli)
			if err == nil {
				fmt.Fprintf(out, "%s\n", res)
			}
			return err
		},
	}
	return command
}

func catResources(resource string, cli *elasticsearch.Client) (res *esapi.Response, err error) {
	return NewCatStrategy(resource, cli).Cat()
}

type CatStrategy struct {
	Strategy CatResource
}

func (c *CatStrategy) Cat() (res *esapi.Response, err error) {
	return c.Strategy.cat()
}

type CatResource interface {
	cat() (res *esapi.Response, err error)
}

type health struct {
	Client *elasticsearch.Client
}

func (o *health) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Health(o.Client.Cat.Health.WithV(true))
}

type nodes struct {
	Client *elasticsearch.Client
}

func (o *nodes) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(o.Client.Cat.Nodes.WithV(true))
}

type allocations struct {
	Client *elasticsearch.Client
}

func (o *allocations) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Allocation(o.Client.Cat.Allocation.WithV(true))
}

type threadpool struct {
	Client *elasticsearch.Client
}

func (o *threadpool) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.ThreadPool(o.Client.Cat.ThreadPool.WithV(true))
}

type cacheMemory struct {
	Client *elasticsearch.Client
}

func (o *cacheMemory) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(
		o.Client.Cat.Nodes.WithV(true),
		o.Client.Cat.Nodes.WithH("name", "queryCacheMemory", "queryCacheEvictions", "requestCacheMemory", "requestCacheHitCount", "request_cache.miss_count"),
	)
}

type segmentsMemory struct {
	Client *elasticsearch.Client
}

func (o *segmentsMemory) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(
		o.Client.Cat.Nodes.WithV(true),
		o.Client.Cat.Nodes.WithH("name", "segments.memory", "segments.index_writer_memory", "fielddata.memory_size", "query_cache.memory_size", "request_cache.memory_size"),
	)
}

type largeIndices struct {
	Client *elasticsearch.Client
}

func (o *largeIndices) cat() (res *esapi.Response, err error) {
	return o.Client.Cat.Indices(
		o.Client.Cat.Indices.WithV(true),
		o.Client.Cat.Indices.WithH("store.size", "index"),
		o.Client.Cat.Indices.WithBytes("gb"),
	)
}

func NewCatStrategy(resource string, cli *elasticsearch.Client) *CatStrategy {
	strategy := new(CatStrategy)
	switch resource {
	case "health":
		strategy.Strategy = &health{Client: cli}
	case "nodes":
		strategy.Strategy = &nodes{Client: cli}
	case "allocations":
		strategy.Strategy = &allocations{Client: cli}
	case "threadpool":
		strategy.Strategy = &threadpool{Client: cli}
	case "cachemem":
		strategy.Strategy = &cacheMemory{Client: cli}
	case "segmem":
		strategy.Strategy = &segmentsMemory{Client: cli}
	case "largeindices":
		strategy.Strategy = &largeIndices{Client: cli}
	}
	return strategy
}
