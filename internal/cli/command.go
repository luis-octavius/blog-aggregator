package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/luis-octavius/blog-aggregator/internal/database"
	"github.com/luis-octavius/blog-aggregator/internal/types"
)

// Command represents a CLI command 
type Command struct {
	Name string
	Args []string
}

// HandlerLogin authenticates a user by username and sets them as the current user.
// it validates command-line arguments, checks user existence in the database,
// and updates the configuration with the authenticated user.
// returns an error if username is not provided, user doesn't exist, or config update fails. 
func HandlerLogin(s *types.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: go run . login <username>")
		return fmt.Errorf("username not provided")
	}

	name := cmd.Args[0]

	ctx := context.Background()
	queries := s.Db

	// verify if user exists in database 
	_, err := queries.GetUser(ctx, name)
	if err != nil {
		fmt.Printf("the user %v does not exist\n", name)
		os.Exit(1)
	}

	// update configuration with authenticated user 
	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user %v: %v", name, err)
	}

	fmt.Printf("username %v has been set\n", name)
	return nil
}

// HandleRegister creates a new user in the database and sets them as the current user. 
// if the username already exists, the operation fails and the program exits. 
// returns an error if username is not provided or user creation fails.
func HandlerRegister(s *types.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: go run . register <username>")
		return fmt.Errorf("username not provided")
	}

	name := cmd.Args[0]

	ctx := context.Background()

	queries := s.Db

	// create a new user with generated UUID and current timestamp 
	// if user already exists, it will fail due to unique constraint
	insertedUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		fmt.Printf("the user %v already exists: %v\n", name, err)
		os.Exit(1)
	}

	// update configuration with authenticated user 
	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user %v: %v", name, err)
	}

	fmt.Printf("user %v was created\n", name)

	// logs info about the user created for debugging 
	fmt.Printf("User: %v\nCreatedAt: %v\nUpdated At: %v\nName: %v\n", insertedUser.ID, insertedUser.CreatedAt, insertedUser.UpdatedAt, insertedUser.Name)

	return nil
}

// HandlerDelete remove all user records from the database. 
// this is a destructive operation intended for reset purpose. 
func HandlerDelete(s *types.State, cmd Command) error {
	ctx := context.Background()

	queries := s.Db

	// execute deletion 
	err := queries.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}

	fmt.Println("rows succesfully deleted")
	return nil
}

func HandlerUsers(s *types.State, cmd Command) error {
	ctx := context.Background() 

	queries := s.Db 

	users, err := queries.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting users from database: %w", err)
	}

	currentUser := s.Config.Current_user_name
	for _, user := range users {
		if currentUser == user.Name {
			fmt.Printf(" - %s (current)\n", user.Name)
		} else {
			fmt.Printf(" - %s\n", user.Name)
		}
	}

	return nil
}
