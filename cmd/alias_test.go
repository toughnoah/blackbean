package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestCreateAlias(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("alias create test-* noah-test", mock)
	require.NoError(t, err)
}

func TestCreateAliasFromRaw(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`alias create test-* noah-test -d '{"test":"1"}' `, mock)
	require.NoError(t, err)
	_, err = executeCommand(`alias create test-* noah-test -f '../pkg/testdata/query.yaml' `, mock)
	require.NoError(t, err)
}

func TestDeleteAlias(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("alias delete test-* noah-test", mock)
	require.NoError(t, err)
}

func TestGetAlias(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "test get alias with flag",
			cmd:  "alias get noah-test",
		},
		{
			name: "test get alias with no args",
			cmd:  "alias get noah-test --is_alias=true",
		},
	}
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, mock)
		require.NoError(t, err)
	}
}
