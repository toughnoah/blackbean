package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"log"
	"strings"
)

const (
	AllocateReplicasOps = "allocate_replica"
	CancelOps           = "cancel"
)

func reroute(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "reroute [subcommand]",
			Short:             "reroute for cluster",
			Long:              "reroute for cluster ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
		}
	)
	command.AddCommand(rerouteMoveIndex(cli, out, osArgs))
	command.AddCommand(cancel(cli, out, osArgs))
	command.AddCommand(rerouteAllocateReplicas(cli, out, osArgs))
	command.AddCommand(failed(cli, out))
	return command
}

func rerouteMoveIndex(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var (
		shard    string
		fromNode string
		toNode   string
		r        = rerouteObject{Client: cli}
		i        = Indices{client: cli}
		req      = new(es.RequestBody)
		command  = &cobra.Command{
			Use:   "move [index]",
			Short: "move index",
			Long:  "move index ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := r.move(args[0], shard, fromNode, toNode, req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&shard, "shard", "", "the shard of the index to be moved")
	f.StringVar(&fromNode, "from_node", "", "from which node the index should be moved")
	f.StringVar(&toNode, "to_node", "", "to which node the index should be moved")
	if err := command.RegisterFlagCompletionFunc("from_node", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllNodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	if err := command.RegisterFlagCompletionFunc("to_node", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllNodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	f = es.AddRequestBodyFlag(command, req)
	f.Parse(osArgs)
	if es.NoRawRequestBodySet(command) {
		_ = command.MarkFlagRequired("shard")
		_ = command.MarkFlagRequired("from_node")
		_ = command.MarkFlagRequired("to_node")
	}
	return command
}

func rerouteAllocateReplicas(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var (
		shard   string
		node    string
		r       = rerouteObject{Client: cli}
		i       = Indices{client: cli}
		req     = new(es.RequestBody)
		command = &cobra.Command{
			Use:   "allocateReplicas [index]",
			Short: "allocate replicas index",
			Long:  "allocate replicas index ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := r.allocateReplicaOrCancel(AllocateReplicasOps, args[0], shard, node, req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&shard, "shard", "", "the shard of the index to be moved")
	f.StringVar(&node, "node", "", "to which node the index should be moved")
	if err := command.RegisterFlagCompletionFunc("node", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllNodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	f = es.AddRequestBodyFlag(command, req)
	f.Parse(osArgs)
	if es.NoRawRequestBodySet(command) {
		_ = command.MarkFlagRequired("shard")
		_ = command.MarkFlagRequired("node")
	}
	return command
}

func cancel(cli *elasticsearch.Client, out io.Writer, osArgs []string) *cobra.Command {
	var (
		shard   string
		node    string
		r       = rerouteObject{Client: cli}
		i       = Indices{client: cli}
		req     = new(es.RequestBody)
		command = &cobra.Command{
			Use:   "cancel [index]",
			Short: "cancel allocating index",
			Long:  "cancel allocating index ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := r.allocateReplicaOrCancel(CancelOps, args[0], shard, node, req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&shard, "shard", "", "the shard of the index to be moved")
	f.StringVar(&node, "node", "", "to which node the index should be moved")
	if err := command.RegisterFlagCompletionFunc("node", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllNodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	f = es.AddRequestBodyFlag(command, req)
	f.Parse(osArgs)
	if es.NoRawRequestBodySet(command) {
		_ = command.MarkFlagRequired("shard")
		_ = command.MarkFlagRequired("node")
	}
	return command
}

func failed(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		r       = rerouteObject{Client: cli}
		command = &cobra.Command{
			Use:               "failed ",
			Short:             "retry failed allocation",
			Long:              "retry failed allocation ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := r.retryFailed()
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

type rerouteObject struct {
	Client *elasticsearch.Client
}

func (o *rerouteObject) move(index, shard, fromNode, toNode string, req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return o.Client.Cluster.Reroute(o.Client.Cluster.Reroute.WithBody(bytes.NewReader(rawBody)))
	}
	body := fmt.Sprintf(`{"commands": [{"move": {"index": "%s", "shard": %s,"from_node": "%s", "to_node": "%s"}}]}`,
		index, shard, fromNode, toNode)
	return o.Client.Cluster.Reroute(o.Client.Cluster.Reroute.WithBody(strings.NewReader(body)))
}

func (o *rerouteObject) allocateReplicaOrCancel(ops, index, shard, node string, req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return o.Client.Cluster.Reroute(o.Client.Cluster.Reroute.WithBody(bytes.NewReader(rawBody)))
	}
	body := fmt.Sprintf(`{"commands": [{ "%s": {"index": "%s", "shard": %s,"node": "%s"}}]}`, ops, index, shard, node)
	return o.Client.Cluster.Reroute(o.Client.Cluster.Reroute.WithBody(strings.NewReader(body)))
}

func (o *rerouteObject) retryFailed() (*esapi.Response, error) {
	return o.Client.Cluster.Reroute(o.Client.Cluster.Reroute.WithRetryFailed(true), o.Client.Cluster.Reroute.WithPretty())
}

func (o *rerouteObject) getAllNodes() []string {
	var (
		resMap   map[string]map[string]map[string]interface{}
		resSlice []string
	)

	ret, err := o.Client.Nodes.Stats()
	if err != nil {
		return nil
	}
	json.NewDecoder(ret.Body).Decode(&resMap)
	for _, v := range resMap["nodes"] {
		resSlice = append(resSlice, v["name"].(string))
	}
	return resSlice
}
