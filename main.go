package main

import (
	"fmt"
	"github.com/mattn/go-mastodon"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var publishedStories []Story

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
			_, stories := retrieveTopStoriesAndPostThemToMastodon(hn, mastodonClient, hackerNewsNumberOfStories)
			publishedStories = append(publishedStories, stories...)
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

func retrieveTopStoriesAndPostThemToMastodon(hn *HackerNews, mastodonClient *MastodonClient, numberOfStories int) ([]*mastodon.Status, []Story) {
	postedStatuses := make([]*mastodon.Status, 0, numberOfStories)

	stories, err := hn.GetTopStories(numberOfStories)
	if err != nil {
		log.Println("Could not get the latest stories", err)
		return nil, nil
	}

	deDuplicateStories(&stories)

	for _, story := range stories {
		log.Printf("About to publish to Mastodon this story %v", story)
		status, err := mastodonClient.sendStatus(fmt.Sprintf("%s\nLink: %v\nComments: %v", story.Title, story.URL, story.CommentURL))
		if err != nil {
			log.Println("Could not publish to Mastodon", err)
		}
		postedStatuses = append(postedStatuses, status)
	}
	return postedStatuses, stories
}

func deDuplicateStories(stories *[]Story) {
	storyMap := make(map[Story]bool)
	for _, story := range publishedStories {
		storyMap[story] = true
	}

	var uniqueStories []Story
	for _, story := range *stories {
		if _, value := storyMap[story]; !value {
			uniqueStories = append(uniqueStories, story)
		}
	}
	*stories = uniqueStories
}
