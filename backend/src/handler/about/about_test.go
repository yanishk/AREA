package about_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

type MockAbout struct {
	mock.Mock
}

func (m *MockAbout) GetAboutServer(about entities.About) (entities.About, error) {
	args := m.Called(about)
	return args.Get(0).(entities.About), args.Error(1)
}

func TestCreateUser(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		mockAbout := new(MockAbout)
		handler := &AboutHandler{AboutService: mockAbout}

		router := gin.Default()
		router.GET("/about.json", handler.getAbout)

		var about entities.About

		mockAbout.On("GetAboutServer", about).
			Return(entities.About{}, nil)

		req, _ := http.NewRequest("GET", "/about.json", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail", func(test *testing.T) {
		mockAbout := new(MockAbout)
		handler := &AboutHandler{AboutService: mockAbout}

		router := gin.Default()
		router.GET("/about.json", handler.getAbout)

		var about entities.About

		mockAbout.On("GetAboutServer", about).
			Return(entities.About{}, errors.New("Fail"))

		req, _ := http.NewRequest("GET", "/about.json", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve informations about the server"}`, w.Body.String())
	})
}
