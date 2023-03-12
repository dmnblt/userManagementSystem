package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "go-crud/controllers"
	"log"
)

var db *sql.DB

// Connect to the database
func Connect() {
	var err error

	db, err = sql.Open("mysql", "admin:nimda@tcp(127.0.0.1:3306)/go_midterm")
	if err != nil {
		log.Fatalf("Error connecting to database: %v \n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error checking the database connection: %v \n", err)
	}

	fmt.Println("Connected to database!")

}

// GetDB returns the database object
func GetDB() *sql.DB {
	return db
}
