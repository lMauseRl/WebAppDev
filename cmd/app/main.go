package main

import (
	"log"

	"awesomeProject/internal/api"
)

func main() {
	log.Println("Application start")
	api.StartServer()
	log.Println("Application terminated")
}
