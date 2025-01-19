package workflow_service

import (
	"fmt"

	"backend/src/entities"
)

func checkEveryHourParamValidity(minute interface{}) (float64, error) {
	minuteNumber, minuteIsNumber := minute.(float64)
	if !minuteIsNumber {
		return minuteNumber, fmt.Errorf(errorMissingField)
	}
	return minuteNumber, nil
}

func checkEveryDayParamValidity(hour interface{}, minute interface{}) (float64, float64, error) {
	var hourNumber, minuteNumber float64

	hourNumber, hourIsNumber := hour.(float64)
	if !hourIsNumber {
		return hourNumber, minuteNumber, fmt.Errorf(errorMissingField)
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
		return dayNumber, hourNumber, minuteNumber, fmt.Errorf(errorMissingField)
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
		return monthNumber, dayNumber, hourNumber, minuteNumber, fmt.Errorf(errorMissingField)
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
		return fmt.Errorf(errorMissingField)
	}

	minuteNumber, err := checkEveryHourParamValidity(minute)
	if err != nil {
		return err
	}

	if int(minuteNumber) == timeRes.Minute {
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
