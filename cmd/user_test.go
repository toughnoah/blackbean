package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"io"
	"testing"
)

type MockTerminal struct {
	toSend       []byte
	bytesPerRead int
	received     []byte
}

func (c *MockTerminal) Read(data []byte) (n int, err error) {
	n = len(data)
	if n == 0 {
		return
	}
	if n > len(c.toSend) {
		n = len(c.toSend)
	}
	if n == 0 {
		return 0, io.EOF
	}
	if c.bytesPerRead > 0 && n > c.bytesPerRead {
		n = c.bytesPerRead
	}
	copy(data, c.toSend[:n])
	c.toSend = c.toSend[n:]
	return
}

func (c *MockTerminal) Write(data []byte) (n int, err error) {
	c.received = append(c.received, data...)
	return len(data), nil
}

func TestCreateUser(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"superuser":true}`,
	}
	_, err := executeCommand("user create noah-test --roles=superuser", mock)
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"superuser":true}`,
	}
	_, err := executeCommand("user delete noah-test", mock)
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "test get user with args",
			cmd:  "user get",
		},
		{
			name: "test get user with no args",
			cmd:  "user get noah",
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

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "test change password and roles",
			cmd:  "user update noah-test --roles=superuser --change_password",
		},
		{
			name: "test only change password",
			cmd:  "user update noah-test --change_password",
		},
		{
			name: "test only change roles",
			cmd:  "user update noah-test --roles=superuser",
		},
		{
			name: "test only change roles and add only ",
			cmd:  "user update noah-test --roles=superuser --add_only=false",
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

func TestMergeRoles(t *testing.T) {
	roles := []interface{}{
		"superuser", "test-app", "kibana"}

	newRoles := []interface{}{
		"superuser", "test-sre",
	}
	require.Equal(t, len(mergeRoles(roles, newRoles)), 4)
}
