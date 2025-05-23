package main

import (
	"github.com/369-123/api-students/api"
	"github.com/rs/zerolog/log"
)

func main() {
	server := api.NewServer()
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msgf("failed to start server")
	}
}
