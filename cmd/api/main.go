package main

import (
	"log"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/configs"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Db"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// Load configuration
	config, err := config.LoadConfig()

	// check the configuration
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	err = db.Connect(config.Database.URI)
	// check the connection
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	// Create a new gin router
	router := gin.Default()

	// Register routes
	Routes.RegisterUserRoutes(router)
	Routes.RegisterProductRoutes(router)


	// Run the server
	router.Run(":" + config.Server.Port)

	// Print the configuration
	log.Printf("Server running on port %s", config.Server.Port)

}
