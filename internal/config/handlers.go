package config

import (
	"errors"
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: login
		// Args[2] is the username
		s.Cfg.CurrentUsername = cmd.Args[2]
		SetUser(s.Cfg.CurrentUsername)
	default:
		return errors.New("you need to pass the username only after login")
	}

	fmt.Printf("user %s has been set!", s.Cfg.CurrentUsername)

	return nil
}
