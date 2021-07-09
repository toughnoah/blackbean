package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestCancleReroute(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("reroute cancel test-* --shard=0 --node=A", mock)
	require.NoError(t, err)
}

func TestReplicaAllocateReroute(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("reroute allocateReplicas test-* --shard=0 --node=A", mock)
	require.NoError(t, err)
}
func TestMoveReroute(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("reroute move test-* --shard=0 --from_node=A --to_node=B", mock)
	require.NoError(t, err)
}

func TestFailedReroute(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("reroute failed", mock)
	require.NoError(t, err)
}
