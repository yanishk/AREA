package entities

type Workflow struct {
	Id            string
	Name          string
	OwnerId       string
	ActionId      string
	ReactionId    string
	IsActivated   bool
	CreatedAt     string
	ActionParam   map[string]interface{} `json:"actionparam"`
	ReactionParam map[string]interface{} `json:"Reactionparam"`
}
