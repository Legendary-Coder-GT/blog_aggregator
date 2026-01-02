package main

import _ "github.com/lib/pq"

import (
	"os"
	"log"
	"database/sql"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/config"
	"github.com/Legendary-Coder-GT/blog_aggregator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	cfg := config.Config{}
	cfg, err := cfg.Read()
	if err != nil {
		log.Fatal("Error reading JSON file: ", err, "\n")
		return
	}
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Fatal("Error opening SQL Database: ", err, "\n")
		return
	}
	dbQueries := database.New(db)
	s := state{dbQueries, &cfg}
	c := commands{make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
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