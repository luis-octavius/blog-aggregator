package cli 

import (
	"context"	

	"github.com/luis-octavius/blog-aggregator/internal/database"
	"github.com/luis-octavius/blog-aggregator/internal/types"
)

// MiddlewareLoggedIn wraps command handlers that require an authenticated user 
// it validates if current user exists in database before executing the handler 
// 
// protected commands: 
// - follow 
// - following 
// - addfeed 
// 
// returns a new handler function with user authentication pre-validated
func MiddlewareLoggedIn(handler func(s *types.State, cmd Command, user database.User) error) func(*types.State, Command) error {	
	return func(s *types.State, cmd Command) error {
		username := s.Config.Current_user_name

		// fetch user from database to validate authentication 
		fetchedUser, err := s.Db.GetUser(context.Background(), username)
		if err != nil {
			return err
		}
		
		// execute the original handler with authenticated user
		return handler(s, cmd, fetchedUser)
	}
}
