package cli

import (
	"fmt"

	"github.com/luis-octavius/blog-aggregator/internal/types"
)

type Commands struct {
	Commands map[string]func(*types.State, Command) error
}

func (c *Commands) Run(s *types.State, cmd Command) error {
	command, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	err := command(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*types.State, Command) error) error {
	c.Commands[name] = f
	return nil
}
