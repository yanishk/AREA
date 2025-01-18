package workflow_service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestMoveFileOrFolder(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.moveFileOrFolder("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCreateTextFile(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.createTextFile("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestAppendTextFile(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.appendTextFile("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}
