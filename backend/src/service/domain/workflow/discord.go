package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"backend/src/entities"
)

const botBearer = "Bot "

func (self *WorkflowService) postDiscordMessage(tokenBot string, workflow entities.Workflow) error {
	message, messageErr := workflow.ReactionParam["message"]
	channelID, channelErr := workflow.ReactionParam["channel"]
	if !messageErr || !channelErr {
		return fmt.Errorf(errorMissingField)
	}

	url := "https://discord.com/api/v10/channels/" + channelID.(string) + "/messages"
	content := map[string]string{
		"content": message.(string),
	}

	jsonBody, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "POST", botBearer, tokenBot, bytes.NewBuffer(jsonBody))
	if errRequest != nil {
		return errRequest
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) createDiscordThread(tokenBot string, workflow entities.Workflow) error {
	title, titleErr := workflow.ReactionParam["title"]
	channelID, channelErr := workflow.ReactionParam["channel"]
	if !titleErr || !channelErr {
		return fmt.Errorf(errorMissingField)
	}

	url := "https://discord.com/api/v10/channels/" + channelID.(string) + "/threads"
	content := map[string]interface{}{
		"name": title,
		"type": 11,
	}

	jsonBody, err := json.Marshal(content)
	if err != nil {
		return err
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "POST", botBearer, tokenBot, bytes.NewBuffer(jsonBody))
	if errRequest != nil {
		return errRequest
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) checkDiscordReactions(workflow entities.Workflow) error {
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	switch reactionFound.Name {
	case "Post a message to a channel":
		return self.postDiscordMessage(tokenBot, workflow)
	case "Create a thread in a channel":
		return self.createDiscordThread(tokenBot, workflow)
	}

	return fmt.Errorf("Unknown reaction")
}
