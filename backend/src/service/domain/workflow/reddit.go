package workflow_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"backend/src/entities"
)

const redditUserRoute = "https://oauth.reddit.com/user/"

func redditReactions() []string {
	return []string{
		"Submit a comment on a post ",
		"Upvote a post",
		"Downvote a post",
		"Submit a post",
		"Submit a post with a link",
	}
}

func (self *WorkflowService) postRedditComment(accessToken string, workflow entities.Workflow) error {
	comment, commentExists := workflow.ReactionParam["comment"]
	id, idExists := workflow.ReactionParam["id"]
	if !commentExists || !idExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("thing_id", "t3_"+id.(string))
	body.Set("text", comment.(string))
	url := "https://oauth.reddit.com/api/comment"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) upvoteDowvoteRedditPost(postId, accessToken string, vote int) error {
	body := url.Values{}
	body.Set("id", "t3_"+postId)
	body.Set("dir", strconv.Itoa(vote))

	url := "https://oauth.reddit.com/api/vote"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) upvoteRedditPost(accessToken string, workflow entities.Workflow) error {
	id, idExists := workflow.ReactionParam["id"]
	if !idExists {
		return fmt.Errorf(errorMissingField)
	}

	return self.upvoteDowvoteRedditPost(id.(string), accessToken, 1)
}

func (self *WorkflowService) downvoteRedditPost(accessToken string, workflow entities.Workflow) error {
	id, idExists := workflow.ReactionParam["id"]
	if !idExists {
		return fmt.Errorf(errorMissingField)
	}

	return self.upvoteDowvoteRedditPost(id.(string), accessToken, -1)
}

func (self *WorkflowService) submitRedditPost(accessToken string, workflow entities.Workflow) error {
	title, titleExists := workflow.ReactionParam["title"]
	content, contentExists := workflow.ReactionParam["content"]
	subreddit, subredditExists := workflow.ReactionParam["subreddit"]
	if !titleExists || !contentExists || !subredditExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("title", title.(string))
	body.Set("text", content.(string))
	body.Set("sr", subreddit.(string))
	body.Set("kind", "self")
	url := "https://oauth.reddit.com/api/submit"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) submitRedditPostLink(accessToken string, workflow entities.Workflow) error {
	title, titleExists := workflow.ReactionParam["title"]
	link, linkExists := workflow.ReactionParam["link"]
	subreddit, subredditExists := workflow.ReactionParam["subreddit"]
	if !titleExists || !linkExists || !subredditExists {
		return fmt.Errorf(errorMissingField)
	}

	body := url.Values{}
	body.Set("title", title.(string))
	body.Set("url", link.(string))
	body.Set("sr", subreddit.(string))
	body.Set("kind", "link")
	url := "https://oauth.reddit.com/api/submit"

	resp, err := self.executeRedditRequest("POST", url, accessToken, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (self *WorkflowService) checkRedditReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Reddit", reactionFound.Name, redditReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	switch reactionFound.Name {
	case "Submit a comment on a post ":
		return self.postRedditComment(accessToken, workflow)
	case "Downvote a post":
		return self.downvoteRedditPost(accessToken, workflow)
	case "Upvote a post":
		return self.upvoteRedditPost(accessToken, workflow)
	case "Submit a post":
		return self.submitRedditPost(accessToken, workflow)
	case "Submit a post with a link":
		return self.submitRedditPostLink(accessToken, workflow)
	}
	return nil
}

func (self *WorkflowService) executeRedditRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", bearerType+accessToken)
	req.Header.Add("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	resp, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (self *WorkflowService) getRedditPost(url, accessToken string) (entities.RedditPostResponse, error) {
	var result entities.RedditPostResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve post")
	}
	return result, nil
}

func (self *WorkflowService) getRedditComments(url, accessToken string) (entities.RedditCommentResponse, error) {
	var result entities.RedditCommentResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve comments")
	}
	return result, nil
}

func (self *WorkflowService) getRedditVotes(url, accessToken string) (entities.RedditVoteResponse, error) {
	var result entities.RedditVoteResponse

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}

	if len(result.Data.Children) == 0 {
		return result, fmt.Errorf("Could not retrieve votes")
	}
	return result, nil
}

func (self *WorkflowService) getUsernameReddit(accessToken string, workflow entities.Workflow) (entities.RedditUsername, error) {
	var result entities.RedditUsername
	url := "https://oauth.reddit.com/api/v1/me"

	resp, err := self.executeRedditRequest("GET", url, accessToken, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&result)
	if errDecode != nil {
		return result, errDecode
	}
	return result, nil
}

func (self *WorkflowService) isANewPost(workflow entities.Workflow, result entities.RedditPostResponse) error {
	latestPost := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestPost.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestPost.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestPost.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) isANewComment(workflow entities.Workflow, result entities.RedditCommentResponse) error {
	latestComment := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestComment.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestComment.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestComment.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) isANewVote(workflow entities.Workflow, result entities.RedditVoteResponse) error {
	latestVote := result.Data.Children[0].Data
	idLatest, idLatestExist := workflow.ActionData["id"]

	if !idLatestExist || idLatest != latestVote.ID {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["id"] = latestVote.ID

		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return fmt.Errorf(errorUpdatingWorkflow)
		}
	}

	if idLatestExist && idLatest != latestVote.ID {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkRedditNewPostInSubredditAction(accessToken string, workflow entities.Workflow) error {
	subreddit, subredditExists := workflow.ActionParam["subreddit"]
	if !subredditExists {
		return fmt.Errorf(errorMissingField)
	}

	url := "https://oauth.reddit.com/" + subreddit.(string) + "/new.json?limit=1"
	result, err := self.getRedditPost(url, accessToken)
	if err != nil {
		return err
	}

	if len(result.Data.Children) == 0 {
		return fmt.Errorf("No posts found in subreddit")
	}

	return self.isANewPost(workflow, result)
}

func (self *WorkflowService) checkRedditNewCommentByMeAction(accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	url := redditUserRoute + resultUsername.Name + "/comments?limit=1"
	result, err := self.getRedditComments(url, accessToken)
	if err != nil {
		return err
	}

	return self.isANewComment(workflow, result)
}

func (self *WorkflowService) checkRedditVoteByMeAction(voteType, accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	urlPostMe := redditUserRoute + resultUsername.Name + voteType + "?limit=1"
	resultPost, err := self.getRedditVotes(urlPostMe, accessToken)
	if err != nil {
		return err
	}

	return self.isANewVote(workflow, resultPost)
}

func (self *WorkflowService) checkRedditPostAction(postType, accessToken string, workflow entities.Workflow) error {
	resultUsername, err := self.getUsernameReddit(accessToken, workflow)
	if err != nil {
		return err
	}

	url := redditUserRoute + resultUsername.Name + postType + "?limit=1"
	result, err := self.getRedditPost(url, accessToken)
	if err != nil {
		return err
	}

	if len(result.Data.Children) == 0 {
		return fmt.Errorf("No posts found")
	}

	return self.isANewPost(workflow, result)
}

func (self *WorkflowService) checkWorkflowsWithRedditActions(action entities.Action) error {
	allWorkflowsFound, errWorkflow := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range allWorkflowsFound {
		if !workflow.IsActivated {
			continue
		}

		accessToken, err := self.getAccessToken("Reddit", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "Any new post in subreddit":
			self.checkRedditNewPostInSubredditAction(accessToken, workflow)
		case "New post by you":
			self.checkRedditPostAction("/submitted", accessToken, workflow)
		case "New comment by you":
			self.checkRedditNewCommentByMeAction(accessToken, workflow)
		case "New downvoted post by you":
			self.checkRedditVoteByMeAction("/downvoted", accessToken, workflow)
		case "New upvoted post by you":
			self.checkRedditVoteByMeAction("/upvoted", accessToken, workflow)
		case "New post saved by you":
			self.checkRedditPostAction("/saved", accessToken, workflow)
		}
	}
	return nil
}

func (self *WorkflowService) CheckRedditActions() error {
	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Reddit")
	if errFindingService != nil {
		return errFindingService
	}

	allActionsFound, errFindingActions := self.ActionRepository.FindActionsByServiceId(serviceFound.Id)
	if errFindingActions != nil {
		return errFindingActions
	}

	for _, action := range allActionsFound {
		errCheckName := self.checkWorkflowsWithRedditActions(action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}
