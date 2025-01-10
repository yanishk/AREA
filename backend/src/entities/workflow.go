package entities

type Workflow struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	OwnerId       string                 `json:"ownerid"`
	ActionId      string                 `json:"actionid"`
	ReactionId    string                 `json:"reactionid"`
	IsActivated   bool                   `json:"isactivated"`
	CreatedAt     string                 `json:"createdat"`
	ActionParam   map[string]interface{} `json:"actionparam"`
	ReactionParam map[string]interface{} `json:"reactionparam"`
	ActionData    map[string]interface{} `json:"actiondata"`
}

type NewWorkflow struct {
	Name          string                 `json:"name"`
	ActionId      string                 `json:"actionid"`
	ReactionId    string                 `json:"reactionid"`
	ActionParam   map[string]interface{} `json:"actionparam"`
	ReactionParam map[string]interface{} `json:"reactionparam"`
	ActionData    map[string]interface{} `json:"actiondata"`
}

type UpdatedWorkflow struct {
	Name          *string                 `json:"name"`
	ActionId      *string                 `json:"actionid"`
	ReactionId    *string                 `json:"reactionid"`
	IsActivated   *bool                   `json:"isactivated"`
	ActionParam   *map[string]interface{} `json:"actionparam"`
	ReactionParam *map[string]interface{} `json:"reactionparam"`
}
