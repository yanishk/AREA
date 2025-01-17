package docs_userservice

import "backend/src/entities"

type DiscordServerExample struct {
	Banner      string
	Features    []string
	Icon        string
	Id          string
	Name        string
	Owner       bool
	Permissions string
}

// Service Callback Responses
type UserServiceServiceCallbackSuccessResponse struct {
	Msg string `json:"success"example:"Token generated"`
}

type UserServiceServiceCallbackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Invalid code authorization-Invalid app type"`
}

type UserServiceServiceCallbackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to get email-Failed to update token"`
}

// User Service Authentication Status Responses
type UserServiceGetUserServiceAuthenticationStatusSuccessResponse struct {
	Msg bool `json:"authenticated"`
}

type UserServiceGetUserServiceAuthenticationStatusBadRequestResponse struct {
	Msg string `json:"error"example:"Could not find requested user-Unknown service"`
}

// Github Get User Repositories
type UserServiceGetGithubUserRepositoriesSuccessResponse struct {
	Repositories []entities.GithubRepository
}

type UserServiceGetGithubUserRepositoriesInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the user's repositories"`
}

// Gitlab Get User Projects
type UserServiceGetGitlabUserProjectsSuccessResponse struct {
	Projects []entities.GitlabProject
}

type UserServiceGetGitlabUserProjectsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the user's projects"`
}

// Discord Get User Servers
type UserServiceGetDiscordUserServersSuccessResponse struct {
	Servers []DiscordServerExample
}

type UserServiceGetDiscordUserServersInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the user's repositories"`
}

// Asana Get User Workspaces
type UserServiceGetAsanaUserWorkspacesSuccessResponse struct {
	Workspaces entities.AsanaWorkspacesInfo
}

type UserServiceGetAsanaUserWorkspacesInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the user's workspaces"`
}

// Asana Get Workspace Assignees
type UserServiceGetAsanaWorkspaceAssigneesSuccessResponse struct {
	Assignees entities.AsanaWorkspacesInfo
}

type UserServiceGetAsanaWorkspaceAssigneesInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's assignees"`
}

// Asana Get Workspace Projects
type UserServiceGetAsanaWorkspaceProjectsSuccessResponse struct {
	Projects entities.AsanaWorkspacesInfo
}

type UserServiceGetAsanaWorkspaceProjectsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's projects"`
}

// Asana Get Workspace Tags
type UserServiceGetAsanaWorkspaceTagsSuccessResponse struct {
	Tags entities.AsanaWorkspacesInfo
}

type UserServiceGetAsanaWorkspaceTagsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's tags"`
}
