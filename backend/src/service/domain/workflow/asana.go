package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"backend/src/entities"
)

func asanaReactions() []string {
	return []string{
		"Create task",
		"Create project",
	}
}

func (self *WorkflowService) createTaskAsana(accessToken string, workflow entities.Workflow) error {
	workspace, workspaceExists := workflow.ReactionParam["workspace"]
	project, projectExists := workflow.ReactionParam["project"]
	taskName, taskNameExists := workflow.ReactionParam["name"]
	notes, notesExists := workflow.ReactionParam["notes"]
	due, dueExists := workflow.ReactionParam["due"]
	assignee, assigneeExists := workflow.ReactionParam["assignee"]
	tag, tagExists := workflow.ReactionParam["tag"]

	if !taskNameExists || !workspaceExists {
		return fmt.Errorf(errorMissingField)
	}

	content := map[string]interface{}{
		"data": map[string]interface{}{
			"workspace": workspace,
			"name":      taskName,
		},
	}

	if projectExists {
		content["data"].(map[string]interface{})["projects"] = project
	}
	if notesExists {
		content["data"].(map[string]interface{})["notes"] = notes
	}
	if dueExists {
		content["data"].(map[string]interface{})["due_on"] = due
	}
	if assigneeExists {
		content["data"].(map[string]interface{})["assignee"] = assignee
	}
	if tagExists {
		content["data"].(map[string]interface{})["tags"] = tag
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	url := "https://app.asana.com/api/1.0/tasks"
	resp, err := self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, accessToken, bytes.NewBuffer(contentBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) createProjectAsana(accessToken string, workflow entities.Workflow) error {
	workspace, workspaceExists := workflow.ReactionParam["workspace"]
	name, nameExists := workflow.ReactionParam["name"]
	description, descriptionExists := workflow.ReactionParam["description"]
	due, dueExists := workflow.ReactionParam["due"]

	if !workspaceExists || !nameExists {
		return fmt.Errorf(errorMissingField)
	}

	content := map[string]interface{}{
		"data": map[string]interface{}{
			"workspace": workspace,
			"name":      name,
		},
	}

	if descriptionExists {
		content["data"].(map[string]interface{})["notes"] = description
	}
	if dueExists {
		content["data"].(map[string]interface{})["due_on"] = due
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	url := "https://app.asana.com/api/1.0/projects"
	resp, err := self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, accessToken, bytes.NewBuffer(contentBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) checkAsanaReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Asana", reactionFound.Name, asanaReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	switch reactionFound.Name {
	case "Create task":
		return self.createTaskAsana(accessToken, workflow)
	case "Create project":
		return self.createProjectAsana(accessToken, workflow)
	}

	return nil
}
