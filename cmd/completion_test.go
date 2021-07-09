package cmd

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"strings"
	"testing"
)

var yamlExample = []byte(`cluster:
  default: 
    url: https://a.es.com:9200
    username: Noah
    password: abc
  backup: 
    url: https://a.es.com:9200
    username: Noah
    password: abc
current: default`)

var _ = Describe("cat resources test", func() {
	Context("test no FileCompletion", func() {
		It("test noCompletions", func() {
			noCompletions(nil, nil, "")
		})
	})
})

func TestCompletion(t *testing.T) {
	testCases := []struct {
		cmd string
	}{
		{
			cmd: "completion bash ",
		},
		{
			cmd: "completion zsh",
		},
		{
			cmd: "completion fish",
		},
		{
			cmd: "completion powershell",
		},
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, nil)
		require.NoError(t, err)
	}
}
func TestCompletion2(t *testing.T) {
	r := bytes.NewReader(yamlExample)
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(r)
	require.NoError(t, err)
	testCases := []struct {
		name     string
		cmd      string
		checkOut string
		mock     *fake.MockEsResponse
	}{
		{
			cmd:      "__complete cat ''",
			checkOut: "health\nnodes\nallocations\nthreadpool\ncachemem\nsegmem\nlargeindices",
		},
		{
			cmd:      "__complete apply settings --allocation_enable ''",
			checkOut: "primaries\nnull\n",
		},
		{
			cmd:      "__complete repo get ''",
			checkOut: "repoA",
			mock: &fake.MockEsResponse{
				ResponseString: `{"repoA":"a","repoB":"b"}`,
			},
		},
		{
			cmd:      "__complete repo get test --snapshot ''",
			checkOut: "snapshot01\nsnapshot01\n",
			mock: &fake.MockEsResponse{
				ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
			},
		},
		{
			cmd:      "__complete repo delete ''",
			checkOut: "repoA",
			mock: &fake.MockEsResponse{
				ResponseString: `{"repoA":"a","repoB":"b"}`,
			},
		},
		{
			cmd:      "__complete snapshot get --repo ''",
			checkOut: "repoA",
			mock: &fake.MockEsResponse{
				ResponseString: `{"repoA":"a","repoB":"b"}`,
			},
		},
		{
			cmd:      "__complete snapshot get --repo repoA ''",
			checkOut: "snapshot01\nsnapshot01\n",
			mock: &fake.MockEsResponse{
				ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
			},
		},
		{
			cmd:      "__complete snapshot delete --repo repoA ''",
			checkOut: "snapshot01",
			mock: &fake.MockEsResponse{
				ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot01"}]}`,
			},
		},
		{
			cmd:      "__complete snapshot delete --repo ''",
			checkOut: "repoA",
			mock: &fake.MockEsResponse{
				ResponseString: `{"repoA":"a","repoB":"b"}`,
			},
		},
		{
			cmd:      "__complete index search ''",
			checkOut: "index1",
			mock: &fake.MockEsResponse{
				ResponseString: `{"index1":"a","index2":"b"}`,
			},
		},
		{
			name:     "test alias completion",
			cmd:      "__complete alias get ''",
			checkOut: "A",
			mock: &fake.MockEsResponse{
				ResponseString: `{"A":{"aliases":{"a":{}}},"B":{"aliases":{"b":{}}}}`,
			},
		},
		{
			name:     "test alias completion",
			cmd:      "__complete alias delete ''",
			checkOut: "A",
			mock: &fake.MockEsResponse{
				ResponseString: `{"A":{"aliases":{"a":{}}},"B":{"aliases":{"b":{}}}}`,
			},
		},
		{
			name:     "test reroute completion",
			cmd:      "__complete reroute move noah-test --from_node ''",
			checkOut: "elasticsearch-master-0",
			mock: &fake.MockEsResponse{
				ResponseString: `{"nodes":{"A": {"name":"elasticsearch-master-0"}}}`,
			},
		},
		{
			name:     "test reroute completion",
			cmd:      "__complete reroute cancel noah-test --node ''",
			checkOut: "elasticsearch-master-0",
			mock: &fake.MockEsResponse{
				ResponseString: `{"nodes":{"A": {"name":"elasticsearch-master-0"}}}`,
			},
		},
		{
			name:     "test role completion",
			cmd:      "__complete role get ''",
			checkOut: "superuser",
			mock: &fake.MockEsResponse{
				ResponseString: `{"superuser":{"cluster":["all"],"indices":[{"names":["*"],"privileges":["all"],"allow_restricted_indices":true}],"applications":[{"application":"*","privileges":["*"],"resources":["*"]}],"run_as":["*"],"metadata":{"_reserved":true},"transient_metadata":{}}}`,
			},
		},
		{
			name:     "test user completion",
			cmd:      "__complete user get ''",
			checkOut: "Noah.Lu",
			mock: &fake.MockEsResponse{
				ResponseString: `{"Noah.Lu":{"username":"Noah.Lu","roles":["superuser"],"full_name":"Noah.Lu","email":"noah.lu@163.com","metadata":{"intelligence":7},"enabled":true}}`,
			},
		},
		{
			name:     "test user completion",
			cmd:      "__complete user delete ''",
			checkOut: "Noah.Lu",
			mock: &fake.MockEsResponse{
				ResponseString: `{"Noah.Lu":{"username":"Noah.Lu","roles":["superuser"],"full_name":"Noah.Lu","email":"noah.lu@163.com","metadata":{"intelligence":7},"enabled":true}}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role create noah --cluster_privilege ''",
			checkOut: "create_snapshot",
			mock: &fake.MockEsResponse{
				ResponseString: `{"cluster":["all","create_snapshot","delegate_pki","manage","manage_api_key","manage_ccr","manage_data_frame_transforms","manage_enrich","manage_ilm","manage_index_templates","manage_ingest_pipelines","manage_ml","manage_oidc","manage_own_api_key","manage_pipeline","manage_rollup","manage_saml","manage_security","manage_slm","manage_token","manage_transform","manage_watcher","monitor","monitor_data_frame_transforms","monitor_ml","monitor_rollup","monitor_transform","monitor_watcher","none","read_ccr","read_ilm","read_slm","transport_client"],"index":["all","create","create_doc","create_index","delete","delete_index","index","manage","manage_follow_index","manage_ilm","manage_leader_index","monitor","none","read","read_cross_cluster","view_index_metadata","write"]}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role create noah --indices ''",
			checkOut: "index1",
			mock: &fake.MockEsResponse{
				ResponseString: `{"index1":"a","index2":"b"}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role create noah --indices test --indices_privilege ''",
			checkOut: "read",
			mock: &fake.MockEsResponse{
				ResponseString: `{"cluster":["all","create_snapshot","delegate_pki","manage","manage_api_key","manage_ccr","manage_data_frame_transforms","manage_enrich","manage_ilm","manage_index_templates","manage_ingest_pipelines","manage_ml","manage_oidc","manage_own_api_key","manage_pipeline","manage_rollup","manage_saml","manage_security","manage_slm","manage_token","manage_transform","manage_watcher","monitor","monitor_data_frame_transforms","monitor_ml","monitor_rollup","monitor_transform","monitor_watcher","none","read_ccr","read_ilm","read_slm","transport_client"],"index":["all","create","create_doc","create_index","delete","delete_index","index","manage","manage_follow_index","manage_ilm","manage_leader_index","monitor","none","read","read_cross_cluster","view_index_metadata","write"]}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role update noah --indices test --indices_privilege ''",
			checkOut: "read",
			mock: &fake.MockEsResponse{
				ResponseString: `{"cluster":["all","create_snapshot","delegate_pki","manage","manage_api_key","manage_ccr","manage_data_frame_transforms","manage_enrich","manage_ilm","manage_index_templates","manage_ingest_pipelines","manage_ml","manage_oidc","manage_own_api_key","manage_pipeline","manage_rollup","manage_saml","manage_security","manage_slm","manage_token","manage_transform","manage_watcher","monitor","monitor_data_frame_transforms","monitor_ml","monitor_rollup","monitor_transform","monitor_watcher","none","read_ccr","read_ilm","read_slm","transport_client"],"index":["all","create","create_doc","create_index","delete","delete_index","index","manage","manage_follow_index","manage_ilm","manage_leader_index","monitor","none","read","read_cross_cluster","view_index_metadata","write"]}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role update noah --indices ''",
			checkOut: "index1",
			mock: &fake.MockEsResponse{
				ResponseString: `{"index1":"a","index2":"b"}`,
			},
		},
		{
			name:     "test privilege completion",
			cmd:      "__complete role update noah --cluster_privilege ''",
			checkOut: "create_snapshot",
			mock: &fake.MockEsResponse{
				ResponseString: `{"cluster":["all","create_snapshot","delegate_pki","manage","manage_api_key","manage_ccr","manage_data_frame_transforms","manage_enrich","manage_ilm","manage_index_templates","manage_ingest_pipelines","manage_ml","manage_oidc","manage_own_api_key","manage_pipeline","manage_rollup","manage_saml","manage_security","manage_slm","manage_token","manage_transform","manage_watcher","monitor","monitor_data_frame_transforms","monitor_ml","monitor_rollup","monitor_transform","monitor_watcher","none","read_ccr","read_ilm","read_slm","transport_client"],"index":["all","create","create_doc","create_index","delete","delete_index","index","manage","manage_follow_index","manage_ilm","manage_leader_index","monitor","none","read","read_cross_cluster","view_index_metadata","write"]}`,
			},
		},
		{
			name:     "test reroute completion",
			cmd:      "__complete reroute allocateReplicas noah-test --node ''",
			checkOut: "elasticsearch-master-0",
			mock: &fake.MockEsResponse{
				ResponseString: `{"nodes":{"A": {"name":"elasticsearch-master-0"}}}`,
			},
		},
	}
	for _, tc := range testCases {
		out, err := executeCommand(tc.cmd, tc.mock)
		require.NoError(t, err)
		require.Equal(t, strings.Contains(out, tc.checkOut), true)
	}
}

func TestCompletion3(t *testing.T) {
	testCases := []struct {
		cmd string
	}{
		{
			cmd: "completion bash ",
		},
		{
			cmd: "completion zsh",
		},
		{
			cmd: "completion fish",
		},
		{
			cmd: "completion powershell",
		},
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, nil)
		require.NoError(t, err)
	}
}
