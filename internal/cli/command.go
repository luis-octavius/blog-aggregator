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

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *types.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: go run . login <username>")
		return fmt.Errorf("username not provided")
	}

	name := cmd.Args[0]

	ctx := context.Background()
	queries := s.Db

	_, err := queries.GetUser(ctx, name)
	if err != nil {
		return fmt.Errorf("the user %v does not exist", name)
		os.Exit(1)
	}

	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user %v: %v", name, err)
	}

	fmt.Printf("username %v has been set\n", name)
	return nil
}

func HandlerRegister(s *types.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: go run . register <username>")
		return fmt.Errorf("username not provided")
	}

	name := cmd.Args[0]

	ctx := context.Background()

	queries := s.Db

	insertedUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("the user %v already exists: %v", name, err)
		os.Exit(1)
	}

	s.Config.SetUser(name)
	fmt.Printf("user %v was created\n", name)

	fetchedUser, err := queries.GetUser(ctx, insertedUser.Name)

	fmt.Printf("User: %v\nCreatedAt: %v\nUpdated At: %v\nName: %v\n", fetchedUser.ID, fetchedUser.CreatedAt, fetchedUser.UpdatedAt, fetchedUser.Name)

	return nil
}

func HandlerDelete(s *types.State, cmd Command) error {
	ctx := context.Background()

	queries := s.Db

	err := queries.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}

	fmt.Println("rows succesfully deleted")
	return nil
}
