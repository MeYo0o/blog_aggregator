package main

import (
	"fmt"

	"github.com/MeYo0o/blog_aggregator/internal/config"
)

func main() {
	fmt.Println(config.Read())
	config.SetUser("meyo")
	fmt.Println(config.Read())
}
