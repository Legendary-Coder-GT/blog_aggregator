package main

import (
	"fmt"
	"context"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		usr, err := s.db.GetUser(ctx, s.cfg.Current_user_name)
		if err != nil {
			fmt.Print("No user currently logged in")
			return err
		}
		return handler(s, cmd, usr)
	}
}