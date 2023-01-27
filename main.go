package main

import (
	"github.com/mattn/go-mastodon"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	log.Println("newsy-mastodon application starting up...")

	hackerNewsNumberOfStories, err := strconv.Atoi(os.Getenv("HN_NUMBER_OF_STORIES"))
	if err != nil {
		log.Fatal("The required environment variable HN_NUMBER_OF_STORIES could not be parsed into an int ", err)
	}

	mastodonBaseUrl := lookupEnvAndFailIfNotPresent("MASTODON_BASE_URL")
	mastodonAccessToken := lookupEnvAndFailIfNotPresent("MASTODON_ACCESS_TOKEN")
	mastodonClientSecret := lookupEnvAndFailIfNotPresent("MASTODON_CLIENT_SECRET")
	mastodonClientId := lookupEnvAndFailIfNotPresent("MASTODON_CLIENT_ID")

	hn := &HackerNews{
		httpClient: &http.Client{},
	}

	mastodonClient := &MastodonClient{
		mastodonClient: NewWrappingMastodonClient(&mastodon.Config{
			Server:       mastodonBaseUrl,
			ClientID:     mastodonClientId,
			ClientSecret: mastodonClientSecret,
			AccessToken:  mastodonAccessToken,
		})}

	log.Printf("Publishing %d top HN stories to Mastodon %v every hour, on top of the hour",
		hackerNewsNumberOfStories,
		mastodonBaseUrl)
	ticker := time.Tick(time.Minute)
	for range ticker {
		now := time.Now()
		if now.Minute() == 0 {
			retrieveTopStoriesAndPostThemToMastodon(hn, mastodonClient, hackerNewsNumberOfStories)
		}
	}
}

func lookupEnvAndFailIfNotPresent(key string) string {
	value, envVarSet := os.LookupEnv(key)
	if !envVarSet {
		log.Fatalf("The required environment variable %v was not set", key)
	}
	return value
}

func retrieveTopStoriesAndPostThemToMastodon(hn *HackerNews, mastodonClient *MastodonClient, numberOfStories int) []*mastodon.Status {
	postedStatuses := make([]*mastodon.Status, 0, numberOfStories)

	stories, err := hn.GetTopStories(numberOfStories)
	if err != nil {
		log.Println("Could not get the latest stories", err)
		return nil
	}

	for _, story := range stories {
		log.Printf("About to publish to Mastodon this story %v", story)
		status, err := mastodonClient.sendStatus(story.Title + "\n" + story.URL)
		if err != nil {
			log.Println("Could not publish to Mastodon", err)
		}
		postedStatuses = append(postedStatuses, status)
	}
	return postedStatuses
}
