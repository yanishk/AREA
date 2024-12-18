package docs_service

type ServiceExample struct {
	Name  string `json:"Name" example:"Google"`
	Color string `json:"Color" example:"#E84A4A"`
	Logo  string `json:"Logo" example:"assets/logos/Google.webp"`
}

type ActionExample struct {
	Id          string `json:"Id" example:"2685fc54-d1e2-43d8-9525-56ab7a91e351"`
	Name        string `json:"Name" example:"Every day at"`
	Description string `json:"Description" example:"This Trigger fires every single day at a specific time set by you."`
	ServiceId   string `json:"ServiceId" example:"203a1f6a-b2d7-4011-a73b-e9585907835b"`
	NbParam     int    `json"NbParam" example: 2`
}

type ReactionExample struct {
	Id          string `json:"Id" example:"5dccb0f3-64db-4456-8286-4678857ae394"`
	Name        string `json:"Name" example:"Start playback"`
	Description string `json:"Description" example:"Start a new context or resume current playback on the user’s active device. Only available for Spotify premium accounts."`
	ServiceId   string `json:"ServiceId" example:"4f6f8655-8e0f-4f96-ba7f-d3cab1577b74"`
	NbParam     int    `json"NbParam" example: 0`
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
	Msg string `json:"error"example:"Invalid callback type"`
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
