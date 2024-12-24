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
	errCreationWorkflow := self.WorkflowRepository.CreateWorkflow(newWorkflow.NameWorkflow,
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

func (self *WorkflowService) ActivateDeactivateWorkflow(workflowActivation entities.ChangeActivationWorkflow) error {
	foundWorkflow, errFetchWorkflow := self.WorkflowRepository.FindWorkflowById(workflowActivation.WorkflowId)
	if errFetchWorkflow != nil {
		return errFetchWorkflow
	}

	errUpdate := self.WorkflowRepository.ActivateDeactivateWorkflow(workflowActivation.WorkflowId, foundWorkflow.IsActivated)
	if errUpdate != nil {
		return errUpdate
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
	defer resp.Body.Close()
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

	dayArray, err := days.([]interface{})
	if !err {
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
