package main

import (
	"github.com/joho/godotenv"
	"membership-system/routes"
)

func main() {
	if godotenv.Load() != nil {
		panic("Error loading .env file")
	}

	router := routes.SetupRoutes()
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
