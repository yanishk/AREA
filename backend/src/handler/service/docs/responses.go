package docs_service

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

type DiscordChannelIconEmojiExample struct {
	Id   string
	Name string
}

type DiscordChannelPermissionOverwritesExample struct {
	Allow string
	Deny  string
	Id    string
	Type  int
}

type DiscordChannelExample struct {
	Flags                 int
	Guild_id              string
	Icon_emoji            []DiscordChannelIconEmojiExample
	Id                    string
	Name                  string
	Nsfw                  bool
	Parent_id             string
	Permission_overwrites []DiscordChannelPermissionOverwritesExample
	Position              int
	Rate_limite_per_user  int
	Theme_color           string
	Topic                 string
	Type                  int
}

// General Responses
type ServiceInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Internal server error"`
}

// OAuth2 Responses
type ServiceOAuth2SuccessResponse struct {
	Msg string `json:"auth-url"example:"url-example.com"`
}

type ServiceOAuth2BadRequestResponse struct {
	Msg string `json:"error"example:"Invalid callback type-Invalid app type-Unknown service"`
}

// Retrieve Service by Id
type ServiceRetrieveServiceByIdSuccessResponse struct {
	Service entities.Service
}

type ServiceRetrieveServiceByIdBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid service id"`
}

type ServiceRetrieveServiceByIdInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve requested service"`
}

// Retrieve Actions/Reactions Services Responses
type ServiceRetrieveActionsReactionsServicesSuccessResponse struct {
	Services []entities.Service
}

// Retrieve Actions from Service Responses
type ServiceRetrieveActionsFromServiceSuccessResponse struct {
	Actions []entities.Action
}

// Retrieve Reactions from Service Responses
type ServiceRetrieveReactionsFromServiceSuccessResponse struct {
	Reactions []entities.Reaction
}

// Discord Get User Servers
type ServiceGetDiscordUserServersSuccessResponse struct {
	Servers []DiscordServerExample
}

type ServiceGetDiscordUserServersInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve or refresh the token-Could not retrieve the user's servers"`
}

// Discord Get Guild Channels
type ServiceGetDiscordGuildChannelsSuccessResponse struct {
	Channels []DiscordChannelExample
}

type ServiceGetDiscordGuildChannelsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the server's channels"`
}

// Asana Get User Workspaces
type ServiceGetAsanaUserWorkspacesSuccessResponse struct {
	Workspaces entities.AsanaUserWorkspaces
}

type ServiceGetAsanaUserWorkspacesInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve or refresh the token-Could not retrieve the user's workspaces-Could not decode the user's workspaces"`
}

// Asana Get Workspace Assignees
type ServiceGetAsanaWorkspaceAssigneesSuccessResponse struct {
	Assignees entities.AsanaWorkspacesInfo
}

type ServiceGetAsanaWorkspaceAssigneesInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's assignees"`
}

// Asana Get Workspace Projects
type ServiceGetAsanaWorkspaceProjectsSuccessResponse struct {
	Projects entities.AsanaWorkspacesInfo
}

type ServiceGetAsanaWorkspaceProjectsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's projects"`
}

// Asana Get Workspace Tags
type ServiceGetAsanaWorkspaceTagsSuccessResponse struct {
	Tags entities.AsanaWorkspacesInfo
}

type ServiceGetAsanaWorkspaceTagsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve workspace's tags"`
}
