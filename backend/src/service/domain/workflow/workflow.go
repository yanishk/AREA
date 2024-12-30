package workflow_service

import (
	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
const accessBearer = "Bearer "

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

func (self *WorkflowService) CreateWorkflow(userEmail string, userConnectionType string, newWorkflow entities.NewWorkflow) error {
	userFound, errFindingUser := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if errFindingUser != nil {
		return errFindingUser
	}
	errCreationWorkflow := self.WorkflowRepository.CreateWorkflow(newWorkflow.Name,
		userFound.Id, newWorkflow.ActionId, newWorkflow.ReactionId,
		newWorkflow.ActionParam, newWorkflow.ReactionParam)
	if errCreationWorkflow != nil {
		return errCreationWorkflow
	}
	return nil
}

func (self *WorkflowService) GetUserWorkflows(email string, connectionType string) ([]entities.Workflow, error) {
	userFound, errUserFound := self.UserRepository.FindUserByEmail(email, connectionType)
	if errUserFound != nil {
		return nil, errUserFound
	}

	retrievedWorkflow, err := self.WorkflowRepository.FindWorkflowByUserId(userFound.Id)
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

func (self *WorkflowService) DeleteWorkflow(email string, connectionType string, workflowId string) error {
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

func (self *WorkflowService) executeSpotifyReactions(workflow entities.Workflow) error {

	reactionFound, errReaction := self.ReactionRepository.FindReactionByID(workflow.ReactionId)
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errReaction != nil || errUser != nil {
		return fmt.Errorf("Error finding reaction")
	}
	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, "Spotify")
	if errRefreshing != nil {
		return errRefreshing
	}

	var resp *http.Response
	if reactionFound.Name == "Start playback" {
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/play", "PUT", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Pause playback" {
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Skip track" {
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Add track to playback queue" {
		uri, uriErr := workflow.ReactionParam["uri"]
		if !uriErr {
			return fmt.Errorf("Incorrect URI")
		}
		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlQueue, "POST", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Save a track" {
		id, idErr := workflow.ReactionParam["id"]
		if !idErr {
			return fmt.Errorf("Incorrect ID")
		}
		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlSave, "PUT", accessBearer, accessToken, nil)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	return fmt.Errorf("Unknown reaction")
}

func (self *WorkflowService) executeDiscordReactions(workflow entities.Workflow) error {
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	reactionFound, errReaction := self.ReactionRepository.FindReactionByID(workflow.ReactionId)
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errReaction != nil || errUser != nil {
		return fmt.Errorf("Error finding reaction")
	}

	_, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, "Discord")
	if errRefreshing != nil {
		return errRefreshing
	}

	if reactionFound.Name == "Post a message to a channel" {
		message, messageErr := workflow.ReactionParam["message"]
		channelID, channelErr := workflow.ReactionParam["channelid"]
		if !messageErr || !channelErr {
			return fmt.Errorf("Failed to get message body or channel id")
		}

		urlDiscordAuth := "https://discord.com/api/v10/channels/" + channelID.(string) + "/messages"
		content, errContent := json.Marshal(map[string]string{
			"content": message.(string),
		})
		if errContent != nil {
			return fmt.Errorf("Failed to marshall message content")
		}

		resp, errRequest := self.ServiceService.ExecuteApiRequest(urlDiscordAuth, "POST", botBearer, tokenBot, bytes.NewBuffer(content))
		if errRequest != nil {
			return errRequest
		}
		defer resp.Body.Close()
	}
	return fmt.Errorf("Unknown reaction")
}

func (self *WorkflowService) executeActiveActionsEveryDay(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !hourExists || !minuteExists {
		return fmt.Errorf("Error in the hour and minute sent")
	}

	if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) executeActiveActionsEveryHour(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	if int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) executeActiveActionsEveryWeekAt(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	days, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf("Error in the minute, hour or day sent")
	}

	dayArray, ok := days.([]interface{})
	if !ok {
		return fmt.Errorf("Error parsing the day array")
	}

	for _, day := range dayArray {
		if int(day.(float64)) == timeRes.WeekDay && int(hour.(float64)) == timeRes.Hour &&
			int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
			self.executeSpotifyReactions(workflow)
			self.executeDiscordReactions(workflow)
		}
	}

	return nil
}

func (self *WorkflowService) executeActiveActionsEveryMonthOnThe(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf("Error in the minute, hour or day sent")
	}

	if int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) executeActiveActionsEveryYearOn(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	month, monthExists := workflow.ActionParam["month"]
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists || !monthExists {
		return fmt.Errorf("Error in the minute, hour, day or month sent")
	}

	if int(month.(float64)) == timeRes.Month && int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkNameOfAction(timeRes entities.TimeResponse, action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if action.Name == "Every day at" {
			self.executeActiveActionsEveryDay(timeRes, workflow)
		}
		if action.Name == "Every hour at" {
			self.executeActiveActionsEveryHour(timeRes, workflow)
		}
		if action.Name == "Every day of the week at" {
			self.executeActiveActionsEveryWeekAt(timeRes, workflow)
		}
		if action.Name == "Every month on the" {
			self.executeActiveActionsEveryMonthOnThe(timeRes, workflow)
		}
		if action.Name == "Every year on" {
			self.executeActiveActionsEveryYearOn(timeRes, workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CallTimeApi() error {
	var timeRes entities.TimeResponse

	timeRes, errRequest := self.ServiceService.RequestToTimeApi()
	if errRequest != nil {
		return errRequest
	}

	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Time & Date")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingaction := self.ActionRepository.FindActionByServiceId(serviceFound.Id)
	if errFindingaction != nil {
		return errFindingaction
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkNameOfAction(timeRes, action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

var lenRepo int

func (self *WorkflowService) CheckGithubNewRepositoryAction(workflow entities.Workflow, accessToken string) (bool, error) {
	var isActionTriggered bool

	resp, err := self.ServiceService.ExecuteApiRequest("https://api.github.com/user/repos", "GET", accessBearer, accessToken, nil)
	if err != nil {
		return isActionTriggered, err
	}

	var repositories []entities.GithubRepository
	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		return isActionTriggered, err
	}

	if lenRepo < len(repositories)  {
		lenRepo = len(repositories)
		isActionTriggered = true
	} else {
		lenRepo = len(repositories)
	}
	return isActionTriggered, nil
}

var lenPullRequests int

func (self *WorkflowService) CheckGithubNewPullRequestAction(workflow entities.Workflow, accessToken string) (bool, error) {
	var isActionTriggered bool

	repository, ok := workflow.ActionParam["repository"]
	if !ok {
		return isActionTriggered, fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "pulls"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return isActionTriggered, err
	}

	var pullRequests []entities.GithubRepository
	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		return isActionTriggered, err
	}

	if lenPullRequests < len(pullRequests) {
		lenPullRequests = len(pullRequests)
		isActionTriggered = true
	} else {
		lenPullRequests = len(pullRequests)
	}
	return isActionTriggered, nil
}

var lenBranches int

func (self *WorkflowService) CheckGithubNewBranchAction(workflow entities.Workflow, accessToken string) (bool, error) {
	var isActionTriggered bool

	repository, ok := workflow.ActionParam["repository"]
	if !ok {
		return isActionTriggered, fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "branches"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return isActionTriggered, err
	}

	var branches []entities.GithubBranch
	err = json.NewDecoder(resp.Body).Decode(&branches)
	if err != nil {
		return isActionTriggered, err
	}

	if lenBranches < len(branches) {
		lenBranches = len(branches)
		isActionTriggered = true
		fmt.Println("new branch!")
	} else {
		lenBranches = len(branches)
	}
	return isActionTriggered, nil
}

var lenCommits int

func (self *WorkflowService) CheckGithubNewCommitAction(workflow entities.Workflow, accessToken string) (bool, error) {
	var isActionTriggered bool

	repository, ok := workflow.ActionParam["repository"]
	if !ok {
		return isActionTriggered, fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "commits"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return isActionTriggered, err
	}

	var commits []entities.GithubCommit
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return isActionTriggered, err
	}

	if lenCommits < len(commits) {
		lenCommits = len(commits)
		isActionTriggered = true
		fmt.Println("new commit!")
	} else {
		lenCommits = len(commits)
	}
	return isActionTriggered, nil
}

var lenIssues int

func (self *WorkflowService) CheckGithubNewIssueAction(workflow entities.Workflow, accessToken string) (bool, error) {
	var isActionTriggered bool

	resp, err := self.ServiceService.ExecuteApiRequest("https://api.github.com/issues", "GET", accessBearer, accessToken, nil)
	if err != nil {
		return isActionTriggered, err
	}

	var issues []entities.GithubIssue
	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return isActionTriggered, err
	}

	if lenIssues < len(issues) {
		lenIssues = len(issues)
		isActionTriggered = true
		fmt.Println("new commit!")
	} else {
		lenIssues = len(issues)
	}
	return isActionTriggered, nil
}

func (self *WorkflowService) CheckWorkflowsWithGithubActions(action entities.Action) error {
	var isActionTriggered bool

	workflows, err := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		user, err := self.UserRepository.FindUserById(workflow.OwnerId)
		if err != nil {
			return err
		}

		accessToken, err := self.UserServiceService.CallApiAndRefresh(user.Email, user.ConnectionType, "Github")
		if err != nil {
			return err
		}

		if action.Name == "New repository" {
			isActionTriggered, err = self.CheckGithubNewRepositoryAction(workflow, accessToken)
			if err != nil {
				return err
			}
		}
		if action.Name == "New pull request" {
			isActionTriggered, err = self.CheckGithubNewPullRequestAction(workflow, accessToken)
			if err != nil {
				return err
			}
		}
		if action.Name == "New branch" {
			isActionTriggered, err = self.CheckGithubNewBranchAction(workflow, accessToken)
			if err != nil {
				return err
			}
		}
		if action.Name == "New commit" {
			isActionTriggered, err = self.CheckGithubNewCommitAction(workflow, accessToken)
			if err != nil {
				return err
			}
		}
		if action.Name == "New issue assignated" {
			isActionTriggered, err = self.CheckGithubNewIssueAction(workflow, accessToken)
			if err != nil {
				return err
			}
		}
		if isActionTriggered {
			self.executeSpotifyReactions(workflow)
			self.executeDiscordReactions(workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckGithubActions() error {
	service, err := self.ServiceService.FindServiceByName("Github")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.CheckWorkflowsWithGithubActions(action)
		if err != nil {
			return err
		}
	}
	return nil
}
