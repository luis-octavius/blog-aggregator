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
	dbUrl := "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

	// opens connection with database
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	state := types.State{
		Db:     dbQueries,
		Config: &cfg,
	}

	commandsHandler := cli.Commands{
		Commands: map[string]func(*types.State, cli.Command) error{},
	}

	commandsHandler.Register("login", cli.HandlerLogin)
	commandsHandler.Register("register", cli.HandlerRegister)
	commandsHandler.Register("reset", cli.HandlerDelete)

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
