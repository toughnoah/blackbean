package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestSearchIndex(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
		mock *fake.MockEsResponse
	}{
		{
			name: "test json",
			cmd:  "index search test-* -f ../pkg/testdata/query.json",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test": "hello json"}`,
			},
		},
		{
			name: "test json",
			cmd:  "index search test-* -f ../pkg/testdata/query.json",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test": "hello yaml"}`,
			},
		},
		{
			name: "test json",
			cmd:  "index get test-*",
			mock: &fake.MockEsResponse{
				ResponseString: `{"test": "hello yaml"}`,
			},
		},
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, tc.mock)
		require.NoError(t, err)
	}
}

func TestCreateIndex(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index create test-*  -d "{"mapping":{}}"`, mock)
	require.NoError(t, err)
}

func TestDeleteIndex(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index delete test-* `, mock)
	require.NoError(t, err)
}

func TestReIndex(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index reindex test-* noah-test-*`, mock)
	require.NoError(t, err)
}

func TestWriteIndex(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index write test -d "{"test":"abc"}"`, mock)
	require.NoError(t, err)
	_, err = executeCommand("index write test -f ../pkg/testdata/write.json", mock)
	require.NoError(t, err)
	_, err = executeCommand("index write test", mock)
	require.Error(t, err)
}

func TestBulk(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index bulk -d "{"test":"abc"}"`, mock)
	require.NoError(t, err)
	_, err = executeCommand("index bulk --raw_file ../pkg/testdata/bulk.json --pipeline test --require_alias true", mock)
	require.NoError(t, err)
	_, err = executeCommand("index bulk", mock)
	require.Error(t, err)
}

func TestMsearch(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`index msearch -d "{"test":"abc"}"`, mock)
	require.NoError(t, err)
	_, err = executeCommand("index msearch --raw_file ../pkg/testdata/bulk.json --max_concurrent_searches 1 --max_concurrent_shard_requests 1", mock)
	require.NoError(t, err)
	_, err = executeCommand("index msearch", mock)
	require.Error(t, err)
}
