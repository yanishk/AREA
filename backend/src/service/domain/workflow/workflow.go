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
		newWorkflow.ActionParam, newWorkflow.ReactionParam, newWorkflow.ActionData)
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

func (self *WorkflowService) CheckSpotifyReactions(workflow entities.Workflow) error {

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

func (self *WorkflowService) CheckDiscordReactions(workflow entities.Workflow) error {
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

func (self *WorkflowService) CheckTimeAndDateEveryDayAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !hourExists || !minuteExists {
		return fmt.Errorf("Error in the hour and minute sent")
	}

	if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckTimeAndDateEveryHourAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	if int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckTimeAndDateEveryWeekAtAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
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
			self.CheckSpotifyReactions(workflow)
			self.CheckDiscordReactions(workflow)
		}
	}

	return nil
}

func (self *WorkflowService) CheckTimeAndDateEveryMonthOnTheAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf("Error in the minute, hour or day sent")
	}

	if int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckTimeAndDateEveryYearOnAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	month, monthExists := workflow.ActionParam["month"]
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists || !monthExists {
		return fmt.Errorf("Error in the minute, hour, day or month sent")
	}

	if int(month.(float64)) == timeRes.Month && int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckWorkflowsWithTimeAndDateActions(timeRes entities.TimeResponse, action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if action.Name == "Every day at" {
			self.CheckTimeAndDateEveryDayAction(timeRes, workflow)
		}
		if action.Name == "Every hour at" {
			self.CheckTimeAndDateEveryHourAction(timeRes, workflow)
		}
		if action.Name == "Every day of the week at" {
			self.CheckTimeAndDateEveryWeekAtAction(timeRes, workflow)
		}
		if action.Name == "Every month on the" {
			self.CheckTimeAndDateEveryMonthOnTheAction(timeRes, workflow)
		}
		if action.Name == "Every year on" {
			self.CheckTimeAndDateEveryYearOnAction(timeRes, workflow)
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

	allActionsFound, errFindingActions := self.ActionRepository.FindActionByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.CheckWorkflowsWithTimeAndDateActions(timeRes, action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) retrieveSubredditInfo(subreddit, accessToken string) (entities.SubredditResponse, error) {
	var result entities.SubredditResponse
	url := "https://oauth.reddit.com/" + subreddit + "/new.json?limit=1"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", accessBearer+accessToken)
	req.Header.Add("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	resp, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve post from the specified subreddit")
	}
	return result, nil
}

func (self *WorkflowService) CheckRedditNewPostInSubredditAction(workflow entities.Workflow) error {
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errUser != nil {
		return fmt.Errorf("Error finding user")
	}
	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, "Reddit")
	if errRefreshing != nil {
		return errRefreshing
	}

	subreddit, subredditExists := workflow.ActionParam["subreddit"]
	if !subredditExists {
		return fmt.Errorf("Could not find the subreddit field")
	}

	result, err := self.retrieveSubredditInfo(subreddit.(string), accessToken)
	if err != nil {
		return err
	}

	latestPost := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestPost.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestPost.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf("Could not update workflow")
		}
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckWorkflowsWithRedditActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if action.Name == "Any new post in subreddit" {
			return self.CheckRedditNewPostInSubredditAction(workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckRedditActions() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Reddit")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.CheckWorkflowsWithRedditActions(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) CheckActionDataLen(workflow entities.Workflow, lenActualTurn float64) error {
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
		self.CheckSpotifyReactions(workflow)
		self.CheckDiscordReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckGithubNewRepositoryAction(workflow entities.Workflow, accessToken string) error {
	var repositories []entities.GithubRepository

	resp, err := self.ServiceService.ExecuteApiRequest("https://api.github.com/user/repos", "GET", accessBearer, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		return err
	}

	return self.CheckActionDataLen(workflow, float64(len(repositories)))
}

func (self *WorkflowService) CheckGithubNewIssueAction(workflow entities.Workflow, accessToken string) error {
	var issues []entities.GithubIssue

	resp, err := self.ServiceService.ExecuteApiRequest("https://api.github.com/issues", "GET", accessBearer, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return err
	}

	return self.CheckActionDataLen(workflow, float64(len(issues)))
}

func (self *WorkflowService) CheckGithubNewPullRequestAction(workflow entities.Workflow, accessToken string) error {
	var pullRequests []entities.GithubRepository

	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "pulls"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		return err
	}

	return self.CheckActionDataLen(workflow, float64(len(pullRequests)))
}

func (self *WorkflowService) CheckGithubNewBranchAction(workflow entities.Workflow, accessToken string) error {
	var branches []entities.GithubBranch

	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "branches"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&branches)
	if err != nil {
		return err
	}

	return self.CheckActionDataLen(workflow, float64(len(branches)))
}

func (self *WorkflowService) CheckGithubNewCommitAction(workflow entities.Workflow, accessToken string) error {
	var commits []entities.GithubCommit

	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return fmt.Errorf("Invalid repository")
	}

	uri :=  "https://api.github.com/repos/" + repository.(string) + "/" + "commits"
	resp, err := self.ServiceService.ExecuteApiRequest(uri, "GET", accessBearer, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return err
	}

	return self.CheckActionDataLen(workflow, float64(len(commits)))
}

func (self *WorkflowService) CheckWorkflowsWithGithubActions(action entities.Action) error {
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
			self.CheckGithubNewRepositoryAction(workflow, accessToken)
		}
		if action.Name == "New issue assignated" {
			self.CheckGithubNewIssueAction(workflow, accessToken)
		}
		if action.Name == "New pull request" {
			self.CheckGithubNewPullRequestAction(workflow, accessToken)
		}
		if action.Name == "New branch" {
			self.CheckGithubNewBranchAction(workflow, accessToken)
		}
		if action.Name == "New commit" {
			self.CheckGithubNewCommitAction(workflow, accessToken)
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
