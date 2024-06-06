package main

import (
	"github.com/joho/godotenv"
	"membership-system/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	router := routes.SetupRoutes()
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
