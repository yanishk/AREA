package docs_service

type ServiceExample struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
	HasActions   bool   `json:"hasactions"`
	HasReactions bool   `json:"hasreactions"`
}

type ActionExample struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceId   string `json:"serviceid"`
	NbParam     int    `json:"nbparam"`
}

type ReactionExample struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceId   string `json:"serviceid"`
	NbParam     int    `json:"nbparam"`
}

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
	Service ServiceExample
}

type ServiceRetrieveServiceByIdBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid service id"`
}

type ServiceRetrieveServiceByIdInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve requested service"`
}

// Retrieve Actions/Reactions Services Responses
type ServiceRetrieveActionsReactionsServicesSuccessResponse struct {
	Services []ServiceExample
}

// Retrieve Actions from Service Responses
type ServiceRetrieveActionsFromServiceSuccessResponse struct {
	Actions []ActionExample
}

// Retrieve Reactions from Service Responses
type ServiceRetrieveReactionsFromServiceSuccessResponse struct {
	Reactions []ReactionExample
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
