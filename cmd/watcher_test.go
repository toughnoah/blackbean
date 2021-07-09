package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestWatcher(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	testCases := []struct {
		name string
		cmd  string
	}{
		{
			name: "start watcher",
			cmd:  "watcher start",
		},
		{
			name: "stop watcher",
			cmd:  "watcher stop",
		},
		{
			name: "get watcher stats",
			cmd:  "watcher stats",
		},
	}
	for _, tc := range testCases {
		_, err := executeCommand(tc.cmd, mock)
		require.NoError(t, err)
	}

}
