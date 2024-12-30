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

const errorUpdatingWorkflows = "Could not update workflow"

const redditUserRoute = "https://oauth.reddit.com/user/"

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

func (self *WorkflowService) executeSpotifyReactions(workflow entities.Workflow) error {

	reactionFound, errReaction := self.ReactionRepository.FindReactionByID(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf("Error finding reaction")
	}

	accessToken, err := self.getAccessToken("Spotify", workflow)
	if err != nil {
		return fmt.Errorf("Could not update token")
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
	if errReaction != nil {
		return fmt.Errorf("Error finding reaction")
	}

	_, err := self.getAccessToken("Discord", workflow)
	if err != nil {
		return fmt.Errorf("Could not update token")
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

func (self *WorkflowService) executeGetRedditRequest(url, accessToken string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", accessBearer+accessToken)
	req.Header.Add("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	resp, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (self *WorkflowService) getRedditPost(url, accessToken string) (entities.RedditPostResponse, error) {
	var result entities.RedditPostResponse

	resp, err := self.executeGetRedditRequest(url, accessToken)
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

	resp, err := self.executeGetRedditRequest(url, accessToken)
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

	resp, err := self.executeGetRedditRequest(url, accessToken)
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

	resp, err := self.executeGetRedditRequest(url, accessToken)
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
			return fmt.Errorf(errorUpdatingWorkflows)
		}
	}

	if idLatestExist && idLatest != latestPost.ID {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
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
			return fmt.Errorf(errorUpdatingWorkflows)
		}
	}

	if idLatestExist && idLatest != latestComment.ID {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
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
			return fmt.Errorf(errorUpdatingWorkflows)
		}
	}

	if idLatestExist && idLatest != latestVote.ID {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
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

func (self *WorkflowService) executeNewPostInSubredditAction(accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) executeNewPostByMeAction(accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	urlPostMe := redditUserRoute + resultUsername.Name + "/submitted?limit=1"
	resultPost, err := self.getRedditPost(urlPostMe, accessToken)
	if err != nil {
		return err
	}

	return self.isANewPost(workflow, resultPost)
}

func (self *WorkflowService) executeNewCommentByMeAction(accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	urlPostMe := redditUserRoute + resultUsername.Name + "/comments?limit=1"
	resultPost, err := self.getRedditComments(urlPostMe, accessToken)
	if err != nil {
		return err
	}

	return self.isANewComment(workflow, resultPost)
}

func (self *WorkflowService) executeVoteByMeAction(voteType, accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) checkNameOfActionReddit(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		accessToken, err := self.getAccessToken("Reddit", workflow)
		if err != nil {
			continue
		}
		if action.Name == "Any new post in subreddit" {
			self.executeNewPostInSubredditAction(accessToken, workflow)
		}
		if action.Name == "New post by you" {
			self.executeNewPostByMeAction(accessToken, workflow)
		}
		if action.Name == "New comment by you" {
			self.executeNewCommentByMeAction(accessToken, workflow)
		}
		if action.Name == "New downvoted post by you" {
			self.executeVoteByMeAction("/downvoted", accessToken, workflow)
		}
		if action.Name == "New upvoted post by you" {
			self.executeVoteByMeAction("/upvoted", accessToken, workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CallRedditApi() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Reddit")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkNameOfActionReddit(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}
