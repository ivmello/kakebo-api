package main

import (
	"github.com/ivmello/kakebo-go-api/internal/adapters/database"
	"github.com/ivmello/kakebo-go-api/internal/adapters/webserver"
	"github.com/ivmello/kakebo-go-api/internal/config"
	"github.com/ivmello/kakebo-go-api/internal/provider"
)

func main() {
	config := config.New()
	database := database.New(config)

	db := database.Connect()
	defer db.Close()

	provider := provider.New(config, db)
	webserver := webserver.New(provider)

	webserver.Start()
}
