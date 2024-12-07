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
	FindActionsServices() ([]entities.ServiceUsefulInformations, error)
	FindReactionsServices() ([]entities.ServiceUsefulInformations, error)
}

type UserServiceRepository interface {
	CreateUserService(idUser string, token string, tokenRefresh string, expiryDate string, idService string) error
	FindUserServiceByIdServiceandIdUser(idUser, idService string) (entities.UserService, error)
	UpdateUserServiceByServiceIdAndUserId(idUser string, token string, tokenRefresh string, expiryDate string, idService string) error
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
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	UserServiceRepository UserServiceRepository
	ActionRepository      ActionRepository
	TriggerRepository     TriggerRepository
	WorkflowRepository    WorkflowRepository
}
