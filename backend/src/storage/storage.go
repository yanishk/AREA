package storage

import (
	"backend/src/entities"
)

type UserRepository interface {
	CreateUser(email, password, connectionType string) error
	FindUserByEmail(email, connectionType string) (entities.User, error)
	FindUserById(userId string) (entities.User, error)
	UpdateUser(email, password, connectionType string) error
	DeleteUser(email, connectionType string) error
}

type ServiceRepository interface {
	CreateService(name, color, logo string) error
	FindServiceById(id string) (entities.Service, error)
	FindServiceByName(name string) (entities.Service, error)
	FindAllServices() ([]entities.Service, error)
	FindActionsServices() ([]entities.Service, error)
	FindReactionsServices() ([]entities.Service, error)
}

type UserServiceRepository interface {
	CreateUserService(userId, token, tokenRefresh, expiryDate, serviceId string) error
	FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error)
	UpdateUserServiceByServiceIdAndUserId(userId, accessToken, refreshToken, expiryDate, serviceId string) error
	DeleteUserServiceByUserId(userId string) error
}

type ReactionRepository interface {
	CreateReaction(name, description, serviceId string, nbParam int) error
	FindReactionById(id string) (entities.Reaction, error)
	FindReactionByName(name string) (entities.Reaction, error)
	FindReactionsByServiceId(serviceId string) ([]entities.Reaction, error)
}

type ActionRepository interface {
	CreateAction(name, description, serviceId string, nbParam int) error
	FindActionById(id string) (entities.Action, error)
	FindActionByName(name string) (entities.Action, error)
	FindActionsByServiceId(serviceId string) ([]entities.Action, error)
	FindActionByNameAndServiceId(name, serviceId string) (entities.Action, error)
}

type WorkflowRepository interface {
	CreateWorkflow(name, ownerId, actionId, reactionId string, actionParam, reactionParam, actionData map[string]interface{}) error
	FindWorkflowById(id string) (entities.Workflow, error)
	FindWorkflowsByActionId(actionId string) ([]entities.Workflow, error)
	FindWorkflowsByOwnerId(ownerId string) ([]entities.Workflow, error)
	UpdateWorkflow(id string, updatedWorkflow entities.Workflow) error
	DeleteWorkflow(id, ownerId string) error
	DeleteWorkflowByOwnerId(ownerId string) error
}

type Repository struct {
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	UserServiceRepository UserServiceRepository
	ReactionRepository    ReactionRepository
	ActionRepository      ActionRepository
	WorkflowRepository    WorkflowRepository
}
