package main

import (
	"context"
	"github.com/mattn/go-mastodon"
	"testing"
)

func TestMastodonClientSendStatus(t *testing.T) {

	spyingMastodonClient := SpyingMastodonClient{}
	client := &MastodonClient{
		mastodonClient:       &spyingMastodonClient,
		mastodonBaseUrl:      "https://one.mastodon.server",
		mastodonClientId:     "client_id",
		mastodonClientSecret: "client_secret",
		mastodonAccessToken:  "access_token",
	}

	statusContent := "Hello, World!"
	status, err := client.sendStatus(statusContent)

	if err != nil {
		t.Errorf("sendStatus() failed: %v", err)
	}

	if status.Content != statusContent {
		t.Errorf("status content is not as expected: %v", status.Content)
	}
	if status.Language != "en" {
		t.Errorf("status language is not 'en' as expected: %v", status.Language)
	}
}

type SpyingMastodonClient struct {
}

func (receiver SpyingMastodonClient) PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error) {
	return &mastodon.Status{Content: toot.Status, Language: "en"}, nil
}
