package workflow_service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestGetSpotifyUrl(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{"key": "value"},
		}

		baseUrl := "base"
		parameterName := "key"
		endUrl := "endurl"

		expectedUrl := baseUrl + workflow.ReactionParam[parameterName].(string) + endUrl

		res, err := getSpotifyUrl(workflow, baseUrl, parameterName, endUrl)

		require.NoError(test, err)
		assert.Equal(test, expectedUrl, res)
	})

	test.Run("Failure", func(test *testing.T) {
		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{"key": "value"},
		}

		baseUrl := "base"
		endUrl := "endurl"

		_, err := getSpotifyUrl(workflow, baseUrl, "fail", endUrl)

		require.Error(test, err)
	})
}
