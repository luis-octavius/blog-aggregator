package types

import (
	"github.com/luis-octavius/blog-aggregator/internal/config"
	"github.com/luis-octavius/blog-aggregator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}
