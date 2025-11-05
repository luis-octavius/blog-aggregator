package main

import (
	"fmt"
	"github.com/luis-octavius/blog-aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("luis")

	cfg = config.Read()
}
