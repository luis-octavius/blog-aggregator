package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/luis-octavius/blog-aggregator/internal/database"
	"github.com/luis-octavius/blog-aggregator/internal/types"
	"github.com/luis-octavius/blog-aggregator/internal/feed"
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

// HandlerRegister creates a new user in the database and sets them as the current user. 
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

// HandlerUsers lists all users from the database and displays their status. 
// it highlights the currently authenticated user with a special marker. 
// returns an error if the database query fails 
func HandlerUsers(s *types.State, cmd Command) error {
	ctx := context.Background() 

	queries := s.Db 

	// retrieve all users from database - fails if query execution errors
	users, err := queries.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting users from database: %w", err)
	}

	// get currently authenticated user from configuration 
	currentUser := s.Config.Current_user_name

	// display users with visual indicator for current user 
	for _, user := range users {
		if currentUser == user.Name {
			fmt.Printf(" - %s (current)\n", user.Name)
		} else {
			fmt.Printf(" - %s\n", user.Name)
		}
	}

	return nil
}

// HandlerAgg print RSSFeed data by calling FetchFeed(ctx context, feedURl string)
// returns an error if FetchFeed fails to retrieve a RSS Struct 
func HandlerAgg(s *types.State, cmd Command) error {
	ctx := context.Background()
	rss, err := feed.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Println("Rss:", rss)

	return nil 
} 

func HandlerAddFeed(s *types.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}

	ctx := context.Background()

	name := cmd.Args[0]
	url := cmd.Args[1]

	queries := s.Db

	queryActualUser, err := queries.GetUser(ctx, user.Name) 
	if err != nil {
		return fmt.Errorf("error getting current user in query GetUser: %w", err)
	}

	insertedFeed, err := queries.CreateFeed(ctx, database.CreateFeedParams{
		Name: name, 
		Url: url, 
		UserID: queryActualUser.ID, 
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error inserting feed in query CreateFeed: %w", err)
	}

	_, err = queries.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: queryActualUser.ID, 
		FeedID: insertedFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed in list of following feeds by user %v: %w", user, err)
	}

	fmt.Println("Feed follows added succesfully")

	fmt.Println("feed recorded succesfully!")
	fmt.Printf("ID: %v\nName: %v\nUrl: %v\nCreated At: %v\nUpdated At: %v\n", insertedFeed.ID, insertedFeed.Url, insertedFeed.UserID, insertedFeed.CreatedAt, insertedFeed.UpdatedAt)

	return nil 
}

// HandlerListFeeds fetchs all feeds and prints all the 
// records one by one showing name, url and the user 
// that owns the feed 
// 
// returns an error if the query GetFeeds fails 
func HandlerListFeeds(s *types.State, cmd Command) error {
	ctx := context.Background()

	queries := s.Db 

	listFeeds, err := queries.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error fetching the list of feeds: %w", err)
	}

	for _, feed := range listFeeds {
		fmt.Println("")
		fmt.Printf("Name: %v\nURL: %v\nUsername: %v\n", feed.Name, feed.Url, feed.Name_2)
	}

	return nil
}

// HandlerFollow creates a feed_follows relationship between the current user and a feed. 
// it validates the feed exists by URL and the user is authenticated, then creates 
// the association in the database. On success, it displays the feed name and username. 
//
// returns error if:
// - feed lookup by URL fails (feed doesn't exist)
// - user retrieval fails (user not authenticated)
// - feed follow creation fails (duplicate violation)
func HandlerFollow(s *types.State, cmd Command, user database.User) error {
	ctx := context.Background()

	url := cmd.Args[0]

	queries := s.Db 
	
	// lookup feed by URL to ensures it exists 
	feed, err := queries.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed by provided url: %w", err)
	}

	// retrieve current logged user 
	_, err = queries.GetUser(ctx, user.Name)
	if err != nil {
		return fmt.Errorf("error getting user by current user name: %w", err)
	}

	// create feed_follows association between user and feed 
	insertFeedFollow, err := queries.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("Feed's name: %v\nCurrent user: %v\n", insertFeedFollow.FeedName, insertFeedFollow.UserName)

	return nil 
}

func HandlerFollowing(s *types.State, cmd Command, user database.User) error {
	ctx := context.Background() 

	queries := s.Db

	feedFollows, err := queries.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return fmt.Errorf("error getting the feed followed by user %v: %w", user.Name, err)
	}

	fmt.Printf("Current user: %v\n", user.Name)
	for _, feed := range feedFollows {
			fmt.Printf("Feed: %v", feed.FeedName)
		}	
	
	return nil 
}
