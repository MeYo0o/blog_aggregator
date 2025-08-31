package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MeYo0o/blog_aggregator/internal/config"
	"github.com/MeYo0o/blog_aggregator/internal/database"
	st "github.com/MeYo0o/blog_aggregator/internal/state"
	"github.com/google/uuid"
)

func HandlerLogin(s *st.State, cmd Command) error {
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
			config.SetUser(s.Cfg.CurrentUsername)
		}
	default:
		return errors.New("you need to pass the username only after login")
	}

	fmt.Printf("user %s has been set!\n", s.Cfg.CurrentUsername)

	return nil
}

func HandlerRegister(s *st.State, cmd Command) error {
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
			config.SetUser(registrationName)
		}
	default:
		return errors.New("you need to pass the username only after register")
	}

	fmt.Printf("user %s is stored in DB!\n", s.Cfg.CurrentUsername)
	fmt.Printf("User: %v\n", user)

	return nil
}

func HandleResetUsers(s *st.State, cmd Command) error {
	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: reset
		//!!!!!!
		if err := s.DB.ResetUsers(context.Background()); err != nil {
			return fmt.Errorf("reset Users Failed: %w", err)
		}
	default:
		return errors.New("you don't need any arguments, just the reset command will do")
	}

	fmt.Println("All users have been deleted successfully!")

	return nil
}

func HandleUsers(s *st.State, cmd Command) error {
	var users []database.User
	var err error

	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: users
		//!!!!!!
		users, err = s.DB.GetUsers(context.Background())
		if err != nil {
			return errors.New("couldn't retrieve users from DB")
		}

	default:
		return errors.New("you don't need any arguments, just the users command will do")
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
