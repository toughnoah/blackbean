package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/es"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestGetRepo(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
		mock *fake.MockEsResponse
	}{
		{
			name: "get repo",
			cmd:  "repo get test --snapshot test",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test":"get repo"}`,
			},
		},
		{
			name: "create repo",
			cmd:  "repo create test --type azure --container test --path /abc",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test":"create repo"}`,
			},
		},
		{
			name: "delete repo",
			cmd:  "repo delete test",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test":"delete repo"}`,
			},
		},
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, tc.mock)
		require.NoError(t, err)
	}
}
func TestGetAllRepos(t *testing.T) {
	testCases := []struct {
		name string
		mock fake.Mock
		want []string
	}{
		{
			name: "test right response",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test":"delete repo","abc":"create repo"}`,
			},
			want: []string{
				"test", "abc",
			},
		},
		{
			name: "test wrong response",
			mock: &fake.MockEsResponse{
				ResponseString: `a`,
			},
			want: nil,
		},
		{
			name: "test send request failed",
			mock: &fake.MockErrorEsResponse{},
			want: nil,
		},
	}
	for _, tc := range testCases {
		fakeClient, err := es.NewEsClient("https://test.com", "a", "b", tc.mock)
		require.NoError(t, err)
		so := Snapshot{
			client: fakeClient,
		}
		if tc.want != nil {
			for _, c := range tc.want {
				require.Equal(t, true, es.Check(c, so.getAllRepos()))
			}
		} else {
			require.Equal(t, []string(nil), so.getAllRepos())
		}
	}
}
