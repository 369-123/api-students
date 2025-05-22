package main

import (
	"log"
	"github.com/369-123/api-students/api"
)

func main() {
	server := api.NewServer()
	server.ConfigureRoutes()
	
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}