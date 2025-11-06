package main

import (
	"fmt"
	"os"

	"github.com/luis-octavius/blog-aggregator/internal/config"
	"github.com/luis-octavius/blog-aggregator/internal/types"
	"github.com/luis-octavius/blog-aggregator/internal/cli"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Errorf("error creating config: %v", err)
	}

	state := types.State{
		Config: &cfg,
	} 	

	commandsHandler := cli.Commands{
		Commands: map[string]func(*types.State, cli.Command) error{},
	}

	commandsHandler.Register("login", cli.HandlerLogin)
	
	args := os.Args 

	fmt.Println(args)

	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = commandsHandler.Run(&state, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
