package workflow_service

import (
	"fmt"
	"net/http"

	"backend/src/entities"
)

func spotifyReactions() []string {
	return []string{
		"Start playback",
		"Pause playback",
		"Skip track",
		"Add track to playback queue",
		"Save a track",
	}
}

func (self *WorkflowService) checkSpotifyReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Spotify", reactionFound.Name, spotifyReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	var resp *http.Response
	switch reactionFound.Name {
	case "Start playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/play", "PUT", bearerType, accessToken, nil)
	case "Pause playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", bearerType, accessToken, nil)
	case "Skip track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", bearerType, accessToken, nil)
	case "Add track to playback queue":
		uri, uriExists := workflow.ReactionParam["link"]
		if !uriExists {
			return fmt.Errorf("Incorrect URI")
		}

		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlQueue, "POST", bearerType, accessToken, nil)
	case "Save a track":
		id, idExists := workflow.ReactionParam["id"]
		if !idExists {
			return fmt.Errorf("Incorrect ID")
		}

		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		resp, _ = self.ServiceService.ExecuteApiRequest(urlSave, "PUT", bearerType, accessToken, nil)
	}

	if resp != nil {
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("Unknown reaction")
}
