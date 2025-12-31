package main

import (
	"os"
	"fmt"
	"errors"
	"log"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	m map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.m[cmd.name]
	if !ok {
		return errors.New("Command not found\n")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.m[name] = f
}

func main() {
	cfg := config.Config{}
	cfg, err := cfg.Read()
	if err != nil {
		log.Fatal("Error reading JSON file: ", err, "\n")
		return
	}
	s := state{}
	s.cfg = &cfg
	c := commands{make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments\n")
		return
	}
	cmd := command{args[1], args[2:]}
	err = c.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Username required\n")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Print("User has been set\n")
	return nil
}