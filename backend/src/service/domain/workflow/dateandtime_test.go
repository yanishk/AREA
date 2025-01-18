package workflow_service

import (
	"backend/src/entities"
	"testing"

	"github.com/stretchr/testify/require"
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
			Minute: 1,
		}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"minute": 1.0,
			},
			IsActivated: false,
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

	test.Run("Invalid Minute", func(test *testing.T) {
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

	test.Run("Invalid Hour/Minute", func(test *testing.T) {
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

	test.Run("Invalid Hour/Minute", func(test *testing.T) {
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

	test.Run("Invalid Hour/Minute", func(test *testing.T) {
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
