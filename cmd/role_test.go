package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestCreateRole(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"superuser":true}`,
	}
	_, err := executeCommand("role create noah-test-role --cluster_privilege=all --indices=test-*", mock)
	require.NoError(t, err)
}

func TestGetRole(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "test get role with args",
			cmd:  "role get",
		},
		{
			name: "test get role with no args",
			cmd:  "role get noah",
		},
	}
	mock := &fake.MockEsResponse{
		ResponseString: `{"noah-test":{"roles":["superuser"]}, "superuser":{"roles":["superuser"]}}`,
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, mock)
		require.NoError(t, err)
	}

}

func TestDeleteRole(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("role delete noah-test", mock)
	require.NoError(t, err)
}

func TestUpdateRole(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "test change cluster privilege",
			cmd:  "role update noah-test --cluster_privilege=all",
		},
		{
			name: "test change cluster privilege and add only is false",
			cmd:  "role update noah-test --cluster_privilege=all --add_only=false",
		},
		{
			name: "test add indices",
			cmd:  "role update noah-test --indices=api,kpi,test",
		},
		{
			name: "test add indices and override",
			cmd:  "role update noah-test --indices=api,kpi,test --add_only=false ",
		},
		{
			name: "test add indices and override",
			cmd:  "role update noah-test --indices=api,kpi,test --add_only=false --indices_privilege=all ",
		},
	}
	mock := &fake.MockEsResponse{
		ResponseString: `{"noah-test":{"cluster":["create_snapshot"], "indices":[{}]}}`,
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, mock)
		require.NoError(t, err)
	}

}

func TestUpdateRoleFailedWithNoFlag(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("role update noah-test", mock)
	require.Error(t, err)
}
