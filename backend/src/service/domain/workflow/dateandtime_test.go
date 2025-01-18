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
