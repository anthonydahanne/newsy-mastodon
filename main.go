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

var publishedHNStories []Story
var publishedSlashdotArticles []Story

func main() {

	log.Println("newsy-mastodon application starting up...")

	hackerNewsNumberOfStories, err := strconv.Atoi(os.Getenv("HN_NUMBER_OF_STORIES"))
	if err != nil {
		log.Fatal("The required environment variable HN_NUMBER_OF_STORIES could not be parsed into an int ", err)
	}

	hnMastodonBaseUrl := lookupEnvAndFailIfNotPresent("HN_MASTODON_BASE_URL")
	hnMastodonClient := &MastodonClient{
		mastodonClient: NewWrappingMastodonClient(&mastodon.Config{
			Server:       hnMastodonBaseUrl,
			ClientID:     lookupEnvAndFailIfNotPresent("HN_MASTODON_CLIENT_ID"),
			ClientSecret: lookupEnvAndFailIfNotPresent("HN_MASTODON_CLIENT_SECRET"),
			AccessToken:  lookupEnvAndFailIfNotPresent("HN_MASTODON_ACCESS_TOKEN"),
		}),
	}
	slashdotMastodonBaseUrl := lookupEnvAndFailIfNotPresent("SLASHDOT_MASTODON_BASE_URL")
	slashdotMastodonClient := &MastodonClient{
		mastodonClient: NewWrappingMastodonClient(&mastodon.Config{
			Server:       slashdotMastodonBaseUrl,
			ClientID:     lookupEnvAndFailIfNotPresent("SLASHDOT_MASTODON_CLIENT_ID"),
			ClientSecret: lookupEnvAndFailIfNotPresent("SLASHDOT_MASTODON_CLIENT_SECRET"),
			AccessToken:  lookupEnvAndFailIfNotPresent("SLASHDOT_MASTODON_ACCESS_TOKEN"),
		}),
	}

	hn := &HackerNews{
		httpClient: &http.Client{},
	}

	slashdot := &Slashdot{
		httpClient: &http.Client{},
	}

	log.Printf("Publishing %d top HN stories to Mastodon %v every hour, on top of the hour",
		hackerNewsNumberOfStories,
		hnMastodonBaseUrl)
	log.Printf("Publishing latest Slashdot articles to Mastodon %v every hour, past 30",
		slashdotMastodonBaseUrl)
	ticker := time.Tick(time.Minute)
	for range ticker {
		now := time.Now()
		if now.Minute() == 0 {
			_, stories := retrieveTopStoriesAndPostThemToMastodon(hn, hnMastodonClient, hackerNewsNumberOfStories)
			publishedHNStories = append(publishedHNStories, stories...)
		} else if now.Minute() == 30 {
			_, articles := retrieveLatestArticlesAndPostThemToMastodon(slashdot, slashdotMastodonClient)
			publishedSlashdotArticles = append(publishedHNStories, articles...)

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
	postedStatuses := make([]*mastodon.Status, 0)

	stories, err := hn.GetTopStories(numberOfStories)
	if err != nil {
		log.Println("Could not get the latest stories", err)
		return nil, nil
	}

	return deDuplicateAndPostStoriesToMastodon(stories, publishedHNStories, mastodonClient, postedStatuses)
}

func retrieveLatestArticlesAndPostThemToMastodon(slashdot *Slashdot, mastodonClient *MastodonClient) ([]*mastodon.Status, []Story) {
	postedStatuses := make([]*mastodon.Status, 0)

	stories, err := slashdot.ScrapeSlashdot()
	if err != nil {
		log.Printf("Error scraping Slashdot: %s", err)
		return nil, nil
	}

	return deDuplicateAndPostStoriesToMastodon(stories, publishedSlashdotArticles, mastodonClient, postedStatuses)
}

func deDuplicateAndPostStoriesToMastodon(stories []Story, publishedStories []Story, mastodonClient *MastodonClient, postedStatuses []*mastodon.Status) ([]*mastodon.Status, []Story) {
	deDuplicateStories(&stories, publishedStories)

	for _, story := range stories {
		log.Printf("About to publish to Mastodon this story %v", story)
		statusToPost := fmt.Sprintf("%s\nLink: %v", story.Title, story.URL)
		if &story.CommentURL != nil && story.CommentURL != "" {
			statusToPost = statusToPost + fmt.Sprintf("\nComments: %v", story.CommentURL)
		}
		status, err := mastodonClient.sendStatus(statusToPost)
		if err != nil {
			log.Println("Could not publish to Mastodon", err)
		}
		postedStatuses = append(postedStatuses, status)
	}
	return postedStatuses, stories
}

func deDuplicateStories(stories *[]Story, publishedStories []Story) {
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

type Story struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	CommentURL string
}
