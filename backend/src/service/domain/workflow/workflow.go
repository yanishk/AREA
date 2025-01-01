package workflow_service

import (
	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

const errorUpdatingWorkflow = "Could not update workflow"

const redditUserRoute = "https://oauth.reddit.com/user/"
const githubBaseUrl = "https://api.github.com/"

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

func (self *WorkflowService) getReactionAndToken(serviceName string, workflow entities.Workflow) (entities.Reaction, string, error) {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return reactionFound, "", fmt.Errorf("Error finding reaction")
	}

	accessToken, err := self.getAccessToken(serviceName, workflow)
	if err != nil {
		return reactionFound, accessToken, fmt.Errorf("Could not update token")
	}
	return reactionFound, accessToken, nil
}

func (self *WorkflowService) CheckSpotifyReactions(workflow entities.Workflow) error {
	reactionFound, accessToken, err := self.getReactionAndToken("Spotify", workflow)
	if err != nil {
		return fmt.Errorf("Could not update token or find reaction")
	}

	var resp *http.Response
	switch reactionFound.Name {
	case "Start playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/play", "PUT", accessBearer, accessToken, nil)
	case "Pause playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", accessBearer, accessToken, nil)
	case "Skip track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", accessBearer, accessToken, nil)
	case "Add track to playback queue":
		uri, uriErr := workflow.ReactionParam["uri"]
		if !uriErr {
			return fmt.Errorf("Incorrect URI")
		}
		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlQueue, "POST", accessBearer, accessToken, nil)
	case "Save a track":
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
	reactionFound, _, err := self.getReactionAndToken("Discord", workflow)
	if err != nil {
		return fmt.Errorf("Could not find reaction")
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

func (self *WorkflowService) ExecuteLinkedinRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
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
		return fmt.Errorf("Failed to get message body")
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
		return fmt.Errorf("Could not encode the body")
	}

	_, errLinkedinRequest := self.ExecuteLinkedinRequest("POST", url, accessToken, bytes.NewBuffer([]byte(jsonBody)))
	if errLinkedinRequest != nil {
		return errLinkedinRequest
	}
	return nil
}

func (self *WorkflowService) publishURLLinkedinPost(urn, accessToken string, workflow entities.Workflow) error {
	message, messageErr := workflow.ReactionParam["message"]
	urlPost, urlErr := workflow.ReactionParam["url"]
	if !messageErr || !urlErr {
		return fmt.Errorf("Failed to get message body or URL")
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
		return fmt.Errorf("Could not encode the body")
	}

	_, errLinkedinRequest := self.ExecuteLinkedinRequest("POST", url, accessToken, bytes.NewBuffer([]byte(jsonBody)))
	if errLinkedinRequest != nil {
		return errLinkedinRequest
	}
	return nil
}

func (self *WorkflowService) CheckLinkedinReactions(workflow entities.Workflow) error {
	reactionFound, accessToken, err := self.getReactionAndToken("Linkedin", workflow)
	if err != nil {
		return fmt.Errorf("Could not update token or find reaction")
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

func (self *WorkflowService) CheckReactions(workflow entities.Workflow) {
	self.CheckSpotifyReactions(workflow)
	self.CheckDiscordReactions(workflow)
	self.CheckLinkedinReactions(workflow)
}

func (self *WorkflowService) CheckTimeAndDateEveryDayAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !hourExists || !minuteExists {
		return fmt.Errorf("Error in the hour and minute sent")
	}

	if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckTimeAndDateEveryHourAction(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	if int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.CheckReactions(workflow)
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
			self.CheckReactions(workflow)
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
		self.CheckReactions(workflow)
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
		self.CheckReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckWorkflowsWithTimeAndDateActions(timeRes entities.TimeResponse, action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		switch action.Name {
		case "Every day at":
			self.CheckTimeAndDateEveryDayAction(timeRes, workflow)
		case "Every hour at":
			self.CheckTimeAndDateEveryHourAction(timeRes, workflow)
		case "Every day of the week at":
			self.CheckTimeAndDateEveryWeekAtAction(timeRes, workflow)
		case "Every month on the":
			self.CheckTimeAndDateEveryMonthOnTheAction(timeRes, workflow)
		case "Every year on":
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
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestPost.ID {
		self.CheckReactions(workflow)
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
		self.CheckReactions(workflow)
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
		self.CheckReactions(workflow)
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

func (self *WorkflowService) CheckRedditNewPostInSubredditAction(accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) CheckRedditNewCommentByMeAction(accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) CheckRedditVoteByMeAction(voteType, accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) CheckRedditPostAction(postType, accessToken string, workflow entities.Workflow) error {
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

func (self *WorkflowService) CheckWorkflowsWithRedditActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		accessToken, err := self.getAccessToken("Reddit", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "Any new post in subreddit":
			self.CheckRedditNewPostInSubredditAction(accessToken, workflow)
		case "New post by you":
			self.CheckRedditPostAction("/submitted", accessToken, workflow)
		case "New comment by you":
			self.CheckRedditNewCommentByMeAction(accessToken, workflow)
		case "New downvoted post by you":
			self.CheckRedditVoteByMeAction("/downvoted", accessToken, workflow)
		case "New upvoted post by you":
			self.CheckRedditVoteByMeAction("/upvoted", accessToken, workflow)
		case "New post saved by you":
			self.CheckRedditPostAction("/saved", accessToken, workflow)
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

func (self *WorkflowService) CheckWeatherCurrentWeatherComparisonAction(checkType string, workflow entities.Workflow) error {
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
		self.CheckReactions(workflow)
	} else if checkType == "belowCurrent" && weatherData.Current.Temperature < temp.(float64) {
		self.CheckReactions(workflow)
	} else if checkType == "aboveForecast" && weatherData.Forecast.ForecastDay[0].Day.MaxTemperature < temp.(float64) {
		self.CheckReactions(workflow)
	} else if checkType == "belowForecast" && weatherData.Forecast.ForecastDay[0].Day.MinTemperature > temp.(float64) {
		self.CheckReactions(workflow)
	} else {
		fmt.Errorf("Check type isn't valid")
	}
	return nil
}

func (self *WorkflowService) CheckWorkflowsWithWeatherActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		switch action.Name {
		case "Current temperature rises above":
			self.CheckWeatherCurrentWeatherComparisonAction("aboveCurrent", workflow)
		case "Current temperature drops below":
			self.CheckWeatherCurrentWeatherComparisonAction("belowCurrent", workflow)
		case "Tomorrow's low drops below":
			self.CheckWeatherCurrentWeatherComparisonAction("belowForecast", workflow)
		case "Tomorrow's high rises above":
			self.CheckWeatherCurrentWeatherComparisonAction("aboveForecast", workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckWeatherActions() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("FreeWeather")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.CheckWorkflowsWithWeatherActions(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

func (self *WorkflowService) ExecuteGithubRequest(workflow entities.Workflow, accessToken string, endpoint string) (*http.Response, error) {
	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return nil, fmt.Errorf("Invalid repository")
	}

	url := githubBaseUrl + "repos/" + repository.(string) + endpoint
	return self.ServiceService.ExecuteApiRequest(url, "GET", accessBearer, accessToken, nil)
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
		self.CheckReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) CheckGithubNewRepositoryAction(workflow entities.Workflow, accessToken string) error {
	var repositories []entities.GithubRepository

	resp, err := self.ServiceService.ExecuteApiRequest(githubBaseUrl+"user/repos", "GET", accessBearer, accessToken, nil)
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

	resp, err := self.ServiceService.ExecuteApiRequest(githubBaseUrl+"issues", "GET", accessBearer, accessToken, nil)
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

	resp, err := self.ExecuteGithubRequest(workflow, accessToken, "/pulls")
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

	resp, err := self.ExecuteGithubRequest(workflow, accessToken, "/branches")
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

	resp, err := self.ExecuteGithubRequest(workflow, accessToken, "/commits")
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
		accessToken, err := self.getAccessToken("Github", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "New repository":
			self.CheckGithubNewRepositoryAction(workflow, accessToken)
		case "New issue assignated":
			self.CheckGithubNewIssueAction(workflow, accessToken)
		case "New pull request":
			self.CheckGithubNewPullRequestAction(workflow, accessToken)
		case "New branch":
			self.CheckGithubNewBranchAction(workflow, accessToken)
		case "New commit":
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
