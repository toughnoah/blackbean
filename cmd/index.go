package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"log"
)

func index(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:   "index [subcommand]",
		Short: "index operations ",
		Long:  "index operations ... wordless",
	}
	command.AddCommand(getIndex(cli, out))
	command.AddCommand(searchIndex(cli, out))
	command.AddCommand(createIndex(cli, out))
	command.AddCommand(deleteIndex(cli, out))
	command.AddCommand(reindex(cli, out))
	return command
}

func getIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	i := Indices{client: cli}
	var command = &cobra.Command{
		Use:   "get [index]",
		Short: "get index from cluster",
		Long:  "get index from cluster ... wordless",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := i.getIndices(args[0])
			if err == nil {
				fmt.Fprintf(out, "%s\n", res)
			}
			return err
		},
	}
	return command
}

func searchIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	req := new(es.RequestBody)
	i := Indices{client: cli}
	var command = &cobra.Command{
		Use:   "search [index]",
		Short: "search index from cluster",
		Long:  "search index from cluster ... wordless",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := i.searchIndex(args[0], req)
			if err == nil {
				fmt.Fprintf(out, "%s\n", res)
			}
			return err
		},
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func createIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	req := new(es.RequestBody)
	i := Indices{client: cli}
	var command = &cobra.Command{
		Use:               "create [index]",
		Short:             "create index from command",
		Long:              "create index from command ... wordless",
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := i.createIndex(args[0], req)
			if err == nil {
				fmt.Fprintf(out, "%s\n", res)
			}
			return err
		},
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func deleteIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	i := Indices{client: cli}
	var command = &cobra.Command{
		Use:   "delete [index]",
		Short: "delete index from command",
		Long:  "delete index from command ... wordless",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := i.deleteIndices(args[0])
			if err == nil {
				fmt.Fprintf(out, "%s\n", res)
			}
			return err
		},
	}
	return command
}

func reindex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	req := new(es.RequestBody)
	i := Indices{client: cli}
	var (
		command = &cobra.Command{
			Use:   "reindex [index] [newIndex]",
			Short: "do reindex",
			Long:  "do reindex ... wordless",
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			Args: cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := i.reIndex(args[0], args[1], req)
				if err == nil {
					fmt.Fprintf(out, "%s\n", res)
				}
				return err
			},
		}
	)
	es.AddRequestBodyFlag(command, req)
	return command
}

type Indices struct {
	client *elasticsearch.Client
}

func (i *Indices) getAllIndices() []string {
	var (
		indicesMap   map[string]interface{}
		indicesSlice []string
	)
	res, err := i.client.Indices.Get([]string{"_all"}, i.client.Indices.Get.WithPretty())
	if err != nil {
		log.Printf("error sending request to es: %s", err)
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&indicesMap); err != nil {
		log.Printf("error parsing the response body: %s", err)
		return nil
	}
	for index, _ := range indicesMap {
		indicesSlice = append(indicesSlice, index)
	}
	return indicesSlice
}

func (i *Indices) getIndices(indices string) (res *esapi.Response, err error) {
	return i.client.Indices.Get(splitWords(indices), i.client.Indices.Get.WithPretty())
}

func (i *Indices) deleteIndices(indices string) (res *esapi.Response, err error) {
	return i.client.Indices.Delete(splitWords(indices), i.client.Indices.Delete.WithIgnoreUnavailable(true))
}

func (i *Indices) searchIndex(index string, req *es.RequestBody) (res *esapi.Response, err error) {
	var raw []byte
	if req.Filename != es.EmptyFile {
		raw, err = es.DecodeFromFile(req.Filename)
		if err != nil {
			return nil, err
		}
	} else {
		raw = []byte(req.Data)
	}
	res, err = i.client.Search(
		i.client.Search.WithContext(context.Background()),
		i.client.Search.WithIndex(index),
		i.client.Search.WithBody(bytes.NewReader(raw)),
		i.client.Search.WithTrackTotalHits(true),
		i.client.Search.WithPretty(),
	)
	return
}

func (i *Indices) createIndex(index string, req *es.RequestBody) (res *esapi.Response, err error) {
	body, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, nil
	}
	return i.client.Indices.Create(index, i.client.Indices.Create.WithBody(bytes.NewReader(body)))
}

func (i *Indices) reIndex(source, dest string, req *es.RequestBody) (res *esapi.Response, err error) {
	body, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if body == nil {
		body = []byte(fmt.Sprintf(`{"source":{"index":"%s"}, "dest":{"index": "%s"}}`, source, dest))
	}
	return i.doReindex(bytes.NewReader(body))
}

func (i *Indices) doReindex(body io.Reader) (res *esapi.Response, err error) {
	res, err = i.client.Reindex(body, i.client.Reindex.WithWaitForCompletion(false))
	return
}
