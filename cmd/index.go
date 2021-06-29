package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
	"log"
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

type Indices struct {
	client *elasticsearch.Client
}

func (i *Indices) getAllIndices() []string {
	var (
		indicesMap   map[string]interface{}
		indicesSlice []string
	)
	log.Println("aaa")
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
