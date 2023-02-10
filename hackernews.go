package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const hackerNewsBaseUrl = "https://hacker-news.firebaseio.com/"
const hackerNewsTopStoriesEndpoint = hackerNewsBaseUrl + "v0/topstories.json"
const hackerNewsStoryEndpoint = hackerNewsBaseUrl + "v0/item/%d.json"
const hackerNewsPublicItemEndpoint = "https://news.ycombinator.com/item?id=%d"

type HackerNews struct {
	httpClient *http.Client
}

func (h *HackerNews) GetTopStories(numberOfStories int) ([]Story, error) {
	resp, err := h.httpClient.Get(hackerNewsTopStoriesEndpoint)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server did not reply with HTTP 200 OK, but %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var storyIds []int
	err = json.Unmarshal(body, &storyIds)
	if err != nil {
		return nil, err
	}

	var stories []Story
	for i := 0; i < numberOfStories; i++ {
		story, err := h.GetStory(storyIds[i])
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}
	return stories, nil
}

func (h *HackerNews) GetStory(id int) (Story, error) {
	url := fmt.Sprintf(hackerNewsStoryEndpoint, id)
	resp, err := h.httpClient.Get(url)
	if err != nil {
		return Story{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Story{}, err
	}

	var story Story
	err = json.Unmarshal(body, &story)
	if err != nil {
		return Story{}, err
	}
	story.addCommentURL()
	return story, nil
}

func (s *Story) addCommentURL() {
	s.CommentURL = fmt.Sprintf(hackerNewsPublicItemEndpoint, s.Id)
}
