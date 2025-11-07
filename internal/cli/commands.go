package cli

import (
	"fmt"

	"github.com/luis-octavius/blog-aggregator/internal/types"
)

// Commands represents a registry of available CLI commands. 
// It maps command names to their corresponding handler functions. 
type Commands struct {
	Commands map[string]func(*types.State, Command) error
}

// Run executes a command by looking up its name in the registry.
// returns an error if the command doesn't exist or the handler function fails. 
func (c *Commands) Run(s *types.State, cmd Command) error {
	// lookup command in registry, return error if not registered 
	command, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}

	// execute the command handler with provided state and arguments  
	err := command(s, cmd)
	if err != nil {
		return fmt.Errorf("command %s failed: %w", cmd.Name, err)
	}
	return nil
}

// Register adds a new command to the registry. 
// This allows dynamic registration of commmand handlers at runtime.
func (c *Commands) Register(name string, f func(*types.State, Command) error) error {
	c.Commands[name] = f
	return nil
}
