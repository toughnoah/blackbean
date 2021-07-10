package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/toughnoah/blackbean/pkg/fake"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test":{}}`,
	}
	_, err := executeCommand("template get", mock)
	require.NoError(t, err)
	_, err = executeCommand("template get test", mock)
	require.NoError(t, err)
}

func TestApplyTemplate(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"test":{}}`,
	}
	_, err := executeCommand(`template apply test -d "{"test":"abc"}"`, mock)
	require.NoError(t, err)
	_, err = executeCommand("template apply test -f ../pkg/testdata/template.json", mock)
	require.NoError(t, err)
	_, err = executeCommand("template apply test", mock)
	require.Error(t, err)
}

func TestDeleteTemplate(t *testing.T) {
	mock := &fake.MockEsResponse{
		ResponseString: `{"acknowledge":true}`,
	}
	_, err := executeCommand(`template delete test`, mock)
	require.NoError(t, err)
}
