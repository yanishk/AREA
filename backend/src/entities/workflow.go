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
	ReactionParam map[string]interface{} `json:"reactionparam"`
}

type NewWorkflow struct {
	NameWorkflow  string                 `json:"name"`
	ActionId      string                 `json:"actionid"`
	ReactionId    string                 `json:"reactionid"`
	ActionParam   map[string]interface{} `json:"actionparam"`
	ReactionParam map[string]interface{} `json:"reactionparam"`
}

type ChangeActivationWorkflow struct {
	WorkflowId string `json:"id"`
}
