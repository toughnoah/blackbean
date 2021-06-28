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

func snapshot(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var command = &cobra.Command{
		Use:               "snapshot [subcommand]",
		Short:             "snapshot operations ",
		Long:              "snapshot operations ... wordless",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: noCompletions,
	}
	command.AddCommand(restoreSnapshot(cli, out))
	command.AddCommand(createSnapshot(cli, out))
	command.AddCommand(deleteSnapshot(cli, out))
	command.AddCommand(getSnapshot(cli, out))
	return command
}

func createSnapshot(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		repository string
	)
	var command = &cobra.Command{
		Use:               "create [snapshot]",
		Short:             "create specific snapshots ",
		Long:              "create specific snapshots ... wordless",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.createSnapshot(repository, args[0])
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&repository, "repo", "r", "", "to specify repo")
	err := command.RegisterFlagCompletionFunc("repo", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = command.MarkFlagRequired("repo")
	return command
}

func deleteSnapshot(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		repository string
	)
	var command = &cobra.Command{
		Use:   "delete [snapshot]",
		Short: "delete specific snapshots ",
		Long:  "delete specific snapshots ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getRepoAllSnapshotsForFlag(repository), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.deleteSnapshot(repository, args[0])
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&repository, "repo", "r", "", "to specify repo")
	err := command.RegisterFlagCompletionFunc("repo", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(-1)
	}
	_ = command.MarkFlagRequired("repo")
	return command
}

func getSnapshot(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		repository string
	)
	var command = &cobra.Command{
		Use:   "get [snapshot]",
		Short: "get specific snapshots ",
		Long:  "get specific snapshots ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getRepoAllSnapshotsForFlag(repository), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.getSnapshot(repository, args[0])
			fmt.Fprintf(out, "%s\n", res)
			return err
		},
	}
	f := command.Flags()
	f.StringVarP(&repository, "repo", "r", "", "to specify repo")
	err := command.RegisterFlagCompletionFunc("repo", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = command.MarkFlagRequired("repo")
	return command
}

func restoreSnapshot(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		snapshots         string
		index             string
		renamePattern     string
		renameReplacement string
	)
	var command = &cobra.Command{
		Use:   "restore [repository]",
		Short: "get specific index to restore ",
		Long:  "get specific index to restore ...wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getAllRepos(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.recoverIndices(args[0], snapshots, index, renamePattern, renameReplacement)
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
		log.Fatal(err)
	}
	f.StringVarP(&index, "index", "i", "", "to get specific index to restore")
	f.StringVarP(&renamePattern, "rename_pattern", "p", "", "to specify rename_pattern")
	f.StringVarP(&renameReplacement, "rename_replacement", "r", "", "to specify rename_replacement")
	_ = command.MarkFlagRequired("index")
	_ = command.MarkFlagRequired("rename_pattern")
	_ = command.MarkFlagRequired("rename_replacement")
	return command
}

type Snapshot struct {
	client *elasticsearch.Client
}

func (S *Snapshot) getRepoAllSnapshots(repos string, snapshot string) (res *esapi.Response, err error) {
	res, err = S.client.Snapshot.Get(repos, splitWords(snapshot), S.client.Snapshot.Get.WithPretty())
	return
}

func (S *Snapshot) getRepoAllSnapshotsForFlag(repos string) []string {
	var resMap = make(map[string][]map[string]interface{})
	var resSlice []string
	res, err := S.client.Snapshot.Get(repos, []string{"_all"}, S.client.Snapshot.Get.WithPretty())
	if err != nil {
		log.Printf("error sending request to es: %s", err)
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Printf("error parsing the response body: %s", err)
		return nil
	}
	if resMap["snapshots"] == nil {
		return nil
	}
	for _, snapshot := range resMap["snapshots"] {
		resSlice = append(resSlice, snapshot["snapshot"].(string))
	}
	return resSlice
}

func (S *Snapshot) recoverIndices(repo, snapshot, index, renamePattern, renameReplacement string) (res *esapi.Response, err error) {
	var body = fmt.Sprintf(`{"indices": "%s","include_global_state": true,"rename_pattern": "%s","rename_replacement": "%s","include_aliases": false}`,
		index,
		renamePattern,
		renameReplacement)
	res, err = S.client.Snapshot.Restore(repo, snapshot, S.client.Snapshot.Restore.WithBody(strings.NewReader(body)))
	return
}

func (S *Snapshot) createSnapshot(repo, snapshot string) (res *esapi.Response, err error) {
	res, err = S.client.Snapshot.Create(repo, snapshot)
	return res, err
}

func (S *Snapshot) deleteSnapshot(repo, snapshot string) (res *esapi.Response, err error) {
	res, err = S.client.Snapshot.Delete(repo, snapshot)
	return
}

func (S *Snapshot) getSnapshot(repo, snapshot string) (res *esapi.Response, err error) {
	res, err = S.client.Snapshot.Get(repo, splitWords(snapshot), S.client.Snapshot.Get.WithPretty())
	return
}

func splitWords(words string) []string {
	return strings.Split(words, ",")
}
