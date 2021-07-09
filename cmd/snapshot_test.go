package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/es"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestGetRepoAllSnapshotsForFlag(t *testing.T) {
	testCases := []struct {
		name string
		mock fake.Mock
		want []string
	}{
		{
			name: "test right response",
			mock: &fake.MockEsResponse{
				ResponseString: `{"snapshots":[{"snapshot":"snapshot01"}, {"snapshot":"snapshot02"}]}`,
			},
			want: []string{
				"snapshot01", "snapshot02",
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
			name: "test wrong response",
			mock: &fake.MockEsResponse{
				ResponseString: `{"snapshots":[]}`,
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
		require.Equal(t, tc.want, so.getRepoAllSnapshotsForFlag("test"))
	}
}

func TestRecoverIndices(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test": "recovery indices"}`,
	}
	fakeClient, err := es.NewEsClient("https://test.com", "a", "b", mock)
	require.NoError(t, err)
	so := Snapshot{
		client: fakeClient,
	}
	_, err = so.recoverIndices("", "", "", "", "")
	require.NoError(t, err)
}

func TestCreateSnapshot(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test": "create snapshot"}`,
	}
	fakeClient, err := es.NewEsClient("https://test.com", "a", "b", mock)
	require.NoError(t, err)
	so := Snapshot{
		client: fakeClient,
	}
	_, err = so.createSnapshot("", "")
	require.NoError(t, err)
}
func TestDeleteSnapshot(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test": "delete snapshot"}`,
	}
	fakeClient, err := es.NewEsClient("https://test.com", "a", "b", mock)
	require.NoError(t, err)
	so := Snapshot{
		client: fakeClient,
	}
	_, err = so.deleteSnapshot("", "")
	require.NoError(t, err)
}
func TestGetSnapshot(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test": "get snapshot"}`,
	}
	fakeClient, err := es.NewEsClient("https://test.com", "a", "b", mock)
	require.NoError(t, err)
	so := Snapshot{
		client: fakeClient,
	}
	_, err = so.getSnapshot("", "")
	require.NoError(t, err)
}

func TestSnapshotCommand(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
		mock *fake.MockEsResponse
	}{
		{
			name: "test get cmd",
			cmd:  "snapshot get snapshot01 --repo repo",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test":" test get"}`,
			},
		},
		{
			name: "test create cmd",
			cmd:  "snapshot create snapshot01 --repo repo",
			mock: &fake.MockEsResponse{
				ResponseString: `{"acknowledge":"true"}`,
			},
		},
		{
			name: "test delete cmd",
			cmd:  "snapshot delete snapshot01 --repo repo",
			mock: &fake.MockEsResponse{
				ResponseString: `{"acknowledge":"true"}`,
			},
		},
	}
	for _, tc := range testCases {
		out, err := executeCommand(tc.cmd, tc.mock)
		require.NoError(t, err)
		require.Equal(t, out, "[200 OK] "+tc.mock.ResponseString+"\n")
	}

}
