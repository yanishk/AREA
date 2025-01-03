package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
)

type WorkflowService struct {
	WorkflowRepository storage.WorkflowRepository
	UserRepository     storage.UserRepository
	ActionRepository   storage.ActionRepository
	ReactionRepository storage.ReactionRepository
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
}

const botBearer = "Bot "
const bearerType = "Bearer "
const contentType = "Content-Type"

const errorUpdatingWorkflow = "Could not update workflow"
const errorUpdatingToken = "Could not update token"
const errorRetrievingReaction = "Error finding reaction"
const errorMissingField = "Missing required field"
const errorMarshaling = "Could not marshal JSON"

const redditUserRoute = "https://oauth.reddit.com/user/"
const githubBaseUrl = "https://api.github.com/"

func spotifyReactions() []string {
	return []string{
		"Start playback",
		"Pause playback",
		"Skip track",
		"Add track to playback queue",
		"Save a track",
	}
}

func linkedinReactions() []string {
	return []string{
		"Share an update",
		"Share a link",
	}
}

func asanaReactions() []string {
	return []string{
		"Create task",
		"Create project",
	}
}

func dropboxReactions() []string {
	return []string{
		"Move file or folder",
		"Create a text file",
		"Append to a text file",
	}
}

func redditReactions() []string {
	return []string{
		"Submit a comment on a post ",
		"Upvote a post",
		"Downvote a post",
		"Submit a post",
		"Submit a post with a link",
	}
}

func webhooksEventsToActions() map[string]string {
	return map[string]string{
		"Push Hook": "New push",
		"Merge Request Hook": "Merge request update",
		"Issue Hook": "Issue update",
		"Note Hook": "Comment update",
		"Tag Push Hook": "New tag push",
		"Wiki Page Hook": "Wiki page update",
		"Release Hook": "Release update",
		"Feature Flag Hook": "Feature flag update",
		"Pipeline Hook": "Pipeline update",
		"Job Hook": "Job update",
		"Deployment Hook": "Deployment update",
		"Emoji Hook": "Emoji update",
	}
}

func gitlabWebhookEventsToState() map[string]string {
	return map[string]string{
		"push_events": "true",
		"issues_events": "true",
		"merge_requests_events": "true",
		"tag_push_events": "true",
		"note_events": "true",
		"job_events": "true",
		"pipeline_events": "true",
		"wiki_page_events": "true",
		"deployment_events": "true",
		"releases_events": "true",
		"feature_flag_events": "true",
		"emoji_events": "true",
	}
}

func NewWorkflowService(WorkflowRepository storage.WorkflowRepository, UserRepository storage.UserRepository,
	ActionRepository storage.ActionRepository, ReactionRepository storage.ReactionRepository,
	ServiceService service.ServiceService, UserServiceService service.UserServiceService) *WorkflowService {
	return &WorkflowService{
		WorkflowRepository: WorkflowRepository,
		UserRepository:     UserRepository,
		ActionRepository:   ActionRepository,
		ReactionRepository: ReactionRepository,
		ServiceService:     ServiceService,
		UserServiceService: UserServiceService,
	}
}

func (self *WorkflowService) CreateWorkflow(userEmail, userConnectionType string, newWorkflow entities.NewWorkflow) error {
	userFound, errFindingUser := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if errFindingUser != nil {
		return errFindingUser
	}
	errCreationWorkflow := self.WorkflowRepository.CreateWorkflow(newWorkflow.Name,
		userFound.Id, newWorkflow.ActionId, newWorkflow.ReactionId,
		newWorkflow.ActionParam, newWorkflow.ReactionParam, newWorkflow.ActionData)
	if errCreationWorkflow != nil {
		return errCreationWorkflow
	}
	return nil
}

func (self *WorkflowService) GetUserWorkflows(email, connectionType string) ([]entities.Workflow, error) {
	userFound, errUserFound := self.UserRepository.FindUserByEmail(email, connectionType)
	if errUserFound != nil {
		return nil, errUserFound
	}

	retrievedWorkflow, err := self.WorkflowRepository.FindWorkflowsByOwnerId(userFound.Id)
	if err != nil {
		return nil, err
	}
	return retrievedWorkflow, nil
}

func (self *WorkflowService) UpdateWorkflow(workflowId string, workflow entities.UpdatedWorkflow) error {
	updatedWorkflow, err := self.WorkflowRepository.FindWorkflowById(workflowId)
	if err != nil {
		return err
	}

	if workflow.Name != nil {
		updatedWorkflow.Name = *workflow.Name
	}
	if workflow.ActionId != nil {
		updatedWorkflow.ActionId = *workflow.ActionId
	}
	if workflow.ReactionId != nil {
		updatedWorkflow.ReactionId = *workflow.ReactionId
	}
	if workflow.IsActivated != nil {
		updatedWorkflow.IsActivated = *workflow.IsActivated
	}
	if workflow.ActionParam != nil {
		updatedWorkflow.ActionParam = *workflow.ActionParam
	}
	if workflow.ReactionParam != nil {
		updatedWorkflow.ReactionParam = *workflow.ReactionParam
	}

	err = self.WorkflowRepository.UpdateWorkflow(workflowId, updatedWorkflow)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowService) DeleteWorkflow(email, connectionType, workflowId string) error {
	user, err := self.UserRepository.FindUserByEmail(email, connectionType)
	if err != nil {
		return err
	}

	err = self.WorkflowRepository.DeleteWorkflow(workflowId, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowService) refreshTokenForService(serviceName, reactionFoundName string, reactionsPossible []string, workflow entities.Workflow) (string, error) {
	var accessToken string
	var err error

	for i := range reactionsPossible {
		if reactionFoundName == reactionsPossible[i] {
			accessToken, err = self.getAccessToken(serviceName, workflow)
			if err != nil {
				return accessToken, fmt.Errorf(errorUpdatingToken)
			}
		}
	}
	return accessToken, nil
}

func (self *WorkflowService) checkSpotifyReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Spotify", reactionFound.Name, spotifyReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	var resp *http.Response
	switch reactionFound.Name {
	case "Start playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/play", "PUT", bearerType, accessToken, nil)
	case "Pause playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", bearerType, accessToken, nil)
	case "Skip track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", bearerType, accessToken, nil)
	case "Add track to playback queue":
		uri, uriExists := workflow.ReactionParam["uri"]
		if !uriExists {
			return fmt.Errorf("Incorrect URI")
		}

		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlQueue, "POST", bearerType, accessToken, nil)
	case "Save a track":
		id, idExists := workflow.ReactionParam["id"]
		if !idExists {
			return fmt.Errorf("Incorrect ID")
		}

		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlSave, "PUT", bearerType, accessToken, nil)
	}

	if resp != nil {
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("Unknown reaction")
}

func (self *WorkflowService) postDiscordMessage(tokenBot string, workflow entities.Workflow) error {
	message, messageErr := workflow.ReactionParam["message"]
	channelID, channelErr := workflow.ReactionParam["channelid"]
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
	channelID, channelErr := workflow.ReactionParam["channelid"]
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
		return fmt.Errorf(errorMarshaling, err)
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

func concatanateSlashToPath(path, name string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	return path + name
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

func (self *WorkflowService) postRedditComment(accessToken string, workflow entities.Workflow) error {
	comment, commentExists := workflow.ReactionParam["comment"]
	id, idExists := workflow.ReactionParam["id"]
	if !commentExists || !idExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("thing_id", "t3_"+id.(string))
	body.Set("text", comment.(string))
	url := "https://oauth.reddit.com/api/comment"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) upvoteDowvoteRedditPost(postId, accessToken string, vote int) error {
	body := url.Values{}
	body.Set("id", "t3_"+postId)
	body.Set("dir", strconv.Itoa(vote))

	url := "https://oauth.reddit.com/api/vote"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) upvoteRedditPost(accessToken string, workflow entities.Workflow) error {
	id, idExists := workflow.ReactionParam["id"]
	if !idExists {
		return fmt.Errorf(errorMissingField)
	}

	return self.upvoteDowvoteRedditPost(id.(string), accessToken, 1)
}

func (self *WorkflowService) downvoteRedditPost(accessToken string, workflow entities.Workflow) error {
	id, idExists := workflow.ReactionParam["id"]
	if !idExists {
		return fmt.Errorf(errorMissingField)
	}

	return self.upvoteDowvoteRedditPost(id.(string), accessToken, -1)
}

func (self *WorkflowService) submitRedditPost(accessToken string, workflow entities.Workflow) error {
	title, titleExists := workflow.ReactionParam["title"]
	content, contentExists := workflow.ReactionParam["content"]
	subreddit, subredditExists := workflow.ReactionParam["subreddit"]
	if !titleExists || !contentExists || !subredditExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("title", title.(string))
	body.Set("text", content.(string))
	body.Set("sr", subreddit.(string))
	body.Set("kind", "self")
	url := "https://oauth.reddit.com/api/submit"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) submitRedditPostLink(accessToken string, workflow entities.Workflow) error {
	title, titleExists := workflow.ReactionParam["title"]
	link, linkExists := workflow.ReactionParam["link"]
	subreddit, subredditExists := workflow.ReactionParam["subreddit"]
	if !titleExists || !linkExists || !subredditExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("title", title.(string))
	body.Set("url", link.(string))
	body.Set("sr", subreddit.(string))
	body.Set("kind", "link")
	url := "https://oauth.reddit.com/api/submit"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) checkRedditReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Reddit", reactionFound.Name, redditReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	switch reactionFound.Name {
	case "Submit a comment on a post ":
		return self.postRedditComment(accessToken, workflow)
	case "Downvote a post":
		return self.downvoteRedditPost(accessToken, workflow)
	case "Upvote a post":
		return self.upvoteRedditPost(accessToken, workflow)
	case "Submit a post":
		return self.submitRedditPost(accessToken, workflow)
	case "Submit a post with a link":
		return self.submitRedditPostLink(accessToken, workflow)
	}
	return nil
}

func (self *WorkflowService) checkReactions(workflow entities.Workflow) {
	self.checkSpotifyReactions(workflow)
	self.checkDiscordReactions(workflow)
	self.checkLinkedinReactions(workflow)
	self.checkAsanaReactions(workflow)
	self.checkSMSReactions(workflow)
	self.checkSendEmailReactions(workflow)
	self.checkDropboxReactions(workflow)
	self.checkRedditReactions(workflow)
}

func checkEveryHourParamValidity(minute interface{}) (float64, error) {
	minuteNumber, minuteIsNumber := minute.(float64)
	if !minuteIsNumber {
		return minuteNumber, fmt.Errorf("Error in the minute sent")
	}
	return minuteNumber, nil
}

func checkEveryDayParamValidity(hour interface{}, minute interface{}) (float64, float64, error) {
	var hourNumber, minuteNumber float64

	hourNumber, hourIsNumber := hour.(float64)
	if !hourIsNumber {
		return hourNumber, minuteNumber, fmt.Errorf("Error in the hour sent")
	}

	minuteNumber, err := checkEveryHourParamValidity(minute)
	if err != nil {
		return hourNumber, minuteNumber, err
	}

	return hourNumber, minuteNumber, nil
}

func checkEveryMonthOnTheParamValidity(day interface{}, hour interface{}, minute interface{}) (float64, float64, float64, error) {
	var dayNumber, hourNumber, minuteNumber float64

	dayNumber, dayIsNumber := day.(float64)
	if !dayIsNumber {
		return dayNumber, hourNumber, minuteNumber, fmt.Errorf("Error in the day sent")
	}

	hourNumber, minuteNumber, err := checkEveryDayParamValidity(hour, minute)
	if err != nil {
		return dayNumber, hourNumber, minuteNumber, err
	}

	return dayNumber, hourNumber, minuteNumber, nil
}

func checkEveryWeekAtParamValidity(day interface{}, hour interface{}, minute interface{}) (float64, float64, float64, error) {
	return checkEveryMonthOnTheParamValidity(day, hour, minute)
}

func checkEveryYearParamValidity(month interface{}, day interface{}, hour interface{}, minute interface{}) (float64, float64, float64, float64, error) {
	var monthNumber, dayNumber, hourNumber, minuteNumber float64

	monthNumber, monthIsNumber := month.(float64)
	if !monthIsNumber {
		return monthNumber, dayNumber, hourNumber, minuteNumber, fmt.Errorf("Error in the month sent")
	}

	dayNumber, hourNumber, minuteNumber, err := checkEveryMonthOnTheParamValidity(day, hour, minute)
	if err != nil {
		return monthNumber, dayNumber, hourNumber, minuteNumber, err
	}

	return monthNumber, dayNumber, hourNumber, minuteNumber, nil
}

func (self *WorkflowService) checkTimeAndDateEveryHourAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	minuteNumber, err := checkEveryHourParamValidity(minute)
	if err != nil {
		return err
	}

	if int(minuteNumber) == timeRes.Minute && workflow.IsActivated {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkTimeAndDateEveryDayAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !hourExists || !minuteExists {
		return fmt.Errorf(errorMissingField)
	}

	hourNumber, minuteNumber, err := checkEveryDayParamValidity(hour, minute)
	if err != nil {
		return err
	}

	if int(hourNumber) == timeRes.Hour && int(minuteNumber) == timeRes.Minute {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkTimeAndDateEveryWeekAtAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	days, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf(errorMissingField)
	}

	dayArray, ok := days.([]interface{})
	if !ok {
		return fmt.Errorf("Error parsing the day array")
	}

	for _, day := range dayArray {
		dayNumber, hourNumber, minuteNumber, err := checkEveryWeekAtParamValidity(day, hour, minute)
		if err != nil {
			return err
		}

		if int(dayNumber) == timeRes.WeekDay && int(hourNumber) == timeRes.Hour &&
			int(minuteNumber) == timeRes.Minute {
			self.checkReactions(workflow)
		}
	}

	return nil
}

func (self *WorkflowService) checkTimeAndDateEveryMonthOnTheAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf(errorMissingField)
	}

	dayNumber, hourNumber, minuteNumber, err := checkEveryMonthOnTheParamValidity(day, hour, minute)
	if err != nil {
		return err
	}

	if int(dayNumber) == timeRes.Day && int(hourNumber) == timeRes.Hour &&
		int(minuteNumber) == timeRes.Minute {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkTimeAndDateEveryYearOnAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	month, monthExists := workflow.ActionParam["month"]
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists || !monthExists {
		return fmt.Errorf(errorMissingField)
	}

	monthNumber, dayNumber, hourNumber, minuteNumber, err := checkEveryYearParamValidity(month, day, hour, minute)
	if err != nil {
		return err
	}

	if int(monthNumber) == timeRes.Month && int(dayNumber) == timeRes.Day && int(hourNumber) == timeRes.Hour &&
		int(minuteNumber) == timeRes.Minute {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkWorkflowsWithTimeAndDateActions(timeRes entities.TimeResponse, action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if !workflow.IsActivated {
			continue
		}

		switch action.Name {
		case "Every day at":
			self.checkTimeAndDateEveryDayAction(timeRes, workflow)
		case "Every hour at":
			self.checkTimeAndDateEveryHourAction(timeRes, workflow)
		case "Every day of the week at":
			self.checkTimeAndDateEveryWeekAtAction(timeRes, workflow)
		case "Every month on the":
			self.checkTimeAndDateEveryMonthOnTheAction(timeRes, workflow)
		case "Every year on":
			self.checkTimeAndDateEveryYearOnAction(timeRes, workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckTimeAndDateActions() error {
	var timeRes entities.TimeResponse

	timeRes, errRequest := self.ServiceService.RequestToTimeApi()
	if errRequest != nil {
		return errRequest
	}

	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Time & Date")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionsByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkWorkflowsWithTimeAndDateActions(timeRes, action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) executeRedditRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", bearerType+accessToken)
	req.Header.Add("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	resp, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (self *WorkflowService) getRedditPost(url, accessToken string) (entities.RedditPostResponse, error) {
	var result entities.RedditPostResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve post")
	}
	return result, nil
}

func (self *WorkflowService) getRedditComments(url, accessToken string) (entities.RedditCommentResponse, error) {
	var result entities.RedditCommentResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve comments")
	}
	return result, nil
}

func (self *WorkflowService) getRedditVotes(url, accessToken string) (entities.RedditVoteResponse, error) {
	var result entities.RedditVoteResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve votes")
	}
	return result, nil
}

func (self *WorkflowService) getUsernameReddit(accessToken string, workflow entities.Workflow) (entities.RedditUsername, error) {
	var result entities.RedditUsername
	url := "https://oauth.reddit.com/api/v1/me"

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}
	return result, nil
}

func (self *WorkflowService) isANewPost(workflow entities.Workflow, result entities.RedditPostResponse) error {
	latestPost := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestPost.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestPost.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestPost.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) isANewComment(workflow entities.Workflow, result entities.RedditCommentResponse) error {
	latestComment := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestComment.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestComment.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestComment.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) isANewVote(workflow entities.Workflow, result entities.RedditVoteResponse) error {
	latestVote := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestVote.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestVote.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestVote.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) getAccessToken(serviceName string, workflow entities.Workflow) (string, error) {
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errUser != nil {
		return "", fmt.Errorf("Error finding user")
	}
	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, serviceName)
	if errRefreshing != nil {
		return "", errRefreshing
	}
	return accessToken, nil
}

func (self *WorkflowService) checkRedditNewPostInSubredditAction(accessToken string, workflow entities.Workflow) error {
	subreddit, subredditExists := workflow.ActionParam["subreddit"]
	if !subredditExists {
		return fmt.Errorf("Could not find the subreddit field")
	}

	url := "https://oauth.reddit.com/" + subreddit.(string) + "/new.json?limit=1"
	result, err := self.getRedditPost(url, accessToken)
	if err != nil {
		return err
	}

	if len(result.Data.Children) == 0 {
		return fmt.Errorf("No posts found in subreddit")
	}

	return self.isANewPost(workflow, result)
}

func (self *WorkflowService) checkRedditNewCommentByMeAction(accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	url := redditUserRoute + resultUsername.Name + "/comments?limit=1"
	result, err := self.getRedditComments(url, accessToken)
	if err != nil {
		return err
	}

	return self.isANewComment(workflow, result)
}

func (self *WorkflowService) checkRedditVoteByMeAction(voteType, accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	urlPostMe := redditUserRoute + resultUsername.Name + voteType + "?limit=1"
	resultPost, err := self.getRedditVotes(urlPostMe, accessToken)
	if err != nil {
		return err
	}

	return self.isANewVote(workflow, resultPost)
}

func (self *WorkflowService) checkRedditPostAction(postType, accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	url := redditUserRoute + resultUsername.Name + postType + "?limit=1"
	result, err := self.getRedditPost(url, accessToken)
	if err != nil {
		return err
	}

	if len(result.Data.Children) == 0 {
		return fmt.Errorf("No posts found")
	}

	return self.isANewPost(workflow, result)
}

func (self *WorkflowService) checkWorkflowsWithRedditActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if !workflow.IsActivated {
			continue
		}

		accessToken, err := self.getAccessToken("Reddit", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "Any new post in subreddit":
			self.checkRedditNewPostInSubredditAction(accessToken, workflow)
		case "New post by you":
			self.checkRedditPostAction("/submitted", accessToken, workflow)
		case "New comment by you":
			self.checkRedditNewCommentByMeAction(accessToken, workflow)
		case "New downvoted post by you":
			self.checkRedditVoteByMeAction("/downvoted", accessToken, workflow)
		case "New upvoted post by you":
			self.checkRedditVoteByMeAction("/upvoted", accessToken, workflow)
		case "New post saved by you":
			self.checkRedditPostAction("/saved", accessToken, workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckRedditActions() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Reddit")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionsByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkWorkflowsWithRedditActions(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) currentWeather(city string, workflow entities.Workflow) (entities.WeatherResponse, error) {
	var weatherData entities.WeatherResponse
	baseUrl := "http://api.weatherapi.com/v1/forecast.json?key=" + os.Getenv("WEATHER_API_KEY") + "&q=" + city

	req, _ := http.NewRequest("GET", baseUrl, nil)
	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return weatherData, err
	}
	defer res.Body.Close()

	errDecode := json.NewDecoder(res.Body).Decode(&weatherData)
	if errDecode != nil {
		return weatherData, errDecode
	}
	return weatherData, nil
}

func (self *WorkflowService) checkWeatherCurrentWeatherComparisonAction(checkType string, workflow entities.Workflow) error {
	city, cityExists := workflow.ActionParam["city"]
	temp, tempExists := workflow.ActionParam["temperature"]
	if !cityExists || !tempExists {
		return fmt.Errorf("Could not retrieve the necessary fields")
	}

	weatherData, err := self.currentWeather(city.(string), workflow)
	if err != nil {
		return err
	}

	if checkType == "aboveCurrent" && weatherData.Current.Temperature > temp.(float64) {
		self.checkReactions(workflow)
	} else if checkType == "belowCurrent" && weatherData.Current.Temperature < temp.(float64) {
		self.checkReactions(workflow)
	} else if checkType == "aboveForecast" && weatherData.Forecast.ForecastDay[0].Day.MaxTemperature < temp.(float64) {
		self.checkReactions(workflow)
	} else if checkType == "belowForecast" && weatherData.Forecast.ForecastDay[0].Day.MinTemperature > temp.(float64) {
		self.checkReactions(workflow)
	} else {
		fmt.Errorf("Check type isn't valid")
	}
	return nil
}

func (self *WorkflowService) checkWorkflowsWithWeatherActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if !workflow.IsActivated {
			continue
		}

		switch action.Name {
		case "Current temperature rises above":
			self.checkWeatherCurrentWeatherComparisonAction("aboveCurrent", workflow)
		case "Current temperature drops below":
			self.checkWeatherCurrentWeatherComparisonAction("belowCurrent", workflow)
		case "Tomorrow's low drops below":
			self.checkWeatherCurrentWeatherComparisonAction("belowForecast", workflow)
		case "Tomorrow's high rises above":
			self.checkWeatherCurrentWeatherComparisonAction("aboveForecast", workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckWeatherActions() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("FreeWeather")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionsByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkWorkflowsWithWeatherActions(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) executeGithubRequest(workflow entities.Workflow, accessToken, endpoint string) (*http.Response, error) {
	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return nil, fmt.Errorf("Invalid repository")
	}

	url := githubBaseUrl + "repos/" + repository.(string) + endpoint
	return self.ServiceService.ExecuteApiRequest(url, "GET", bearerType, accessToken, nil)
}

func (self *WorkflowService) checkActionDataLen(workflow entities.Workflow, lenActualTurn float64) error {
	workflow.ActionData = make(map[string]interface{})
	lenLastTurn, lenLastTurnExists := workflow.ActionData["len"]
	if !lenLastTurnExists {
		workflow.ActionData["len"] = lenActualTurn
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
		return nil
	}

	if lenLastTurn.(float64) < lenActualTurn {
		workflow.ActionData["len"] = lenActualTurn
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkGithubNewRepositoryAction(workflow entities.Workflow, accessToken string) error {
	repositories, err := self.ServiceService.RequestGithubUserRepositories(accessToken)
	if err != nil {
		return err
	}
	return self.checkActionDataLen(workflow, float64(len(repositories)))
}

func (self *WorkflowService) checkGithubNewIssueAction(workflow entities.Workflow, accessToken string) error {
	var issues []entities.GithubIssue

	resp, err := self.ServiceService.ExecuteApiRequest(githubBaseUrl+"issues", "GET", bearerType, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(issues)))
}

func (self *WorkflowService) checkGithubNewPullRequestAction(workflow entities.Workflow, accessToken string) error {
	var pullRequests []entities.GithubRepository

	resp, err := self.executeGithubRequest(workflow, accessToken, "/pulls")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(pullRequests)))
}

func (self *WorkflowService) checkGithubNewBranchAction(workflow entities.Workflow, accessToken string) error {
	var branches []entities.GithubBranch

	resp, err := self.executeGithubRequest(workflow, accessToken, "/branches")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&branches)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(branches)))
}

func (self *WorkflowService) checkGithubNewCommitAction(workflow entities.Workflow, accessToken string) error {
	var commits []entities.GithubCommit

	resp, err := self.executeGithubRequest(workflow, accessToken, "/commits")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(commits)))
}

func (self *WorkflowService) checkWorkflowsWithGithubActions(action entities.Action) error {
	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		if !workflow.IsActivated {
			continue
		}
		accessToken, err := self.getAccessToken("Github", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "New repository":
			self.checkGithubNewRepositoryAction(workflow, accessToken)
		case "New issue assignated":
			self.checkGithubNewIssueAction(workflow, accessToken)
		case "New pull request":
			self.checkGithubNewPullRequestAction(workflow, accessToken)
		case "New branch":
			self.checkGithubNewBranchAction(workflow, accessToken)
		case "New commit":
			self.checkGithubNewCommitAction(workflow, accessToken)
		}
	}
	return nil
}

func (self *WorkflowService) CheckGithubActions() error {
	service, err := self.ServiceService.FindServiceByName("Github")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionsByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.checkWorkflowsWithGithubActions(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGitlabIdActionParam(workflow entities.Workflow) (string, error) {
	projectId, projectIdExists := workflow.ActionParam["project"]
	if !projectIdExists {
		return "", fmt.Errorf("Error in the id sent")
	}

	projectIdString, projectIdIsString := projectId.(string)
	if !projectIdIsString {
		return "", fmt.Errorf("Error in the project id sent")
	}
	return projectIdString, nil
}

func (self *WorkflowService) checkGitlabWebhooksWorkflow(workflow entities.Workflow, webhookProjectId int) error {
	webhookProjectIdString := strconv.Itoa(webhookProjectId)
	userProjectIdString, err := getGitlabIdActionParam(workflow)
	if err != nil {
		return err
	}

	if (userProjectIdString == webhookProjectIdString) {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkGitlabWebhooksWorkflows(headers http.Header, webhookJsonDataBytes []byte, eventsToActions map[string]string) error {
	var webhookResponse entities.GitlabWebhookResponse
	if len(webhookJsonDataBytes) > 0 {
		err := json.Unmarshal(webhookJsonDataBytes, &webhookResponse)
		if err != nil {
			return err
		}
	}
	if (webhookResponse.Project.Id == 0) {
		return fmt.Errorf("Incorrect project id")
	}

	eventName := headers["X-Gitlab-Event"][0]
	actionName, actionNameExists := eventsToActions[eventName]
	if !actionNameExists {
		return fmt.Errorf("No action mapped with this event")
	}

	action, err := self.ActionRepository.FindActionByName(actionName)
	if err != nil {
		return err
	}

	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		if !workflow.IsActivated {
			continue
		}
		self.checkGitlabWebhooksWorkflow(workflow, webhookResponse.Project.Id)
	}
	return nil
}

func unmarshalJsonToMap(jsonDataBytes []byte) ([]map[string]interface{}, error) {
	var jsonData []map[string]interface{}
	if len(jsonDataBytes) > 0 {
		err := json.Unmarshal(jsonDataBytes, &jsonData)
		if err != nil {
			return jsonData, err
		}
	}
	return jsonData, nil
}

func (self *WorkflowService) CheckWebhooksWorkflows(serviceName string, request *http.Request) error {
	webhookJsonDataBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	switch serviceName {
	case "Gitlab":
        return self.checkGitlabWebhooksWorkflows(request.Header, webhookJsonDataBytes, webhooksEventsToActions())
    }
	return nil
}

func (self *WorkflowService) ExecuteGitlabRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerType+accessToken)
	req.Header.Set(contentType, "application/x-www-form-urlencoded")

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self*WorkflowService) createNewWorkflowGitlabWebhook(projectId, accessToken string) error {
	createWebhookUrl := "https://gitlab.com/api/v4/projects/" + projectId + "/hooks"
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Gitlab"
	params := url.Values{}

	params.Set("url", webhookUrl)
	for webhookEvent, webhookState := range gitlabWebhookEventsToState() {
		params.Set(webhookEvent, webhookState)
	}

	res, err := self.ExecuteGitlabRequest("POST", createWebhookUrl, accessToken, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) isProjectWebhookPresent(projectId, accessToken string) (bool, error) {
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Gitlab"
	getWebhooksUrl := "https://gitlab.com/api/v4/projects/" + projectId + "/hooks"

	res, err := self.ExecuteGitlabRequest("GET", getWebhooksUrl, accessToken, nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	webhooksJsonDataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	webhooksJsonData, err := unmarshalJsonToMap(webhooksJsonDataBytes)
	if err != nil {
		return false, err
	}

	if (len(webhooksJsonData) == 0) {
		return false, nil
	}

	for _, webhookJsonData := range webhooksJsonData {
		url, urlExists := webhookJsonData["url"]
		if !urlExists {
			return false, fmt.Errorf("Error in the url received")
		}

		urlString, urlIsString := url.(string)
		if !urlIsString {
			return false, fmt.Errorf("Error in the url received")
		}

		if (urlString == webhookUrl) {
			return true, nil
		}
	}

	return false, nil
}

func (self *WorkflowService) checkNewWorkflowGitlabWebhook(workflow entities.Workflow, accessToken string) error {
	projectIdString, err := getGitlabIdActionParam(workflow)
	if err != nil {
		return err
	}

	isWebhookPresent, err := self.isProjectWebhookPresent(projectIdString, accessToken)
	if err != nil {
		return err
	}

	if !isWebhookPresent {
		self.createNewWorkflowGitlabWebhook(projectIdString, accessToken)
		workflow.ActionData = make(map[string]interface{})
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WorkflowService) checkNewWorkflowsGitlabWebhook(action entities.Action) error {
	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		accessToken, err := self.getAccessToken("Gitlab", workflow)
		if err != nil {
			continue
		}
		self.checkNewWorkflowGitlabWebhook(workflow, accessToken)
	}
	return nil
}

func (self *WorkflowService) CheckNewGitlabWorkflows() error {
	service, err := self.ServiceService.FindServiceByName("Gitlab")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionsByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.checkNewWorkflowsGitlabWebhook(action)
		if err != nil {
			return err
		}
	}
	return nil
}
