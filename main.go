package main 

import (
	"github.com/luis-octavius/blog-aggregator/internal/config"
	"fmt"
)
func main() {
	cfg := config.Read() 
	cfg.SetUser("luis")
	
	cfg = config.Read()
}
