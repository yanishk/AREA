package docs_service

import "backend/src/entities"

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

// Retrieve All Services
type ServiceRetrieveAllServicesSuccessResponse struct {
	Services []entities.Service
}

// Retrieve Actions/Reactions Services Responses
type ServiceRetrieveActionsReactionsServicesSuccessResponse struct {
	Services []entities.Service
}

// Retrieve Service by Id
type ServiceRetrieveServiceByIdSuccessResponse struct {
	Service entities.Service
}

type ServiceRetrieveServiceByIdInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve requested service"`
}

// Retrieve Actions from Service Responses
type ServiceRetrieveActionsFromServiceSuccessResponse struct {
	Actions []entities.Action
}

type ServiceRetrieveActionsFromServiceBadRequestResponse struct {
	Msg string `json:"error"example:"Unknown service"`
}

// Retrieve Reactions from Service Responses
type ServiceRetrieveReactionsFromServiceSuccessResponse struct {
	Reactions []entities.Reaction
}

type ServiceRetrieveReactionsFromServiceBadRequestResponse struct {
	Msg string `json:"error"example:"Unknown service"`
}

// Discord Get Guild Channels
type ServiceGetDiscordGuildChannelsSuccessResponse struct {
	Channels []DiscordChannelExample
}

type ServiceGetDiscordGuildChannelsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve the server's channels"`
}
