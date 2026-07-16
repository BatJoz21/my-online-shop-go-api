package main

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/BatJoz21/my-online-shop-go-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found, using environment variables")
	}

	database.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
