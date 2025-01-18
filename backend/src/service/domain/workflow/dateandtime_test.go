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
		_, err := checkEveryHourParamValidity(15)

		require.EqualError(test, err, errorMissingField)
	})
}
