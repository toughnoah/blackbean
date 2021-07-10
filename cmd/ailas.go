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
)

func alias(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:   "alias",
		Short: "alias index",
		Long:  "alias index ... wordless",
	}
	command.AddCommand(createAlias(cli, out))
	command.AddCommand(getAlias(cli, out))
	command.AddCommand(deleteAlias(cli, out))
	return command
}

func createAlias(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		req     = new(es.RequestBody)
		i       = Indices{client: cli}
		a       = Alias{client: cli}
		command = &cobra.Command{
			Use:   "create [index] [alias]",
			Short: "create alias for index",
			Long:  "create alias for index ... wordless",
			Args:  cobra.ExactArgs(2),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := a.createAlias(args[0], args[1], req)
				fmt.Fprintln(out, res)
				return err
			},
		}
	)
	es.AddRequestBodyFlag(command, req)
	return command
}

func getAlias(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		i       = Indices{client: cli}
		a       = Alias{client: cli}
		isAlias bool
		command = &cobra.Command{
			Use:   "get [index/alias]",
			Short: "get alias for index or get alias list",
			Long:  "get alias for index or get alias list",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				res := i.getAllIndices()
				res = append(res, a.getAllAlias()...)
				return res, cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := a.getAlias(args[0], isAlias)
				fmt.Fprintln(out, res)
				return err
			},
		}
	)
	f := command.Flags()
	f.BoolVar(&isAlias, "is_alias", false, "to specify the args is alias.")
	return command
}

func deleteAlias(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		i       = Indices{client: cli}
		a       = Alias{client: cli}
		command = &cobra.Command{
			Use:   "delete [index] [alias]",
			Short: "delete alias for index",
			Long:  "delete alias for index ... wordless",
			Args:  cobra.ExactArgs(2),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) == 1 {
					return a.getAllAlias(), cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := a.deleteAlias(args[0], args[1])
				fmt.Fprintln(out, res)
				return err
			},
		}
	)
	return command
}

type Alias struct {
	client *elasticsearch.Client
}

func (a *Alias) createAlias(indices, alias string, req *es.RequestBody) (res *esapi.Response, err error) {
	body, err := es.GetRawRequestBody(req)
	if err != nil {
		log.Println("failed to get raw request body")
		return nil, err
	}
	return a.client.Indices.PutAlias(splitWords(indices), alias, a.client.Indices.PutAlias.WithBody(bytes.NewReader(body)))
}

func (a *Alias) getAlias(indicesOrAlias string, isAlias bool) (res *esapi.Response, err error) {
	if isAlias {
		return a.client.Indices.GetAlias(a.client.Indices.GetAlias.WithName(splitWords(indicesOrAlias)...))
	}
	return a.client.Indices.GetAlias(a.client.Indices.GetAlias.WithIndex(splitWords(indicesOrAlias)...))
}

func (a *Alias) deleteAlias(indices, name string) (res *esapi.Response, err error) {
	return a.client.Indices.DeleteAlias(splitWords(indices), splitWords(name))
}

func (a *Alias) getAllAlias() []string {
	var (
		resMap   map[string]map[string]interface{}
		resSlice []string
	)
	res, err := a.client.Indices.GetAlias()
	if err != nil {
		log.Printf("error sending request to es: %s", err)
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Printf("error parsing the response body: %s", err)
		return nil
	}
	for _, aliasMap := range resMap {
		for aliasName, _ := range aliasMap["aliases"].(map[string]interface{}) {
			resSlice = append(resSlice, aliasName)
		}
	}
	return resSlice
}
