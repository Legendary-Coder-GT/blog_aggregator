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
		fmt.Print("Error fetching feed")
		return err
	}
	fmt.Print(*feed, "\n")
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Insufficient inputs, need name and url of feed")
	} else if len(cmd.args) == 1 {
		return fmt.Errorf("Insufficient inputs, need url of feed")
	}
	ctx := context.Background()
	usr, _ := s.db.GetUser(ctx, s.cfg.Current_user_name)
	params := database.CreateFeedParams{
		uuid.New(), 
		time.Now(), 
		time.Now(), 
		cmd.args[0],
		cmd.args[1],
		usr.ID,
	}
	feed, err := s.db.CreateFeed(ctx, params)
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