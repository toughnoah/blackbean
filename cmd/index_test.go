package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func Test_decodeFromFile(t *testing.T) {
	testCase := []struct {
		name string
		file string
		want string
	}{
		{
			name: "test json",
			file: "../pkg/testdata/query.json",
			want: `{"query":{"match_all": {}}}`,
		},
		{
			name: "test yaml",
			file: "../pkg/testdata/query.yaml",
			want: `{"query":{"match":{"name":"test"}}}`,
		},
	}
	for _, tc := range testCase {
		file, err := decodeFromFile(tc.file)
		if err != nil {
			return
		}
		require.Equal(t, tc.want, string(file))
	}
}

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
