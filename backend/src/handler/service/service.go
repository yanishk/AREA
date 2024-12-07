package service_handler

import (
	"backend/src/entities"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
	UserService        service.UserService
}

func NewServiceHandler(ServiceService service.ServiceService, UserServiceService service.UserServiceService, UserService service.UserService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{
		ServiceService:     ServiceService,
		UserServiceService: UserServiceService,
		UserService:        UserService,
	}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	service := router.Group("/:service")
	{
		service.GET("/actions", self.RetrieveActionsFromService)
		service.GET("/reactions", self.RetrieveReactionsFromService)
	}
	router.GET("/actions-services", self.RetrieveActionsServices)
	router.GET("/reactions-services", self.RetrieveReactionsServices)
}

func (self *ServiceHandler) ReactionStartPlayback(reactionName string, emailUser string) error {
	spotifyResumePlayURL := "https://api.spotify.com/v1/me/player/play"
	if reactionName != "Start playback" {
		return nil
	}

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(emailUser, "Spotify")
	if errRefreshing != nil {
		return errRefreshing
	}

	req, err := http.NewRequest("PUT", spotifyResumePlayURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Spotify API call failed")
	}
	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryDay(timeRes entities.TimeResponse, AllworkflowFound []entities.Workflow) error {
	for _, workflow := range AllworkflowFound {
		hour, hourExists := workflow.ActionParam["hour"]
		minute, minuteExists := workflow.ActionParam["minute"]
		if !hourExists || !minuteExists {
			return fmt.Errorf("missing hour or minute in ActionParam")
		}

		if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
			reactionFound, errReaction := self.ServiceService.FindReactionByID(workflow.ReactionId)
			if errReaction != nil {
				return errReaction
			}
			foundUser, errUser := self.UserService.FindUserById(workflow.OwnerId)
			if errUser != nil {
				return errUser
			}

			return self.ReactionStartPlayback(reactionFound.Name, foundUser.Email)
		}
	}
	return nil
}

func (self *ServiceHandler) checkNameOfAction(timeRes entities.TimeResponse, action entities.Action) error {
	if action.Name == "Every day at" {

		AllworkflowFound, errWorkflow := self.ServiceService.FindWorkflowByActionId(action.Id)
		if errWorkflow != nil {
			return errWorkflow
		}
		return self.executeActiveActionsEveryDay(timeRes, AllworkflowFound)
	}
	return nil
}

func (self *ServiceHandler) requestToTimeAPI() (entities.TimeResponse, error) {
	timeUrl := "https://tools.aimylogic.com/api/now?tz=Europe/Paris"
	var timeRes entities.TimeResponse

	request, errRequest := http.NewRequest("GET", timeUrl, nil)
	if errRequest != nil {
		return timeRes, errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return timeRes, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return timeRes, fmt.Errorf("HTTP request failed: %s", res.Status)
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&timeRes)
	if errDecoder != nil {
		return timeRes, errDecoder
	}
	return timeRes, nil
}

func (self *ServiceHandler) CallTimeApi() error {
	var timeRes entities.TimeResponse

	timeRes, errRequest := self.requestToTimeAPI()
	if errRequest != nil {
		return errRequest
	}

	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Time & Date")
	if errFindingService != nil {
		return errFindingService
	}

	allactionFound, errFindingaction := self.ServiceService.FindActionByServiceId(serviceFound.Id)
	if errFindingaction != nil {
		return errFindingaction
	}

	for _, action := range allactionFound {
		errCheckName := self.checkNameOfAction(timeRes, action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

// @Summary		Actions Services
// @Description	Retrieve services with actions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/actions-services [get]
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	actionsServices, err := self.ServiceService.RetrieveActionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, actionsServices)
	}
}

// @Summary		Reactions Services
// @Description	Retrieve services with reactions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/reactions-services [get]
func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, reactionsServices)
	}
}

// @Summary		Actions From Service
// @Description	Retrieve actions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveActionsFromServiceSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveActionsReactionsFromServiceBadRequestResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/actions [get]
func (self *ServiceHandler) RetrieveActionsFromService(context *gin.Context) {
	serviceName := context.Param("service")
	foundService, errFoundService := self.ServiceService.FindServiceByName(serviceName)
	if errFoundService != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested service",
		})
		return
	}

	actions, errFindActions := self.ServiceService.FindActionByServiceId(foundService.Id)
	if errFindActions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{
			"actions": actions,
		})
	}
}

// @Summary		Reactions From Service
// @Description	Retrieve reactions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveReactionsFromServiceSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveActionsReactionsFromServiceBadRequestResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/reactions [get]
func (self *ServiceHandler) RetrieveReactionsFromService(context *gin.Context) {
	serviceName := context.Param("service")
	foundService, errFoundService := self.ServiceService.FindServiceByName(serviceName)
	if errFoundService != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested service",
		})
		return
	}

	reactions, errFindReactions := self.ServiceService.FindReactionByServiceId(foundService.Id)
	if errFindReactions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{
			"reactions": reactions,
		})
	}
}
