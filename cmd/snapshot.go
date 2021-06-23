package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func getSnapshot(cli *elasticsearch.Client) *cobra.Command {
	var so = Snapshot{client: cli}
	var (
		snapshots string
	)
	var command = &cobra.Command{
		Use:   "snapshots [repository]",
		Short: "get specific snapshots ",
		Long:  "get specific snapshots ... wordless",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return so.getAllSnapshots(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := so.getRepoAllSnapshots(args[0], getSnapshotsFromFlag(snapshots))
			fmt.Println(res)
			if err != nil {
				return err
			}
			return nil
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	return command
}

type Snapshot struct {
	client *elasticsearch.Client
}

func (S *Snapshot) getAllSnapshots() []string {
	var (
		reposMap  = make(map[string]interface{})
		repoSlice []string
	)
	repos, err := S.client.Snapshot.GetRepository(S.client.Snapshot.GetRepository.WithRepository("_all"))
	if err != nil {
		log.Fatalf("error sending request to es: %s", err)
	}
	if err = json.NewDecoder(repos.Body).Decode(&reposMap); err != nil {
		log.Fatalf("error parsing the response body: %s", err)
	}
	for repo, _ := range reposMap {
		repoSlice = append(repoSlice, repo)
	}
	return repoSlice
}

func (S *Snapshot) getRepoAllSnapshots(repos string, snapshot []string) (res *esapi.Response, err error) {
	res, err = S.client.Snapshot.Get(repos, snapshot, S.client.Snapshot.Get.WithPretty())
	return
}

func (S *Snapshot) getRepoAllSnapshotsForFlag(repos string) []string {
	var resMap = make(map[string][]map[string]interface{})
	var resSlice []string
	res, err := S.client.Snapshot.Get(repos, []string{"_all"}, S.client.Snapshot.Get.WithPretty())
	if err != nil {
		log.Fatalf("error sending request to es: %s", err)
	}
	if err = json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Fatalf("error parsing the response body: %s", err)
	}
	if resMap["snapshots"] == nil {
		return resSlice
	}
	for _, snapshot := range resMap["snapshots"] {
		resSlice = append(resSlice, snapshot["snapshot"].(string))
	}
	return resSlice
}

func getSnapshotsFromFlag(snapshots string) []string {
	return strings.Split(snapshots, ",")
}
