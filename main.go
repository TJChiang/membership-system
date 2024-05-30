package main

import "membership-system/routes"

func main() {
	router := routes.SetupRoutes()
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
