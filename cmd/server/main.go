package main

import (
	"distributed-db-sharding/config"
	"distributed-db-sharding/database/connections"
)


func main() {

	// Configure the environment
	config , err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	// connect to the Postgres database
	_, _, err = Connection.ConnectToPostgres(&config.Database)
	if err != nil {
		panic(err)
	}
}