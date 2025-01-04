package workflow_service

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"backend/src/entities"
)

func (self *WorkflowService) sendAnSMS(workflow entities.Workflow) error {
	receiver, receiverExists := workflow.ReactionParam["to"]
	body, bodyExists := workflow.ReactionParam["body"]
	if !receiverExists || !bodyExists {
		return fmt.Errorf(errorMissingField)
	}

	urlBase := "https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("SMS_ACCOUNT_SID") + "/Messages.json"

	data := url.Values{}
	data.Set("To", receiver.(string))
	data.Set("From", os.Getenv("TWILIO_VIRT_NUMBER"))
	data.Set("Body", body.(string))

	req, err := http.NewRequest("POST", urlBase, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set(contentType, "application/json")
	req.SetBasicAuth(os.Getenv("SMS_ACCOUNT_SID"), os.Getenv("AUTH_TOKEN"))

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) checkSMSReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	switch reactionFound.Name {
	case "Send an SMS":
		return self.sendAnSMS(workflow)
	}
	return nil
}
