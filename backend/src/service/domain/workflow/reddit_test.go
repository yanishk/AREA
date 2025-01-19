package workflow_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestPostRedditComment(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.postRedditComment("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestUpvoteRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.upvoteRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestDownvoteRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.downvoteRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestSubmitRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.submitRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckRedditReactions(test *testing.T) {
	workflow := entities.Workflow{
		ReactionId: "1",
	}

	test.Run("Success", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		reactionFound := entities.Reaction{
			Name: "Test",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(reactionFound, nil)

		err := reddit.checkRedditReactions(workflow)

		require.NoError(test, err)
	})

	test.Run("Fail Find Reaction", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := reddit.checkRedditReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})
}

func TestIsANewPost(test *testing.T) {
	test.Run("Fail Update Workflow", func(t *testing.T) {
		var result entities.RedditPostResponse

		mockWorkflowRepo := new(MockWorkflowRepository)

		reddit := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		workflow := entities.Workflow{
			Id:         "1",
			ActionData: map[string]interface{}{},
		}

		redditResponse := `{
			"data": {
				"children": [
					{
						"data": {
							"id": "id",
							"title": "title",
							"author": "author",
							"url": "url"
						}
					}
				]
			}
		}`

		json.Unmarshal([]byte(redditResponse), &result)

		expectedWorkflow := workflow
		expectedWorkflow.ActionData = map[string]interface{}{"id": "id"}

		mockWorkflowRepo.On("UpdateWorkflow", expectedWorkflow.Id, expectedWorkflow).
			Return(fmt.Errorf("update error"))

		err := reddit.isANewPost(workflow, result)

		require.EqualError(t, err, errorUpdatingWorkflow)
	})
}

func TestIsANewComment(test *testing.T) {
	test.Run("Fail Update Workflow", func(t *testing.T) {
		var result entities.RedditCommentResponse

		mockWorkflowRepo := new(MockWorkflowRepository)

		reddit := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		workflow := entities.Workflow{
			Id:         "1",
			ActionData: map[string]interface{}{},
		}

		redditResponse := `{
			"data": {
				"children": [
					{
						"data": {
							"id": "id",
							"author": "author",
							"body": "body"
						}
					}
				]
			}
		}`

		json.Unmarshal([]byte(redditResponse), &result)

		expectedWorkflow := workflow
		expectedWorkflow.ActionData = map[string]interface{}{"id": "id"}

		mockWorkflowRepo.On("UpdateWorkflow", expectedWorkflow.Id, expectedWorkflow).
			Return(fmt.Errorf("update error"))

		err := reddit.isANewComment(workflow, result)

		require.EqualError(t, err, errorUpdatingWorkflow)
	})
}

func TestIsANewVote(test *testing.T) {
	test.Run("Fail Update Workflow", func(t *testing.T) {
		var result entities.RedditVoteResponse

		mockWorkflowRepo := new(MockWorkflowRepository)

		reddit := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		workflow := entities.Workflow{
			Id:         "1",
			ActionData: map[string]interface{}{},
		}

		redditResponse := `{
			"data": {
				"children": [
					{
						"data": {
							"id": "id"
						}
					}
				]
			}
		}`

		json.Unmarshal([]byte(redditResponse), &result)

		expectedWorkflow := workflow
		expectedWorkflow.ActionData = map[string]interface{}{"id": "id"}

		mockWorkflowRepo.On("UpdateWorkflow", expectedWorkflow.Id, expectedWorkflow).
			Return(fmt.Errorf("update error"))

		err := reddit.isANewVote(workflow, result)

		require.EqualError(t, err, errorUpdatingWorkflow)
	})
}

func TestCheckRedditNewPostInSubredditAction(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		reddit := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{"key": "value"},
		}

		err := reddit.checkRedditNewPostInSubredditAction("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckWorkflowsWithRedditActions(test *testing.T) {
	test.Run("Fail Find Workflow", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		reddit := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, errors.New("Fail find workflows"))

		err := reddit.checkWorkflowsWithRedditActions(action)

		require.EqualError(test, err, "Fail find workflows")
	})

	test.Run("Fail Find Workflow", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		reddit := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, nil)

		err := reddit.checkWorkflowsWithRedditActions(action)

		require.NoError(test, err)
	})
}

func TestCheckRedditActions(test *testing.T) {
	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		reddit := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("FindServiceByName", "Reddit").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := reddit.CheckRedditActions()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		reddit := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		serviceFound := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Reddit").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", serviceFound.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		err := reddit.CheckRedditActions()

		require.EqualError(test, err, "Fail find actions")
	})
}
