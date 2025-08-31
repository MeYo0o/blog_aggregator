package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MeYo0o/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	var err error
	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: login
		// Args[2] is the username
		loginUsername := cmd.Args[2]
		_, err = s.DB.GetUser(context.Background(), loginUsername)
		if err != nil {
			return errors.New("user doesn't exist in DB")
		} else {
			s.Cfg.CurrentUsername = loginUsername
			SetUser(s.Cfg.CurrentUsername)
		}
	default:
		return errors.New("you need to pass the username only after login")
	}

	fmt.Printf("user %s has been set!\n", s.Cfg.CurrentUsername)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	var user database.User
	var err error

	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: register
		// Args[2] is the username to be stored inside db
		registrationName := cmd.Args[2]
		user, err = s.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      registrationName,
		})
		if err != nil {
			return errors.New("name already exists")
		} else {
			s.Cfg.CurrentUsername = registrationName
			SetUser(registrationName)
		}
	default:
		return errors.New("you need to pass the username only after register")
	}

	fmt.Printf("user %s is stored in DB!\n", s.Cfg.CurrentUsername)
	fmt.Printf("User: %v\n", user)

	return nil
}
