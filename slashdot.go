package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Slashdot struct {
	httpClient *http.Client
}

func (s *Slashdot) ScrapeSlashdot() ([]Story, error) {
	res, err := s.httpClient.Get("https://slashdot.org/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var stories []Story
	doc.Find("article h2.story").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".story-title a").Text())
		idAsRawString, existsId := s.Find(".story-title").Attr("id")
		id, err := strconv.Atoi(strings.Split(idAsRawString, "-")[1])
		link, existsLink := s.Find("a").Attr("href")
		if existsId && existsLink && err == nil {
			stories = append(stories, Story{Id: id, Title: title, URL: "https:" + link})
		} else {
			log.Printf("Could not fetch the story with Id: %d, Title: %v, URL: %v", id, title, link)
		}
	})

	return stories, nil
}
