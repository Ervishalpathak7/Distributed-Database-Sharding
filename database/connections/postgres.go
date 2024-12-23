package Connection

import (
	"database/sql"
	"distributed-db-sharding/config"
	"fmt"
	"time"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


// function to create the connection string for the Postgres database
func connectionString(Databaseconfig *config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Databaseconfig.Host, Databaseconfig.Port, Databaseconfig.User, Databaseconfig.Password, Databaseconfig.Database)
}

// function to connect to the Postgres database using the connection string
func ConnectToPostgres(Databaseconfig *config.DatabaseConfig)(conn *gorm.DB ,  db *sql.DB , err error) {

	// connect to the Postgres database
	db , err = sql.Open("postgres", connectionString(Databaseconfig))
	if err != nil {
		fmt.Println("Error connecting to Postgres database", err)
		return nil, nil, err
	}
	defer db.Close()

	// check if the connection is successful
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to Postgres database", err)
		return nil, nil, err
	}

	// setup connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)


	// establish the connection using gorm
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}),
											&gorm.Config{ })
	if err != nil {
		fmt.Println("Error connecting to Postgres database", err)
		return nil, nil, err
	}

	// return the connection
	fmt.Println("Successfully connected to Postgres database")
	return gormDB ,db,nil
}

