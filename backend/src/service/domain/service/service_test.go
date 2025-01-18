package service_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) CreateService(name, color, logo string) error {
	args := m.Called(name, color, logo)
	return args.Error(0)
}

func (m *MockServiceRepository) FindServiceById(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindAllServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindActionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindReactionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

type MockActionRepository struct {
	mock.Mock
}

func (m *MockActionRepository) CreateAction(name, description, serviceId string, nbParam int) error {
	args := m.Called(name, description, serviceId, nbParam)
	return args.Error(0)
}

func (m *MockActionRepository) FindActionById(id string) (entities.Action, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionByName(name string) (entities.Action, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionsByServiceId(serviceId string) ([]entities.Action, error) {
	args := m.Called(serviceId)
	return args.Get(0).([]entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionByNameAndServiceId(name, serviceId string) (entities.Action, error) {
	args := m.Called(name, serviceId)
	return args.Get(0).(entities.Action), args.Error(1)
}

type MockReactionRepository struct {
	mock.Mock
}

func (m *MockReactionRepository) CreateReaction(name, description, serviceId string, nbParam int) error {
	args := m.Called(name, description, serviceId, nbParam)
	return args.Error(0)
}

func (m *MockReactionRepository) FindReactionById(id string) (entities.Reaction, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Reaction), args.Error(1)
}

func (m *MockReactionRepository) FindReactionByName(name string) (entities.Reaction, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Reaction), args.Error(1)
}

func (m *MockReactionRepository) FindReactionsByServiceId(serviceId string) ([]entities.Reaction, error) {
	args := m.Called(serviceId)
	return args.Get(0).([]entities.Reaction), args.Error(1)
}

func TestFindServiceByName(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindServiceByName", "service").
		Return(entities.Service{}, nil)

	_, err := serviceservice.FindServiceByName("service")

	require.NoError(test, err)
}

func TestFindServiceById(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindServiceById", "id").
		Return(entities.Service{}, nil)

	_, err := serviceservice.FindServiceById("id")

	require.NoError(test, err)
}

func TestFindServiceByActionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var action entities.Action

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		action.ServiceId = "id"

		mockActionRepo.On("FindActionById", "id").
			Return(action, nil)

		mockServiceRepo.On("FindServiceById", "id").
			Return(entities.Service{}, nil)

		_, err := serviceservice.FindServiceByActionId("id")

		require.NoError(test, err)
	})

	test.Run("Find action failure", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		mockActionRepo.On("FindActionById", "id").
			Return(entities.Action{}, errors.New("Fail find action"))

		_, err := serviceservice.FindServiceByActionId("id")

		require.EqualError(test, err, "Fail find action")
	})
}

func TestFindServiceByReactionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var reaction entities.Reaction

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		reaction.ServiceId = "id"

		mockReactionRepo.On("FindReactionById", "id").
			Return(reaction, nil)

		mockServiceRepo.On("FindServiceById", "id").
			Return(entities.Service{}, nil)

		_, err := serviceservice.FindServiceByReactionId("id")

		require.NoError(test, err)
	})

	test.Run("Find action failure", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", "id").
			Return(entities.Reaction{}, errors.New("Fail find reaction"))

		_, err := serviceservice.FindServiceByReactionId("id")

		require.EqualError(test, err, "Fail find reaction")
	})
}

func TestFindAllServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindAllServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.FindAllServices()

	require.NoError(test, err)
}

func TestRetrieveActionsFromService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, nil)

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, errors.New("Fail find service"))

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.EqualError(test, err, "Fail find actions")
	})
}

func TestRetrieveReactionsFromService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockReactionRepo.On("FindReactionsByServiceId", service.Id).
			Return([]entities.Reaction{}, nil)

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, errors.New("Fail find service"))

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockReactionRepo.On("FindReactionsByServiceId", service.Id).
			Return([]entities.Reaction{}, errors.New("Fail find reactions"))

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.EqualError(test, err, "Fail find reactions")
	})
}

func TestRetrieveActionsServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindActionsServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.RetrieveActionsServices()

	require.NoError(test, err)
}

func TestRetrieveReactionsServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindReactionsServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.RetrieveReactionsServices()

	require.NoError(test, err)
}

func TestGetCallbackAndClientId(test *testing.T) {
	test.Run("Login", func(test *testing.T) {

		serviceName := "SPOTIFY"

		callback, clientId := getCallbackAndClientId("login", "SPOTIFY", true)

		resCallback := serviceName + "_LOGIN_CALLBACK"
		resClientId := serviceName + "_LOGIN_CLIENT_ID"

		assert.Equal(test, os.Getenv(resCallback), callback)
		assert.Equal(test, os.Getenv(resClientId), clientId)
	})

	test.Run("Service", func(test *testing.T) {

		serviceName := "SPOTIFY"

		callback, clientId := getCallbackAndClientId("service", "SPOTIFY", true)

		resCallback := serviceName + "_SERVICE_CALLBACK"
		resClientId := serviceName + "_SERVICE_CLIENT_ID"

		assert.Equal(test, os.Getenv(resCallback), callback)
		assert.Equal(test, os.Getenv(resClientId), clientId)
	})
}

func TestOauth2Google(test *testing.T) {
	test.Run("Login web", func(test *testing.T) {
		callbackLink := os.Getenv("GOOGLE_LOGIN_CALLBACK")
		appTypeLink := os.Getenv("GOOGLE_WEB_CLIENT_ID")

		res := oauth2Google("login", "web")

		expectedRes := "https://accounts.google.com/o/oauth2/auth?" +
			clientIdParam + appTypeLink +
			codeResponseType +
			"&access_type=offline" +
			redirectUriParam + callbackLink +
			"&scope=email"

		assert.Equal(test, res, expectedRes)
	})

	test.Run("Service mobile", func(test *testing.T) {
		callbackLink := os.Getenv("GOOGLE_WEB_CLIENT_ID")
		appTypeLink := os.Getenv("GOOGLE_MOBILE_CLIENT_ID")

		res := oauth2Google("service", "mobile")

		expectedRes := "https://accounts.google.com/o/oauth2/auth?" +
			clientIdParam + appTypeLink +
			codeResponseType +
			"&access_type=offline" +
			redirectUriParam + callbackLink +
			"&scope=email"

		assert.Equal(test, res, expectedRes)
	})
}

func TestOauth2Spotify(test *testing.T) {
	expectedRes := "https://accounts.spotify.com/authorize?" +
		clientIdParam + os.Getenv("SPOTIFY_CLIENT_ID") +
		codeResponseType +
		redirectUriParam + os.Getenv("SPOTIFY_LOGIN_CALLBACK") +
		"&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify playlist-modify-public"

	res := oauth2Spotify("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Discord(test *testing.T) {
	test.Run("Login", func(test *testing.T) {
		scope := "identify+email"

		expectedRes := "https://discord.com/oauth2/authorize?" +
			clientIdParam + os.Getenv("DISCORD_CLIENT_ID") +
			"&permissions=141376" +
			codeResponseType +
			redirectUriParam + os.Getenv("DISCORD_LOGIN_CALLBACK") +
			"&scope=" + scope

		res := oauth2Discord("login")

		assert.Equal(test, res, expectedRes)
	})

	test.Run("Service", func(test *testing.T) {
		scope := "identify+email+guilds+bot"

		expectedRes := "https://discord.com/oauth2/authorize?" +
			clientIdParam + os.Getenv("DISCORD_CLIENT_ID") +
			"&permissions=141376" +
			codeResponseType +
			redirectUriParam + os.Getenv("DISCORD_LOGIN_CALLBACK") +
			"&scope=" + scope

		res := oauth2Discord("service")

		assert.Equal(test, res, expectedRes)
	})
}

func TestOauth2Github(test *testing.T) {
	clientID := os.Getenv("GITHUB_LOGIN_CLIENT_ID")
	callbackLink := os.Getenv("GITHUB_LOGIN_CALLBACK")

	expectedRes := "https://github.com/login/oauth/authorize?" +
		clientIdParam + clientID +
		redirectUriParam + callbackLink +
		"&scope=repo admin:org user" +
		stateParam + os.Getenv("GITHUB_STATE")

	res := oauth2Github("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Reddit(test *testing.T) {
	clientID := os.Getenv("REDDIT_LOGIN_CLIENT_ID")
	callbackLink := os.Getenv("REDDIT_LOGIN_CALLBACK")

	expectedRes := "https://www.reddit.com/api/v1/authorize?" +
		clientIdParam + clientID +
		codeResponseType +
		stateParam + os.Getenv("REDDIT_STATE") +
		redirectUriParam + callbackLink +
		"&duration=permanent" +
		"&scope=identity read submit vote mysubreddits history account edit privatemessages"

	res := oauth2Reddit("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Asana(test *testing.T) {
	callbackLink := os.Getenv("ASANA_LOGIN_CALLBACK")

	expectedRes := "https://app.asana.com/-/oauth_authorize?" +
		clientIdParam + os.Getenv("ASANA_CLIENT_ID") +
		redirectUriParam + callbackLink +
		codeResponseType +
		stateParam + os.Getenv("ASANA_STATE") +
		"&scope=default email profile"

	res := oauth2Asana("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Linkedin(test *testing.T) {
	callbackLink := os.Getenv("LINKEDIN_LOGIN_CALLBACK")

	expectedRes := "https://www.linkedin.com/oauth/v2/authorization?" +
		clientIdParam + url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_ID")) +
		redirectUriParam + callbackLink +
		codeResponseType +
		"&scope=openid email profile w_member_social" +
		stateParam + os.Getenv("LINKEDIN_STATE")

	res := oauth2Linkedin("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Dropbox(test *testing.T) {
	callbackLink := os.Getenv("DROPBOX_LOGIN_CALLBACK")

	expectedRes := "https://www.dropbox.com/oauth2/authorize?" +
		clientIdParam + url.QueryEscape(os.Getenv("DROPBOX_CLIENT_ID")) +
		redirectUriParam + callbackLink +
		codeResponseType + "&token_access_type=offline"

	res := oauth2Dropbox("login")

	assert.Equal(test, res, expectedRes)
}

func TestOauth2Gitlab(test *testing.T) {
	callbackLink := os.Getenv("GITLAB_LOGIN_CALLBACK")

	expectedRes := "https://gitlab.com/oauth/authorize?" +
		clientIdParam + os.Getenv("GITLAB_CLIENT_ID") +
		codeResponseType +
		redirectUriParam + callbackLink +
		stateParam + os.Getenv("GITLAB_STATE") +
		"&scope=api read_api read_user read_repository write_repository"

	res := oauth2Gitlab("login")

	assert.Equal(test, res, expectedRes)
}

func TestOAuth2Service(test *testing.T) {
	test.Run("Google", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Google", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Spotify", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Spotify", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Discord", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Discord", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Github", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Github", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Reddit", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Reddit", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Asana", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Asana", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Linkedin", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Linkedin", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Dropbox", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Dropbox", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Gitlab", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Gitlab", "login", "web")

		require.NoError(test, err)
	})

	test.Run("Unknown", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.OAuth2Service("Unknown", "login", "web")

		require.EqualError(test, err, unknownServiceMessage)
	})
}

func TestGenericAccessTokenRequest(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		_, err := genericAccessTokenRequest("url", "code", "callback")

		assert.NoError(test, err, "Error should be nil")
	})

	test.Run("Failure", func(test *testing.T) {
		_, err := genericAccessTokenRequest(":", "code", "callback")

		require.Error(test, err)
	})
}

func TestGetGoogleOAuth2AccessTokenWebRequestBody(test *testing.T) {
	expectedRes := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",
		"code",
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"callback",
		grantTypeAuthorization,
	)
	res := getGoogleOAuth2AccessTokenWebRequestBody("code", "callback")

	assert.Equal(test, res, expectedRes)
}

func TestGetGoogleOAuth2AccessTokenMobileRequestBody(test *testing.T) {
	expectedRes := fmt.Sprintf(
		"code=%s&client_id=%s&redirect_uri=%s&grant_type=%s",
		"code",
		os.Getenv("GOOGLE_MOBILE_CLIENT_ID"),
		"callback",
		grantTypeAuthorization,
	)
	res := getGoogleOAuth2AccessTokenMobileRequestBody("code", "callback")

	assert.Equal(test, res, expectedRes)
}

func TestGetGoogleOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service Web", func(test *testing.T) {
		_, err := getGoogleOAuth2AccessTokenRequest("code", "service", "web")

		require.NoError(test, err)
	})

	test.Run("Success Login Mobile", func(test *testing.T) {
		_, err := getGoogleOAuth2AccessTokenRequest("code", "login", "mobile")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getGoogleOAuth2AccessTokenRequest("code", "invalid", "mobile")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})

	test.Run("Invalid App Type", func(test *testing.T) {
		_, err := getGoogleOAuth2AccessTokenRequest("code", "service", "invalid")

		require.EqualError(test, err, "Invalid app type")
	})
}

func TestGetGithubOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getGithubOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getGithubOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getGithubOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetGitlabOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getGitlabOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getGitlabOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getGitlabOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetAsanaOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getAsanaOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getAsanaOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getAsanaOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetLinkedinOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getLinkedinOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getLinkedinOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getLinkedinOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetDropboxOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getDropboxOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getDropboxOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getDropboxOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetSpotifyOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getSpotifyOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getSpotifyOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getSpotifyOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetDiscordOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getDiscordOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getDiscordOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getDiscordOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGetRedditOAuth2AccessTokenRequest(test *testing.T) {
	test.Run("Success Service", func(test *testing.T) {
		_, err := getRedditOAuth2AccessTokenRequest("code", "service")

		require.NoError(test, err)
	})

	test.Run("Success Login", func(test *testing.T) {
		_, err := getRedditOAuth2AccessTokenRequest("code", "login")

		require.NoError(test, err)
	})

	test.Run("Invalid Callback Type", func(test *testing.T) {
		_, err := getRedditOAuth2AccessTokenRequest("code", "invalid")

		require.EqualError(test, err, invalidCallbackTypeMessage)
	})
}

func TestGenericRefreshTokenRequest(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.genericRefreshTokenRequest("tokenurl", "refreshToken")

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.genericRefreshTokenRequest(":", "refreshToken")

		require.Error(test, err)
	})
}

func TestGenericJsonBodyRefreshToken(test *testing.T) {
	serviceservice := &ServiceService{}

	expectedRes := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		"clientId",
		"secretId",
		grantTypeRefreshToken,
		"refreshToken",
	)

	res := serviceservice.genericJsonBodyRefreshToken("clientId", "secretId", "refreshToken")

	assert.Equal(test, res, expectedRes)
}

func TestGetGoogleRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetGoogleRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetAsanaRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetAsanaRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetDropboxRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetDropboxRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetSpotifyRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetSpotifyRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetDiscordRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetDiscordRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetRedditRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetRedditRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestGetGitlabRefreshTokenRequest(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetGitlabRefreshTokenRequest("refreshToken")

	require.NoError(test, err)
}

func TestExecuteRequest(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		serviceservice := &ServiceService{}

		req, _ := http.NewRequest("GET", "https://tools.aimylogic.com/api/now?tz=Europe/Paris", nil)

		_, err := serviceservice.ExecuteRequest(req)

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		serviceservice := &ServiceService{}

		req, _ := http.NewRequest("GET", "https://tools.aimylogic.", nil)

		_, err := serviceservice.ExecuteRequest(req)

		require.Error(test, err)
	})
}

func TestExecuteApiRequest(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.ExecuteApiRequest("https://tools.aimylogic.com/api/now?tz=Europe/Paris", "GET", "token", "accessToken", nil)

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.ExecuteApiRequest(":", "GET", "token", "accessToken", nil)

		require.Error(test, err)
	})
}

func TestGetResultTokenFromCode(test *testing.T) {
	serviceservice := &ServiceService{}

	_, err := serviceservice.GetResultTokenFromCode("code", "Unknown", "service", "web")

	require.EqualError(test, err, unknownServiceMessage)
}

func TestGetOAuth2UserEmailRequest(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		_, err := getOAuth2UserEmailRequest("https://tools.aimylogic.com/api/now?tz=Europe/Paris", "accessToken")

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		_, err := getOAuth2UserEmailRequest(":", "accessToken")

		require.Error(test, err)
	})
}

func TestDecodeUserInfoWithMultipleResults(test *testing.T) {
	userInfos := []entities.UserInfo{
		{Email: "test@test.com"},
	}

	responseBody, _ := json.Marshal(userInfos)

	res := &http.Response{
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
		StatusCode: http.StatusOK,
	}

	_, err := decodeUserInfoWithMultipleResults(res)

	require.NoError(test, err)
}

func TestGetUserInfoFromService(test *testing.T) {
	test.Run("Unknown", func(test *testing.T) {
		serviceservice := &ServiceService{}

		_, err := serviceservice.GetUserInfoFromService("accessToken", "invalid")

		require.EqualError(test, err, unknownServiceMessage)
	})
}
