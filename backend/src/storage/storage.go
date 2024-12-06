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
	CreateService(name string, color string, logo string) error
	FindServiceByName(name string) (entities.Service, error)
}

type UserServiceRepository interface {
	UpdateUserServiceByServiceIdAndUserId(uuidService string, uuidUser string, accessToken string, refreshToken string, expiryToken string) error
	CreateUserService(idUser string, token string, refreshToken string, expiryDate string, idService string) error
	FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error)
}

type ActionRepository interface {
	CreateAction(name string, description string, serviceId string, nbParam int) error
	FindActionByName(name string) (entities.Action, error)
}

type TriggerRepository interface {
	CreateTrigger(name string, description string, serviceId string, nbParam int) error
	FindTriggerByName(name string) (entities.Trigger, error)
	FindTriggerByIdService(idService string) ([]entities.Trigger, error)
}

type WorkflowRepository interface {
	CreateWorkflow(name string, ownerId string, triggerId string, actionId string, actionParam map[string]interface{}, reactionParam map[string]interface{}) error
	FindWorkflowByNameAndIdUser(name string, ownerId string) (entities.Workflow, error)
	FindWorkflowByIdTrigger(triggerId string) ([]entities.Workflow, error)
}

type Repository struct {
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	UserServiceRepository UserServiceRepository
	ActionRepository      ActionRepository
	TriggerRepository     TriggerRepository
	WorkflowRepository    WorkflowRepository
}
