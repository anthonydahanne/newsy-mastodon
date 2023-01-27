# Newsy Mastodon

A Mastodon Bot that toots the top 5 stories from HackerNews, at the top of every hour.

## How to build

```go build .```

## How to run

You need to export the following environment variables:

* `HN_NUMBER_OF_STORIES`, for example 3 - the number of top stories you want to publish to Mastodon
* `MASTODON_BASE_URL` , for example https://one.mastodon.server
* `MASTODON_CLIENT_ID`, you'll get it after creating a new App in your mastodon instance, associated to a mastodon account
* `MASTODON_CLIENT_SECRET`, you'll get it after creating a new App in your mastodon instance, associated to a mastodon account
* `MASTODON_ACCESS_TOKEN`, you'll get it after creating a new App in your mastodon instance, associated to a mastodon account
