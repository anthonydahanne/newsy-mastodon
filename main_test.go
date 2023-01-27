package main

import (
	"fmt"
	"net/http"
	"os"
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

	statuses := retrieveTopStoriesAndPostThemToMastodon(hn, client, 1)

	if len(statuses) != 1 {
		t.Errorf("There should be exactly 1 status, but there were : %d", len(statuses))
	}

	expectedStory := Story{
		Id:         34469378,
		Title:      "Paper map sales are booming",
		URL:        "https://www.wsj.com/articles/why-paper-map-sales-are-booming-11674164824",
		CommentURL: "https://news.ycombinator.com/item?id=34469378",
	}
	if statuses[0].Content != fmt.Sprintf("%v\n%v", expectedStory.Title, expectedStory.URL) {
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
