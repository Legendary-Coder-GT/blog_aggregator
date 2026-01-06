package main

import (
	"fmt"
	"context"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/database"
	"time"
	"database/sql"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed_to_fetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	sql_time := sql.NullTime{time.Now(), true}
	params := database.MarkFeedFetchedParams{sql_time, feed_to_fetch.ID}
	err = s.db.MarkFeedFetched(ctx, params)
	if err != nil {
		return err
	}
	feed, err := fetchFeed(ctx, feed_to_fetch.Url)
	if err != nil {
		return err
	}
	fmt.Print(feed_to_fetch.Name, "\n")
	for _, item := range feed.Channel.Item {
		parsed_time, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		post_params := database.CreatePostParams{
			uuid.New(),
			time.Now(),
			time.Now(),
			item.Title,
			item.Link,
			sql.NullString{item.Description, true},
			sql.NullTime{parsed_time, true},
			feed_to_fetch.ID,
		}
		_, err = s.db.CreatePost(ctx, post_params)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "posts_url_key"` {
				fmt.Print(err)
			}
		}
	}
	return nil
}