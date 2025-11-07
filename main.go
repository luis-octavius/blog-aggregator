package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/luis-octavius/blog-aggregator/internal/cli"
	"github.com/luis-octavius/blog-aggregator/internal/config"
	"github.com/luis-octavius/blog-aggregator/internal/database"
	"github.com/luis-octavius/blog-aggregator/internal/types"
)

func main() {
	// load application configuration from file 	
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	dbUrl := cfg.Db_url

	// establish connection with database
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	dbQueries := database.New(db)

	// initialize application state with dependencies 
	state := types.State{
		Db:     dbQueries,
		Config: &cfg,
	}

	// CLI command registry - maps command names to handler functions 
	commandsHandler := cli.Commands{
		Commands: map[string]func(*types.State, cli.Command) error{},
	}

	// register available commands 
	commandsHandler.Register("login", cli.HandlerLogin)
	commandsHandler.Register("register", cli.HandlerRegister)
	commandsHandler.Register("reset", cli.HandlerDelete)
	commandsHandler.Register("users", cli.HandlerUsers)

	args := os.Args

	fmt.Println(args)

	// validate command-line arguments 
	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	// parse command from command-line arguments 
	cmd := cli.Command{
		Name: args[1],
		Args: args[2:],
	}

	// execute requested command 
	err = commandsHandler.Run(&state, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
