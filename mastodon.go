package main

import (
	"context"
	"github.com/mattn/go-mastodon"
	"log"
)

type MastodonClient struct {
	mastodonClient       MastodonClientInterface
	mastodonBaseUrl      string
	mastodonClientId     string
	mastodonClientSecret string
	mastodonAccessToken  string
}

func (m *MastodonClient) sendStatus(statusContent string) (*mastodon.Status, error) {
	toot := &mastodon.Toot{Status: statusContent}
	status, err := m.mastodonClient.PostStatus(context.Background(), toot)
	if err != nil {
		return nil, err
	}
	log.Println("status posted!")
	return status, nil
}

type MastodonClientInterface interface {
	PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error)
}

type WrappingMastodonClient struct {
	mastodonClient *mastodon.Client
}

func (receiver WrappingMastodonClient) PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error) {
	return receiver.mastodonClient.PostStatus(ctx, toot)
}

func NewWrappingMastodonClient(config *mastodon.Config) *WrappingMastodonClient {
	return &WrappingMastodonClient{mastodonClient: mastodon.NewClient(config)}
}
