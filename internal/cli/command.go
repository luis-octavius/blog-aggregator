package cli 

import (
	"fmt"
	"strings"

	"github.com/luis-octavius/blog-aggregator/internal/types"
)

type Command struct {
	Name	string 
	Args  []string 
}

func HandlerLogin(s *types.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: go run . login <username>")
		return fmt.Errorf("username not provided")
	}

	joinedArgsSlice := strings.Join(cmd.Args, " ")

	config := s.Config
	err := config.SetUser(joinedArgsSlice)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("username %v has been set\n", joinedArgsSlice)
	return nil 
}


