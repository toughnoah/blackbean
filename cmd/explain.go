package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
)

func explain(cli *elasticsearch.Client, out io.Writer, args []string) *cobra.Command {
	var (
		shard   string
		node    string
		primary bool
		e       = Explain{Client: cli}
		i       = Indices{client: cli}
		command = &cobra.Command{
			Use:   "explain [index]",
			Short: "explain index allocation",
			Long:  "explain index allocation ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := e.allocationExplain(args, shard, node, primary)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&shard, "shard", "", "specifies the ID of the shard that you would like an explanation for.")
	f.StringVar(&node, "current_node", "", "specifies the node ID or the name of the node to only explain a shard that is currently located on the specified node.")
	f.BoolVar(&primary, "primary", false, "if true, returns explanation for the primary shard for the given shard ID.")
	if ArgSet(args) {
		_ = command.MarkFlagRequired("shard")
		_ = command.MarkFlagRequired("primary")
	}
	return command
}

type Explain struct {
	Client *elasticsearch.Client
}

type ExplainBody struct {
	Index   string `json:"index,omitempty"`
	Shard   string `json:"shard,omitempty"`
	Node    string `json:"current_node,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

func (e *Explain) allocationExplain(args []string, shard, node string, primary bool) (*esapi.Response, error) {
	if len(args) == 1 {
		explainBody := new(ExplainBody)
		explainBody.Index = args[0]
		explainBody.Node = node
		explainBody.Shard = shard
		explainBody.Primary = primary
		bytesBody, err := json.Marshal(&explainBody)
		if err != nil {
			return nil, err
		}
		return e.Client.Cluster.AllocationExplain(
			e.Client.Cluster.AllocationExplain.WithBody(bytes.NewReader(bytesBody)),
			e.Client.Cluster.AllocationExplain.WithPretty())
	}
	return e.Client.Cluster.AllocationExplain(e.Client.Cluster.AllocationExplain.WithPretty())
}

func ArgSet(args []string) bool {
	return len(args) >= 2
}
