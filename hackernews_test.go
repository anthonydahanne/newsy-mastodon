package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestAddCommentURL(t *testing.T) {
	story := &Story{
		Id:    43,
		Title: "A Story",
		URL:   "https://my.url",
	}
	story.addCommentURL()

	if story.CommentURL != fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.Id) {
		t.Errorf("AddCommentURL() problem: comment url not added properly, got %v", story.CommentURL)
	}
}

func TestGetStory(t *testing.T) {
	httpClientWithMockTransport := &http.Client{
		Transport: &mockTransport{},
	}
	// Create a new HackerNews client with a mock transport returning a valid Json
	hn := &HackerNews{httpClientWithMockTransport}

	// Test getting 2 stories
	story, err := hn.GetStory(43)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedStory := Story{
		Id:         34469378,
		Title:      "Paper map sales are booming",
		URL:        "https://www.wsj.com/articles/why-paper-map-sales-are-booming-11674164824",
		CommentURL: "https://news.ycombinator.com/item?id=34469378",
	}
	if story != expectedStory {
		t.Errorf("Expected story to be equals to the JSON story, but they're not")
	}
}

func TestGetTopStories(t *testing.T) {
	httpClientWithMockTransport := &http.Client{
		Transport: &mockTransport{},
	}
	// Create a new HackerNews client with a mock transport returning a valid Json
	hn := &HackerNews{httpClientWithMockTransport}

	// Test getting 2 stories
	stories, err := hn.GetTopStories(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(stories) != 2 {
		t.Errorf("Expected 2 stories, got %d", len(stories))
	}

	// Test getting 5 stories
	stories, err = hn.GetTopStories(5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(stories) != 5 {
		t.Errorf("Expected 3 stories, got %d", len(stories))
	}
}

func TestGetTopStoriesError(t *testing.T) {
	// Create a mock http client with a faulty transport
	httpClient := &http.Client{
		Transport: &mockErrorTransport{},
	}

	hn := &HackerNews{httpClient}

	_, err := hn.GetTopStories(2)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "Server did not reply with HTTP 200 OK, but 500" {
		t.Errorf("Expected specific error message, got %v", err.Error())
	}
}

type mockErrorTransport struct{}

func (*mockErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil
}

type mockTransport struct{}

func (*mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var testJsonFile io.Reader
	var err error
	if strings.Contains(hackerNewsTopStoriesEndpoint, req.URL.Path) {
		filePath := "testdata/topstories.json"
		testJsonFile, err = os.Open(filePath)
		if err != nil {
			log.Fatalf("Impossible to read file at path: %v", filePath)
		}
	} else {
		filePath := "testdata/story.json"
		testJsonFile, err = os.Open(filePath)
		if err != nil {
			log.Fatalf("Impossible to read file at path: %v", filePath)
		}
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(testJsonFile),
	}, nil
}
