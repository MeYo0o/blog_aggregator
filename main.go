package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/MeYo0o/blog_aggregator/internal/config"
	"github.com/MeYo0o/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}

	dbQueries := database.New(db)

	state := config.State{
		DB:  dbQueries,
		Cfg: &cfg,
	}

	var commands config.Commands
	commands.Cmds = make(map[string]func(*config.State, config.Command) error, 0)

	//* Register Commands & Handlers
	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	commands.Register("reset", config.HandleResetUsers)

	//* Run Commands
	if len(os.Args) == 1 {
		log.Fatal("empty arguments!... need to pass arguments for the program to run")
	}

	if err := commands.Run(&state, config.Command{
		Name: os.Args[1],
		Args: os.Args,
	}); err != nil {
		log.Fatal(err)
	}

}
