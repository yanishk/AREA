package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"backend/src/entities"
)

func (self *WorkflowService) sendMeEmail(workflow entities.Workflow) error {
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errUser != nil {
		return fmt.Errorf("Error finding user")
	}

	subject, subjectExists := workflow.ReactionParam["subject"]
	body, bodyExists := workflow.ReactionParam["body"]
	if !subjectExists || !bodyExists {
		return fmt.Errorf(errorMissingField)
	}

	username, _, stringSplit := strings.Cut(foundUser.Email, "@")
	if !stringSplit {
		username = "Default Name"
	}

	content := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]interface{}{
					{
						"email": foundUser.Email,
						"name":  username,
					},
				},
				"subject": subject,
			},
		},
		"from": map[string]interface{}{
			"email": os.Getenv("SENDER_EMAIL"),
			"name":  "AREA",
		},
		"content": []map[string]interface{}{
			{
				"type":  "text/plain",
				"value": body,
			},
		},
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	url := "https://api.sendgrid.com/v3/mail/send"

	res, err := self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, os.Getenv("API_KEY"), bytes.NewBuffer(contentBytes))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) checkSendEmailReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	switch reactionFound.Name {
	case "Send me an email":
		return self.sendMeEmail(workflow)
	}
	return nil
}
