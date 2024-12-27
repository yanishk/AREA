package about_service

import (
	"time"

	"backend/src/entities"
	_ "backend/src/handler/about/docs"
	"backend/src/storage"
)

type AboutService struct {
	ServiceRepository  storage.ServiceRepository
	ActionRepository   storage.ActionRepository
	ReactionRepository storage.ReactionRepository
}

func NewAboutService(ServiceRepository storage.ServiceRepository, ActionRepository storage.ActionRepository,
	ReactionRepository storage.ReactionRepository) *AboutService {
	return &AboutService{
		ServiceRepository:  ServiceRepository,
		ActionRepository:   ActionRepository,
		ReactionRepository: ReactionRepository,
	}
}

func (self *AboutService) GetAboutServer(about entities.About) (entities.About, error) {
	about.Server.CurrentTime = time.Now().Unix()

	services, err := self.ServiceRepository.FindServices()
	if err != nil {
		return about, err
	}

	for _, service := range services {
		var aboutService entities.AboutService
		actions, err := self.ActionRepository.FindActionByServiceId(service.Id)
		if err != nil {
			return about, err
		}
		for _, action := range actions {
			var aboutAction entities.AboutAction
			aboutAction.Name = action.Name
			aboutAction.Description = action.Description
			aboutService.Actions = append(aboutService.Actions, aboutAction)
		}

		reactions, err := self.ReactionRepository.FindReactionByServiceId(service.Id)
		if err != nil {
			return about, err
		}
		for _, reaction := range reactions {
			var aboutReaction entities.AboutReaction
			aboutReaction.Name = reaction.Name
			aboutReaction.Description = reaction.Description
			aboutService.Reactions = append(aboutService.Reactions, aboutReaction)
		}

		aboutService.Name = service.Name
		about.Server.Services = append(about.Server.Services, aboutService)
	}
	return about, nil
}
