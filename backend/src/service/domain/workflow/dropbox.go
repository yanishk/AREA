package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend/src/entities"
)

func dropboxReactions() []string {
	return []string{
		"Move file or folder",
		"Create a text file",
		"Append to a text file",
	}
}

func (self *WorkflowService) moveFileOrFolder(accessToken string, workflow entities.Workflow) error {
	path, pathExists := workflow.ReactionParam["path"]
	destination, destinationExists := workflow.ReactionParam["destination"]
	if !pathExists || !destinationExists {
		return fmt.Errorf(errorMissingField)
	}

	if path.(string)[0] != '/' {
		path = "/" + path.(string)
	}
	if destination.(string)[0] != '/' {
		destination = "/" + destination.(string)
	}

	url := "https://api.dropboxapi.com/2/files/move_v2"

	content := map[string]interface{}{
		"allow_ownership_transfer": false,
		"allow_shared_folder":      false,
		"autorename":               false,
		"from_path":                path,
		"to_path":                  destination,
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	res, err := self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, accessToken, bytes.NewBuffer(contentBytes))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) requestFileCreationModificationDropbox(url, accessToken, apiArg string, body io.Reader, workflow entities.Workflow) error {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", bearerType+accessToken)
	req.Header.Set(contentType, "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", apiArg)

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) createTextFile(accessToken string, workflow entities.Workflow) error {
	name, nameExists := workflow.ReactionParam["name"]
	content, contentExists := workflow.ReactionParam["content"]
	path, pathExists := workflow.ReactionParam["path"]
	if !nameExists || !contentExists || !pathExists {
		return fmt.Errorf(errorMissingField)
	}

	filePath := concatanateSlashToPath(path.(string), name.(string))
	url := "https://content.dropboxapi.com/2/files/upload"
	arg := fmt.Sprintf(`{"autorename":false,"mode":"add","mute":false,"path":"%s","strict_conflict":false}`, filePath)

	return self.requestFileCreationModificationDropbox(url, accessToken, arg, bytes.NewBuffer([]byte(content.(string))), workflow)
}

func (self *WorkflowService) downloadFileFromDropbox(accessToken, filePath string) (string, error) {
	url := "https://content.dropboxapi.com/2/files/download"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", bearerType+accessToken)
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf(`{"path": "%s"}`, filePath))

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	fileContent, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

func (self *WorkflowService) appendTextFile(accessToken string, workflow entities.Workflow) error {
	name, nameExists := workflow.ReactionParam["name"]
	content, contentExists := workflow.ReactionParam["content"]
	path, pathExists := workflow.ReactionParam["path"]
	if !nameExists || !contentExists || !pathExists {
		return fmt.Errorf(errorMissingField)
	}

	filePath := concatanateSlashToPath(path.(string), name.(string))
	contentAppend, err := self.downloadFileFromDropbox(accessToken, filePath)
	if err != nil {
		return err
	}

	contentAppend = contentAppend + content.(string)
	url := "https://content.dropboxapi.com/2/files/upload"
	arg := fmt.Sprintf(`{"autorename":false,"mode":"overwrite","mute":false,"path":"%s","strict_conflict":false}`, filePath)

	return self.requestFileCreationModificationDropbox(url, accessToken, arg, bytes.NewBuffer([]byte(contentAppend)), workflow)
}

func (self *WorkflowService) checkDropboxReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Dropbox", reactionFound.Name, dropboxReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	switch reactionFound.Name {
	case "Move file or folder":
		return self.moveFileOrFolder(accessToken, workflow)
	case "Create a text file":
		return self.createTextFile(accessToken, workflow)
	case "Append to a text file":
		return self.appendTextFile(accessToken, workflow)
	}
	return nil
}
