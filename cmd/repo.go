package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
)

func repo(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:               "repo [subcommand]",
		Short:             "repo operations ",
		Long:              "repo operations ... wordless",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: noCompletions,
	}
	command.AddCommand(getRepos(cli, out))
	command.AddCommand(createRepo(cli, out))
	command.AddCommand(deleteRepo(cli, out))
	return command
}

func getRepos(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		snapshots string
	)
	var command = &cobra.Command{
		Use:   "get [repository]",
		Short: "get specific repository ",
		Long:  "get specific repository ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.getRepoAllSnapshots(args[0], snapshots)
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&snapshots, "snapshot", "s", "_all", "to get specific snapshot")
	err := command.RegisterFlagCompletionFunc("snapshot", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return so.getRepoAllSnapshotsForFlag(args[0]), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(-1)
	}
	return command
}

func createRepo(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		containType string
		container   string
		path        string
	)
	var command = &cobra.Command{
		Use:               "create [repository]",
		Short:             "create specific snapshots ",
		Long:              "create specific snapshots ... wordless",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.createSnapshotRepo(containType, container, path, args[0])
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&containType, "type", "e", "", "to specify repo type")
	f.StringVarP(&container, "container", "c", "", "to specify repo container")
	f.StringVarP(&path, "path", "p", "", "to specify repo path")
	_ = command.MarkFlagRequired("type")
	_ = command.MarkFlagRequired("container")
	_ = command.MarkFlagRequired("path")
	return command
}

func deleteRepo(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var command = &cobra.Command{
		Use:   "delete [repository]",
		Short: "delete specific snapshots ",
		Long:  "delete specific snapshots ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.deleteSnapshotRepo(args[0])
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	return command
}

func (S *Snapshot) createSnapshotRepo(containType, container, path, repo string) (res *esapi.Response, err error) {
	var body = fmt.Sprintf(`{"type": "%s","settings": {"container": "%s","base_path": "%s","chunk_size": "32m","compress": true,"max_snapshot_bytes_per_sec" : "50mb","max_restore_bytes_per_sec" : "50mb"}}`,
		containType, container, path)
	res, err = S.client.Snapshot.CreateRepository(repo, strings.NewReader(body))
	return
}

func (S *Snapshot) deleteSnapshotRepo(repo string) (res *esapi.Response, err error) {
	fmt.Println(splitWords(repo))
	res, err = S.client.Snapshot.DeleteRepository(splitWords(repo))
	return
}

func (S *Snapshot) getAllRepos() []string {
	var (
		reposMap  = make(map[string]interface{})
		repoSlice []string
	)
	repos, err := S.client.Snapshot.GetRepository(S.client.Snapshot.GetRepository.WithRepository("_all"))
	if err != nil {
		log.Printf("error sending request to es: %s", err)
		return nil
	}
	if err = json.NewDecoder(repos.Body).Decode(&reposMap); err != nil {
		log.Printf("error parsing the response body: %s", err)
		return nil
	}
	for rep, _ := range reposMap {
		repoSlice = append(repoSlice, rep)
	}
	return repoSlice
}
