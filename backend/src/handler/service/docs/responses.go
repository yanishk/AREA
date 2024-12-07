package docs_service

type ServiceActionsReactionsExample struct {
	Name  string `json:"Name" example:"Google"`
	Color string `json:"Color" example:"#E84A4A"`
	Logo  string `json:"Logo" example:"assets/logos/Google.png"`
}

type ServiceActionsReactionsSuccessResponse struct {
	Services []ServiceActionsReactionsExample
}

type ServiceActionReactionInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Internal Server Error"`
}
