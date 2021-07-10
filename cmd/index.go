package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"io/ioutil"
	"log"
	"strings"
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
	command.AddCommand(writeIndex(cli, out))
	command.AddCommand(bulk(cli, out))
	command.AddCommand(mSearch(cli, out))
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
				fmt.Fprintln(out, res)
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
				fmt.Fprintln(out, res)
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
				fmt.Fprintln(out, res)
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
				fmt.Fprintln(out, res)
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
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	es.AddRequestBodyFlag(command, req)
	return command
}

func writeIndex(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	req := new(es.RequestBody)
	i := Indices{client: cli}
	var command = &cobra.Command{
		Use:   "write [index]",
		Short: "write index from command",
		Long:  "write index from command ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if es.NoRawRequestBodySet(cmd) {
				return es.NoRawRequestFlagError()
			}
			res, err := i.writeIndex(args[0], req)
			if err == nil {
				fmt.Fprintln(out, res)
			}
			return err
		},
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func bulk(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		i            = Indices{client: cli}
		requireAlias bool
		pipeline     string
		rawFile      string
		data         string
		command      = &cobra.Command{
			Use:   "bulk",
			Short: "send bulk request",
			Long:  "send bulk request ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				i.rawFile = rawFile
				fmt.Println(i.rawFile)
				if len(args) != 0 {
					i.index = args[0]
				}
				if es.GetFlagValue(cmd, "data") == es.EmptyData && es.GetFlagValue(cmd, "raw_file") == es.EmptyFile {
					return errors.New("one of --data and --raw_file should be specified")
				}
				res, err := i.bulk(requireAlias, pipeline)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&pipeline, "pipeline", "", "ID of the pipeline to use to preprocess incoming documents.")
	f.BoolVar(&requireAlias, "require_alias", false, "if true, the request’s actions must target an index alias.")
	f.StringVar(&rawFile, "raw_file", es.EmptyFile, "the path to raw file with request body")
	f.StringVarP(&data, "data", "d", es.EmptyData, "the path to raw file with request body")
	return command
}

func mSearch(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		i                          = Indices{client: cli}
		maxConcurrentSearches      int
		maxConcurrentShardRequests int
		rawFile                    string
		data                       string
		command                    = &cobra.Command{
			Use:   "msearch",
			Short: "send msearch request",
			Long:  "send msearch request ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				i.rawFile = rawFile
				if len(args) != 0 {
					i.index = args[0]
				}
				if es.GetFlagValue(cmd, "data") == es.EmptyData && es.GetFlagValue(cmd, "raw_file") == es.EmptyFile {
					return errors.New("one of --data and --raw_file should be specified")
				}
				res, err := i.msearch(maxConcurrentSearches, maxConcurrentShardRequests)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.IntVar(&maxConcurrentSearches, "max_concurrent_searches", 0, "ID of the pipeline to use to preprocess incoming documents.")
	f.IntVar(&maxConcurrentShardRequests, "max_concurrent_shard_requests", 0, "if true, the request’s actions must target an index alias.")
	f.StringVar(&rawFile, "raw_file", es.EmptyFile, "the path to raw file with request body")
	f.StringVarP(&data, "data", "d", es.EmptyData, "the path to raw file with request body")
	return command
}

type Indices struct {
	client  *elasticsearch.Client
	index   string
	rawFile string
	data    string
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

func (i *Indices) writeIndex(index string, req *es.RequestBody) (res *esapi.Response, err error) {
	body, err := es.GetRawRequestBody(req)
	if err != nil {
		log.Println("failed to get raw request body")
		return nil, err
	}
	return i.client.Index(index, bytes.NewReader(body), i.client.Index.WithPretty())
}

func (i *Indices) bulk(requireAlias bool, pipeline string) (res *esapi.Response, err error) {
	var bulkRequest []func(*esapi.BulkRequest)
	if i.index != "" {
		bulkRequest = append(bulkRequest, i.client.Bulk.WithIndex(i.index))
	}
	if pipeline != "" {
		bulkRequest = append(bulkRequest, i.client.Bulk.WithPipeline(pipeline))
	}
	if requireAlias != false {
		bulkRequest = append(bulkRequest, i.client.Bulk.WithRequireAlias(true))
	}
	bulkRequest = append(bulkRequest, i.client.Bulk.WithPretty())
	if i.rawFile != "" {
		body, err := i.readFromRawFile()
		if err != nil {
			return nil, err
		}
		return i.client.Bulk(bytes.NewReader(body),
			bulkRequest...)
	}
	return i.client.Bulk(strings.NewReader(i.data),
		bulkRequest...)
}

func (i *Indices) msearch(maxConcurrentSearches, maxConcurrentShardRequests int) (res *esapi.Response, err error) {
	var mSearchRequest []func(*esapi.MsearchRequest)
	if i.index != "" {
		mSearchRequest = append(mSearchRequest, i.client.Msearch.WithIndex(i.index))
	}
	if maxConcurrentSearches != 0 {
		mSearchRequest = append(mSearchRequest, i.client.Msearch.WithMaxConcurrentSearches(maxConcurrentSearches))
	}
	if maxConcurrentShardRequests != 0 {
		mSearchRequest = append(mSearchRequest, i.client.Msearch.WithMaxConcurrentShardRequests(maxConcurrentShardRequests))
	}
	mSearchRequest = append(mSearchRequest, i.client.Msearch.WithPretty())
	if i.rawFile != "" {
		body, err := i.readFromRawFile()
		if err != nil {
			return nil, err
		}
		return i.client.Msearch(bytes.NewReader(body),
			mSearchRequest...)
	}
	return i.client.Msearch(strings.NewReader(i.data),
		mSearchRequest...)
}

func (i *Indices) doReindex(body io.Reader) (res *esapi.Response, err error) {
	res, err = i.client.Reindex(body, i.client.Reindex.WithWaitForCompletion(false))
	return
}

func (i *Indices) readFromRawFile() ([]byte, error) {
	fmt.Println(i.rawFile)
	file, err := ioutil.ReadFile(i.rawFile)
	if err != nil {
		return nil, err
	}
	return file, nil
}
