package storage

import (
	"backend/src/entities"
)

type UserRepository interface {
	CreateUser(email string, password string) error
	FindUserByEmail(email string) (entities.User, error)
	UpdateUser(email string, password string) error
	DeleteUser(email string) error
	FindUserById(userId string) (entities.User, error)
}

type ServiceRepository interface {
	CreateService(name string, color string, logo string) error
	FindServiceByName(name string) (entities.Service, error)
	FindActionsServices() ([]entities.ServiceUsefulInformations, error)
	FindReactionsServices() ([]entities.ServiceUsefulInformations, error)
}

type UserServiceRepository interface {
	CreateUserService(userId string, token string, tokenRefresh string, expiryDate string, serviceId string) error
	FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error)
	UpdateUserServiceByServiceIdAndUserId(userId string, accessToken string, refreshToken string, expiryDate string, serviceId string) error
}

type ReactionRepository interface {
	CreateReaction(name string, description string, serviceId string, nbParam int) error
	FindReactionByName(name string) (entities.Reaction, error)
	FindReactionByServiceId(serviceId string) ([]entities.Reaction, error)
	FindReactionByID(id string) (entities.Reaction, error)
}

type ActionRepository interface {
	CreateAction(name string, description string, serviceId string, nbParam int) error
	FindActionByName(name string) (entities.Action, error)
	FindActionByServiceId(serviceId string) ([]entities.Action, error)
}

type WorkflowRepository interface {
	CreateWorkflow(name string, ownerId string, actionId string, reactionId string, actionParam map[string]interface{}, reactionParam map[string]interface{}) error
	FindWorkflowByNameAndUserId(name string, ownerId string) (entities.Workflow, error)
	FindWorkflowByActionId(actionId string) ([]entities.Workflow, error)
	FindWorkflowByUserId(userId string) ([]entities.Workflow, error)
}

type Repository struct {
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	UserServiceRepository UserServiceRepository
	ReactionRepository    ReactionRepository
	ActionRepository      ActionRepository
	WorkflowRepository    WorkflowRepository
}
