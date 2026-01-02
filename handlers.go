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