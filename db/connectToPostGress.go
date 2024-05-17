package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectToPostGress() (*sql.DB, error) {
	POSTGRESS_URL := os.Getenv("POSTGRESS_URL")
	db, err := sql.Open("postgres", POSTGRESS_URL)
	if err != nil {
		log.Printf("unable to connect to db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Printf("unable to ping db: %v", err)
	}
	fmt.Println("Connected to PostGress")
	return db, err
}


/*

stateless vs statefull backend
stateless: no session data stored on server
statefull: session data stored on server

// inmemort db

*/