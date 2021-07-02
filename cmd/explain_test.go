package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestExplain(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("explain ", mock)
	require.NoError(t, err)
}

func TestExplainIndex(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("explain test-* --primary --shard=0 ", mock)
	require.NoError(t, err)
}

func TestExplainIndexError(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand("explain test-*", mock)
	require.Error(t, err)
}
