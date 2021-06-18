package cmd

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"log"
	"os"
)

var resources = []string{"health", "nodes", "allocations", "threadpool", "cachemem", "segmem", "largeindices"}

func catClusterResources() *cobra.Command {
	var Cluster string
	var command = &cobra.Command{
		Use:   "get [resource]",
		Short: "get allocation/nodes/health/nodes/threadpool/cache memory/segments memory/large indices.",
		Long:  "get es cluster info ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return resources, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := es.NewEsClient(Cluster)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
				os.Exit(-1)
			}
			co := &CatObject{
				Client:   cli,
				Resource: args[0],
			}
			if err = co.catResources(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
				os.Exit(-1)
			}
		},
	}
	f := command.Flags()
	f.StringVarP(&Cluster, "cluster", "c", "default", "to specify a es cluster")

	err := command.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return es.GetConfigEnv(toComplete), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	})

	if err != nil {
		log.Fatal(err)
	}
	_ = rootCmd.MarkFlagRequired("cluster")
	return command
}

type CatObject struct {
	Client   *elasticsearch.Client
	Resource string
}

func (o *CatObject) catResources() (err error) {

	var res *esapi.Response
	switch o.Resource {
	case "health":
		res, err = o.catHealth()
	case "nodes":
		res, err = o.catNodes()
	case "allocations":
		res, err = o.catAllocation()
	case "threadpool":
		res, err = o.catThreadpool()
	case "cachemem":
		res, err = o.catCacheMemory()
	case "segmem":
		res, err = o.catSegmentsMemory()
	case "largeindices":
		res, err = o.catLargeIndices()
	}
	fmt.Println(res)
	return
}

func (o *CatObject) catHealth() (res *esapi.Response, err error) {
	return o.Client.Cat.Health(o.Client.Cat.Health.WithV(true))
}

func (o *CatObject) catNodes() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(o.Client.Cat.Nodes.WithV(true))
}

func (o *CatObject) catAllocation() (res *esapi.Response, err error) {
	return o.Client.Cat.Allocation(o.Client.Cat.Allocation.WithV(true))
}

func (o *CatObject) catThreadpool() (res *esapi.Response, err error) {
	return o.Client.Cat.ThreadPool(o.Client.Cat.ThreadPool.WithV(true))
}
func (o *CatObject) catCacheMemory() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(
		o.Client.Cat.Nodes.WithV(true),
		o.Client.Cat.Nodes.WithH("name", "queryCacheMemory", "queryCacheEvictions", "requestCacheMemory", "requestCacheHitCount", "request_cache.miss_count"),
	)
}

func (o *CatObject) catSegmentsMemory() (res *esapi.Response, err error) {
	return o.Client.Cat.Nodes(
		o.Client.Cat.Nodes.WithV(true),
		o.Client.Cat.Nodes.WithH("name", "segments.memory", "segments.index_writer_memory", "fielddata.memory_size", "query_cache.memory_size", "request_cache.memory_size"),
	)
}
func (o *CatObject) catLargeIndices() (res *esapi.Response, err error) {
	return o.Client.Cat.Indices(
		o.Client.Cat.Indices.WithV(true),
		o.Client.Cat.Indices.WithH("store.size", "index"),
		o.Client.Cat.Indices.WithBytes("gb"),
	)
}

func init() {
	rootCmd.AddCommand(catClusterResources())
}
