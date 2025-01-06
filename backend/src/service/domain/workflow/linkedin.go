package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend/src/entities"
)

func linkedinReactions() []string {
	return []string{
		"Share an update",
		"Share a link",
	}
}

func (self *WorkflowService) ExecuteLinkedinRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerType+accessToken)
	req.Header.Set(contentType, "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self *WorkflowService) getUserUrnLinkedin(accessToken string, workflow entities.Workflow) (string, error) {
	var sub entities.LinkedinSub
	url := "https://api.linkedin.com/v2/userinfo"

	res, err := self.ExecuteLinkedinRequest("GET", url, accessToken, nil)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	errDecode := json.NewDecoder(res.Body).Decode(&sub)
	if errDecode != nil {
		return "", errDecode
	}

	urn := "urn:li:person:" + sub.Sub
	return urn, nil
}

func (self *WorkflowService) publishTextLinkedinPost(urn, accessToken string, workflow entities.Workflow) error {
	message, messageErr := workflow.ReactionParam["message"]
	if !messageErr {
		return fmt.Errorf(errorMissingField)
	}

	url := "https://api.linkedin.com/v2/ugcPosts"
	content := map[string]interface{}{
		"author":         urn,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": message,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	jsonBody, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	res, errLinkedinRequest := self.ExecuteLinkedinRequest("POST", url, accessToken, bytes.NewBuffer([]byte(jsonBody)))
	if errLinkedinRequest != nil {
		return errLinkedinRequest
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) publishURLLinkedinPost(urn, accessToken string, workflow entities.Workflow) error {
	message, messageErr := workflow.ReactionParam["message"]
	urlPost, urlErr := workflow.ReactionParam["url"]
	if !messageErr || !urlErr {
		return fmt.Errorf(errorMissingField)
	}

	url := "https://api.linkedin.com/v2/ugcPosts"
	content := map[string]interface{}{
		"author":         urn,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareMediaCategory": "ARTICLE",
				"shareCommentary": map[string]interface{}{
					"text": message,
				},
				"media": []map[string]interface{}{
					{
						"status":      "READY",
						"originalUrl": urlPost,
					},
				},
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	jsonBody, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	_, errLinkedinRequest := self.ExecuteLinkedinRequest("POST", url, accessToken, bytes.NewBuffer([]byte(jsonBody)))
	if errLinkedinRequest != nil {
		return errLinkedinRequest
	}
	return nil
}

func (self *WorkflowService) checkLinkedinReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Linkedin", reactionFound.Name, linkedinReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	switch reactionFound.Name {
	case "Share an update":
		urn, err := self.getUserUrnLinkedin(accessToken, workflow)
		if err != nil {
			return err
		}
		return self.publishTextLinkedinPost(urn, accessToken, workflow)
	case "Share a link":
		urn, err := self.getUserUrnLinkedin(accessToken, workflow)
		if err != nil {
			return err
		}
		return self.publishURLLinkedinPost(urn, accessToken, workflow)
	}
	return nil
}
