package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	cmds "github.com/MeYo0o/blog_aggregator/internal/commands"
	"github.com/MeYo0o/blog_aggregator/internal/config"
	"github.com/MeYo0o/blog_aggregator/internal/database"
	st "github.com/MeYo0o/blog_aggregator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}

	dbQueries := database.New(db)

	state := st.State{
		DB:  dbQueries,
		Cfg: &cfg,
	}

	var commands cmds.Commands
	commands.Cmds = make(map[string]func(*st.State, cmds.Command) error, 0)

	//* Register Commands & Handlers
	commands.Register("login", cmds.HandlerLogin)
	commands.Register("register", cmds.HandlerRegister)
	commands.Register("reset", cmds.HandleResetUsers)
	commands.Register("users", cmds.HandleGetUsers)
	commands.Register("agg", cmds.HandleAgg)
	commands.Register("feeds", cmds.HandleGetFeeds)
	commands.Register("addfeed", middlewareLoggedIn(cmds.HandleAddFeed))
	commands.Register("follow", middlewareLoggedIn(cmds.HandleFollowFeed))
	commands.Register("following", middlewareLoggedIn(cmds.HandleFollowing))
	commands.Register("unfollow", middlewareLoggedIn(cmds.HandleUnfollow))
	commands.Register("browse", middlewareLoggedIn(cmds.HandleBrowse))

	//* Run Commands
	if len(os.Args) == 1 {
		log.Fatal("empty arguments!... need to pass arguments for the program to run")
	}

	if err := commands.Run(&state, cmds.Command{
		Name: os.Args[1],
		Args: os.Args,
	}); err != nil {
		log.Fatal(err)
	}

}

// * Middlewares
func middlewareLoggedIn(handler func(s *st.State, cmd cmds.Command, user database.User) error) func(*st.State, cmds.Command) error {

	return func(s *st.State, c cmds.Command) error {

		currentUser, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("couldn't retrieve user data: %w", err)
		}

		return handler(s, c, currentUser)
	}

}
