package workflow_service

import (
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
