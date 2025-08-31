package commands

import (
	"fmt"

	st "github.com/MeYo0o/blog_aggregator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmds map[string]func(*st.State, Command) error
}

func (c *Commands) Run(s *st.State, cmd Command) error {
	if command, ok := c.Cmds[cmd.Name]; ok {
		if err := command(s, cmd); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("command: %s not found", cmd.Name)
	}

	return nil
}

func (c *Commands) Register(name string, f func(*st.State, Command) error) {
	c.Cmds[name] = f
}
