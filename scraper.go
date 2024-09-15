package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/MishaYanov/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenReq time.Duration,
){
	log.Printf("SCraping on %v gorutines every %s duration", concurrency, timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C{
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched", err)
		return
	}
	log.Print(feed.Url)
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("failed fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Post: ", item.Title, " on feed ", feed.Name)
	}
	log.Printf("Feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Item))
}