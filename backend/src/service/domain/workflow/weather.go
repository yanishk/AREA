package workflow_service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"backend/src/entities"
)

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
		return fmt.Errorf(errorMissingField)
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
