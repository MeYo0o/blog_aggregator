package main

import (
	"log"
	"os"

	"github.com/MeYo0o/blog_aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	state := config.State{
		Cfg: &cfg,
	}

	var commands config.Commands
	commands.Cmds = make(map[string]func(*config.State, config.Command) error, 0)

	//* Register Commands & Handlers
	commands.Register("login", config.HandlerLogin)

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
