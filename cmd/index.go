package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/util"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
)

func Index(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:               "index [subcommand]",
		Short:             "index operations ",
		Long:              "index operations ... wordless",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: noCompletions,
	}
	command.AddCommand(getIndex(cli, out))
	command.AddCommand(searchIndex(cli, out))
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
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	return command
}

func searchIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		filename string
		data     string
	)
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
			res, err := i.queryIndex(args[0], data, filename)
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&filename, "filename", "f", "", "get query body from specific file.")
	f.StringVarP(&data, "data", "d", "{}", "specify query body")
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
	log.Println(res)
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
	res, err = i.client.Indices.Get(splitWords(indices), i.client.Indices.Get.WithPretty())
	return
}

func (i *Indices) queryIndex(index, data, filename string) (res *esapi.Response, err error) {
	var raw []byte
	if filename != "" {
		raw, err = decodeFromFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		raw = []byte(data)
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

func decodeFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	utf16bom := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	reader := transform.NewReader(f, utf16bom)
	raw := new(json.RawMessage)
	d := util.NewYAMLOrJSONDecoder(reader, 4096)
	if err = d.Decode(raw); err != nil {
		if err == io.EOF {
			return []byte(`{}`), nil
		}
		return nil, fmt.Errorf("error parsing %s: %v", filename, err)
	}
	return *raw, nil
}
