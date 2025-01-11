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
		"Activate playback shuffle",
		"Deactivate playback shuffle",
		"Skip track",
		"Skip to previous track",
		"Add track to playback queue",
		"Save a track",
		"Save an album",
		"Save an audiobook",
		"Save an episode",
		"Save a show",
		"Set playback volume",
		"Follow a playlist",
		"Unfollow a playlist",
	}
}

func getSpotifyUrl(workflow entities.Workflow, baseUrl, parameterName, endUrl string) (string, error) {
	parameter, parameterExists := workflow.ReactionParam[parameterName]
	if !parameterExists {
		return "", fmt.Errorf("Incorrect parameter")
	}

	return (baseUrl + parameter.(string) + endUrl), nil
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
	case "Activate playback shuffle":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/shuffle?state=true", "PUT", bearerType, accessToken, nil)
	case "Deactivate playback shuffle":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/shuffle?state=false", "PUT", bearerType, accessToken, nil)
	case "Skip track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", bearerType, accessToken, nil)
	case "Skip to previous track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/previous", "POST", bearerType, accessToken, nil)
	case "Add track to playback queue":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/player/queue?uri=", "uri", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, accessToken, nil)
	case "Save a track":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/tracks?ids=", "id", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Save an album":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/albums?ids=", "id", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Save an audiobook":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/audiobooks?ids=", "id", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Save an episode":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/episodes?ids=", "id", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Save a show":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/shows?ids=", "id", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Set playback volume":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/player/volume?volume_percent=", "volume", "")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Follow a playlist":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
	case "Unfollow a playlist":
		url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers")
		if err != nil {
			return err
		}
		resp, _ = self.ServiceService.ExecuteApiRequest(url, "DELETE", bearerType, accessToken, nil)
	}

	if resp != nil {
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("Unknown reaction")
}
