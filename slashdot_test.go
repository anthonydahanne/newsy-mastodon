package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestGetLatestArticles(t *testing.T) {
	httpClientWithMockTransport := &http.Client{
		Transport: &mockTransportSlashDot{},
	}

	slashdot := &Slashdot{httpClientWithMockTransport}

	stories, err := slashdot.ScrapeSlashdot()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(stories) != 3 {
		t.Errorf("Expected 2 stories, got %d", len(stories))
	}
	if stories[1].Id != 170306193 ||
		stories[1].Title != "Kraken Settles With SEC For $30 Million, Agrees To Shutter Crypto-Staking Operation (coindesk.com)" ||
		stories[1].URL != "https://slashdot.org/story/23/02/09/2127238/kraken-settles-with-sec-for-30-million-agrees-to-shutter-crypto-staking-operation" {
		t.Errorf("Did not get expected story, got Id: %d, %v, %v", stories[1].Id, stories[1].Title, stories[1].URL)
	}

}

func TestGetLatestArticlesError(t *testing.T) {
	// Create a mock http client with a faulty transport
	httpClient := &http.Client{
		Transport: &mockErrorTransport{},
	}

	slashdot := &Slashdot{httpClient}

	_, err := slashdot.ScrapeSlashdot()
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "status code error: 500 " {
		t.Errorf("Expected specific error message, got %v", err.Error())
	}
}

type mockTransportSlashDot struct{}

func (*mockTransportSlashDot) RoundTrip(req *http.Request) (*http.Response, error) {
	var testJsonFile io.Reader
	var err error
	filePath := "testdata/slashdot.html"
	testJsonFile, err = os.Open(filePath)
	if err != nil {
		log.Fatalf("Impossible to read file at path: %v", filePath)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(testJsonFile),
	}, nil
}
