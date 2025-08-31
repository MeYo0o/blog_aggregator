package config

import "github.com/MeYo0o/blog_aggregator/internal/database"

type State struct {
	DB  *database.Queries
	Cfg *Config
}
