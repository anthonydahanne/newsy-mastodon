package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestRetrieveTopStoriesAndPostThemToMastodon(t *testing.T) {
	spyingMastodonClient := SpyingMastodonClient{}
	client := &MastodonClient{
		mastodonClient:       &spyingMastodonClient,
		mastodonBaseUrl:      "https://one.mastodon.server",
		mastodonClientId:     "client_id",
		mastodonClientSecret: "client_secret",
		mastodonAccessToken:  "access_token",
	}
	httpClientWithMockTransport := &http.Client{
		Transport: &mockTransport{},
	}
	hn := &HackerNews{httpClientWithMockTransport}

	statuses, stories := retrieveTopStoriesAndPostThemToMastodon(hn, client, 1)

	if len(statuses) != 1 {
		t.Errorf("There should be exactly 1 status, but there were : %d", len(statuses))
	}

	expectedStory := Story{
		Id:         34469378,
		Title:      "Paper map sales are booming",
		URL:        "https://www.wsj.com/articles/why-paper-map-sales-are-booming-11674164824",
		CommentURL: "https://news.ycombinator.com/item?id=34469378",
	}
	if statuses[0].Content != fmt.Sprintf("%v\n%v", expectedStory.Title, expectedStory.URL) || expectedStory != stories[0] {
		t.Errorf("The posted status did not match the story : %v", statuses[0].Content)
	}
}

func TestLookupEnvAndFailIfNotPresent(t *testing.T) {
	os.Setenv("PIF", "pof")
	value := lookupEnvAndFailIfNotPresent("PIF")
	if value != "pof" {
		t.Errorf("Expected value was: %v", value)
	}
}

func TestDeDuplicateStories(t *testing.T) {
	publishedStories = []Story{
		{Id: 1, Title: "Story 1", URL: "URL 1", CommentURL: "Comment URL 1"},
		{Id: 2, Title: "Story 2", URL: "URL 2", CommentURL: "Comment URL 2"},
		{Id: 3, Title: "Story 3", URL: "URL 3", CommentURL: "Comment URL 3"},
	}

	input := []Story{
		{Id: 4, Title: "Story 4", URL: "URL 4", CommentURL: "Comment URL 4"},
		{Id: 2, Title: "Story 2", URL: "URL 2", CommentURL: "Comment URL 2"},
		{Id: 5, Title: "Story 5", URL: "URL 5", CommentURL: "Comment URL 5"},
	}

	expected := []Story{
		{Id: 4, Title: "Story 4", URL: "URL 4", CommentURL: "Comment URL 4"},
		{Id: 5, Title: "Story 5", URL: "URL 5", CommentURL: "Comment URL 5"},
	}

	deDuplicateStories(&input)

	if !reflect.DeepEqual(input, expected) {
		t.Errorf("Expected %v but got %v", expected, input)
	}
}
