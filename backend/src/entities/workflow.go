package entities

type Workflow struct {
	Id          string
	Name        string
	OwnerId     string
	ActionId    string
	TriggerId   string
	IsActivated bool
	CreatedAt   string
}
