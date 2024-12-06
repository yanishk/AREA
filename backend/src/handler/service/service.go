package service_handler

import (
	"backend/src/entities"
	"backend/src/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService service.ServiceService
}

func NewServiceHandler(ServiceService service.ServiceService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{ServiceService: ServiceService}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
}

func (self *ServiceHandler) CallTimeApi() error {
	timeUrl := "https://tools.aimylogic.com/api/now?tz=Europe/Paris"
	var timeRes entities.TimeResponse

	request, errRequest := http.NewRequest("GET", timeUrl, nil)
	if errRequest != nil {
		return errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed: %s", res.Status)
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&timeRes)
	if errDecoder != nil {
		return errDecoder
	}

	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Time & Date")
	if errFindingService != nil {
		return errFindingService
	}

	allTriggerFound, errFindingTrigger := self.ServiceService.FindTriggerByIdService(serviceFound.Id)
	if errFindingTrigger != nil {
		return errFindingTrigger
	}

	for _, trigger := range allTriggerFound {
		if trigger.Name == "Every day at" {

			AllworkflowFound, errWorkflow := self.ServiceService.FindWorkflowByIdTrigger(trigger.Id)
			if errWorkflow != nil {
				return errWorkflow
			}

			for _, workflow := range AllworkflowFound {
				hour, hourExists := workflow.ActionParam["hour"]
				minute, minuteExists := workflow.ActionParam["minute"]
				if !hourExists || !minuteExists {
					return fmt.Errorf("missing hour or minute in ActionParam")
				}

				if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute {
					fmt.Println("TRIGGER CALLED: matching time found")
				}
			}
		}
	}
	return nil
}
