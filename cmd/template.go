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

func template(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:   "template [subcommand]",
		Short: "template operations",
		Long:  "template operations ... wordless",
	}
	command.AddCommand(getTemplate(cli, out))
	command.AddCommand(applyTemplate(cli, out))
	command.AddCommand(deleteTemplate(cli, out))
	return command
}

func getTemplate(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		t       = Template{Client: cli}
		command = &cobra.Command{
			Use:   "get [template]",
			Short: "get template",
			Long:  "get template ...wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return t.getAllTemplateName(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					res *esapi.Response
					err error
				)
				if len(args) == 0 {
					res, err = t.getAllTemplate()
				} else {
					res, err = t.getIndexTemplate(args[0])
				}
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

func applyTemplate(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		t       = Template{Client: cli}
		req     = &es.RequestBody{}
		command = &cobra.Command{
			Use:   "apply [template]",
			Short: "create or update template",
			Long:  "create or update template...wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return t.getAllTemplateName(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if es.NoRawRequestBodySet(cmd) {
					return es.NoRawRequestFlagError()
				}
				res, err := t.applyIndexTemplate(args[0], req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	es.AddRequestBodyFlag(command, req)
	return command
}

func deleteTemplate(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		t       = Template{Client: cli}
		command = &cobra.Command{
			Use:   "delete [template]",
			Short: "delete or update template",
			Long:  "delete or update template...wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return t.getAllTemplateName(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := t.deleteIndexTemplate(args[0])
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

type Template struct {
	Client *elasticsearch.Client
}

func (t *Template) getAllTemplate() (res *esapi.Response, err error) {
	return t.Client.Indices.GetTemplate(t.Client.Indices.GetTemplate.WithPretty())
}

func (t *Template) getIndexTemplate(name string) (res *esapi.Response, err error) {
	return t.Client.Indices.GetTemplate(t.Client.Indices.GetTemplate.WithName(name), t.Client.Indices.GetTemplate.WithPretty())
}

func (t *Template) applyIndexTemplate(name string, req *es.RequestBody) (res *esapi.Response, err error) {
	body, err := es.GetRawRequestBody(req)
	if err != nil {
		log.Println("failed to get raw request body")
		return nil, err
	}
	return t.Client.Indices.PutTemplate(name, bytes.NewReader(body))
}

func (t *Template) deleteIndexTemplate(name string) (res *esapi.Response, err error) {
	return t.Client.Indices.DeleteTemplate(name)
}

func (t *Template) getAllTemplateName() []string {
	var resMap map[string]interface{}
	var resSlice []string
	res, err := t.Client.Indices.GetTemplate()
	if err != nil {
		return nil
	}
	_ = json.NewDecoder(res.Body).Decode(&resMap)
	for k := range resMap {
		resSlice = append(resSlice, k)
	}
	return resSlice
}
