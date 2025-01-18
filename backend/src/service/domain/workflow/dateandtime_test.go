package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestCheckEveryHourParamValidity(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		res, err := checkEveryHourParamValidity(1.0)

		require.NoError(test, err)
		require.Equal(test, 1.0, res)
	})

	test.Run("Failure", func(test *testing.T) {
		_, err := checkEveryHourParamValidity(1)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckEveryDayParamValidity(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		hour, minute, err := checkEveryDayParamValidity(1.0, 1.0)

		require.NoError(test, err)
		require.Equal(test, 1.0, hour)
		require.Equal(test, 1.0, minute)
	})

	test.Run("Hour Failure", func(test *testing.T) {
		_, _, err := checkEveryDayParamValidity(1, 1.0)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Minute Failure", func(test *testing.T) {
		_, _, err := checkEveryDayParamValidity(1.0, 1)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckEveryMonthOnTheParamValidity(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		day, hour, minute, err := checkEveryMonthOnTheParamValidity(1.0, 1.0, 1.0)

		require.NoError(test, err)
		require.Equal(test, 1.0, day)
		require.Equal(test, 1.0, hour)
		require.Equal(test, 1.0, minute)
	})

	test.Run("Fail Day", func(test *testing.T) {
		_, _, _, err := checkEveryMonthOnTheParamValidity(1, 1.0, 1.0)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Fail Hour/Minute", func(test *testing.T) {
		_, _, _, err := checkEveryMonthOnTheParamValidity(1.0, 1, 1)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckEveryYearParamValidity(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		month, day, hour, minute, err := checkEveryYearParamValidity(1.0, 1.0, 1.0, 1.0)

		require.NoError(test, err)
		require.Equal(test, 1.0, month)
		require.Equal(test, 1.0, day)
		require.Equal(test, 1.0, hour)
		require.Equal(test, 1.0, minute)
	})

	test.Run("Month Failure", func(test *testing.T) {
		_, _, _, _, err := checkEveryYearParamValidity(1, 1.0, 1.0, 1.0)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Failure Day/Hour/Minute", func(test *testing.T) {
		_, _, _, _, err := checkEveryYearParamValidity(1.0, 1, 1, 1)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckTimeAndDateEveryHourAction(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		timeDate := &WorkflowService{}

		time := entities.TimeResponse{
			Minute: 2,
		}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"minute": 1.0,
			},
		}

		err := timeDate.checkTimeAndDateEveryHourAction(time, workflow)

		require.NoError(test, err)
	})

	test.Run("Missing Field", func(test *testing.T) {
		timeDate := &WorkflowService{}

		time := entities.TimeResponse{
			Minute: 1.0,
		}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{},
			IsActivated: false,
		}

		err := timeDate.checkTimeAndDateEveryHourAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Invalid Param", func(test *testing.T) {
		timeDate := &WorkflowService{}

		time := entities.TimeResponse{
			Minute: 1,
		}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"minute": "value",
			},
			IsActivated: false,
		}

		err := timeDate.checkTimeAndDateEveryHourAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckTimeAndDateEveryDayAction(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"hour":   1.0,
				"minute": 1.0,
			},
		}

		time := entities.TimeResponse{
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryDayAction(time, workflow)

		require.NoError(test, err)
	})

	test.Run("Missing Field", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{},
		}

		time := entities.TimeResponse{
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryDayAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Invalid Param", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"hour":   "1.0",
				"minute": "1.0",
			},
		}

		time := entities.TimeResponse{
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryDayAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckTimeAndDateEveryWeekAtAction(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"day":    []interface{}{float64(1)},
				"hour":   1.0,
				"minute": 1.0,
			},
		}

		time := entities.TimeResponse{
			WeekDay: 2,
			Hour:    2,
			Minute:  2,
		}

		err := timeDate.checkTimeAndDateEveryWeekAtAction(time, workflow)

		require.NoError(test, err)
	})

	test.Run("Fail Parsing Day", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"day":    1,
				"hour":   1.0,
				"minute": 1.0,
			},
		}

		time := entities.TimeResponse{
			WeekDay: 2,
			Hour:    2,
			Minute:  2,
		}

		err := timeDate.checkTimeAndDateEveryWeekAtAction(time, workflow)

		require.EqualError(test, err, "Error parsing the day array")
	})

	test.Run("Missing Field", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{},
		}

		time := entities.TimeResponse{
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryWeekAtAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Invalid Param", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"day":    []interface{}{float64(1)},
				"hour":   "1.0",
				"minute": "1.0",
			},
		}

		time := entities.TimeResponse{
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryWeekAtAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckTimeAndDateEveryMonthOnTheAction(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"day":    1.0,
				"hour":   1.0,
				"minute": 1.0,
			},
		}

		time := entities.TimeResponse{
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryMonthOnTheAction(time, workflow)

		require.NoError(test, err)
	})

	test.Run("Missing Field", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{},
		}

		time := entities.TimeResponse{
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryMonthOnTheAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Invalid Param", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"day":    "1.0",
				"hour":   "1.0",
				"minute": "1.0",
			},
		}

		time := entities.TimeResponse{
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryMonthOnTheAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckTimeAndDateEveryYearOnAction(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"month":  1.0,
				"day":    1.0,
				"hour":   1.0,
				"minute": 1.0,
			},
		}

		time := entities.TimeResponse{
			Month:  2,
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryYearOnAction(time, workflow)

		require.NoError(test, err)
	})

	test.Run("Missing Field", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{},
		}

		time := entities.TimeResponse{
			Month:  2,
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryYearOnAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Invalid Param", func(test *testing.T) {
		timeDate := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"month":  "1.0",
				"day":    "1.0",
				"hour":   "1.0",
				"minute": "1.0",
			},
		}

		time := entities.TimeResponse{
			Month:  2,
			Day:    2,
			Hour:   2,
			Minute: 2,
		}

		err := timeDate.checkTimeAndDateEveryYearOnAction(time, workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckWorkflowsWithTimeAndDateActions(test *testing.T) {
	test.Run("Success Deactivated", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		timeDate := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		workflows := []entities.Workflow{
			{Id: "1", IsActivated: false},
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return(workflows, nil)

		err := timeDate.checkWorkflowsWithTimeAndDateActions(entities.TimeResponse{}, action)

		require.NoError(test, err)
	})

	test.Run("Success Activated", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		timeDate := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		workflows := []entities.Workflow{
			{Id: "1", IsActivated: true},
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return(workflows, nil)

		err := timeDate.checkWorkflowsWithTimeAndDateActions(entities.TimeResponse{}, action)

		require.NoError(test, err)
	})

	test.Run("Fail Find Workflows", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		timeDate := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		workflows := []entities.Workflow{
			{Id: "1", IsActivated: false},
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return(workflows, errors.New("Fail find workflows"))

		err := timeDate.checkWorkflowsWithTimeAndDateActions(entities.TimeResponse{}, action)

		require.EqualError(test, err, "Fail find workflows")
	})
}

func TestCheckTimeAndDateActions(test *testing.T) {
	test.Run("Fail Request Time API", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		timeDate := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("RequestToTimeApi").
			Return(entities.TimeResponse{}, errors.New("Fail time request"))

		err := timeDate.CheckTimeAndDateActions()

		require.EqualError(test, err, "Fail time request")
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		timeDate := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("RequestToTimeApi").
			Return(entities.TimeResponse{}, nil)

		mockServiceServiceRepo.On("FindServiceByName", "Time & Date").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := timeDate.CheckTimeAndDateActions()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		timeDate := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("RequestToTimeApi").
			Return(entities.TimeResponse{}, nil)

		mockServiceServiceRepo.On("FindServiceByName", "Time & Date").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		err := timeDate.CheckTimeAndDateActions()

		require.EqualError(test, err, "Fail find actions")
	})

	test.Run("Success", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		timeDate := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("RequestToTimeApi").
			Return(entities.TimeResponse{}, nil)

		mockServiceServiceRepo.On("FindServiceByName", "Time & Date").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, nil)

		err := timeDate.CheckTimeAndDateActions()

		require.NoError(test, err)
	})
}
