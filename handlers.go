package main

import (
	"context"
	"fmt"
	"time"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username required\n")
	}
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Print("User has been set\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username required\n")
	}
	ctx := context.Background()
	params := database.CreateUserParams{uuid.New(), time.Now(), time.Now(), cmd.args[0]}
	usr, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Print("User has been created and set\n")
	fmt.Print(usr, "\n")
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		fmt.Print("Deletion was unsuccessfun\n")
		return err
	}
	fmt.Print("Table is cleared\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, usr := range users {
		if usr.Name == s.cfg.Current_user_name {
			fmt.Print("* ", usr.Name, " (current)\n")
		} else {
			fmt.Print("* ", usr.Name, "\n")
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Print("Error fetching feed\n")
		return err
	}
	fmt.Print(*feed, "\n")
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Insufficient inputs, need name and url of feed\n")
	} else if len(cmd.args) == 1 {
		return fmt.Errorf("Insufficient inputs, need url of feed\n")
	}
	ctx := context.Background()
	params := database.CreateFeedParams{
		uuid.New(), 
		time.Now(), 
		time.Now(), 
		cmd.args[0],
		cmd.args[1],
		user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}
	cmd.args = []string{cmd.args[1]}
	err = middlewareLoggedIn(handlerFollow)(s, cmd)
	if err != nil {
		return err
	}
	fmt.Print(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Print("Feed_name\tURL\tUser_name\n-------------------------\n")
	for _, row := range feeds {
		fmt.Print(row.FeedName, "\t", row.Url, "\t", row.UserName, "\n")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Insufficient inputs, need a URL\n")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Feed does not exist\n")
	}
	params := database.CreateFeedFollowParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		user.ID,
		feed.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("Issue creating feed follow record\n")
	}
	fmt.Print("Feed follow created successfully!\n")
	fmt.Print("Feed Name: ", feed.Name, "\n")
	fmt.Print("User Name: ", user.Name, "\n")
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Issue retrieving feed information\n")
	}
	if len(feeds) == 0 {
		fmt.Print("No feeds for user currently\n")
		return nil
	}
	for _, feed := range feeds {
		fmt.Print("* ", feed.FeedName, "\n")
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No URL provided\n")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.Unfollow(ctx, database.UnfollowParams{user.ID, feed.ID})
	if err != nil {
		return err
	}
	fmt.Print("Successfully unfollowed ", feed.Name, "\n")
	return nil
}