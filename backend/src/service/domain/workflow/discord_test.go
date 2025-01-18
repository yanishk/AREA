package workflow_service

import (
	"backend/src/entities"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostDiscordMessage(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		discord := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := discord.postDiscordMessage("tokenBot", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCreateDiscordThread(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		discord := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := discord.createDiscordThread("tokenBot", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}
