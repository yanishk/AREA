package storage

import (
	"backend/src/entities"
)

type UserRepository interface {
	CreateUser(email string, password string) error
	FindUserByEmail(email string) (entities.User, error)
	UpdateUser(email string, password string) error
	DeleteUser(email string) error
}

type ServiceRepository interface {
	CreateService(name string) error
	FindServiceByName(name string) (entities.Service, error)
}

type ActionRepository interface {
	CreateAction(name string, description string, serviceId string) error
	FindActionByName(name string) (entities.Action, error)
}

type TriggerRepository interface {
	CreateTrigger(name string, description string, serviceId string) error
	FindTriggerByName(name string) (entities.Trigger, error)
}

type WorkflowRepository interface {
	CreateWorkflow(name string, ownerId string, triggerId string, actionId string) error
	FindWorkflowByName(name string) (entities.Workflow, error)
}

type Repository struct {
	UserRepository UserRepository
	ServiceRepository ServiceRepository
	ActionRepository ActionRepository
	TriggerRepository TriggerRepository
	WorkflowRepository WorkflowRepository
}
